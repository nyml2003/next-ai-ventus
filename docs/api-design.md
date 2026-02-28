# Ventus API 契约设计

本文档描述 Ventus 的 API 契约优先架构，使用 OpenAPI 作为单一真相源，生成前后端代码。

---

## 1. 核心思想

**契约优先（Contract First）**

```
┌─────────────────────────────────────────────────────────────┐
│                    OpenAPI 3.0 (api.yml)                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  Schemas        │  │  Paths          │  │  Tags       │ │
│  │  - Post         │  │  - /api/page    │  │  - BFF      │ │
│  │  - Author       │  │  - /api/posts   │  │  - Admin    │ │
│  │  - Error        │  │  - /api/login   │  │  - Public   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└──────────────────────────────┬──────────────────────────────┘
                               │
           ┌───────────────────┴───────────────────┐
           ▼                                       ▼
┌──────────────────────────────┐    ┌──────────────────────────────┐
│  前端 SDK 生成               │    │  后端代码生成                │
│  @ventus/api-client          │    │  internal/api/               │
│                              │    │                              │
│  - types.ts (TypeScript)     │    │  - types.go (Go struct)      │
│  - client.ts (Axios/fetch)   │    │  - handlers.go (接口骨架)    │
│  - hooks.ts (React Query)    │    │  - validators.go (校验)      │
└──────────────────────────────┘    └──────────────────────────────┘
```

### 1.1 为什么用 OpenAPI

| 优势 | 说明 |
|------|------|
| **单一真相源** | 前后端共用一份契约，避免类型不一致 |
| **类型安全** | 自动生成 TS/Go 类型，编译期检查 |
| **文档即代码** | API 文档与实现同步，无需单独维护 |
| **工具生态** | Swagger UI、代码生成、Mock 数据 |

---

## 2. 目录结构

```
ventus/
├── api/                          # API 契约定义（DSL）
│   ├── openapi.yml               # OpenAPI 主文件
│   ├── schemas/                  # 数据模型
│   │   ├── post.yml
│   │   ├── author.yml
│   │   └── common.yml
│   ├── paths/                    # 接口路径
│   │   ├── page.yml              # BFF 接口
│   │   ├── posts.yml             # 文章管理
│   │   └── auth.yml              # 认证
│   └── generate.py               # Python 代码生成脚本
│
├── sdk/                          # 生成的 SDK（独立维护）
│   ├── ts/                       # TypeScript SDK
│   │   ├── package.json
│   │   ├── src/
│   │   │   ├── types.ts          # 生成的类型
│   │   │   ├── client.ts         # 生成的客户端
│   │   │   └── hooks.ts          # React Query hooks
│   │   └── README.md
│   │
│   └── go/                       # Go SDK（可选）
│       └── ...
│
├── frontend/
│   └── packages/
│       └── api-client -> ../../sdk/ts    # 软链接或 npm link
│
└── server/
    └── internal/
        └── api/                    # 生成的 Go 代码
            ├── types.go            # 生成的 struct
            └── validators.go       # 生成的校验器
```

---

## 3. OpenAPI 契约定义

### 3.1 主文件结构

```yaml
# api/openapi.yml
openapi: 3.0.3
info:
  title: Ventus API
  version: 1.0.0

servers:
  - url: /api
    description: 本地开发

paths:
  /page:
    $ref: './paths/page.yml#/Page'
  
  /posts/{id}:
    $ref: './paths/posts.yml#/PostByID'
  
  /admin/posts:
    $ref: './paths/posts.yml#/AdminPosts'
  
  /login:
    $ref: './paths/auth.yml#/Login'

components:
  schemas:
    Post:
      $ref: './schemas/post.yml#/Post'
    PostListResponse:
      $ref: './schemas/post.yml#/PostListResponse'
    Error:
      $ref: './schemas/common.yml#/Error'
```

### 3.2 Schema 定义示例

```yaml
# api/schemas/post.yml
Post:
  type: object
  required:
    - id
    - title
    - slug
    - status
  properties:
    id:
      type: string
      description: 文章唯一标识
    title:
      type: string
      description: 文章标题
    slug:
      type: string
      description: URL 短链接
    content:
      type: string
      description: Markdown 内容
    excerpt:
      type: string
      description: 摘要
    tags:
      type: array
      items:
        type: string
    status:
      type: string
      enum: [draft, published]
    createdAt:
      type: string
      format: date-time
    updatedAt:
      type: string
      format: date-time
    version:
      type: integer
      description: 乐观锁版本号

PostListResponse:
  type: object
  properties:
    items:
      type: array
      items:
        $ref: '#/Post'
    pagination:
      type: object
      properties:
        page:
          type: integer
        pageSize:
          type: integer
        total:
          type: integer
```

