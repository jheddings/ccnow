import type { SegmentNode } from './types.js';

interface JsonSegment {
  segment: string;
  style?: Record<string, unknown>;
  props?: Record<string, unknown>;
  enabled?: boolean;
  children?: JsonSegment[];
}

interface JsonConfig {
  segments?: JsonSegment[];
}

const NO_PROVIDER = new Set(['literal', 'sep']);

function inferProvider(type: string): string | undefined {
  if (NO_PROVIDER.has(type)) return undefined;
  const dotIndex = type.indexOf('.');
  if (dotIndex > 0) return type.slice(0, dotIndex);
  return undefined;
}

function mapSegment(json: JsonSegment): SegmentNode {
  const node: SegmentNode = {
    type: json.segment,
  };

  const provider = inferProvider(json.segment);
  if (provider) node.provider = provider;

  if (json.style) node.style = json.style as any;
  if (json.props) node.props = json.props;
  if (json.enabled !== undefined) node.enabled = json.enabled;
  if (json.children) {
    node.children = json.children.map(mapSegment);
  }

  return node;
}

export function parseConfig(config: Record<string, unknown>): SegmentNode[] {
  const typed = config as JsonConfig;
  if (!typed.segments || !Array.isArray(typed.segments)) return [];
  return typed.segments.map(mapSegment);
}
