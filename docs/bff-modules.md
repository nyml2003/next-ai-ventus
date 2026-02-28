# Ventus BFF 模块映射

本文档说明前端组件与后端 BFF 模块的映射关系。

---

## 问题背景

前端容易混淆**组件名**和**BFF 模块名**：

```typescript
// ❌ 错误：使用前端组件名
request.call({
  page: 'home',
  modules: ['Logo', 'Nav', 'PostList', 'Footer']  // 这是组件名！
})

// ✅ 正确：使用 BFF 模块名
request.call({
  page: 'home',
  modules: ['header', 'postList', 'footer']  // 这是 BFF 模块名！
})
```

---

## BFF 模块列表

### C 端页面模块

| BFF 模块 | 前端组件 | 说明 | 数据内容 |
|---------|---------|------|---------|
| `header` | Header, Logo, Nav | 站点头部 | Logo, 导航链接, 用户信息 |
| `postList` | PostList | 文章列表 | 文章数组, 分页信息 |
| `article` | Article | 文章详情 | 标题, 内容, 元信息, TOC |
| `footer` | Footer | 站点底部 | 版权信息, 备案号 |

### B 端页面模块

| BFF 模块 | 前端组件 | 说明 | 数据内容 |
|---------|---------|------|---------|
| `adminSidebar` | Sidebar | 管理后台侧边栏 | 菜单, 用户信息 |
| `adminFilter` | FilterBar | 筛选选项 | 状态选项, 标签列表 |
| `adminPostList` | PostTable | 文章列表（管理视图）| 文章列表, 统计数字 |
| `editor` | Editor | 编辑器内容 | 文章内容, 版本信息 |
| `editorSettings` | Settings | 编辑器设置 | 所有标签列表 |

### P1 扩展模块（未实现）

| BFF 模块 | 说明 | 优先级 |
|---------|------|--------|
| `hero` | 首页横幅 | P1 |
| `sidebar` | 侧边栏（标签云等） | P1 |
| `toc` | 文章目录 | P1 |
| `related` | 相关文章 | P1 |

---

## 页面与模块映射

### 首页 (`/`)

```typescript
// orchestration.ts
const homeConfig = {
  id: 'home',
  modules: ['header', 'postList', 'footer'],  // ✅ BFF 模块名
  regions: [
    {
      id: 'header',
      type: 'header',
      block: {
        children: [
          { type: 'module', name: 'Logo' },    // 前端组件名
          { type: 'module', name: 'Nav' }      // 前端组件名
        ]
      }
    },
    {
      id: 'content',
      type: 'content',
      block: {
        children: [
          { type: 'module', name: 'PostList' } // 前端组件名
        ]
      }
    }
  ]
}
```

**关键区别**：
- `modules: ['header', ...]` → 告诉后端需要哪些 BFF 模块
- `name: 'Logo'` → 前端自己渲染的组件名

### 文章详情页 (`/post/:slug`)

```typescript
const postConfig = {
  id: 'post',
  modules: ['header', 'article', 'footer'],
  // ...
}
```

### 文章管理页 (`/admin/posts`)

```typescript
const adminPostsConfig = {
  id: 'adminPosts',
  modules: ['adminSidebar', 'adminFilter', 'adminPostList'],
  // ...
}
```

---

## 使用示例

### 前端请求

```typescript
import { createRequest } from '@ventus/request';

const request = createRequest();

// 获取首页数据
async function loadHomePage() {
  const modules = await request.call({
    page: 'home',
    modules: ['header', 'postList', 'footer'],
    params: { page: 1 }
  });
  
  // modules.header.data - 头部数据
  // modules.postList.data - 文章列表
  // modules.footer.data - 底部数据
  return modules;
}

// 获取文章详情页
async function loadPostPage(slug: string) {
  const modules = await request.call({
    page: 'post',
    modules: ['header', 'article', 'footer'],
    params: { slug }
  });
  
  return modules;
}
```

---

## 常见错误

### 错误 1：模块名不存在

```typescript
// ❌ 错误
modules: ['Logo', 'Nav']  // 后端返回 404

// ✅ 正确
modules: ['header']  // header 模块包含 Logo 和 Nav 的数据
```

### 错误 2：模块名拼写错误

```typescript
// ❌ 错误
modules: ['postlist']  // 驼峰错误

// ✅ 正确
modules: ['postList']  // 注意大小写
```

### 错误 3：使用了未实现的模块

```typescript
// ❌ 错误
modules: ['sidebar']  // P1 模块，MVP 未实现

// ✅ 正确
modules: ['postList']  // MVP 已实现
```

---

## 检查清单

### 添加新页面时

- [ ] 确认需要的 BFF 模块名（查阅本文档）
- [ ] 确认模块已实现（MVP/P1）
- [ ] 使用正确的命名（驼峰）
- [ ] 测试请求返回的数据结构

### 添加新模块时

- [ ] 后端注册模块名（`bff/handler.go`）
- [ ] 更新本文档
- [ ] 通知前端模块可用
