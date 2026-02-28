/**
 * Store 单元测试
 */
import { createOrchestrationStore } from '../index';
import type { PageOrchestrationConfig, PageProps } from '@ventus/types';
import type { RequestInstance } from '@ventus/request';

describe('OrchestrationStore', () => {
  const mockConfig: PageOrchestrationConfig = {
    id: 'test-page',
    regions: []
  };

  const mockRequest = {
    call: jest.fn()
  } as unknown as RequestInstance;

  const mockPageProps: PageProps = {
    getParam: (key: string) => undefined,
    getQuery: (key: string) => undefined,
    params: {},
    query: {}
  };

  const mockResolver = (token: string) => {
    const map: Record<string, string> = { lg: '24px', md: '16px' };
    return map[token] || '0px';
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('createOrchestrationStore', () => {
    it('应该创建独立的 store 实例', () => {
      const store1 = createOrchestrationStore();
      const store2 = createOrchestrationStore();
      
      // 两个实例应该是不同的对象
      expect(store1).not.toBe(store2);
      
      // 初始状态应该为 null
      expect(store1.getState().config).toBeNull();
      expect(store2.getState().config).toBeNull();
    });

    it('应该在独立实例中存储不同的状态', () => {
      const store1 = createOrchestrationStore();
      const store2 = createOrchestrationStore();
      
      const config1: PageOrchestrationConfig = { id: 'page-1', regions: [] };
      const config2: PageOrchestrationConfig = { id: 'page-2', regions: [] };
      
      store1.getState().setConfig(config1);
      store2.getState().setConfig(config2);
      
      expect(store1.getState().config?.id).toBe('page-1');
      expect(store2.getState().config?.id).toBe('page-2');
    });
  });

  describe('setConfig / getConfig', () => {
    it('应该正确设置和获取配置', () => {
      const store = createOrchestrationStore();
      
      expect(store.getState().config).toBeNull();
      
      store.getState().setConfig(mockConfig);
      
      expect(store.getState().config).toEqual(mockConfig);
    });
  });

  describe('setRequest / getRequest', () => {
    it('应该正确设置和获取请求实例', () => {
      const store = createOrchestrationStore();
      
      expect(store.getState().request).toBeNull();
      
      store.getState().setRequest(mockRequest);
      
      expect(store.getState().request).toBe(mockRequest);
    });
  });

  describe('setPageProps / getPageProps', () => {
    it('应该正确设置和获取页面属性', () => {
      const store = createOrchestrationStore();
      
      expect(store.getState().pageProps).toBeNull();
      
      store.getState().setPageProps(mockPageProps);
      
      expect(store.getState().pageProps).toBe(mockPageProps);
    });
  });

  describe('setResolver / getResolver', () => {
    it('应该正确设置和获取间距解析器', () => {
      const store = createOrchestrationStore();
      
      expect(store.getState().resolver).toBeNull();
      
      store.getState().setResolver(mockResolver);
      
      expect(store.getState().resolver).toBe(mockResolver);
      expect(store.getState().resolver?.('lg')).toBe('24px');
      expect(store.getState().resolver?.('md')).toBe('16px');
    });
  });

  describe('requestCache', () => {
    it('应该正确设置和获取请求状态', () => {
      const store = createOrchestrationStore();
      
      // 初始状态
      expect(store.getState().getRequestState('test-key')).toEqual({
        data: null,
        loading: false,
        error: null
      });
      
      // 设置加载状态
      store.getState().setRequestState('test-key', { loading: true });
      expect(store.getState().getRequestState('test-key')).toEqual({
        data: null,
        loading: true,
        error: null
      });
      
      // 设置数据和完成状态
      store.getState().setRequestState('test-key', { 
        data: { id: 1, name: 'Test' }, 
        loading: false 
      });
      expect(store.getState().getRequestState('test-key')).toEqual({
        data: { id: 1, name: 'Test' },
        loading: false,
        error: null
      });
    });

    it('应该为不同的 key 维护独立的状态', () => {
      const store = createOrchestrationStore();
      
      store.getState().setRequestState('key-1', { data: 'data-1' });
      store.getState().setRequestState('key-2', { data: 'data-2' });
      
      expect(store.getState().getRequestState('key-1').data).toBe('data-1');
      expect(store.getState().getRequestState('key-2').data).toBe('data-2');
    });
  });

  describe('subscribe', () => {
    it('应该监听状态变化', () => {
      const store = createOrchestrationStore();
      const listener = jest.fn();
      
      const unsubscribe = store.subscribe(listener);
      
      // 触发状态变化
      store.getState().setConfig(mockConfig);
      
      expect(listener).toHaveBeenCalled();
      
      // 取消订阅
      unsubscribe();
      
      // 再次触发状态变化
      store.getState().setConfig({ id: 'new-page', regions: [] });
      
      // 由于已经取消订阅，不应该再被调用
      expect(listener).toHaveBeenCalledTimes(1);
    });
  });
});
