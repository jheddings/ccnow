# New Payload Segments Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `claude.version`, `claude.style`, `model.id`, and `context.remaining` segments sourced from the existing Claude Code statusline JSON payload.

**Architecture:** Extend `SessionData` types to capture `version`, `output_style`, and `remaining_percentage` fields. Add a new `claude` provider for app-level metadata. Extend existing `model` and `context` providers with new fields. Register four new segments.

**Tech Stack:** Go, standard library only

---

## File Structure

| File | Action | Responsibility |
|---|---|---|
| `internal/types/types.go` | Modify | Add `Version`, `OutputStyle`, `RemainingPercentage` fields |
| `internal/provider/claude.go` | Create | New `claude` provider returning `ClaudeData` |
| `internal/provider/model.go` | Modify | Add `ID` field to `ModelData` |
| `internal/provider/context.go` | Modify | Add `Remaining` field to `ContextData` |
| `internal/provider/provider.go` | Modify | Register `claude` provider |
| `internal/provider/provider_test.go` | Modify | Tests for `claude`, extended `model`, extended `context` |
| `internal/segment/segment.go` | Modify | Add 4 segment structs, register them |
| `internal/segment/segment_test.go` | Modify | Tests for all 4 new segments |
| `internal/preset/full.json` | Modify | Add new segments to full preset |

---

## Chunk 1: Types and Providers

### Task 1: Add new fields to SessionData types

**Files:**
- Modify: `internal/types/types.go:4-9` (SessionData), `internal/types/types.go:27-33` (ContextWindow)

- [ ] **Step 1: Add OutputStyleInfo struct and extend SessionData**

In `internal/types/types.go`, add the `OutputStyleInfo` struct after `ModelInfo`, and add `Version` and `OutputStyle` fields to `SessionData`:

```go
// SessionData represents the JSON session data piped from Claude Code.
type SessionData struct {
	CWD           string         `json:"cwd"`
	Model         *ModelInfo     `json:"model,omitempty"`
	Cost          *CostInfo      `json:"cost,omitempty"`
	ContextWindow *ContextWindow `json:"context_window,omitempty"`
	Version       string         `json:"version,omitempty"`
	OutputStyle   *OutputStyleInfo `json:"output_style,omitempty"`
}

// ModelInfo contains model identification from the session.
type ModelInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// OutputStyleInfo contains the output style configuration.
type OutputStyleInfo struct {
	Name string `json:"name"`
}
```

- [ ] **Step 2: Add RemainingPercentage to ContextWindow**

In the same file, add `RemainingPercentage` to the `ContextWindow` struct:

```go
// ContextWindow contains token usage data.
type ContextWindow struct {
	UsedPercentage      int           `json:"used_percentage"`
	RemainingPercentage int           `json:"remaining_percentage"`
	ContextWindowSize   int           `json:"context_window_size,omitempty"`
	TotalInputTokens    *int          `json:"total_input_tokens,omitempty"`
	TotalOutputTokens   *int          `json:"total_output_tokens,omitempty"`
	CurrentUsage        *CurrentUsage `json:"current_usage,omitempty"`
}
```

- [ ] **Step 3: Run tests to verify no regressions**

Run: `go test ./...`
Expected: All existing tests pass (type additions are backward-compatible).

- [ ] **Step 4: Commit**

```bash
git add internal/types/types.go
git commit -m "feat(types): add Version, OutputStyle, and RemainingPercentage fields"
```

---

### Task 2: Create claude provider with tests (TDD)

**Files:**
- Create: `internal/provider/claude.go`
- Modify: `internal/provider/provider.go:6-14`
- Modify: `internal/provider/provider_test.go`

- [ ] **Step 1: Write failing tests for claude provider**

Add to `internal/provider/provider_test.go`:

```go
func TestClaudeProvider(t *testing.T) {
	p := &claudeProvider{}
	sess := &types.SessionData{
		CWD:     "/tmp",
		Version: "2.1.75",
		OutputStyle: &types.OutputStyleInfo{Name: "concise"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ClaudeData)
	if data.Version == nil || *data.Version != "2.1.75" {
		t.Errorf("expected version 2.1.75, got %v", data.Version)
	}
	if data.Style == nil || *data.Style != "concise" {
		t.Errorf("expected style concise, got %v", data.Style)
	}
}

func TestClaudeProviderEmpty(t *testing.T) {
	p := &claudeProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ClaudeData)
	if data.Version != nil {
		t.Errorf("expected nil version, got %v", data.Version)
	}
	if data.Style != nil {
		t.Errorf("expected nil style, got %v", data.Style)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/provider/...`
Expected: FAIL — `claudeProvider` and `ClaudeData` undefined.

- [ ] **Step 3: Implement claude provider**

Create `internal/provider/claude.go`:

