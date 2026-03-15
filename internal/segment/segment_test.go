package segment

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestNewlineSegment(t *testing.T) {
	seg := &newlineSegment{}

	if seg.Name() != "newline" {
		t.Errorf("expected name newline, got %s", seg.Name())
	}

	result := seg.Render(&types.SegmentContext{})
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if *result != "\n" {
		t.Errorf("expected newline character, got %q", *result)
	}
}
