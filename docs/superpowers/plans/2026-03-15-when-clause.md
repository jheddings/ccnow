# When Clause Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `when` conditional expressions to segments using `expr-lang/expr`, enabling data-driven visibility. Update the minimal preset as a showcase.

**Architecture:** New `internal/condition/` package handles expression compilation and evaluation. `BuildEnv` builds an expr environment from provider data, `value`, and `text`. The render pipeline evaluates `when` after formatting but before styling. Compilation is cached per expression string.

**Tech Stack:** Go, `github.com/expr-lang/expr`

**Spec:** `docs/superpowers/specs/2026-03-15-when-clause-design.md`

---

## Chunk 1: Condition package and dependency

### Task 1: Add expr dependency and When field

**Files:**
- Modify: `go.mod`
- Modify: `internal/types/types.go`

- [ ] **Step 1: Add expr-lang/expr dependency**

Run: `go get github.com/expr-lang/expr`

- [ ] **Step 2: Add When field to SegmentNode**

In `internal/types/types.go`, add `When` after `Format`:

```go
type SegmentNode struct {
	Type     string         `json:"segment,omitempty"`
	Provider string         `json:"provider,omitempty"`
	Format   string         `json:"format,omitempty"`
	When     string         `json:"when,omitempty"`
	Enabled  *bool          `json:"enabled,omitempty"`
	Style    *StyleAttrs    `json:"style,omitempty"`
	Props    map[string]any `json:"props,omitempty"`
	Children []SegmentNode  `json:"children,omitempty"`

	// EnabledFn is set programmatically (presets) and takes precedence over Enabled.
	EnabledFn func(*SessionData) bool `json:"-"`
}
```

- [ ] **Step 3: Run existing tests**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum internal/types/types.go
git commit -m "feat(types): add When field and expr-lang/expr dependency"
```

### Task 2: Implement condition package

**Files:**
- Create: `internal/condition/condition.go`
- Create: `internal/condition/condition_test.go`

- [ ] **Step 1: Write failing tests**

Create `internal/condition/condition_test.go`:

```go
package condition

import (
	"testing"
)

// Test struct matching provider data shape
type testProvider struct {
	Name    *string
	Count   *int
	Score   *float64
	Label   string
	Enabled bool
}

func strPtr(s string) *string { return &s }
func intPtr(n int) *int       { return &n }
func floatPtr(f float64) *float64 { return &f }

// --- Compile tests ---

func TestCompile_Valid(t *testing.T) {
	c, err := Compile(".count > 0")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Condition")
	}
}

func TestCompile_Empty(t *testing.T) {
	c, err := Compile("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c != nil {
		t.Fatal("expected nil Condition for empty expression")
	}
}

