# Ventus 编排系统设计文档

## 1. 概述

Ventus 编排系统是一套用于 MPA（多页应用）的页面结构管理系统，支持四层结构：

```
Page (页面)
  └── Region (区域) - 页面垂直分区
        └── Block (区块) - Flex 布局容器（可嵌套）
              └── Module (模块) - 业务组件（叶子节点）
```

## 2. 核心设计原则

- **页面级隔离**：每个 MPA 页面拥有独立的注册表，互不影响
- **纯数据配置**：`orchestration.ts` 只包含纯数据，可 JSON 序列化，方便后端下发
- **语义化 Spacing**：使用业务语义 token（如 `navPadding`），非固定数值，由页面级解析器转换
- **单解析器**：每个页面只有一个 Spacing 解析器

## 3. 数据结构

### 3.1 Module（模块）

```typescript
interface ModuleConfig {
  type: 'module';
  name: string;  // 模块名称，对应注册表
}
```

### 3.2 Block（区块）

仅支持 Flex 布局，配置内外边距使用语义化 token：

```typescript
interface BlockConfig {
  type: 'block';
  flexDirection: 'row' | 'column' | 'row-reverse' | 'column-reverse';
  gap?: string;           // 子元素间距 token
  padding?: string;       // 内边距 token
  margin?: string;        // 外边距 token
  justifyContent?: 'start' | 'center' | 'end' | 'between' | 'around';
  alignItems?: 'start' | 'center' | 'end' | 'stretch';
  children: (BlockConfig | ModuleConfig)[];
}
```

### 3.3 Region（区域）

```typescript
interface RegionConfig {
  id: string;
  type: 'header' | 'content' | 'footer' | 'sidebar';
  padding?: string;   // 语义化 token
  margin?: string;    // 语义化 token
  block: BlockConfig; // 唯一根区块
}
```

### 3.4 Page（页面）

```typescript
interface PageOrchestrationConfig {
  id: string;
  meta?: {
    title?: string;
    description?: string;
  };
  regions: RegionConfig[];
}
```

## 4. 页面结构

每个 MPA 页面独立维护自己的编排系统：

```
pages/home/
├── main.tsx              # 入口：创建注册表、注册 resolver 和 modules
├── orchestration.ts      # 纯数据配置（可迁到后端）
├── App.tsx               # 使用 OrchestrationRenderer
├── modules/              # 页面专用模块
│   ├── Logo.tsx
│   ├── Nav.tsx
│   └── PostList.tsx
└── types.ts              # 页面类型（可选）
```

## 5. Spacing 解析器

每个页面在 `createOrchestration` 时传入解析器，将语义化 token 转为具体数值：

```typescript
// main.tsx
const orchestration = createOrchestration({
  config: homeConfig,
  request,
  modules,
  resolver: (token) => {
    const map: Record<string, string> = {
      // 首页业务语义
      navPadding: '24px',
      navGap: '16px',
      pagePadding: '24px',
      sectionGap: '32px',
      cardGap: '16px'
    };
    return map[token] || '0px';
  }
});
```

不同页面可以定义同名 token 但不同值，完全隔离。

## 7. 配置示例

### 7.1 纯数据配置

```typescript
// orchestration.ts
export const homeConfig: PageOrchestrationConfig = {
  id: 'home',
  meta: { title: '首页' },
  regions: [
    {
      id: 'header',
      type: 'header',
      padding: 'navPadding',  // 语义化 token
      block: {
        type: 'block',
        flexDirection: 'row',
        gap: 'navGap',
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
      padding: 'pagePadding',
      block: {
        type: 'block',
        flexDirection: 'row',
        gap: 'sectionGap',
        children: [
          {
            // 嵌套区块：主内容
            type: 'block',
            flexDirection: 'column',
            gap: 'cardGap',
            children: [
              { type: 'module', name: 'PostList' }
            ]
          },
          {
            // 嵌套区块：侧边栏
            type: 'block',
            flexDirection: 'column',
            gap: 'cardGap',
            padding: 'cardPadding',
            children: [
              { type: 'module', name: 'TagCloud' }
            ]
          }
        ]
      }
    }
  ]
};
```

### 7.2 入口文件

