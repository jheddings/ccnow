import type { DataProvider, SessionData } from '../types.js';

export interface ContextData {
  tokens: string | null;
  percent: number | null;
}

function formatTokens(total: number): string {
  if (total >= 1_000_000) {
    return `${(total / 1_000_000).toFixed(1)}M`;
  }
  if (total >= 1_000) {
    return `${Math.round(total / 1_000)}K`;
  }
  return `${total}`;
}

export const contextProvider: DataProvider = {
  name: 'context',
  async resolve(session: SessionData): Promise<ContextData> {
    const cw = session.context_window;
    if (!cw) {
      return { tokens: null, percent: null };
    }

    const usage = cw.current_usage;
    let totalTokens = 0;
    if (usage) {
      totalTokens =
        (usage.input_tokens ?? 0) +
        (usage.cache_creation_input_tokens ?? 0) +
        (usage.cache_read_input_tokens ?? 0);
    }

    return {
      tokens: usage ? formatTokens(totalTokens) : null,
      percent: cw.used_percentage ?? null,
    };
  },
};
