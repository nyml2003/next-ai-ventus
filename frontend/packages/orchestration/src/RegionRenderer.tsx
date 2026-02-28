/**
 * RegionRenderer - 区域渲染器
 */

import * as React from 'react';
import type { RegionConfig } from '@ventus/types';
import { useResolver } from '@ventus/store';
import { BlockRenderer } from './BlockRenderer';

interface RegionRendererProps {
  region: RegionConfig;
}

export const RegionRenderer: React.FC<RegionRendererProps> = ({ region }) => {
  const resolver = useResolver();
  
  // 解析区域样式
  const style: React.CSSProperties = {
    padding: region.padding ? resolver(region.padding) : undefined,
    margin: region.margin ? resolver(region.margin) : undefined
  };
  
  // 使用 React.createElement 替代 JSX
  return React.createElement(
    'section',
    {
      className: `region region-${region.type}`,
      'data-region-id': region.id,
      'data-region-type': region.type,
      style: style
    },
    React.createElement(BlockRenderer, { block: region.block })
  );
};
