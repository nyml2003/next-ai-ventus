# 前端详细设计

## 目录

- [工程化配置](#工程化配置)
- [主题系统](#主题系统)
- [状态管理](#状态管理)
- [路由策略](#路由策略)
- [浏览器兼容性](#浏览器兼容性)
- [构建配置](#构建配置)

---

## 工程化配置

### 代码规范

- **ESLint 9** (flat config) + **Prettier**
- **simple-git-hooks** 提交前检查
- TypeScript 严格模式

```bash
pnpm lint:fix      # 自动修复
pnpm format        # 格式化
```

### 路径别名

```json
{
  "@/*": ["./src/*"],
  "@ui/*": ["../packages/ui/src/*"],
  "@utils/*": ["../packages/utils/src/*"],
  "@types/*": ["../packages/types/src/*"],
  "@markdown/*": ["../packages/markdown/src/*"]
}
```

### 测试（Jest）

```bash
pnpm test          # 运行测试
pnpm test:coverage # 覆盖率
pnpm test:watch    # Watch 模式
```

---

## 主题系统

自定义设计系统，不依赖 Tailwind。

### 设计令牌

```typescript
// tokens/colors.ts
export const colors = {
  primary: { 50: '#eff6ff', 500: '#3b82f6', 600: '#2563eb' },
  neutral: { 0: '#fff', 100: '#f5f5f5', 900: '#171717' },
}

// themes/light.ts
export const lightTheme = {
  colors: {
    bgPrimary: colors.neutral[0],
    textPrimary: colors.neutral[900],
    primary: colors.primary[500],
  },
  spacing: { sm: '8px', md: '16px', lg: '24px' },
  radius: { sm: '4px', md: '8px' },
}
```

### CSS 变量

```css
:root {
  --color-bg-primary: #ffffff;
  --color-text-primary: #171717;
  --color-primary: #3b82f6;
}

[data-theme="dark"] {
  --color-bg-primary: #171717;
  --color-text-primary: #ffffff;
}
```

---

## 状态管理

**方案：Zustand**

### BFF 数据获取（首屏模块必须）

**封装 Hook：usePageData**

```typescript
// 获取页面 BFF 数据
function usePageData(
  page: string,                    // 页面标识：home, post, admin, adminPosts...
  modules: string[],              // 所需模块列表
  params?: Record<string, any>    // 页面参数：{ page: 1, tag: 'go' }
) {
  const { data, error, isLoading } = useSWR(
    ['/api/page', page, modules, params],
    () => fetchPageData(page, modules, params),
    { revalidateOnFocus: false }
  );
  
  return {
    modules: data?.modules || {},   // 各模块数据
    meta: data?.meta || {},         // 页面元数据
    isLoading,
    error
  };
}

// 使用示例
function HomePage() {
  const { modules, isLoading } = usePageData('home', [
    'header', 'hero', 'postList', 'sidebar', 'footer'
  ], { page: 1 });
  
  if (isLoading) return <PageSkeleton />;
  
  return (
    <Layout>
      <Header data={modules.header.data} />
      <Hero data={modules.hero.data} />
      <PostList data={modules.postList.data} />
      <Sidebar data={modules.sidebar.data} />
      <Footer data={modules.footer.data} />
    </Layout>
  );
}
```

### 服务端状态（BFF 数据）

```typescript
const usePageStore = create(() => ({
  modules: {},      // BFF 返回的模块数据
  meta: {},         // 页面元数据
}))
```

### 客户端状态（UI 状态）

```typescript
const useUIStore = create(() => ({
  theme: 'light',   // 主题
  sidebarOpen: false,
  // 非首屏交互状态
  selectedIds: [],  // 表格选中项
  previewOpen: false,
}))
```

### 分工

| 状态类型 | 管理方式 | 来源 | 说明 |
|---------|---------|------|------|
| **首屏数据** | BFF Hook | 后端 | 必须通过 `/api/page` 获取 |
| **服务端状态** | PageStore | BFF 返回 | 只读，模块自取 |
| **客户端状态** | UIStore | 前端 | 可读可写，UI 交互 |
| **表单状态** | 组件内 State | 前端 | 表单输入等临时状态 |

---

## 路由策略

**无子路由设计**

- 每个页面是独立 HTML，无 React Router
- 独立页面通过不同 HTML 入口实现：

```
/                  # 首页（文章列表）
/login             # 登录页（独立页面，非弹窗）
/admin             # 管理首页（仪表盘）
/admin/posts       # 文章管理页（独立页面）
/admin/editor      # 新建文章
/admin/editor?id=123  # 编辑文章（通过 id 区分）
```

> 登录页使用独立页面而非弹窗，原因：
> 1. 管理后台需要登录态，独立页面便于做统一登录检查
> 2. 登录页可以作为未授权访问时的统一跳转目标
> 3. 简单直接，不需要处理弹窗的遮罩、关闭等交互

---

## 浏览器兼容性

- **目标**：PC Chrome（最新版）
- **不使用**：ES2022+ 新特性（避免 transpile 问题）
- **CSS**：使用成熟特性，避免 Container Query 等较新特性
- **不做**：移动端适配、响应式

---

## 构建配置

### 多入口构建

```typescript
// shell/vite.config.ts
export default {
  build: {
    rollupOptions: {
      input: {
        // C 端页面
        home: './pages/home/index.html',
        post: './pages/post/index.html',
        // B 端页面
        login: './pages/login/index.html',
        admin: './pages/admin/index.html',
        adminPosts: './pages/admin-posts/index.html',
        adminEditor: './pages/admin-editor/index.html',
        adminImages: './pages/admin-images/index.html',
      },
      output: {
        manualChunks: {
          'vendor-react': ['react', 'react-dom'],
          'ui': ['@next-ai-ventus/ui'],
        }
      }
    }
  }
}
```

### 页面初始化流程

每个页面入口统一处理 BFF 数据获取：

```typescript
// pages/home/entry.tsx
import { usePageData } from '@utils/bff';

function HomeEntry() {
  // 首屏模块列表与产品文档保持一致
  const { modules, isLoading, error } = usePageData('home', [
    'header',     // 顶部导航
    'hero',       // 首屏横幅
    'postList',   // 文章列表（首屏核心）
    'sidebar',    // 侧边栏
    'footer'      // 页脚
  ]);

  if (isLoading) return <HomeSkeleton />;
  if (error) return <ErrorPage error={error} />;

  return <HomePage modules={modules} />;
}
```

### 共享包结构

```
frontend/packages/ui/src/
├── theme/                  # 主题系统
│   ├── tokens/             # 设计令牌
│   ├── themes/             # 主题定义
│   ├── ThemeProvider.tsx
│   └── global.css
├── components/
│   ├── Button/
│   ├── Layout/
│   └── PostList/
└── index.ts
```

### 页面组装

```tsx
// frontend/pages/post/entry.tsx
import { Layout, ArticleMeta } from '@ui/components'
import { MarkdownRenderer } from '@markdown/renderer'

function PostPage() {
  return (
    <Layout>
      <ArticleMeta />
      <MarkdownRenderer />
    </Layout>
  )
}
```
