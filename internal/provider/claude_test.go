package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestClaudeProvider(t *testing.T) {
	p := &claudeProvider{}
	sess := &types.SessionData{
		CWD:         "/tmp",
		Version:     "2.1.75",
		OutputStyle: &types.OutputStyleInfo{Name: "concise"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ClaudeData)
	if data.Version == nil || *data.Version != "2.1.75" {
		t.Errorf("expected version 2.1.75, got %v", data.Version)
	}
	if data.Style == nil || *data.Style != "concise" {
		t.Errorf("expected style concise, got %v", data.Style)
	}
}

func TestClaudeProviderEmpty(t *testing.T) {
	p := &claudeProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ClaudeData)
	if data.Version != nil {
		t.Errorf("expected nil version, got %v", data.Version)
	}
	if data.Style != nil {
		t.Errorf("expected nil style, got %v", data.Style)
	}
}
