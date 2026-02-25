# next-ai-ventus 架构设计

## 核心思想

**Monorepo + 多页面构建 + 共享 Chunk 复用**

- 各页面（URL）完全独立，有自己的 HTML 入口
- 页面通过组装 npm 包（`@next-ai-ventus/ui`, `@next-ai-ventus/markdown` 等）构建
- 构建时自动提取共享代码为独立 chunk，实现跨页面复用

---

## 目录结构

```
next-ai-ventus/
├── frontend/                   # 前端代码
│   ├── pnpm-workspace.yaml     # 工作区配置
│   ├── package.json            # 根脚本
│   │
│   ├── packages/               # 共享 npm 包（被页面组装）
│   │   ├── ui/                 # @next-ai-ventus/ui - 组件库
│   │   ├── utils/              # @next-ai-ventus/utils - 工具库
│   │   ├── markdown/           # @next-ai-ventus/markdown - Markdown 渲染
│   │   └── types/              # @next-ai-ventus/types - 共享类型/Schema
│   │
│   ├── pages/                  # 页面入口（组装共享包）
│   │   ├── home/               # 首页：UI + PostList
│   │   ├── post/               # 文章页：UI + Markdown
│   │   ├── admin/              # 管理后台：UI + Markdown(Editor)
│   │   └── playground/         # WASM 运行页：UI + Canvas
│   │
│   └── shell/                  # 构建壳（统一打包配置）
│       └── vite.config.ts      # 多入口 + 代码分割配置
│
└── server/                     # Go 后端（无数据库）
    ├── cmd/server/main.go
    ├── internal/
    │   ├── store/              # 文件存储层
    │   ├── service/            # 业务逻辑层
    │   ├── handler/            # HTTP 处理器
    │   ├── middleware/         # 中间件
    │   └── router/             # 路由
    ├── content/                # 文章内容（Git 管理）
    └── storage/                # 上传文件
```

---

## 关键设计

### 1. 共享包（packages）

每个包独立开发，通过 `workspace:*` 被页面引用。

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

### 2. 页面组装（pages）

每个页面是独立应用，通过 `entry.tsx` 组装所需模块。

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

### 3. 构建壳（shell）

统一配置多入口构建和代码分割。

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

### 核心设计

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

## 后端架构（无数据库版）

### 核心思想

**有后端，无数据库，文件即数据**

- 后端 Go 服务处理 HTTP 请求、业务逻辑、权限控制
- 数据以文件形式存储在 `content/` 目录
- 启动时加载索引到内存，运行时内存查询
- 文件变更可触发 Git 提交，实现版本管理

### 技术选型

| 组件 | 选择 | 理由 |
|------|------|------|
| Web 框架 | Gin | 轻量、高性能 |
| 数据存储 | 文件系统 | 无需数据库服务 |
| 搜索 | 内存倒排索引 | 轻量，零依赖 |
| 认证 | JWT | 无状态 |

### 数据模型

```
content/posts/2024-06-hello-world/
├── meta.json           # 元数据
└── content.md          # Markdown 正文
```

**meta.json:**
```json
{
  "id": "2024-06-hello-world",
  "title": "Hello World",
  "slug": "hello-world",
  "date": "2024-06-15",
  "tags": ["go", "架构"],
  "status": "published"
}
```

### Store 层

```go
type PostStore interface {
    List(opts ListOptions) ([]Post, int)     // 内存查询
    Get(id string) (*Post, error)
    Create(post *Post) error                 // 写文件
    Update(id string, post *Post) error
    Delete(id string) error
    ReloadIndex() error                      // 刷新索引
}
```

**内存索引:**
```go
type Index struct {
    Posts   []Post            // 文章列表（按时间排序）
    SlugMap map[string]string // slug -> id
    TagMap  map[string][]string // tag -> ids
    mu      sync.RWMutex
}
```

### API 设计

```
公开 API
├── GET  /api/posts              # 列表
├── GET  /api/posts?tag=go       # 按标签筛选
├── GET  /api/posts/:slug        # 详情
├── POST /api/posts/:id/view     # 记录阅读
├── GET  /api/search?q=xxx       # 搜索
└── GET  /api/tags               # 标签列表

管理 API（需 JWT 认证）
├── POST   /api/admin/posts
├── PUT    /api/admin/posts/:id
├── DELETE /api/admin/posts/:id
└── POST   /api/admin/upload     # 图片上传
```

### 图片上传

```
storage/uploads/
├── 2024/
│   └── 06/
│       └── abc123.jpg    # 重命名存储
```

---

## 部署配置

### Nginx

```nginx
server {
    listen 80;
    root /var/www/blog/frontend/dist;
    
    # 首页
    location = / {
        try_files /pages/home/index.html =404;
    }
    
    # 文章页
    location /post/ {
        alias /var/www/blog/frontend/dist/pages/post/;
        try_files $uri $uri/ /pages/post/index.html;
    }
    
    # 管理后台
    location /admin/ {
        alias /var/www/blog/frontend/dist/pages/admin/;
        try_files $uri $uri/ /pages/admin/index.html;
    }
    
    # 共享 chunk 长期缓存
    location /shared/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # API
    location /api/ {
        proxy_pass http://localhost:8080;
    }
}
```

### Docker Compose

```yaml
version: '3.8'

services:
  frontend:
    build: ./frontend

  server:
    build: ./server
    environment:
      - JWT_SECRET=${JWT_SECRET}
    volumes:
      - ./content:/app/content:rw
      - ./storage:/app/storage:rw

  nginx:
    image: nginx:alpine
    ports: ["80:80"]
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
```

---

## 决策记录

### 为什么不用 Single Page Application (SPA)？

- 各页面完全独立，一个页面崩溃不影响其他
- 首屏加载更快（只加载该页面所需代码）
- 更好的 SEO（独立 URL + HTML）

### 为什么不用数据库？

1. **部署简单** - 单二进制文件，无需数据库服务
2. **数据可视** - 文章就是 Markdown 文件，随时可编辑
3. **版本管理** - Git 天然支持内容版本回滚
4. **备份简单** - `git push` 即备份
5. **性能足够** - 个人博客文章数有限，内存索引足够快

### 什么时候应该加数据库？

- 文章数 > 10,000
- 多用户同时编辑（冲突频繁）
- 需要复杂查询（多表关联、聚合统计）
- 需要事务保证

### 2核2G 服务器能否支撑？

- **可以** - 静态文件走 Nginx，不经过后端
- Go 后端仅处理 API 请求，内存占用 < 20MB
- 共享 chunk 长期缓存在浏览器
