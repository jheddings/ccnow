# Segments

Every segment renders a single piece of data — a branch name, a token count, a
dollar amount. Compose them into any layout you want.

Segments are identified by `provider.field` expressions. The prefix determines
which provider fetches the data; the suffix picks the specific value. If a
segment has nothing to show (empty string, no data available), it silently
collapses out of the output.

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

| Segment      | Description                              | Example Output |
| ------------ | ---------------------------------------- | -------------- |
| `cost.usd`   | Session cost formatted USD               | `$12.50`       |
| `cost.total` | Raw numeric cost (float64, format $%.2f) | `$12.50`       |

`cost.total` is the raw numeric value suitable for `when` conditions
(e.g. `"when": "cost.total > 5"`).

## Speed — `speed`

| Segment        | Description                        | Example Output       |
| -------------- | ---------------------------------- | -------------------- |
| `speed.input`  | Input token throughput             | `45 t/s`, `1.2K t/s` |
| `speed.output` | Output token throughput            | `82 t/s`             |
| `speed.total`  | Combined input + output throughput | `127 t/s`            |

Speed is calculated from total tokens divided by API duration. Formatting
scales the same way as tokens: raw below 1K, `n.nK t/s` above.

## Session — `session`

| Segment                      | Description                      | Example Output  |
| ---------------------------- | -------------------------------- | --------------- |
| `session.duration.total`     | Wall-clock session time          | `2h 15m`, `45m` |
| `session.duration.api`       | Time spent on API calls          | `8m`, `1h 2m`   |
| `session.duration.total_min` | Wall-clock time in minutes (int) | `135`           |
| `session.duration.api_min`   | API time in minutes (int)        | `8`             |
| `session.id`                 | Session identifier               | `abc-123`       |
| `session.lines-added`        | Total lines added this session   | `1380`          |
| `session.lines-removed`      | Total lines removed this session | `21`            |

The `_min` variants are raw integers suitable for `when` conditions
(e.g. `"when": "session.duration.total_min > 60"`).

## Claude — `claude`

| Segment          | Description                     | Example Output |
| ---------------- | ------------------------------- | -------------- |
| `claude.version` | Claude Code application version | `2.1.75`       |
| `claude.style`   | Current output style            | `concise`      |

## System — `system`

| Segment                  | Description                        | Example Output    |
| ------------------------ | ---------------------------------- | ----------------- |
| `system.load.avg1`       | 1-minute load average (format %.2f)  | `1.42`          |
| `system.load.avg5`       | 5-minute load average (format %.2f)  | `2.10`          |
| `system.load.avg15`      | 15-minute load average (format %.2f) | `1.87`          |
| `system.mem.used`        | Used memory, human-formatted       | `12.4G`           |
| `system.mem.total`       | Total memory, human-formatted      | `32G`             |
| `system.mem.percent`     | Memory usage percentage            | `39%`             |
| `system.disk.used`       | Used disk space, human-formatted   | `234G`            |
| `system.disk.total`      | Total disk space, human-formatted  | `1T`              |
| `system.disk.percent`    | Disk usage percentage              | `23%`             |
| `system.battery.percent` | Battery charge percentage          | `85%`             |
| `system.battery.state`   | Battery state                      | `charging`        |
| `system.uptime`          | System uptime, human-formatted     | `3d 14h`, `2h 5m` |

Disk usage is measured at the mount point of the current working directory.
Battery segments return zero values on machines without a battery.

## Node Types

There are two kinds of atomic nodes:

### `expr` — Expression nodes

Evaluate an expression against the provider data environment. This is the
primary way to display provider values.

```json
{ "expr": "git.branch", "style": { "bold": true } }
{ "expr": "context.percent.used", "style": { "prefix": " (" , "suffix": ")" } }
```

### `value` — Static value nodes

Render a fixed string. Use these for separators, icons, and line breaks.

```json
{ "value": "|", "style": { "color": "240" } }
{ "value": "\n" }
{ "value": "\ue0b0", "style": { "color": "#DC0000", "bgcolor": "#3A3A3A" } }
```

## Node Properties

### `format`

Expression nodes accept an optional `format` string that controls how the raw
value is displayed. Uses Go's `fmt.Sprintf` syntax.

```json
{ "expr": "git.insertions", "format": "+%d" }
{ "expr": "context.percent.used", "format": "(%d%%)" }
```

If no format is specified, the node uses its default format (declared by
the provider) or falls back to the raw value as a string.

### `when`

Any node can conditionally show or hide based on data. See
**[Conditional Visibility](WHEN.md)** for the full reference.

```json
{ "expr": "git.branch", "when": "git.branch != '' && git.branch != 'main'" }
{ "expr": "context.percent.used", "when": "context.percent.used >= 50" }
{ "expr": "git.modified", "when": "value > 0" }
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
    { "expr": "git.branch", "style": { "bold": true } },
    {
      "expr": "git.insertions",
      "when": "value > 0",
      "style": { "color": "green", "prefix": " +" }
    },
    { "expr": "git.deletions", "when": "value > 0", "style": { "color": "red", "prefix": " -" } }
  ]
}
```

Composites support `when` expressions that can reference any provider's data,
allowing you to gate entire sections on any condition. See
[WHEN.md](WHEN.md) for details.
