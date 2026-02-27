import type {
  APIRequest,
  APIResponse,
  BFFPageRequest,
  BFFPageResponse,
  BusinessCode,
} from '@ventus/types';

const API_BASE = import.meta.env.VITE_API_BASE || '/api';

// 统一 API 请求函数
export async function apiRequest<T = unknown, D = Record<string, unknown>>(
  endpoint: string,
  body: APIRequest<D>
): Promise<APIResponse<T>> {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const result = await response.json();
  
  if (result.code !== 0) {
    throw new APIError(result.code, result.message);
  }

  return result;
}

// API 错误类
export class APIError extends Error {
  code: number;
  
  constructor(code: number, message: string) {
    super(message);
    this.code = code;
    this.name = 'APIError';
  }
}

// BFF 页面数据获取
export async function fetchPageData(
  request: BFFPageRequest
): Promise<BFFPageResponse> {
  const response = await apiRequest<BFFPageResponse>('/public', {
    sceneCode: 'page.get',
    data: request as Record<string, unknown>,
  });
  
  return response.data!;
}

// 公开 API 调用
export async function publicAPI<T = unknown, D = Record<string, unknown>>(
  sceneCode: string,
  data?: D
): Promise<T> {
  const response = await apiRequest<T>('/public', {
    sceneCode,
    data,
  });
  return response.data!;
}

// 管理 API 调用
export async function adminAPI<T = unknown, D = Record<string, unknown>>(
  sceneCode: string,
  data?: D
): Promise<T> {
  const response = await apiRequest<T>('/admin', {
    sceneCode,
    data,
  });
  return response.data!;
}
