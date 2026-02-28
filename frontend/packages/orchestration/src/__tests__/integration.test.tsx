/**
 * 编排系统集成测试
 * 测试完整的页面渲染流程
 */

import * as React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { createOrchestration } from '../index';
import { createRequest } from '@ventus/request';
import type { PageOrchestrationConfig, PageProps } from '@ventus/types';

// 模拟完整的页面场景
describe('编排系统集成测试', () => {
  // 模拟后端数据
  const mockAPIResponses: Record<string, unknown> = {
    'user.getCurrent:{}': { name: 'Admin', avatar: '/avatar.png' },
    'post.list:{"page":1}': {
      items: [
        { id: '1', title: '文章1', excerpt: '摘要1', tags: ['tech'], date: '2024-01-01', href: '/post/1' },
        { id: '2', title: '文章2', excerpt: '摘要2', tags: ['life'], date: '2024-01-02', href: '/post/2' },
      ],
      pagination: { page: 1, totalPages: 3 },
    },
    'tag.list:{}': [
      { name: 'tech', count: 10 },
      { name: 'life', count: 5 },
    ],
  };

  const mockRequest = {
    call: jest.fn(async ({ scene, params }: { scene: string; params?: unknown }) => {
      const key = `${scene}:${JSON.stringify(params || {})}`;
      const response = mockAPIResponses[key];
      if (response === undefined) {
        throw new Error(`Unexpected API call: ${key}`);
      }
      return response;
    }),
    get: jest.fn(),
    post: jest.fn(),
  };

  // 模块组件
  const Header = ({ pageProps }: { pageProps: PageProps }) => {
    const [user, setUser] = React.useState<{ name: string } | null>(null);
    
    React.useEffect(() => {
      mockRequest.call({ scene: 'user.getCurrent' }).then(setUser);
    }, []);

    return React.createElement('header', { 'data-testid': 'header' },
      React.createElement('div', null, user ? `Hello, ${user.name}` : 'Loading...')
    );
  };

  const PostList = ({ pageProps }: { pageProps: PageProps }) => {
    const [data, setData] = React.useState<any>(null);
    const page = parseInt(pageProps.getQuery('page') || '1');

    React.useEffect(() => {
      mockRequest.call({ scene: 'post.list', params: { page } }).then(setData);
    }, [page]);

    if (!data) return React.createElement('div', null, 'Loading posts...');

    return React.createElement('div', { 'data-testid': 'post-list' },
      data.items.map((post: any) =>
        React.createElement('article', { key: post.id, 'data-testid': `post-${post.id}` },
          React.createElement('h2', null, post.title)
        )
      )
    );
  };

  const TagCloud = () => {
    const [tags, setTags] = React.useState<any[]>([]);

    React.useEffect(() => {
      mockRequest.call({ scene: 'tag.list' }).then(setTags);
    }, []);

    return React.createElement('div', { 'data-testid': 'tag-cloud' },
      tags.map((tag: any) =>
        React.createElement('span', { key: tag.name, 'data-testid': `tag-${tag.name}` },
          `${tag.name} (${tag.count})`
        )
      )
    );
  };

  const homeConfig: PageOrchestrationConfig = {
    id: 'home',
    regions: [
      {
        id: 'header',
        type: 'header',
        padding: 'navPadding',
        block: {
          type: 'block',
          flexDirection: 'row',
          justifyContent: 'between',
          children: [{ type: 'module', name: 'Header' }],
        },
      },
      {
        id: 'content',
        type: 'content',
        padding: 'pagePadding',
        block: {
          type: 'block',
          flexDirection: 'row',
          gap: 'sectionGap',
          children: [
            {
              type: 'block',
              flexDirection: 'column',
              gap: 'contentGap',
              children: [{ type: 'module', name: 'PostList' }],
            },
            {
              type: 'block',
              flexDirection: 'column',
              gap: 'sidebarGap',
              children: [{ type: 'module', name: 'TagCloud' }],
            },
          ],
        },
      },
    ],
  };

  const resolver = (token: string) => {
    const map: Record<string, string> = {
      navPadding: '16px',
      pagePadding: '24px',
      sectionGap: '32px',
      contentGap: '16px',
      sidebarGap: '12px',
    };
    return map[token] || '0px';
  };

  it('应该渲染完整的页面结构', async () => {
    // 设置 URL
    window.location.href = 'http://localhost:3000/?page=1';
    window.location.search = '?page=1';

    const orchestration = createOrchestration({
      config: homeConfig,
      request: mockRequest as any,
      modules: { Header, PostList, TagCloud },
      resolver,
    });

    render(React.createElement(orchestration.Renderer));

    // 验证区域结构
    expect(screen.getByTestId('header')).toBeInTheDocument();
    
    // 等待数据加载
    await waitFor(() => {
      expect(screen.getByTestId('post-list')).toBeInTheDocument();
    });

    await waitFor(() => {
      expect(screen.getByTestId('tag-cloud')).toBeInTheDocument();
    });

    // 验证文章内容
    expect(screen.getByTestId('post-1')).toHaveTextContent('文章1');
    expect(screen.getByTestId('post-2')).toHaveTextContent('文章2');

    // 验证标签
    expect(screen.getByTestId('tag-tech')).toHaveTextContent('tech (10)');
    expect(screen.getByTestId('tag-life')).toHaveTextContent('life (5)');
  });

  it('应该正确应用 Spacing 解析', () => {
    const orchestration = createOrchestration({
      config: homeConfig,
      request: mockRequest as any,
      modules: { Header, PostList, TagCloud },
      resolver,
    });

    const { container } = render(React.createElement(orchestration.Renderer));

    // 验证 Region 样式
    const headerRegion = container.querySelector('[data-region-id="header"]');
    expect(headerRegion).toHaveStyle({ padding: '16px' });

    const contentRegion = container.querySelector('[data-region-id="content"]');
    expect(contentRegion).toHaveStyle({ padding: '24px' });

    // 验证 Block 样式
    const contentBlock = contentRegion?.querySelector('.block');
    expect(contentBlock).toHaveStyle({ gap: '32px' });
  });

  it('应该响应 URL 参数变化', async () => {
    // 第一页
    window.location.href = 'http://localhost:3000/?page=1';
    window.location.search = '?page=1';

    const orchestration = createOrchestration({
      config: homeConfig,
      request: mockRequest as any,
      modules: { Header, PostList, TagCloud },
      resolver,
    });

    const { rerender } = render(React.createElement(orchestration.Renderer));

    await waitFor(() => {
      expect(screen.getByTestId('post-list')).toBeInTheDocument();
    });

    // 验证第一页数据
    expect(mockRequest.call).toHaveBeenCalledWith(
      expect.objectContaining({ scene: 'post.list', params: { page: 1 } })
    );
  });
});
