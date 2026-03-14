import {
  StatusLine, Sep, Git, Group, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const defaultPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.smart', provider: 'pwd', style: { color: 'cyan' } },
  { type: 'pwd.name', provider: 'pwd', style: { color: 'cyanBright', bold: true } },
  Sep({ char: '|', dim: true }),
  Git()(() => [
    Branch({ color: 'whiteBright', bold: true, icon: '\ue0a0 ' }),
    Group({ prefix: ' [', suffix: ']' })(() => [
      Insertions({ color: 'green', prefix: '+' }),
      Deletions({ color: 'red', prefix: ' -' }),
    ]),
  ]),
  Sep({ char: '|', dim: true }),
  Context()(() => [
    Literal({ text: 'ctx: ' }),
    Tokens({ bold: true }),
    Literal({ text: ' (' }),
    Percent(),
    Literal({ text: ')' }),
  ]),
]);
