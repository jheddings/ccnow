package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{500, "500"},
		{1000, "1K"},
		{24500, "24K"},
		{1000000, "1M"},
		{1500000, "1.5M"},
		{2000000, "2M"},
	}

	for _, tt := range tests {
		result := FormatTokens(tt.input)
		if result != tt.expected {
			t.Errorf("FormatTokens(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func contextValues(result *types.ProviderResult) map[string]any {
	return result.Values["context"].(map[string]any)
}

func contextPercent(result *types.ProviderResult) map[string]any {
	return contextValues(result)["percent"].(map[string]any)
}

func TestContextProvider(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage:    36,
			ContextWindowSize: 1000000,
			CurrentUsage: &types.CurrentUsage{
				InputTokens:              100,
				CacheCreationInputTokens: 200,
				CacheReadInputTokens:     300,
			},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	ctx := contextValues(result)
	if ctx["tokens"] != "600" {
		t.Errorf("expected 600 tokens, got %s", ctx["tokens"])
	}
	if ctx["size"] != "1M" {
		t.Errorf("expected 1M size, got %s", ctx["size"])
	}
	pct := contextPercent(result)
	if pct["used"] != 36 {
		t.Errorf("expected 36%%, got %v", pct["used"])
	}
}

func TestContextProviderWithTotalTokens(t *testing.T) {
	p := &contextProvider{}
	inputTokens := 50000
	outputTokens := 8000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage:    36,
			ContextWindowSize: 1000000,
			TotalInputTokens:  &inputTokens,
			TotalOutputTokens: &outputTokens,
			CurrentUsage: &types.CurrentUsage{
				InputTokens:              100,
				CacheCreationInputTokens: 200,
				CacheReadInputTokens:     300,
			},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	ctx := contextValues(result)
	if ctx["input"] != "50K" {
		t.Errorf("expected Input 50K, got %s", ctx["input"])
	}
	if ctx["output"] != "8K" {
		t.Errorf("expected Output 8K, got %s", ctx["output"])
	}
}

func TestContextProviderInputFallback(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage: 10,
			CurrentUsage: &types.CurrentUsage{
				InputTokens:              100,
				CacheCreationInputTokens: 200,
				CacheReadInputTokens:     300,
			},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	ctx := contextValues(result)
	if ctx["input"] != "600" {
		t.Errorf("expected Input 600, got %s", ctx["input"])
	}
	if ctx["output"] != "" {
		t.Errorf("expected empty Output, got %s", ctx["output"])
	}
}

func TestContextProviderRemaining(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage:      36,
			RemainingPercentage: 64,
			ContextWindowSize:   1000000,
			CurrentUsage:        &types.CurrentUsage{InputTokens: 100},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	pct := contextPercent(result)
	if pct["remaining"] != 64 {
		t.Errorf("expected remaining 64, got %v", pct["remaining"])
	}
}

func TestContextProviderNoRemaining(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	pct := contextPercent(result)
	if pct["remaining"] != 0 {
		t.Errorf("expected remaining 0, got %v", pct["remaining"])
	}
}

func TestContextProviderZeroRemaining(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD:           "/tmp",
		ContextWindow: &types.ContextWindow{},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	pct := contextPercent(result)
	if pct["remaining"] != 0 {
		t.Errorf("expected remaining 0 for zero value with no usage, got %v", pct["remaining"])
	}
}
