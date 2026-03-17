package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		ms       float64
		expected string
	}{
		{0, "0m"},
		{30000, "0m"},
		{60000, "1m"},
		{300000, "5m"},
		{3600000, "1h 0m"},
		{5400000, "1h 30m"},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.ms)
		if result != tt.expected {
			t.Errorf("FormatDuration(%f) = %q, want %q", tt.ms, result, tt.expected)
		}
	}
}

func sessionValues(result *types.ProviderResult) map[string]any {
	return result.Values["session"].(map[string]any)
}

func sessionDuration(result *types.ProviderResult) map[string]any {
	return sessionValues(result)["duration"].(map[string]any)
}

func TestSessionProvider(t *testing.T) {
	p := &sessionProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		Cost: &types.CostInfo{
			TotalDurationMS:    5400000,
			TotalAPIDurationMS: 522771,
			TotalLinesAdded:    100,
			TotalLinesRemoved:  50,
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	dur := sessionDuration(result)
	if dur["total"] != "1h 30m" {
		t.Errorf("expected 1h 30m, got %s", dur["total"])
	}
	if dur["api"] != "8m" {
		t.Errorf("expected 8m, got %s", dur["api"])
	}

	s := sessionValues(result)
	if s["lines-added"] != 100 {
		t.Errorf("expected 100 lines added, got %v", s["lines-added"])
	}
	if s["lines-removed"] != 50 {
		t.Errorf("expected 50 lines removed, got %v", s["lines-removed"])
	}
}

func TestSessionProviderID(t *testing.T) {
	p := &sessionProvider{}
	sess := &types.SessionData{
		CWD:       "/tmp",
		SessionID: "abc-123",
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	s := sessionValues(result)
	if s["id"] != "abc-123" {
		t.Errorf("expected abc-123, got %v", s["id"])
	}
}

func TestSessionProviderNoID(t *testing.T) {
	p := &sessionProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	s := sessionValues(result)
	if s["id"] != "" {
		t.Errorf("expected empty ID, got %v", s["id"])
	}
}
