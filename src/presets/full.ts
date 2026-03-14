import {
  StatusLine, Sep, Git, Group, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const fullPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.path', provider: 'pwd', style: { color: 'cyanBright' } },
  { type: 'pwd.name', provider: 'pwd', style: { color: 'cyanBright', bold: true } },
  Git()(() => [
    Sep({ char: '|', color: '240' }),
    Branch({ color: 'whiteBright', bold: true, prefix: '\ue0a0 ' }),
    Group({ prefix: ' [', suffix: ']' })(() => [
      Insertions({ color: 'green', prefix: '+' }),
      Deletions({ color: 'red', prefix: ' -' }),
    ]),
  ]),
  Sep({ char: '|', color: '240' }),
  Context()(() => [
    Literal({ text: 'ctx: ' }),
    Tokens({ color: 'white', bold: true }),
    Literal({ text: ' (' }),
    Percent({ color: 'white' }),
    Literal({ text: ')' }),
  ]),
]);
