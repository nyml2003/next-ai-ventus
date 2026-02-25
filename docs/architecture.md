# next-ai-ventus 架构设计

## 核心思想

**Monorepo + MPA + 无数据库 + BFF**

- 各页面（URL）完全独立，有自己的 HTML 入口
- 页面通过组装 npm 包构建，共享代码自动提取为独立 chunk
- 数据以文件形式存储，Git 管理，启动时加载索引到内存
- BFF 统一接口 `/api/page`，后端按需组装模块数据

---

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                         Nginx                               │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐ │
│  │  /      │  │ /post/  │  │ /admin/ │  │ /api, /uploads  │ │
│  │ 首页    │  │ 文章页  │  │ 管理后台│  │ → Go Backend    │ │
│  └─────────┘  └─────────┘  └─────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┴─────────────────────┐
        ▼                                           ▼
┌─────────────────────────┐           ┌─────────────────────────┐
│     Frontend (MPA)      │           │    Backend (Go + Gin)   │
│  ┌─────────────────┐    │           │  ┌─────────────────┐    │
│  │  pages/home     │    │           │  │  BFF Handler    │    │
│  │  pages/post     │    │◄──────────┼──┤  /api/page      │    │
│  │  pages/admin    │    │  /api/page│  └─────────────────┘    │
│  │                 │    │           │           │             │
│  │  packages/ui    │    │           │  ┌────────┴────────┐    │
│  │  packages/mark  │    │           │  ▼                 ▼    │
│  └─────────────────┘    │           │ ┌─────────┐   ┌────────┐│
└─────────────────────────┘           │ │ Modules │   │ Store  ││
                                      │ │(并行执  │   │(文件  ││
                                      │ │ 行)     │   │ 系统) ││
                                      │ └─────────┘   └────────┘│
                                      │       │                 │
                                      │       ▼                 │
                                      │ ┌─────────────────────┐ │
                                      │ │ content/            │ │
                                      │ │  ├── posts/         │ │
                                      │ │  └── authors/       │ │
                                      │ └─────────────────────┘ │
                                      └─────────────────────────┘
```

---

## 目录结构

```
next-ai-ventus/
├── content/                    # 文章内容（Git 管理）
│   ├── posts/                  # 文章目录
│   │   └── 2024-06-hello/      # 每篇文章一个目录（目录名 = 文章ID）
│   │       ├── meta.json       # 元数据（标题、slug、标签、状态、版本等）
│   │       └── content.md      # Markdown 正文
│   └── authors/                # 作者信息
│
├── frontend/                   # 前端代码
│   ├── packages/               # 共享 npm 包
│   │   ├── ui/                 # @next-ai-ventus/ui
│   │   ├── utils/              # @next-ai-ventus/utils
│   │   ├── markdown/           # @next-ai-ventus/markdown
│   │   └── types/              # @next-ai-ventus/types
│   ├── pages/                  # 页面入口
│   │   ├── home/
│   │   ├── post/
│   │   ├── admin/
│   │   └── playground/
│   └── shell/                  # 统一构建配置
│
├── server/                     # Go 后端
│   ├── internal/
│   │   ├── store/              # 文件存储层
│   │   ├── service/            # 业务逻辑层
│   │   ├── handler/            # HTTP 处理器
│   │   ├── middleware/         # 中间件
│   │   └── router/             # 路由
│   └── storage/                # 上传文件（运行时生成）
│
├── docker-compose.yml
└── nginx.conf
```

---

## BFF 架构（核心）

**前端声明需求，后端按需组装**

```
前端配置模块列表 → 发起请求 → 后端并行执行模块 → 统一返回 → 前端缓存 → 模块自取数据
```

### 统一接口

```http
POST /api/page

{
  "page": "home",              // 页面标识
  "modules": ["header", "hero", "postList", "sidebar", "footer"],
  "params": { "page": 1, "tag": "go" }
}
```

### 响应格式

```json
{
  "page": "home",
  "meta": { "title": "...", "description": "..." },
  "modules": {
    "header":   { "code": 200, "data": { ... } },
    "postList": { "code": 200, "data": { "items": [...], "pagination": {...} } },
    "sidebar":  { "code": 500, "error": "..." }   // 单个模块失败不影响其他
  }
}
```

### 模块注册

```go
type ModuleHandler func(ctx *ModuleContext) (interface{}, error)

