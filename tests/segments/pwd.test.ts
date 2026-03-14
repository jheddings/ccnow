import { describe, it, expect } from '@jest/globals';
import { pwdNameSegment } from '../../src/segments/pwd.name.js';
import { pwdPathSegment } from '../../src/segments/pwd.path.js';
import { pwdSmartSegment } from '../../src/segments/pwd.smart.js';
import type { SessionData } from '../../src/types.js';

const session: SessionData = { cwd: '/Users/test/project' };

describe('pwd.name segment', () => {
  it('returns directory name from provider data', () => {
    const provider = { name: 'project', path: '/Users/test', smart: '~/t' };
    expect(pwdNameSegment.render({ session, provider })).toBe('project');
  });

  it('returns null when no provider data', () => {
    expect(pwdNameSegment.render({ session })).toBeNull();
  });

  it('declares pwd provider', () => {
    expect(pwdNameSegment.provider).toBe('pwd');
  });
});

describe('pwd.path segment', () => {
  it('returns path prefix from provider data', () => {
    const provider = { name: 'project', path: '/Users/test', smart: '~/t' };
    expect(pwdPathSegment.render({ session, provider })).toBe('/Users/test');
  });

  it('returns null when path is empty (root dir)', () => {
    const provider = { name: '/', path: '', smart: '' };
    expect(pwdPathSegment.render({ session, provider })).toBeNull();
  });
});

describe('pwd.smart segment', () => {
  it('returns smart-truncated prefix from provider data', () => {
    const provider = { name: 'project', path: '/Users/test', smart: '~/t' };
    expect(pwdSmartSegment.render({ session, provider })).toBe('~/t');
  });

  it('returns null when smart is empty (root dir)', () => {
    const provider = { name: '/', path: '', smart: '' };
    expect(pwdSmartSegment.render({ session, provider })).toBeNull();
  });
});
