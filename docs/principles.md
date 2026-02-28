# Ventus 设计原则

本文档记录项目核心的设计原则，指导后续开发和维护。

---

## 1. BFF 模块设计原则

### 1.1 首屏模块必须走 BFF

**原则**：所有首屏可见的信息展示模块，必须通过 BFF `/api/page` 接口获取数据。

**为什么**：
- 首屏速度决定用户体验
- 统一数据获取方式，避免碎片化
- 便于缓存和性能优化

**判定标准**：
```
模块是否首屏可见？
  ├─ 是 → 必须有对应的 BFF 模块
  │       └─ 是否纯信息展示？
  │         ├─ 是 → BFF 返回完整数据
  │         └─ 否（有交互）→ BFF 返回初始数据，交互前端处理
  │
  └─ 否（二次交互后显示）→ 可不走 BFF
          └─ 如：弹窗、下拉菜单、Toast
```

**示例**：
| 模块 | 是否首屏 | BFF 模块 | 说明 |
|------|---------|----------|------|
| 文章列表 | ✅ | `postList` | 首屏核心内容 |
| 侧边栏标签云 | ✅ | `sidebar` | 首屏可见 |
| 登录表单 | ✅ | 无 | 首屏可见但纯表单，用户输入前端处理 |
| 删除确认框 | ❌ | 无 | 点击后才出现 |

### 1.2 一个页面模块对应一个 BFF 模块

**原则**：产品文档中的每个「模块」尽量对应一个独立的 BFF 模块名。

**命名规范**：
```
C 端模块：camelCase，描述功能
  header, hero, postList, sidebar, article, toc, related

B 端模块：admin 前缀 + 功能
  adminSidebar, adminStats, adminPostList, adminFilter

编辑器相关：editor 前缀
  editor, editorSettings
```

### 1.3 BFF 模块只返回数据，不处理交互状态

**原则**：BFF 返回纯数据，交互状态（选中、展开、loading）由前端管理。

**正确示例**：
```json
// BFF 返回
{
  "postList": {
    "data": {
      "items": [...],
      "pagination": { "page": 1, "total": 100 }
    }
  }
}

// 前端管理
const [selectedIds, setSelectedIds] = useState([])  // 选中状态前端管
```

---

## 2. 模块划分原则

### 2.1 什么时候用 BFF 模块

| 场景 | 处理方式 | 示例 |
|------|---------|------|
| 首屏信息展示 | ✅ BFF 模块 | Header, PostList, Stats |
| 首屏数据依赖 | ✅ BFF 模块 | Sidebar（依赖全局标签）、TOC（依赖文章内容）|
| 用户输入表单 | ❌ 纯前端 | LoginForm, Editor 输入 |
| 二次交互内容 | ❌ 可选 | Modal 内容可前端静态或懒加载 |
| 实时预览 | ❌ 前端处理 | Markdown 预览前端渲染 |

### 2.2 模块边界划分

**原则**：一个 BFF 模块应该是「自包含」的，不依赖其他模块数据。

**例外情况**：
- `toc` 模块依赖 `article` 的内容，但两者可并行执行
- 这种依赖关系由后端处理，前端不感知

**禁止**：
```
❌ 前端先调用模块 A，拿到结果后再调用模块 B

✅ 正确做法：
   前端同时请求 A 和 B
   后端如果 B 依赖 A，自行处理依赖关系
```

---

## 3. 数据流原则

### 3.1 单向数据流

```
BFF 返回数据 → PageStore（只读）→ 组件渲染
                    ↓
            用户交互 → UIStore（可读可写）→ 局部更新
                    ↓
            需要服务端同步 → 调用 API → 刷新 PageStore
```

### 3.2 服务端状态 vs 客户端状态

| 类型 | 存储 | 可变性 | 示例 |
|------|------|--------|------|
| 服务端状态 | PageStore | 只读 | 文章列表、用户信息 |
| 客户端状态 | UIStore | 可读可写 | 主题、侧边栏展开、选中项 |
| 临时状态 | 组件 State | 可读可写 | 表单输入、弹窗显隐 |

