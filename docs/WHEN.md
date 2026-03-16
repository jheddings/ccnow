# Conditional Visibility (`when`)

Segments can show or hide based on their data. Add a `when` expression to any
segment node — if it evaluates to `true`, the segment renders; if `false`, it
collapses as if it had no data.

```json
{
  "segment": "context.percent.used",
  "when": ".percent >= 50",
  "style": { "color": "yellow" }
}
```

This only shows the context percentage when usage hits 50% or higher. Below
that, the segment disappears from the statusline.

## Field References

Expressions can reference three kinds of values:

### `.field` — Provider fields

Dot-prefixed names access fields from the segment's provider. These are the
Go struct field names, lowercased.

| Provider | Available Fields |
|----------|-----------------|
| git      | `.branch`, `.insertions`, `.deletions`, `.modified`, `.staged`, `.untracked`, `.owner`, `.repo`, `.worktree` |
| context  | `.tokens`, `.size`, `.percent`, `.remaining`, `.input`, `.output` |
| model    | `.name`, `.id` |
| cost     | `.usd` |
| speed    | `.input`, `.output`, `.total` |
| session  | `.duration`, `.apiduration`, `.linesadded`, `.linesremoved`, `.id` |
| pwd      | `.name`, `.path`, `.smart` |
| claude   | `.version`, `.style` |

### `value` — Raw data value

The raw value from the provider field this segment maps to. The type matches
the provider field — `int` for counts, `string` for text, etc.

```json
{ "segment": "git.modified", "when": "value > 0" }
```

### `text` — Formatted display text

The formatted but unstyled text the segment would display. Always a string.

```json
{ "segment": "cost.usd", "when": "text != ''" }
```

## Operators

Full expression syntax is available, powered by
[expr-lang/expr](https://github.com/expr-lang/expr).

### Comparison

```
.percent >= 50
.branch != 'main'
value == 0
```

Operators: `==`, `!=`, `>`, `>=`, `<`, `<=`

### Boolean

```
.insertions > 0 || .deletions > 0
.branch != '' && .branch != 'main'
!(value == 0)
```

Operators: `&&`, `||`, `!`

### Arithmetic

```
.insertions + .deletions > 10
```

Operators: `+`, `-`, `*`, `/`, `%`

### String

```
.branch contains 'feat'
.branch startsWith 'fix/'
.branch matches '^(feat|fix)/'
```

Functions: `contains`, `startsWith`, `endsWith`, `matches` (regex)

## Nil Handling

Optional provider fields (pointer types like `*int`, `*string`) are
automatically coerced when nil:

- Nil `*int` → `0`
- Nil `*string` → `""`

This means you don't need nil guards. `.percent >= 50` just works — if
there's no context data, `.percent` is `0`, the comparison is `0 >= 50`
→ `false`, and the segment hides.

## Groups and Composites

`when` works on composite nodes too. Set `provider` explicitly so the
expression knows which data to evaluate against:

```json
{
  "provider": "git",
  "when": ".repo != ''",
  "children": [
    { "segment": "git.owner", "style": { "color": "240" } },
    { "segment": "git.repo", "style": { "color": "39", "prefix": "/" } }
  ]
}
```

If the condition fails, the entire group collapses — no children are
evaluated.

This is useful for "show one thing or another" patterns:

```json
{
  "provider": "git",
  "when": ".repo != ''",
  "children": [
    { "segment": "git.repo", "style": { "color": "39" } }
  ]
},
{
  "provider": "git",
  "when": ".repo == ''",
  "children": [
    { "segment": "pwd.name", "style": { "color": "39" } }
  ]
}
```

Two mutually exclusive groups: show the repo name if available, otherwise
fall back to the directory name.

## Examples

**Show branch only when not on main:**
```json
{ "segment": "git.branch", "when": ".branch != '' && .branch != 'main'" }
```

**Show dirty indicators only when non-zero:**
```json
{ "segment": "git.modified", "when": "value > 0", "style": { "prefix": " ~" } }
{ "segment": "git.untracked", "when": "value > 0", "style": { "prefix": " ?" } }
```

**Context warning at high usage:**
```json
{ "segment": "context.percent.used", "when": ".percent >= 80", "style": { "color": "red" } }
{ "segment": "context.percent.used", "when": ".percent >= 50 && .percent < 80", "style": { "color": "yellow" } }
```

Two copies of the same segment with mutually exclusive conditions — red
above 80%, yellow between 50-80%, hidden below 50%.

**Show speed only when available:**
```json
{ "segment": "speed.output", "when": "text != ''" }
```

For more on available segments, see [SEGMENTS.md](SEGMENTS.md).
For styling options, see [STYLE.md](STYLE.md).
