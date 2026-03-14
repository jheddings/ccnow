import type { Segment, SegmentContext } from '../types.js';
import type { PwdData } from '../providers/pwd.js';

export const pwdSmartSegment: Segment = {
  name: 'pwd.smart',
  provider: 'pwd',
  render(context: SegmentContext): string | null {
    const data = context.provider as PwdData | undefined;
    return data?.smart || null;
  },
};
