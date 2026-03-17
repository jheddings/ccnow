package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func modelValues(result *types.ProviderResult) map[string]any {
	return result.Values["model"].(map[string]any)
}

func TestModelProvider(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{
		CWD:   "/tmp",
		Model: &types.ModelInfo{DisplayName: "Opus 4.6"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	m := modelValues(result)
	if m["name"] != "Opus 4.6" {
		t.Errorf("expected Opus 4.6, got %s", m["name"])
	}
}

func TestModelProviderID(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{
		CWD:   "/tmp",
		Model: &types.ModelInfo{ID: "claude-opus-4-6[1m]", DisplayName: "Opus 4.6"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	m := modelValues(result)
	if m["id"] != "claude-opus-4-6[1m]" {
		t.Errorf("expected claude-opus-4-6[1m], got %v", m["id"])
	}
}

func TestModelProviderNoModel(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	m := modelValues(result)
	if m["id"] != "" {
		t.Errorf("expected empty ID, got %v", m["id"])
	}
	if m["name"] != "" {
		t.Errorf("expected empty Name, got %v", m["name"])
	}
}
