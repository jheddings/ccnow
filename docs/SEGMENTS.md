# Segments

Every segment renders a single piece of data — a branch name, a token count, a
dollar amount. Compose them into any layout you want.

Segments are identified by `provider.field` names. The prefix determines which
provider fetches the data; the suffix picks the specific value. If a segment has
nothing to show (empty string, no data available), it silently collapses out of
the output.

## Directory — `pwd`

| Segment     | Description                                                                | Example Output |
| ----------- | -------------------------------------------------------------------------- | -------------- |
| `pwd.name`  | Directory basename                                                         | `ccglow`       |
| `pwd.path`  | Full path prefix (everything before the basename, with trailing slash)     | `~/Projects/`  |
| `pwd.smart` | Smart-truncated path — abbreviates intermediate directories for deep paths | `~/P/…/`       |

`pwd.smart` keeps the first and last path components readable and abbreviates
the middle when nesting gets deep. Pair it with `pwd.name` for a compact but
navigable path display.

## Git — `git`

| Segment          | Description                                       | Example Output |
| ---------------- | ------------------------------------------------- | -------------- |
| `git.branch`     | Current branch name                               | `main`         |
| `git.insertions` | Lines added (staged + unstaged combined)          | `42`           |
| `git.deletions`  | Lines removed (staged + unstaged combined)        | `17`           |
| `git.modified`   | Count of modified (unstaged) files                | `3`            |
| `git.staged`     | Count of staged files                             | `2`            |
| `git.untracked`  | Count of untracked files                          | `5`            |
| `git.owner`      | Repository owner extracted from the remote URL    | `jheddings`    |
| `git.repo`       | Repository name extracted from the remote URL     | `ccglow`       |
| `git.worktree`   | Linked worktree name (empty in main working copy) | `docs-update`  |

All git segments require a git repository in the current working directory.
Remote-based segments (`git.owner`, `git.repo`) parse the `origin` remote URL
and handle both SSH and HTTPS formats. When not in a git repo, all git segments
return their zero values (`""` for strings, `0` for integers).

## Context — `context`

| Segment                     | Description                          | Example Output |
| --------------------------- | ------------------------------------ | -------------- |
| `context.tokens`            | Total token count, human-formatted   | `360K`, `1.2M` |
| `context.size`              | Context window capacity              | `1M`, `200K`   |
| `context.percent.used`      | Usage as integer percentage          | `36%`          |
| `context.percent.remaining` | Remaining capacity as percentage     | `64%`          |
| `context.input`             | Total input tokens, human-formatted  | `162K`         |
| `context.output`            | Total output tokens, human-formatted | `45K`          |

Token formatting scales automatically: raw count below 1K, `nK` for thousands,
`n.nM` for millions.

## Model — `model`

| Segment      | Description        | Example Output          |
| ------------ | ------------------ | ----------------------- |
| `model.name` | Model display name | `Opus 4.6 (1M context)` |
| `model.id`   | Model identifier   | `claude-opus-4-6`       |

## Cost — `cost`

| Segment    | Description                | Example Output |
| ---------- | -------------------------- | -------------- |
| `cost.usd` | Session cost formatted USD | `$12.50`       |

## Speed — `speed`

| Segment        | Description                        | Example Output       |
| -------------- | ---------------------------------- | -------------------- |
| `speed.input`  | Input token throughput             | `45 t/s`, `1.2K t/s` |
| `speed.output` | Output token throughput            | `82 t/s`             |
| `speed.total`  | Combined input + output throughput | `127 t/s`            |

Speed is calculated from total tokens divided by API duration. Formatting
scales the same way as tokens: raw below 1K, `n.nK t/s` above.

## Session — `session`

| Segment                  | Description                      | Example Output  |
| ------------------------ | -------------------------------- | --------------- |
| `session.duration.total` | Wall-clock session time          | `2h 15m`, `45m` |
| `session.duration.api`   | Time spent on API calls          | `8m`, `1h 2m`   |
| `session.id`             | Session identifier               | `abc-123`       |
| `session.lines-added`    | Total lines added this session   | `1380`          |
| `session.lines-removed`  | Total lines removed this session | `21`            |

## Claude — `claude`

| Segment          | Description                     | Example Output |
| ---------------- | ------------------------------- | -------------- |
| `claude.version` | Claude Code application version | `2.1.75`       |
| `claude.style`   | Current output style            | `concise`      |

## Utility Segments

These segments don't use a provider — they're structural.

| Segment   | Description                                                        |
| --------- | ------------------------------------------------------------------ |
| `literal` | Renders static text. Requires a `text` property (see below).       |
| `newline` | Renders a line break — use this for multi-line statusline layouts. |

### The `literal` segment

`literal` is the only segment that requires a property. Set `text` in the
`props` object:

```json
{
  "segment": "literal",
  "props": { "text": "|" },
  "style": { "color": "240" }
}
```

## Segment Properties

### `format`

Data segments accept an optional `format` string that controls how the raw
value is displayed. Uses Go's `fmt.Sprintf` syntax.

```json
{ "segment": "git.insertions", "format": "+%d" }
{ "segment": "context.percent.used", "format": "(%d%%)" }
```

If no format is specified, the segment uses its default format (declared by
the provider) or falls back to the raw value as a string.

### `when`

Any segment can conditionally show or hide based on data. See
**[Conditional Visibility](WHEN.md)** for the full reference.

```json
{ "segment": "git.branch", "when": "git.branch != '' && git.branch != 'main'" }
{ "segment": "context.percent.used", "when": "context.percent.used >= 50" }
{ "segment": "git.modified", "when": "value > 0" }
```

### `enabled`

Set `"enabled": false` on any node to exclude it from rendering. Disabled nodes
are skipped entirely, as if they weren't in the tree. Unlike `when`, this is a
static setting — it doesn't evaluate at runtime.

## Groups and Composites

Any node can have `children`. When it does, it acts as a composite — rendering
all children depth-first and collapsing entirely if every child produces empty
output. Use composites for sections that should appear or disappear together.

```json
{
  "when": "git.branch != ''",
  "style": { "prefix": " | " },
  "children": [
    { "segment": "git.branch", "style": { "bold": true } },
    { "segment": "git.insertions", "when": "value > 0", "style": { "color": "green", "prefix": " +" } },
    { "segment": "git.deletions", "when": "value > 0", "style": { "color": "red", "prefix": " -" } }
  ]
}
```

Composites support `when` expressions that can reference any provider's data,
allowing you to gate entire sections on any condition. See
[WHEN.md](WHEN.md) for details.
