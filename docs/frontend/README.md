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
}))
```

### 分工

- **服务端状态**：BFF 返回，只读，存入 PageStore
- **客户端状态**：Zustand 管理，可读可写

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
        home: './pages/home/index.html',
        post: './pages/post/index.html',
        admin: './pages/admin/index.html',
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
