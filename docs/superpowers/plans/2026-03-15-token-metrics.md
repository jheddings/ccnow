# Token Metrics Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add input/output token count segments and a speed provider with input/output/total token speed segments.

**Architecture:** Extend `SessionData` types to parse new JSON fields. Add `Input`/`Output` fields to the existing context provider. Create a new `speed` provider that computes session-average token throughput from cross-struct data (ContextWindow + Cost). Add five new segments and update the F1 preset.

**Tech Stack:** Go, standard library testing

**Spec:** `docs/superpowers/specs/2026-03-15-token-metrics-design.md`

---

## Chunk 1: Types and context provider

### Task 1: Add new fields to SessionData types

**Files:**
- Modify: `internal/types/types.go:27-38`

- [ ] **Step 1: Add fields to ContextWindow and CurrentUsage**

```go
// ContextWindow contains token usage data.
type ContextWindow struct {
	UsedPercentage    int           `json:"used_percentage"`
	ContextWindowSize int           `json:"context_window_size,omitempty"`
	TotalInputTokens  *int          `json:"total_input_tokens,omitempty"`
	TotalOutputTokens *int          `json:"total_output_tokens,omitempty"`
	CurrentUsage      *CurrentUsage `json:"current_usage,omitempty"`
}

// CurrentUsage breaks down token counts by category.
type CurrentUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}
```

- [ ] **Step 2: Run existing tests to verify nothing breaks**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/types/types.go
git commit -m "feat(types): add input/output token fields to session types"
```

### Task 2: Add Input/Output to context provider

**Files:**
- Modify: `internal/provider/context.go`
- Modify: `internal/provider/provider_test.go`

- [ ] **Step 1: Write failing tests for context Input/Output**

Add to `internal/provider/provider_test.go`:

```go
func TestContextProviderWithTotalTokens(t *testing.T) {
	p := &contextProvider{}
	inputTokens := 50000
	outputTokens := 8000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage:    36,
			ContextWindowSize: 1000000,
			TotalInputTokens:  &inputTokens,
			TotalOutputTokens: &outputTokens,
			CurrentUsage: &types.CurrentUsage{
				InputTokens:              100,
				CacheCreationInputTokens: 200,
				CacheReadInputTokens:     300,
			},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	// Input should prefer TotalInputTokens over summing cache fields
	if data.Input != "50K" {
		t.Errorf("expected Input 50K, got %s", data.Input)
	}
	if data.Output != "8K" {
		t.Errorf("expected Output 8K, got %s", data.Output)
	}
}

