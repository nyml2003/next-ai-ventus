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

| 页面 | 文档 | MVP 核心模块 | P1 扩展模块 | BFF 模块（MVP）|
|------|------|-------------|-------------|---------------|
| 首页 | [home.md](./home.md) | Header, PostList, Footer | Hero, Sidebar | `header`, `postList`, `footer` |
| 文章详情 | [post.md](./post.md) | Header, Article, Footer | TOC, Related | `header`, `article`, `footer` |
| 登录页 | [login.md](./login.md) | LoginForm | - | 无 |

### B 端页面

| 页面 | 文档 | MVP/P1 | 核心模块 | BFF 模块 |
|------|------|--------|---------|----------|
| 管理首页 | [admin-home.md](./admin-home.md) | **P1** | Sidebar, Stats, QuickActions, RecentPosts | `adminSidebar`, `adminStats`, `recentPosts` |
| 文章管理 | [admin-posts.md](./admin-posts.md) | **MVP** | Sidebar, FilterBar, PostTable, Pagination, Stats | `adminSidebar`, `adminFilter`, `adminPostList` |
| 文章编辑 | [admin-editor.md](./admin-editor.md) | **MVP** | Toolbar, Editor, Preview, Settings | `editor`, `editorSettings` |
| 图片管理 | [admin-images.md](./admin-images.md) | **P1** | Sidebar, FilterBar, ImageGrid, Pagination | `adminSidebar`, `imageFilter`, `imageList` |

---

## 模块复用关系

### MVP 版本模块

```
┌─────────────────────────────────────────────────────────┐
│  Header（全站通用）- BFF: header                        │
│  - Logo                                                 │
│  - 导航链接                                             │
│  - 登录入口（B 端显示）                                  │
└─────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐    ┌─────────────────┐    ┌──────────────┐
│   C 端首页    │    │   文章详情页     │    │   管理后台    │
│              │    │                 │    │              │
│  PostList    │    │  Article        │    │  Sidebar     │
│(BFF:postList)│    │(BFF:article)    │    │(BFF:adminSidebar)│
│  Footer      │    │  Footer         │    │  PostTable   │
│(BFF:footer)  │    │(BFF:footer)     │    │(BFF:adminPostList)│
└──────────────┘    └─────────────────┘    │  Editor      │
                                           │(BFF:editor)  │
                                           └──────────────┘
```

### P1 扩展模块

| 模块 | 所属页面 | BFF 模块 | 说明 |
|------|---------|----------|------|
| Hero | 首页 | `hero` | 首屏横幅 |
| Sidebar | 首页 | `sidebar` | 博主简介、标签云 |
| TOC | 文章详情 | `toc` | 目录导航 |
| Related | 文章详情 | `related` | 相关文章 |
| Stats | 管理首页 | `adminStats` | 统计数据 |

**说明**：括号内标注的是对应的 BFF 模块名，首屏模块必须走 BFF 获取数据。

---

## 编排系统集成

页面通过 [编排系统](../orchestration.md) 声明式组装模块：

```typescript
// pages/home/orchestration.ts
export const homeConfig: PageOrchestrationConfig = {
  id: 'home',
  modules: ['header', 'postList', 'footer'],  // 声明所需 BFF 模块
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
      block: {
        type: 'block',
        flexDirection: 'column',
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
        children: [
          { type: 'module', name: 'Footer' }
        ]
      }
    }
  ]
};
```

编排系统自动：
1. 请求声明的 BFF 模块数据
2. 将数据存入 PageStore
3. 根据配置渲染布局结构
4. 模块通过 `useModuleData` 从 Store 读取数据
