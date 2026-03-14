export interface SessionData {
  cwd: string;
  context_window?: {
    used_percentage: number;
    current_usage?: {
      input_tokens: number;
      cache_creation_input_tokens: number;
      cache_read_input_tokens: number;
    };
  };
  [key: string]: unknown;
}

export interface StyleAttrs {
  color?: string;
  bold?: boolean;
  italic?: boolean;
  prefix?: string;
  suffix?: string;
}

export type EnabledFn = (session: SessionData) => boolean;

export interface SegmentNode {
  type: string;
  provider?: string;
  enabled?: boolean | EnabledFn;
  style?: StyleAttrs;
  props?: Record<string, unknown>;
  children?: SegmentNode[];
}

export interface SegmentContext {
  session: SessionData;
  provider?: unknown;
  props?: Record<string, unknown>;
}

export interface Segment {
  name: string;
  provider?: string;
  render(context: SegmentContext): string | null;
}

export interface DataProvider {
  name: string;
  resolve(session: SessionData): Promise<unknown>;
}
