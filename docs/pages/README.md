# 页面设计文档

本文档从**产品视角**描述每个页面的模块组成、用户场景和交互设计。

> **重要约定**：
> - 这里的「模块」指页面上的**功能区块/组件**
> - **首屏信息展示模块**必须通过 BFF 获取数据，对应具体的 BFF 模块名
> - 各页面的 BFF 模块映射见 [architecture.md 的页面与 BFF 模块映射表](../architecture.md#页面与-bff-模块映射)

---

## 文档结构

每个页面文档包含：
- **页面目标** - 这个页面解决什么问题
- **用户场景** - 用户在这个页面做什么
- **模块清单** - 页面上有哪些功能区块
- **布局说明** - 模块的排列方式
- **交互说明** - 用户如何与页面交互

---

## 页面索引

### C 端页面

| 页面 | 文档 | 核心模块 | BFF 模块 |
|------|------|---------|----------|
| 首页 | [home.md](./home.md) | Header, Hero, PostList, Sidebar, Footer | `header`, `hero`, `postList`, `sidebar`, `footer` |
| 文章详情 | [post.md](./post.md) | Header, Article, TOC, Related, Footer | `header`, `article`, `toc`, `related` |
| 登录页 | [login.md](./login.md) | LoginForm | 无（纯前端表单）|

### B 端页面

| 页面 | 文档 | 核心模块 | BFF 模块 |
|------|------|---------|----------|
| 管理首页 | [admin-home.md](./admin-home.md) | Sidebar, Stats, QuickActions, RecentPosts | `adminSidebar`, `adminStats`, `recentPosts` |
| 文章管理 | [admin-posts.md](./admin-posts.md) | Sidebar, FilterBar, PostTable, Pagination | `adminSidebar`, `adminFilter`, `adminPostList` |
| 文章编辑 | [admin-editor.md](./admin-editor.md) | Toolbar, Editor, Preview, Settings | `editor`, `editorSettings` |
| 图片管理 | [admin-images.md](./admin-images.md) | Sidebar, FilterBar, ImageGrid, Pagination | `adminSidebar`, `imageFilter`, `imageList` |

---

## 模块复用关系

```
┌─────────────────────────────────────────────────────────┐
│  Header（全站通用）- BFF: header                        │
│  - Logo                                                 │
│  - 导航链接（首页、关于、搜索）                          │
│  - 主题切换按钮                                          │
│  - 登录/头像（B 端显示）                                 │
└─────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐    ┌─────────────────┐    ┌──────────────┐
│   C 端首页    │    │   文章详情页     │    │   管理后台    │
│              │    │                 │    │              │
│  Hero        │    │  ArticleMeta    │    │  Sidebar     │
│  (BFF:hero)  │    │  (BFF:article)  │    │(BFF:adminSidebar)│
│  PostList    │    │  Article        │    │  Stats       │
│(BFF:postList)│    │  TOC            │    │(BFF:adminStats)  │
│  Sidebar     │    │(BFF:toc)        │    │  PostTable   │
│(BFF:sidebar) │    │  Related        │    │(BFF:adminPostList)│
│  Footer      │    │(BFF:related)    │    │  Editor      │
│(BFF:footer)  │    │  Footer         │    │(BFF:editor)  │
└──────────────┘    │(BFF:footer)     │    │  Preview     │
                    └─────────────────┘    │(同BFF:editor)│
                                           └──────────────┘
```

**说明**：括号内标注的是对应的 BFF 模块名，首屏模块必须走 BFF 获取数据。
