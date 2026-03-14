import { describe, it, expect } from '@jest/globals';
import { gitBranchSegment } from '../../src/segments/git.branch.js';
import { gitInsertionsSegment } from '../../src/segments/git.insertions.js';
import { gitDeletionsSegment } from '../../src/segments/git.deletions.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/tmp' };

describe('git.branch segment', () => {
  it('returns branch name', () => {
    const provider = { branch: 'main', insertions: 5, deletions: 3 };
    expect(gitBranchSegment.render({ session, provider })).toBe('main');
  });

  it('returns null when branch is null', () => {
    const provider = { branch: null, insertions: null, deletions: null };
    expect(gitBranchSegment.render({ session, provider })).toBeNull();
  });

  it('declares git provider', () => {
    expect(gitBranchSegment.provider).toBe('git');
  });
});

describe('git.insertions segment', () => {
  it('returns insertion count as string', () => {
    const provider = { branch: 'main', insertions: 12, deletions: 3 };
    expect(gitInsertionsSegment.render({ session, provider })).toBe('12');
  });

  it('returns null when insertions is 0', () => {
    const provider = { branch: 'main', insertions: 0, deletions: 0 };
    expect(gitInsertionsSegment.render({ session, provider })).toBeNull();
  });

  it('returns null when insertions is null', () => {
    const provider = { branch: null, insertions: null, deletions: null };
    expect(gitInsertionsSegment.render({ session, provider })).toBeNull();
  });
});

describe('git.deletions segment', () => {
  it('returns deletion count as string', () => {
    const provider = { branch: 'main', insertions: 0, deletions: 7 };
    expect(gitDeletionsSegment.render({ session, provider })).toBe('7');
  });

  it('returns null when deletions is 0', () => {
    const provider = { branch: 'main', insertions: 0, deletions: 0 };
    expect(gitDeletionsSegment.render({ session, provider })).toBeNull();
  });
});
