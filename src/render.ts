import type { SegmentNode, SessionData } from './types.js';
import { SegmentRegistry } from './registry.js';
import { applyStyle } from './style.js';

function isEnabled(node: SegmentNode, session: SessionData): boolean {
  if (node.enabled === undefined) return true;
  if (typeof node.enabled === 'boolean') return node.enabled;
  try {
    return node.enabled(session);
  } catch {
    return false;
  }
}

function renderNode(
  node: SegmentNode,
  registry: SegmentRegistry,
  session: SessionData,
  providerData: Map<string, unknown>,
): string | null {
  if (!isEnabled(node, session)) return null;

  // Composite node: render children
  if (node.children) {
    const parts: string[] = [];
    for (const child of node.children) {
      const rendered = renderNode(child, registry, session, providerData);
      if (rendered !== null) parts.push(rendered);
    }
    if (parts.length === 0) return null;
    const joined = parts.join('');
    return applyStyle(joined, node.style);
  }

  // Atomic node: look up segment and render
  const segment = registry.get(node.type);
  if (!segment) return null;

  const context = {
    session,
    provider: node.provider ? providerData.get(node.provider) : undefined,
    props: node.props,
  };

  const value = segment.render(context);
  if (value === null) return null;

  return applyStyle(value, node.style);
}

export function renderTree(
  tree: SegmentNode[],
  registry: SegmentRegistry,
  session: SessionData,
  providerData: Map<string, unknown>,
): string {
  const parts: string[] = [];
  for (const node of tree) {
    const rendered = renderNode(node, registry, session, providerData);
    if (rendered !== null) parts.push(rendered);
  }
  return parts.join('');
}
