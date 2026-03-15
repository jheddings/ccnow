# Segments

Every segment renders a single piece of data — a branch name, a token count, a
dollar amount. Compose them into any layout you want.

Segments are identified by `provider.field` names. The prefix determines which
provider fetches the data; the suffix picks the specific value. If a segment has
nothing to show (no git repo, no cost data yet), it returns empty and silently
collapses out of the output.

## Directory — `pwd`

| Segment     | Description                                                                 | Example Output        |
| ----------- | --------------------------------------------------------------------------- | --------------------- |
| `pwd.name`  | Directory basename                                                          | `ccglow`              |
| `pwd.path`  | Full path prefix (everything before the basename, with trailing slash)      | `~/Projects/`         |
| `pwd.smart` | Smart-truncated path — abbreviates intermediate directories for deep paths  | `~/P/…/`              |

`pwd.smart` keeps the first and last path components readable and abbreviates
the middle when nesting gets deep. Pair it with `pwd.name` for a compact but
navigable path display.

## Git — `git`

| Segment          | Description                                          | Example Output |
| ---------------- | ---------------------------------------------------- | -------------- |
| `git.branch`     | Current branch name                                  | `main`         |
| `git.insertions` | Lines added (staged + unstaged combined)              | `42`           |
| `git.deletions`  | Lines removed (staged + unstaged combined)            | `17`           |
| `git.modified`   | Count of modified (unstaged) files                   | `3`            |
| `git.staged`     | Count of staged files                                | `2`            |
| `git.untracked`  | Count of untracked files                             | `5`            |
| `git.owner`      | Repository owner extracted from the remote URL       | `jheddings`    |
| `git.repo`       | Repository name extracted from the remote URL        | `ccglow`       |
| `git.worktree`   | Linked worktree name (empty in main working copy)    | `docs-update`  |

All git segments require a git repository in the current working directory.
Remote-based segments (`git.owner`, `git.repo`) parse the `origin` remote URL
and handle both SSH and HTTPS formats.

## Context — `context`

| Segment           | Description                            | Example Output |
| ----------------- | -------------------------------------- | -------------- |
| `context.tokens`  | Total token count, human-formatted     | `360K`, `1.2M` |
| `context.size`    | Context window capacity                | `1M`, `200K`   |
| `context.percent` | Usage as integer percentage            | `36%`          |

Token formatting scales automatically: raw count below 1K, `nK` for thousands,
`n.nM` for millions.

## Model — `model`

| Segment      | Description        | Example Output          |
| ------------ | ------------------ | ----------------------- |
| `model.name` | Model display name | `Opus 4.6 (1M context)` |

## Cost — `cost`

| Segment    | Description                | Example Output |
| ---------- | -------------------------- | -------------- |
| `cost.usd` | Session cost formatted USD | `$12.50`       |

## Session — `session`

| Segment                | Description                          | Example Output |
| ---------------------- | ------------------------------------ | -------------- |
| `session.duration`     | Wall-clock session time              | `2h 15m`, `45m`|
| `session.lines-added`  | Total lines added this session       | `1380`         |
| `session.lines-removed`| Total lines removed this session     | `21`           |

## Utility Segments

These segments don't use a provider — they're structural.

| Segment   | Description                                                        |
| --------- | ------------------------------------------------------------------ |
| `literal` | Renders static text. Requires a `text` property (see below).       |
| `newline`  | Renders a line break — use this for multi-line statusline layouts. |
| `group`   | Container for child segments. Auto-collapses if all children are empty. |

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

### Groups and composites

Any segment node can have `children`. When it does, it acts as a composite —
rendering all children depth-first and collapsing entirely if every child
produces empty output. Use the `group` type for pure containers, or attach
children to a data segment to gate a whole section on that provider's
availability.

```json
{
  "segment": "git",
  "style": { "prefix": " | " },
  "children": [
    { "segment": "git.branch", "style": { "bold": true } },
    { "segment": "git.insertions", "style": { "color": "green", "prefix": " +" } },
    { "segment": "git.deletions", "style": { "color": "red", "prefix": " -" } }
  ]
}
```

If there's no git repo, the entire group disappears — no stray separators, no
empty brackets.

### Disabling segments

Set `"enabled": false` on any node to exclude it from rendering. Disabled nodes
are skipped entirely, as if they weren't in the tree.

## Provider Auto-Wiring

You don't need to set `"provider"` explicitly. The segment type prefix
determines the provider automatically:

- `git.branch` → provider `git`
- `context.tokens` → provider `context`
- `pwd.name` → provider `pwd`

Each provider fetches its data once and caches it — so ten git segments don't
mean ten calls to `git status`.
