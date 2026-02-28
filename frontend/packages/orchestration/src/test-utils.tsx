/**
 * 编排系统测试工具库
 */

import * as React from 'react';
import { render as rtlRender } from '@testing-library/react';
import { createOrchestration } from './index';
import { createRequest } from '@ventus/request';
import type { PageOrchestrationConfig, ModuleConfig } from '@ventus/types';

export interface RenderWithOrchestrationOptions {
  /** 页面配置 */
  config?: Partial<PageOrchestrationConfig>;
  /** 要测试的模块 */
  modules: Record<string, React.ComponentType<any>>;
  /** Spacing 解析器 */
  resolver?: (token: string) => string;
  /** 模拟的 API 响应 */
  mockResponses?: Record<string, unknown>;
  /** 初始 URL */
  url?: string;
}

/**
 * 渲染带编排系统的组件
 * 
 * @example
 * ```tsx
 * const { container } = renderWithOrchestration({
 *   modules: { MyModule: MyComponent },
 *   config: {
 *     regions: [{
 *       id: 'content',
 *       type: 'content',
 *       block: {
 *         type: 'block',
 *         flexDirection: 'row',
 *         children: [{ type: 'module', name: 'MyModule' }]
 *       }
 *     }]
 *   }
 * });
 * ```
 */
export function renderWithOrchestration({
  config = {},
  modules,
  resolver = (token) => ({ gap: '16px', padding: '24px' }[token] || '0px'),
  mockResponses = {},
  url = 'http://localhost:3000/',
}: RenderWithOrchestrationOptions) {
  // 设置 URL
  window.location.href = url;
  window.location.search = new URL(url).search;

  // 创建模拟 request
  const request = createRequest();
  if (Object.keys(mockResponses).length > 0) {
    jest.spyOn(request, 'call').mockImplementation(async ({ scene, params }) => {
      const key = `${scene}:${JSON.stringify(params || {})}`;
      const response = mockResponses[key];
      if (response === undefined) {
        throw new Error(`Unexpected API call: ${key}`);
      }
      return response;
    });
  }

  // 合并默认配置
  const fullConfig: PageOrchestrationConfig = {
    id: 'test',
    regions: [],
    ...config,
  };

  // 创建编排系统
  const orchestration = createOrchestration({
    config: fullConfig,
    request,
    modules,
    resolver,
  });

  return rtlRender(React.createElement(orchestration.Renderer));
}

/**
 * 创建简单的测试模块配置
 */
export function createModuleConfig(name: string): ModuleConfig {
  return { type: 'module', name };
}

/**
 * 创建简单的区块配置
 */
export function createBlockConfig(
  children: ModuleConfig[],
  options: Partial<{ flexDirection: 'row' | 'column'; gap: string; padding: string }> = {}
) {
  return {
    type: 'block' as const,
    flexDirection: options.flexDirection || 'row',
    gap: options.gap,
    padding: options.padding,
    children,
  };
}

/**
 * 等待所有请求完成
 */
export async function waitForRequests() {
  await new Promise((resolve) => setTimeout(resolve, 0));
}
