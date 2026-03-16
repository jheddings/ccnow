package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

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

	if result.Values["model.name"] != "Opus 4.6" {
		t.Errorf("expected Opus 4.6, got %s", result.Values["model.name"])
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

	if result.Values["model.id"] != "claude-opus-4-6[1m]" {
		t.Errorf("expected claude-opus-4-6[1m], got %v", result.Values["model.id"])
	}
}

func TestModelProviderNoModel(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["model.id"] != "" {
		t.Errorf("expected empty ID, got %v", result.Values["model.id"])
	}
	if result.Values["model.name"] != "" {
		t.Errorf("expected empty Name, got %v", result.Values["model.name"])
	}
}
