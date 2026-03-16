# Conditional `when` Expression Syntax

Adds a `when` clause to `SegmentNode` that controls segment visibility
based on provider data or rendered values, using `expr-lang/expr` for
expression evaluation.

Resolves #1.

## Syntax

The `when` field is a string expression on `SegmentNode`, evaluated by
`expr-lang/expr`.

```json
{
  "segment": "context.percent.used",
  "when": ".percent >= 50",
  "style": { "color": "red" }
}
```

### Field references

- **`.field`** — resolves against the segment's provider data struct.
  Registered as variables in the expr environment using dot-prefixed
  lowercase names derived from the provider's exported Go field names.
  Examples: `.percent`, `.insertions`, `.branch`.
- **`value`** — the segment's raw data value from the provider field
  it maps to (via struct tags). For a data segment like
  `context.percent.used`, `value` is the `*int` Percent field
  (dereferenced). Type depends on the provider field.
- **`text`** — the segment's formatted text output (the result of
  `FormatValue(value, format)`), available as a string. This is the
  unstyled display text.

### Nil pointer coercion

Provider data structs use pointer types for optional fields (`*int`,
`*string`). To keep expressions simple, `BuildEnv` coerces nil
pointers to typed zero values:

- `*int` nil → `0`
- `*string` nil → `""`
- `*float64` nil → `0.0`

Non-nil pointers are dereferenced to their underlying value.

This means `.percent >= 50` works without nil guards — a nil
`*int` becomes `0`, so the comparison is `0 >= 50` → false, and the
segment hides. This matches the intuition that "no data = don't show."

### Expression language

Full `expr-lang/expr` syntax is available:

- **Comparison:** `>=`, `<=`, `>`, `<`, `==`, `!=`
- **Boolean:** `&&`, `||`, `!`
- **Arithmetic:** `+`, `-`, `*`, `/` (e.g., `.insertions + .deletions > 10`)
- **String ops:** `contains`, `startsWith`, `endsWith`, `matches`

Examples:
- `".percent >= 50"` — show when context usage is >= 50%
- `".insertions > 0 || .deletions > 0"` — show when there are changes
- `".branch != ''"` — show when branch data is available
- `"text != ''"` — show when segment produces display text
- `"value > 0"` — show when the raw data value is positive

### Composite nodes

`when` is valid on composite nodes (nodes with `Children`). To use
`.field` references on a composite, the node must have an explicit
`provider` field set so the render pipeline knows which provider's
data to use for the environment.

`value` and `text` are not meaningful on composite nodes. `value`
evaluates as `nil` (coerced to zero) and `text` as `""`.

Example — hide the entire git group when not in a repo:
```json
{
  "provider": "git",
  "when": ".branch != ''",
  "children": [
    { "segment": "git.branch", "style": { "prefix": " " } },
    { "segment": "git.insertions", "style": { "prefix": " +" } }
  ]
}
```

### Behavior

- `when` is empty → segment always renders (current behavior)
- Expression evaluates to true → segment renders normally
- Expression evaluates to false → segment returns nil (collapses)
- Expression compilation fails → treat as false, log to stderr
- Provider data is nil → all dot-field variables get zero values;
  `value` is coerced to zero

## SegmentNode Changes

Add to `SegmentNode` in `internal/types/types.go`:

```go
When string `json:"when,omitempty"`
```

## New Package: `internal/condition/`

Create `internal/condition/condition.go` with:

### `Compile`

```go
func Compile(expression string) (*Condition, error)
```

Compiles the expression string into a reusable `Condition`. Returns an
error if the expression is syntactically invalid. Empty string returns
a nil `*Condition` (always true).

### `Condition.Evaluate`

```go
func (c *Condition) Evaluate(env map[string]any) bool
```

Runs the compiled expression against the provided environment map.
Returns true if the result is boolean `true`. Any other result type
or runtime error returns false.

### `BuildEnv`

```go
func BuildEnv(providerData any, value any, text string) map[string]any
```

Builds the variable environment for expression evaluation:

1. If `providerData` is a pointer, dereference to the struct.
2. For each exported field, register a dot-prefixed lowercase variable
   using the Go field name (e.g., `Branch` → `.branch`).
