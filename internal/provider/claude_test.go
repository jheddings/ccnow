package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func claudeValues(result *types.ProviderResult) map[string]any {
	return result.Values["claude"].(map[string]any)
}

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

	c := claudeValues(result)
	if c["version"] != "2.1.75" {
		t.Errorf("expected version 2.1.75, got %v", c["version"])
	}
	if c["style"] != "concise" {
		t.Errorf("expected style concise, got %v", c["style"])
	}
}

func TestClaudeProviderEmpty(t *testing.T) {
	p := &claudeProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	c := claudeValues(result)
	if c["version"] != "" {
		t.Errorf("expected empty version, got %v", c["version"])
	}
	if c["style"] != "" {
		t.Errorf("expected empty style, got %v", c["style"])
	}
}
