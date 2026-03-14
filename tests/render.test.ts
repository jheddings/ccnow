import { describe, it, expect } from '@jest/globals';
import { renderTree } from '../src/render.js';
import { SegmentRegistry } from '../src/registry.js';
import type { SessionData, SegmentNode, Segment } from '../src/types.js';
import type { GitData } from '../src/providers/git.js';

const session: SessionData = { cwd: '/Users/test/project' };
const providerData = new Map<string, unknown>();

function makeRegistry(...segments: Segment[]): SegmentRegistry {
  const reg = new SegmentRegistry();
  for (const seg of segments) reg.register(seg);
  return reg;
}

describe('renderTree', () => {
  it('renders a single atomic segment', () => {
    const reg = makeRegistry({
      name: 'literal',
      render: (ctx) => (ctx.props?.text as string) ?? null,
    });
    const tree: SegmentNode[] = [{ type: 'literal', props: { text: 'hello' } }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('hello');
  });

  it('applies style to segment output', () => {
    const reg = makeRegistry({
      name: 'literal',
      render: (ctx) => (ctx.props?.text as string) ?? null,
    });
    const tree: SegmentNode[] = [
      { type: 'literal', props: { text: 'hello' }, style: { bold: true } },
    ];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toContain('hello');
    expect(result.length).toBeGreaterThan('hello'.length); // ANSI codes added
  });

  it('skips segments that render null', () => {
    const reg = makeRegistry({ name: 'empty', render: () => null });
    const tree: SegmentNode[] = [{ type: 'empty' }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('skips segments with enabled=false', () => {
    const reg = makeRegistry({
      name: 'literal',
      render: (ctx) => (ctx.props?.text as string) ?? null,
    });
    const tree: SegmentNode[] = [{ type: 'literal', props: { text: 'hidden' }, enabled: false }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('evaluates enabled function', () => {
    const reg = makeRegistry({
      name: 'literal',
      render: (ctx) => (ctx.props?.text as string) ?? null,
    });
    const tree: SegmentNode[] = [
      {
        type: 'literal',
        props: { text: 'shown' },
        enabled: (s) => s.cwd === '/Users/test/project',
      },
      { type: 'literal', props: { text: 'hidden' }, enabled: (s) => s.cwd === '/nope' },
    ];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('shown');
  });

  it('renders composite children and concatenates', () => {
    const reg = makeRegistry({
      name: 'literal',
      render: (ctx) => (ctx.props?.text as string) ?? null,
    });
    const tree: SegmentNode[] = [
      {
        type: 'group',
        children: [
          { type: 'literal', props: { text: 'a' } },
          { type: 'literal', props: { text: 'b' } },
        ],
      },
    ];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('ab');
  });

  it('collapses composite to empty when all children are null', () => {
    const reg = makeRegistry({ name: 'empty', render: () => null });
    const tree: SegmentNode[] = [
      {
        type: 'group',
        children: [{ type: 'empty' }, { type: 'empty' }],
      },
    ];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('collapses composite with disabled enabled flag', () => {
    const reg = makeRegistry({
      name: 'literal',
      render: (ctx) => (ctx.props?.text as string) ?? null,
    });
    const tree: SegmentNode[] = [
      {
        type: 'group',
        enabled: false,
        children: [{ type: 'literal', props: { text: 'should not appear' } }],
      },
    ];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('passes provider data to segment context', () => {
    const providers = new Map<string, unknown>([['git', { branch: 'main' }]]);
    const reg = makeRegistry({
      name: 'git.branch',
      provider: 'git',
      render: (ctx) => (ctx.provider as GitData)?.branch ?? null,
    });
    const tree: SegmentNode[] = [{ type: 'git.branch', provider: 'git' }];
    const result = renderTree(tree, reg, session, providers);
    expect(result).toBe('main');
  });

  it('skips segment when unknown type', () => {
    const reg = makeRegistry();
    const tree: SegmentNode[] = [{ type: 'unknown.thing' }];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('collapses composite with prefix/suffix when all children are null', () => {
    const reg = makeRegistry({ name: 'empty', render: () => null });
    const tree: SegmentNode[] = [
      {
        type: 'group',
        style: { prefix: ' [', suffix: ']' },
        children: [{ type: 'empty' }, { type: 'empty' }],
      },
    ];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe('');
  });

  it('renders composite prefix/suffix when children have data', () => {
    const reg = makeRegistry({ name: 'val', render: () => '5' });
    const tree: SegmentNode[] = [
      {
        type: 'group',
        style: { prefix: ' [', suffix: ']' },
        children: [{ type: 'val' }],
      },
    ];
    const result = renderTree(tree, reg, session, providerData);
    expect(result).toBe(' [5]');
  });

  it('renders git composite without brackets when diff is clean', () => {
    const gitData = { branch: 'main', insertions: 0, deletions: 0 };
    const providers = new Map<string, unknown>([['git', gitData]]);
    const reg = makeRegistry(
      {
        name: 'git.branch',
        provider: 'git',
        render: (ctx) => (ctx.provider as GitData)?.branch ?? null,
      },
      {
        name: 'git.insertions',
        provider: 'git',
        render: (ctx) => {
          const d = ctx.provider as GitData;
          return d?.insertions ? `${d.insertions}` : null;
        },
      },
      {
        name: 'git.deletions',
        provider: 'git',
        render: (ctx) => {
          const d = ctx.provider as GitData;
          return d?.deletions ? `${d.deletions}` : null;
        },
      }
    );
    const tree: SegmentNode[] = [
      {
        type: 'git',
        children: [
          { type: 'git.branch', provider: 'git', style: { prefix: '\ue0a0 ' } },
          {
            type: 'group',
            style: { prefix: ' [', suffix: ']' },
            children: [
              { type: 'git.insertions', provider: 'git', style: { prefix: '+' } },
              { type: 'git.deletions', provider: 'git', style: { prefix: ' -' } },
            ],
          },
        ],
      },
    ];
    const result = renderTree(tree, reg, session, providers);
    expect(result).toBe('\ue0a0 main');
    expect(result).not.toContain('[');
    expect(result).not.toContain(']');
  });

  it('renders git composite with brackets when diff has changes', () => {
    const gitData = { branch: 'feat/test', insertions: 10, deletions: 3 };
    const providers = new Map<string, unknown>([['git', gitData]]);
    const reg = makeRegistry(
      {
        name: 'git.branch',
        provider: 'git',
        render: (ctx) => (ctx.provider as GitData)?.branch ?? null,
      },
      {
        name: 'git.insertions',
        provider: 'git',
        render: (ctx) => {
          const d = ctx.provider as GitData;
          return d?.insertions ? `${d.insertions}` : null;
        },
      },
      {
        name: 'git.deletions',
        provider: 'git',
        render: (ctx) => {
          const d = ctx.provider as GitData;
          return d?.deletions ? `${d.deletions}` : null;
        },
      }
    );
    const tree: SegmentNode[] = [
      {
        type: 'git',
        children: [
          { type: 'git.branch', provider: 'git' },
          {
            type: 'group',
            style: { prefix: ' [', suffix: ']' },
            children: [
              { type: 'git.insertions', provider: 'git', style: { prefix: '+' } },
              { type: 'git.deletions', provider: 'git', style: { prefix: ' -' } },
            ],
          },
        ],
      },
    ];
    const result = renderTree(tree, reg, session, providers);
    expect(result).toBe('feat/test [+10 -3]');
  });

  it('renders git composite with only insertions', () => {
    const gitData = { branch: 'main', insertions: 7, deletions: 0 };
    const providers = new Map<string, unknown>([['git', gitData]]);
    const reg = makeRegistry(
      {
        name: 'git.branch',
        provider: 'git',
        render: (ctx) => (ctx.provider as GitData)?.branch ?? null,
      },
      {
        name: 'git.insertions',
        provider: 'git',
        render: (ctx) => {
          const d = ctx.provider as GitData;
          return d?.insertions ? `${d.insertions}` : null;
        },
      },
      {
        name: 'git.deletions',
        provider: 'git',
        render: (ctx) => {
          const d = ctx.provider as GitData;
          return d?.deletions ? `${d.deletions}` : null;
        },
      }
    );
    const tree: SegmentNode[] = [
      {
        type: 'git',
        children: [
          { type: 'git.branch', provider: 'git' },
          {
            type: 'group',
            style: { prefix: ' [', suffix: ']' },
            children: [
              { type: 'git.insertions', provider: 'git', style: { prefix: '+' } },
              { type: 'git.deletions', provider: 'git', style: { prefix: ' -' } },
            ],
          },
        ],
      },
    ];
    const result = renderTree(tree, reg, session, providers);
    expect(result).toBe('main [+7]');
  });
});
