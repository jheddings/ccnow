import { describe, it, expect } from '@jest/globals';
import {
  StatusLine, Pwd, Sep, Git, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from '../../src/dsl/index.js';
import type { SegmentNode } from '../../src/types.js';

describe('DSL factory functions', () => {
  it('Literal creates a literal node', () => {
    const node = Literal({ text: 'hello' });
    expect(node.type).toBe('literal');
    expect(node.props?.text).toBe('hello');
  });

  it('Sep creates a separator node', () => {
    const node = Sep({ char: '|', color: '240' });
    expect(node.type).toBe('sep');
    expect(node.props?.char).toBe('|');
    expect(node.style?.color).toBe('240');
  });

  it('Pwd creates a pwd.smart node by default', () => {
    const node = Pwd({ color: 'cyan' });
    expect(node.type).toBe('pwd.smart');
    expect(node.provider).toBe('pwd');
    expect(node.style?.color).toBe('cyan');
  });

  it('Pwd respects style override', () => {
    const node = Pwd({ style: 'name', color: 'blue' });
    expect(node.type).toBe('pwd.name');
  });

  it('Branch creates a git.branch node', () => {
    const node = Branch({ color: 'white', prefix: '\ue0a0 ' });
    expect(node.type).toBe('git.branch');
    expect(node.provider).toBe('git');
    expect(node.style?.prefix).toBe('\ue0a0 ');
  });

  it('Git creates a composite node with trailing closure', () => {
    const node = Git()(() => [
      Branch({ color: 'white' }),
      Literal({ text: ' ' }),
    ]);
    expect(node.type).toBe('git');
    expect(node.children).toHaveLength(2);
    expect(node.children![0].type).toBe('git.branch');
  });

  it('Git accepts enabled function', () => {
    const enabledFn = (s: any) => s.cwd !== '';
    const node = Git({ enabled: enabledFn })(() => [Branch()]);
    expect(typeof node.enabled).toBe('function');
  });

  it('Context creates a composite node', () => {
    const node = Context()(() => [
      Tokens({ bold: true }),
      Percent(),
    ]);
    expect(node.type).toBe('context');
    expect(node.children).toHaveLength(2);
  });

  it('StatusLine returns flat array of nodes', () => {
    const tree = StatusLine(() => [
      Pwd({ color: 'cyan' }),
      Sep({ char: '|' }),
    ]);
    expect(tree).toHaveLength(2);
    expect(tree[0].type).toBe('pwd.smart');
    expect(tree[1].type).toBe('sep');
  });
});