```go
package provider

import "github.com/jheddings/ccglow/internal/types"

// ClaudeData holds resolved Claude Code application metadata.
type ClaudeData struct {
	Version *string
	Style   *string
}

type claudeProvider struct{}

func (p *claudeProvider) Name() string { return "claude" }

func (p *claudeProvider) Resolve(session *types.SessionData) (any, error) {
	data := &ClaudeData{}
	if session.Version != "" {
		data.Version = &session.Version
	}
	if session.OutputStyle != nil && session.OutputStyle.Name != "" {
		data.Style = &session.OutputStyle.Name
	}
	return data, nil
}
```

- [ ] **Step 4: Register claude provider**

In `internal/provider/provider.go`, add to `RegisterBuiltin`:

```go
registry.Register(&claudeProvider{})
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/provider/...`
Expected: All tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/provider/claude.go internal/provider/provider.go internal/provider/provider_test.go
git commit -m "feat(provider): add claude provider for version and output style"
```

---

### Task 3: Extend model provider with ID (TDD)

**Files:**
- Modify: `internal/provider/model.go`
- Modify: `internal/provider/provider_test.go`

- [ ] **Step 1: Write failing test for model ID**

Add to `internal/provider/provider_test.go`:

```go
func TestModelProviderID(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{
		CWD:   "/tmp",
		Model: &types.ModelInfo{ID: "claude-opus-4-6[1m]", DisplayName: "Opus 4.6"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ModelData)
	if data.ID == nil || *data.ID != "claude-opus-4-6[1m]" {
		t.Errorf("expected claude-opus-4-6[1m], got %v", data.ID)
	}
}

func TestModelProviderNoModel(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ModelData)
	if data.ID != nil {
		t.Errorf("expected nil ID, got %v", data.ID)
	}
	if data.Name != nil {
		t.Errorf("expected nil Name, got %v", data.Name)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/provider/...`
Expected: FAIL — `ModelData` has no field `ID`.

- [ ] **Step 3: Add ID to ModelData and populate in Resolve**

Modify `internal/provider/model.go`:

```go
package provider

import "github.com/jheddings/ccglow/internal/types"

// ModelData holds resolved model information.
type ModelData struct {
	Name *string
	ID   *string
}

type modelProvider struct{}

func (p *modelProvider) Name() string { return "model" }

func (p *modelProvider) Resolve(session *types.SessionData) (any, error) {
	data := &ModelData{}
	if session.Model != nil && session.Model.DisplayName != "" {
		data.Name = &session.Model.DisplayName
	}
	if session.Model != nil && session.Model.ID != "" {
		data.ID = &session.Model.ID
	}
	return data, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/provider/...`
Expected: All tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/provider/model.go internal/provider/provider_test.go
git commit -m "feat(provider): add model ID to model provider"
```

---

### Task 4: Extend context provider with Remaining (TDD)

**Files:**
- Modify: `internal/provider/context.go`
- Modify: `internal/provider/provider_test.go`

- [ ] **Step 1: Write failing test for context remaining**

Add to `internal/provider/provider_test.go`:

```go
func TestContextProviderRemaining(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage:      36,
			RemainingPercentage: 64,
			ContextWindowSize:   1000000,
			CurrentUsage:        &types.CurrentUsage{InputTokens: 100},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	if data.Remaining == nil || *data.Remaining != 64 {
		t.Errorf("expected remaining 64, got %v", data.Remaining)
	}
}

func TestContextProviderNoRemaining(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	if data.Remaining != nil {
		t.Errorf("expected nil remaining, got %v", data.Remaining)
	}
}

func TestContextProviderZeroRemaining(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD:           "/tmp",
		ContextWindow: &types.ContextWindow{},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	if data.Remaining != nil {
		t.Errorf("expected nil remaining for zero value with no usage, got %v", data.Remaining)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/provider/...`
Expected: FAIL — `ContextData` has no field `Remaining`.

- [ ] **Step 3: Add Remaining to ContextData and populate in Resolve**

In `internal/provider/context.go`, add `Remaining *int` to `ContextData`:

```go
type ContextData struct {
	Tokens    string
	Size      string
	Percent   *int
	Remaining *int
	Input     string
	Output    string
}
```

At the end of the `Resolve` method (after the existing `Percent` block), add:

```go
if cw.RemainingPercentage > 0 || cw.CurrentUsage != nil {
	rem := cw.RemainingPercentage
	data.Remaining = &rem
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/provider/...`
Expected: All tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/provider/context.go internal/provider/provider_test.go
git commit -m "feat(provider): add remaining percentage to context provider"
```

---

## Chunk 2: Segments and Preset

### Task 5: Add four new segments with tests (TDD)

**Files:**
- Modify: `internal/segment/segment.go`
- Modify: `internal/segment/segment_test.go`

- [ ] **Step 1: Write failing tests for all four segments**

Add to `internal/segment/segment_test.go`:

```go
// Add strPtr alongside the existing intPtr helper at the top of the file.
func strPtr(s string) *string { return &s }

func TestClaudeVersionSegment(t *testing.T) {
	seg := &claudeVersionSegment{}
	if seg.Name() != "claude.version" {
		t.Errorf("expected name claude.version, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ClaudeData{Version: strPtr("2.1.75")}}
	result := seg.Render(ctx)
	if result == nil || *result != "2.1.75" {
		t.Errorf("expected '2.1.75', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ClaudeData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestClaudeStyleSegment(t *testing.T) {
	seg := &claudeStyleSegment{}
	if seg.Name() != "claude.style" {
		t.Errorf("expected name claude.style, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ClaudeData{Style: strPtr("concise")}}
	result := seg.Render(ctx)
	if result == nil || *result != "concise" {
		t.Errorf("expected 'concise', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ClaudeData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestModelIDSegment(t *testing.T) {
	seg := &modelIDSegment{}
	if seg.Name() != "model.id" {
		t.Errorf("expected name model.id, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ModelData{ID: strPtr("claude-opus-4-6[1m]")}}
	result := seg.Render(ctx)
	if result == nil || *result != "claude-opus-4-6[1m]" {
		t.Errorf("expected 'claude-opus-4-6[1m]', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ModelData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestContextRemainingSegment(t *testing.T) {
	seg := &contextRemainingSegment{}
	if seg.Name() != "context.remaining" {
		t.Errorf("expected name context.remaining, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ContextData{Remaining: intPtr(96)}}
	result := seg.Render(ctx)
	if result == nil || *result != "96%" {
		t.Errorf("expected '96%%', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ContextData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/segment/...`
Expected: FAIL — undefined segment types.

- [ ] **Step 3: Implement all four segments**

Add to `internal/segment/segment.go`, after the existing `modelNameSegment`:

```go
type modelIDSegment struct{}

func (s *modelIDSegment) Name() string { return "model.id" }
func (s *modelIDSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ModelData); ok && data != nil {
		return data.ID
	}
	return nil
}
```

After the existing `contextOutputSegment`:

```go
type contextRemainingSegment struct{}

func (s *contextRemainingSegment) Name() string { return "context.remaining" }
func (s *contextRemainingSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ContextData); ok && data != nil && data.Remaining != nil {
		v := fmt.Sprintf("%d%%", *data.Remaining)
		return &v
	}
	return nil
}
```

At the end of the file, add a new `// --- Claude ---` section:

```go
// --- Claude ---

type claudeVersionSegment struct{}

func (s *claudeVersionSegment) Name() string { return "claude.version" }
func (s *claudeVersionSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ClaudeData); ok && data != nil {
		return data.Version
	}
	return nil
}

type claudeStyleSegment struct{}

func (s *claudeStyleSegment) Name() string { return "claude.style" }
func (s *claudeStyleSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ClaudeData); ok && data != nil {
		return data.Style
	}
	return nil
}
```

- [ ] **Step 4: Register all four segments in RegisterBuiltin**

In `RegisterBuiltin`, add after `&modelNameSegment{}`:

```go
registry.Register(&modelIDSegment{})
```

Add after `&contextOutputSegment{}`:

```go
registry.Register(&contextRemainingSegment{})
```

Add at the end:

```go
registry.Register(&claudeVersionSegment{})
registry.Register(&claudeStyleSegment{})
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/segment/...`
Expected: All tests PASS.

- [ ] **Step 6: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: All pass, no vet warnings.

- [ ] **Step 7: Commit**

```bash
git add internal/segment/segment.go internal/segment/segment_test.go
git commit -m "feat(segment): add claude.version, claude.style, model.id, and context.remaining segments"
```

---

### Task 6: Update full preset

**Files:**
- Modify: `internal/preset/full.json`

- [ ] **Step 1: Add model.id to the model group**

In `internal/preset/full.json`, find the group containing `model.name` (currently lines 62-66). Replace it with a group that includes both `model.name` and `model.id`:

```json
{
  "segment": "group",
  "style": {
    "prefix": " | "
  },
  "children": [
    {
      "segment": "model.name"
    },
    {
      "segment": "model.id",
      "style": {
        "color": "240",
        "prefix": " "
      }
    }
  ]
}
```

- [ ] **Step 2: Add context.remaining near context.percent**

In the context group, add `context.remaining` after `context.percent`:

```json
{
  "segment": "context.remaining",
  "style": {
    "color": "white",
    "prefix": " [",
    "suffix": " free]"
  }
}
```

- [ ] **Step 3: Add claude group at the end (before session groups)**

Add a new group before the session duration group:

```json
{
  "segment": "group",
  "style": {
    "prefix": " · "
  },
  "children": [
    {
      "segment": "claude.version",
      "style": {
        "color": "240",
        "prefix": "v"
      }
    },
    {
      "segment": "claude.style",
      "style": {
        "color": "cyan",
        "prefix": " · "
      }
    }
  ]
}
```

- [ ] **Step 4: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: All pass.

- [ ] **Step 5: Commit**

```bash
git add internal/preset/full.json
git commit -m "feat(preset): add new segments to full preset"
```
