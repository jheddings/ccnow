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

func TestGitOwnerSegment(t *testing.T) {
	seg := &gitOwnerSegment{}
	if seg.Name() != "git.owner" {
		t.Errorf("expected name git.owner, got %s", seg.Name())
	}

	owner := "jheddings"
	ctx := &types.SegmentContext{Provider: &provider.GitData{Owner: &owner}}
	result := seg.Render(ctx)
	if result == nil || *result != "jheddings" {
		t.Errorf("expected 'jheddings', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.GitData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestGitRepoSegment(t *testing.T) {
	seg := &gitRepoSegment{}
	if seg.Name() != "git.repo" {
		t.Errorf("expected name git.repo, got %s", seg.Name())
	}

	repo := "ccglow"
	ctx := &types.SegmentContext{Provider: &provider.GitData{Repo: &repo}}
	result := seg.Render(ctx)
	if result == nil || *result != "ccglow" {
		t.Errorf("expected 'ccglow', got %v", result)
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

func TestContextInputSegment(t *testing.T) {
	seg := &contextInputSegment{}
	if seg.Name() != "context.input" {
		t.Errorf("expected name context.input, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ContextData{Input: "50K"}}
	result := seg.Render(ctx)
	if result == nil || *result != "50K" {
		t.Errorf("expected '50K', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ContextData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil for empty Input, got %v", result)
	}
}

func TestContextOutputSegment(t *testing.T) {
	seg := &contextOutputSegment{}
	if seg.Name() != "context.output" {
		t.Errorf("expected name context.output, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ContextData{Output: "8K"}}
	result := seg.Render(ctx)
	if result == nil || *result != "8K" {
		t.Errorf("expected '8K', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ContextData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil for empty Output, got %v", result)
	}
}

func TestSpeedInputSegment(t *testing.T) {
	seg := &speedInputSegment{}
	if seg.Name() != "speed.input" {
		t.Errorf("expected name speed.input, got %s", seg.Name())
	}

	v := "2K t/s"
	ctx := &types.SegmentContext{Provider: &provider.SpeedData{Input: &v}}
	result := seg.Render(ctx)
	if result == nil || *result != "2K t/s" {
		t.Errorf("expected '2K t/s', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.SpeedData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestSpeedOutputSegment(t *testing.T) {
	seg := &speedOutputSegment{}
	if seg.Name() != "speed.output" {
		t.Errorf("expected name speed.output, got %s", seg.Name())
	}

	v := "1K t/s"
	ctx := &types.SegmentContext{Provider: &provider.SpeedData{Output: &v}}
	result := seg.Render(ctx)
	if result == nil || *result != "1K t/s" {
		t.Errorf("expected '1K t/s', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.SpeedData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestSpeedTotalSegment(t *testing.T) {
	seg := &speedTotalSegment{}
	if seg.Name() != "speed.total" {
		t.Errorf("expected name speed.total, got %s", seg.Name())
	}

	v := "3K t/s"
	ctx := &types.SegmentContext{Provider: &provider.SpeedData{Total: &v}}
	result := seg.Render(ctx)
	if result == nil || *result != "3K t/s" {
		t.Errorf("expected '3K t/s', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.SpeedData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}
