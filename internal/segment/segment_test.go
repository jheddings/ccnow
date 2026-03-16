package segment

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestLiteralSegment(t *testing.T) {
	seg := &literalSegment{}
	if seg.Name() != "literal" {
		t.Errorf("expected name literal, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Props: map[string]any{"text": "hello"}}
	result := seg.Render(ctx)
	if result == nil || *result != "hello" {
		t.Errorf("expected 'hello', got %v", result)
	}

	ctx = &types.SegmentContext{}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestNewlineSegment(t *testing.T) {
	seg := &newlineSegment{}
	if seg.Name() != "newline" {
		t.Errorf("expected name newline, got %s", seg.Name())
	}

	result := seg.Render(&types.SegmentContext{})
	if result == nil || *result != "\n" {
		t.Errorf("expected newline, got %v", result)
	}
}
