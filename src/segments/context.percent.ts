import type { Segment, SegmentContext } from '../types.js';
import type { ContextData } from '../providers/context.js';

export const contextPercentSegment: Segment = {
  name: 'context.percent',
  provider: 'context',
  render(context: SegmentContext): string | null {
    const data = context.provider as ContextData | undefined;
    if (data?.percent == null) return null;
    return `${data.percent}%`;
  },
};
