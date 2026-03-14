import {
  StatusLine, Pwd, Sep, Git, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const defaultPreset: SegmentNode[] = StatusLine(() => [
  Pwd({ color: 'cyan', bold: true }),
  Sep({ char: '|', dim: true }),
  Git()(() => [
    Branch({ color: 'white', bold: true, icon: '\ue0a0 ' }),
    Literal({ text: ' [' }),
    Insertions({ color: 'green', prefix: '+' }),
    Literal({ text: ' ' }),
    Deletions({ color: 'red', prefix: '-' }),
    Literal({ text: ']' }),
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
