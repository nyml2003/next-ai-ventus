/**
 * 从 URL 查询参数中获取指定 key 的值
 */
export function getQueryParam(key: string): string | null {
  const url = new URL(window.location.href);
  return url.searchParams.get(key);
}

/**
 * 获取所有查询参数
 */
export function getQueryParams(): Record<string, string> {
  const url = new URL(window.location.href);
  const params: Record<string, string> = {};
  url.searchParams.forEach((value, key) => {
    params[key] = value;
  });
  return params;
}

/**
 * 通用的页面 Props 类型
 * 每个页面通过 URL 参数接收输入
 */
export interface PageProps {
  /** URL 查询参数 */
  query: Record<string, string>;
  /** 获取单个参数 */
  getParam: (key: string) => string | null;
}

/**
 * 创建页面 Props（在组件内调用）
 */
export function createPageProps(): PageProps {
  return {
    query: getQueryParams(),
    getParam: getQueryParam,
  };
}
