# New Payload Segments Design

Add four new segments sourced directly from the Claude Code statusline JSON
payload: `claude.version`, `claude.style`, `model.id`, and `context.remaining`.

## Motivation

Several fields in the session JSON payload are not yet exposed as segments.
These are all straightforward to add since the data is already available —
no external I/O or transcript parsing required.

## New Segments

| Segment | Provider | Source field | Example output |
|---|---|---|---|
| `claude.version` | `claude` | `version` | `2.1.75` |
| `claude.style` | `claude` | `output_style.name` | `default` |
| `model.id` | `model` | `model.id` | `claude-opus-4-6[1m]` |
| `context.remaining` | `context` | `context_window.remaining_percentage` | `96%` |

## Type Changes

Add to `SessionData`:

```go
Version     string           `json:"version,omitempty"`
OutputStyle *OutputStyleInfo  `json:"output_style,omitempty"`
```

New struct:

```go
type OutputStyleInfo struct {
    Name string `json:"name"`
}
```

Add to `ContextWindow`:

```go
RemainingPercentage int `json:"remaining_percentage"`
```

## New Provider: `claude`

File: `internal/provider/claude.go`

```go
type ClaudeData struct {
    Version *string
    Style   *string
}
```

Resolve reads `session.Version` and `session.OutputStyle.Name`. Returns nil
pointers for missing data so segments fail silent.

## Extended Providers

**model** — Add `ID *string` to `ModelData`. Populate from `session.Model.ID`.

**context** — Add `Remaining *string` to `ContextData`. Format as `"96%"`
matching the existing `Percent` field pattern.

## Segments

All four segments follow the existing pattern: cast provider data, return
the relevant pointer (nil when no data).

Register all four in `RegisterBuiltin()` in `segment.go`.

## Presets

Add `claude.version` and `claude.style` to the `full` preset. Leave
`default` and `minimal` unchanged.

## Testing

- Unit tests for `claude` provider (version present, style present, both nil)
- Unit tests for extended `model` provider (ID present, ID empty)
- Unit tests for extended `context` provider (remaining percentage)
- Unit tests for all four new segments (data present, data nil)
