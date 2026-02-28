/**
 * useRequest Hook 测试
 */

import * as React from 'react';
import { renderHook, waitFor } from '@testing-library/react';
import { useRequest, StoreProvider } from '../index';
import type { PageOrchestrationConfig, PageProps } from '@ventus/types';
import type { RequestInstance } from '@ventus/request';

// Mock request
const createMockRequest = (responses: Record<string, unknown>): RequestInstance => ({
  call: jest.fn(async ({ scene, params }) => {
    const key = `${scene}:${JSON.stringify(params || {})}`;
    if (responses[key] instanceof Error) {
      throw responses[key];
    }
    return responses[key];
  }),
  get: jest.fn(),
  post: jest.fn(),
});

const mockPageProps: PageProps = {
  getParam: () => undefined,
  getQuery: () => undefined,
  params: {},
  query: {},
};

const mockConfig: PageOrchestrationConfig = {
  id: 'test',
  regions: [],
};

const createWrapper = (request: RequestInstance) => {
  return function Wrapper({ children }: { children: React.ReactNode }) {
    return React.createElement(StoreProvider, {
      config: mockConfig,
      request,
      resolver: () => '0px',
      pageProps: mockPageProps,
    }, children);
  };
};

describe('useRequest', () => {
  it('应该正确获取数据', async () => {
    const mockData = { items: [1, 2, 3] };
    const request = createMockRequest({
      'post.list:{"page":1}': mockData,
    });

    const { result } = renderHook(
      () => useRequest({ scene: 'post.list', params: { page: 1 } }),
      { wrapper: createWrapper(request) }
    );

    // 初始状态是 loading
    expect(result.current.loading).toBe(true);
    expect(result.current.data).toBeNull();

    // 等待请求完成
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.data).toEqual(mockData);
    expect(result.current.error).toBeNull();
  });

  it('应该处理请求错误', async () => {
    const request = createMockRequest({
      'post.list:{}': new Error('Network error'),
    });

    const { result } = renderHook(
      () => useRequest({ scene: 'post.list' }),
      { wrapper: createWrapper(request) }
    );

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBeInstanceOf(Error);
    expect(result.current.error?.message).toBe('Network error');
  });

  it('应该缓存相同请求', async () => {
    const mockData = { items: [] };
    const request = createMockRequest({
      'tag.list:{}': mockData,
    });

    // 第一次调用
    const { result: result1 } = renderHook(
      () => useRequest({ scene: 'tag.list' }),
      { wrapper: createWrapper(request) }
    );

    await waitFor(() => {
      expect(result1.current.loading).toBe(false);
    });

    // 第二次调用相同请求
    const { result: result2 } = renderHook(
      () => useRequest({ scene: 'tag.list' }),
      { wrapper: createWrapper(request) }
    );

    // 应该立即返回缓存数据，不需要 loading
    expect(result2.current.data).toEqual(mockData);
    expect(result2.current.loading).toBe(false);

    // request.call 应该只被调用一次
    expect(request.call).toHaveBeenCalledTimes(1);
  });

  it('refetch 应该重新获取数据', async () => {
    let callCount = 0;
    const request = createMockRequest({
      'post.list:{}': new Proxy({}, {
        get() {
          callCount++;
          return { count: callCount };
        },
      }),
    });

    const { result } = renderHook(
      () => useRequest({ scene: 'post.list' }),
      { wrapper: createWrapper(request) }
    );

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(callCount).toBe(1);

    // 调用 refetch
    result.current.refetch();

    await waitFor(() => {
      expect(result.current.data).toEqual({ count: 2 });
    });

    expect(callCount).toBe(2);
  });

  it('不同参数应该触发不同请求', async () => {
    const request = createMockRequest({
      'post.list:{"page":1}': { items: [1] },
      'post.list:{"page":2}': { items: [2] },
    });

    const { result: result1 } = renderHook(
      () => useRequest({ scene: 'post.list', params: { page: 1 } }),
      { wrapper: createWrapper(request) }
    );

    const { result: result2 } = renderHook(
      () => useRequest({ scene: 'post.list', params: { page: 2 } }),
      { wrapper: createWrapper(request) }
    );

    await waitFor(() => {
      expect(result1.current.data).toEqual({ items: [1] });
      expect(result2.current.data).toEqual({ items: [2] });
    });

    expect(request.call).toHaveBeenCalledTimes(2);
  });
});
