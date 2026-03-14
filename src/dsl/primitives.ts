import type { SegmentNode, EnabledFn, StyleAttrs } from '../types.js';

export interface BaseProps extends Partial<StyleAttrs> {
  enabled?: boolean | EnabledFn;
}

export interface CompositeProps extends BaseProps {}

interface LiteralProps {
  text: string;
}

export function extractStyle(props: Partial<StyleAttrs>): StyleAttrs | undefined {
  const style: StyleAttrs = {};
  let hasStyle = false;
  for (const key of ['color', 'bold', 'italic', 'prefix', 'suffix'] as const) {
    if (props[key] !== undefined) {
      const styleRecord = style as Record<string, unknown>;
      styleRecord[key] = props[key];
      hasStyle = true;
    }
  }
  return hasStyle ? style : undefined;
}

export function Literal(props: LiteralProps): SegmentNode {
  return { type: 'literal', props: { text: props.text } };
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
