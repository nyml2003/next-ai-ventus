# Ventus API 方案对比与业界实践

本文档对比不同 API 契约管理方案的优劣，供 Ventus 项目参考。

---

## 1. Ventus 当前方案

### 架构
```
┌─────────────────────────────────────────────────────────────┐
│  后端 (Go)                    前端 (TypeScript)              │
│  ───────────                  ───────────────               │
│  手写 enum                     手写 enum                     │
│  手写 handler                  生成 client (openapi)         │
│  手写 response                 手写 hooks                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  OpenAPI 契约     │
                    │  (api/openapi.yml)│
                    └──────────────────┘
```

### 特点
- ✅ 简单可控，适合小团队
- ✅ 错误码人工同步（28 个，数量可控）
- ⚠️ 需要人工维护一致性
- ⚠️ 无运行时类型保证

---

## 2. 业界方案对比

### 方案 A: gRPC + Protocol Buffers（推荐用于中大型项目）

```
┌─────────────────────────────────────────────────────────────┐
│                    .proto 文件（IDL）                        │
│  service PostService {                                      │
│    rpc GetPost(GetPostReq) returns (GetPostResp);           │
│  }                                                          │
└────────────────────────────┬────────────────────────────────┘
                             │
         ┌───────────────────┴───────────────────┐
         ▼                                       ▼
┌────────────────────────┐         ┌──────────────────────────────┐
│  Go 生成代码            │         │  TypeScript 生成代码          │
│  - post.pb.go          │         │  - post_pb.ts                │
│  - post_grpc.pb.go     │         │  - post.client.ts            │
└────────────────────────┘         └──────────────────────────────┘
         │                                       │
         ▼                                       ▼
┌────────────────────────┐         ┌──────────────────────────────┐
│  gRPC Server           │◄────────┤  gRPC-Web / Connect          │
│  直接实现生成接口       │  HTTP/2 │  生成客户端                   │
└────────────────────────┘         └──────────────────────────────┘
```

**优势**:
- ✅ 强类型，编译期检查
- ✅ 自动生成前后端代码
- ✅ 支持流式传输
- ✅ 二进制协议，性能高
- ✅ 生态完善（gateway, load balancer 等）

**劣势**:
- ❌ 需要引入 gRPC 运行时
- ❌ 浏览器需要 gRPC-Web 转换层
- ❌ 学习成本较高
- ❌ 不适合简单项目

**适用**: 中大型团队，微服务架构

**工具链**:
- `protoc` - 编译器
- `protoc-gen-go` - Go 代码生成
- `protoc-gen-ts` / `ts-proto` - TypeScript 代码生成
- `grpc-gateway` - RESTful 网关
- `Connect` - 现代 gRPC 替代方案（推荐）

---

### 方案 B: OpenAPI + 代码生成（Ventus 当前方向）

```
┌─────────────────────────────────────────────────────────────┐
│                    OpenAPI 3.0 (api/openapi.yml)           │
└────────────────────────────┬────────────────────────────────┘
                             │
         ┌───────────────────┴───────────────────┐
         ▼                                       ▼
┌────────────────────────┐         ┌──────────────────────────────┐
│  oapi-codegen          │         │  oazapfts / openapi-generator│
│  生成 Go 类型+接口      │         │  生成 TS 类型+客户端          │
└──────────┬─────────────┘         └──────────┬─────────────────┘
           │                                  │
           ▼                                  ▼
┌────────────────────────┐         ┌──────────────────────────────┐
│  Go Handler 实现接口    │         │  直接使用生成客户端           │
└────────────────────────┘         └──────────────────────────────┘
```

**优势**:
- ✅ HTTP/REST 原生支持，无需额外层
- ✅ 工具生态丰富（Swagger UI, Postman 等）
- ✅ 前后端都可以生成代码
- ✅ 适合 Web 项目

**劣势**:
- ⚠️ 生成代码质量参差不齐
- ⚠️ 需要额外维护 OpenAPI 文件
- ⚠️ 类型安全弱于 gRPC

**适用**: Web 项目，RESTful API

**工具链**:
- `oapi-codegen` - Go 代码生成
- `oazapfts` - TypeScript 客户端生成
- `openapi-generator` - 多语言生成
- `Redoc` / `Swagger UI` - 文档展示

---

### 方案 C: tRPC（全栈 TypeScript 项目推荐）

