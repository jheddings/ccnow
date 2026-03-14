import type { SessionData } from './types.js';

export function parseSession(input: string): SessionData | null {
  if (!input.trim()) return null;

  let parsed: unknown;
  try {
    parsed = JSON.parse(input);
  } catch {
    return null;
  }

  if (typeof parsed !== 'object' || parsed === null) return null;

  const obj = parsed as Record<string, unknown>;
  if (typeof obj.cwd !== 'string') return null;

  return obj as SessionData;
}
