import { describe, it, expect } from '@jest/globals';
import { literalSegment } from '../../src/segments/literal.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/tmp' };

describe('literal segment', () => {
  it('returns text from props', () => {
    expect(literalSegment.render({ session, props: { text: 'hello' } })).toBe('hello');
  });

  it('returns null when no text prop', () => {
    expect(literalSegment.render({ session, props: {} })).toBeNull();
  });

  it('returns null when no props', () => {
    expect(literalSegment.render({ session })).toBeNull();
  });

  it('has correct name', () => {
    expect(literalSegment.name).toBe('literal');
  });
});
