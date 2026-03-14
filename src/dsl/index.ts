import type { SegmentNode, EnabledFn, StyleAttrs } from '../types.js';

interface BaseProps extends Partial<StyleAttrs> {
  enabled?: boolean | EnabledFn;
}

interface PwdProps extends BaseProps {
  style?: 'name' | 'path' | 'smart';
}

interface LiteralProps {
  text: string;
}

interface SepProps extends Partial<StyleAttrs> {
  char?: string;
}

interface CompositeProps extends BaseProps {}

function extractStyle(props: Partial<StyleAttrs>): StyleAttrs | undefined {
  const style: StyleAttrs = {};
  let hasStyle = false;
  for (const key of ['color', 'bold', 'italic', 'prefix', 'suffix'] as const) {
    if (props[key] !== undefined) {
      (style as any)[key] = props[key];
      hasStyle = true;
    }
  }
  return hasStyle ? style : undefined;
}

export function Literal(props: LiteralProps): SegmentNode {
  return { type: 'literal', props: { text: props.text } };
}

export function Sep(props: SepProps = {}): SegmentNode {
  const { char, ...styleProps } = props;
  return {
    type: 'sep',
    props: char !== undefined ? { char } : undefined,
    style: extractStyle(styleProps),
  };
}

export function Pwd(props: PwdProps = {}): SegmentNode {
  const { style: variant = 'smart', enabled, ...styleProps } = props;
  return {
    type: `pwd.${variant}`,
    provider: 'pwd',
    enabled,
    style: extractStyle(styleProps),
  };
}

export function Branch(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'git.branch', provider: 'git', enabled, style: extractStyle(styleProps) };
}

export function Insertions(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'git.insertions', provider: 'git', enabled, style: extractStyle(styleProps) };
}

export function Deletions(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'git.deletions', provider: 'git', enabled, style: extractStyle(styleProps) };
}

export function Tokens(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'context.tokens', provider: 'context', enabled, style: extractStyle(styleProps) };
}

export function Size(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'context.size', provider: 'context', enabled, style: extractStyle(styleProps) };
}

export function Percent(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'context.percent', provider: 'context', enabled, style: extractStyle(styleProps) };
}

export function Git(props: CompositeProps = {}): (children: () => SegmentNode[]) => SegmentNode {
  const { enabled, ...styleProps } = props;
  return (children) => ({
    type: 'git',
    enabled,
    style: extractStyle(styleProps),
    children: children(),
  });
}

export function Context(props: CompositeProps = {}): (children: () => SegmentNode[]) => SegmentNode {
  const { enabled, ...styleProps } = props;
  return (children) => ({
    type: 'context',
    enabled,
    style: extractStyle(styleProps),
    children: children(),
  });
}

export function Group(props: CompositeProps = {}): (children: () => SegmentNode[]) => SegmentNode {
  const { enabled, ...styleProps } = props;
  return (children) => ({
    type: 'group',
    enabled,
    style: extractStyle(styleProps),
    children: children(),
  });
}

export function StatusLine(children: () => SegmentNode[]): SegmentNode[] {
  return children();
}
