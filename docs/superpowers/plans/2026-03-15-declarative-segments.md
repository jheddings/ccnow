# Declarative Segments Refactor Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace 25 per-segment Go structs with data-driven segment resolution using provider struct tags, making segments a config concern.

**Architecture:** Providers declare segment-to-field mappings via `segment` struct tags. A tag index built once at startup maps segment names to field accessors. The render pipeline resolves data segments generically — extract value, apply format, style output. Only `literal` and `newline` remain as registered Go segment types.

**Tech Stack:** Go, `reflect` for struct tag parsing (startup only)

**Spec:** `docs/superpowers/specs/2026-03-15-declarative-segments-design.md`

---

## Chunk 1: Foundation — types, tag index, format

### Task 1: Add FieldProvider interface and Format field

**Files:**
- Modify: `internal/types/types.go`

- [ ] **Step 1: Add FieldProvider interface and Format field to SegmentNode**

Add the `FieldProvider` interface after the existing `DataProvider` interface:

```go
// FieldProvider extends DataProvider with struct tag discovery.
// Providers that implement this expose their data struct type
// for the tag index to map segment names to fields.
type FieldProvider interface {
	DataProvider
	Fields() any // returns a zero-value struct pointer, e.g., &GitData{}
}
```

Add `Format` to `SegmentNode`:

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

- [ ] **Step 2: Run existing tests**

Run: `go vet ./... && go test ./...`
Expected: PASS (additive changes only)

- [ ] **Step 3: Commit**

```bash
git add internal/types/types.go
git commit -m "refactor(types): add FieldProvider interface and Format field"
```

### Task 2: Add segment struct tags to all providers

**Files:**
- Modify: `internal/provider/git.go`
- Modify: `internal/provider/context.go`
- Modify: `internal/provider/model.go`
- Modify: `internal/provider/cost.go`
- Modify: `internal/provider/speed.go`
- Modify: `internal/provider/session.go`
- Modify: `internal/provider/pwd.go`

- [ ] **Step 1: Add segment tags and Fields() to git provider**

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

func (p *gitProvider) Fields() any { return &GitData{} }
```

- [ ] **Step 2: Add segment tags and Fields() to context provider**

```go
type ContextData struct {
	Tokens  string `segment:"context.tokens"`
	Size    string `segment:"context.size"`
	Percent *int   `segment:"context.percent,format:%d%%"`
	Input   string `segment:"context.input"`
	Output  string `segment:"context.output"`
}

func (p *contextProvider) Fields() any { return &ContextData{} }
```

- [ ] **Step 3: Add segment tags and Fields() to remaining providers**

`internal/provider/model.go`:
```go
type ModelData struct {
	Name *string `segment:"model.name"`
}

func (p *modelProvider) Fields() any { return &ModelData{} }
```

`internal/provider/cost.go`:
```go
type CostData struct {
	USD *string `segment:"cost.usd"`
}

func (p *costProvider) Fields() any { return &CostData{} }
```

`internal/provider/speed.go`:
```go
type SpeedData struct {
	Input  *string `segment:"speed.input"`
	Output *string `segment:"speed.output"`
	Total  *string `segment:"speed.total"`
}

func (p *speedProvider) Fields() any { return &SpeedData{} }
```

`internal/provider/session.go`:
```go
type SessionData struct {
	Duration     *string `segment:"session.duration"`
	LinesAdded   *int    `segment:"session.lines-added"`
	LinesRemoved *int    `segment:"session.lines-removed"`
}

func (p *sessionProvider) Fields() any { return &SessionData{} }
```

`internal/provider/pwd.go`:
```go
type PwdData struct {
	Name  string `segment:"pwd.name"`
	Path  string `segment:"pwd.path"`
	Smart string `segment:"pwd.smart"`
}

func (p *pwdProvider) Fields() any { return &PwdData{} }
```

- [ ] **Step 4: Run existing tests**

Run: `go vet ./... && go test ./...`
Expected: PASS (tags and Fields() are additive)

- [ ] **Step 5: Commit**

```bash
git add internal/provider/
git commit -m "refactor(provider): add segment struct tags and Fields() to all providers"
```

### Task 3: Implement tag index

**Files:**
- Create: `internal/render/tagindex.go`
- Create: `internal/render/tagindex_test.go`

- [ ] **Step 1: Write failing tests for BuildTagIndex**

Create `internal/render/tagindex_test.go`:

```go
package render

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

