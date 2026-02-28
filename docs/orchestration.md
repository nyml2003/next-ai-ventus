# Ventus 编排系统设计文档

本文档定义 Ventus 的编排系统，负责页面布局结构的声明式配置。

> **与 BFF 架构的关系**：编排系统专注于「布局结构」，BFF 负责「数据获取」。两者配合实现页面的声明式渲染。

---

## 1. 核心思想

### 1.1 职责分离

| 层级 | 职责 | 说明 |
|------|------|------|
| **BFF** | 数据获取 | 首屏数据通过 `/api/page` 统一获取，存入 PageStore |
| **编排系统** | 布局结构 | 声明式配置页面模块的组织和排列方式 |
| **模块** | 数据消费 | 从 PageStore 读取数据，负责自身渲染 |

### 1.2 设计原则

- **页面级隔离**：每个 MPA 页面拥有独立的编排配置，互不影响
- **纯数据配置**：`orchestration.ts` 只包含纯数据，可 JSON 序列化，方便后端下发
- **布局与数据解耦**：布局配置不感知数据来源，只声明模块名称
- **BFF 首屏优先**：所有首屏可见模块必须通过 BFF 获取数据（与 [principles.md](./principles.md) 一致）

---

## 2. 四层结构

```
Page (页面)
  └── Region (区域) - 页面垂直分区
        └── Block (区块) - Flex 布局容器（可嵌套）
              └── Module (模块) - 业务组件（叶子节点）
```

### 2.1 Module（模块）

业务组件，对应 BFF 返回的数据模块：

```typescript
interface ModuleConfig {
  type: 'module';
  name: string;  // 模块名称，对应 BFF 模块名（如 'postList'）
}
```

### 2.2 Block（区块）

仅支持 Flex 布局，使用 CSS 变量控制间距：

```typescript
interface BlockConfig {
  type: 'block';
  flexDirection: 'row' | 'column' | 'row-reverse' | 'column-reverse';
  gap?: string;           // CSS 变量，如 'var(--spacing-md)'
  padding?: string;       // CSS 变量，如 'var(--spacing-lg)'
  margin?: string;        // CSS 变量
  justifyContent?: 'start' | 'center' | 'end' | 'between' | 'around';
  alignItems?: 'start' | 'center' | 'end' | 'stretch';
  style?: Record<string, string>;  // 扩展样式（如固定宽度）
  children: (BlockConfig | ModuleConfig)[];
}
```

### 2.3 Region（区域）

页面的垂直分区，语义化标识：

```typescript
interface RegionConfig {
  id: string;
  type: 'header' | 'content' | 'footer' | 'sidebar';
  padding?: string;   // CSS 变量
  margin?: string;    // CSS 变量
  block: BlockConfig; // 唯一根区块
}
```

### 2.4 Page（页面）

顶层配置，声明所需 BFF 模块和布局结构：

```typescript
interface PageOrchestrationConfig {
  id: string;
  // 首屏所需 BFF 模块（编排系统统一请求）
  modules: string[];
  meta?: {
    title?: string;
    description?: string;
  };
  regions: RegionConfig[];
}
```

---

## 3. 页面结构

每个 MPA 页面独立维护自己的编排系统：

```
pages/home/
├── main.tsx              # 入口：声明 BFF 模块、创建编排
├── orchestration.ts      # 纯布局配置（可迁到后端）
├── App.tsx               # 页面根组件（可选）
└── modules/              # 页面专用模块
    ├── Logo.tsx
    ├── Nav.tsx
    └── PostList.tsx
```

---

## 4. 数据流

### 4.1 完整数据流

```
┌─────────────────────────────────────────────────────────────────┐
│  main.tsx                                                        │
│  - 声明所需 BFF 模块: ['header', 'postList', 'footer']          │
│  - 创建编排系统                                                  │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│  编排系统初始化                                                  │
│  - 调用 POST /api/page 获取首屏数据                              │
│  - 数据存入 PageStore                                            │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│  OrchestrationRenderer                                          │
│  - 读取编排配置（regions/blocks/modules）                        │
│  - 递归渲染布局结构                                              │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│  Module 组件（如 PostList）                                      │
│  - 通过 useModuleData('postList') 从 PageStore 读取数据          │
│  - 渲染自身 UI                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 模块数据消费

模块通过 Hook 从 PageStore 读取数据：

```typescript
// modules/PostList.tsx
import { useModuleData } from '@ventus/store';

