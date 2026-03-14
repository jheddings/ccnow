import { StatusLine, Git, Branch } from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const minimalPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.name', provider: 'pwd', style: { color: 'cyanBright', bold: true } },
  Git({ prefix: ' | ', color: '240' })(() => [
    Branch({ color: 'whiteBright', bold: true }),
  ]),
]);
