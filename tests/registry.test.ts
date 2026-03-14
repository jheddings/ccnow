import { describe, it, expect } from '@jest/globals';
import { SegmentRegistry } from '../src/registry.js';
import type { Segment, SegmentContext } from '../src/types.js';

const mockSegment: Segment = {
  name: 'test.seg',
  render: (_ctx: SegmentContext) => 'hello',
};

describe('SegmentRegistry', () => {
  it('registers and retrieves a segment by name', () => {
    const reg = new SegmentRegistry();
    reg.register(mockSegment);
    expect(reg.get('test.seg')).toBe(mockSegment);
  });

  it('returns undefined for unknown segment', () => {
    const reg = new SegmentRegistry();
    expect(reg.get('nope')).toBeUndefined();
  });

  it('registers multiple segments', () => {
    const reg = new SegmentRegistry();
    const seg2: Segment = { name: 'other', render: () => 'world' };
    reg.register(mockSegment);
    reg.register(seg2);
    expect(reg.get('test.seg')).toBe(mockSegment);
    expect(reg.get('other')).toBe(seg2);
  });

  it('later registration overwrites earlier', () => {
    const reg = new SegmentRegistry();
    const replacement: Segment = { name: 'test.seg', render: () => 'replaced' };
    reg.register(mockSegment);
    reg.register(replacement);
    expect(reg.get('test.seg')?.render({ session: { cwd: '' } })).toBe('replaced');
  });
});
