import { describe, it, expect } from '@jest/globals';
import type { SegmentNode, Segment, SegmentContext, StyleAttrs, SessionData, DataProvider } from '../src/types.js';

describe('type contracts', () => {
  it('SegmentNode accepts atomic segment config', () => {
    const node: SegmentNode = {
      type: 'git.branch',
      provider: 'git',
      style: { color: 'white', bold: true },
    };
    expect(node.type).toBe('git.branch');
    expect(node.provider).toBe('git');
    expect(node.children).toBeUndefined();
  });

  it('SegmentNode accepts composite segment config', () => {
    const node: SegmentNode = {
      type: 'git',
      enabled: true,
      children: [
        { type: 'git.branch', provider: 'git' },
        { type: 'literal', props: { text: ' [' } },
      ],
    };
    expect(node.children).toHaveLength(2);
  });

  it('SegmentNode accepts enabled as function', () => {
    const node: SegmentNode = {
      type: 'git',
      enabled: (session) => session.cwd !== '',
      children: [],
    };
    expect(typeof node.enabled).toBe('function');
  });

  it('StyleAttrs supports all style properties', () => {
    const style: StyleAttrs = {
      color: 'cyan',
      bold: true,
      italic: true,
      prefix: '+',
      suffix: '%',
    };
    expect(style.color).toBe('cyan');
  });

  it('SessionData matches Claude Code stdin shape', () => {
    const session: SessionData = {
      cwd: '/Users/test/project',
      context_window: {
        used_percentage: 42,
        current_usage: {
          input_tokens: 38000,
          cache_creation_input_tokens: 2000,
          cache_read_input_tokens: 1500,
        },
      },
    };
    expect(session.cwd).toBe('/Users/test/project');
    expect(session.context_window?.used_percentage).toBe(42);
  });

  it('SessionData allows missing context_window', () => {
    const session: SessionData = { cwd: '/tmp' };
    expect(session.context_window).toBeUndefined();
  });
});
