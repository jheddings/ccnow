import type { Segment, SegmentContext } from '../types.js';
import type { ContextData } from '../providers/context.js';

export const contextSizeSegment: Segment = {
  name: 'context.size',
  provider: 'context',
  render(context: SegmentContext): string | null {
    const data = context.provider as ContextData | undefined;
    return data?.size ?? null;
  },
};
