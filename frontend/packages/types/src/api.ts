// API 请求/响应类型

// 统一 API 请求
export interface APIRequest<T = Record<string, unknown>> {
  sceneCode: string;
  data?: T;
}

// 统一 API 响应
export interface APIResponse<T = unknown> {
  code: number;
  message: string;
  data?: T;
}

// 业务错误码
export const BusinessCode = {
  // 成功
  SUCCESS: 0,
  // 通用错误
  INVALID_PARAM: 1,
  INTERNAL_ERROR: 2,
  UNAUTHORIZED: 3,
  FORBIDDEN: 4,
  NOT_FOUND: 5,
  // 认证错误
  AUTH_FAILED: 100,
  TOKEN_INVALID: 101,
  TOKEN_EXPIRED: 102,
  TOKEN_MISSING: 103,
  INVALID_CREDENTIALS: 104,
  // 文章错误
  POST_NOT_FOUND: 200,
  POST_ALREADY_EXISTS: 201,
  SLUG_EXISTS: 202,
  INVALID_TITLE: 203,
  INVALID_CONTENT: 204,
  VERSION_CONFLICT: 206,
} as const;

// BFF 模块请求
export interface BFFPageRequest {
  page: string;
  modules: string[];
  params?: Record<string, unknown>;
}

// BFF 模块结果
export interface BFFModuleResult<T = unknown> {
  code: number;
  data?: T;
  error?: string;
}

// BFF 页面响应
export interface BFFPageResponse {
  page: string;
  modules: Record<string, BFFModuleResult>;
}

// 文章
export interface Post {
  id: string;
  title: string;
  slug: string;
  content?: string;
  html?: string;
  excerpt?: string;
  tags: string[];
  status: 'draft' | 'published';
  createdAt: string;
  updatedAt: string;
  publishedAt?: string;
  wordCount?: number;
}

// 文章列表项
export interface PostListItem {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  tags: string[];
  date: string;
  href: string;
}

// 分页信息
export interface PaginationInfo {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

// 登录请求
export interface LoginRequest {
  username: string;
  password: string;
}

// 登录响应
export interface LoginResponse {
  token: string;
}
