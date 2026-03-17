package eval

import (
	"testing"
)

func TestCompile_Valid(t *testing.T) {
	c, err := Compile("git.branch != ''")
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

func TestCompileCached_Valid(t *testing.T) {
	c := CompileCached("git.branch != ''")
	if c == nil {
		t.Fatal("expected non-nil Condition")
	}
}

func TestCompileCached_Empty(t *testing.T) {
	c := CompileCached("")
	if c != nil {
		t.Fatal("expected nil for empty expression")
	}
}

func TestCompileCached_Invalid(t *testing.T) {
	c := CompileCached(">>>bad<<<")
	if c != nil {
		t.Fatal("expected nil for invalid expression")
	}
}

func TestBuildSegmentEnv(t *testing.T) {
	nested := map[string]any{
		"git": map[string]any{"branch": "main"},
	}
	env := BuildSegmentEnv(nested, 42, "42")

	if env["value"] != 42 {
		t.Errorf("expected value=42, got %v", env["value"])
	}
	if env["text"] != "42" {
		t.Errorf("expected text='42', got %v", env["text"])
	}
	git, ok := env["git"].(map[string]any)
	if !ok || git["branch"] != "main" {
		t.Error("expected nested env preserved")
	}
}

func TestBuildSegmentEnv_NilValue(t *testing.T) {
	env := BuildSegmentEnv(map[string]any{}, nil, "")
	if env["value"] != nil {
		t.Errorf("expected value=nil, got %v", env["value"])
	}
}

func TestEvaluate_NumericComparisons(t *testing.T) {
	tests := []struct {
		expr     string
		expected bool
	}{
		{"context.percent.used >= 50", false},
		{"context.percent.used >= 36", true},
		{"context.percent.used > 35", true},
		{"context.percent.used > 36", false},
		{"context.percent.used < 37", true},
		{"context.percent.used <= 36", true},
		{"context.percent.used == 36", true},
		{"context.percent.used != 36", false},
	}

	nested := map[string]any{
		"context": map[string]any{
			"percent": map[string]any{"used": 36},
		},
	}

	for _, tt := range tests {
		c, err := Compile(tt.expr)
		if err != nil {
			t.Fatalf("Compile(%q) error: %v", tt.expr, err)
		}
		env := BuildSegmentEnv(nested, 36, "36")
		result := c.Evaluate(env)
		if result != tt.expected {
			t.Errorf("Evaluate(%q) = %v, want %v", tt.expr, result, tt.expected)
		}
	}
}

func TestEvaluate_StringComparisons(t *testing.T) {
	nested := map[string]any{
		"git": map[string]any{"branch": "main"},
	}

	c, _ := Compile("git.branch == 'main'")
	env := BuildSegmentEnv(nested, "main", "main")
	if !c.Evaluate(env) {
		t.Error("expected true for git.branch == 'main'")
	}

	nested = map[string]any{
		"git": map[string]any{"branch": "feat"},
	}
	env = BuildSegmentEnv(nested, "feat", "feat")
	if c.Evaluate(env) {
		t.Error("expected false for git.branch == 'main' when branch is 'feat'")
	}
}

func TestEvaluate_BooleanCombinators(t *testing.T) {
	nested := map[string]any{
		"git": map[string]any{
			"modified": 5,
			"branch":   "test",
		},
	}

	c, _ := Compile("git.modified > 0 && git.branch != ''")
	env := BuildSegmentEnv(nested, nil, "")
	if !c.Evaluate(env) {
		t.Error("expected true for both conditions met")
	}

	nested = map[string]any{
		"git": map[string]any{
			"modified": 0,
			"branch":   "test",
		},
	}
	env = BuildSegmentEnv(nested, nil, "")
	if c.Evaluate(env) {
		t.Error("expected false when modified is 0")
	}
}

func TestEvaluate_ZeroValueDefaults(t *testing.T) {
	nested := map[string]any{
		"context": map[string]any{
			"percent": map[string]any{"used": 0},
		},
	}

	c, _ := Compile("context.percent.used >= 50")
	env := BuildSegmentEnv(nested, nil, "")
	if c.Evaluate(env) {
		t.Error("expected false: zero percent, 0 >= 50 is false")
	}
}

func TestEvaluate_ValueKeyword(t *testing.T) {
	c, _ := Compile("value > 0")
	env := BuildSegmentEnv(map[string]any{}, 42, "42")
	if !c.Evaluate(env) {
		t.Error("expected true for value > 0 with value=42")
	}

	env = BuildSegmentEnv(map[string]any{}, 0, "0")
	if c.Evaluate(env) {
		t.Error("expected false for value > 0 with value=0")
	}
}

func TestEvaluate_TextKeyword(t *testing.T) {
	c, _ := Compile("text != ''")
	env := BuildSegmentEnv(map[string]any{}, nil, "hello")
	if !c.Evaluate(env) {
		t.Error("expected true for text != '' with text='hello'")
	}

	env = BuildSegmentEnv(map[string]any{}, nil, "")
	if c.Evaluate(env) {
		t.Error("expected false for text != '' with text=''")
	}
}

func TestEvaluate_NilCondition(t *testing.T) {
	var c *Condition
	if !c.Evaluate(map[string]any{}) {
		t.Error("expected true for nil Condition")
	}
}

func TestEvaluate_NonBoolResult(t *testing.T) {
	c, _ := Compile("value + 1")
	env := BuildSegmentEnv(map[string]any{}, 5, "")
	if c.Evaluate(env) {
		t.Error("expected false for non-bool result")
	}
}

func TestEvaluate_CrossProviderReference(t *testing.T) {
	nested := map[string]any{
		"git": map[string]any{"repo": ""},
		"pwd": map[string]any{"name": "mydir"},
	}

	c, err := Compile("git.repo == ''")
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}
	env := BuildSegmentEnv(nested, "mydir", "mydir")
	if !c.Evaluate(env) {
		t.Error("expected true for git.repo == '' cross-provider reference")
	}

	nested = map[string]any{
		"git": map[string]any{"repo": "ccglow"},
		"pwd": map[string]any{"name": "mydir"},
	}
	env = BuildSegmentEnv(nested, "mydir", "mydir")
	if c.Evaluate(env) {
		t.Error("expected false for git.repo == '' when repo is set")
	}
}

func TestEval_Expression(t *testing.T) {
	env := map[string]any{
		"git": map[string]any{"branch": "main"},
	}

	result, err := Eval("git.branch", env)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if result != "main" {
		t.Errorf("expected 'main', got %v", result)
	}
}

func TestEval_Arithmetic(t *testing.T) {
	env := map[string]any{
		"git": map[string]any{"insertions": 10, "deletions": 5},
	}

	result, err := Eval("git.insertions + git.deletions", env)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if result != 15 {
		t.Errorf("expected 15, got %v", result)
	}
}

func TestEval_Empty(t *testing.T) {
	result, err := Eval("", map[string]any{})
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestEval_Invalid(t *testing.T) {
	_, err := Eval(">>>bad<<<", map[string]any{})
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestEval_Cached(t *testing.T) {
	env := map[string]any{"x": 1}
	r1, _ := Eval("x + 1", env)
	r2, _ := Eval("x + 1", env)
	if r1 != r2 {
		t.Error("cached results should match")
	}
}

func TestEval_UndefinedVariable(t *testing.T) {
	_, err := Eval("missing.field", map[string]any{})
	if err == nil {
		t.Error("expected error for undefined variable field access")
	}
}

func TestEval_UndefinedTopLevel(t *testing.T) {
	result, err := Eval("missing", map[string]any{})
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil for undefined top-level variable, got %v", result)
	}
}