```typescript
// main.tsx
import { createOrchestration } from '@ventus/orchestration';
import { createRequest } from '@ventus/request';
import { createRoot } from 'react-dom/client';
import { homeConfig } from './orchestration';
import { Logo, Nav, PostList, TagCloud } from './modules';

const request = createRequest({ baseURL: '/api' });

const orchestration = createOrchestration({
  config: homeConfig,
  request,
  modules: {
    Logo,
    Nav,
    PostList,
    TagCloud
  },
  resolver: (token) => {
    const map = {
      navPadding: '24px',
      navGap: '16px',
      pagePadding: '24px',
      sectionGap: '32px',
      cardGap: '16px',
      cardPadding: '16px'
    };
    return map[token] || '0px';
  }
});

// 渲染（OrchestrationRenderer 从 store 读取配置）
createRoot(document.getElementById('root')!).render(
  <orchestration.Renderer />
);
```

## 8. 向后端迁移

由于配置是纯数据的，未来迁移到后端只需修改数据获取方式：

```typescript
// 阶段1：前端静态配置（当前）
import { homeConfig } from './orchestration';

const orchestration = createOrchestration({
  config: homeConfig,
  request,
  modules,
  resolver
});

// 阶段2：后端下发配置（未来）
const homeConfig = await fetch('/api/page-config?id=home').then(r => r.json());

const orchestration = createOrchestration({
  config: homeConfig,  // 后端下发的配置
  request,
  modules,
  resolver  // 解析器仍在前端，或也从后端下发
});

// 渲染层代码完全不变
<orchestration.Renderer />
```

## 9. 与现有代码对比

| 现有方案 | 新编排系统 |
|---------|-----------|
| 直接调用 `fetchPageData` | 通过编排配置声明式获取 |
| 组件内手动布局 | Block 配置化布局 |
| 硬编码间距值 | 语义化 token + 解析器 |
| 模块直接 import | 通过注册表懒加载 |

## 10. 包架构设计

编排系统与请求、状态管理打通，分为三个包：

### 10.1 @ventus/request

负责 HTTP 请求和数据获取：

```typescript
import { createRequest } from '@ventus/request';

const request = createRequest({ baseURL: '/api' });

// 调用 API
const data = await request.call<PostListData>({
  scene: 'post.list',
  params: { page: 1, pageSize: 10 }
});
```

### 10.2 @ventus/store

基于 Zustand + React Context 的状态管理：

```typescript
import { useRequest } from '@ventus/store';
import type { PageProps } from '@ventus/types';

// 模块通过 props 接收跳链参数
interface PostListProps {
  pageProps: PageProps;  // 编排系统注入
}

function PostList({ pageProps }: PostListProps) {
  const slug = pageProps.getParam('slug');
  const page = pageProps.getQuery('page');
  
  // 模块自主获取数据
  const { data, loading } = useRequest<PostListData>({
    scene: 'post.list',
    params: { slug, page }
  });
  
  if (loading) return <Loading />;
  return <div>{data.items.map(...)}</div>;
}
```

### 10.3 @ventus/orchestration

编排渲染引擎，负责：
1. 初始化 store（写入 config、registry、request）
2. 解析 URL 跳链参数
3. 渲染页面结构

```typescript
// main.tsx
import { createOrchestration } from '@ventus/orchestration';
import { createRequest } from '@ventus/request';
import { homeConfig } from './orchestration';
import * as modules from './modules';

const request = createRequest({ baseURL: '/api' });

// 创建编排系统（初始化 store）
const orchestration = createOrchestration({
  config: homeConfig,
  request,
  modules,
  resolver: (token) => ({
    navPadding: '24px',
    navGap: '16px'
  }[token] || '0px')
});

// 渲染（OrchestrationRenderer 从 store 读取配置）
createRoot(document.getElementById('root')!).render(
  <orchestration.Renderer />
);
```

## 11. 数据流设计

### 11.1 跳链参数注入

编排系统自动解析 URL，通过 props 传递给所有 Module：

```
URL: /post/my-article?page=2&source=wechat
       ↓
Orchestration 解析为 pageProps
       ↓
渲染 Module 时通过 props 传入
       ↓
模块通过 props.pageProps 获取
```

```typescript
// modules/PostList.tsx
import { useRequest } from '@ventus/store';
import type { PageProps } from '@ventus/types';

interface PostListProps {
  pageProps: PageProps;  // 编排系统通过 props 注入
}

export const PostList = ({ pageProps }: PostListProps) => {
  const currentPage = pageProps.getQuery('page') || '1';
  const source = pageProps.getQuery('source');
  
  // 模块自主请求数据
  const { data } = useRequest<PostListData>({
    scene: 'post.list',
    params: { 
      page: parseInt(currentPage),
      source 
    }
  });
  
  return <div>...</div>;
};
```

