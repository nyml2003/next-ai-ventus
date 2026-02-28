/**
 * BlockRenderer - 区块渲染器
 */

import * as React from 'react';
import type { BlockConfig, ModuleConfig, PageProps } from '@ventus/types';
import { useOrchestrationContext } from './index';
import { ModuleContext } from '@ventus/store';
import { usePageProps, useResolver } from '@ventus/store';

interface BlockRendererProps {
  block: BlockConfig;
}

// ==================== 模块渲染器 ====================

interface ModuleRendererProps {
  module: ModuleConfig;
  pageProps: PageProps;
}

const ModuleRenderer: React.FC<ModuleRendererProps> = ({ module, pageProps }) => {
  const { modules } = useOrchestrationContext();
  
  const moduleLoader = modules[module.name];
  
  if (!moduleLoader) {
    // 使用 React.createElement 替代 JSX
    return React.createElement(
      'div',
      { className: 'module-error', style: { padding: '16px', color: 'red' } },
      `Module "${module.name}" not found`
    );
  }
  
  // 提供 ModuleContext，让子组件知道自己在哪个模块中
  const ModuleWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) =>
    React.createElement(
      ModuleContext.Provider,
      { value: { moduleName: module.name } },
      children
    );
  
  // 判断是组件还是懒加载函数
  const isLazyLoader = typeof moduleLoader === 'function' && 
    moduleLoader.prototype === undefined &&
    !('displayName' in moduleLoader || 'name' in moduleLoader);
  
  // 统一的 props，包含 pageProps 和 moduleName
  const componentProps = { 
    pageProps: pageProps,
    __moduleName: module.name  // 注入模块名，供 useModuleData 使用
  };
  
  if (isLazyLoader) {
    const LazyModule = React.lazy(moduleLoader as () => Promise<{ default: React.ComponentType<any> }>);
    
    const WrappedModule: React.FC = () =>
      React.createElement(
        ModuleWrapper,
        {},
        React.createElement(
          React.Suspense,
          { fallback: React.createElement('div', { style: { padding: '16px' } }, `Loading ${module.name}...`) },
          React.createElement(
            'div',
            { className: `module module-${module.name}` },
            React.createElement(LazyModule, componentProps)
          )
        )
      );
    
    return React.createElement(WrappedModule);
  }
  
  // 同步组件
  const Component = moduleLoader as React.ComponentType<{ pageProps: PageProps }>;
  
  return React.createElement(
    ModuleWrapper,
    {},
    React.createElement(
      'div',
      { className: `module module-${module.name}` },
      React.createElement(Component, componentProps)
    )
  );
};

// ==================== 区块渲染器 ====================

export const BlockRenderer: React.FC<BlockRendererProps> = ({ block }) => {
  const resolver = useResolver();
  const pageProps = usePageProps();
  
  // 解析样式
  const style: React.CSSProperties = {
    display: 'flex',
    flexDirection: block.flexDirection,
    gap: block.gap ? resolver(block.gap) : undefined,
    padding: block.padding ? resolver(block.padding) : undefined,
    margin: block.margin ? resolver(block.margin) : undefined,
    justifyContent: block.justifyContent === 'between' ? 'space-between' 
      : block.justifyContent === 'around' ? 'space-around'
      : block.justifyContent,
    alignItems: block.alignItems
  };
  
  // 渲染子元素
  const children = block.children.map((child, index) => {
    if (child.type === 'block') {
      return React.createElement(BlockRenderer, { 
        key: `block-${index}`, 
        block: child 
      });
    }
    
    if (child.type === 'module') {
      return React.createElement(ModuleRenderer, { 
        key: `module-${child.name}`, 
        module: child, 
        pageProps: pageProps 
      });
    }
    
    return null;
  });
  
  // 使用 React.createElement 替代 JSX
  return React.createElement(
    'div',
    {
      className: 'block',
      'data-block-type': 'flex',
      style: style
    },
    children
  );
};
