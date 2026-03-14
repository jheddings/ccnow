import type { Segment, SegmentContext } from '../types.js';
import type { PwdData } from '../providers/pwd.js';

export const pwdPathSegment: Segment = {
  name: 'pwd.path',
  provider: 'pwd',
  render(context: SegmentContext): string | null {
    const data = context.provider as PwdData | undefined;
    return data?.path ?? null;
  },
};
