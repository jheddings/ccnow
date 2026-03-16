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

	if result.Values["claude.version"] != "2.1.75" {
		t.Errorf("expected version 2.1.75, got %v", result.Values["claude.version"])
	}
	if result.Values["claude.style"] != "concise" {
		t.Errorf("expected style concise, got %v", result.Values["claude.style"])
	}
}

func TestClaudeProviderEmpty(t *testing.T) {
	p := &claudeProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["claude.version"] != "" {
		t.Errorf("expected empty version, got %v", result.Values["claude.version"])
	}
	if result.Values["claude.style"] != "" {
		t.Errorf("expected empty style, got %v", result.Values["claude.style"])
	}
}
