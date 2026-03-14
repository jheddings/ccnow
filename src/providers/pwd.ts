import path from 'node:path';
import os from 'node:os';
import type { DataProvider, SessionData } from '../types.js';

export interface PwdData {
  name: string;
  path: string;
  smart: string;
}

function smartTruncate(cwd: string): string {
  const home = os.homedir();
  let p = cwd;

  // Replace home dir with ~
  if (p.startsWith(home)) {
    p = '~' + p.slice(home.length);
  }

  const parts = p.split('/');
  if (parts.length <= 3) return p;

  // Keep first component and last component, truncate middle to initials
  const first = parts[0]; // '' for absolute, '~' for home-relative
  const last = parts[parts.length - 1];
  const middle = parts.slice(1, -1).map((part) => part[0] ?? '');

  return [first, ...middle, last].join('/');
}

export const pwdProvider: DataProvider = {
  name: 'pwd',
  async resolve(session: SessionData): Promise<PwdData> {
    const cwd = session.cwd;
    return {
      name: cwd === '/' ? '/' : path.basename(cwd),
      path: cwd,
      smart: smartTruncate(cwd),
    };
  },
};