// Test provider and data struct
type testData struct {
	Name  *string `segment:"test.name"`
	Count *int    `segment:"test.count"`
	Plain string  `segment:"test.plain"`
	Skip  string  // no tag
}

type testProvider struct{}

func (p *testProvider) Name() string                                   { return "test" }
func (p *testProvider) Resolve(session *types.SessionData) (any, error) { return &testData{}, nil }
func (p *testProvider) Fields() any                                    { return &testData{} }

// Provider with default format in tag
type fmtData struct {
	Pct *int `segment:"fmt.pct,format:%d%%"`
}

type fmtProvider struct{}

func (p *fmtProvider) Name() string                                   { return "fmt" }
func (p *fmtProvider) Resolve(session *types.SessionData) (any, error) { return &fmtData{}, nil }
func (p *fmtProvider) Fields() any                                    { return &fmtData{} }

// Provider without FieldProvider interface
type plainProvider struct{}

func (p *plainProvider) Name() string                                   { return "plain" }
func (p *plainProvider) Resolve(session *types.SessionData) (any, error) { return nil, nil }

func TestBuildTagIndex(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}

	idx, err := BuildTagIndex(providers)
	if err != nil {
		t.Fatal(err)
	}

	// Tagged fields should be in the index
	if _, ok := idx["test.name"]; !ok {
		t.Error("expected test.name in index")
	}
	if _, ok := idx["test.count"]; !ok {
		t.Error("expected test.count in index")
	}
	if _, ok := idx["test.plain"]; !ok {
		t.Error("expected test.plain in index")
	}

	// Field without tag should not be in the index
	if _, ok := idx["test.skip"]; ok {
		t.Error("expected test.skip NOT in index")
	}

	// Provider name should be set
	if idx["test.name"].Provider != "test" {
		t.Errorf("expected provider 'test', got %q", idx["test.name"].Provider)
	}
}

func TestBuildTagIndex_DefaultFormat(t *testing.T) {
	providers := map[string]types.DataProvider{
		"fmt": &fmtProvider{},
	}

	idx, err := BuildTagIndex(providers)
	if err != nil {
		t.Fatal(err)
	}

	if idx["fmt.pct"].DefaultFormat != "%d%%" {
		t.Errorf("expected default format '%%d%%%%', got %q", idx["fmt.pct"].DefaultFormat)
	}
}

func TestBuildTagIndex_SkipsNonFieldProvider(t *testing.T) {
	providers := map[string]types.DataProvider{
		"plain": &plainProvider{},
	}

	idx, err := BuildTagIndex(providers)
	if err != nil {
		t.Fatal(err)
	}

	if len(idx) != 0 {
		t.Errorf("expected empty index for non-FieldProvider, got %d entries", len(idx))
	}
}

func TestBuildTagIndex_DuplicateErrors(t *testing.T) {
	// Two providers claiming the same segment name
	providers := map[string]types.DataProvider{
		"test":  &testProvider{},
		"test2": &testProvider{}, // same Fields() → same tags
	}

	_, err := BuildTagIndex(providers)
	if err == nil {
		t.Error("expected error for duplicate segment names")
	}
}

func TestResolveSegmentValues(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}

	idx, _ := BuildTagIndex(providers)

	name := "hello"
	count := 42
	providerData := map[string]any{
		"test": &testData{Name: &name, Count: &count, Plain: "raw"},
	}

	values := ResolveSegmentValues(idx, providerData)

	// Pointer fields dereferenced
	if values["test.name"] != "hello" {
		t.Errorf("expected 'hello', got %v", values["test.name"])
	}
	if values["test.count"] != 42 {
		t.Errorf("expected 42, got %v", values["test.count"])
	}

	// Non-pointer field as-is
	if values["test.plain"] != "raw" {
		t.Errorf("expected 'raw', got %v", values["test.plain"])
	}
}

func TestResolveSegmentValues_NilPointer(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}

	idx, _ := BuildTagIndex(providers)

	providerData := map[string]any{
		"test": &testData{}, // all pointer fields nil
	}

	values := ResolveSegmentValues(idx, providerData)

	if values["test.name"] != nil {
		t.Errorf("expected nil for nil *string, got %v", values["test.name"])
	}
	if values["test.count"] != nil {
		t.Errorf("expected nil for nil *int, got %v", values["test.count"])
	}
}

