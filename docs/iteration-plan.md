# Ventus 开发迭代计划

## 总体策略

**测试驱动开发（TDD）**：每一阶段都先写测试，后写实现，确保核心代码测试覆盖率 >= 90%。

**分层推进**：
1. **第一阶段**：Go 核心领域层（纯业务逻辑，可独立测试）
2. **第二阶段**：基础设施层（HTTP API、文件存储、BFF模块）
3. **第三阶段**：MVP 前端（React + TypeScript）

---

## 第一阶段：Go 核心领域层（Core Domain）

**目标**：实现不依赖 HTTP 和数据库的可测试核心业务逻辑

**交付物**：
- 文章实体与值对象（Go struct + 方法）
- 内存 Store 实现（Repository 接口）
- 索引管理器
- Markdown/TOC 处理工具
- 测试覆盖率 >= 90%

### 1.1 目录结构

与现有文档对齐：

```
server/
├── internal/
│   ├── domain/               # 领域层 - 纯业务逻辑，无外部依赖
│   │   ├── post.go          # Post 实体
│   │   ├── author.go        # Author 实体
│   │   └── valueobject/     # 值对象
│   │       ├── slug.go      # Slug 验证与生成
│   │       ├── tag.go       # Tag 值对象
│   │       └── status.go    # 文章状态
│   │
│   ├── repository/           # 仓库接口（抽象）
│   │   └── post_repository.go
│   │
│   └── service/              # 领域服务（纯内存实现，无IO）
│       ├── post_service.go      # 文章业务逻辑
│       ├── index_service.go     # 索引管理
│       └── slug_service.go      # Slug 生成与冲突检测
│
└── pkg/                      # 可复用工具（无业务逻辑）
    ├── markdown/            # Markdown 解析与 TOC 生成
    │   ├── parser.go
    │   ├── toc.go
    │   └── sanitize.go
    └── validator/           # 通用校验
```

### 1.2 详细任务

#### Task 1.1: 值对象与实体 (2天)

**Slug 值对象** (`domain/valueobject/slug.go`)：
```go
package valueobject

type Slug struct {
    value string
}

func NewSlug(raw string) (Slug, error)  // 验证并创建
func GenerateFromTitle(title string, existing []string) Slug  // 自动生成
func (s Slug) String() string
func (s Slug) Equals(other Slug) bool
```

**Tag 值对象** (`domain/valueobject/tag.go`)：
```go
type Tag struct {
    name string
}

func NewTag(name string) (Tag, error)
func (t Tag) String() string
```

**Post 实体** (`domain/post.go`)：
```go
type Post struct {
    ID          string
    Title       string
    Slug        valueobject.Slug
    Content     string
    Excerpt     string
    Tags        []valueobject.Tag
    Status      valueobject.Status
    CreatedAt   time.Time
    UpdatedAt   time.Time
    PublishedAt *time.Time
    Version     int  // 乐观锁
    Cover       string
}

// 领域方法
func (p *Post) Publish() error
func (p *Post) Unpublish() error
func (p *Post) UpdateContent(content string)
func (p *Post) UpdateTags(tags []valueobject.Tag)
func (p *Post) GenerateExcerpt(maxLen int) string
func (p *Post) IncrementVersion()
```

**测试要求**（覆盖率 100%）：
- Slug 格式验证（只允许小写字母、数字、连字符）
- Slug 冲突检测与自动编号（hello-world → hello-world-2）
- Tag 去重与排序
- Post 状态转换（draft → published，约束校验）
- Excerpt 生成（从 Markdown 提取纯文本）

#### Task 1.2: 仓库接口与内存实现 (2天)

**仓库接口** (`repository/post_repository.go`)：
```go
package repository

type PostRepository interface {
    FindByID(id string) (*domain.Post, error)
    FindBySlug(slug string) (*domain.Post, error)
    FindAll(opts ListOptions) ([]*domain.Post, int, error)
    FindByTag(tag string) ([]*domain.Post, error)
    Save(post *domain.Post) error
    Delete(id string) error
    Exists(slug string) (bool, error)
    Count(opts CountOptions) (int, error)
}

type ListOptions struct {
    Page       int
    PageSize   int
    Tag        string
    Status     string
    OrderBy    string
}
```

