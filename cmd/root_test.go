package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_ReturnsOutputAndEnv(t *testing.T) {
	stdin := `{"cwd": "/tmp"}`

	output, env := run("default", "", "plain", stdin)

	if output == "" {
		t.Fatal("expected rendered output, got empty string")
	}

	if env == nil {
		t.Fatal("expected env map, got nil")
	}

	data, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("expected env to be JSON-serializable, got error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty JSON from env")
	}
}

func TestDump_WritesEnvToFile(t *testing.T) {
	stdin := `{"cwd": "/tmp"}`
	dumpPath := filepath.Join(t.TempDir(), "env.json")

	_, env := run("default", "", "plain", stdin)

	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal env: %v", err)
	}

	if err := os.WriteFile(dumpPath, data, 0644); err != nil {
		t.Fatalf("failed to write dump file: %v", err)
	}

	contents, err := os.ReadFile(dumpPath)
	if err != nil {
		t.Fatalf("failed to read dump file: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(contents, &parsed); err != nil {
		t.Fatalf("dump file is not valid JSON: %v", err)
	}
}
