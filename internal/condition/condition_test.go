package condition

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

func TestBuildNestedEnv(t *testing.T) {
	values := map[string]any{
		"git.branch":            "main",
		"git.repo":              "ccglow",
		"pwd.name":              "project",
		"context.percent.used":  36,
		"context.tokens":        "360K",
	}
	env := BuildNestedEnv(values)

	git, ok := env["git"].(map[string]any)
	if !ok {
		t.Fatal("expected git namespace")
	}
	if git["branch"] != "main" {
		t.Errorf("expected git.branch='main', got %v", git["branch"])
	}
	if git["repo"] != "ccglow" {
		t.Errorf("expected git.repo='ccglow', got %v", git["repo"])
	}

	pwd, ok := env["pwd"].(map[string]any)
	if !ok {
		t.Fatal("expected pwd namespace")
	}
	if pwd["name"] != "project" {
		t.Errorf("expected pwd.name='project', got %v", pwd["name"])
	}

	ctx, ok := env["context"].(map[string]any)
	if !ok {
		t.Fatal("expected context namespace")
	}
	pct, ok := ctx["percent"].(map[string]any)
	if !ok {
		t.Fatal("expected context.percent namespace")
	}
	if pct["used"] != 36 {
		t.Errorf("expected context.percent.used=36, got %v", pct["used"])
	}
	if ctx["tokens"] != "360K" {
		t.Errorf("expected context.tokens='360K', got %v", ctx["tokens"])
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

	values := map[string]any{
		"context.percent.used": 36,
	}
	nested := BuildNestedEnv(values)

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
	values := map[string]any{
		"git.branch": "main",
	}
	nested := BuildNestedEnv(values)

	c, _ := Compile("git.branch == 'main'")
	env := BuildSegmentEnv(nested, "main", "main")
	if !c.Evaluate(env) {
		t.Error("expected true for git.branch == 'main'")
	}

	values["git.branch"] = "feat"
	nested = BuildNestedEnv(values)
	env = BuildSegmentEnv(nested, "feat", "feat")
	if c.Evaluate(env) {
		t.Error("expected false for git.branch == 'main' when branch is 'feat'")
	}
}

func TestEvaluate_BooleanCombinators(t *testing.T) {
	values := map[string]any{
		"git.modified": 5,
		"git.branch":   "test",
	}
	nested := BuildNestedEnv(values)

	c, _ := Compile("git.modified > 0 && git.branch != ''")
	env := BuildSegmentEnv(nested, nil, "")
	if !c.Evaluate(env) {
		t.Error("expected true for both conditions met")
	}

	values["git.modified"] = 0
	nested = BuildNestedEnv(values)
	env = BuildSegmentEnv(nested, nil, "")
	if c.Evaluate(env) {
		t.Error("expected false when modified is 0")
	}
}

func TestEvaluate_ZeroValueDefaults(t *testing.T) {
	// With no nil values, zeros are the defaults
	values := map[string]any{
		"context.percent.used": 0,
	}
	nested := BuildNestedEnv(values)

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
	// This is the #42 feature: pwd segment referencing git data
	values := map[string]any{
		"git.repo": "",
		"pwd.name": "mydir",
	}
	nested := BuildNestedEnv(values)

	c, err := Compile("git.repo == ''")
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}
	env := BuildSegmentEnv(nested, "mydir", "mydir")
	if !c.Evaluate(env) {
		t.Error("expected true for git.repo == '' cross-provider reference")
	}

	// Now with a repo set
	values["git.repo"] = "ccglow"
	nested = BuildNestedEnv(values)
	env = BuildSegmentEnv(nested, "mydir", "mydir")
	if c.Evaluate(env) {
		t.Error("expected false for git.repo == '' when repo is set")
	}
}
