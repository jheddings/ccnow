import type { Segment, SegmentContext } from '../types.js';
import type { GitData } from '../providers/git.js';

export const gitInsertionsSegment: Segment = {
  name: 'git.insertions',
  provider: 'git',
  render(context: SegmentContext): string | null {
    const data = context.provider as GitData | undefined;
    if (!data?.insertions) return null;
    return `${data.insertions}`;
  },
};
