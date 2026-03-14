import { describe, it, expect } from '@jest/globals';
import { pwdProvider } from '../../src/providers/pwd.js';
import type { SessionData } from '../../src/types.js';

describe('pwd provider', () => {
  it('resolves name, path, and smart from cwd', async () => {
    const session: SessionData = { cwd: '/Users/jheddings/Projects/ccnow' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.name).toBe('ccnow');
    expect(data.path).toBe('/Users/jheddings/Projects/ccnow');
    expect(data.smart).toBeDefined();
  });

  it('handles root path', async () => {
    const session: SessionData = { cwd: '/' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.name).toBe('/');
    expect(data.path).toBe('/');
  });

  it('smart truncates long paths', async () => {
    const session: SessionData = { cwd: '/Users/jheddings/Projects/very/deep/nested/path' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart.length).toBeLessThan(data.path.length);
    expect(data.smart).toContain('path'); // always keeps last component
  });

  it('smart keeps short paths as-is', async () => {
    const session: SessionData = { cwd: '/tmp' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart).toBe('/tmp');
  });

  it('replaces home dir with ~ in smart', async () => {
    const home = process.env.HOME ?? '/Users/test';
    const session: SessionData = { cwd: `${home}/Projects/ccnow` };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart.startsWith('~')).toBe(true);
  });

  it('has correct name', () => {
    expect(pwdProvider.name).toBe('pwd');
  });
});
