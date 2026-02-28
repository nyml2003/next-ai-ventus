/**
 * @ventus/request - HTTP 请求库
 */

export interface RequestConfig {
  baseURL?: string;
  timeout?: number;
  headers?: Record<string, string>;
}

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  headers?: Record<string, string>;
  body?: unknown;
}

export interface APIResponse<T = unknown> {
  code: number;
  message: string;
  data?: T;
}

export interface APIRequest<T = unknown> {
  sceneCode: string;
  data?: T;
}

/** BFF 聚合请求 */
export interface BFFRequest {
  /** 页面标识 */
  page: string;
  /** 模块列表 */
  modules: string[];
  /** 请求参数 */
  params?: Record<string, unknown>;
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

export interface RequestInstance {
  /** 发起 BFF 聚合请求 */
  call<T = unknown>(config: {
    page: string;
    modules: string[];
    params?: Record<string, unknown>;
  }): Promise<Record<string, BFFModuleResult<T>>>;
  get<T = unknown>(url: string): Promise<T>;
  post<T = unknown>(url: string, data?: unknown): Promise<T>;
}

class Request implements RequestInstance {
  private config: RequestConfig;

  constructor(config: RequestConfig = {}) {
    this.config = {
      baseURL: '/api',
      timeout: 10000,
      ...config
    };
  }

  private async fetch<T>(url: string, options: RequestOptions = {}): Promise<T> {
    const fullUrl = url.startsWith('http') ? url : `${this.config.baseURL}${url}`;
    
    const response = await fetch(fullUrl, {
      method: options.method || 'GET',
      headers: {
        'Content-Type': 'application/json',
        ...this.config.headers,
        ...options.headers
      },
      credentials: 'include',
      body: options.body ? JSON.stringify(options.body) : undefined
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result: APIResponse<T> = await response.json();
    
    if (result.code !== 0) {
      throw new Error(result.message || `API error! code: ${result.code}`);
    }

    return result.data as T;
  }

  async call<T = unknown>({
    page,
    modules,
    params
  }: {
    page: string;
    modules: string[];
    params?: Record<string, unknown>;
  }): Promise<Record<string, BFFModuleResult<T>>> {
    return this.fetch<BFFResponse>('/public', {
      method: 'POST',
      body: {
        sceneCode: 'page.get',
        data: {
          page,
          modules,
          params
        }
      }
    }).then(res => res.modules as Record<string, BFFModuleResult<T>>);
  }

  async get<T = unknown>(url: string): Promise<T> {
    return this.fetch<T>(url, { method: 'GET' });
  }

  async post<T = unknown>(url: string, data?: unknown): Promise<T> {
    return this.fetch<T>(url, { method: 'POST', body: data });
  }
}

export function createRequest(config?: RequestConfig): RequestInstance {
  return new Request(config);
}

export default createRequest;