### 3.3 数据获取封装

**原则**：页面通过编排系统声明所需 BFF 模块，模块从 Store 读取数据。

```typescript
// 编排配置声明所需 BFF 模块
// pages/home/orchestration.ts
export const homeConfig: PageOrchestrationConfig = {
  id: 'home',
  modules: ['header', 'postList', 'footer'],  // 声明所需 BFF 模块
  regions: [/* ... */]
};

// 模块内读取数据
// modules/PostList.tsx
const { data, loading, error } = useModuleData('postList');
```

**禁止**：
```typescript
// ❌ 不要在组件内直接 fetch
useEffect(() => {
  fetch('/api/page').then(...)  // 错误！
}, []);

// ❌ 不要手动调用 BFF 获取首屏数据
const { data } = usePageData('home', [...]);  // 首屏数据由编排系统获取

// ✅ 模块从 Store 读取数据
const { data } = useModuleData('postList');
```

更多详情参见 [编排系统文档](./orchestration.md)。

---

## 4. API 设计原则

### 4.1 接口分层

```
POST /api/page          # BFF 接口：获取页面数据（首屏必须）
POST /api/login         # 认证接口：登录
POST /api/logout        # 认证接口：登出
POST /api/posts/:id/view # 埋点接口：阅读量统计

# 管理 API（需 JWT）
POST   /api/admin/posts
PUT    /api/admin/posts/:id
DELETE /api/admin/posts/:id
POST   /api/admin/upload
GET    /api/admin/images
```

### 4.2 BFF 请求格式

```typescript
// 请求
{
  "page": "home",                           // 页面标识
  "modules": ["header", "postList"],       // 所需模块列表
  "params": { "page": 1, "tag": "go" }      // 页面参数
}

// 响应
{
  "page": "home",
  "meta": { "title": "...", "description": "..." },
  "modules": {
    "header": { "code": 200, "data": {...} },
    "postList": { "code": 200, "data": {...} },
    "sidebar": { "code": 500, "error": "..." }  // 单个失败不影响其他
  }
}
```

### 4.3 错误处理

**原则**：单个 BFF 模块失败不影响其他模块渲染。

```typescript
// 响应中部分模块失败
{
  "modules": {
    "header": { "code": 200, "data": {...} },     // ✅ 成功
    "postList": { "code": 500, "error": "..." },  // ❌ 失败
    "sidebar": { "code": 200, "data": {...} }     // ✅ 成功
  }
}

// 前端处理
<Header data={modules.header.data} />              // 正常渲染
<ErrorFallback error={modules.postList.error} />   // 显示错误占位
<Sidebar data={modules.sidebar.data} />            // 正常渲染
```

### 4.4 乐观锁（并发控制）

**原则**：编辑操作使用乐观锁防止并发覆盖。

```typescript
// 请求携带 version
PUT /api/admin/posts/:id
{
  "title": "新标题",
  "content": "新内容",
  "version": 3  // 必须携带当前版本号
}

// 冲突响应
{
  "code": 409001,
  "message": "文章已被修改",
  "data": {
    "currentVersion": 4  // 服务端最新版本
  }
}
```

---

## 5. 前后端协作原则

### 5.1 产品文档优先

**原则**：产品文档是唯一的「真相源」，技术文档必须与产品文档保持一致。

**工作流**：
```
产品需求 → 产品文档（pages/*.md）→ 技术实现
                ↓
           发现不一致时
                ↓
         以产品文档为准更新技术文档
```

### 5.2 模块命名对齐

**原则**：产品文档中的模块名 = BFF 模块名 = 前端组件名。

| 产品文档 | BFF 模块 | 前端组件 | 说明 |
|---------|---------|---------|------|
| PostList | `postList` | `<PostList />` | 完全一致 |
| Sidebar | `sidebar` | `<Sidebar />` | 完全一致 |
| 数据统计 | `adminStats` | `<AdminStats />` | B 端加前缀 |

### 5.3 变更同步

