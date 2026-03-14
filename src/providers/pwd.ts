import path from 'node:path';
import os from 'node:os';
import type { DataProvider, SessionData } from '../types.js';

export interface PwdData {
  name: string;
  path: string;   // full prefix path (without basename), e.g. '/Users/test/'
  smart: string;  // smart-truncated prefix (without basename), e.g. '~/t/'
}

const MAX_ABBREVIATED = 2;

function smartPrefix(cwd: string): string {
  if (cwd === '/') return '';

  const home = os.homedir();
  let p = cwd;

  // Replace home dir with ~
  if (p.startsWith(home)) {
    p = '~' + p.slice(home.length);
  }

  const parts = p.split('/');
  // 2 segments (e.g. /tmp) — just the leading slash
  if (parts.length <= 2) return parts[0] + '/';

  const first = parts[0]; // '' for absolute, '~' for home-relative
  const middle = parts.slice(1, -1);

  // 3 segments (e.g. ~/Projects/ccnow) — keep prefix as-is
  if (parts.length <= 3) return first + '/' + middle.join('/') + '/';

  // Abbreviate up to MAX_ABBREVIATED middle segments, then ellipsis if more remain
  if (middle.length <= MAX_ABBREVIATED) {
    const abbreviated = middle.map((part) => part[0] ?? '');
    return [first, ...abbreviated, ''].join('/');
  }

  const abbreviated = middle.slice(0, MAX_ABBREVIATED).map((part) => part[0] ?? '');
  return [first, ...abbreviated, '\u2026', ''].join('/');
}

export const pwdProvider: DataProvider = {
  name: 'pwd',
  async resolve(session: SessionData): Promise<PwdData> {
    const cwd = session.cwd;
    const name = cwd === '/' ? '/' : path.basename(cwd);
    const fullPrefix = cwd === '/' ? '' : path.dirname(cwd) + '/';

    return {
      name,
      path: fullPrefix,
      smart: smartPrefix(cwd),
    };
  },
};
