/**
 * @ventus/request - HTTP 请求库
 * 
 * 特性：
 * 1. 统一的 BFF 请求接口
 * 2. 完整的错误码支持
 * 3. 类型安全的响应处理
 */

import { 
  ErrorCode, 
  getErrorMessage, 
  isReloginError, 
  APIError 
} from './errors';

// ===== 类型定义 =====

export interface RequestConfig {
  /** 基础 URL */
  baseURL?: string;
  /** 超时时间（毫秒） */
  timeout?: number;
  /** 请求头 */
  headers?: Record<string, string>;
  /** 未登录时的跳转路径 */
  loginPath?: string;
}

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  headers?: Record<string, string>;
  body?: unknown;
}

/** 统一 API 响应结构 - 与后端对齐 */
export interface APIResponse<T = unknown> {
  /** 错误码，0 表示成功 */
  code: number;
  /** 错误消息 */
  message: string;
  /** 响应数据 */
  data?: T;
}

/** BFF 模块结果 */
export interface BFFModuleResult<T = unknown> {
  code: number;
  data?: T;
  error?: string;
}

/** BFF 聚合响应 */
export interface BFFResponse {
  page: string;
  modules: Record<string, BFFModuleResult>;
}

/** 请求实例接口 */
export interface RequestInstance {
  /**
   * 发起 BFF 聚合请求
   * @param page - 页面标识
   * @param modules - 所需模块列表
   * @param params - 请求参数
   */
  call<T = unknown>(config: {
    page: string;
    modules: string[];
    params?: Record<string, unknown>;
  }): Promise<Record<string, BFFModuleResult<T>>>;

  /** GET 请求 */
  get<T = unknown>(url: string): Promise<T>;
  
  /** POST 请求 */
  post<T = unknown>(url: string, data?: unknown): Promise<T>;
  
  /** PUT 请求 */
  put<T = unknown>(url: string, data?: unknown): Promise<T>;
  
  /** DELETE 请求 */
  delete<T = unknown>(url: string): Promise<T>;
}

// ===== 实现 =====

class Request implements RequestInstance {
  private config: Required<RequestConfig>;

  constructor(config: RequestConfig = {}) {
    this.config = {
      baseURL: '/api',
      timeout: 10000,
      headers: {},
      loginPath: '/login',
      ...config
    };
  }

  /**
   * 底层 fetch 方法
   */
  private async fetch<T>(url: string, options: RequestOptions = {}): Promise<T> {
    const fullUrl = url.startsWith('http') ? url : `${this.config.baseURL}${url}`;
    
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.config.timeout);

    try {
      const response = await fetch(fullUrl, {
        method: options.method || 'GET',
        headers: {
          'Content-Type': 'application/json',
          ...this.config.headers,
          ...options.headers
        },
        credentials: 'include',
        body: options.body ? JSON.stringify(options.body) : undefined,
        signal: controller.signal
      });

      clearTimeout(timeoutId);

      // HTTP 错误处理
      if (!response.ok) {
        if (response.status === 401) {
          this.handleUnauthorized();
        }
        throw new APIError(
          ErrorCode.INTERNAL_ERROR,
          `HTTP ${response.status}: ${response.statusText}`
        );
      }

      const result: APIResponse<T> = await response.json();
      
      // 业务错误处理
      if (result.code !== ErrorCode.SUCCESS) {
        this.handleBusinessError(result.code, result.message);
      }

      return result.data as T;
    } catch (error) {
      clearTimeout(timeoutId);
      
      // 请求超时
      if (error instanceof Error && error.name === 'AbortError') {
        throw new APIError(ErrorCode.TIMEOUT, '请求超时，请稍后重试');
      }
      
      // 已处理的 APIError 直接抛出
      if (error instanceof APIError) {
        throw error;
      }
      
      // 其他网络错误
      throw new APIError(
        ErrorCode.INTERNAL_ERROR,
        error instanceof Error ? error.message : '网络请求失败'
      );
    }
  }

  /**
   * 处理业务错误
   */
  private handleBusinessError(code: number, message?: string): never {
    // 未登录处理
    if (isReloginError(code)) {
      this.handleUnauthorized();
    }

    throw new APIError(
      code as ErrorCode,
      message || getErrorMessage(code)
    );
  }

  /**
   * 处理未登录
   */
  private handleUnauthorized(): void {
    // 清除登录态
    // 跳转登录页
    if (typeof window !== 'undefined') {
      window.location.href = `${this.config.loginPath}?redirect=${encodeURIComponent(window.location.pathname)}`;
    }
  }

  // ===== 公开方法 =====

  /**
   * BFF 聚合请求
   * @param page - 页面标识，如 'home', 'post', 'adminPosts'
   * @param modules - BFF 模块名列表，如 ['header', 'postList', 'footer']
   * @param params - 页面参数，如 { page: 1, tag: 'go' }
   * 
   * 注意：后端 BFF 模块名与前端组件名不同：
   * - 后端模块：'header', 'postList', 'footer', 'article', 'adminSidebar', ...
   * - 不是前端组件名：'Logo', 'Nav', 'PostList'
   */
  async call<T = unknown>({
    page,
    modules,
    params
  }: {
    page: string;
    modules: string[];
    params?: Record<string, unknown>;
  }): Promise<Record<string, BFFModuleResult<T>>> {
    // 后端 BFF 接口通过 sceneCode 路由，参数嵌套在 data 中
    return this.fetch<BFFResponse>('/public', {
      method: 'POST',
      body: {
        sceneCode: 'page.get',
        data: { page, modules, params }
      }
    }).then(res => res.modules as Record<string, BFFModuleResult<T>>);
  }

  /**
   * 管理端 API 调用
   */
  async admin<T = unknown>(
    sceneCode: string,
    data?: unknown
  ): Promise<T> {
    return this.fetch<T>('/admin', {
      method: 'POST',
      body: { sceneCode, data }
    });
  }

  async get<T = unknown>(url: string): Promise<T> {
    return this.fetch<T>(url, { method: 'GET' });
  }

  async post<T = unknown>(url: string, data?: unknown): Promise<T> {
    return this.fetch<T>(url, { method: 'POST', body: data });
  }

  async put<T = unknown>(url: string, data?: unknown): Promise<T> {
    return this.fetch<T>(url, { method: 'PUT', body: data });
  }

  async delete<T = unknown>(url: string): Promise<T> {
    return this.fetch<T>(url, { method: 'DELETE' });
  }
}

// ===== 导出 =====

export function createRequest(config?: RequestConfig): RequestInstance {
  return new Request(config);
}

// 导出错误码相关
export { 
  ErrorCode, 
  getErrorMessage, 
  isRetryableError, 
  isRefreshableError, 
  isReloginError,
  APIError 
} from './errors';

// 类型已从上面导出，避免重复导出

export default createRequest;
