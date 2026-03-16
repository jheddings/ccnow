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

	if result.Values["session.duration.total"] != "1h 30m" {
		t.Errorf("expected 1h 30m, got %s", result.Values["session.duration.total"])
	}
	if result.Values["session.duration.api"] != "8m" {
		t.Errorf("expected 8m, got %s", result.Values["session.duration.api"])
	}
	if result.Values["session.lines-added"] != 100 {
		t.Errorf("expected 100 lines added, got %v", result.Values["session.lines-added"])
	}
	if result.Values["session.lines-removed"] != 50 {
		t.Errorf("expected 50 lines removed, got %v", result.Values["session.lines-removed"])
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

	if result.Values["session.id"] != "abc-123" {
		t.Errorf("expected abc-123, got %v", result.Values["session.id"])
	}
}

func TestSessionProviderNoID(t *testing.T) {
	p := &sessionProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["session.id"] != "" {
		t.Errorf("expected empty ID, got %v", result.Values["session.id"])
	}
}
