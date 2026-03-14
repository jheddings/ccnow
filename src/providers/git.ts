import { execFile } from 'node:child_process';
import { promisify } from 'node:util';
import type { DataProvider, SessionData } from '../types.js';

const execFileAsync = promisify(execFile);

export interface GitData {
  branch: string | null;
  insertions: number | null;
  deletions: number | null;
}

async function exec(cmd: string, args: string[], cwd: string): Promise<string> {
  try {
    const { stdout } = await execFileAsync(cmd, args, { cwd, timeout: 5000 });
    return stdout.trim();
  } catch {
    return '';
  }
}

export async function gitAvailable(cwd: string): Promise<boolean> {
  const result = await exec('git', ['-C', cwd, 'rev-parse', '--git-dir'], cwd);
  return result !== '';
}

export const gitProvider: DataProvider = {
  name: 'git',
  async resolve(session: SessionData): Promise<GitData> {
    const cwd = session.cwd;
    const isGit = await gitAvailable(cwd);

    if (!isGit) {
      return { branch: null, insertions: null, deletions: null };
    }

    const branch = await exec('git', ['branch', '--show-current'], cwd) || null;

    const diffstat = await exec('git', ['diff', '--shortstat', 'HEAD'], cwd);
    let insertions: number | null = null;
    let deletions: number | null = null;

    if (diffstat) {
      const insMatch = diffstat.match(/(\d+) insertion/);
      const delMatch = diffstat.match(/(\d+) deletion/);
      if (insMatch) insertions = parseInt(insMatch[1], 10);
      if (delMatch) deletions = parseInt(delMatch[1], 10);
    }

    return { branch, insertions, deletions };
  },
};
