import { StatusLine, Git, Branch, Context, Tokens, Literal, Size } from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const minimalPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.name', provider: 'pwd', style: { color: '39' } },
  Git({ prefix: ' | ', color: '240' })(() => [
    Branch({ color: 'whiteBright', bold: true }),
  ]),
  Context({ prefix: ' | ', color: '240' })(() => [
    Tokens({ color: 'white' }),
    Literal({ text: '/' }),
    Size({ color: 'white' }),
  ]),
]);