func TestContextProviderInputFallback(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage: 10,
			CurrentUsage: &types.CurrentUsage{
				InputTokens:              100,
				CacheCreationInputTokens: 200,
				CacheReadInputTokens:     300,
			},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	// No TotalInputTokens — should fall back to sum of cache fields
	if data.Input != "600" {
		t.Errorf("expected Input 600, got %s", data.Input)
	}
	// No TotalOutputTokens — should be empty
	if data.Output != "" {
		t.Errorf("expected empty Output, got %s", data.Output)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/provider/ -run "TestContextProvider(WithTotalTokens|InputFallback)" -v`
Expected: FAIL (Input/Output fields don't exist)

- [ ] **Step 3: Implement Input/Output in context provider**

Update `internal/provider/context.go`:

```go
type ContextData struct {
	Tokens  string
	Size    string
	Percent *int
	Input   string
	Output  string
}
```

Add to `Resolve()` after the existing `Percent` logic, before the return:

```go
	// Input tokens: prefer TotalInputTokens, fall back to cache field sum
	if cw.TotalInputTokens != nil {
		data.Input = FormatTokens(*cw.TotalInputTokens)
	} else if totalTokens > 0 {
		data.Input = FormatTokens(totalTokens)
	}

	// Output tokens: only available via TotalOutputTokens
	if cw.TotalOutputTokens != nil {
		data.Output = FormatTokens(*cw.TotalOutputTokens)
	}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/provider/ -run "TestContextProvider" -v`
Expected: PASS (all context tests including existing ones)

- [ ] **Step 5: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/provider/context.go internal/provider/provider_test.go
git commit -m "feat(context): add input and output token count fields"
```

### Task 3: Create speed provider

**Files:**
- Create: `internal/provider/speed.go`
- Create: `internal/provider/speed_test.go`
- Modify: `internal/provider/provider.go`

- [ ] **Step 1: Write failing tests for speed provider**

Create `internal/provider/speed_test.go`:

```go
package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestSpeedProvider(t *testing.T) {
	p := &speedProvider{}
	inputTokens := 10000
	outputTokens := 5000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			TotalInputTokens:  &inputTokens,
			TotalOutputTokens: &outputTokens,
		},
		Cost: &types.CostInfo{
			TotalAPIDurationMS: 5000, // 5 seconds
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*SpeedData)
	// 10000 tokens / 5s = 2000 t/s = "2.0K t/s"
	if data.Input == nil || *data.Input != "2.0K t/s" {
		t.Errorf("expected Input '2.0K t/s', got %v", data.Input)
	}
	// 5000 tokens / 5s = 1000 t/s = "1.0K t/s"
	if data.Output == nil || *data.Output != "1.0K t/s" {
		t.Errorf("expected Output '1.0K t/s', got %v", data.Output)
	}
	// 15000 tokens / 5s = 3000 t/s = "3.0K t/s"
	if data.Total == nil || *data.Total != "3.0K t/s" {
		t.Errorf("expected Total '3.0K t/s', got %v", data.Total)
	}
}

func TestSpeedProviderZeroDuration(t *testing.T) {
	p := &speedProvider{}
	inputTokens := 10000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			TotalInputTokens: &inputTokens,
		},
		Cost: &types.CostInfo{
			TotalAPIDurationMS: 0,
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*SpeedData)
	if data.Input != nil {
		t.Errorf("expected nil Input for zero duration, got %v", *data.Input)
	}
	if data.Output != nil {
		t.Errorf("expected nil Output for zero duration, got %v", *data.Output)
	}
	if data.Total != nil {
		t.Errorf("expected nil Total for zero duration, got %v", *data.Total)
	}
}

func TestSpeedProviderNilContextWindow(t *testing.T) {
	p := &speedProvider{}
	sess := &types.SessionData{
		CWD:  "/tmp",
		Cost: &types.CostInfo{TotalAPIDurationMS: 5000},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*SpeedData)
	if data.Input != nil || data.Output != nil || data.Total != nil {
		t.Error("expected all nil fields when ContextWindow is nil")
	}
}

func TestSpeedProviderNilCost(t *testing.T) {
	p := &speedProvider{}
	inputTokens := 10000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			TotalInputTokens: &inputTokens,
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*SpeedData)
	if data.Input != nil || data.Output != nil || data.Total != nil {
		t.Error("expected all nil fields when Cost is nil")
	}
}

func TestSpeedProviderPartialTokens(t *testing.T) {
	p := &speedProvider{}
	inputTokens := 10000
	// No output tokens
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			TotalInputTokens: &inputTokens,
		},
		Cost: &types.CostInfo{
			TotalAPIDurationMS: 2000,
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*SpeedData)
	// 10000 / 2s = 5000 t/s = "5.0K t/s"
	if data.Input == nil || *data.Input != "5.0K t/s" {
		t.Errorf("expected Input '5.0K t/s', got %v", data.Input)
	}
	if data.Output != nil {
		t.Errorf("expected nil Output, got %v", *data.Output)
	}
	if data.Total != nil {
		t.Errorf("expected nil Total when output missing, got %v", *data.Total)
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0 t/s"},
		{42, "42 t/s"},
		{999, "999 t/s"},
		{1000, "1.0K t/s"},
		{1500, "1.5K t/s"},
		{2000, "2.0K t/s"},
		{10500, "10.5K t/s"},
	}

	for _, tt := range tests {
		result := FormatSpeed(tt.input)
		if result != tt.expected {
			t.Errorf("FormatSpeed(%g) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/provider/ -run "TestSpeed|TestFormatSpeed" -v`
Expected: FAIL (types not defined)

- [ ] **Step 3: Implement speed provider**

Create `internal/provider/speed.go`:

```go
package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

// SpeedData holds resolved token speed information.
type SpeedData struct {
	Input  *string
	Output *string
	Total  *string
}

type speedProvider struct{}

func (p *speedProvider) Name() string { return "speed" }

func (p *speedProvider) Resolve(session *types.SessionData) (any, error) {
	data := &SpeedData{}

	cw := session.ContextWindow
	cost := session.Cost
	if cw == nil || cost == nil || cost.TotalAPIDurationMS == 0 {
		return data, nil
	}

	durationSec := cost.TotalAPIDurationMS / 1000.0

	if cw.TotalInputTokens != nil {
		speed := float64(*cw.TotalInputTokens) / durationSec
		s := FormatSpeed(speed)
		data.Input = &s
	}

	if cw.TotalOutputTokens != nil {
		speed := float64(*cw.TotalOutputTokens) / durationSec
		s := FormatSpeed(speed)
		data.Output = &s
	}

	if cw.TotalInputTokens != nil && cw.TotalOutputTokens != nil {
		speed := float64(*cw.TotalInputTokens+*cw.TotalOutputTokens) / durationSec
		s := FormatSpeed(speed)
		data.Total = &s
	}

	return data, nil
}

// FormatSpeed formats a tokens-per-second value for display.
func FormatSpeed(tokensPerSec float64) string {
	if tokensPerSec >= 1000 {
		return fmt.Sprintf("%.1fK t/s", tokensPerSec/1000)
	}
	return fmt.Sprintf("%d t/s", int(tokensPerSec))
}
```

- [ ] **Step 4: Register the speed provider**

Add to `RegisterBuiltin()` in `internal/provider/provider.go`:

```go
	registry.Register(&speedProvider{})
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/provider/ -run "TestSpeed|TestFormatSpeed" -v`
Expected: PASS

- [ ] **Step 6: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/provider/speed.go internal/provider/speed_test.go internal/provider/provider.go
git commit -m "feat(speed): add token speed provider with input/output/total"
```

## Chunk 2: Segments and preset

### Task 4: Add new segments

**Files:**
- Modify: `internal/segment/segment.go`
- Modify: `internal/segment/segment_test.go`

- [ ] **Step 1: Write failing tests for new segments**

Add to `internal/segment/segment_test.go`:

```go
func TestContextInputSegment(t *testing.T) {
	seg := &contextInputSegment{}
	if seg.Name() != "context.input" {
		t.Errorf("expected name context.input, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ContextData{Input: "50K"}}
	result := seg.Render(ctx)
	if result == nil || *result != "50K" {
		t.Errorf("expected '50K', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ContextData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil for empty Input, got %v", *result)
	}
}

func TestContextOutputSegment(t *testing.T) {
	seg := &contextOutputSegment{}
	if seg.Name() != "context.output" {
		t.Errorf("expected name context.output, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ContextData{Output: "8K"}}
	result := seg.Render(ctx)
	if result == nil || *result != "8K" {
		t.Errorf("expected '8K', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ContextData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil for empty Output, got %v", *result)
	}
}

func TestSpeedInputSegment(t *testing.T) {
	seg := &speedInputSegment{}
	if seg.Name() != "speed.input" {
		t.Errorf("expected name speed.input, got %s", seg.Name())
	}

	v := "2K t/s"
	ctx := &types.SegmentContext{Provider: &provider.SpeedData{Input: &v}}
	result := seg.Render(ctx)
	if result == nil || *result != "2K t/s" {
		t.Errorf("expected '2K t/s', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.SpeedData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestSpeedOutputSegment(t *testing.T) {
	seg := &speedOutputSegment{}
	if seg.Name() != "speed.output" {
		t.Errorf("expected name speed.output, got %s", seg.Name())
	}

	v := "1K t/s"
	ctx := &types.SegmentContext{Provider: &provider.SpeedData{Output: &v}}
	result := seg.Render(ctx)
	if result == nil || *result != "1K t/s" {
		t.Errorf("expected '1K t/s', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.SpeedData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestSpeedTotalSegment(t *testing.T) {
	seg := &speedTotalSegment{}
	if seg.Name() != "speed.total" {
		t.Errorf("expected name speed.total, got %s", seg.Name())
	}

	v := "3K t/s"
	ctx := &types.SegmentContext{Provider: &provider.SpeedData{Total: &v}}
	result := seg.Render(ctx)
	if result == nil || *result != "3K t/s" {
		t.Errorf("expected '3K t/s', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.SpeedData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/segment/ -run "Test(Context(Input|Output)|Speed(Input|Output|Total))Segment" -v`
Expected: FAIL (types not defined)

- [ ] **Step 3: Implement the five segment types**

Add to `internal/segment/segment.go` after the existing context segments:

```go
type contextInputSegment struct{}

func (s *contextInputSegment) Name() string { return "context.input" }
func (s *contextInputSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ContextData); ok && data != nil && data.Input != "" {
		return &data.Input
	}
	return nil
}

type contextOutputSegment struct{}

func (s *contextOutputSegment) Name() string { return "context.output" }
func (s *contextOutputSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ContextData); ok && data != nil && data.Output != "" {
		return &data.Output
	}
	return nil
}
```

Add after the cost segment section:

```go
// --- Speed ---

type speedInputSegment struct{}

func (s *speedInputSegment) Name() string { return "speed.input" }
func (s *speedInputSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.SpeedData); ok && data != nil {
		return data.Input
	}
	return nil
}

type speedOutputSegment struct{}

func (s *speedOutputSegment) Name() string { return "speed.output" }
func (s *speedOutputSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.SpeedData); ok && data != nil {
		return data.Output
	}
	return nil
}

type speedTotalSegment struct{}

func (s *speedTotalSegment) Name() string { return "speed.total" }
func (s *speedTotalSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.SpeedData); ok && data != nil {
		return data.Total
	}
	return nil
}
```

Register all five in `RegisterBuiltin()`:

```go
	registry.Register(&contextInputSegment{})
	registry.Register(&contextOutputSegment{})
	registry.Register(&speedInputSegment{})
	registry.Register(&speedOutputSegment{})
	registry.Register(&speedTotalSegment{})
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/segment/ -run "Test(Context(Input|Output)|Speed(Input|Output|Total))Segment" -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/segment/segment.go internal/segment/segment_test.go
git commit -m "feat(segment): add context.input, context.output, speed.input, speed.output, speed.total segments"
```

### Task 5: Update F1 preset

**Files:**
- Modify: `internal/preset/f1.json`

- [ ] **Step 1: Add speed.output to F1 preset**

Insert after the `context.percent` entry and before the powerline separator
that transitions to the cost block. The new segment stays within the
context block's `#2D2D2D` background:

```json
    {
      "segment": "speed.output",
      "style": { "color": "#50FA7B", "bgcolor": "#2D2D2D", "prefix": " \uf0e7 " }
    },
```

- [ ] **Step 2: Run preset tests and full suite**

Run: `go vet ./... && go test ./... && go build ./...`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/preset/f1.json
git commit -m "feat(preset): add speed.output to F1 preset"
```

### Task 6: Final verification

- [ ] **Step 1: Run full build and tests**

Run: `go vet ./... && go test ./... && go build ./...`
Expected: all PASS, binary builds

- [ ] **Step 2: Smoke test**

Run:
```bash
echo '{"cwd":"/tmp"}' | go run . --preset f1
```
Expected: renders without error, speed segments collapse

Run:
```bash
echo '{"cwd":"'$(pwd)'","context_window":{"total_input_tokens":50000,"total_output_tokens":8000,"used_percentage":36,"context_window_size":200000},"cost":{"total_api_duration_ms":5000}}' | go run . --preset f1
```
Expected: renders with speed.output visible (e.g., "1.6K t/s")