### 3.3 Path 定义示例

```yaml
# api/paths/page.yml
Page:
  post:
    summary: BFF 页面数据接口
    tags:
      - BFF
    operationId: getPageData
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - page
              - modules
            properties:
              page:
                type: string
                enum: [home, post, adminPosts, adminEditor]
              modules:
                type: array
                items:
                  type: string
              params:
                type: object
                additionalProperties: true
    responses:
      '200':
        description: 成功
        content:
          application/json:
            schema:
              type: object
              properties:
                page:
                  type: string
                meta:
                  type: object
                  properties:
                    title:
                      type: string
                    description:
                      type: string
                modules:
                  type: object
                  additionalProperties:
                    type: object
                    properties:
                      code:
                        type: integer
                      data:
                        type: object
                      error:
                        type: string
```

---

## 4. 代码生成

> 使用 **Python** 编写生成脚本，优势：
> - 跨平台兼容（Windows/Linux/macOS）
> - 易于扩展（解析 OpenAPI、自定义模板）
> - 丰富的字符串处理和文件操作能力

### 4.1 代码生成脚本（Python）

```python
#!/usr/bin/env python3
# api/generate.py

import subprocess
import sys
import os
from pathlib import Path

def run_command(cmd: list[str], cwd: Path = None, description: str = "") -> bool:
    """运行命令并处理错误"""
    if description:
        print(f"\n🔄 {description}...")
    
    try:
        result = subprocess.run(
            cmd,
            cwd=cwd,
            capture_output=True,
            text=True,
            check=True
        )
        print(f"✅ {description} 完成")
        if result.stdout:
            print(result.stdout)
        return True
    except subprocess.CalledProcessError as e:
        print(f"❌ {description} 失败")
        print(f"   错误: {e.stderr}")
        return False

def generate_ts_sdk():
    """生成 TypeScript SDK"""
    api_file = Path("api/openapi.yml")
    output_dir = Path("sdk/ts/src")
    output_file = output_dir / "client.ts"
    
    # 确保输出目录存在
    output_dir.mkdir(parents=True, exist_ok=True)
    
    # 使用 oazapfts 生成 TS 客户端
    success = run_command(
        [
            "oazapfts",
            str(api_file),
            str(output_file),
            "--useEnumType",
            "--preferUnknown"
        ],
        description="生成 TypeScript 客户端"
    )
    
    if not success:
        return False
    
    # 生成 React Query Hooks（自定义模板）
    hooks_file = output_dir / "hooks.ts"
    generate_react_hooks(hooks_file)
    
    # 格式化代码
    sdk_dir = Path("sdk/ts")
    run_command(["pnpm", "format"], cwd=sdk_dir, description="格式化代码")
    
    print(f"\n📦 TypeScript SDK 生成完成: {output_file}")
    return True

def generate_react_hooks(output_file: Path):
    """生成 React Query Hooks（基于生成的 client）"""
    hooks_content = '''// 生成的 React Query Hooks
import { useQuery, useMutation } from '@tanstack/react-query';
import * as api from './client';

// BFF 页面数据 Hook
export function usePageData(
  page: api.PageRequest['page'],
  modules: string[],
  params?: Record<string, any>,
  options?: any
) {
  return useQuery({
    queryKey: ['page', page, modules, params],
    queryFn: () => api.getPageData({ page, modules, params }),
    staleTime: 5 * 60 * 1000,
    ...options
  });
}

// 文章详情 Hook
export function usePost(id: string, options?: any) {
  return useQuery({
    queryKey: ['post', id],
    queryFn: () => api.getPostById(id),
    enabled: !!id,
    ...options
  });
}

// 更多 hooks...
'''
    output_file.write_text(hooks_content, encoding='utf-8')
    print(f"✅ 生成 React Hooks: {output_file}")

def generate_go_types():
    """生成 Go 类型"""
    api_file = Path("api/openapi.yml")
    output_dir = Path("server/internal/api")
    output_file = output_dir / "types.go"
    
    # 确保输出目录存在
    output_dir.mkdir(parents=True, exist_ok=True)
    
    # 使用 oapi-codegen 生成 Go 类型和接口
    success = run_command(
        [
            "oapi-codegen",
            "-generate", "types,server,spec",
            "-package", "api",
            str(api_file),
            str(output_file)
        ],
        description="生成 Go 类型和接口"
    )
    
    if success:
        print(f"📦 Go 类型生成完成: {output_file}")
    
    return success

def generate_go_validators():
    """生成 Go 校验器（可选）"""
    # 可以基于 OpenAPI 的 validation 生成自定义校验逻辑
    pass

def main():
    """主入口"""
    print("=" * 60)
    print("Ventus API SDK 生成工具")
    print("=" * 60)
    
    # 检查必要工具
    required_tools = {
        "oazapfts": "npm install -g oazapfts",
        "oapi-codegen": "go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest"
    }
    
    print("\n🔍 检查必要工具...")
    for tool, install_cmd in required_tools.items():
        try:
            subprocess.run([tool, "--version"], capture_output=True, check=True)
            print(f"  ✅ {tool}")
        except (subprocess.CalledProcessError, FileNotFoundError):
            print(f"  ❌ {tool} 未安装")
            print(f"     安装命令: {install_cmd}")
            sys.exit(1)
    
    # 生成前端 SDK
    if not generate_ts_sdk():
        print("\n❌ TypeScript SDK 生成失败")
        sys.exit(1)
    
    # 生成后端代码
    if not generate_go_types():
        print("\n❌ Go 类型生成失败")
        sys.exit(1)
    
    print("\n" + "=" * 60)
    print("✨ 所有 SDK 生成完成！")
    print("=" * 60)
    print("\n前端 SDK: sdk/ts/src/")
    print("后端类型: server/internal/api/")

if __name__ == "__main__":
    main()
```

