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

	data := result.(*SpeedData)
	if data.Input == nil || *data.Input != "2.0K t/s" {
		t.Errorf("expected Input '2.0K t/s', got %v", data.Input)
	}
	if data.Output == nil || *data.Output != "1.0K t/s" {
		t.Errorf("expected Output '1.0K t/s', got %v", data.Output)
	}
	if data.Total == nil || *data.Total != "3.0K t/s" {
		t.Errorf("expected Total '3.0K t/s', got %v", data.Total)
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

	data := result.(*SpeedData)
	if data.Input != nil {
		t.Errorf("expected nil Input for zero duration, got %v", *data.Input)
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

	data := result.(*SpeedData)
	if data.Input != nil || data.Output != nil || data.Total != nil {
		t.Error("expected all nil fields when ContextWindow is nil")
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

	data := result.(*SpeedData)
	if data.Input != nil || data.Output != nil || data.Total != nil {
		t.Error("expected all nil fields when Cost is nil")
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

	data := result.(*SpeedData)
	if data.Input == nil || *data.Input != "5.0K t/s" {
		t.Errorf("expected Input '5.0K t/s', got %v", data.Input)
	}
	if data.Output != nil {
		t.Errorf("expected nil Output, got %v", *data.Output)
	}
	if data.Total != nil {
		t.Errorf("expected nil Total when output missing, got %v", *data.Total)
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