3. Dereference pointer fields with nil coercion:
   - `*int` nil → `0`, non-nil → `int` value
   - `*string` nil → `""`, non-nil → `string` value
   - `*float64` nil → `0.0`, non-nil → `float64` value
4. Non-pointer fields registered as-is.
5. Register `value` from the `value` parameter (nil-coerced same
   as pointer fields).
6. Register `text` from the `text` parameter.

If `providerData` is nil or not a struct/pointer-to-struct, return
a map containing only `value` and `text`.

Note: field names use the Go exported name lowercased, not JSON tags
or segment tags. Provider field names are already short and clear
(e.g., `Branch`, `Percent`, `Tokens`).

## Render Pipeline Changes

Modify `internal/render/render.go` to evaluate `when` during tree
traversal.

### Data segments

For data segments, the pipeline already has:
- `value` from `segmentValues[node.Type]`
- `text` from `FormatValue(value, format)`
- Provider data from `providerData[tagIdx[node.Type].Provider]`

Evaluation happens after format but before style:

```
value = segmentValues[node.Type]
if value is nil → collapse
text = FormatValue(value, format)
if text is empty → collapse
if node.When != "" →
    provider = providerData[tagIdx[node.Type].Provider]
    env = BuildEnv(provider, value, text)
    if !condition.Evaluate(env) → collapse
apply style → output
```

### Literal segments

Literal segments can use `when` with `text` (their rendered output).
`value` is nil. Provider data is nil. Only `text`-based conditions
are meaningful.

### Composite nodes

For composites with `when`, evaluate before recursing into children.
Provider data is resolved from the explicit `provider` field on the
node. `value` is nil (coerced to zero), `text` is `""`.

```
if node.When != "" →
    provider = providerData[node.Provider]
    env = BuildEnv(provider, nil, "")
    if !condition.Evaluate(env) → collapse (skip children)
recurse children...
```

### Interaction with `Enabled` / `EnabledFn`

`when` is an additional gate. Evaluation order:
1. `Enabled` / `EnabledFn` → if false, skip (existing behavior)
2. For composites: evaluate `when` → if false, skip children
3. For data segments: extract value, format, evaluate `when`
4. Apply style

### Compilation caching

Compile `when` expressions into a `map[string]*Condition` keyed by
the expression string, built once during render tree setup. This
avoids adding runtime state to `SegmentNode` (which stays a pure
config struct) and deduplicates identical expressions across nodes.

### Error handling

Compilation errors are logged to stderr with `log.Printf` and the
segment is treated as hidden. Runtime evaluation errors also return
false (segment hidden).

## Dependency

Add `github.com/expr-lang/expr` to `go.mod`.

## Testing

### Condition package tests (`internal/condition/condition_test.go`)

**Compile:**
- Valid expression compiles without error
- Invalid expression returns error
- Empty expression returns nil Condition

**BuildEnv:**
- Provider struct with string, int, `*int`, `*string` fields →
  correct dot-prefixed variables
- Nil pointer fields → coerced zero values (0, "")
- Non-nil pointer fields → dereferenced values
- Nil provider → map with only `value` and `text`
- Non-struct provider → map with only `value` and `text`
- `value` parameter passed through correctly
- `text` parameter passed through correctly

**Evaluate:**
- Numeric comparisons: `>=`, `>`, `<`, `<=`, `==`, `!=`
- String comparisons: `==`, `!=`
- Boolean combinators: `&&`, `||`
- Nil-coerced field: `.field >= 50` where field was nil `*int`
- `value` against raw data
- `text` against formatted output
- Empty expression (nil Condition) → true
- Expression returning non-bool → false
- Runtime error → false

### Render integration tests (`internal/render/render_test.go`)

- Data segment with `when` that passes → renders normally
- Data segment with `when` that fails → collapses
- Data segment with `value` condition
- Data segment with `text` condition
- Composite with `when` and `provider` → evaluates against provider
- Composite with `when` that fails → subtree skipped
- Literal segment with `when` on `text`
- Segment with no `when` → unchanged behavior
- Segment with invalid `when` → treated as hidden

## What Doesn't Change

- Existing `Enabled` / `EnabledFn` behavior — `when` is additive
- Tag index, segment value resolution, formatting — all unchanged
- All existing presets and configs — `when` is optional
