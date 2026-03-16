package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestSpeedProvider(t *testing.T) {
	p := &speedProvider{}
	inputTokens := 10000
	outputTokens := 5000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			TotalInputTokens:  &inputTokens,
			TotalOutputTokens: &outputTokens,
		},
		Cost: &types.CostInfo{
			TotalAPIDurationMS: 5000,
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["speed.input"] != "2.0K t/s" {
		t.Errorf("expected Input '2.0K t/s', got %v", result.Values["speed.input"])
	}
	if result.Values["speed.output"] != "1.0K t/s" {
		t.Errorf("expected Output '1.0K t/s', got %v", result.Values["speed.output"])
	}
	if result.Values["speed.total"] != "3.0K t/s" {
		t.Errorf("expected Total '3.0K t/s', got %v", result.Values["speed.total"])
	}
}

func TestSpeedProviderZeroDuration(t *testing.T) {
	p := &speedProvider{}
	inputTokens := 10000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			TotalInputTokens: &inputTokens,
		},
		Cost: &types.CostInfo{TotalAPIDurationMS: 0},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["speed.input"] != "" {
		t.Errorf("expected empty Input for zero duration, got %v", result.Values["speed.input"])
	}
}

func TestSpeedProviderNilContextWindow(t *testing.T) {
	p := &speedProvider{}
	sess := &types.SessionData{
		CWD:  "/tmp",
		Cost: &types.CostInfo{TotalAPIDurationMS: 5000},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["speed.input"] != "" || result.Values["speed.output"] != "" || result.Values["speed.total"] != "" {
		t.Error("expected all empty fields when ContextWindow is nil")
	}
}

func TestSpeedProviderNilCost(t *testing.T) {
	p := &speedProvider{}
	inputTokens := 10000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			TotalInputTokens: &inputTokens,
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["speed.input"] != "" || result.Values["speed.output"] != "" || result.Values["speed.total"] != "" {
		t.Error("expected all empty fields when Cost is nil")
	}
}

func TestSpeedProviderPartialTokens(t *testing.T) {
	p := &speedProvider{}
	inputTokens := 10000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			TotalInputTokens: &inputTokens,
		},
		Cost: &types.CostInfo{TotalAPIDurationMS: 2000},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["speed.input"] != "5.0K t/s" {
		t.Errorf("expected Input '5.0K t/s', got %v", result.Values["speed.input"])
	}
	if result.Values["speed.output"] != "" {
		t.Errorf("expected empty Output, got %v", result.Values["speed.output"])
	}
	if result.Values["speed.total"] != "" {
		t.Errorf("expected empty Total when output missing, got %v", result.Values["speed.total"])
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0 t/s"},
		{42, "42 t/s"},
		{999, "999 t/s"},
		{1000, "1.0K t/s"},
		{1500, "1.5K t/s"},
		{2000, "2.0K t/s"},
		{10500, "10.5K t/s"},
	}

	for _, tt := range tests {
		result := FormatSpeed(tt.input)
		if result != tt.expected {
			t.Errorf("FormatSpeed(%g) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
