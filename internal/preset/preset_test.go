package preset

import (
	"encoding/json"
	"testing"
)

func TestGet_Default(t *testing.T) {
	nodes := Get("default")
	if nodes == nil {
		t.Fatal("expected default preset, got nil")
	}
	if len(nodes) == 0 {
		t.Fatal("expected non-empty segment tree")
	}
}

func TestGet_Minimal(t *testing.T) {
	nodes := Get("minimal")
	if nodes == nil {
		t.Fatal("expected minimal preset, got nil")
	}
}

func TestGet_Full(t *testing.T) {
	nodes := Get("full")
	if nodes == nil {
		t.Fatal("expected full preset, got nil")
	}
}

func TestGet_Unknown(t *testing.T) {
	nodes := Get("nonexistent")
	if nodes != nil {
		t.Errorf("expected nil for unknown preset, got %v", nodes)
	}
}

func TestGet_InfersProviders(t *testing.T) {
	nodes := Get("minimal")
	if nodes == nil {
		t.Fatal("expected minimal preset")
	}
	// First segment is pwd.name — provider should be inferred as "pwd"
	if nodes[0].Provider != "pwd" {
		t.Errorf("expected inferred provider pwd, got %q", nodes[0].Provider)
	}
}

func TestList(t *testing.T) {
	names := List()
	if len(names) < 3 {
		t.Fatalf("expected at least 3 presets, got %d: %v", len(names), names)
	}

	expected := map[string]bool{"default": false, "minimal": false, "full": false}
	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected preset %q in list", name)
		}
	}
}

func TestDump_ReturnsValidJSON(t *testing.T) {
	data, err := Dump("default")
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(data) {
		t.Error("expected valid JSON from Dump")
	}
}

func TestDump_Unknown(t *testing.T) {
	_, err := Dump("nonexistent")
	if err == nil {
		t.Error("expected error for unknown preset")
	}
}
