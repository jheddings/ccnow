import { describe, it, expect } from '@jest/globals';
import { parseSession } from '../src/session.js';

describe('parseSession', () => {
  it('parses valid JSON with all fields', () => {
    const input = JSON.stringify({
      cwd: '/Users/test/project',
      context_window: {
        used_percentage: 42,
        current_usage: { input_tokens: 38000, cache_creation_input_tokens: 2000, cache_read_input_tokens: 1500 },
      },
    });
    const session = parseSession(input)!;
    expect(session.cwd).toBe('/Users/test/project');
    expect(session.context_window?.used_percentage).toBe(42);
  });

  it('parses minimal JSON with only cwd', () => {
    const session = parseSession('{"cwd":"/tmp"}')!;
    expect(session.cwd).toBe('/tmp');
    expect(session.context_window).toBeUndefined();
  });

  it('returns null for empty string', () => {
    expect(parseSession('')).toBeNull();
  });

  it('returns null for invalid JSON', () => {
    expect(parseSession('not json')).toBeNull();
  });

  it('returns null for JSON missing cwd', () => {
    expect(parseSession('{"foo":"bar"}')).toBeNull();
  });

  it('preserves extra fields via index signature', () => {
    const session = parseSession('{"cwd":"/tmp","model":"opus"}');
    expect(session?.['model']).toBe('opus');
  });
});
