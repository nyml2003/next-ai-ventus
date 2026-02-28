import type { PageOrchestrationConfig } from '@ventus/types';

/**
 * 首页编排配置
 * 纯数据配置，可迁移到后端
 */
export const homeConfig: PageOrchestrationConfig = {
  id: 'home',
  meta: {
    title: '首页 - Ventus Blog',
    description: 'Ventus 博客首页'
  },
  // BFF 模块列表 - 编排系统会统一请求这些数据
  modules: ['Logo', 'Nav', 'UserAction', 'PostList', 'TagCloud', 'Footer'],
  regions: [
    {
      id: 'header',
      type: 'header',
      padding: 'navPadding',
      block: {
        type: 'block',
        flexDirection: 'row',
        gap: 'navGap',
        justifyContent: 'between',
        alignItems: 'center',
        children: [
          { type: 'module', name: 'Logo' },
          { type: 'module', name: 'Nav' },
          { type: 'module', name: 'UserAction' }
        ]
      }
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
            // 主内容区
            type: 'block',
            flexDirection: 'column',
            gap: 'contentGap',
            padding: 'contentPadding',
            children: [
              { type: 'module', name: 'PostList' }
            ]
          },
          {
            // 侧边栏
            type: 'block',
            flexDirection: 'column',
            gap: 'sidebarGap',
            padding: 'sidebarPadding',
            children: [
              { type: 'module', name: 'TagCloud' }
            ]
          }
        ]
      }
    },
    {
      id: 'footer',
      type: 'footer',
      padding: 'footerPadding',
      block: {
        type: 'block',
        flexDirection: 'row',
        justifyContent: 'center',
        children: [
          { type: 'module', name: 'Footer' }
        ]
      }
    }
  ]
};
