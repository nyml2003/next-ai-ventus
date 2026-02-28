/**
 * Ventus 编排系统类型定义
 * 支持 Page - Region - Block - Module 四层结构
 */

// ==================== 配置层（纯数据，JSON-friendly） ====================

/** Spacing 解析器函数 */
export type SpacingResolver = (token: string) => string;

/** 模块配置 */
export interface ModuleConfig {
  type: 'module';
  /** 模块名称，对应注册表中的 key */
  name: string;
}

/** 区块布局方向 */
export type FlexDirection = 'row' | 'column' | 'row-reverse' | 'column-reverse';

/** 主轴对齐方式 */
export type FlexJustify = 'start' | 'center' | 'end' | 'between' | 'around';

/** 交叉轴对齐方式 */
export type FlexAlign = 'start' | 'center' | 'end' | 'stretch';

/** 区块配置 - 仅支持 Flex 布局 */
export interface BlockConfig {
  type: 'block';
  /** Flex 方向 */
  flexDirection: FlexDirection;
  /** 子元素间距 - 语义化 token */
  gap?: string;
  /** 内边距 - 语义化 token */
  padding?: string;
  /** 外边距 - 语义化 token */
  margin?: string;
  /** 主轴对齐 */
  justifyContent?: FlexJustify;
  /** 交叉轴对齐 */
  alignItems?: FlexAlign;
  /** 子元素（区块或模块） */
  children: (BlockConfig | ModuleConfig)[];
}

/** 区域类型 */
export type RegionType = 'header' | 'content' | 'footer' | 'sidebar';

/** 区域配置 */
export interface RegionConfig {
  id: string;
  type: RegionType;
  /** 区域内边距 - 语义化 token */
  padding?: string;
  /** 区域外边距 - 语义化 token */
  margin?: string;
  /** 根区块 */
  block: BlockConfig;
}

/** 页面编排配置 */
export interface PageOrchestrationConfig {
  id: string;
  /** 页面元信息 */
  meta?: {
    title?: string;
    description?: string;
  };
  /** 区域列表 */
  regions: RegionConfig[];
}

// ==================== 运行时层 ====================

/** 带数据的模块 */
export interface ModuleData extends ModuleConfig {
  data: unknown;
}

/** 带数据的区块 */
export interface BlockData extends Omit<BlockConfig, 'children'> {
  children: (BlockData | ModuleData)[];
  /** 解析后的样式值 */
  computedStyle?: {
    gap?: string;
    padding?: string;
    margin?: string;
  };
}

/** 带数据的区域 */
export interface RegionData extends Omit<RegionConfig, 'block'> {
  block: BlockData;
  computedStyle?: {
    padding?: string;
    margin?: string;
  };
}

/** 编排渲染数据 */
export interface OrchestrationData {
  page: {
    id: string;
    meta?: PageOrchestrationConfig['meta'];
  };
  regions: RegionData[];
}

// ==================== 注册表 ====================

/** 模块加载函数 */
export type ModuleLoader = () => Promise<{ default: React.ComponentType<any> }>;

/** 页面级注册表接口 */
export interface PageRegistry {
  /** 设置 Spacing 解析器 */
  setResolver: (resolver: SpacingResolver) => void;
  /** 获取 Spacing 解析器 */
  getResolver: () => SpacingResolver;
  /** 注册模块 */
  registerModule: (name: string, loader: ModuleLoader) => void;
  /** 获取模块 */
  getModule: (name: string) => ModuleLoader | undefined;
  /** 获取所有模块名称 */
  getModuleNames: () => string[];
}

// ==================== PageProps ====================

/** 跳链参数类型 */
export interface PageProps {
  /** 获取路由参数 /post/:slug -> getParam('slug') */
  getParam: (key: string) => string | undefined;
  /** 获取查询参数 ?page=2 -> getQuery('page') */
  getQuery: (key: string) => string | undefined;
  /** 获取所有路由参数 */
  params: Record<string, string>;
  /** 获取所有查询参数 */
  query: Record<string, string>;
}