```
┌─────────────────────────────────────────────────────────────┐
│                    Router 定义（TS 文件）                     │
│  const appRouter = router({                                 │
│    post: postRouter,                                        │
│  });                                                        │
└────────────────────────────┬────────────────────────────────┘
                             │
         ┌───────────────────┴───────────────────┐
         │ 类型推导 + 代码生成                     │
         ▼                                       ▼
┌────────────────────────┐         ┌──────────────────────────────┐
│  后端 (Node.js)        │         │  前端 (Browser/React)        │
│  - 直接实现 router      │         │  - 类型安全调用               │
│  - 类型自动推导         │◄────────┤    client.post.getById()     │
│  - 支持 middleware      │  HTTP   │  - 自动类型补全               │
└────────────────────────┘         └──────────────────────────────┘
```

**优势**:
- ✅ 极致的类型安全（端到端类型）
- ✅ 无需代码生成，类型自动推导
- ✅ 支持订阅（WebSocket）
- ✅ 开发体验极佳

**劣势**:
- ❌ 后端必须是 Node.js / TypeScript
- ❌ 生态相对封闭
- ❌ 不适合异构语言栈

**适用**: 全栈 TypeScript 项目

**工具链**:
- `@trpc/server` - 服务端
- `@trpc/client` - 客户端
- `@trpc/react-query` - React 集成

---

### 方案 D: GraphQL（灵活查询场景）

```
┌─────────────────────────────────────────────────────────────┐
│                    Schema 定义（.graphql）                    │
│  type Post {                                                │
│    id: ID!                                                  │
│    title: String!                                           │
│  }                                                          │
└────────────────────────────┬────────────────────────────────┘
                             │
         ┌───────────────────┴───────────────────┐
         ▼                                       ▼
┌────────────────────────┐         ┌──────────────────────────────┐
│  Go 实现 Resolver      │         │  TypeScript 生成代码          │
│  - graphql-go          │◄────────┤  - CodeGen                   │
│  - gqlgen              │  HTTP   │  - Apollo Client             │
└────────────────────────┘         └──────────────────────────────┘
```

**优势**:
- ✅ 客户端决定查询字段，减少过度获取
- ✅ 强类型 Schema
- ✅ 一个端点，所有查询

**劣势**:
- ❌ 学习曲线陡峭
- ❌ 缓存复杂
- ❌ N+1 查询问题
- ❌ 不适合简单 CRUD

**适用**: 复杂数据关联，前端需求多变

---

## 3. Ventus 演进建议

### 当前（MVP）
```
OpenAPI 契约
├── 生成 TS Client（oazapfts）
└── 生成 Go Types（oapi-codegen）

错误码：人工维护 + Python 检查脚本
```

### 未来可选演进

#### 方案 1: 保持现状（推荐）
- Ventus 是个人博客，API 数量可控（< 50 个）
- 当前方案足够简单有效
- 维护成本低

#### 方案 2: 引入 Connect（如果 API 增长）
```
Connect (gRPC + REST 双协议)
├── 从 .proto 生成 Go 代码
├── 生成 TS 代码
└── 同时支持 gRPC 和 HTTP/JSON
```

#### 方案 3: 全类型安全（如果后端改 Node.js）
```
tRPC
├── 共享类型定义
├── 无代码生成步骤
└── 极致开发体验
```

---

## 4. 决策树

```
选择 API 方案？
├── 团队语言栈统一为 TypeScript？
│   └── 是 → tRPC（最佳选择）
│
├── 需要极致性能（二进制协议）？
│   └── 是 → gRPC
│
├── 前端查询需求多变？
│   └── 是 → GraphQL
│
├── 传统 Web 项目，RESTful？
│   └── 是 → OpenAPI + 代码生成（Ventus 当前）
│
└── 简单项目，API 数量 < 50？
    └── 是 → 人工维护（Ventus 当前简化版）
```

---

## 5. 推荐资源

### OpenAPI 生态
- [oapi-codegen](https://github.com/deepmap/oapi-codegen) - Go 代码生成
- [oazapfts](https://github.com/oazapfts/oazapfts) - TypeScript 生成
- [OpenAPI Generator](https://openapi-generator.tech/) - 多语言生成
- [Swagger Editor](https://editor.swagger.io/) - 在线编辑器

### gRPC 生态
- [Connect](https://connectrpc.com/) - 现代 gRPC（推荐）
- [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) - REST 网关
- [Buf](https://buf.build/) - Protobuf 工具链

### tRPC
- [tRPC 文档](https://trpc.io/)
- [create-t3-app](https://create.t3.gg/) - 全栈 TS 模板

### GraphQL
- [gqlgen](https://gqlgen.com/) - Go GraphQL 框架
- [Apollo Client](https://www.apollographql.com/docs/react/) - 前端客户端
