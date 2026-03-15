package segment

import (
	"testing"

	"github.com/jheddings/ccglow/internal/provider"
	"github.com/jheddings/ccglow/internal/types"
)

func intPtr(n int) *int { return &n }
func strPtr(s string) *string { return &s }

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

func TestClaudeVersionSegment(t *testing.T) {
	seg := &claudeVersionSegment{}
	if seg.Name() != "claude.version" {
		t.Errorf("expected name claude.version, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ClaudeData{Version: strPtr("2.1.75")}}
	result := seg.Render(ctx)
	if result == nil || *result != "2.1.75" {
		t.Errorf("expected '2.1.75', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ClaudeData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestClaudeStyleSegment(t *testing.T) {
	seg := &claudeStyleSegment{}
	if seg.Name() != "claude.style" {
		t.Errorf("expected name claude.style, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ClaudeData{Style: strPtr("concise")}}
	result := seg.Render(ctx)
	if result == nil || *result != "concise" {
		t.Errorf("expected 'concise', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ClaudeData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestModelIDSegment(t *testing.T) {
	seg := &modelIDSegment{}
	if seg.Name() != "model.id" {
		t.Errorf("expected name model.id, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ModelData{ID: strPtr("claude-opus-4-6[1m]")}}
	result := seg.Render(ctx)
	if result == nil || *result != "claude-opus-4-6[1m]" {
		t.Errorf("expected 'claude-opus-4-6[1m]', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ModelData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestSessionDurationTotalSegment(t *testing.T) {
	seg := &sessionDurationTotalSegment{}
	if seg.Name() != "session.duration.total" {
		t.Errorf("expected name session.duration.total, got %s", seg.Name())
	}

	dur := "1h 30m"
	ctx := &types.SegmentContext{Provider: &provider.SessionData{Duration: &dur}}
	result := seg.Render(ctx)
	if result == nil || *result != "1h 30m" {
		t.Errorf("expected '1h 30m', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.SessionData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestSessionDurationAPISegment(t *testing.T) {
	seg := &sessionDurationAPISegment{}
	if seg.Name() != "session.duration.api" {
		t.Errorf("expected name session.duration.api, got %s", seg.Name())
	}

	dur := "8m"
	ctx := &types.SegmentContext{Provider: &provider.SessionData{APIDuration: &dur}}
	result := seg.Render(ctx)
	if result == nil || *result != "8m" {
		t.Errorf("expected '8m', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.SessionData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestSessionIDSegment(t *testing.T) {
	seg := &sessionIDSegment{}
	if seg.Name() != "session.id" {
		t.Errorf("expected name session.id, got %s", seg.Name())
	}

	id := "abc-123"
	ctx := &types.SegmentContext{Provider: &provider.SessionData{ID: &id}}
	result := seg.Render(ctx)
	if result == nil || *result != "abc-123" {
		t.Errorf("expected 'abc-123', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.SessionData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestContextRemainingSegment(t *testing.T) {
	seg := &contextRemainingSegment{}
	if seg.Name() != "context.remaining" {
		t.Errorf("expected name context.remaining, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.ContextData{Remaining: intPtr(96)}}
	result := seg.Render(ctx)
	if result == nil || *result != "96%" {
		t.Errorf("expected '96%%', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.ContextData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}
