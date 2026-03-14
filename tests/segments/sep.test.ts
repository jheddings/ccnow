import { describe, it, expect } from '@jest/globals';
import { sepSegment } from '../../src/segments/sep.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/tmp' };

describe('sep segment', () => {
  it('returns char from props with spaces', () => {
    expect(sepSegment.render({ session, props: { char: '|' } })).toBe(' | ');
  });

  it('defaults to pipe when no char prop', () => {
    expect(sepSegment.render({ session })).toBe(' | ');
  });

  it('has correct name', () => {
    expect(sepSegment.name).toBe('sep');
  });
});