**内存实现** (`repository/memory/post_repository.go`)：
```go
type MemoryPostRepository struct {
    posts     map[string]*domain.Post  // id -> post
    slugIndex map[string]string        // slug -> id
    tagIndex  map[string]map[string]struct{}  // tag -> set(id)
    mu        sync.RWMutex
}

func NewMemoryPostRepository() *MemoryPostRepository
```

**测试要求**：
- 所有 CRUD 操作正确性
- 标签索引自动维护（添加/删除/更新文章时更新 tagIndex）
- Slug 索引一致性
- 分页查询（边界条件）
- 并发安全（读写锁正确性）

#### Task 1.3: 索引管理器 (1天)

**索引结构** (`service/index_service.go`)：
```go
type Index struct {
    PostIDs      []string                 // 按时间排序的文章ID列表
    SlugToID     map[string]string        // slug -> id
    TagToIDs     map[string][]string      // tag -> ids（按时间排序）
    DateToIDs    map[string][]string      // "2024-06" -> ids
}

type IndexService struct {
    repo PostRepository
}

func (s *IndexService) BuildIndex() (*Index, error)
func (s *IndexService) AddToIndex(index *Index, post *domain.Post)
func (s *IndexService) SearchByTag(index *Index, tag string) []string
func (s *IndexService) SearchByDateRange(index *Index, start, end time.Time) []string
```

**测试要求**：
- 索引构建性能（1000 篇文章 < 50ms）
- 标签搜索准确性
- 日期范围搜索正确性

#### Task 1.4: Markdown 处理 (1天)

**位置**: `pkg/markdown/`

```go
package markdown

type TOCItem struct {
    Level    int
    Text     string
    Anchor   string
    Children []*TOCItem
}

type ParsedResult struct {
    HTML            string
    TOC             []*TOCItem
    Excerpt         string
    WordCount       int
    ReadTimeMinutes int
}

func Parse(content string) (*ParsedResult, error)
func ExtractTOC(content string) ([]*TOCItem, error)
func GenerateAnchor(headingText string) string
func SanitizeHTML(html string) string
```

**测试要求**：
- 各种 Markdown 语法解析（标题、列表、代码块、表格等）
- TOC 层级正确性（h1-h6）
- 锚点生成唯一性（处理特殊字符、中文）
- XSS 防护（script 标签过滤）

#### Task 1.5: 领域服务层 (2天)

**文章服务** (`service/post_service.go`)：
```go
type PostService struct {
    repo        repository.PostRepository
    slugService *SlugService
}

type CreatePostInput struct {
    Title   string
    Content string
    Tags    []string
    Status  string
}

func (s *PostService) CreatePost(input CreatePostInput) (*domain.Post, error)
func (s *PostService) UpdatePost(id string, input UpdatePostInput, version int) (*domain.Post, error)  // 乐观锁
func (s *PostService) DeletePost(id string) error
func (s *PostService) GetPost(id string) (*domain.Post, error)
func (s *PostService) ListPosts(opts ListOptions) ([]*domain.Post, int, error)
func (s *PostService) PublishPost(id string) (*domain.Post, error)
```

**Slug 服务** (`service/slug_service.go`)：
```go
type SlugService struct {
    repo repository.PostRepository
}

func (s *SlugService) GenerateUniqueSlug(title string) (valueobject.Slug, error)
func (s *SlugService) CheckConflict(slug string, excludeID string) (bool, error)
```

**测试要求**：
- 完整 CRUD 流程测试
- 乐观锁冲突处理（version 不匹配时返回特定错误）
- Slug 自动生成与冲突解决
- 标签关联正确性

### 1.3 阶段交付检查清单

- [ ] 所有值对象单元测试（覆盖率 100%）
- [ ] Post 实体单元测试（状态转换、Excerpt 生成）
- [ ] 内存仓库单元测试（含并发测试）
- [ ] 索引服务单元测试
- [ ] Markdown 处理器单元测试
- [ ] 领域服务集成测试
- [ ] 整体测试覆盖率 >= 90%

