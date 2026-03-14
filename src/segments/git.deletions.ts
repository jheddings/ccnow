import type { Segment, SegmentContext } from '../types.js';
import type { GitData } from '../providers/git.js';

export const gitDeletionsSegment: Segment = {
  name: 'git.deletions',
  provider: 'git',
  render(context: SegmentContext): string | null {
    const data = context.provider as GitData | undefined;
    if (!data?.deletions) return null;
    return `${data.deletions}`;
  },
};
