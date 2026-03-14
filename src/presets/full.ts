import {
  StatusLine, GitGroup, Group, GitBranch, GitInsertions, GitDeletions,
  ContextGroup, ContextTokens, ContextSize, ContextPercent, Literal,
  ModelName, CostUSD, SessionDuration, SessionLinesAdded, SessionLinesRemoved,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const fullPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.smart', provider: 'pwd', style: { color: '31' } },
  { type: 'pwd.name', provider: 'pwd', style: { color: '39', bold: true } },
  GitGroup({ prefix: ' | ', color: '240' })(() => [
    GitBranch({ color: 'whiteBright', bold: true, prefix: '\ue0a0 ' }),
    Group({ prefix: ' [', suffix: ']' })(() => [
      GitInsertions({ color: 'green', prefix: '+' }),
      GitDeletions({ color: 'red', prefix: ' -' }),
    ]),
  ]),
  { type: 'literal', props: { text: ' | ' }, style: { color: '240' } },
  ModelName({ color: '240' }),
  { type: 'literal', props: { text: ' · ' } },
  ContextGroup()(() => [
    ContextTokens({ color: 'white', bold: true }),
    Literal({ text: '/' }),
    ContextSize({ color: 'white' }),
    Literal({ text: ' (' }),
    ContextPercent({ color: 'white' }),
    Literal({ text: ')' }),
  ]),
  { type: 'literal', props: { text: ' · ' } },
  CostUSD({ color: 'yellow', bold: true }),
  { type: 'literal', props: { text: ' · ' } },
  SessionDuration({ color: 'magenta' }),
  { type: 'literal', props: { text: ' · ' } },
  Group()(() => [
    SessionLinesAdded({ color: 'green', prefix: '+' }),
    SessionLinesRemoved({ color: 'red', prefix: ' -' }),
  ]),
]);
