import { StatusLine, Pwd, Sep, Branch } from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const minimalPreset: SegmentNode[] = StatusLine(() => [
  Pwd({ style: 'name', color: 'cyan', bold: true }),
  Sep({ char: '|', dim: true }),
  Branch({ color: 'white', bold: true, icon: '\ue0a0 ' }),
]);
