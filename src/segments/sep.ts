import type { Segment, SegmentContext } from '../types.js';

export const sepSegment: Segment = {
  name: 'sep',
  render(context: SegmentContext): string | null {
    const char = (context.props?.char as string) ?? '|';
    return ` ${char} `;
  },
};
