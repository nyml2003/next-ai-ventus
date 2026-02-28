# Ventus 前后端 API 对齐检查

本文档记录前后端接口的对齐情况。

---

## 1. 接口概览

### 后端实现（Go）

| 端点 | SceneCode | 状态 | 描述 |
|------|-----------|------|------|
| `POST /api/public` | `auth.login` | ⚠️ | 登录 |
| `POST /api/public` | `page.get` | ✅ | BFF 页面数据 |
| `POST /api/public` | `post.recordView` | ⚠️ | 记录阅读 |
| `POST /api/admin` | `post.create` | ✅ | 创建文章 |
| `POST /api/admin` | `post.update` | ✅ | 更新文章 |
| `POST /api/admin` | `post.delete` | ✅ | 删除文章 |
| `POST /api/admin` | `post.get` | ✅ | 获取文章详情 |
| `POST /api/admin` | `post.list` | ✅ | 获取文章列表 |
| `POST /api/admin` | `file.upload` | ⚠️ | 文件上传（特殊处理） |

### 前端封装（TS）

| 方法 | SceneCode | 状态 | 描述 |
|------|-----------|------|------|
| `request.call()` | `page.get` | ✅ | BFF 请求 |
| `request.admin()` | 通用 | ✅ | 管理端通用 |
| `login()` | ❌ | **缺失** | 未封装 |
| `recordView()` | ❌ | **缺失** | 未封装 |
| `uploadFile()` | ❌ | **缺失** | 未封装（需要 multipart） |

---

## 2. 问题分析

### 2.1 缺失的接口

#### 登录接口
```typescript
// 前端需要添加
async login(username: string, password: string): Promise<{ token: string }> {
  return this.fetch('/public', {
    method: 'POST',
    body: { sceneCode: 'auth.login', data: { username, password } }
  });
}
```

#### 记录阅读
```typescript
// 前端需要添加
async recordView(postId: string): Promise<void> {
  return this.fetch('/public', {
    method: 'POST',
    body: { sceneCode: 'post.recordView', data: { id: postId } }
  });
}
```

#### 文件上传
```typescript
// 前端需要添加（特殊：multipart/form-data）
async uploadFile(file: File): Promise<{ url: string }> {
  const formData = new FormData();
  formData.append('file', file);
  
  return fetch('/api/admin', {
    method: 'POST',
    headers: {
      'X-Scene-Code': 'file.upload'  // 可能需要特殊 header
    },
    body: formData,
    credentials: 'include'
  });
}
```

### 2.2 请求格式对比

**标准请求（JSON）**:
```typescript
// 前端发送
{
  sceneCode: 'post.create',
  data: { title: 'xxx', content: 'yyy' }
}

// 后端接收
type APIRequest struct {
    SceneCode string                 `json:"sceneCode"`
    Data      map[string]interface{} `json:"data"`
}
```

**文件上传（Multipart）**:
- 后端 `file.upload` 使用 `c.Request.FormFile("file")` 接收
- 前端需要特殊处理，不能使用 JSON

---

## 3. 参数对齐检查

### 3.1 `post.create`

| 字段 | 前端期望 | 后端期望 | 状态 |
|------|----------|----------|------|
| `title` | string | string | ✅ |
| `content` | string | string | ✅ |
| `tags` | string[] | []interface{} | ⚠️ 类型转换 |
| `status` | 'draft' \| 'published' | string | ⚠️ 建议前端 enum |

### 3.2 `post.update`

| 字段 | 前端期望 | 后端期望 | 状态 |
|------|----------|----------|------|
| `id` | string | string | ✅ |
| `version` | number | float64 | ✅ |
| `title` | string? | *string | ✅ 可选 |
| `content` | string? | *string | ✅ 可选 |
| `status` | string? | *string | ✅ 可选 |
| `tags` | string[]? | []string | ✅ 可选 |

### 3.3 `post.list`

| 字段 | 前端期望 | 后端期望 | 状态 |
|------|----------|----------|------|
| `page` | number | float64 | ✅ |
| `pageSize` | number | float64 | ✅ |
| `status` | string? | string | ✅ 可选 |
| `tag` | string? | string | ✅ 可选 |

---

## 4. 响应格式对齐

### 4.1 统一响应结构 ✅

前后端一致：
```typescript
interface APIResponse<T> {
  code: number;      // 0 = 成功
  message: string;   // 错误消息
  data?: T;          // 响应数据
}
```

### 4.2 BFF 响应结构 ✅

```typescript
interface BFFResponse {
  page: string;
  modules: Record<string, {
    code: number;
    data?: any;
    error?: string;
  }>;
}
```

### 4.3 文章对象对齐 ⚠️

需要确认字段命名是否一致：
- `id` vs `ID`
- `createdAt` vs `CreatedAt` vs `created_at`
- `slug` 格式

---

## 5. 修复建议

### 5.1 前端 request 包增强

```typescript
// frontend/packages/request/src/index.ts

export interface RequestInstance {
  // 已有方法...
  
  // 需要添加
  login(username: string, password: string): Promise<{ token: string }>;
  recordView(postId: string): Promise<void>;
  uploadFile(file: File): Promise<{ url: string }>;
  
  // Admin API 类型安全封装
  createPost(data: CreatePostInput): Promise<Post>;
  updatePost(id: string, data: UpdatePostInput, version: number): Promise<Post>;
  deletePost(id: string): Promise<void>;
  getPost(id: string): Promise<Post>;
  listPosts(params: ListPostsInput): Promise<ListPostsResult>;
}
```

### 5.2 类型定义文件

```typescript
// frontend/packages/request/src/types.ts

export interface Post {
  id: string;
  title: string;
  slug: string;
  content?: string;
  excerpt?: string;
  tags?: string[];
  status: 'draft' | 'published';
  createdAt?: string;
  updatedAt?: string;
  version?: number;
}

export interface CreatePostInput {
  title: string;
  content: string;
  tags?: string[];
  status?: 'draft' | 'published';
}

export interface UpdatePostInput {
  title?: string;
  content?: string;
  tags?: string[];
  status?: 'draft' | 'published';
}

export interface ListPostsInput {
  page?: number;
  pageSize?: number;
  status?: string;
  tag?: string;
}

export interface ListPostsResult {
  items: Post[];
  total: number;
  page: number;
  pageSize: number;
}
```

### 5.3 后端字段命名检查

需要确认 Go 返回的 JSON 字段名是否与前端期望一致：

```go
// 后端当前实现
type Post struct {
    ID      string   `json:"id"`       // 前端期望 "id" ✅
    Title   string   `json:"title"`    // 前端期望 "title" ✅
    CreatedAt time.Time `json:"createdAt"` // 前端期望 "createdAt" ✅
}
```

---

## 6. 检查清单

### 开发前
- [ ] 确认接口路径一致 (`/api/public`, `/api/admin`)
- [ ] 确认 SceneCode 命名一致
- [ ] 确认请求参数字段名一致
- [ ] 确认响应字段命名风格（camelCase）

### 联调时
- [ ] 测试每个 SceneCode 的完整流程
- [ ] 验证错误码返回是否正确
- [ ] 验证可选参数处理（空值 vs 未传）
- [ ] 验证文件上传（特殊处理）

### 维护期
- [ ] 修改后端接口时同步更新前端
- [ ] 保持类型定义文件同步
- [ ] 更新本文档
