package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestPresetList(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"preset", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, name := range []string{"default", "minimal", "full"} {
		if !strings.Contains(output, name) {
			t.Errorf("expected %q in output, got: %s", name, output)
		}
	}
}

func TestPresetShow_Valid(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"preset", "show", "default"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "{") {
		t.Errorf("expected JSON output, got: %s", output)
	}
}

func TestPresetShow_Invalid(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"preset", "show", "nonexistent"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for unknown preset")
	}
}