---

## 第二阶段：基础设施与 API（Infrastructure & API）

**目标**：实现 HTTP API 和文件存储，可独立运行后端服务

**交付物**：
- 文件系统 Store 实现（替换内存实现）
- BFF 架构实现（Go + Gin）
- JWT 认证
- 图片上传
- 集成测试

### 2.1 目录结构

```
server/
├── internal/
│   ├── domain/              # 第一阶段已完成（复用）
│   ├── repository/          # 新增文件实现
│   │   └── file/
│   │       └── post_repository.go
│   │
│   ├── service/             # 第一阶段已完成（复用）
│   │
│   ├── store/               # 文件存储层（与现有文档对齐）
│   │   ├── file_store.go    # 文件读写抽象
│   │   └── index_loader.go  # 索引加载器
│   │
│   ├── interfaces/          # 接口适配层
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── post_handler.go
│   │   │   │   └── upload_handler.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go      # JWT 验证
│   │   │   │   ├── cors.go
│   │   │   │   └── logger.go
│   │   │   └── router.go
│   │   │
│   │   └── bff/             # BFF 模块实现
│   │       ├── handler.go   # POST /api/page 处理
│   │       └── modules/
│   │           ├── header.go          # MVP
│   │           ├── post_list.go       # MVP
│   │           ├── footer.go          # MVP
│   │           ├── article.go         # MVP
│   │           ├── admin_sidebar.go   # MVP
│   │           ├── admin_filter.go    # MVP
│   │           ├── admin_post_list.go # MVP
│   │           ├── editor.go          # MVP
│   │           ├── editor_settings.go # MVP
│   │           ├── hero.go            # P1
│   │           ├── sidebar.go         # P1
│   │           ├── toc.go             # P1
│   │           ├── related.go         # P1
│   │           ├── admin_stats.go     # P1
│   │           └── image_list.go      # P1
│   │
│   └── config/
│       └── config.go
│
├── storage/uploads/         # 上传文件目录（.gitignore）
│
├── content/                 # 文章内容目录
│   ├── posts/
│   └── authors/
│
└── tests/
    └── integration/         # 集成测试
```

### 2.2 详细任务

#### Task 2.1: 文件存储实现 (2天)

**数据目录结构**：
```
content/
├── posts/
│   └── 2024-06-hello-world/
│       ├── meta.json
│       └── content.md
└── authors/
    └── nyml.json
```

**meta.json 格式**：
```json
{
  "id": "2024-06-hello-world",
  "title": "Hello World",
  "slug": "hello-world",
  "date": "2024-06-15",
  "tags": ["go", "architecture"],
  "status": "published",
  "version": 1,
  "cover": "/uploads/2024/06/abc123.jpg"
}
```

**FilePostRepository** (`repository/file/post_repository.go`)：
```go
type FilePostRepository struct {
    basePath string
    index    *service.Index
    mu       sync.RWMutex
}

func (r *FilePostRepository) LoadIndex() error  // 启动时加载
func (r *FilePostRepository) reloadIndex() error // 热更新
func (r *FilePostRepository) Save(post *domain.Post) error
func (r *FilePostRepository) FindByID(id string) (*domain.Post, error)
// ... 其他接口方法
```

**测试要求**：
- 文件读写正确性
- 索引加载与热更新
- 并发读写安全

#### Task 2.2: HTTP 基础框架 (1天)

**Router** (`interfaces/http/router.go`)：
```go
func SetupRouter(
    postService *service.PostService,
    authService *service.AuthService,
    bffHandler *bff.Handler,
) *gin.Engine {
    r := gin.Default()
    
    // 公开 API
    r.POST("/api/login", authHandler.Login)
    r.POST("/api/page", bffHandler.Handle)  // BFF 统一接口
    r.POST("/api/posts/:id/view", postHandler.RecordView)  // 埋点
    
    // 需认证 API
    admin := r.Group("/api/admin")
    admin.Use(middleware.JWTAuth())
    {
        admin.POST("/posts", postHandler.Create)
        admin.PUT("/api/admin/posts/:id", postHandler.Update)
        admin.DELETE("/api/admin/posts/:id", postHandler.Delete)
        admin.POST("/upload", uploadHandler.Upload)
    }
    
    return r
}
```

