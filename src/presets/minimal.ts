import { StatusLine, Pwd, Sep, Branch } from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const minimalPreset: SegmentNode[] = StatusLine(() => [
  Pwd({ style: 'name', color: 'cyan' }),
  Sep({ char: '|', dim: true }),
  Branch({ color: 'whiteBright', bold: true, icon: '\ue0a0 ' }),
]);
