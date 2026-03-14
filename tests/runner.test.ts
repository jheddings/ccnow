import { describe, it, expect } from '@jest/globals';
import { readFileSync } from 'node:fs';
import { run } from '../src/runner.js';

const sessionFull = readFileSync(new URL('./data/session-v2.1.75.json', import.meta.url), 'utf-8');

describe('run', () => {
  it('renders default preset with session data', async () => {
    const output = await run({ preset: 'default', format: 'ansi' }, sessionFull);
    expect(output.length).toBeGreaterThan(0);
    expect(output).toContain('123K');
    expect(output).toContain('12%');
  });

  it('renders minimal preset', async () => {
    const output = await run({ preset: 'minimal', format: 'plain' }, sessionFull);
    expect(output).toContain('my-app');
    expect(output).toContain('123K/1M');
  });

  it('renders full preset', async () => {
    const output = await run({ preset: 'full', format: 'plain' }, sessionFull);
    expect(output).toContain('Opus 4.6');
    expect(output).toContain('$7.93');
    expect(output).toContain('6h 53m');
    expect(output).toContain('+610');
    expect(output).toContain('-59');
  });

  it('renders plain format without ANSI codes', async () => {
    const output = await run({ preset: 'minimal', format: 'plain' }, sessionFull);
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
    expect(output).not.toContain('| |');
  });
});
