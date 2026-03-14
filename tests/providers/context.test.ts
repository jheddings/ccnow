import { describe, it, expect } from '@jest/globals';
import { contextProvider } from '../../src/providers/context.js';
import type { SessionData } from '../../src/types.js';

describe('context provider', () => {
  it('formats tokens as K for thousands', async () => {
    const session: SessionData = {
      cwd: '/tmp',
      context_window: {
        used_percentage: 42,
        current_usage: { input_tokens: 38000, cache_creation_input_tokens: 2000, cache_read_input_tokens: 1500 },
      },
    };
    const data = await contextProvider.resolve(session) as any;
    expect(data.tokens).toBe('42K');
    expect(data.percent).toBe(42);
  });

  it('formats tokens as M for millions', async () => {
    const session: SessionData = {
      cwd: '/tmp',
      context_window: {
        used_percentage: 85,
        current_usage: { input_tokens: 1000000, cache_creation_input_tokens: 100000, cache_read_input_tokens: 100000 },
      },
    };
    const data = await contextProvider.resolve(session) as any;
    expect(data.tokens).toBe('1.2M');
  });

  it('returns raw number for small token counts', async () => {
    const session: SessionData = {
      cwd: '/tmp',
      context_window: {
        used_percentage: 1,
        current_usage: { input_tokens: 500, cache_creation_input_tokens: 0, cache_read_input_tokens: 0 },
      },
    };
    const data = await contextProvider.resolve(session) as any;
    expect(data.tokens).toBe('500');
  });

  it('returns null fields when context_window missing', async () => {
    const session: SessionData = { cwd: '/tmp' };
    const data = await contextProvider.resolve(session) as any;
    expect(data.tokens).toBeNull();
    expect(data.percent).toBeNull();
  });

  it('has correct name', () => {
    expect(contextProvider.name).toBe('context');
  });
});
