/**
 * @ventus/orchestration - 编排渲染引擎
 */

import * as React from 'react';
import type { 
  PageOrchestrationConfig, 
  SpacingResolver,
  PageProps 
} from '@ventus/types';
import type { RequestInstance } from '@ventus/request';
import { StoreProvider, createPagePropsFromURL } from '@ventus/store';
import { BlockRenderer } from './BlockRenderer';
import { RegionRenderer } from './RegionRenderer';

// ==================== 类型定义 ====================

export type ModuleComponent = React.ComponentType<{ pageProps: PageProps }>;

export interface ModulesRegistry {
  [name: string]: ModuleComponent | (() => Promise<{ default: ModuleComponent }>);
}

export interface CreateOrchestrationOptions {
  config: PageOrchestrationConfig;
  request: RequestInstance;
  modules: ModulesRegistry;
  resolver: SpacingResolver;
}

export interface OrchestrationInstance {
  Renderer: React.FC;
}

// ==================== Context 传递模块注册表 ====================

interface OrchestrationContextValue {
  modules: ModulesRegistry;
  resolver: SpacingResolver;
}

const OrchestrationContext = React.createContext<OrchestrationContextValue | null>(null);

export function useOrchestrationContext(): OrchestrationContextValue {
  const context = React.useContext(OrchestrationContext);
  if (!context) {
    throw new Error('useOrchestrationContext must be used within OrchestrationProvider');
  }
  return context;
}

// ==================== 页面渲染器 ====================

interface PageRendererProps {
  config: PageOrchestrationConfig;
  modules: ModulesRegistry;
  resolver: SpacingResolver;
}

const PageRenderer: React.FC<PageRendererProps> = ({ config, modules, resolver }) => {
  // 使用 React.createElement 替代 JSX
  return React.createElement(
    OrchestrationContext.Provider,
    { value: { modules, resolver } },
    React.createElement(
      'div',
      { className: 'page', 'data-page-id': config.id },
      config.regions.map((region) =>
        React.createElement(RegionRenderer, {
          key: region.id,
          region: region
        })
      )
    )
  );
};

// ==================== 创建编排系统 ====================

export function createOrchestration({
  config,
  request,
  modules,
  resolver
}: CreateOrchestrationOptions): OrchestrationInstance {
  // 解析 URL 参数
  const pageProps = createPagePropsFromURL();
  
  const Renderer: React.FC = () => {
    // 使用 React.createElement 替代 JSX
    return React.createElement(
      StoreProvider,
      {
        config: config,
        request: request,
        resolver: resolver,
        pageProps: pageProps
      },
      React.createElement(PageRenderer, {
        config: config,
        modules: modules,
        resolver: resolver
      })
    );
  };
  
  return { Renderer };
}

export default createOrchestration;
