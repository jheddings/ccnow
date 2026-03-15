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

	data := result.(*SessionData)
	if *data.Duration != "1h 30m" {
		t.Errorf("expected 1h 30m, got %s", *data.Duration)
	}
	if *data.APIDuration != "8m" {
		t.Errorf("expected 8m, got %s", *data.APIDuration)
	}
	if *data.LinesAdded != 100 {
		t.Errorf("expected 100 lines added, got %d", *data.LinesAdded)
	}
	if *data.LinesRemoved != 50 {
		t.Errorf("expected 50 lines removed, got %d", *data.LinesRemoved)
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

	data := result.(*SessionData)
	if data.ID == nil || *data.ID != "abc-123" {
		t.Errorf("expected abc-123, got %v", data.ID)
	}
}

func TestSessionProviderNoID(t *testing.T) {
	p := &sessionProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*SessionData)
	if data.ID != nil {
		t.Errorf("expected nil ID, got %v", data.ID)
	}
}