export const PostList = () => {
  // 从 PageStore 读取对应 BFF 模块的数据
  const { data, loading, error } = useModuleData('postList');
  
  if (loading) return <Skeleton />;
  if (error) return <ErrorFallback error={error} />;
  
  return (
    <div>
      {data.items.map(post => (
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  );
};
```

### 4.3 跳链参数获取

编排系统自动解析 URL，注入到所有模块：

```typescript
// modules/PostList.tsx
import { usePageProps } from '@ventus/store';

export const PostList = () => {
  const { data } = useModuleData('postList');
  const pageProps = usePageProps();  // 获取 URL 解析结果
  
  // 获取查询参数 ?tag=xxx
  const currentTag = pageProps.getQuery('tag');
  
  return <div>...</div>;
};
```

---

## 5. 配置示例

### 5.1 首页配置

```typescript
// pages/home/orchestration.ts
import type { PageOrchestrationConfig } from '@ventus/types';

export const homeConfig: PageOrchestrationConfig = {
  id: 'home',
  // 声明首屏所需 BFF 模块（编排系统统一请求）
  modules: ['header', 'postList', 'footer'],
  meta: { title: '首页' },
  regions: [
    {
      id: 'header',
      type: 'header',
      block: {
        type: 'block',
        flexDirection: 'row',
        gap: 'var(--spacing-md)',
        justifyContent: 'between',
        children: [
          { type: 'module', name: 'Logo' },
          { type: 'module', name: 'Nav' }
        ]
      }
    },
    {
      id: 'content',
      type: 'content',
      padding: 'var(--spacing-lg)',
      block: {
        type: 'block',
        flexDirection: 'column',
        gap: 'var(--spacing-md)',
        children: [
          { type: 'module', name: 'PostList' }
        ]
      }
    },
    {
      id: 'footer',
      type: 'footer',
      block: {
        type: 'block',
        flexDirection: 'row',
        justifyContent: 'center',
        children: [
          { type: 'module', name: 'Footer' }
        ]
      }
    }
  ]
};
```

### 5.2 入口文件

```typescript
// pages/home/main.tsx
import { createOrchestration } from '@ventus/orchestration';
import { createRequest } from '@ventus/request';
import { createRoot } from 'react-dom/client';
import { homeConfig } from './orchestration';
import { Logo, Nav, PostList, Footer } from './modules';

const request = createRequest({ baseURL: '/api' });

// 创建编排系统
const orchestration = createOrchestration({
  config: homeConfig,
  request,
  modules: {
    Logo,
    Nav,
    PostList,
    Footer
  }
});

// 渲染（OrchestrationRenderer 从 Store 读取配置和数据）
createRoot(document.getElementById('root')!).render(
  <orchestration.Renderer />
);
```

### 5.3 带侧边栏的复杂布局

```typescript
// 文章详情页示例
export const postConfig: PageOrchestrationConfig = {
  id: 'post',
  modules: ['header', 'article', 'toc', 'footer'],
  regions: [
    {
      id: 'header',
      type: 'header',
      block: {
        type: 'block',
        flexDirection: 'row',
        children: [
          { type: 'module', name: 'Logo' },
          { type: 'module', name: 'Nav' }
        ]
      }
    },
    {
      id: 'content',
      type: 'content',
      padding: 'var(--spacing-lg)',
      block: {
        type: 'block',
        flexDirection: 'row',
        gap: 'var(--spacing-xl)',
        children: [
          // 主内容区
          {
            type: 'block',
            flexDirection: 'column',
            style: { flex: '1' },
            children: [
              { type: 'module', name: 'Article' }
            ]
          },
          // 侧边栏
          {
            type: 'block',
            flexDirection: 'column',
            style: { width: '280px' },
            children: [
              { type: 'module', name: 'TOC' }
            ]
          }
        ]
      }
    },
    {
      id: 'footer',
      type: 'footer',
      block: {
        type: 'block',
        children: [
          { type: 'module', name: 'Footer' }
        ]
      }
    }
  ]
};
```

---

## 6. Store 结构

编排系统依赖的 Store 结构：

```typescript
interface OrchestrationStore {
  // 布局配置
  config: PageOrchestrationConfig;
  
  // 模块注册表
  registry: Map<string, ComponentType>;
  
  // BFF 数据（只读）
  modules: Record<string, {
    data?: any;
    loading: boolean;
    error?: Error;
  }>;
  
  // URL 解析结果
  pageProps: {
    getParam: (key: string) => string | undefined;
    getQuery: (key: string) => string | undefined;
    params: Record<string, string>;
    query: Record<string, string>;
  };
  
  // Request 实例
  request: RequestInstance;
}
```

---

## 7. 向后端迁移

由于配置是纯数据的，未来迁移到后端只需修改数据获取方式：

```typescript
// 阶段1：前端静态配置（当前）
import { homeConfig } from './orchestration';

const orchestration = createOrchestration({
  config: homeConfig,
  request,
  modules
});

// 阶段2：后端下发配置（未来）
const config = await request.call('page.config', { page: 'home' });

const orchestration = createOrchestration({
  config,  // 后端下发的配置
  request,
  modules  // 模块组件仍需前端注册
});

// 渲染层代码完全不变
<orchestration.Renderer />
```

---

## 8. 与现有代码对比

| 现有方案 | 新编排系统 |
|---------|-----------|
| 页面内手动组装组件 | 声明式布局配置 |
| 硬编码布局结构 | Block/Region 配置化布局 |
| 直接调用 BFF Hook | 编排系统自动请求，模块从 Store 读取 |
| 硬编码间距值 | CSS 变量控制间距 |
| 模块直接 import | 通过注册表映射 |

---

## 9. 与 API SDK 的关系

编排系统基于 @ventus/api-client（由 OpenAPI 生成）进行数据请求。

```
┌─────────────────────────────────────────────────────────────┐
│                     前端架构关系图                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐     ┌─────────────────────────────┐   │
│  │ @ventus/api-client│◄──│ OpenAPI 契约 (api/openapi.yml)│   │
│  │  - types.ts     │     └─────────────────────────────┘   │
│  │  - client.ts    │              生成                     │
│  │  - hooks.ts     │                                         │
│  └────────┬────────┘                                         │
│           │                                                  │
│           ▼                                                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │           @ventus/orchestration                       │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │ init store  │  │ fetch BFF   │  │ render      │  │   │
│  │  │ parse URL   │  │ data        │  │ layout      │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  └────────────────────────┬─────────────────────────────┘   │
│                           │                                  │
│                           ▼                                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              @ventus/store                            │   │
│  │  ┌──────────────────────────────────────────────┐   │   │
│  │  │ PageStore (BFF data)                         │   │   │
│  │  │  modules.postList → { data, loading, error } │   │   │
│  │  └──────────────────────────────────────────────┘   │   │
│  └────────────────────────┬─────────────────────────────┘   │
│                           │                                  │
│                           ▼                                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Module Components (Logo, Nav, PostList...)          │   │
│  │  - useModuleData('postList')                         │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 9.1 @ventus/api-client（由 OpenAPI 生成）

类型安全的 API 客户端，从 OpenAPI 契约自动生成。

```typescript
import { getPageData, PageRequest, Post } from '@ventus/api-client';

// 完全类型安全
const response = await getPageData({
  page: 'home',
  modules: ['header', 'postList'],
  params: { page: 1 }
});

// response.modules.postList 有完整类型提示
```

详见 [api-design.md](./api-design.md)

### 9.2 @ventus/store

基于 Zustand 的状态管理：

```typescript
import { useModuleData, usePageProps } from '@ventus/store';

// 模块从 Store 读取 BFF 数据
function PostList() {
  const { data, loading, error } = useModuleData('postList');
  const pageProps = usePageProps();
  
  const currentPage = pageProps.getQuery('page') || '1';
  
  if (loading) return <Loading />;
  return <div>{data.items.map(...)}</div>;
}
```

### 9.3 @ventus/orchestration

编排渲染引擎，内部使用 @ventus/api-client：

```typescript
import { createOrchestration } from '@ventus/orchestration';
import { QueryClient } from '@tanstack/react-query';
import { homeConfig } from './orchestration';

const queryClient = new QueryClient();

const orchestration = createOrchestration({
  config: homeConfig,
  queryClient,  // 用于数据获取和缓存
  modules: { Logo, Nav, PostList, Footer }
});

createRoot(document.getElementById('root')!).render(
  <orchestration.Renderer />
);
```

---

## 10. 类型定义

### 10.1 编排配置

```typescript
// orchestration.ts
import type { PageOrchestrationConfig } from '@ventus/types';

export const homeConfig: PageOrchestrationConfig = {
  id: 'home',
  modules: ['header', 'postList', 'footer'],
  regions: [
    {
      id: 'header',
      type: 'header',
      block: {
        type: 'block',
        flexDirection: 'row',
        children: [
          { type: 'module', name: 'Logo' },
          { type: 'module', name: 'Nav' }
        ]
      }
    },
    // ...
  ]
};
```

### 10.2 模块组件

```typescript
// modules/PostList.tsx
import { useModuleData } from '@ventus/store';
import type { PostListData } from '@ventus/types';

export const PostList = () => {
  const { data, loading, error } = useModuleData<PostListData>('postList');
  
  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  
  return (
    <div>
      {data?.items.map(post => (
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  );
};
```

### 10.3 入口文件

```typescript
// main.tsx
import { createOrchestration } from '@ventus/orchestration';
import { createRequest } from '@ventus/request';
import { createRoot } from 'react-dom/client';
import { homeConfig } from './orchestration';
import { Logo, Nav, PostList, Footer } from './modules';

const request = createRequest({ baseURL: '/api' });

const orchestration = createOrchestration({
  config: homeConfig,
  request,
  modules: {
    Logo,
    Nav,
    PostList,
    Footer
  }
});

createRoot(document.getElementById('root')!).render(
  <orchestration.Renderer />
);
```

---

## 11. 待确认事项

1. **Block 是否需要支持 `width` / `minWidth` 等属性？**
   - 如侧边栏固定宽度 300px
   - 建议：通过 `style` 字段支持

2. **是否需要内置通用模块库？**
   - 如 `Container`、`Divider` 等布局辅助模块
   - 建议：P1 再做，MVP 只提供基础 Block

3. **是否需要支持条件渲染？**
   - 如根据登录状态显示不同 Module
   - 建议：P1 再做，MVP 通过多个页面实现

4. **错误处理策略**
   - 单个模块失败时，是显示错误占位还是整个页面失败？
   - 建议：模块级错误处理，不影响其他模块

5. **Loading 策略**
   - 页面级 Skeleton 还是模块级 Skeleton？
   - 建议：模块级，每个模块自行处理 loading 状态
