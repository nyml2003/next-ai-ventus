import * as React from 'react';
import './style.css';

// ==================== Title / Heading ====================

export interface TitleProps {
  level?: 1 | 2 | 3 | 4 | 5 | 6;
  children: React.ReactNode;
  className?: string;
  style?: React.CSSProperties;
}

export const Title: React.FC<TitleProps> = ({
  level = 1,
  children,
  className = '',
  style,
}) => {
  const Tag = `h${level}` as keyof JSX.IntrinsicElements;
  return (
    <Tag className={`vent-title vent-title--level-${level} ${className}`} style={style}>
      {children}
    </Tag>
  );
};

// ==================== Text / Paragraph ====================

export type TextVariant = 'body' | 'secondary' | 'tertiary' | 'caption' | 'code';
export type TextSize = 'xs' | 'sm' | 'base' | 'lg' | 'xl';

export interface TextProps {
  children: React.ReactNode;
  variant?: TextVariant;
  size?: TextSize;
  bold?: boolean;
  ellipsis?: boolean;
  className?: string;
  style?: React.CSSProperties;
  as?: 'p' | 'span' | 'div' | 'label';
}

export const Text: React.FC<TextProps> = ({
  children,
  variant = 'body',
  size = 'base',
  bold = false,
  ellipsis = false,
  className = '',
  style,
  as: Component = 'span',
}) => {
  const classes = [
    'vent-text',
    `vent-text--${variant}`,
    `vent-text--${size}`,
    bold && 'vent-text--bold',
    ellipsis && 'vent-text--ellipsis',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <Component className={classes} style={style}>
      {children}
    </Component>
  );
};

// ==================== Link ====================

export interface LinkProps {
  href: string;
  children: React.ReactNode;
  external?: boolean;
  underline?: boolean;
  className?: string;
  style?: React.CSSProperties;
  onClick?: (e: React.MouseEvent<HTMLAnchorElement>) => void;
}

export const Link: React.FC<LinkProps> = ({
  href,
  children,
  external = false,
  underline = true,
  className = '',
  style,
  onClick,
}) => {
  const classes = [
    'vent-link',
    underline && 'vent-link--underline',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <a
      href={href}
      className={classes}
      style={style}
      onClick={onClick}
      {...(external && { target: '_blank', rel: 'noopener noreferrer' })}
    >
      {children}
    </a>
  );
};

// ==================== Container / Box ====================

export interface BoxProps {
  children: React.ReactNode;
  padding?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  margin?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  className?: string;
  style?: React.CSSProperties;
}

export const Box: React.FC<BoxProps> = ({
  children,
  padding = 'none',
  margin = 'none',
  className = '',
  style,
}) => {
  const classes = [
    'vent-box',
    padding !== 'none' && `vent-box--padding-${padding}`,
    margin !== 'none' && `vent-box--margin-${margin}`,
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classes} style={style}>
      {children}
    </div>
  );
};

// ==================== Flex / Stack ====================

export type FlexDirection = 'row' | 'column' | 'row-reverse' | 'column-reverse';
export type FlexJustify = 'start' | 'center' | 'end' | 'between' | 'around' | 'evenly';
export type FlexAlign = 'start' | 'center' | 'end' | 'stretch' | 'baseline';
export type FlexGap = 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';

export interface FlexProps {
  children: React.ReactNode;
  direction?: FlexDirection;
  justify?: FlexJustify;
  align?: FlexAlign;
  gap?: FlexGap;
  wrap?: boolean;
  className?: string;
  style?: React.CSSProperties;
}

export const Flex: React.FC<FlexProps> = ({
  children,
  direction = 'row',
  justify = 'start',
  align = 'stretch',
  gap = 'none',
  wrap = false,
  className = '',
  style,
}) => {
  const classes = [
    'vent-flex',
    `vent-flex--direction-${direction}`,
    `vent-flex--justify-${justify}`,
    `vent-flex--align-${align}`,
    gap !== 'none' && `vent-flex--gap-${gap}`,
    wrap && 'vent-flex--wrap',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classes} style={style}>
      {children}
    </div>
  );
};

// ==================== Stack (sugar for Flex column) ====================

export interface StackProps {
  children: React.ReactNode;
  gap?: FlexGap;
  align?: FlexAlign;
  className?: string;
  style?: React.CSSProperties;
}

export const Stack: React.FC<StackProps> = ({
  children,
  gap = 'md',
  align = 'stretch',
  className = '',
  style,
}) => {
  return (
    <Flex
      direction="column"
      gap={gap}
      align={align}
      className={className}
      style={style}
    >
      {children}
    </Flex>
  );
};

// ==================== Inline (sugar for Flex row) ====================

export interface InlineProps {
  children: React.ReactNode;
  gap?: FlexGap;
  align?: FlexAlign;
  className?: string;
  style?: React.CSSProperties;
}

export const Inline: React.FC<InlineProps> = ({
  children,
  gap = 'sm',
  align = 'center',
  className = '',
  style,
}) => {
  return (
    <Flex
      direction="row"
      gap={gap}
      align={align}
      wrap
      className={className}
      style={style}
    >
      {children}
    </Flex>
  );
};