**原则**：产品设计变更时，同步更新以下文档：
1. `docs/pages/*.md` - 产品文档
2. `docs/architecture.md` - 架构映射表
3. `docs/server/README.md` - BFF 模块注册表
4. `docs/principles.md` - 如有原则性变化

---

## 6. 性能原则

### 6.1 首屏加载

**原则**：首屏内容必须在一次 BFF 请求内返回。

**限制**：
- 单个页面 BFF 模块数建议不超过 6 个
- 单个模块返回数据量不超过 100KB
- BFF 整体响应时间 P99 < 100ms

### 6.2 并行执行

**原则**：BFF 模块之间无依赖的，必须并行执行。

```go
// 后端实现示例
var wg sync.WaitGroup
results := make(map[string]ModuleResult)

for _, moduleName := range req.Modules {
    wg.Add(1)
    go func(name string) {
        defer wg.Done()
        handler := ModuleRegistry[name]
        data, err := handler(ctx)
        results[name] = ModuleResult{data, err}
    }(moduleName)
}

wg.Wait()
```

### 6.3 缓存策略

| 数据类型 | 缓存位置 | 缓存时间 | 示例 |
|---------|---------|---------|------|
| 站点配置 | 服务端内存 | 长期 | Header、Footer |
| 文章列表 | 服务端内存 | 1 分钟 | PostList |
| 文章详情 | 服务端内存 | 5 分钟 | Article |
| 统计数据 | 服务端内存 | 10 秒 | Stats（实时性要求高）|

---

## 7. 安全原则

### 7.1 鉴权分层

| 接口 | 鉴权要求 | 说明 |
|------|---------|------|
| `GET /api/page` (C 端页面) | 无需鉴权 | 公开访问 |
| `POST /api/page` (B 端页面) | JWT Cookie | 管理后台需要登录 |
| `/api/admin/*` | JWT Cookie | 所有管理 API 需要登录 |
| `/uploads/*` | 公开 | 文章配图需要公开访问 |

### 7.2 输入校验

**原则**：所有用户输入必须在服务端校验。

```go
// 文件路径校验（防目录遍历）
if strings.Contains(path, "..") {
    return error("invalid path")
}

// 图片类型白名单
allowedTypes := []string{"jpg", "jpeg", "png", "webp", "gif"}
```

---

## 8. 演进原则

### 8.1 向后兼容

**原则**：BFF 接口变更必须向后兼容。

```
✅ 允许：新增字段、新增模块、新增可选参数
❌ 禁止：删除字段、重命名字段、改变数据类型
```

### 8.2 渐进增强

**原则**：新功能优先作为独立模块添加，不影响现有模块。

```
新增「热门文章」功能：
  ✅ 新增 `hotPosts` BFF 模块
  ✅ 页面按需引入，不影响现有 `postList`
  
  ❌ 不要在 `postList` 里加一个 `isHot` 字段混用
```

---

## 附录：检查清单

### 新增页面时检查

- [ ] 产品文档（pages/*.md）已更新
- [ ] 编排配置（orchestration.ts）已创建
- [ ] 各模块已标注 BFF 模块名
- [ ] architecture.md 映射表已更新
- [ ] server/README.md 模块注册表已更新
- [ ] 首屏模块都已定义对应的 BFF 模块
- [ ] 非首屏模块已说明原因

### 新增 BFF 模块时检查

- [ ] 模块名符合命名规范
- [ ] 模块功能单一、自包含
- [ ] 已考虑缓存策略
- [ ] 错误处理已定义
- [ ] 已更新模块注册表

### 新增编排模块（Module）时检查

- [ ] 模块已在 `orchestration.ts` 的 `modules` 中注册
- [ ] 模块组件已创建并正确导出
- [ ] 模块使用 `useModuleData` 读取数据
- [ ] 模块处理了 loading 和 error 状态

### 接口变更时检查

- [ ] 是否向后兼容
- [ ] 产品文档是否同步更新
- [ ] 前端调用方是否同步更新

### 错误码变更时

- [ ] 前后端同时修改（同一 PR）
- [ ] 命名语义保持一致（仅风格差异）
- [ ] Code Review 确认两端一致
- [ ] 检查业务代码是否依赖该错误码
