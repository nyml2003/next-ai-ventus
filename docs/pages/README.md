# 页面设计文档

本文档从**产品视角**描述每个页面的模块组成、用户场景和交互设计。

> 注意：这里的「模块」指页面上的**功能区块/组件**，不是 BFF 的后端模块。

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

| 页面 | 文档 | 核心模块 |
|------|------|---------|
| 首页 | [home.md](./home.md) | Header, Hero, PostList, Sidebar, Footer |
| 文章详情 | [post.md](./post.md) | Header, ArticleMeta, Article, TOC, Related, Footer |

### B 端页面

| 页面 | 文档 | 核心模块 |
|------|------|---------|
| 管理首页 | [admin-home.md](./admin-home.md) | Sidebar, Stats, QuickActions, RecentPosts |
| 文章管理 | [admin-posts.md](./admin-posts.md) | Sidebar, FilterBar, PostTable, Pagination |
| 文章编辑 | [admin-editor.md](./admin-editor.md) | Toolbar, Editor, Preview, Settings |

---

## 模块复用关系

```
┌─────────────────────────────────────────────────────────┐
│  Header（全站通用）                                      │
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
│  PostList    │    │  Article        │    │  Stats       │
│  Sidebar     │    │  TOC            │    │  PostTable   │
│  Footer      │    │  Related        │    │  Editor      │
└──────────────┘    │  Footer         │    │  Preview     │
                    └─────────────────┘    └──────────────┘
```
