import type { DataProvider, SessionData, SegmentNode } from './types.js';

export class ProviderRegistry {
  private providers = new Map<string, DataProvider>();

  register(provider: DataProvider): void {
    this.providers.set(provider.name, provider);
  }

  collectProviderNames(tree: SegmentNode[]): Set<string> {
    const names = new Set<string>();
    const walk = (nodes: SegmentNode[]) => {
      for (const node of nodes) {
        // Skip statically disabled nodes (and their children)
        if (node.enabled === false) continue;
        if (node.provider) names.add(node.provider);
        if (node.children) walk(node.children);
      }
    };
    walk(tree);
    return names;
  }

  async resolveAll(
    names: string[],
    session: SessionData,
  ): Promise<Map<string, unknown>> {
    const results = new Map<string, unknown>();

    const entries = names
      .map((name) => ({ name, provider: this.providers.get(name) }))
      .filter((e): e is { name: string; provider: DataProvider } => e.provider !== undefined);

    const settled = await Promise.allSettled(
      entries.map(async ({ name, provider }) => {
        const data = await provider.resolve(session);
        return { name, data };
      }),
    );

    for (const result of settled) {
      if (result.status === 'fulfilled') {
        results.set(result.value.name, result.value.data);
      }
    }

    return results;
  }
}
