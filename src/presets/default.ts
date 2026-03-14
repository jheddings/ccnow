import {
  StatusLine, Sep, Git, Group, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const defaultPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.smart', provider: 'pwd', style: { color: '31', suffix: '/' } },
  { type: 'pwd.name', provider: 'pwd', style: { color: '39', bold: true } },
  Sep({ char: '|', color: '240' }),
  Git()(() => [
    Branch({ color: 'whiteBright', bold: true, icon: '\ue0a0 ' }),
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
