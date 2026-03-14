import { describe, it, expect } from '@jest/globals';
import { contextTokensSegment } from '../../src/segments/context.tokens.js';
import { contextPercentSegment } from '../../src/segments/context.percent.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/tmp' };

describe('context.tokens segment', () => {
  it('returns formatted token string', () => {
    const provider = { tokens: '42K', percent: 42 };
    expect(contextTokensSegment.render({ session, provider })).toBe('42K');
  });

  it('returns null when tokens is null', () => {
    const provider = { tokens: null, percent: null };
    expect(contextTokensSegment.render({ session, provider })).toBeNull();
  });

  it('declares context provider', () => {
    expect(contextTokensSegment.provider).toBe('context');
  });
});

describe('context.percent segment', () => {
  it('returns percent as string', () => {
    const provider = { tokens: '42K', percent: 42 };
    expect(contextPercentSegment.render({ session, provider })).toBe('42%');
  });

  it('returns null when percent is null', () => {
    const provider = { tokens: null, percent: null };
    expect(contextPercentSegment.render({ session, provider })).toBeNull();
  });
});