使用方式：

```bash
# 一键生成所有 SDK
cd api && python generate.py

# 输出示例：
# ============================================================
# Ventus API SDK 生成工具
# ============================================================
#
# 🔍 检查必要工具...
#   ✅ oazapfts
#   ✅ oapi-codegen
#
# 🔄 生成 TypeScript 客户端...
# ✅ 生成 TypeScript 客户端 完成
# ✅ 生成 React Hooks: sdk/ts/src/hooks.ts
#
# 🔄 生成 Go 类型和接口...
# ✅ 生成 Go 类型和接口 完成
#
# ============================================================
# ✨ 所有 SDK 生成完成！
# ============================================================

# 带参数（未来扩展）
python generate.py --target ts --watch
python generate.py --target go --validate-only
```

生成的 `client.ts`：

```typescript
// sdk/ts/src/client.ts (生成文件)
import * as Oazapfts from "@oazapfts/runtime";
import * as QS from "qs";

export type Post = {
  id: string;
  title: string;
  slug: string;
  content?: string;
  excerpt?: string;
  tags?: string[];
  status: "draft" | "published";
  createdAt?: string;
  updatedAt?: string;
  version?: number;
};

export type PageRequest = {
  page: "home" | "post" | "adminPosts" | "adminEditor";
  modules: string[];
  params?: Record<string, any>;
};

export type PageResponse = {
  page: string;
  meta?: {
    title?: string;
    description?: string;
  };
  modules?: Record<string, ModuleResult>;
};

// 自动生成的 API 客户端
export function getPageData(body: PageRequest): Promise<PageResponse> {
  return Oazapfts.fetchJson<PageResponse>("/page", {
    method: "POST",
    body: Oazapfts.json(body),
  });
}

export function getPostById(id: string): Promise<Post> {
  return Oazapfts.fetchJson<Post>(`/posts/${id}`, {
    method: "GET",
  });
}

// ... 更多生成的方法
```

### 4.2 自定义 React Hooks

Python 脚本自动生成 hooks.ts：

```typescript
// sdk/ts/src/hooks.ts (自动生成)
import { useQuery, useMutation } from '@tanstack/react-query';
import * as api from './client';

// BFF 页面数据 Hook
export function usePageData(
  page: api.PageRequest['page'],
  modules: string[],
  params?: Record<string, any>
) {
  return useQuery({
    queryKey: ['page', page, modules, params],
    queryFn: () => api.getPageData({ page, modules, params }),
    staleTime: 5 * 60 * 1000,
  });
}

// 文章详情 Hook
export function usePost(id: string) {
  return useQuery({
    queryKey: ['post', id],
    queryFn: () => api.getPostById(id),
    enabled: !!id,
  });
}

// 更新文章 Mutation
export function useUpdatePost() {
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: api.UpdatePostInput }) =>
      api.updatePost(id, data),
  });
}
```

### 4.3 后端 Go 代码生成

Python 脚本自动调用 `oapi-codegen`：

```python
# 已在 generate.py 中实现
def generate_go_types():
    success = run_command(
        [
            "oapi-codegen",
            "-generate", "types,server,spec",
            "-package", "api",
            "api/openapi.yml",
            "server/internal/api/types.go"
        ],
        description="生成 Go 类型和接口"
    )
    return success
```