#### Task 2.3: BFF 模块实现 (3天)

**MVP BFF 模块**：

| 模块 | 页面 | 功能 |
|------|------|------|
| `header` | 通用 | 站点 Logo、导航 |
| `footer` | 通用 | 版权信息 |
| `postList` | 首页 | 文章列表（分页、标签筛选）|
| `article` | 文章页 | 文章详情（含 Markdown 渲染）|
| `adminSidebar` | 管理页 | 管理后台导航 |
| `adminFilter` | 文章管理 | 筛选选项：状态、标签列表 |
| `adminPostList` | 文章管理 | 文章列表（含统计数字）|
| `editor` | 编辑器 | 文章内容获取/保存 |
| `editorSettings` | 编辑器 | 所有标签列表（供选择）|

**P1 扩展模块**：`hero`, `sidebar`, `toc`, `related`, `adminStats`, `imageList`, `imageFilter`

**BFF Handler 结构** (`interfaces/bff/handler.go`)：
```go
type ModuleHandler func(ctx *ModuleContext) (interface{}, error)

var ModuleRegistry = map[string]ModuleHandler{
    "header":        modules.HandleHeader,
    "postList":      modules.HandlePostList,
    "article":       modules.HandleArticle,
    // ...
}

type Handler struct {
    services *Services
}

func (h *Handler) Handle(c *gin.Context) {
    var req PageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{...})
        return
    }
    
    // 并行执行模块
    results := h.executeModules(req.Page, req.Modules, req.Params)
    
    c.JSON(200, PageResponse{
        Page:    req.Page,
        Modules: results,
    })
}
```

#### Task 2.4: JWT 认证 (1天)

**位置**: `interfaces/http/middleware/auth.go`

```go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token, err := c.Cookie("token")
        if err != nil {
            c.AbortWithStatus(401)
            return
        }
        
        claims, err := validateToken(token)
        if err != nil {
            c.AbortWithStatus(401)
            return
        }
        
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

#### Task 2.5: 图片上传 (1天)

**位置**: `interfaces/http/handlers/upload_handler.go`

```go
func (h *UploadHandler) Upload(c *gin.Context) {
    // 1. JWT 验证（由中间件处理）
    // 2. 获取文件
    file, header, err := c.Request.FormFile("file")
    // 3. 校验文件类型（jpg/png/webp/gif）
    // 4. 校验文件大小（< 5MB）
    // 5. 生成文件名 {hash}.{ext}
    // 6. 保存到 storage/uploads/2024/06/
    // 7. 返回访问 URL /uploads/2024/06/xxx.jpg
}
```

### 2.3 阶段交付检查清单

- [ ] 文件存储集成测试
- [ ] API 端点测试（使用 httptest）
- [ ] BFF 模块测试（各模块独立测试 + 组合测试）
- [ ] JWT 认证测试
- [ ] 图片上传测试
- [ ] 整体测试覆盖率 >= 80%

---

## 第三阶段：MVP 前端（Frontend MVP）

**目标**：实现可使用的最小可用产品

**交付物（5 个页面）**：
- 登录页 (`/login`)
- 首页 (`/`) - 极简布局，无 Hero/Sidebar
- 文章详情页 (`/post/:slug`) - 无 TOC/Related
- 文章管理 (`/admin/posts`) - 含顶部统计
- 文章编辑器 (`/admin/editor`)

**不做**：
- `/admin` 仪表盘页面（直接进入文章管理）
- `/admin/images` 图片管理页面（P1 再做）
- Hero、Sidebar、TOC、Related 模块（P1 再做）

### 3.1 目录结构

与现有文档对齐：

```
frontend/
├── packages/                    # 共享 npm 包
│   ├── ui/                     # @next-ai-ventus/ui
│   │   ├── components/
│   │   │   ├── Button/
│   │   │   ├── Layout/
│   │   │   ├── PostCard/
│   │   │   ├── Tag/
│   │   │   ├── Modal/
│   │   │   └── Toast/
│   │   └── theme/              # 主题系统
│   │
│   ├── utils/                  # @next-ai-ventus/utils
│   │   ├── bff.ts             # BFF 请求封装
│   │   ├── auth.ts            # 认证工具
│   │   └── storage.ts         # Cookie/localStorage
│   │
│   ├── types/                  # @next-ai-ventus/types
│   │   └── api.ts             # API 类型定义
│   │
│   └── markdown/               # @next-ai-ventus/markdown
│       ├── renderer.tsx
│       └── highlight.ts
│
├── pages/                      # 页面入口（MPA）
│   ├── home/
│   │   ├── index.html
│   │   └── main.tsx
│   │
│   ├── post/                   # 文章详情
│   │   └── main.tsx            # 路由 /post/:slug
│   │
│   ├── login/
│   │   ├── index.html
│   │   └── main.tsx
│   │
│   ├── admin-posts/            # 文章管理（MVP 入口）
│   │   ├── index.html
│   │   └── main.tsx
│   │
│   └── admin-editor/           # 文章编辑器
│       ├── index.html
│       └── main.tsx
│
└── shell/                      # 统一构建配置
    └── vite.config.ts
