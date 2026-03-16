package condition

import (
	"testing"
)

type testProvider struct {
	Name    *string
	Count   *int
	Score   *float64
	Label   string
	Enabled bool
}

func strPtr(s string) *string    { return &s }
func intPtr(n int) *int          { return &n }
func floatPtr(f float64) *float64 { return &f }

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
	data := &testProvider{}
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
	env := BuildEnv(&testProvider{}, nil, "")
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
