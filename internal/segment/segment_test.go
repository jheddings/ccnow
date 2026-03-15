package segment

import (
	"testing"

	"github.com/jheddings/ccglow/internal/provider"
	"github.com/jheddings/ccglow/internal/types"
)

func intPtr(n int) *int { return &n }

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

func TestGitModifiedSegment(t *testing.T) {
	seg := &gitModifiedSegment{}
	if seg.Name() != "git.modified" {
		t.Errorf("expected name git.modified, got %s", seg.Name())
	}

	// With data
	ctx := &types.SegmentContext{Provider: &provider.GitData{Modified: intPtr(3)}}
	result := seg.Render(ctx)
	if result == nil || *result != "3" {
		t.Errorf("expected '3', got %v", result)
	}

	// Nil field
	ctx = &types.SegmentContext{Provider: &provider.GitData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestGitStagedSegment(t *testing.T) {
	seg := &gitStagedSegment{}
	if seg.Name() != "git.staged" {
		t.Errorf("expected name git.staged, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.GitData{Staged: intPtr(5)}}
	result := seg.Render(ctx)
	if result == nil || *result != "5" {
		t.Errorf("expected '5', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.GitData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestGitUntrackedSegment(t *testing.T) {
	seg := &gitUntrackedSegment{}
	if seg.Name() != "git.untracked" {
		t.Errorf("expected name git.untracked, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.GitData{Untracked: intPtr(2)}}
	result := seg.Render(ctx)
	if result == nil || *result != "2" {
		t.Errorf("expected '2', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.GitData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestGitWorktreeSegment(t *testing.T) {
	seg := &gitWorktreeSegment{}
	if seg.Name() != "git.worktree" {
		t.Errorf("expected name git.worktree, got %s", seg.Name())
	}

	name := "my-worktree"
	ctx := &types.SegmentContext{Provider: &provider.GitData{Worktree: &name}}
	result := seg.Render(ctx)
	if result == nil || *result != "my-worktree" {
		t.Errorf("expected 'my-worktree', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.GitData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}
