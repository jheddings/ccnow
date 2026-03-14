import { describe, it, expect } from '@jest/globals';
import { gitProvider, gitAvailable } from '../../src/providers/git.js';
import type { GitData } from '../../src/providers/git.js';
import type { SessionData } from '../../src/types.js';

describe('gitAvailable', () => {
  it('returns true for a git repo directory', async () => {
    // Use the ccnow project dir itself (which is a git repo)
    const result = await gitAvailable(process.cwd());
    expect(result).toBe(true);
  });

  it('returns false for /tmp', async () => {
    const result = await gitAvailable('/tmp');
    expect(result).toBe(false);
  });
});

describe('git provider', () => {
  it('resolves data for a git repo', async () => {
    const session: SessionData = { cwd: process.cwd() };
    const data = (await gitProvider.resolve(session)) as GitData;
    // branch may be null in detached HEAD (e.g. CI checkout)
    expect(data.branch === null || typeof data.branch === 'string').toBe(true);
  });

  it('returns null fields for non-git directory', async () => {
    const session: SessionData = { cwd: '/tmp' };
    const data = (await gitProvider.resolve(session)) as GitData;
    expect(data.branch).toBeNull();
    expect(data.insertions).toBeNull();
    expect(data.deletions).toBeNull();
  });

  it('has correct name', () => {
    expect(gitProvider.name).toBe('git');
  });
});