生成的 `types.go`：

```go
// server/internal/api/types.go (生成文件)
package api

import (
	"time"
)

// Post 文章模型
type Post struct {
	Id        string     `json:"id"`
	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	Content   *string    `json:"content,omitempty"`
	Status    PostStatus `json:"status"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type PostStatus string

const (
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
)

// ServerInterface 生成的接口定义
type ServerInterface interface {
	GetPageData(c *gin.Context)
	GetPostById(c *gin.Context, id string)
	UpdatePost(c *gin.Context, id string)
}
```

生成的 `types.go`：

```go
// server/internal/api/types.go (生成文件)
package api

import (
	"time"
)

// Post 文章模型
type Post struct {
	Id        string     `json:"id"`
	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	Content   *string    `json:"content,omitempty"`
	Excerpt   *string    `json:"excerpt,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	Status    PostStatus `json:"status"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Version   *int       `json:"version,omitempty"`
}

type PostStatus string

const (
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
)

// PageRequest BFF 请求
type PageRequest struct {
	Page    string                 `json:"page"`
	Modules []string               `json:"modules"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// ServerInterface 生成的接口定义
type ServerInterface interface {
	// BFF 页面数据接口
	GetPageData(c *gin.Context)
	
	// 获取文章详情
	GetPostById(c *gin.Context, id string)
	
	// 更新文章
	UpdatePost(c *gin.Context, id string)
	
	// ... 更多接口
}
```

---

## 5. 使用方式

### 5.1 前端使用

```typescript
// pages/home/main.tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { usePageData, usePost } from '@ventus/api-client';

const queryClient = new QueryClient();

function HomePage() {
  // 自动获得类型提示
  const { data, isLoading } = usePageData('home', ['header', 'postList']);
  
  if (isLoading) return <Skeleton />;
  
  // data.modules.postList 有完整类型
  return <PostList data={data.modules.postList.data} />;
}
```

### 5.2 后端使用

```go
// server/internal/interfaces/http/router.go
package http

import (
	"github.com/gin-gonic/gin"
	"ventus/server/internal/api"  // 生成的代码
)

// 实现生成的接口
type Handler struct {
	postService *service.PostService
}

// 确保编译期检查
var _ api.ServerInterface = (*Handler)(nil)

func (h *Handler) GetPageData(c *gin.Context) {
	var req api.PageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, api.Error{Code: 400001, Message: err.Error()})
		return
	}
	
	// 处理 BFF 逻辑
	result, err := h.bffService.GetPageData(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, api.Error{Code: 500001, Message: err.Error()})
		return
	}
	
	c.JSON(200, result)
}

func (h *Handler) GetPostById(c *gin.Context, id string) {
	post, err := h.postService.GetByID(id)
	if err != nil {
		c.JSON(404, api.Error{Code: 404001, Message: "文章不存在"})
		return
	}
	
	c.JSON(200, post)
}

// 路由注册
func SetupRouter(h *Handler) *gin.Engine {
	r := gin.Default()
	
	// 使用生成的路由注册函数
	api.RegisterHandlers(r, h)
	
	return r
}
```

---

## 6. 开发工作流

### 6.1 新增接口流程

```
1. 修改 api/openapi.yml 或子文件
        ↓
2. 运行 api/generate.sh
        ↓
3. 前端：SDK 自动更新，直接使用
   后端：实现生成的接口方法
        ↓
4. 测试联调
```

### 6.2 CI/CD 集成

```yaml
# .github/workflows/api.yml
name: API Contract

on:
  push:
    paths:
      - 'api/**'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      
      - name: Setup Node
        uses: actions/setup-node@v4
        with:
          node-version: 20
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install generators
        run: |
          npm install -g oazapfts
          go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
      
      - name: Generate SDKs
        run: |
          cd api && python generate.py
      
      - name: Check for changes
        run: |
          if [[ -n $(git status --porcelain sdk/ server/internal/api/) ]]; then
            git add sdk/ server/internal/api/
            git commit -m "chore: auto-generate API clients"
            git push
          fi
```

---

## 7. 与编排系统集成

编排系统与 API SDK 配合使用：

```typescript
// pages/home/orchestration.ts
import type { PageOrchestrationConfig } from '@ventus/orchestration';

export const homeConfig: PageOrchestrationConfig = {
  id: 'home',
  // 声明所需 BFF 模块
  modules: ['header', 'postList', 'footer'],
  regions: [/* ... */]
};

// pages/home/main.tsx
import { createOrchestration } from '@ventus/orchestration';
import { QueryClient } from '@tanstack/react-query';
import { homeConfig } from './orchestration';

// 编排系统内部使用 @ventus/api-client 请求 BFF
const orchestration = createOrchestration({
  config: homeConfig,
  queryClient,  // React Query 实例
  modules: { Logo, Nav, PostList, Footer }
});
```

---

## 8. 工具推荐

| 工具 | 用途 | 推荐度 |
|------|------|--------|
| [oazapfts](https://github.com/oazapfts/oazapfts) | TS 客户端生成 | ⭐⭐⭐ |
| [oapi-codegen](https://github.com/deepmap/oapi-codegen) | Go 服务端代码生成 | ⭐⭐⭐ |
| [openapi-generator](https://openapi-generator.tech/) | 多语言生成 | ⭐⭐ |
| [Redoc](https://github.com/Redocly/redoc) | API 文档展示 | ⭐⭐⭐ |
| [Swagger UI](https://swagger.io/tools/swagger-ui/) | 交互式文档 | ⭐⭐ |
| [Prism](https://stoplight.io/open-source/prism) | Mock 服务器 | ⭐⭐ |

---

## 8. 错误码设计

### 8.1 决策：前后端写死

**决策原因**：
- Ventus 错误码数量少（28 个），稳定后几乎不变
- 个人项目，人工维护成本低于自动化成本
- 简单即美德，避免过度设计

**错误码范围**：

| 范围 | 用途 | 示例 |
|------|------|------|
| 0 | 成功 | `SUCCESS` |
| 1-99 | 通用错误 | `INVALID_PARAM`, `INTERNAL_ERROR` |
| 100-199 | 认证错误 | `AUTH_FAILED`, `TOKEN_INVALID` |
| 200-299 | 文章错误 | `POST_NOT_FOUND`, `VERSION_CONFLICT` |
| 300-399 | BFF 模块错误 | `MODULE_NOT_FOUND` |
| 400-499 | 文件上传错误 | `UPLOAD_FAILED` |

### 8.2 前后端定义

**前端**: `frontend/packages/request/src/errors.ts`
```typescript
export enum ErrorCode {
  SUCCESS = 0,
  INVALID_PARAM = 1,
  VERSION_CONFLICT = 206,
}

export const ErrorMessages: Record<ErrorCode, string> = {
  [ErrorCode.SUCCESS]: '成功',
  [ErrorCode.VERSION_CONFLICT]: '版本冲突，文章已被其他人修改',
};
```

**后端**: `server/internal/interfaces/http/response/response.go`
```go
const (
    CodeSuccess = 0
    CodeInvalidParam = 1
    CodeVersionConflict = 206
)
```

### 8.3 同步原则

**如何保持一致**：

1. **命名约定**：
   - Go: `CodeXxxYyy` (驼峰)
   - TS: `XXX_YYY` (大写下划线)
   - 语义相同，仅风格差异

2. **修改流程**：
   - 添加新错误码时，前后端同时修改
   - 通过 Code Review 确保一致
   - 不单独修改一端

3. **数量控制**：
   - 如果错误码超过 50 个，考虑切换到 DSL 生成方案
   - 目前 28 个，人工维护完全可行

### 8.4 为什么不使用检查脚本

尝试过检查脚本，但发现问题：

| 问题 | 说明 |
|------|------|
| **解析脆弱** | 正则解析代码，格式一变就失效 |
| **额外维护** | 需要维护 YAML "真相源"，变成三个地方 |
| **过度设计** | 28 个错误码，人工检查更高效 |
| **信任问题** | 脚本通过≠真的对齐，仍需人工 review |

**结论**：对于 Ventus 规模，**Code Review > 自动化检查**

---

## 9. 演进路线

### Phase 1: 基础契约（MVP）
- [x] 决策：错误码前后端写死
- [ ] 创建 api/openapi.yml 基础结构
- [ ] 编写 api/generate.py 生成脚本
- [ ] 前端实现错误码定义
- [ ] 改造前端 request 包支持错误码
- [ ] 生成 TS SDK @ventus/api-client
- [ ] 生成 Go types
- [ ] 集成到前后端代码

### Phase 2: 完善工具链（P1）
- [ ] CI 自动生成 SDK
- [ ] API 变更检测与通知
- [ ] 版本化发布 SDK

### Phase 3: 高级特性（P2）
- [ ] Mock 服务器（基于 OpenAPI）
- [ ] API 兼容性检查
- [ ] 自动化测试生成
