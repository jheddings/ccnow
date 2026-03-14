import type { Segment, SegmentContext } from '../types.js';
import type { GitData } from '../providers/git.js';

export const gitBranchSegment: Segment = {
  name: 'git.branch',
  provider: 'git',
  render(context: SegmentContext): string | null {
    const data = context.provider as GitData | undefined;
    return data?.branch ?? null;
  },
};