```

### 3.2 详细任务

#### Task 3.1: 基础组件库 (2天)

**@next-ai-ventus/ui** 必要组件：
- Button、Input、TextArea、Select
- Layout（Header、Footer、Container）**Sidebar P1 再做**
- PostCard（文章卡片：标题、摘要、标签、时间）
- Tag（标签组件）
- Modal、Toast（反馈组件）
- Skeleton（加载占位）

**主题系统**：
```typescript
// CSS 变量
:root {
  --color-bg-primary: #ffffff;
  --color-text-primary: #171717;
  --color-primary: #3b82f6;
  // ...
}
```

#### Task 3.2: BFF 客户端封装 (1天)

**@next-ai-ventus/utils/bff.ts**：
```typescript
interface PageRequest {
  page: string;
  modules: string[];
  params?: Record<string, any>;
}

interface PageResponse {
  page: string;
  meta: { title: string; description?: string };
  modules: Record<string, ModuleResult>;
}

async function fetchPageData(req: PageRequest): Promise<PageResponse>;

// React Hook (使用 SWR)
function usePageData(
  page: string, 
  modules: string[], 
  params?: object
): { data; error; isLoading };
```

#### Task 3.3: 页面实现 (4天)

**1. 登录页** (`/login`) - 0.5天
```
模块: LoginForm (纯前端，无 BFF)
- 用户名/密码输入
- JWT 存储到 Cookie
- 登录成功后跳转 /admin
```

**2. 首页** (`/`) - 1天
```
BFF 模块: header, postList, footer
- Header: Logo、导航链接
- PostList: 文章卡片列表（分页）
- 标签筛选: URL query ?tag=xxx
- Footer: 版权信息

布局: 单列居中，无 Sidebar，无 Hero
```

**3. 文章详情页** (`/post/:slug`) - 1天
```
BFF 模块: header, article, footer
- Article: 标题、元信息、Markdown 渲染、代码高亮
- 代码高亮: highlight.js

不做: TOC（P1）、Related（P1）
```

**4. 文章管理** (`/admin/posts`) - 1天
```
BFF 模块: adminSidebar, adminFilter, adminPostList
- 顶部统计卡片: 文章总数、已发布、草稿数（合并 adminStats）
- FilterBar: 状态筛选（全部/已发布/草稿）、标签筛选
- PostTable: 表格展示（标题、状态、时间、操作）
- 删除确认弹窗

说明: MVP 不做独立仪表盘，统计信息放在文章管理页顶部
```

**6. 文章编辑器** (`/admin/editor`) - 1天
```
布局: 分栏（左编辑 / 右预览）
BFF 模块（编辑时）: editor, editorSettings

