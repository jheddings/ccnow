import { describe, it, expect } from '@jest/globals';
import { run } from '../src/runner.js';

describe('run', () => {
  it('renders default preset with basic session data', async () => {
    const input = JSON.stringify({
      cwd: process.cwd(),
      context_window: {
        used_percentage: 42,
        current_usage: {
          input_tokens: 38000,
          cache_creation_input_tokens: 2000,
          cache_read_input_tokens: 1500,
        },
      },
    });
    const output = await run({ preset: 'default', format: 'ansi' }, input);
    expect(output.length).toBeGreaterThan(0);
    expect(output).toContain('42K');
    expect(output).toContain('42%');
  });

  it('renders plain format without ANSI codes', async () => {
    const input = JSON.stringify({ cwd: '/tmp' });
    const output = await run({ preset: 'minimal', format: 'plain' }, input);
    // eslint-disable-next-line no-control-regex
    expect(output).not.toMatch(/\x1b\[/);
  });

  it('returns empty string for invalid stdin', async () => {
    const output = await run({ preset: 'default', format: 'ansi' }, 'not json');
    expect(output).toBe('');
  });

  it('renders non-git directory without empty separators', async () => {
    const input = JSON.stringify({
      cwd: '/tmp',
      context_window: {
        used_percentage: 10,
        current_usage: {
          input_tokens: 5000,
          cache_creation_input_tokens: 0,
          cache_read_input_tokens: 0,
        },
      },
    });
    const output = await run({ preset: 'default', format: 'plain' }, input);
    expect(output).toContain('/tmp');
    expect(output).toContain('5K');
    // Should not have double separators from collapsed git
    expect(output).not.toContain('| |');
  });
});
