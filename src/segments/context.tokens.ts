import type { Segment, SegmentContext } from '../types.js';
import type { ContextData } from '../providers/context.js';

export const contextTokensSegment: Segment = {
  name: 'context.tokens',
  provider: 'context',
  render(context: SegmentContext): string | null {
    const data = context.provider as ContextData | undefined;
    return data?.tokens ?? null;
  },
};
