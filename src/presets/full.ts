import {
  StatusLine, Sep, Git, Group, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const fullPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.path', provider: 'pwd', style: { color: '#00afff' } },
  { type: 'pwd.name', provider: 'pwd', style: { color: '#00afff', bold: true } },
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
