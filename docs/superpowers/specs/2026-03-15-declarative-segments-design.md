# Declarative Segments Refactor

Replace per-segment Go structs with data-driven segment resolution using
provider struct tags. Segments become a config concern, not a code concern.

## Motivation

All 25 data segment types follow the same pattern: type-assert provider,
extract a field, optionally format it. This refactor eliminates the
boilerplate by having providers declare their segment mappings via struct
tags. The render pipeline resolves data segments generically.

This also lays the groundwork for `when` expressions (#1) by making
`value` (raw provider field) and `text` (formatted output) first-class
concepts in the pipeline.

## Three Node Types

The render pipeline handles three kinds of nodes:

### LiteralSegment

Static text from config. Two registered types:
- `literal` — value from `props.text`
- `newline` — hardcoded `\n`

No provider, no `value`, no `format`. Handled by the segment registry.

### DataSegment

Bound to a provider field via struct tags. Has a `value` (raw field
data), optional `format` string, and produces `text` (formatted output).

Collapse chain:
1. `value` is nil → collapse (return nil)
2. `value` is non-nil → apply format → `text`
3. `text` is empty string → collapse (return nil)
4. `text` is non-empty → apply style → output

Not in the segment registry. Resolved generically by the render pipeline
via the tag index.

### SegmentGroup

Composite node with `Children`. Recurses into children, collapses if
all children return nil. Existing behavior, unchanged.

## Provider Struct Tags

Each provider data struct annotates fields with `segment` tags. The tag
format is `segment:"name"` or `segment:"name,format:fmt"` for fields
that need a default display format.

```go
type GitData struct {
    Branch     *string `segment:"git.branch"`
    Insertions *int    `segment:"git.insertions"`
    Deletions  *int    `segment:"git.deletions"`
    Modified   *int    `segment:"git.modified"`
    Staged     *int    `segment:"git.staged"`
    Untracked  *int    `segment:"git.untracked"`
    Owner      *string `segment:"git.owner"`
    Repo       *string `segment:"git.repo"`
    Worktree   *string `segment:"git.worktree"`
}

type ContextData struct {
    Tokens  string `segment:"context.tokens"`
    Size    string `segment:"context.size"`
    Percent *int   `segment:"context.percent,format:%d%%"`
    Input   string `segment:"context.input"`
    Output  string `segment:"context.output"`
}

type ModelData struct {
    Name *string `segment:"model.name"`
}

type CostData struct {
    USD *string `segment:"cost.usd"`
}

type SpeedData struct {
    Input  *string `segment:"speed.input"`
    Output *string `segment:"speed.output"`
    Total  *string `segment:"speed.total"`
}

type SessionData struct {
    Duration     *string `segment:"session.duration"`
    LinesAdded   *int    `segment:"session.lines-added"`
    LinesRemoved *int    `segment:"session.lines-removed"`
}

type PwdData struct {
    Name  string `segment:"pwd.name"`
    Path  string `segment:"pwd.path"`
    Smart string `segment:"pwd.smart"`
}
```

Fields without `segment` tags are not directly renderable.

## Tag Index

Built once at startup in the render package.

```go
type fieldAccessor struct {
    Provider      string
    FieldIndex    int
    DefaultFormat string // from tag, e.g., "%d%%"
}

// segment name → field accessor
type TagIndex map[string]fieldAccessor
```

### Building the index

Providers implement an optional `FieldProvider` interface to expose
their data struct type for reflection without triggering side effects:

```go
type FieldProvider interface {
    DataProvider
    Fields() any // returns a zero-value struct pointer, e.g., &GitData{}
}
```

`BuildTagIndex` iterates registered providers. For each provider that
implements `FieldProvider`, it reflects on the returned struct to
discover `segment` tags, field indices, and default formats.

Providers that do not implement `FieldProvider` are skipped (their
segments cannot be data-driven, but this is not expected to happen
with built-in providers).

### Duplicate segment names

If two providers declare the same segment name, `BuildTagIndex`
returns an error. This is a programming bug and should fail fast at
startup.

### Resolving segment values

At render time, `ResolveSegmentValues` uses the tag index to build a
flat `map[string]any` from resolved provider data:

1. For each entry in the tag index, find the provider data in the
   resolved provider map.
2. Use the cached field index to extract the field value.
3. Dereference pointer fields: nil stays nil, non-nil is dereferenced
   to the underlying value.
4. Store in the map as `segmentName → value`.

This runs once per refresh cycle, not per node.

## SegmentNode Changes

Add `Format` to `SegmentNode` in `internal/types/types.go`:

```go
type SegmentNode struct {
    Type      string         `json:"segment,omitempty"`
    Provider  string         `json:"provider,omitempty"`
    Format    string         `json:"format,omitempty"`
    Enabled   *bool          `json:"enabled,omitempty"`
    Style     *StyleAttrs    `json:"style,omitempty"`
    Props     map[string]any `json:"props,omitempty"`
    Children  []SegmentNode  `json:"children,omitempty"`
    EnabledFn func(*SessionData) bool `json:"-"`
}
```

The `Provider` field is no longer needed for data segments — the tag
index maps segment names to providers authoritatively. `Provider`
remains on the struct for composite groups (for future `when` clause
support). The `config.InferProviders` function is removed.

## Format Behavior

Format resolution order:
1. `node.Format` from config (explicit override)
2. `fieldAccessor.DefaultFormat` from struct tag
3. No format → `fmt.Sprintf("%v", value)` (strings pass through as-is)

Config-level `"format"` always overrides the tag default. This means
users can customize display without changing Go code:

```json
{ "segment": "git.insertions", "format": "+%d" }
```

### FormatValue function

```go
func FormatValue(value any, format string) string
```

If `format` is empty, returns `fmt.Sprintf("%v", value)`.
If `format` is set, returns `fmt.Sprintf(format, value)`.

## Render Pipeline Changes

### renderNode — three branches

```
renderNode(node):
    if !isEnabled(node) → nil

    if node has Children → SegmentGroup
        recurse children, collapse if all nil
        apply style to joined result

    else if segments.Get(node.Type) exists → LiteralSegment
        delegate to registered segment's Render()
        apply style

    else → DataSegment
        value = segmentValues[node.Type]
        if value is nil → nil (collapse)
        format = node.Format or tagIndex[node.Type].DefaultFormat
        text = FormatValue(value, format)
        if text is empty → nil (collapse)
        apply style to text
```

`renderNode` receives `segmentValues map[string]any` and the tag index.
The `segmentValues` map is built once per refresh by
`ResolveSegmentValues`.

### CollectProviderNames

Replaces the current implementation that reads `node.Provider`. Instead,
looks up the segment name in the tag index to find the provider name.
Literal segment types (found in the segment registry) are skipped.
This preserves lazy-provider-resolution behavior.

### config.InferProviders — removed

No longer needed. The tag index is the authoritative mapping from
segment names to providers. The `config.Parse` function drops the
`InferProviders` call. The `noProviderSegments` map is also removed.

## Preset Changes

No preset JSON changes are needed. The `context.percent` format is
declared as a default in the struct tag (`format:%d%%`), so existing
presets that use `context.percent` without an explicit format will
render identically.

Presets that want custom formatting can add `"format"` to any segment
node to override the default.

## What Gets Deleted

- All 25 data segment structs and their `Render()` methods in
  `internal/segment/segment.go`
- All data segment registrations in `RegisterBuiltin()`
- All data segment tests in `internal/segment/segment_test.go`
- `SegmentContext.Provider` field (data segments no longer need it;
  literal segments use `Props`)
- `config.InferProviders` and the `noProviderSegments` map
- The `provider` import from `internal/segment/segment.go`

## What Gets Added

- `internal/render/tagindex.go` — `TagIndex`, `BuildTagIndex`,
  `ResolveSegmentValues`
- `internal/render/format.go` — `FormatValue`
- `segment` struct tags on all provider data types
- `FieldProvider` interface in `internal/types/types.go`
- `Fields()` method on each provider implementation

## What Stays

- `Segment` interface — used by literal and newline only
- `segment.Registry` — contains only `literalSegment` and
  `newlineSegment`
- All provider implementations — unchanged except for struct tags
  and `Fields()` method
- `SegmentNode` — gains `Format`, loses reliance on `Provider` for
  data segments
- `EnabledFn` / presets with programmatic gating — unchanged

## Testing

### Tag index tests (`internal/render/tagindex_test.go`)

- Build index from a provider with tagged fields → correct mappings
- Fields without `segment` tag → not in index
- Pointer fields → correct field indices
- Tag with default format → parsed correctly
- Tag without default format → empty DefaultFormat
- Multiple providers → all segment names indexed
- Duplicate segment name across providers → error
- Provider not implementing FieldProvider → skipped
- No tags on struct → empty index for that provider
- Non-flat struct (embedded) → only top-level fields indexed

### ResolveSegmentValues tests

- Pointer field nil → nil in map
- Pointer field non-nil → dereferenced value in map
- Non-pointer field → value as-is
- Missing provider data → segment not in map

### Format tests (`internal/render/format_test.go`)

- `FormatValue(42, "")` → `"42"`
- `FormatValue("hello", "")` → `"hello"`
- `FormatValue(42, "%d%%")` → `"42%"`
- `FormatValue(nil, "")` → `""`
- `FormatValue(3.14, "%.1f")` → `"3.1"`
- `FormatValue("text", "%s!")` → `"text!"`

### Render integration tests (update `internal/render/render_test.go`)

- DataSegment resolves from segmentValues map → formatted output
- DataSegment with nil value → collapses
- DataSegment with config format → overrides default
- DataSegment with tag default format → applied
- DataSegment with no format → raw value as string
- LiteralSegment → unchanged behavior
- SegmentGroup → collapse behavior preserved
- Unknown segment name (not in registry or tag index) → collapses
- CollectProviderNames uses tag index → correct provider set
- Non-pointer string field with empty value → collapses

### Existing tests

- Provider tests → unchanged (providers don't change behavior)
- Preset tests → unchanged (no preset JSON changes)
- Existing render tests → update to new pipeline signature
