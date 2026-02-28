/**
 * @ventus/store - 基于 Zustand + React Context 的状态管理
 */

import { create } from 'zustand';
import * as React from 'react';
import { createContext, useContext, useCallback } from 'react';
import type { RequestInstance } from '@ventus/request';
import type { PageOrchestrationConfig, PageProps } from '@ventus/types';

// ==================== Store 类型定义 ====================

export interface RequestState<T = unknown> {
  data: T | null;
  loading: boolean;
  error: Error | null;
}

export interface OrchestrationStore {
  // 编排配置
  config: PageOrchestrationConfig | null;
  setConfig: (config: PageOrchestrationConfig) => void;
  
  // 请求实例
  request: RequestInstance | null;
  setRequest: (request: RequestInstance) => void;
  
  // URL 参数
  pageProps: PageProps | null;
  setPageProps: (pageProps: PageProps) => void;
  
  // 请求缓存
  requestCache: Map<string, RequestState>;
  setRequestState: <T>(key: string, state: Partial<RequestState<T>>) => void;
  getRequestState: <T>(key: string) => RequestState<T>;
  
  // Spacing 解析器
  resolver: ((token: string) => string) | null;
  setResolver: (resolver: (token: string) => string) => void;
}

// ==================== Zustand Store ====================

export const useOrchestrationStore = create<OrchestrationStore>((set, get) => ({
  config: null,
  setConfig: (config) => set({ config }),
  
  request: null,
  setRequest: (request) => set({ request }),
  
  pageProps: null,
  setPageProps: (pageProps) => set({ pageProps }),
  
  requestCache: new Map(),
  setRequestState: <T>(key: string, state: Partial<RequestState<T>>) => {
    const cache = get().requestCache;
    const current = cache.get(key) as RequestState<T> | undefined;
    const updated: RequestState<T> = {
      data: current?.data ?? null,
      loading: current?.loading ?? false,
      error: current?.error ?? null,
      ...state
    };
    cache.set(key, updated);
    set({ requestCache: new Map(cache) });
  },
  getRequestState: <T>(key: string): RequestState<T> => {
    const cache = get().requestCache;
    return (cache.get(key) as RequestState<T>) ?? {
      data: null,
      loading: false,
      error: null
    };
  },
  
  resolver: null,
  setResolver: (resolver) => set({ resolver })
}));

// ==================== React Context ====================

type StoreContextType = ReturnType<typeof useOrchestrationStore>;

const StoreContext = createContext<StoreContextType | null>(null);

export interface StoreProviderProps {
  children: React.ReactNode;
  config: PageOrchestrationConfig;
  request: RequestInstance;
  resolver: (token: string) => string;
  pageProps: PageProps;
}

export const StoreProvider: React.FC<StoreProviderProps> = ({
  children,
  config,
  request,
  resolver,
  pageProps
}) => {
  // 创建带有初始值的 store
  const store = React.useMemo(() => {
    const newStore = useOrchestrationStore.getState();
    newStore.setConfig(config);
    newStore.setRequest(request);
    newStore.setResolver(resolver);
    newStore.setPageProps(pageProps);
    return newStore;
  }, []);
  
  // 同步更新 store
  React.useMemo(() => {
    store.setConfig(config);
    store.setRequest(request);
    store.setResolver(resolver);
    store.setPageProps(pageProps);
  }, [config, request, resolver, pageProps]);
  
  // 使用 React.createElement 替代 JSX
  return React.createElement(
    StoreContext.Provider,
    { value: store },
    children
  );
};

// ==================== Hooks ====================

function useStore(): OrchestrationStore {
  const context = useContext(StoreContext);
  if (!context) {
    throw new Error('useStore must be used within a StoreProvider');
  }
  return context;
}

export interface UseRequestOptions<T = unknown> {
  scene: string;
  params?: Record<string, unknown>;
  deps?: unknown[];
}

export function useRequest<T = unknown>({
  scene,
  params,
  deps = []
}: UseRequestOptions<T>): RequestState<T> & { refetch: () => void } {
  const store = useStore();
  const request = store.request;
  
  // 生成缓存 key
  const cacheKey = React.useMemo(() => {
    return `${scene}:${JSON.stringify(params || {})}`;
  }, [scene, params]);
  
  const state = store.getRequestState<T>(cacheKey);
  
  const fetchData = useCallback(async () => {
    if (!request) return;
    
    store.setRequestState(cacheKey, { loading: true, error: null });
    
    try {
      const data = await request.call<T>({ scene, params });
      store.setRequestState(cacheKey, { data, loading: false });
    } catch (err) {
      store.setRequestState(cacheKey, {
        error: err instanceof Error ? err : new Error(String(err)),
        loading: false
      });
    }
  }, [cacheKey, request, scene, params, store]);
  
  React.useEffect(() => {
    // 如果缓存中已有数据，不重复请求
    if (state.data === null && !state.loading) {
      fetchData();
    }
  }, deps);
  
  return {
    ...state,
    refetch: fetchData
  };
}

export function usePageProps(): PageProps {
  const store = useStore();
  const pageProps = store.pageProps;
  
  if (!pageProps) {
    throw new Error('pageProps not initialized in store');
  }
  
  return pageProps;
}

export function useResolver(): (token: string) => string {
  const store = useStore();
  const resolver = store.resolver;
  
  if (!resolver) {
    return () => '0px';
  }
  
  return resolver;
}

// ==================== 辅助函数 ====================

export function createPagePropsFromURL(): PageProps {
  const url = new URL(window.location.href);
  
  // 解析路由参数（简单实现，实际可能需要路由配置）
  const pathParts = url.pathname.split('/').filter(Boolean);
  const params: Record<string, string> = {};
  
  // 解析查询参数
  const query: Record<string, string> = {};
  url.searchParams.forEach((value, key) => {
    query[key] = value;
  });
  
  return {
    getParam: (key: string) => params[key],
    getQuery: (key: string) => query[key],
    params,
    query
  };
}

export default useOrchestrationStore;