func TestResolveSegmentValues_MissingProvider(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}

	idx, _ := BuildTagIndex(providers)

	// No provider data for "test"
	values := ResolveSegmentValues(idx, map[string]any{})

	if _, ok := values["test.name"]; ok {
		t.Error("expected test.name NOT in values when provider data missing")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/render/ -run "TestBuildTagIndex|TestResolveSegmentValues" -v`
Expected: FAIL (functions don't exist)

- [ ] **Step 3: Implement tag index**

Create `internal/render/tagindex.go`:

```go
package render

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jheddings/ccglow/internal/types"
)

// fieldAccessor maps a segment name to a provider field.
type fieldAccessor struct {
	Provider      string
	FieldIndex    int
	DefaultFormat string
}

// TagIndex maps segment names to their provider field accessors.
type TagIndex map[string]fieldAccessor

// BuildTagIndex reflects on provider data struct tags to build the
// segment name → field mapping. Built once at startup.
func BuildTagIndex(providers map[string]types.DataProvider) (TagIndex, error) {
	idx := make(TagIndex)

	for name, p := range providers {
		fp, ok := p.(types.FieldProvider)
		if !ok {
			continue
		}

		fields := fp.Fields()
		t := reflect.TypeOf(fields)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			continue
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("segment")
			if tag == "" {
				continue
			}

			segName, defaultFmt := parseSegmentTag(tag)

			if _, exists := idx[segName]; exists {
				return nil, fmt.Errorf("duplicate segment name %q", segName)
			}

			idx[segName] = fieldAccessor{
				Provider:      name,
				FieldIndex:    i,
				DefaultFormat: defaultFmt,
			}
		}
	}

	return idx, nil
}

// parseSegmentTag parses "name" or "name,format:fmt" from a struct tag.
func parseSegmentTag(tag string) (name, defaultFormat string) {
	parts := strings.SplitN(tag, ",", 2)
	name = parts[0]
	if len(parts) > 1 && strings.HasPrefix(parts[1], "format:") {
		defaultFormat = strings.TrimPrefix(parts[1], "format:")
	}
	return
}

// ResolveSegmentValues builds a flat map of segment name → value from
// resolved provider data. Pointer fields are dereferenced; nil stays nil.
func ResolveSegmentValues(idx TagIndex, providerData map[string]any) map[string]any {
	values := make(map[string]any)

	for segName, accessor := range idx {
		data, ok := providerData[accessor.Provider]
		if !ok || data == nil {
			continue
		}

		v := reflect.ValueOf(data)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				continue
			}
			v = v.Elem()
		}

		field := v.Field(accessor.FieldIndex)

		// Dereference pointer fields
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				values[segName] = nil
				continue
			}
			field = field.Elem()
		}

		values[segName] = field.Interface()
	}

	return values
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/render/ -run "TestBuildTagIndex|TestResolveSegmentValues" -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/render/tagindex.go internal/render/tagindex_test.go
git commit -m "refactor(render): add tag index for segment-to-field mapping"
```

### Task 4: Implement FormatValue

**Files:**
- Create: `internal/render/format.go`
- Create: `internal/render/format_test.go`

- [ ] **Step 1: Write failing tests**

Create `internal/render/format_test.go`:

```go
package render

import "testing"

