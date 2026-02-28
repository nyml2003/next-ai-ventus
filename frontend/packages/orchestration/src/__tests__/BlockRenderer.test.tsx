/**
 * BlockRenderer 单元测试
 */

import * as React from 'react';
import { render, screen } from '@testing-library/react';
import { BlockRenderer } from '../BlockRenderer';
import { createOrchestration } from '../index';
import { createRequest } from '@ventus/request';
import type { PageOrchestrationConfig, PageProps } from '@ventus/types';

// 测试模块组件
const TestModule = ({ pageProps }: { pageProps: PageProps }) => {
  return React.createElement('div', { 'data-testid': 'test-module' }, 'Test Module');
};

const ErrorModule = () => {
  throw new Error('Module Error');
};

describe('BlockRenderer', () => {
  const mockResolver = (token: string) => {
    const map: Record<string, string> = {
      gap: '16px',
      padding: '24px',
    };
    return map[token] || '0px';
  };

  const createTestOrchestration = (config: PageOrchestrationConfig) => {
    return createOrchestration({
      config,
      request: createRequest(),
      modules: {
        TestModule,
        ErrorModule,
      },
      resolver: mockResolver,
    });
  };

  it('应该正确渲染简单区块', () => {
    const config: PageOrchestrationConfig = {
      id: 'test',
      regions: [{
        id: 'content',
        type: 'content',
        block: {
          type: 'block',
          flexDirection: 'row',
          gap: 'gap',
          padding: 'padding',
          children: [
            { type: 'module', name: 'TestModule' },
          ],
        },
      }],
    };

    const orchestration = createTestOrchestration(config);
    const { container } = render(React.createElement(orchestration.Renderer));

    // 验证区块样式
    const block = container.querySelector('.block');
    expect(block).toHaveStyle({
      display: 'flex',
      flexDirection: 'row',
      gap: '16px',
      padding: '24px',
    });
  });

  it('应该正确渲染嵌套区块', () => {
    const config: PageOrchestrationConfig = {
      id: 'test',
      regions: [{
        id: 'content',
        type: 'content',
        block: {
          type: 'block',
          flexDirection: 'row',
          gap: 'gap',
          children: [
            {
              type: 'block',
              flexDirection: 'column',
              gap: 'gap',
              children: [
                { type: 'module', name: 'TestModule' },
              ],
            },
          ],
        },
      }],
    };

    const orchestration = createTestOrchestration(config);
    const { container } = render(React.createElement(orchestration.Renderer));

    // 验证嵌套结构
    const blocks = container.querySelectorAll('.block');
    expect(blocks).toHaveLength(2);
    
    // 外层是 row
    expect(blocks[0]).toHaveStyle({ flexDirection: 'row' });
    // 内层是 column
    expect(blocks[1]).toHaveStyle({ flexDirection: 'column' });
  });

  it('应该正确传递 pageProps 给模块', () => {
    // 设置 URL 参数
    window.location.href = 'http://localhost:3000/?tag=tech';
    window.location.search = '?tag=tech';

    const PagePropsDisplay = ({ pageProps }: { pageProps: PageProps }) => {
      return React.createElement(
        'div',
        { 'data-testid': 'page-props' },
        pageProps.getQuery('tag') || 'no-tag'
      );
    };

    const config: PageOrchestrationConfig = {
      id: 'test',
      regions: [{
        id: 'content',
        type: 'content',
        block: {
          type: 'block',
          flexDirection: 'row',
          children: [
            { type: 'module', name: 'PagePropsDisplay' },
          ],
        },
      }],
    };

    const orchestration = createOrchestration({
      config,
      request: createRequest(),
      modules: { PagePropsDisplay },
      resolver: mockResolver,
    });

    render(React.createElement(orchestration.Renderer));

    expect(screen.getByTestId('page-props')).toHaveTextContent('tech');
  });

  it('应该显示模块未找到的错误', () => {
    const config: PageOrchestrationConfig = {
      id: 'test',
      regions: [{
        id: 'content',
        type: 'content',
        block: {
          type: 'block',
          flexDirection: 'row',
          children: [
            { type: 'module', name: 'NonExistentModule' },
          ],
        },
      }],
    };

    const orchestration = createOrchestration({
      config,
      request: createRequest(),
      modules: {}, // 空模块表
      resolver: mockResolver,
    });

    const { container } = render(React.createElement(orchestration.Renderer));

    expect(container.querySelector('.module-error')).toHaveTextContent(
      'Module "NonExistentModule" not found'
    );
  });

  it('应该正确解析 justifyContent', () => {
    const config: PageOrchestrationConfig = {
      id: 'test',
      regions: [{
        id: 'content',
        type: 'content',
        block: {
          type: 'block',
          flexDirection: 'row',
          justifyContent: 'between',
          children: [],
        },
      }],
    };

    const orchestration = createTestOrchestration(config);
    const { container } = render(React.createElement(orchestration.Renderer));

    const block = container.querySelector('.block');
    expect(block).toHaveStyle({ justifyContent: 'space-between' });
  });
});
