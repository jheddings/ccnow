import { describe, it, expect } from '@jest/globals';
import { defaultPreset } from '../../src/presets/default.js';

describe('default preset', () => {
  it('returns an array of segment nodes', () => {
    expect(Array.isArray(defaultPreset)).toBe(true);
    expect(defaultPreset.length).toBeGreaterThan(0);
  });

  it('starts with pwd.smart', () => {
    expect(defaultPreset[0].type).toBe('pwd.smart');
  });

  it('contains a git composite', () => {
    const git = defaultPreset.find((n) => n.type === 'git');
    expect(git).toBeDefined();
    expect(git?.children?.length).toBeGreaterThan(0);
  });

  it('contains a context composite', () => {
    const ctx = defaultPreset.find((n) => n.type === 'context');
    expect(ctx).toBeDefined();
    expect(ctx?.children?.length).toBeGreaterThan(0);
  });

  it('contains separators', () => {
    const seps = defaultPreset.filter((n) => n.type === 'sep');
    expect(seps.length).toBeGreaterThanOrEqual(2);
  });
});
