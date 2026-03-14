import { describe, it, expect } from '@jest/globals';
import { minimalPreset } from '../../src/presets/minimal.js';

describe('minimal preset', () => {
  it('returns an array of segment nodes', () => {
    expect(Array.isArray(minimalPreset)).toBe(true);
  });

  it('contains pwd and git.branch', () => {
    const types = minimalPreset.map((n) => n.type);
    expect(types).toContain('pwd.name');
    expect(types).toContain('git.branch');
  });

  it('is shorter than default preset', () => {
    expect(minimalPreset.length).toBeLessThanOrEqual(4);
  });
});
