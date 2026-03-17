# Conditional Visibility (`when`)

Nodes can show or hide based on data. Add a `when` expression to any
node — if it evaluates to `true`, the node renders; if `false`, it
collapses as if it had no data.

```json
{
  "expr": "context.percent.used",
  "when": "context.percent.used >= 50",
  "style": { "color": "yellow" }
}
```

This only shows the context percentage when usage hits 50% or higher. Below
that, the node disappears from the statusline.

## Field References

Expressions can reference three kinds of values:

### `provider.field` — Provider values

Use the full dotted name to access any provider's data. All provider data
is available in every expression — you can reference fields from any provider,
not just the one that owns the current node.

| Provider | Available Fields                                                                                                                                              |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| git      | `git.branch`, `git.insertions`, `git.deletions`, `git.modified`, `git.staged`, `git.untracked`, `git.owner`, `git.repo`, `git.worktree`                      |
| context  | `context.tokens`, `context.size`, `context.percent.used`, `context.percent.remaining`, `context.input`, `context.output`                                      |
| model    | `model.name`, `model.id`                                                                                                                                      |
| cost     | `cost.usd`, `cost.total`                                                                                                                                      |
| speed    | `speed.input`, `speed.output`, `speed.total`                                                                                                                  |
| session  | `session.duration.total`, `session.duration.api`, `session.duration.total_min`, `session.duration.api_min`, `session.lines-added`, `session.lines-removed`, `session.id` |
| pwd      | `pwd.name`, `pwd.path`, `pwd.smart`                                                                                                                           |
| claude   | `claude.version`, `claude.style`                                                                                                                               |
| system   | `system.load.avg1`, `system.load.avg5`, `system.load.avg15`, `system.mem.used`, `system.mem.total`, `system.mem.percent`, `system.disk.used`, `system.disk.total`, `system.disk.percent`, `system.battery.percent`, `system.battery.state`, `system.uptime` |

### `value` — Raw data value

The raw value from the expression this node evaluates. The type matches the
provider field — `int` for counts, `string` for text, etc. For `command` nodes,
`value` is the raw stdout string (after whitespace trimming).

```json
{ "expr": "git.modified", "when": "value > 0" }
{ "command": "cat VERSION", "when": "value != ''" }
```

### `text` — Formatted display text

The formatted but unstyled text the node would display. Always a string.
For `command` nodes, `text` is the formatted output (after applying `format`).

```json
{ "expr": "cost.usd", "when": "text != ''" }
{ "command": "echo 42", "format": "v%s", "when": "text != ''" }
```

## Operators

Full expression syntax is available, powered by
[expr-lang/expr](https://github.com/expr-lang/expr).

### Comparison

```
context.percent.used >= 50
git.branch != 'main'
value == 0
```

Operators: `==`, `!=`, `>`, `>=`, `<`, `<=`

### Boolean

```
git.insertions > 0 || git.deletions > 0
git.branch != '' && git.branch != 'main'
!(value == 0)
```

Operators: `&&`, `||`, `!`

### Arithmetic

```
git.insertions + git.deletions > 10
```

Operators: `+`, `-`, `*`, `/`, `%`

### String

```
git.branch contains 'feat'
git.branch startsWith 'fix/'
git.branch matches '^(feat|fix)/'
```

Functions: `contains`, `startsWith`, `endsWith`, `matches` (regex)

## Default Values

All provider fields always have a value — there are no nil fields. When
data is unavailable, providers return typed zero values:

- Strings → `""`
- Integers → `0`

This means you don't need nil guards. `context.percent.used >= 50` just
works — if there's no context data, the value is `0`, the comparison is
`0 >= 50` → `false`, and the node hides.

## Groups and Composites

`when` works on composite nodes too. All provider data is available in
every expression, so you can gate groups on any condition:

```json
{
  "when": "git.repo != ''",
  "children": [
    { "expr": "git.owner", "style": { "color": "240" } },
    { "expr": "git.repo", "style": { "color": "39", "prefix": "/" } }
  ]
}
```

If the condition fails, the entire group collapses — no children are
evaluated.

This is useful for "show one thing or another" patterns:

```json
{
  "when": "git.repo != ''",
  "children": [
    { "expr": "git.repo", "style": { "color": "39" } }
  ]
},
{
  "when": "git.repo == ''",
  "children": [
    { "expr": "pwd.name", "style": { "color": "39" } }
  ]
}
```

Two mutually exclusive groups: show the repo name if available, otherwise
fall back to the directory name.

### Cross-provider expressions

Because all provider data is available everywhere, you can write expressions
that reference any combination of providers:

```json
{
  "expr": "pwd.name",
  "when": "git.repo == ''",
  "style": { "color": "39" }
}
```

This shows the directory name only when there's no git repo — a `pwd`
expression gated on `git` data.

```json
{
  "expr": "context.percent.used",
  "when": "context.percent.used >= 50 && model.name != ''",
  "style": { "color": "yellow" }
}
```

Show context usage only when it's high AND we have model info available.

## Examples

**Show branch only when not on main:**

```json
{ "expr": "git.branch", "when": "git.branch != '' && git.branch != 'main'" }
```

**Show dirty indicators only when non-zero:**

```json
{ "expr": "git.modified", "when": "value > 0", "style": { "prefix": " ~" } }
{ "expr": "git.untracked", "when": "value > 0", "style": { "prefix": " ?" } }
```

**Context warning at high usage:**

```json
{ "expr": "context.percent.used", "when": "context.percent.used >= 80", "style": { "color": "red" } }
{ "expr": "context.percent.used", "when": "context.percent.used >= 50 && context.percent.used < 80", "style": { "color": "yellow" } }
```

Two copies of the same expression with mutually exclusive conditions — red
above 80%, yellow between 50-80%, hidden below 50%.

**Show speed only when available:**

```json
{ "expr": "speed.output", "when": "text != ''" }
```

For more on available segments, see [SEGMENTS.md](SEGMENTS.md).
For styling options, see [STYLE.md](STYLE.md).
