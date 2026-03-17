package command

import (
	"os"
	"path/filepath"
	"testing"
)

// --- Interpolate tests ---

func TestInterpolate_NoVars(t *testing.T) {
	result := Interpolate("echo hello", nil)
	if result != "echo hello" {
		t.Errorf("expected 'echo hello', got %q", result)
	}
}

func TestInterpolate_SingleVar(t *testing.T) {
	env := map[string]any{
		"git": map[string]any{"branch": "main"},
	}
	result := Interpolate("git log ${git.branch}", env)
	if result != "git log main" {
		t.Errorf("expected 'git log main', got %q", result)
	}
}

func TestInterpolate_MultipleVars(t *testing.T) {
	env := map[string]any{
		"git": map[string]any{"owner": "jheddings", "repo": "ccglow"},
	}
	result := Interpolate("gh pr list --repo ${git.owner}/${git.repo}", env)
	expected := "gh pr list --repo jheddings/ccglow"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestInterpolate_MissingVar(t *testing.T) {
	env := map[string]any{}
	result := Interpolate("echo ${nonexistent.var}", env)
	if result != "echo " {
		t.Errorf("expected 'echo ', got %q", result)
	}
}

func TestInterpolate_NumericVar(t *testing.T) {
	env := map[string]any{
		"test": map[string]any{"count": 42},
	}
	result := Interpolate("echo ${test.count}", env)
	if result != "echo 42" {
		t.Errorf("expected 'echo 42', got %q", result)
	}
}

func TestInterpolate_NilEnv(t *testing.T) {
	result := Interpolate("echo ${test.val}", nil)
	if result != "echo " {
		t.Errorf("expected 'echo ', got %q", result)
	}
}

// --- Run tests ---

func TestRun_SimpleEcho(t *testing.T) {
	result := Run("echo hello", nil, "", DefaultTimeout)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestRun_EmptyOutputCollapses(t *testing.T) {
	result := Run("printf ''", nil, "", DefaultTimeout)
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestRun_NonZeroExitCollapses(t *testing.T) {
	result := Run("exit 1", nil, "", DefaultTimeout)
	if result != "" {
		t.Errorf("expected empty for non-zero exit, got %q", result)
	}
}

func TestRun_TimeoutCollapses(t *testing.T) {
	result := Run("sleep 10", nil, "", 1) // 1ns timeout
	if result != "" {
		t.Errorf("expected empty for timeout, got %q", result)
	}
}

func TestRun_WithInterpolation(t *testing.T) {
	env := map[string]any{
		"test": map[string]any{"name": "world"},
	}
	result := Run("echo ${test.name}", env, "", DefaultTimeout)
	if result != "world" {
		t.Errorf("expected 'world', got %q", result)
	}
}

func TestRun_TrimsWhitespace(t *testing.T) {
	result := Run("echo '  hello  '", nil, "", DefaultTimeout)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestRun_EmptyCommandString(t *testing.T) {
	result := Run("", nil, "", DefaultTimeout)
	if result != "" {
		t.Errorf("expected empty for empty command, got %q", result)
	}
}

func TestRun_CWDRespected(t *testing.T) {
	dir := t.TempDir()
	// Resolve symlinks to handle macOS /tmp -> /private/tmp
	dir, _ = filepath.EvalSymlinks(dir)

	result := Run("pwd", nil, dir, DefaultTimeout)

	// Normalize both paths
	resultResolved, _ := filepath.EvalSymlinks(result)
	if resultResolved == "" {
		resultResolved = result
	}

	if resultResolved != dir {
		t.Errorf("expected %q, got %q", dir, resultResolved)
	}
}

func TestRun_CWDEmptyUsesProcess(t *testing.T) {
	result := Run("pwd", nil, "", DefaultTimeout)
	cwd, _ := os.Getwd()
	cwd, _ = filepath.EvalSymlinks(cwd)
	resultResolved, _ := filepath.EvalSymlinks(result)
	if resultResolved != cwd {
		t.Errorf("expected %q, got %q", cwd, resultResolved)
	}
}
