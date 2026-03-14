import { describe, it, expect } from '@jest/globals';
import {
  StatusLine,
  Pwd,
  Literal,
  Group,
  GitGroup,
  GitBranch,
  ContextGroup,
  ContextTokens,
  ContextPercent,
} from '../../src/dsl/index.js';
import type { SessionData } from '../../src/types.js';

describe('DSL factory functions', () => {
  it('Literal creates a literal node', () => {
    const node = Literal({ text: 'hello' });
    expect(node.type).toBe('literal');
    expect(node.props?.text).toBe('hello');
  });

  it('Group creates a composite with trailing closure', () => {
    const node = Group({ prefix: ' [', suffix: ']' })(() => [Literal({ text: 'a' })]);
    expect(node.type).toBe('group');
    expect(node.children).toHaveLength(1);
    expect(node.style?.prefix).toBe(' [');
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

  it('GitBranch creates a git.branch node', () => {
    const node = GitBranch({ color: 'white', prefix: '\ue0a0 ' });
    expect(node.type).toBe('git.branch');
    expect(node.provider).toBe('git');
    expect(node.style?.prefix).toBe('\ue0a0 ');
  });

  it('GitGroup creates a composite with trailing closure', () => {
    const node = GitGroup()(() => [GitBranch({ color: 'white' }), Literal({ text: ' ' })]);
    expect(node.type).toBe('git');
    expect(node.children).toHaveLength(2);
    expect(node.children![0].type).toBe('git.branch');
  });

  it('GitGroup accepts enabled function', () => {
    const enabledFn = (s: SessionData) => s.cwd !== '';
    const node = GitGroup({ enabled: enabledFn })(() => [GitBranch()]);
    expect(typeof node.enabled).toBe('function');
  });

  it('ContextGroup creates a composite node', () => {
    const node = ContextGroup()(() => [ContextTokens({ bold: true }), ContextPercent()]);
    expect(node.type).toBe('context');
    expect(node.children).toHaveLength(2);
  });

  it('StatusLine returns flat array of nodes', () => {
    const tree = StatusLine(() => [Pwd({ color: 'cyan' }), Literal({ text: ' | ' })]);
    expect(tree).toHaveLength(2);
    expect(tree[0].type).toBe('pwd.smart');
    expect(tree[1].type).toBe('literal');
  });
});
