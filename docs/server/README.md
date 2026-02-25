# 后端详细设计

## 目录

- [数据存储](#数据存储)
- [Store 层](#store-层)
- [API 设计](#api-设计)
- [图片上传](#图片上传)
- [模块机制](#模块机制)
- [安全设计](#安全设计)
- [错误处理](#错误处理)
- [日志规范](#日志规范)

---

## 数据存储

**有后端，无数据库，文件即数据**

### 目录结构

```
content/                          # 项目根目录下
├── posts/
│   └── 2024-06-hello-world/      # 文章目录：{YYYY-MM}-{slug}
│       ├── meta.json             # 元数据
│       └── content.md            # Markdown 正文
└── authors/
    └── nyml.json                 # 作者信息
```

**目录命名规则**：
- 格式：`{YYYY-MM}-{slug}`，如 `2024-06-hello-world`
- `slug`：URL 短链接，全局唯一，只能包含字母、数字、连字符
- slug 冲突时拒绝创建/更新，提示用户修改

### meta.json 格式

```json
{
  "id": "2024-06-hello-world",
  "title": "Hello World",
  "slug": "hello-world",
  "date": "2024-06-15",
  "tags": ["go", "架构"],
  "status": "published",
  "version": 1,
  "cover": "/uploads/2024/06/abc123.jpg"
}
```

**字段说明：**
- `id`: 文章唯一标识（目录名）
- `title`: 文章标题
- `slug`: URL 友好的短链接（如 `hello-world`）
- `date`: 发布日期（ISO 8601 格式）
- `tags`: 标签数组
- `status`: 文章状态（`draft` 草稿 / `published` 已发布）
- `version`: 乐观锁版本号（编辑时递增，防并发覆盖）
- `cover`: 封面图路径（可选）

### 内存索引

```go
type Index struct {
    Posts   []Post            // 文章列表（按时间排序）
    SlugMap map[string]string // slug -> id
    TagMap  map[string][]string // tag -> ids
    mu      sync.RWMutex
}

// 内存中的文章结构（包含运行时数据）
type Post struct {
    Meta    PostMeta  // meta.json 内容
    Content string    // content.md 内容
    Views   int       // 阅读量（运行时从埋点统计）
}

type PostMeta struct {
    ID      string   `json:"id"`
    Title   string   `json:"title"`
    Slug    string   `json:"slug"`
    Date    string   `json:"date"`
    Tags    []string `json:"tags"`
    Status  string   `json:"status"`
    Version int      `json:"version"`
    Cover   string   `json:"cover,omitempty"`
}
```

---

## Store 层

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

---

## API 设计

### BFF 统一接口（主要）

```
POST /api/page               # 页面数据聚合（列表、详情、各页面都用这个）
```

### 独立接口（特定场景）

```
POST /api/posts/:id/view     # 阅读统计（纯埋点，不返回页面数据）
POST /api/admin/upload       # 图片上传（文件上传）
```

### 管理 API（需 JWT 认证）

**认证接口**
```
POST /api/login              # 登录，返回 JWT token
POST /api/logout             # 登出，清除 Cookie
```

**文章管理**
```
POST   /api/admin/posts      # 创建文章
PUT    /api/admin/posts/:id  # 更新文章（需携带 version 乐观锁）
DELETE /api/admin/posts/:id  # 删除文章
```

**图片管理**
```
GET    /api/admin/images     # 图片列表（支持筛选、分页）
DELETE /api/admin/images/:id # 删除单张图片
POST   /api/admin/images/batch-delete  # 批量删除图片
```

**乐观锁说明**
- 编辑文章时，请求 body 需包含 `version` 字段
- 冲突时返回 409001，响应包含当前最新版本号

### BFF 请求示例

```http
POST /api/page

{
  "page": "home",
  "modules": ["header", "hero", "postList", "sidebar", "footer"],
  "params": { "page": 1, "tag": "go" }
}
```

---

## 图片上传

### 存储结构

```
server/storage/uploads/           # 物理存储在服务器
├── 2024/
│   └── 06/
│       └── abc123.jpg            # 重命名后的文件名
```

### 上传流程

1. JWT 鉴权（只有登录用户可上传）
2. 校验文件类型（白名单：jpg/png/webp/gif）
3. 校验文件大小（< 5MB）
4. 生成文件名 `{hash}.{ext}`
5. 保存到 `server/storage/uploads/2024/06/`
6. 返回访问 URL：`/uploads/2024/06/abc123.jpg`

### 访问流程

- **公开图片**（文章配图）：直接访问，无需鉴权
- **管理图片**（图片管理列表）：需要 JWT 鉴权
- **防盗链**：通过 Referer 或 Nginx 配置限制外部引用

---

## 模块机制

### 模块接口

```go
type ModuleContext struct {
    Page   string                 // 页面名
    Params map[string]interface{} // 请求参数
    Cache  *Cache                 // 缓存实例
}

type ModuleHandler func(ctx *ModuleContext) (interface{}, error)
```

### 模块注册表

```go
var ModuleRegistry = map[string]ModuleHandler{
    // ═══════════════════════════════════════════════════════════
    // C 端页面模块（首屏必须走 BFF）
    // ═══════════════════════════════════════════════════════════
    
    // 首页模块
    "header":     modules.HandleHeader,      // 站点 Logo、导航链接
    "hero":       modules.HandleHero,        // 站点标题、简介（可关闭）
    "postList":   modules.HandlePostList,    // 文章列表（支持分页、标签筛选）
    "sidebar":    modules.HandleSidebar,     // 博主简介、热门文章、标签云
    "footer":     modules.HandleFooter,      // 版权信息、备案号
    
    // 文章详情页模块
    "article":    modules.HandleArticle,     // 完整文章：标题+元信息+正文（首屏）
    "toc":        modules.HandleTOC,         // 目录导航（从 article 内容生成）
    "related":    modules.HandleRelated,     // 相关文章推荐（基于标签匹配）
    
    // ═══════════════════════════════════════════════════════════
    // B 端页面模块（管理后台，首屏必须走 BFF）
    // ═══════════════════════════════════════════════════════════
    
    // 管理首页模块
    "adminSidebar":   modules.HandleAdminSidebar,   // 博主头像、名称、导航菜单
    "adminStats":     modules.HandleAdminStats,     // 文章数、阅读量、标签数等统计
    "recentPosts":    modules.HandleRecentPosts,    // 最近编辑的 5 篇文章
    
    // 文章管理页模块
    "adminPostList":  modules.HandleAdminPostList,  // 管理端文章列表（支持筛选、搜索）
    "adminFilter":    modules.HandleAdminFilter,    // 筛选数据：状态选项、标签列表
    
    // 文章编辑页模块
    "editor":         modules.HandleEditor,         // 获取/保存文章内容
    "editorSettings": modules.HandleEditorSettings, // 编辑器设置数据：所有标签列表
    
    // 图片管理页模块
    "imageList":      modules.HandleImageList,      // 图片列表（支持时间筛选、搜索）
    "imageFilter":    modules.HandleImageFilter,    // 图片筛选选项：时间范围
}

// ═══════════════════════════════════════════════════════════
// 独立接口处理器（非 BFF 模块，纯行为接口）
// ═══════════════════════════════════════════════════════════
// POST /api/posts/:id/view -> HandlePostView()  // 阅读量统计埋点（无需返回页面数据）
```

### 执行流程

1. 接收请求，解析模块列表
2. 检查模块依赖关系，确定执行顺序
3. 并行执行各模块 Handler
4. 聚合结果，组装响应
5. 返回统一格式

---

## 安全设计

### 预留点

| 层面 | 预留方式 | 后期可接入 |
|------|---------|-----------|
| XSS | Markdown 渲染器接口预留 sanitize 钩子 | DOMPurify |
| CSRF | Middleware 链式结构 | CSRF Token 中间件 |
| CSP | Nginx 配置模板预留 header 位置 | add_header Content-Security-Policy |
| 注入 | Store 层文件路径校验接口 | 更严格的 path validate |
| 审计 | Logger 中间件预留字段 | 审计日志模块 |

### 当前最小防护

- 文件路径校验（禁止 `../`）
- 图片类型白名单（jpg/png/webp/gif）
- JWT 基础鉴权

---

## 错误处理

### 错误码规范

| 错误码 | 含义 | HTTP 状态 |
|--------|------|----------|
| 0 | 成功 | 200 |
| 400001 | 参数错误 | 400 |
| 401001 | 未登录 | 401 |
| 403001 | 无权限 | 403 |
| 404001 | 资源不存在 | 404 |
| 409001 | 版本冲突（乐观锁） | 409 |
| 500001 | 模块执行失败 | 200（body 内 code） |
| 500999 | 内部错误 | 500 |

### 错误响应格式

```json
{
  "code": 404001,
  "message": "文章不存在",
  "data": {
    "slug": "not-exist"
  }
}
```

---

## 日志规范

### 日志格式

```json
{
  "time": "2024-06-15T10:30:00Z",
  "level": "info",
  "trace_id": "uuid",
  "module": "postList",
  "duration_ms": 15,
  "status": 200,
  "path": "/api/page",
  "ip": "..."
}
```

### 日志级别

- **INFO**：正常请求（采样 1%）
- **WARN**：慢请求（> 500ms）、模块降级
- **ERROR**：接口异常、文件读写失败

### 日志切分

按天切分，保留 7 天。
