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

	data := result.(*ModelData)
	if *data.Name != "Opus 4.6" {
		t.Errorf("expected Opus 4.6, got %s", *data.Name)
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

	data := result.(*ModelData)
	if data.ID == nil || *data.ID != "claude-opus-4-6[1m]" {
		t.Errorf("expected claude-opus-4-6[1m], got %v", data.ID)
	}
}

func TestModelProviderNoModel(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ModelData)
	if data.ID != nil {
		t.Errorf("expected nil ID, got %v", data.ID)
	}
	if data.Name != nil {
		t.Errorf("expected nil Name, got %v", data.Name)
	}
}
