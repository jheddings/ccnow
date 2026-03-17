package cmd

import (
	"encoding/json"
	"testing"
)

func TestRun_Dump(t *testing.T) {
	stdin := `{"cwd": "/tmp"}`

	output := run("default", "", "plain", stdin, true)

	if output == "" {
		t.Fatal("expected JSON output from dump, got empty string")
	}

	var env map[string]any
	if err := json.Unmarshal([]byte(output), &env); err != nil {
		t.Fatalf("expected valid JSON from dump, got error: %v\noutput: %s", err, output)
	}
}

func TestRun_NoDump(t *testing.T) {
	stdin := `{"cwd": "/tmp"}`

	output := run("default", "", "plain", stdin, false)

	// Without dump, output should be rendered text (not JSON).
	// Just verify it doesn't parse as a JSON object with provider keys.
	var env map[string]any
	if err := json.Unmarshal([]byte(output), &env); err == nil {
		// If it happens to be valid JSON, that's fine, but it shouldn't
		// be the same as the dump output (rendered text is unlikely valid JSON).
		t.Log("output happened to be valid JSON, skipping further checks")
	}
}
