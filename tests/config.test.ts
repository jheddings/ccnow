import { describe, it, expect } from '@jest/globals';
import { parseConfig } from '../src/config.js';
import type { SegmentNode } from '../src/types.js';

describe('parseConfig', () => {
  it('parses a simple segment list', () => {
    const json = {
      segments: [
        { segment: 'pwd.smart', style: { color: 'cyan' } },
        { segment: 'sep', props: { char: '|' } },
      ],
    };
    const tree = parseConfig(json);
    expect(tree).toHaveLength(2);
    expect(tree[0].type).toBe('pwd.smart');
    expect(tree[0].style?.color).toBe('cyan');
    expect(tree[1].type).toBe('sep');
    expect(tree[1].props?.char).toBe('|');
  });

  it('parses composite segments with children', () => {
    const json = {
      segments: [
        {
          segment: 'git',
          children: [
            { segment: 'git.branch', style: { color: 'white' } },
          ],
        },
      ],
    };
    const tree = parseConfig(json);
    expect(tree[0].type).toBe('git');
    expect(tree[0].children).toHaveLength(1);
    expect(tree[0].children![0].type).toBe('git.branch');
  });

  it('handles enabled boolean', () => {
    const json = {
      segments: [
        { segment: 'git', enabled: false, children: [] },
      ],
    };
    const tree = parseConfig(json);
    expect(tree[0].enabled).toBe(false);
  });

  it('returns empty array for missing segments', () => {
    expect(parseConfig({})).toEqual([]);
    expect(parseConfig({ segments: [] })).toEqual([]);
  });

  it('maps segment field to type and infers provider', () => {
    const json = {
      segments: [{ segment: 'git.branch' }],
    };
    const tree = parseConfig(json);
    expect(tree[0].type).toBe('git.branch');
    expect(tree[0].provider).toBe('git');
  });

  it('does not infer provider for literal and sep', () => {
    const json = {
      segments: [
        { segment: 'literal', props: { text: 'hi' } },
        { segment: 'sep' },
      ],
    };
    const tree = parseConfig(json);
    expect(tree[0].provider).toBeUndefined();
    expect(tree[1].provider).toBeUndefined();
  });
});