var ModuleRegistry = map[string]ModuleHandler{
    "header":    modules.HandleHeader,
    "postList":  modules.HandlePostList,
    "article":   modules.HandleArticle,
    // ...
}
```

---

## 页面跳转设计

**后端拼链，前端只管使用**

后端返回数据时，每个资源自带完整的跳转 URL（`href` 字段）：

```json
{
  "list": [
    {
      "id": "2024-06-hello",
      "title": "Hello World",
      "href": "/post/hello-world?from=home&page=2&pos=0"
    }
  ],
  "backUrl": "/?page=2"
}
```

前端直接使用 `<a href={item.href}>`，不解析 URL 参数，不拼接链接。

---

## 技术栈

| 层 | 技术 |
|----|------|
| 前端 | React + TypeScript + 自定义主题系统 |
| 构建 | Vite + pnpm workspace |
| 后端 | Go + Gin |
| 存储 | 文件系统（content/ + storage/） |
| 部署 | Docker Compose + Nginx |

---

## 非目标（明确不做）

| 项 | 说明 |
|----|------|
| 移动端适配 | 不支持手机浏览器 |
| SSR | 使用 BFF + CSR |
| 数据库 | 文件存储足够 |
| IE/旧浏览器 | 仅 Chrome |
| 第三方服务 | 评论、统计、搜索（自建优先） |

---

## 关键决策

### 为什么不用 SPA？

- 各页面完全独立，一个页面崩溃不影响其他
- 首屏加载更快（只加载该页面所需代码）
- 更好的 SEO（独立 URL + HTML）

### 为什么不用数据库？

1. **部署简单** - 单二进制文件，无需数据库服务
2. **数据可视** - 文章就是 Markdown 文件，随时可编辑
3. **版本管理** - Git 天然支持内容版本回滚
4. **性能足够** - 个人博客文章数有限，内存索引足够快

### 为什么用 BFF 而不是 SSR？

- SSR 需要 Node.js 运行时，内存占用大（> 100MB）
- BFF 纯 API，前端静态部署，Nginx 直接 serve
- 内容不频繁变化，CSR + 缓存足够满足需求

---

## 页面与 BFF 模块映射

| 页面 | MVP 首屏 BFF 模块 | P1 扩展模块 | 说明 |
|------|------------------|-------------|------|
| **首页** `/` | `header`, `postList`, `footer` | `hero`, `sidebar` | MVP 极简布局 |
| **文章页** `/post/:slug` | `header`, `article`, `footer` | `toc`, `related` | article 包含完整内容 |
| **管理首页** `/admin` | - | `adminSidebar`, `adminStats`, `recentPosts` | **MVP 不做**，直接进 `/admin/posts` |
| **文章管理** `/admin/posts` | `adminSidebar`, `adminFilter`, `adminPostList` | - | 含文章统计 |
| **文章编辑** `/admin/editor` | `editor`, `editorSettings` | - | Editor、Settings 面板 |
| **图片管理** `/admin/images` | - | `adminSidebar`, `imageFilter`, `imageList` | P1 功能 |
| **登录页** `/login` | - | - | 纯前端表单，无 BFF |

> **原则**：首屏可见的信息展示模块必须通过 BFF 获取数据，确保首屏渲染性能。

## 详细设计文档

- [设计原则](./principles.md) - **必读**：BFF 模块设计、数据流、性能、安全原则
- [前端详细设计](./frontend/README.md) - 工程化、主题、状态管理、BFF 数据获取
- [后端详细设计](./server/README.md) - 存储、API、模块机制
- [页面详细设计](./pages/README.md) - 各页面模块组成和交互
  - [登录页](./pages/login.md) - P0：JWT 登录
  - [首页](./pages/home.md) - P0：文章列表
  - [文章详情](./pages/post.md) - P0：文章阅读
  - [管理首页](./pages/admin-home.md) - P0：仪表盘
  - [文章管理](./pages/admin-posts.md) - P0：文章列表管理
  - [文章编辑](./pages/admin-editor.md) - P0：Markdown 编辑器
  - [图片管理](./pages/admin-images.md) - P1：图片上传管理
