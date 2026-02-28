/**
 * @ventus/store - 基于 Zustand + React Context 的状态管理
 */

import { create, type StoreApi } from "zustand";
import * as React from "react";
import type { RequestInstance } from "@ventus/request";
import type { PageOrchestrationConfig, PageProps } from "@ventus/types";

// ==================== Store 类型定义 ====================

export interface RequestState<T = unknown> {
  data: T | null;
  loading: boolean;
  error: Error | null;
}

export interface OrchestrationState {
  // 编排配置
  config: PageOrchestrationConfig | null;

  // 请求实例
  request: RequestInstance | null;

  // URL 参数
  pageProps: PageProps | null;

  // 请求缓存
  requestCache: Map<string, RequestState>;

  // Spacing 解析器
  resolver: ((token: string) => string) | null;

  // 模块注册表 - 用于请求时带上模块列表
  modules: Record<string, unknown> | null;
}

export interface OrchestrationActions {
  setConfig: (config: PageOrchestrationConfig) => void;
  setRequest: (request: RequestInstance) => void;
  setPageProps: (pageProps: PageProps) => void;
  setResolver: (resolver: (token: string) => string) => void;
  setModules: (modules: Record<string, unknown>) => void;
  setRequestState: <T>(key: string, state: Partial<RequestState<T>>) => void;
  getRequestState: <T>(key: string) => RequestState<T>;
  /** 获取模块注册表的所有 key */
  getModuleKeys: () => string[];
}

export type OrchestrationStore = OrchestrationState & OrchestrationActions;

// ==================== Store 工厂函数 ====================

export type OrchestrationStoreApi = StoreApi<OrchestrationStore>;

export function createOrchestrationStore(): OrchestrationStoreApi {
  return create<OrchestrationStore>((set, get) => ({
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
        ...state,
      };
      cache.set(key, updated);
      set({ requestCache: new Map(cache) });
    },
    getRequestState: <T>(key: string): RequestState<T> => {
      const cache = get().requestCache;
      return (
        (cache.get(key) as RequestState<T>) ?? {
          data: null,
          loading: false,
          error: null,
        }
      );
    },

    resolver: null,
    setResolver: (resolver) => set({ resolver }),

    modules: null,
    setModules: (modules) => set({ modules }),
    getModuleKeys: () => {
      const modules = get().modules;
      return modules ? Object.keys(modules) : [];
    },
  }));
}

// 全局默认 store（向后兼容）
export const useOrchestrationStore = createOrchestrationStore();

// ==================== React Context ====================

const StoreContext = React.createContext<OrchestrationStoreApi | null>(null);

export interface StoreProviderProps {
  children: React.ReactNode;
  config: PageOrchestrationConfig;
  request: RequestInstance;
  resolver: (token: string) => string;
  pageProps: PageProps;
  /** 模块注册表 */
  modules: Record<string, unknown>;
}

export const StoreProvider: React.FC<StoreProviderProps> = ({
  children,
  config,
  request,
  resolver,
  pageProps,
  modules,
}) => {
  // 每个 Provider 创建独立的 store 实例
  const storeRef = React.useRef<OrchestrationStoreApi | null>(null);

  if (!storeRef.current) {
    storeRef.current = createOrchestrationStore();
    // 立即初始化
    storeRef.current.getState().setConfig(config);
    storeRef.current.getState().setRequest(request);
    storeRef.current.getState().setResolver(resolver);
    storeRef.current.getState().setPageProps(pageProps);
    storeRef.current.getState().setModules(modules);
  }

  return React.createElement(
    StoreContext.Provider,
    { value: storeRef.current },
    children,
  );
};

// ==================== Hooks ====================

function useStore(): OrchestrationStoreApi {
  const context = React.useContext(StoreContext);
  if (!context) {
    throw new Error("useStore must be used within a StoreProvider");
  }
  return context;
}

export interface UseRequestOptions<T = unknown> {
  /** 页面标识 */
  page: string;
  /** 模块列表 - 如果不传则自动从 store 获取 */
  modules?: string[];
  /** 请求参数 */
  params?: Record<string, unknown>;
  /** 依赖列表 */
  deps?: unknown[];
}

