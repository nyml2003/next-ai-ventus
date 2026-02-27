import { publicAPI } from './bff-client';
import type { LoginRequest, LoginResponse } from '@ventus/types';

const TOKEN_KEY = 'ventus_token';

// 登录
export async function login(credentials: LoginRequest): Promise<void> {
  const data = await publicAPI<LoginResponse>('auth.login', credentials);
  localStorage.setItem(TOKEN_KEY, data.token);
}

// 登出
export function logout(): void {
  localStorage.removeItem(TOKEN_KEY);
}

// 获取 token
export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

// 检查是否已登录
export function isAuthenticated(): boolean {
  return !!getToken();
}
