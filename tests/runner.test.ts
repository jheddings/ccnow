import { describe, it, expect } from '@jest/globals';
import { run } from '../src/runner.js';

describe('run', () => {
  it('renders default preset with basic session data', async () => {
    const input = JSON.stringify({
      cwd: process.cwd(),
      context_window: {
        used_percentage: 42,
        current_usage: { input_tokens: 38000, cache_creation_input_tokens: 2000, cache_read_input_tokens: 1500 },
      },
    });
    const output = await run({ preset: 'default', segments: [], format: 'ansi' }, input);
    expect(output.length).toBeGreaterThan(0);
    // Should contain context info
    expect(output).toContain('42K');
    expect(output).toContain('42%');
  });

  it('renders plain format without ANSI codes', async () => {
    const input = JSON.stringify({ cwd: '/tmp' });
    const output = await run({ preset: 'minimal', segments: [], format: 'plain' }, input);
    // Plain format should not contain ANSI escape codes
    expect(output).not.toMatch(/\x1b\[/);
  });

  it('returns empty string for invalid stdin', async () => {
    const output = await run({ preset: 'default', segments: [], format: 'ansi' }, 'not json');
    expect(output).toBe('');
  });

  it('uses segment flags when provided', async () => {
    const input = JSON.stringify({
      cwd: '/tmp',
      context_window: { used_percentage: 10, current_usage: { input_tokens: 5000, cache_creation_input_tokens: 0, cache_read_input_tokens: 0 } },
    });
    const output = await run({ preset: 'default', segments: ['context'], format: 'plain' }, input);
    expect(output).toContain('5K');
  });
});
