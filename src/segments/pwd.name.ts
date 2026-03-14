import type { Segment, SegmentContext } from '../types.js';
import type { PwdData } from '../providers/pwd.js';

export const pwdNameSegment: Segment = {
  name: 'pwd.name',
  provider: 'pwd',
  render(context: SegmentContext): string | null {
    const data = context.provider as PwdData | undefined;
    return data?.name ?? null;
  },
};
