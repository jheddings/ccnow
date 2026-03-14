import type { Segment, SegmentContext } from '../types.js';

export const literalSegment: Segment = {
  name: 'literal',
  render(context: SegmentContext): string | null {
    const text = context.props?.text;
    if (typeof text !== 'string') return null;
    return text;
  },
};