export function useRequest<T = unknown>({
  page,
  modules: customModules,
  params,
  deps = [],
}: UseRequestOptions<T>): RequestState<Record<string, import("@ventus/request").BFFModuleResult<T>>> & { refetch: () => void } {
  const store = useStore();
  const request = store.getState().request;
  
  // 从 store 获取模块列表（如果没有传自定义 modules）
  const moduleKeys = React.useMemo(() => {
    return customModules || store.getState().getModuleKeys();
  }, [customModules, store]);

  // 生成缓存 key
  const cacheKey = React.useMemo(() => {
    return `${page}:${moduleKeys.join(',')}:${JSON.stringify(params || {})}`;
  }, [page, moduleKeys, params]);

  // 使用 subscribe 监听状态变化
  const [state, setState] = React.useState<RequestState<Record<string, import("@ventus/request").BFFModuleResult<T>>>>(() =>
    store.getState().getRequestState<Record<string, import("@ventus/request").BFFModuleResult<T>>>(cacheKey),
  );

  React.useEffect(() => {
    setState(store.getState().getRequestState<Record<string, import("@ventus/request").BFFModuleResult<T>>>(cacheKey));
    const unsubscribe = store.subscribe((newState) => {
      const newRequestState = (newState as OrchestrationState).requestCache.get(
        cacheKey,
      ) as RequestState<Record<string, import("@ventus/request").BFFModuleResult<T>>> | undefined;
      if (newRequestState) {
        setState(newRequestState);
      }
    });
    return unsubscribe;
  }, [cacheKey, store]);

  const fetchData = React.useCallback(async () => {
    if (!request) return;

    store.getState().setRequestState(cacheKey, { loading: true, error: null });

    try {
      const data = await request.call<T>({ page, modules: moduleKeys, params });
      store.getState().setRequestState(cacheKey, { data, loading: false });
    } catch (err) {
      store.getState().setRequestState(cacheKey, {
        error: err instanceof Error ? err : new Error(String(err)),
        loading: false,
      });
    }
  }, [cacheKey, request, page, moduleKeys, params, store]);

  React.useEffect(() => {
    // 如果缓存中已有数据，不重复请求
    if (state.data === null && !state.loading) {
      fetchData();
    }
  }, deps);

  return {
    ...state,
    refetch: fetchData,
  };
}

export function usePageProps(): PageProps {
  const store = useStore();
  const [pageProps, setPageProps] = React.useState<PageProps | null>(
    () => store.getState().pageProps,
  );

  React.useEffect(() => {
    const unsubscribe = store.subscribe((state) => {
      setPageProps((state as OrchestrationState).pageProps);
    });
    return unsubscribe;
  }, [store]);

  if (!pageProps) {
    throw new Error("pageProps not initialized in store");
  }

  return pageProps;
}

export function useResolver(): (token: string) => string {
  const store = useStore();
  const [resolver, setResolver] = React.useState<
    ((token: string) => string) | null
  >(() => store.getState().resolver);

  React.useEffect(() => {
    const unsubscribe = store.subscribe((state) => {
      setResolver((state as OrchestrationState).resolver);
    });
    return unsubscribe;
  }, [store]);

  if (!resolver) {
    return () => "0px";
  }

  return resolver;
}

export function useConfig(): PageOrchestrationConfig {
  const store = useStore();
  const [config, setConfig] = React.useState<PageOrchestrationConfig | null>(
    () => store.getState().config,
  );

  React.useEffect(() => {
    const unsubscribe = store.subscribe((state) => {
      setConfig((state as OrchestrationState).config);
    });
    return unsubscribe;
  }, [store]);

  if (!config) {
    throw new Error("config not initialized in store");
  }

  return config;
}

/** 获取模块注册表的所有 key */
export function useModuleKeys(): string[] {
  const store = useStore();
  const [keys, setKeys] = React.useState<string[]>(() =>
    store.getState().getModuleKeys(),
  );

  React.useEffect(() => {
    const unsubscribe = store.subscribe(() => {
      setKeys(store.getState().getModuleKeys());
    });
    return unsubscribe;
  }, [store]);

  return keys;
}

/** 获取指定模块的数据（从聚合请求结果中） */
export function useModuleData<T = unknown>(
  moduleName: string,
): import("@ventus/request").BFFModuleResult<T> | undefined {
  const store = useStore();
  const [data, setData] = React.useState<
    import("@ventus/request").BFFModuleResult<T> | undefined
  >(() => {
    const state = store.getState();
    // 获取最近一次的聚合请求结果
    for (const [, requestState] of state.requestCache) {
      if (requestState.data && requestState.data[moduleName]) {
        return requestState.data[moduleName];
      }
    }
    return undefined;
  });

  React.useEffect(() => {
    const unsubscribe = store.subscribe((state) => {
      for (const [, requestState] of (state as OrchestrationState).requestCache) {
        if (requestState.data && requestState.data[moduleName]) {
          setData(requestState.data[moduleName]);
          return;
        }
      }
    });
    return unsubscribe;
  }, [store, moduleName]);

  return data;
}

// ==================== 辅助函数 ====================

export function createPagePropsFromURL(): PageProps {
  const url = new URL(window.location.href);

  // 解析路由参数（简单实现，实际可能需要路由配置）
  const pathParts = url.pathname.split("/").filter(Boolean);
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
    query,
  };
}

export default useOrchestrationStore;