- 左侧: Markdown 编辑器（textarea）
- 右侧: 实时预览（Markdown 渲染）
- 顶部工具栏: 保存、发布、返回
- 侧边设置: 标题、标签、状态、封面图

新建: /admin/editor
编辑: /admin/editor?id=xxx
```

### 3.3 测试策略（BFF 前后端都要测）

**前端单元测试（Jest）**：
- 组件渲染测试（Button、PostCard、Modal 等）
- Hook 测试（usePageData、useAuth）
- 工具函数测试（bff-client、auth、storage）

**前端集成测试（Testing Library）**：
- 页面渲染流程（挂载 → 请求 BFF → 渲染模块）
- 用户交互流程（点击标签筛选、表单提交）

**BFF 端到端测试（前后端联调）**：
- 每个页面的 BFF 数据流：前端请求 → 后端组装 → 前端渲染
- 错误处理测试（模块失败时页面降级）

**交付检查清单**：
- [ ] 所有页面可正常访问
- [ ] 文章 CRUD 流程完整
- [ ] 登录/登出功能正常
- [ ] Markdown 渲染正确
- [ ] 代码高亮生效
- [ ] **前端单元测试覆盖率 >= 70%**
- [ ] **BFF 端到端测试通过（5 个页面 × 核心场景）**

---

## 第四阶段：部署与优化（Deployment & Polish）

**目标**：可部署的生产环境

### 4.1 详细任务

#### Task 4.1: Docker 化 (1天)

```dockerfile
# Dockerfile
# 多阶段构建
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/ .
RUN npm install -g pnpm && pnpm install && pnpm build

FROM golang:1.21-alpine AS backend-builder
WORKDIR /app/server
COPY server/ .
RUN go build -o ventus-server cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
COPY --from=backend-builder /app/server/ventus-server .
EXPOSE 8080
CMD ["./ventus-server"]
```

```yaml
# docker-compose.yml
version: '3'
services:
  ventus:
    build: .
    ports:
      - "3000:8080"
    volumes:
      - ./content:/app/content
      - ./storage:/app/storage