func TestCompile_Invalid(t *testing.T) {
	_, err := Compile(">>>invalid<<<")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

// --- BuildEnv tests ---

func TestBuildEnv_Fields(t *testing.T) {
	data := &testProvider{
		Name:  strPtr("hello"),
		Count: intPtr(42),
		Score: floatPtr(3.14),
		Label: "test",
	}
	env := BuildEnv(data, 42, "42")

	if env[".name"] != "hello" {
		t.Errorf("expected .name='hello', got %v", env[".name"])
	}
	if env[".count"] != 42 {
		t.Errorf("expected .count=42, got %v", env[".count"])
	}
	if env[".score"] != 3.14 {
		t.Errorf("expected .score=3.14, got %v", env[".score"])
	}
	if env[".label"] != "test" {
		t.Errorf("expected .label='test', got %v", env[".label"])
	}
	if env["value"] != 42 {
		t.Errorf("expected value=42, got %v", env["value"])
	}
	if env["text"] != "42" {
		t.Errorf("expected text='42', got %v", env["text"])
	}
}

func TestBuildEnv_NilPointers(t *testing.T) {
	data := &testProvider{} // all pointer fields nil
	env := BuildEnv(data, nil, "")

	if env[".name"] != "" {
		t.Errorf("expected .name='', got %v", env[".name"])
	}
	if env[".count"] != 0 {
		t.Errorf("expected .count=0, got %v", env[".count"])
	}
	if env[".score"] != 0.0 {
		t.Errorf("expected .score=0.0, got %v", env[".score"])
	}
}

func TestBuildEnv_NilProvider(t *testing.T) {
	env := BuildEnv(nil, "hello", "hello")

	if env["value"] != "hello" {
		t.Errorf("expected value='hello', got %v", env["value"])
	}
	if env["text"] != "hello" {
		t.Errorf("expected text='hello', got %v", env["text"])
	}
	// Should not have dot-fields
	if _, ok := env[".name"]; ok {
		t.Error("expected no .name for nil provider")
	}
}

func TestBuildEnv_NilValue(t *testing.T) {
	env := BuildEnv(nil, nil, "")

	if env["value"] != nil {
		t.Errorf("expected value=nil, got %v", env["value"])
	}
}

// --- Evaluate tests ---

func TestEvaluate_NumericComparisons(t *testing.T) {
	tests := []struct {
		expr     string
		expected bool
	}{
		{".count >= 50", false},
		{".count >= 42", true},
		{".count > 41", true},
		{".count > 42", false},
		{".count < 43", true},
		{".count <= 42", true},
		{".count == 42", true},
		{".count != 42", false},
	}

	for _, tt := range tests {
		c, err := Compile(tt.expr)
		if err != nil {
			t.Fatalf("Compile(%q) error: %v", tt.expr, err)
		}
		env := BuildEnv(&testProvider{Count: intPtr(42)}, 42, "42")
		result := c.Evaluate(env)
		if result != tt.expected {
			t.Errorf("Evaluate(%q) = %v, want %v", tt.expr, result, tt.expected)
		}
	}
}

func TestEvaluate_StringComparisons(t *testing.T) {
	c, _ := Compile(".name == 'main'")
	env := BuildEnv(&testProvider{Name: strPtr("main")}, "main", "main")
	if !c.Evaluate(env) {
		t.Error("expected true for .name == 'main'")
	}

	env = BuildEnv(&testProvider{Name: strPtr("feat")}, "feat", "feat")
	if c.Evaluate(env) {
		t.Error("expected false for .name == 'main' when name is 'feat'")
	}
}

func TestEvaluate_BooleanCombinators(t *testing.T) {
	c, _ := Compile(".count > 0 && .name != ''")
	env := BuildEnv(&testProvider{Count: intPtr(5), Name: strPtr("test")}, nil, "")
	if !c.Evaluate(env) {
		t.Error("expected true for both conditions met")
	}

	env = BuildEnv(&testProvider{Count: intPtr(0), Name: strPtr("test")}, nil, "")
	if c.Evaluate(env) {
		t.Error("expected false when count is 0")
	}
}

func TestEvaluate_NilCoercion(t *testing.T) {
	c, _ := Compile(".count >= 50")
	env := BuildEnv(&testProvider{}, nil, "") // Count is nil → coerced to 0
	if c.Evaluate(env) {
		t.Error("expected false: nil count coerced to 0, 0 >= 50 is false")
	}
}

func TestEvaluate_ValueKeyword(t *testing.T) {
	c, _ := Compile("value > 0")
	env := BuildEnv(nil, 42, "42")
	if !c.Evaluate(env) {
		t.Error("expected true for value > 0 with value=42")
	}

	env = BuildEnv(nil, 0, "0")
	if c.Evaluate(env) {
		t.Error("expected false for value > 0 with value=0")
	}
}

func TestEvaluate_TextKeyword(t *testing.T) {
	c, _ := Compile("text != ''")
	env := BuildEnv(nil, nil, "hello")
	if !c.Evaluate(env) {
		t.Error("expected true for text != '' with text='hello'")
	}

	env = BuildEnv(nil, nil, "")
	if c.Evaluate(env) {
		t.Error("expected false for text != '' with text=''")
	}
}

func TestEvaluate_NilCondition(t *testing.T) {
	// nil Condition (empty expression) always returns true
	var c *Condition
	if !c.Evaluate(map[string]any{}) {
		t.Error("expected true for nil Condition")
	}
}

func TestEvaluate_NonBoolResult(t *testing.T) {
	c, _ := Compile(".count + 1")
	env := BuildEnv(&testProvider{Count: intPtr(5)}, nil, "")
	if c.Evaluate(env) {
		t.Error("expected false for non-bool result")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/condition/ -v`
Expected: FAIL (package doesn't exist)

- [ ] **Step 3: Implement condition package**

Create `internal/condition/condition.go`:

```go
package condition

import (
	"reflect"
	"strings"

	"github.com/expr-lang/expr"
)

// Condition is a compiled when expression.
type Condition struct {
	program *expr.Program
}

// Compile compiles an expression string into a reusable Condition.
// Returns nil for empty expressions (always true).
func Compile(expression string) (*Condition, error) {
	if expression == "" {
		return nil, nil
	}

	// expr-lang doesn't allow dots in variable names by default,
	// so we compile with AllowUndefinedVariables to handle our
	// dot-prefixed field names.
	program, err := expr.Compile(expression, expr.AllowUndefinedVariables())
	if err != nil {
		return nil, err
	}

	return &Condition{program: program}, nil
}

// Evaluate runs the compiled expression against the environment.
// Returns true only if the result is boolean true.
// Nil receiver (empty expression) returns true.
func (c *Condition) Evaluate(env map[string]any) bool {
	if c == nil {
		return true
	}

	result, err := expr.Run(c.program, env)
	if err != nil {
		return false
	}

	b, ok := result.(bool)
	return ok && b
}

// BuildEnv builds the variable environment for expression evaluation.
// providerData is the provider struct (e.g., *GitData).
// value is the raw segment value. text is the formatted display string.
func BuildEnv(providerData any, value any, text string) map[string]any {
	env := make(map[string]any)

	// Register provider fields with dot-prefixed lowercase names
	if providerData != nil {
		v := reflect.ValueOf(providerData)
		if v.Kind() == reflect.Ptr {
			if !v.IsNil() {
				v = v.Elem()
			} else {
				v = reflect.Value{}
			}
		}

		if v.IsValid() && v.Kind() == reflect.Struct {
			t := v.Type()
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				if !field.IsExported() {
					continue
				}

				key := "." + strings.ToLower(field.Name)
				fv := v.Field(i)

				// Dereference pointers with nil coercion
				if fv.Kind() == reflect.Ptr {
					if fv.IsNil() {
						env[key] = coerceNilPointer(field.Type)
						continue
					}
					fv = fv.Elem()
				}

				env[key] = fv.Interface()
			}
		}
	}

	// Register segment keywords
	env["value"] = value
	env["text"] = text

	return env
}

// coerceNilPointer returns the zero value for the pointer's element type.
func coerceNilPointer(t reflect.Type) any {
	elem := t.Elem()
	switch elem.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0
	case reflect.Float32, reflect.Float64:
		return 0.0
	case reflect.String:
		return ""
	case reflect.Bool:
		return false
	default:
		return nil
	}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/condition/ -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/condition/ go.mod go.sum
git commit -m "feat(condition): add when expression compilation and evaluation"
```

## Chunk 2: Render pipeline integration

### Task 3: Wire when into the render pipeline

**Files:**
- Modify: `internal/render/render.go`
- Modify: `internal/render/render_test.go`

- [ ] **Step 1: Write failing render tests for when**

Add to `internal/render/render_test.go`:

```go
func TestTree_DataSegmentWhenPasses(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.count": 75}
	tagIdx := TagIndex{"test.count": fieldAccessor{Provider: "test", FieldIndex: 1}}

	tree := []types.SegmentNode{
		{Type: "test.count", When: "value >= 50"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "75" {
		t.Errorf("expected '75', got %q", result)
	}
}

func TestTree_DataSegmentWhenFails(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.count": 25}
	tagIdx := TagIndex{"test.count": fieldAccessor{Provider: "test", FieldIndex: 1}}

	tree := []types.SegmentNode{
		{Type: "test.count", When: "value >= 50"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty (when failed), got %q", result)
	}
}

func TestTree_DataSegmentWhenDotField(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	name := "feature"
	count := 5
	type testStruct struct {
		Name  *string `segment:"test.name"`
		Count *int    `segment:"test.count"`
	}
	providerData := map[string]any{
		"test": &testStruct{Name: &name, Count: &count},
	}
	segmentValues := map[string]any{"test.name": "feature"}
	tagIdx := TagIndex{"test.name": fieldAccessor{Provider: "test", FieldIndex: 0}}

	tree := []types.SegmentNode{
		{Type: "test.name", When: ".count > 0"},
	}

	result := Tree(tree, seg, sess, providerData, segmentValues, tagIdx)
	if result != "feature" {
		t.Errorf("expected 'feature', got %q", result)
	}
}

func TestTree_DataSegmentWhenText(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.name": "hello"}
	tagIdx := TagIndex{"test.name": fieldAccessor{Provider: "test", FieldIndex: 0}}

	tree := []types.SegmentNode{
		{Type: "test.name", When: "text != ''"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_CompositeWhen(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	name := "main"
	type testStruct struct {
		Name *string `segment:"test.name"`
	}
	providerData := map[string]any{
		"test": &testStruct{Name: &name},
	}
	segmentValues := map[string]any{"test.name": "main"}
	tagIdx := TagIndex{"test.name": fieldAccessor{Provider: "test", FieldIndex: 0}}

	tree := []types.SegmentNode{
		{
			Provider: "test",
			When:     ".name != ''",
			Children: []types.SegmentNode{
				{Type: "test.name"},
			},
		},
	}

	result := Tree(tree, seg, sess, providerData, segmentValues, tagIdx)
	if result != "main" {
		t.Errorf("expected 'main', got %q", result)
	}
}

func TestTree_CompositeWhenFails(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	name := ""
	type testStruct struct {
		Name *string `segment:"test.name"`
	}
	providerData := map[string]any{
		"test": &testStruct{Name: &name},
	}
	segmentValues := map[string]any{"test.name": ""}
	tagIdx := TagIndex{"test.name": fieldAccessor{Provider: "test", FieldIndex: 0}}

	tree := []types.SegmentNode{
		{
			Provider: "test",
			When:     ".name != ''",
			Children: []types.SegmentNode{
				{Type: "literal", Props: map[string]any{"text": "should not appear"}},
			},
		},
	}

	result := Tree(tree, seg, sess, providerData, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty (composite when failed), got %q", result)
	}
}

func TestTree_WhenNoExpression(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.name": "hello"}
	tagIdx := TagIndex{"test.name": fieldAccessor{Provider: "test", FieldIndex: 0}}

	tree := []types.SegmentNode{
		{Type: "test.name"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_WhenInvalidExpression(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.name": "hello"}
	tagIdx := TagIndex{"test.name": fieldAccessor{Provider: "test", FieldIndex: 0}}

	tree := []types.SegmentNode{
		{Type: "test.name", When: ">>>bad<<<"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty (invalid when), got %q", result)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/render/ -run "TestTree_.*When" -v`
Expected: FAIL (When field not evaluated)

- [ ] **Step 3: Implement when evaluation in renderNode**

Update `internal/render/render.go`. Add import:

```go
import (
	"log"
	"strings"
	"sync"

	"github.com/jheddings/ccglow/internal/condition"
	"github.com/jheddings/ccglow/internal/segment"
	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
)
```

Add a condition cache and helper at the top of the file:

```go
// conditionCache caches compiled when expressions.
var conditionCache = make(map[string]*condition.Condition)
var conditionMu sync.Mutex

func getCondition(expr string) *condition.Condition {
	if expr == "" {
		return nil
	}
	conditionMu.Lock()
	defer conditionMu.Unlock()

	if c, ok := conditionCache[expr]; ok {
		return c
	}

	c, err := condition.Compile(expr)
	if err != nil {
		log.Printf("ccglow: invalid when expression %q: %v", expr, err)
		conditionCache[expr] = nil
		return nil
	}
	conditionCache[expr] = c
	return c
}
```

Update the composite branch in `renderNode` to evaluate `when` before recursing:

```go
	// SegmentGroup: evaluate when, then render children
	if len(node.Children) > 0 {
		if node.When != "" {
			c := getCondition(node.When)
			if c != nil {
				var pd any
				if node.Provider != "" {
					pd = providerData[node.Provider]
				}
				env := condition.BuildEnv(pd, nil, "")
				if !c.Evaluate(env) {
					return nil
				}
			} else if node.When != "" {
				// Compilation failed (nil cached) → hide
				return nil
			}
		}

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
```

Update the data segment branch to evaluate `when` after formatting:

```go
	// DataSegment: resolve from segment values
	value, ok := segmentValues[node.Type]
	if !ok || value == nil {
		return nil
	}

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

	// Evaluate when expression
	if node.When != "" {
		c := getCondition(node.When)
		if c != nil {
			var pd any
			if accessor, exists := tagIdx[node.Type]; exists {
				pd = providerData[accessor.Provider]
			}
			env := condition.BuildEnv(pd, value, text)
			if !c.Evaluate(env) {
				return nil
			}
		} else {
			return nil // compilation failed
		}
	}

	styled := style.Apply(text, node.Style)
	return &styled
```

Also update `collectNames` to include providers referenced by composite
`Provider` fields (needed for `when` evaluation on composites):

```go
func collectNames(nodes []types.SegmentNode, names map[string]bool, idx TagIndex) {
	for _, node := range nodes {
		if node.Enabled != nil && !*node.Enabled {
			continue
		}
		if accessor, ok := idx[node.Type]; ok {
			names[accessor.Provider] = true
		}
		if node.Provider != "" {
			names[node.Provider] = true
		}
		if len(node.Children) > 0 {
			collectNames(node.Children, names, idx)
		}
	}
}
```

- [ ] **Step 4: Run when tests**

Run: `go test ./internal/render/ -run "TestTree_.*When" -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/render/render.go internal/render/render_test.go
git commit -m "feat(render): evaluate when expressions in render pipeline"
```

## Chunk 3: Minimal preset showcase

### Task 4: Update minimal preset with when conditions

**Files:**
- Modify: `internal/preset/minimal.json`

- [ ] **Step 1: Replace minimal preset**

```json
{
  "segments": [
    {
      "segment": "pwd.name",
      "style": { "color": "39" }
    },
    {
      "segment": "git.branch",
      "when": ".branch != '' && .branch != 'main'",
      "style": { "color": "whiteBright", "bold": true, "prefix": " | " }
    },
    {
      "segment": "git.modified",
      "when": "value > 0",
      "style": { "color": "yellow", "prefix": " ~", "format": "%d" }
    },
    {
      "segment": "git.untracked",
      "when": "value > 0",
      "style": { "color": "cyan", "prefix": " ?", "format": "%d" }
    },
    {
      "segment": "context.percent.used",
      "when": ".percent >= 50",
      "style": { "color": "yellow", "prefix": " | " }
    },
    {
      "segment": "cost.usd",
      "when": "text != ''",
      "style": { "color": "green", "prefix": " | " }
    },
    {
      "segment": "model.name",
      "when": "text != ''",
      "style": { "color": "240", "prefix": " | " }
    }
  ]
}
```

- [ ] **Step 2: Run preset tests**

Run: `go test ./internal/preset/ -v`
Expected: PASS

- [ ] **Step 3: Run full suite and build**

Run: `go vet ./... && go test ./... && go build ./...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/preset/minimal.json
git commit -m "feat(preset): update minimal as when-clause showcase"
```

### Task 5: Final verification

- [ ] **Step 1: Full build and tests**

Run: `go vet ./... && go test ./... && go build ./...`
Expected: all PASS

- [ ] **Step 2: Smoke test minimal — quiet state**

Run:
```bash
echo '{"cwd":"/tmp"}' | go run . --preset minimal --format plain
```
Expected: `tmp` — just the directory name, everything else hidden

- [ ] **Step 3: Smoke test minimal — on main branch in a git repo**

Run:
```bash
echo '{"cwd":"'$(pwd)'"}' | go run . --preset minimal --format plain
```
Expected: directory name only (branch hidden because it's `main` or the when-clause branch)

Wait — we're on `feat/when-clause`, not `main`. So the branch should show.
Expected: `refactor-segments | feat/when-clause` (or similar, with branch visible)

- [ ] **Step 4: Smoke test minimal — high context usage**

Run:
```bash
echo '{"cwd":"/tmp","context_window":{"used_percentage":75,"current_usage":{"input_tokens":100}},"cost":{"total_cost_usd":1.23},"model":{"display_name":"Opus 4.6"}}' | go run . --preset minimal --format plain
```
Expected: `tmp | 75% | $1.23 | Opus 4.6` — context percent and cost appear because thresholds met

- [ ] **Step 5: Smoke test minimal — low context usage**

Run:
```bash
echo '{"cwd":"/tmp","context_window":{"used_percentage":20,"current_usage":{"input_tokens":100}},"model":{"display_name":"Opus 4.6"}}' | go run . --preset minimal --format plain
```
Expected: `tmp | Opus 4.6` — context percent hidden (below 50%), no cost data
