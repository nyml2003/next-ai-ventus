/**
 * @ventus/orchestration - 编排渲染引擎
 */

import * as React from "react";
import type {
  PageOrchestrationConfig,
  SpacingResolver,
  PageProps,
} from "@ventus/types";
import type { RequestInstance } from "@ventus/request";
import { StoreProvider, createPagePropsFromURL, useRequest } from "@ventus/store";
import { RegionRenderer } from "./RegionRenderer";

// ==================== 类型定义 ====================

export type ModuleComponent = React.ComponentType<{ pageProps: PageProps }>;

export interface ModulesRegistry {
  [name: string]:
    | ModuleComponent
    | (() => Promise<{ default: ModuleComponent }>);
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

const OrchestrationContext =
  React.createContext<OrchestrationContextValue | null>(null);

export function useOrchestrationContext(): OrchestrationContextValue {
  const context = React.useContext(OrchestrationContext);
  if (!context) {
    throw new Error(
      "useOrchestrationContext must be used within OrchestrationProvider",
    );
  }
  return context;
}

// ==================== Module Context ====================
// 从 store 包重新导出，避免循环依赖
export { ModuleContext, type ModuleContextValue } from '@ventus/store';

export function useModuleContext(): import('@ventus/store').ModuleContextValue {
  const context = React.useContext(ModuleContext);
  if (!context) {
    throw new Error('useModuleContext must be used within a Module component');
  }
  return context;
}

// ==================== 页面渲染器 ====================

interface PageRendererProps {
  config: PageOrchestrationConfig;
  modules: ModulesRegistry;
  resolver: SpacingResolver;
}

const PageRenderer: React.FC<PageRendererProps> = ({
  config,
  modules,
  resolver,
}) => {
  // 从 URL 获取参数
  const pageProps = createPagePropsFromURL();
  
  // 获取 BFF 模块列表（如果配置中没有，则使用空数组）
  const bffModules = config.modules || [];
  
  // 使用 ref 确保依赖数组稳定，避免重复请求
  const requestConfig = React.useRef({
    page: config.id,
    modules: bffModules,
    params: pageProps.query,
  }).current;
  
  // 发起聚合请求（如果配置了 modules）
  const { loading, error } = useRequest({
    ...requestConfig,
    deps: [], // 空依赖，只请求一次
  });

  // 加载中状态
  if (loading) {
    return React.createElement(
      "div",
      { className: "page-loading", style: { padding: '40px', textAlign: 'center' } },
      "加载中..."
    );
  }

  // 错误状态
  if (error) {
    return React.createElement(
      "div",
      { className: "page-error", style: { padding: '40px', textAlign: 'center', color: 'red' } },
      React.createElement("p", null, `加载失败: ${error.message}`),
      React.createElement(
        "button",
        { onClick: () => window.location.reload() },
        "重试"
      )
    );
  }

  // 使用 React.createElement 替代 JSX
  return React.createElement(
    OrchestrationContext.Provider,
    { value: { modules, resolver } },
    React.createElement(
      "div",
      { className: "page", "data-page-id": config.id },
      config.regions.map((region) =>
        React.createElement(RegionRenderer, {
          key: region.id,
          region: region,
        }),
      ),
    ),
  );
};

// ==================== 创建编排系统 ====================

export function createOrchestration({
  config,
  request,
  modules,
  resolver,
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
        pageProps: pageProps,
        modules: modules as Record<string, unknown>,
      },
      React.createElement(PageRenderer, {
        config: config,
        modules: modules,
        resolver: resolver,
      }),
    );
  };

  return { Renderer };
}

export default createOrchestration;