```

#### Task 4.2: Nginx 配置 (0.5天)

```nginx
server {
    listen 80;
    
    # 静态资源（前端页面）
    location / {
        root /var/www/ventus;
        try_files $uri $uri/ /index.html;
        expires 1d;
    }
    
    # API 代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
    }
    
    # 上传文件
    location /uploads/ {
        alias /app/storage/uploads/;
        expires 30d;
        # 防盗链配置
        valid_referers none blocked server_names;
        if ($invalid_referer) {
            return 403;
        }
    }
    
    # Gzip
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;
}
```

#### Task 4.3: 部署验证 (0.5天)

- 端到端冒烟测试（登录 → 创建文章 → 查看文章）
- 文件持久化验证（重启容器后数据仍在）
- 图片访问验证

**P1 再做**：性能优化（代码分割、缓存策略、CDN）

---

## 里程碑与时间表

| 阶段 | 里程碑 | 预计工期 | 累计工期 |
|------|--------|----------|----------|
| Phase 1 | Go 核心领域层完成，测试覆盖率 90%+ | 8 天 | 8 天 |
| Phase 2 | Go + Gin 后端 API 可用，可独立测试 | 7 天 | 15 天 |
| Phase 3 | MVP 前端完成（5 页面）+ 测试 | **6 天** | **21 天** |
| Phase 4 | 可部署的生产环境 | 2 天 | **23 天** |

**总计约 4-5 周**

### 各阶段测试要求

| 阶段 | 后端测试 | 前端测试 | 集成测试 |
|------|---------|---------|---------|
| Phase 1 | 单元测试 90%+ | - | - |
| Phase 2 | API 测试 80%+ | - | BFF 模块测试 |
| Phase 3 | - | 单元测试 70%+ | **E2E 测试（关键路径）** |
| Phase 4 | 部署验证 | 部署验证 | 端到端冒烟测试 |

### MVP 范围确认

**做（P0）**：
- ✅ 首页：Header + PostList + Footer（单列布局）
- ✅ 文章页：Header + Article + Footer（无 TOC）
- ✅ 登录页：纯前端表单
- ✅ 文章管理：Sidebar + Stats（合并）+ Filter + Table
- ✅ 文章编辑器：Markdown 编辑器 + 预览
- ✅ 9 个 BFF 模块：header, footer, postList, article, adminSidebar, adminFilter, adminPostList, editor, editorSettings

**不做（P1）**：
- ❌ Hero 模块
- ❌ Sidebar 模块（C 端侧边栏）
- ❌ TOC 模块（文章目录）
- ❌ Related 模块（相关文章）
- ❌ adminStats 独立模块（合并到 adminPostList）
- ❌ /admin 仪表盘页面
- ❌ /admin/images 图片管理页面

---

## 下一步行动

1. **创建第一阶段目录结构** (`server/internal/domain/`, `server/repository/`, `server/service/`)
2. **编写第一个测试** (`domain/valueobject/slug_test.go`)
3. **开始 Task 1.1 开发**

---

## 附录：测试策略

### 单元测试（Go testing）

```go
// domain/valueobject/slug_test.go
func TestSlugCreate(t *testing.T) {
    tests := []struct {
        name    string
        raw     string
        wantErr bool
    }{
        {"valid", "hello-world", false},
        {"valid with numbers", "hello-world-123", false},
        {"invalid chars", "hello@world", true},
        {"uppercase", "Hello-World", true}, // 必须小写
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewSlug(tt.raw)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewSlug() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 集成测试

```go
// tests/integration/post_api_test.go
func TestCreatePost(t *testing.T) {
    router := setupTestRouter()
    
    body := `{"title":"Test","content":"# Hello"}`
    req := httptest.NewRequest("POST", "/api/admin/posts", strings.NewReader(body))
    req.Header.Set("Cookie", "token="+testToken)
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 201, w.Code)
}
```

### 前端测试（Jest + Testing Library）

```typescript
// packages/ui/components/PostCard/PostCard.test.tsx
import { render, screen } from '@testing-library/react';
import { PostCard } from './PostCard';

describe('PostCard', () => {
  const mockPost = {
    id: '1',
    title: '测试文章',
    excerpt: '这是摘要',
    date: '2024-06-15',
    tags: [{ name: 'go' }],
    href: '/post/test'
  };

  it('renders post title', () => {
    render(<PostCard post={mockPost} />);
    expect(screen.getByText('测试文章')).toBeInTheDocument();
  });

  it('renders tags', () => {
    render(<PostCard post={mockPost} />);
    expect(screen.getByText('go')).toBeInTheDocument();
  });
});
```

```typescript
// utils/bff-client.test.ts
import { fetchPageData } from './bff-client';

describe('bff-client', () => {
  it('fetches page data correctly', async () => {
    const result = await fetchPageData({
      page: 'home',
      modules: ['header', 'postList']
    });
    
    expect(result.page).toBe('home');
    expect(result.modules.header).toBeDefined();
    expect(result.modules.postList).toBeDefined();
  });
});
```

### BFF 端到端测试（Playwright 或 Cypress）

```typescript
// e2e/home.spec.ts
import { test, expect } from '@playwright/test';

test('home page loads with post list', async ({ page }) => {
  await page.goto('/');
  
  // 验证 BFF 数据已渲染
  await expect(page.locator('[data-testid="post-list"]')).toBeVisible();
  await expect(page.locator('[data-testid="post-card"]').first()).toBeVisible();
});

test('tag filter works', async ({ page }) => {
  await page.goto('/');
  
  // 点击标签
  await page.click('[data-testid="tag-go"]');
  
  // 验证 URL 变化
  await expect(page).toHaveURL(/\?tag=go/);
  
  // 验证列表已筛选
  await expect(page.locator('[data-testid="filter-info"]')).toContainText('go');
});
```

### 测试覆盖率检查

```bash
# Go 测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 前端测试覆盖率
pnpm test:coverage

# E2E 测试
pnpm test:e2e
```
