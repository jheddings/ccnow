import {
  StatusLine,
  Group,
  GitGroup,
  GitBranch,
  GitInsertions,
  GitDeletions,
  ContextGroup,
  ContextTokens,
  ContextSize,
  ContextPercent,
  Literal,
  ModelName,
  CostUSD,
  SessionDuration,
  SessionLinesAdded,
  SessionLinesRemoved,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const fullPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.smart', provider: 'pwd', style: { color: '31' } },
  { type: 'pwd.name', provider: 'pwd', style: { color: '39', bold: true } },
  GitGroup({ prefix: ' | ', color: '240' })(() => [
    GitBranch({ color: 'whiteBright', bold: true, prefix: '\ue0a0 ' }),
    Group({ prefix: ' ·' })(() => [
      GitInsertions({ color: 'green', prefix: ' +' }),
      GitDeletions({ color: 'red', prefix: ' -' }),
    ]),
  ]),
  Group({ prefix: ' | ', color: '240' })(() => [ModelName({ color: '240' })]),
  Group({ prefix: ' ·' })(() => [
    ContextGroup({ prefix: ' ' })(() => [
      ContextTokens({ color: 'white', bold: true }),
      Literal({ text: '/' }),
      ContextSize({ color: 'white' }),
      Literal({ text: ' (' }),
      ContextPercent({ color: 'white' }),
      Literal({ text: ')' }),
    ]),
  ]),
  Group({ prefix: ' ·' })(() => [CostUSD({ color: 'yellow', bold: true, prefix: ' ' })]),
  Group({ prefix: ' ·' })(() => [SessionDuration({ color: 'magenta', prefix: ' ' })]),
  Group({ prefix: ' ·' })(() => [
    SessionLinesAdded({ color: 'green', prefix: ' +' }),
    SessionLinesRemoved({ color: 'red', prefix: ' -' }),
  ]),
]);