func TestFormatValue(t *testing.T) {
	tests := []struct {
		value    any
		format   string
		expected string
	}{
		// No format — %v passthrough
		{42, "", "42"},
		{"hello", "", "hello"},
		{3.14, "", "3.14"},

		// With format string
		{42, "%d%%", "42%"},
		{42, "+%d", "+42"},
		{"text", "%s!", "text!"},
		{3.14, "%.1f", "3.1"},

		// Nil — defensive
		{nil, "", ""},
		{nil, "%v", ""},
	}

	for _, tt := range tests {
		result := FormatValue(tt.value, tt.format)
		if result != tt.expected {
			t.Errorf("FormatValue(%v, %q) = %q, want %q", tt.value, tt.format, result, tt.expected)
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/render/ -run TestFormatValue -v`
Expected: FAIL

- [ ] **Step 3: Implement FormatValue**

Create `internal/render/format.go`:

```go
package render

import "fmt"

// FormatValue formats a segment value for display.
// If format is empty, uses fmt.Sprintf("%v", value).
// Nil values return "".
func FormatValue(value any, format string) string {
	if value == nil {
		return ""
	}
	if format == "" {
		return fmt.Sprintf("%v", value)
	}
	return fmt.Sprintf(format, value)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/render/ -run TestFormatValue -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/render/format.go internal/render/format_test.go
git commit -m "refactor(render): add FormatValue for data segment formatting"
```

## Chunk 2: Render pipeline and cleanup

### Task 5: Update render pipeline for data segments

**Files:**
- Modify: `internal/render/render.go`
- Modify: `internal/render/render_test.go`

- [ ] **Step 1: Write new render tests for data segment behavior**

Add to `internal/render/render_test.go`:

```go
func TestTree_DataSegment(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"test.name": "hello",
	}
	tagIdx := TagIndex{
		"test.name": fieldAccessor{Provider: "test", FieldIndex: 0},
	}

	tree := []types.SegmentNode{
		{Type: "test.name"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_DataSegmentNilCollapses(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"test.name": nil,
	}
	tagIdx := TagIndex{
		"test.name": fieldAccessor{Provider: "test", FieldIndex: 0},
	}

	tree := []types.SegmentNode{
		{Type: "test.name"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty (collapsed), got %q", result)
	}
}

func TestTree_DataSegmentWithFormat(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"test.count": 42,
	}
	tagIdx := TagIndex{
		"test.count": fieldAccessor{Provider: "test", FieldIndex: 1},
	}

	tree := []types.SegmentNode{
		{Type: "test.count", Format: "+%d"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "+42" {
		t.Errorf("expected '+42', got %q", result)
	}
}

func TestTree_DataSegmentDefaultFormat(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"fmt.pct": 85,
	}
	tagIdx := TagIndex{
		"fmt.pct": fieldAccessor{Provider: "fmt", FieldIndex: 0, DefaultFormat: "%d%%"},
	}

	tree := []types.SegmentNode{
		{Type: "fmt.pct"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "85%" {
		t.Errorf("expected '85%%', got %q", result)
	}
}

func TestTree_DataSegmentFormatOverridesDefault(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"fmt.pct": 85,
	}
	tagIdx := TagIndex{
		"fmt.pct": fieldAccessor{Provider: "fmt", FieldIndex: 0, DefaultFormat: "%d%%"},
	}

	tree := []types.SegmentNode{
		{Type: "fmt.pct", Format: "(%d)"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "(85)" {
		t.Errorf("expected '(85)', got %q", result)
	}
}

func TestTree_EmptyStringCollapses(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"test.plain": "",
	}
	tagIdx := TagIndex{
		"test.plain": fieldAccessor{Provider: "test", FieldIndex: 2},
	}

	tree := []types.SegmentNode{
		{Type: "test.plain"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty (collapsed), got %q", result)
	}
}
```

- [ ] **Step 2: Run new tests to verify they fail**

Run: `go test ./internal/render/ -run "TestTree_DataSegment" -v`
Expected: FAIL (Tree signature doesn't match)

- [ ] **Step 3: Update renderNode and Tree to handle data segments**

Update `internal/render/render.go`. The `Tree` and `renderNode` functions
gain `segmentValues map[string]any` and `tagIdx TagIndex` parameters:

```go
func renderNode(
	node *types.SegmentNode,
	segments *segment.Registry,
	session *types.SessionData,
	providerData map[string]any,
	segmentValues map[string]any,
	tagIdx TagIndex,
) *string {
	if !isEnabled(node, session) {
		return nil
	}

	// SegmentGroup: render children, collapse if all nil
	if len(node.Children) > 0 {
		var parts []string
		for i := range node.Children {
			rendered := renderNode(&node.Children[i], segments, session, providerData, segmentValues, tagIdx)
			if rendered != nil {
				parts = append(parts, *rendered)
			}
		}
		if len(parts) == 0 {
			return nil
		}
		joined := strings.Join(parts, "")
		styled := style.Apply(joined, node.Style)
		return &styled
	}

	// LiteralSegment: delegate to registered segment
	seg := segments.Get(node.Type)
	if seg != nil {
		ctx := &types.SegmentContext{
			Session: session,
			Props:   node.Props,
		}
		value := seg.Render(ctx)
		if value == nil {
			return nil
		}
		styled := style.Apply(*value, node.Style)
		return &styled
	}

	// DataSegment: resolve from segment values
	value, ok := segmentValues[node.Type]
	if !ok || value == nil {
		return nil
	}

	// Resolve format: config override > tag default > none
	format := node.Format
	if format == "" {
		if accessor, exists := tagIdx[node.Type]; exists {
			format = accessor.DefaultFormat
		}
	}

	text := FormatValue(value, format)
	if text == "" {
		return nil
	}

	styled := style.Apply(text, node.Style)
	return &styled
}

func Tree(
	tree []types.SegmentNode,
	segments *segment.Registry,
	session *types.SessionData,
	providerData map[string]any,
	segmentValues map[string]any,
	tagIdx TagIndex,
) string {
	var parts []string
	for i := range tree {
		rendered := renderNode(&tree[i], segments, session, providerData, segmentValues, tagIdx)
		if rendered != nil {
			parts = append(parts, *rendered)
		}
	}
	return strings.Join(parts, "")
}
```

- [ ] **Step 4: Update existing render tests for new signature**

Update all existing tests in `render_test.go` to pass the new parameters.
For existing tests, pass empty `segmentValues` and `tagIdx`:

```go
// In each existing test, change Tree calls from:
result := Tree(tree, seg, sess, providerData)
// To:
result := Tree(tree, seg, sess, providerData, map[string]any{}, TagIndex{})
```

The existing tests for `pwd.name`, `git.branch`, etc. still work because
those segment types are still in the registry (we haven't removed them yet).

- [ ] **Step 5: Run all render tests**

Run: `go test ./internal/render/ -v`
Expected: PASS (both old and new tests)

- [ ] **Step 6: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: FAIL — `main.go` calls `render.Tree` with old signature. Fix in next step.

- [ ] **Step 7: Update main.go to pass new parameters**

Update `run()` in `main.go`:

```go
func run(presetName, configPath, format, stdin string) string {
	sess := session.Parse(stdin)
	if sess == nil {
		return ""
	}

	if format == "plain" {
		style.SetColorLevel(0)
	} else {
		style.SetColorLevel(1)
	}
	defer style.SetColorLevel(1)

	segments := segment.NewRegistry()
	segment.RegisterBuiltin(segments)

	providers := provider.NewRegistry()
	provider.RegisterBuiltin(providers)

	tagIdx, err := render.BuildTagIndex(providers.All())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ccglow: tag index error: %v\n", err)
		return ""
	}

	tree := resolveTree(presetName, configPath)

	providerNames := render.CollectProviderNames(tree)
	providerData := render.ResolveProviders(providerNames, providers.All(), sess)
	segmentValues := render.ResolveSegmentValues(tagIdx, providerData)

	return render.Tree(tree, segments, sess, providerData, segmentValues, tagIdx)
}
```

- [ ] **Step 8: Update CollectProviderNames to use tag index**

Update `CollectProviderNames` to require a `TagIndex` parameter and use it
to look up provider names for data segment types. Fall back to the
existing `node.Provider` field for composite groups:

```go
func CollectProviderNames(tree []types.SegmentNode, tagIdx TagIndex) map[string]bool {
	names := make(map[string]bool)
	collectNames(tree, names, tagIdx)
	return names
}

func collectNames(nodes []types.SegmentNode, names map[string]bool, idx TagIndex) {
	for _, node := range nodes {
		if node.Enabled != nil && !*node.Enabled {
			continue
		}
		// Try tag index first (data segments)
		if accessor, ok := idx[node.Type]; ok {
			names[accessor.Provider] = true
		}
		// Fall back to explicit Provider field (composite groups)
		if node.Provider != "" {
			names[node.Provider] = true
		}
		if len(node.Children) > 0 {
			collectNames(node.Children, names, idx)
		}
	}
}
```

Add a test for `CollectProviderNames` with the tag index:

```go
func TestCollectProviderNames_TagIndex(t *testing.T) {
	tagIdx := TagIndex{
		"git.branch":     fieldAccessor{Provider: "git", FieldIndex: 0},
		"context.tokens": fieldAccessor{Provider: "context", FieldIndex: 0},
	}

	tree := []types.SegmentNode{
		{Type: "git.branch"},
		{Type: "context.tokens"},
		{Type: "literal", Props: map[string]any{"text": "hi"}},
	}

	names := CollectProviderNames(tree, tagIdx)
	if !names["git"] {
		t.Error("expected git provider")
	}
	if !names["context"] {
		t.Error("expected context provider")
	}
	if len(names) != 2 {
		t.Errorf("expected 2 providers, got %d", len(names))
	}
}
```

Update `main.go` call:

```go
providerNames := render.CollectProviderNames(tree, tagIdx)
```

- [ ] **Step 9: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 10: Commit**

```bash
git add internal/render/render.go internal/render/render_test.go main.go
git commit -m "refactor(render): add data segment resolution to render pipeline"
```

### Task 6: Remove data segment structs and clean up

**Files:**
- Modify: `internal/segment/segment.go`
- Modify: `internal/segment/segment_test.go`
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`
- Modify: `internal/types/types.go`

- [ ] **Step 1: Remove all data segment structs from segment.go**

Replace the entire contents of `internal/segment/segment.go` with only
the registry, literal, and newline:

```go
package segment

import "github.com/jheddings/ccglow/internal/types"

// RegisterBuiltin adds all built-in segment implementations to the registry.
func RegisterBuiltin(registry *Registry) {
	registry.Register(&literalSegment{})
	registry.Register(&newlineSegment{})
}

// Registry maps segment type names to their implementations.
type Registry struct {
	segments map[string]types.Segment
}

// NewRegistry creates an empty segment registry.
func NewRegistry() *Registry {
	return &Registry{segments: make(map[string]types.Segment)}
}

// Register adds a segment implementation.
func (r *Registry) Register(seg types.Segment) {
	r.segments[seg.Name()] = seg
}

// Get returns the segment for the given type name, or nil.
func (r *Registry) Get(name string) types.Segment {
	return r.segments[name]
}

// --- Literal ---

type literalSegment struct{}

func (s *literalSegment) Name() string { return "literal" }
func (s *literalSegment) Render(ctx *types.SegmentContext) *string {
	if ctx.Props == nil {
		return nil
	}
	if text, ok := ctx.Props["text"].(string); ok {
		return &text
	}
	return nil
}

// --- Newline ---

type newlineSegment struct{}

func (s *newlineSegment) Name() string { return "newline" }
func (s *newlineSegment) Render(ctx *types.SegmentContext) *string {
	v := "\n"
	return &v
}
```

- [ ] **Step 2: Replace segment_test.go with literal/newline tests only**

Replace `internal/segment/segment_test.go`:

```go
package segment

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestLiteralSegment(t *testing.T) {
	seg := &literalSegment{}
	if seg.Name() != "literal" {
		t.Errorf("expected name literal, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Props: map[string]any{"text": "hello"}}
	result := seg.Render(ctx)
	if result == nil || *result != "hello" {
		t.Errorf("expected 'hello', got %v", result)
	}

	// No props
	ctx = &types.SegmentContext{}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestNewlineSegment(t *testing.T) {
	seg := &newlineSegment{}
	if seg.Name() != "newline" {
		t.Errorf("expected name newline, got %s", seg.Name())
	}

	result := seg.Render(&types.SegmentContext{})
	if result == nil || *result != "\n" {
		t.Errorf("expected newline, got %v", result)
	}
}
```

- [ ] **Step 3: Remove InferProviders from config.go**

Replace `internal/config/config.go`:

```go
package config

import (
	"encoding/json"

	"github.com/jheddings/ccglow/internal/types"
)

type configFile struct {
	Segments []json.RawMessage `json:"segments"`
}

// Parse parses a JSON config file into a segment tree.
func Parse(data []byte) ([]types.SegmentNode, error) {
	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	var nodes []types.SegmentNode
	for _, raw := range cfg.Segments {
		var node types.SegmentNode
		if err := json.Unmarshal(raw, &node); err != nil {
			continue
		}
		if node.Type == "" {
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}
```

- [ ] **Step 4: Remove SegmentContext.Provider field**

In `internal/types/types.go`, update `SegmentContext`:

```go
// SegmentContext is passed to Segment.Render with resolved data.
type SegmentContext struct {
	Session *SessionData
	Props   map[string]any
}
```

- [ ] **Step 5: Update existing render tests that use old data segments**

Update `TestTree_AtomicNode` to use the data segment path instead of
the registered segment path. Since `pwd.name` is no longer in the
registry, it goes through the data segment branch:

```go
func TestTree_AtomicNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/home/user/project"}

	segmentValues := map[string]any{
		"pwd.name": "project",
	}
	tagIdx := TagIndex{
		"pwd.name": fieldAccessor{Provider: "pwd", FieldIndex: 0},
	}

	tree := []types.SegmentNode{
		{Type: "pwd.name"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "project" {
		t.Errorf("expected project, got %q", result)
	}
}
```

Update `TestTree_CompositeCollapse`:

```go
func TestTree_CompositeCollapse(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"git.branch":     nil,
		"git.insertions": nil,
	}
	tagIdx := TagIndex{
		"git.branch":     fieldAccessor{Provider: "git", FieldIndex: 0},
		"git.insertions": fieldAccessor{Provider: "git", FieldIndex: 1},
	}

	tree := []types.SegmentNode{
		{
			Style: &types.StyleAttrs{Prefix: " | "},
			Children: []types.SegmentNode{
				{Type: "git.branch"},
				{Type: "git.insertions"},
			},
		},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty (collapsed composite), got %q", result)
	}
}
```

Update `TestTree_DisabledNode`:

```go
func TestTree_DisabledNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"pwd.name": "tmp",
	}
	tagIdx := TagIndex{
		"pwd.name": fieldAccessor{Provider: "pwd", FieldIndex: 0},
	}

	disabled := false
	tree := []types.SegmentNode{
		{Type: "pwd.name", Enabled: &disabled},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty for disabled node, got %q", result)
	}
}
```

`TestTree_Literal` and `TestTree_MissingSegment` need only the signature
update (add empty `segmentValues` and `tagIdx`).

Remove the `"github.com/jheddings/ccglow/internal/provider"` import from
`render_test.go` — it is no longer used after the test updates.

- [ ] **Step 6: Update config_test.go**

Remove all `Provider` assertions from config tests since `InferProviders`
is gone. Replace `internal/config/config_test.go`:

```go
package config

import "testing"

func TestParse_Valid(t *testing.T) {
	input := `{
		"segments": [
			{"segment": "pwd.name", "style": {"color": "red"}},
			{"segment": "git.branch"}
		]
	}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Type != "pwd.name" {
		t.Errorf("expected pwd.name, got %s", nodes[0].Type)
	}
}

func TestParse_WithChildren(t *testing.T) {
	input := `{
		"segments": [
			{
				"segment": "group",
				"children": [
					{"segment": "git.branch"},
					{"segment": "git.insertions"}
				]
			}
		]
	}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	if len(nodes[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(nodes[0].Children))
	}
}

func TestParse_LiteralNoProvider(t *testing.T) {
	input := `{"segments": [{"segment": "literal", "props": {"text": "hi"}}]}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParse_WithFormat(t *testing.T) {
	input := `{"segments": [{"segment": "context.percent", "format": "%d%%"}]}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if nodes[0].Format != "%d%%" {
		t.Errorf("expected format '%%d%%%%', got %q", nodes[0].Format)
	}
}
```

- [ ] **Step 7: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/segment/ internal/config/ internal/types/types.go internal/render/render_test.go
git commit -m "refactor(segment): remove data segment structs, use tag-driven resolution"
```

### Task 7: Final verification

- [ ] **Step 1: Run full build and tests**

Run: `go vet ./... && go test ./... && go build ./...`
Expected: all PASS, binary builds

- [ ] **Step 2: Smoke test with presets**

Run:
```bash
echo '{"cwd":"'$(pwd)'","context_window":{"used_percentage":36,"context_window_size":200000,"total_input_tokens":50000,"total_output_tokens":8000,"current_usage":{"input_tokens":100}},"cost":{"total_cost_usd":0.42,"total_duration_ms":300000,"total_api_duration_ms":5000},"model":{"display_name":"Opus 4.6"}}' | go run . --preset default
```
Expected: renders with all segments visible

```bash
echo '{"cwd":"'$(pwd)'"}' | go run . --preset f1
```
Expected: renders with git data, non-git segments collapse

```bash
echo '{"cwd":"/tmp"}' | go run . --preset minimal
```
Expected: renders without git data

- [ ] **Step 3: Verify format override works**

Create a test config `/tmp/test-format.json`:
```json
{
  "segments": [
    { "segment": "context.percent", "format": "(%d)" }
  ]
}
```

Run:
```bash
echo '{"cwd":"/tmp","context_window":{"used_percentage":42,"current_usage":{"input_tokens":100}}}' | go run . --config /tmp/test-format.json --format plain
```
Expected: `(42)` — config format overrides tag default `%d%%`
