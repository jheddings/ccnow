import { describe, it, expect, jest } from '@jest/globals';
import { ProviderRegistry } from '../src/providers.js';
import type { DataProvider, SessionData, SegmentNode } from '../src/types.js';

const session: SessionData = { cwd: '/tmp' };

const mockProvider: DataProvider = {
  name: 'test',
  resolve: async (s) => ({ value: s.cwd }),
};

describe('ProviderRegistry', () => {
  it('registers and resolves a provider', async () => {
    const reg = new ProviderRegistry();
    reg.register(mockProvider);
    const results = await reg.resolveAll(['test'], session);
    expect(results.get('test')).toEqual({ value: '/tmp' });
  });

  it('resolves multiple providers concurrently', async () => {
    const reg = new ProviderRegistry();
    const provider2: DataProvider = {
      name: 'other',
      resolve: async () => ({ data: 42 }),
    };
    reg.register(mockProvider);
    reg.register(provider2);
    const results = await reg.resolveAll(['test', 'other'], session);
    expect(results.get('test')).toEqual({ value: '/tmp' });
    expect(results.get('other')).toEqual({ data: 42 });
  });

  it('returns empty map when no providers requested', async () => {
    const reg = new ProviderRegistry();
    const results = await reg.resolveAll([], session);
    expect(results.size).toBe(0);
  });

  it('skips unknown provider names gracefully', async () => {
    const reg = new ProviderRegistry();
    const results = await reg.resolveAll(['nonexistent'], session);
    expect(results.has('nonexistent')).toBe(false);
  });

  it('catches provider errors and excludes from results', async () => {
    const reg = new ProviderRegistry();
    const failing: DataProvider = {
      name: 'fail',
      resolve: async () => { throw new Error('boom'); },
    };
    reg.register(failing);
    const results = await reg.resolveAll(['fail'], session);
    expect(results.has('fail')).toBe(false);
  });

  it('collects required provider names from a segment tree', () => {
    const reg = new ProviderRegistry();
    const tree: SegmentNode[] = [
      { type: 'pwd.smart', provider: 'pwd' },
      { type: 'sep' },
      { type: 'git', children: [
        { type: 'git.branch', provider: 'git' },
        { type: 'literal', props: { text: ' ' } },
        { type: 'git.insertions', provider: 'git' },
      ]},
    ];
    const names = reg.collectProviderNames(tree);
    expect([...names].sort()).toEqual(['git', 'pwd']);
  });

  it('collectProviderNames skips nodes with enabled=false', () => {
    const reg = new ProviderRegistry();
    const tree: SegmentNode[] = [
      { type: 'pwd.smart', provider: 'pwd' },
      { type: 'git', enabled: false, children: [
        { type: 'git.branch', provider: 'git' },
      ]},
    ];
    const names = reg.collectProviderNames(tree);
    expect([...names]).toEqual(['pwd']);
  });
});
