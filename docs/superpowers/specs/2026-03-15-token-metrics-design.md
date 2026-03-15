# Token Metrics Segments

Adds input/output token count segments to the context provider and a new
speed provider with input/output/total token speed segments.

Resolves #27 (token speed) and #33 (separate input/output token counts).

## SessionData Changes

Update `internal/types/types.go` to parse additional fields from the Claude
Code session JSON.

Add to `ContextWindow`:
```go
TotalInputTokens  *int `json:"total_input_tokens,omitempty"`
TotalOutputTokens *int `json:"total_output_tokens,omitempty"`
```

Add to `CurrentUsage`:
```go
OutputTokens int `json:"output_tokens"`
```

The `OutputTokens` field on `CurrentUsage` is added for JSON parsing
completeness. It is not used by any provider in this spec but ensures we
capture the full session JSON shape for future use.

## Context Provider Changes

Add two new fields to `ContextData` in `internal/provider/context.go`:

```go
type ContextData struct {
    Tokens  string
    Size    string
    Percent *int
    Input   string // formatted input token count
    Output  string // formatted output token count
}
```

These are bare strings following the existing pattern in `ContextData`
(`Tokens` and `Size` are also bare strings). Segments check for `!= ""`
and return nil when empty, consistent with `contextTokensSegment` and
`contextSizeSegment`.

Populating:
- `Input`: use `*TotalInputTokens` if non-nil; fall back to summing
  `InputTokens + CacheCreationInputTokens + CacheReadInputTokens` from
  `CurrentUsage`. Format with `FormatTokens()`. Note: `TotalInputTokens`
  is `*int` and must be dereferenced with a nil check before passing to
  `FormatTokens(int)`.
- `Output`: use `*TotalOutputTokens` if non-nil. Format with
  `FormatTokens()`. Leave empty when `TotalOutputTokens` is nil.

## New Speed Provider

Create `internal/provider/speed.go` with a new `speed` provider.

```go
type SpeedData struct {
    Input  *string // formatted input tokens/sec, e.g., "42 t/s"
    Output *string // formatted output tokens/sec
    Total  *string // formatted total tokens/sec
}
```

The speed provider uses `*string` pointers (unlike `ContextData`'s bare
strings) because speed data may be entirely unavailable — there's no
meaningful zero/empty state for speed the way there is for token counts.

### Cross-struct data dependency

Speed calculation requires data from two independent top-level structs:
- Token counts from `session.ContextWindow` (`TotalInputTokens`,
  `TotalOutputTokens`)
- Duration from `session.Cost` (`TotalAPIDurationMS`)

Both can independently be nil. The provider must guard against:
1. `session.ContextWindow` being nil → all speed fields remain nil
2. `session.Cost` being nil → all speed fields remain nil
3. `TotalAPIDurationMS` being zero → all speed fields remain nil
4. Individual token count fields (`*int`) being nil → that specific
   speed field remains nil

### Speed calculations

All calculations use `float64` arithmetic:
- `Input`: `float64(*TotalInputTokens) * 1000.0 / TotalAPIDurationMS`
- `Output`: `float64(*TotalOutputTokens) * 1000.0 / TotalAPIDurationMS`
- `Total`: `float64(*TotalInputTokens + *TotalOutputTokens) * 1000.0 / TotalAPIDurationMS`

Total is only computed when both input and output token counts are
available.

### FormatSpeed

Add `FormatSpeed(tokensPerSec float64) string` in
`internal/provider/speed.go`. Returns formatted values like:
- Below 1000: integer display, e.g., `"42 t/s"`
- 1000+: one decimal K, e.g., `"1.2K t/s"`

Register the provider in `internal/provider/provider.go`.

## New Segments

Five new segment types in `internal/segment/segment.go`:

| Segment | Provider | Field | Render |
|---------|----------|-------|--------|
| `context.input` | context | Input | String, nil when empty |
| `context.output` | context | Output | String, nil when empty |
| `speed.input` | speed | Input | *string as-is, nil when nil |
| `speed.output` | speed | Output | *string as-is, nil when nil |
| `speed.total` | speed | Total | *string as-is, nil when nil |

Context segments check `data.Input != ""` (bare string pattern).
Speed segments check `data.Input != nil` (pointer pattern).

Provider wiring: the render pipeline resolves providers by matching the
segment name prefix to the provider name (e.g., `speed.input` →
`speed` provider). The new `speed` provider is registered in
`provider.go` alongside existing providers.

## F1 Preset Update

Add `speed.output` to line 2 of the F1 preset, within the existing
context block (same `#2D2D2D` background). Insert after `context.percent`
and before the powerline separator that transitions to the cost block.
No new separator needed — it stays within the same visual block.

## Testing

### Context provider tests (extend `internal/provider/provider_test.go`)

- Verify `Input` is populated from `TotalInputTokens` when available
- Verify `Input` falls back to sum of cache fields when
  `TotalInputTokens` is nil
- Verify `Output` is populated from `TotalOutputTokens`
- Verify `Output` is empty when `TotalOutputTokens` is nil

### Speed provider tests (`internal/provider/speed_test.go`, new file)

- Verify speed calculations with known values (e.g., 10000 input tokens
  over 5000ms → 2000 t/s → "2K t/s")
- Verify nil fields when `TotalAPIDurationMS` is zero
- Verify nil fields when `ContextWindow` is nil
- Verify nil fields when `Cost` is nil
- Verify nil fields when individual token counts are nil
- Verify `FormatSpeed` output: below 1000, above 1000

### Segment tests (extend `internal/segment/segment_test.go`)

- Test each new segment with populated data (verify rendered string)
- Test each new segment with nil/empty fields (verify nil return)

## What Doesn't Change

- Existing `context.tokens`, `context.size`, `context.percent` segments
  unchanged
- Existing provider registry pattern — speed is a new provider alongside
  the existing ones