### 11.2 模块自主数据获取

每个模块通过 `useRequest` 独立获取数据，编排系统不负责数据分发：

```typescript
// modules/TagCloud.tsx
import { useRequest } from '@ventus/store';

export const TagCloud = () => {
  // 独立请求，与其他模块互不干扰
  const { data, loading } = useRequest<TagData>({
    scene: 'tag.hotList',
    params: { limit: 20 }
  });
  
  if (loading) return <Skeleton />;
  return <div>{data.tags.map(...)}</div>;
};
```

### 11.3 数据缓存

`@ventus/store` 基于 Zustand 实现请求缓存：

```typescript
// 相同请求自动复用
const { data } = useRequest({ scene: 'tag.hotList', params: { limit: 20 } });
// 多个模块使用相同配置时，只发一次请求
```

## 12. 类型定义

### 12.1 PageProps

跳链参数类型，编排系统解析 URL 后生成：

```typescript
interface PageProps {
  /** 获取路由参数 /post/:slug -> getParam('slug') */
  getParam: (key: string) => string | undefined;
  /** 获取查询参数 ?page=2 -> getQuery('page') */
  getQuery: (key: string) => string | undefined;
  /** 获取所有路由参数 */
  params: Record<string, string>;
  /** 获取所有查询参数 */
  query: Record<string, string>;
}
```

### 12.2 编排配置

```typescript
// orchestration.ts
import type { PageOrchestrationConfig } from '@ventus/types';

export const homeConfig: PageOrchestrationConfig = {
  id: 'home',
  regions: [
    {
      id: 'header',
      type: 'header',
      padding: 'navPadding',
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
      padding: 'pagePadding',
      block: {
        type: 'block',
        flexDirection: 'row',
        gap: 'sectionGap',
        children: [
          { type: 'module', name: 'PostList' },
          { type: 'module', name: 'TagCloud' }
        ]
      }
    }
  ]
};
```

### 12.3 模块组件

```typescript
// modules/PostList.tsx
import { useRequest } from '@ventus/store';
import type { PageProps, PostListItem } from '@ventus/types';

interface PostListData {
  items: PostListItem[];
  total: number;
}

interface PostListProps {
  pageProps: PageProps;  // 编排系统注入
}

export const PostList = ({ pageProps }: PostListProps) => {
  const tag = pageProps.getQuery('tag');
  
  const { data, loading } = useRequest<PostListData>({
    scene: 'post.list',
    params: { tag, page: 1 }
  });
  
  if (loading) return <div>Loading...</div>;
  
  return (
    <div>
      {data?.items.map(post => (
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  );
};
```

### 12.4 入口文件

```typescript
// main.tsx
import { createOrchestration } from '@ventus/orchestration';
import { createRequest } from '@ventus/request';
import { createRoot } from 'react-dom/client';
import { homeConfig } from './orchestration';
import { Logo, Nav, PostList, TagCloud } from './modules';

const request = createRequest({ baseURL: '/api' });

const orchestration = createOrchestration({
  config: homeConfig,
  request,
  modules: {
    Logo,
    Nav,
    PostList,
    TagCloud
  },
  resolver: (token) => ({
    navPadding: '24px',
    navGap: '16px',
    pagePadding: '24px',
    sectionGap: '32px'
  }[token] || '0px')
});

createRoot(document.getElementById('root')!).render(
  <orchestration.Renderer />
);
```

## 13. 待确认事项

1. **Block 是否需要支持 `width` / `minWidth` 等属性？**
   - 如侧边栏固定宽度 300px

2. **是否需要内置通用模块库？**
   - 如 `Container`、`Divider` 等布局辅助模块

3. **是否需要支持条件渲染？**
   - 如根据登录状态显示不同 Module

4. **Module 是否需要统一的基础 Props？**
   - 如 `className`、`style` 等，方便编排系统控制

5. **Store 结构确认**
   ```typescript
   interface OrchestrationStore {
     config: PageOrchestrationConfig;  // 页面配置
     registry: {
       resolver: SpacingResolver;
       modules: Map<string, ComponentType>;
     };
     request: RequestInstance;  // request 实例
     url: {
       params: Record<string, string>;  // 解析后的路由参数
       query: Record<string, string>;   // 解析后的查询参数
     };
   }
   ```
