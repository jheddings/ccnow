package segment

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/provider"
	"github.com/jheddings/ccglow/internal/types"
)

// RegisterBuiltin adds all built-in segment implementations to the registry.
func RegisterBuiltin(registry *Registry) {
	registry.Register(&literalSegment{})
	registry.Register(&newlineSegment{})
	registry.Register(&pwdNameSegment{})
	registry.Register(&pwdPathSegment{})
	registry.Register(&pwdSmartSegment{})
	registry.Register(&gitBranchSegment{})
	registry.Register(&gitInsertionsSegment{})
	registry.Register(&gitDeletionsSegment{})
	registry.Register(&gitModifiedSegment{})
	registry.Register(&gitStagedSegment{})
	registry.Register(&gitUntrackedSegment{})
	registry.Register(&gitOwnerSegment{})
	registry.Register(&gitRepoSegment{})
	registry.Register(&gitWorktreeSegment{})
	registry.Register(&contextTokensSegment{})
	registry.Register(&contextSizeSegment{})
	registry.Register(&contextPercentSegment{})
	registry.Register(&contextInputSegment{})
	registry.Register(&contextOutputSegment{})
	registry.Register(&contextRemainingSegment{})
	registry.Register(&modelNameSegment{})
	registry.Register(&modelIDSegment{})
	registry.Register(&costUSDSegment{})
	registry.Register(&speedInputSegment{})
	registry.Register(&speedOutputSegment{})
	registry.Register(&speedTotalSegment{})
	registry.Register(&sessionDurationSegment{})
	registry.Register(&sessionLinesAddedSegment{})
	registry.Register(&sessionLinesRemovedSegment{})
	registry.Register(&claudeVersionSegment{})
	registry.Register(&claudeStyleSegment{})
}

// Registry maps segment type names to their implementations.
type Registry struct {
	segments map[string]types.Segment
}

// NewRegistry creates an empty segment registry.
func NewRegistry() *Registry {
	return &Registry{segments: make(map[string]types.Segment)}
}

// Register adds a segment implementation.
func (r *Registry) Register(seg types.Segment) {
	r.segments[seg.Name()] = seg
}

// Get returns the segment for the given type name, or nil.
func (r *Registry) Get(name string) types.Segment {
	return r.segments[name]
}

// --- Literal ---

type literalSegment struct{}

func (s *literalSegment) Name() string { return "literal" }
func (s *literalSegment) Render(ctx *types.SegmentContext) *string {
	if ctx.Props == nil {
		return nil
	}
	if text, ok := ctx.Props["text"].(string); ok {
		return &text
	}
	return nil
}

// --- Newline ---

type newlineSegment struct{}

func (s *newlineSegment) Name() string { return "newline" }
func (s *newlineSegment) Render(ctx *types.SegmentContext) *string {
	v := "\n"
	return &v
}

// --- PWD ---

type pwdNameSegment struct{}

func (s *pwdNameSegment) Name() string { return "pwd.name" }
func (s *pwdNameSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.PwdData); ok && data != nil {
		return &data.Name
	}
	return nil
}

type pwdPathSegment struct{}

func (s *pwdPathSegment) Name() string { return "pwd.path" }
func (s *pwdPathSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.PwdData); ok && data != nil && data.Path != "" {
		return &data.Path
	}
	return nil
}

type pwdSmartSegment struct{}

func (s *pwdSmartSegment) Name() string { return "pwd.smart" }
func (s *pwdSmartSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.PwdData); ok && data != nil && data.Smart != "" {
		return &data.Smart
	}
	return nil
}

// --- Git ---

type gitBranchSegment struct{}

func (s *gitBranchSegment) Name() string { return "git.branch" }
func (s *gitBranchSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil {
		return data.Branch
	}
	return nil
}

type gitInsertionsSegment struct{}

func (s *gitInsertionsSegment) Name() string { return "git.insertions" }
func (s *gitInsertionsSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil && data.Insertions != nil {
		v := fmt.Sprintf("%d", *data.Insertions)
		return &v
	}
	return nil
}

type gitDeletionsSegment struct{}

func (s *gitDeletionsSegment) Name() string { return "git.deletions" }
func (s *gitDeletionsSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil && data.Deletions != nil {
		v := fmt.Sprintf("%d", *data.Deletions)
		return &v
	}
	return nil
}

type gitModifiedSegment struct{}

func (s *gitModifiedSegment) Name() string { return "git.modified" }
func (s *gitModifiedSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil && data.Modified != nil {
		v := fmt.Sprintf("%d", *data.Modified)
		return &v
	}
	return nil
}

type gitStagedSegment struct{}

func (s *gitStagedSegment) Name() string { return "git.staged" }
func (s *gitStagedSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil && data.Staged != nil {
		v := fmt.Sprintf("%d", *data.Staged)
		return &v
	}
	return nil
}

type gitUntrackedSegment struct{}

func (s *gitUntrackedSegment) Name() string { return "git.untracked" }
func (s *gitUntrackedSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil && data.Untracked != nil {
		v := fmt.Sprintf("%d", *data.Untracked)
		return &v
	}
	return nil
}

type gitOwnerSegment struct{}

func (s *gitOwnerSegment) Name() string { return "git.owner" }
func (s *gitOwnerSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil {
		return data.Owner
	}
	return nil
}

type gitRepoSegment struct{}

func (s *gitRepoSegment) Name() string { return "git.repo" }
func (s *gitRepoSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil {
		return data.Repo
	}
	return nil
}

type gitWorktreeSegment struct{}

func (s *gitWorktreeSegment) Name() string { return "git.worktree" }
func (s *gitWorktreeSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil {
		return data.Worktree
	}
	return nil
}

// --- Context ---

type contextTokensSegment struct{}

func (s *contextTokensSegment) Name() string { return "context.tokens" }
func (s *contextTokensSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ContextData); ok && data != nil && data.Tokens != "" {
		return &data.Tokens
	}
	return nil
}

type contextSizeSegment struct{}

func (s *contextSizeSegment) Name() string { return "context.size" }
func (s *contextSizeSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ContextData); ok && data != nil && data.Size != "" {
		return &data.Size
	}
	return nil
}

type contextPercentSegment struct{}

func (s *contextPercentSegment) Name() string { return "context.percent" }
func (s *contextPercentSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ContextData); ok && data != nil && data.Percent != nil {
		v := fmt.Sprintf("%d%%", *data.Percent)
		return &v
	}
	return nil
}

type contextInputSegment struct{}

func (s *contextInputSegment) Name() string { return "context.input" }
func (s *contextInputSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ContextData); ok && data != nil && data.Input != "" {
		return &data.Input
	}
	return nil
}

type contextOutputSegment struct{}

func (s *contextOutputSegment) Name() string { return "context.output" }
func (s *contextOutputSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ContextData); ok && data != nil && data.Output != "" {
		return &data.Output
	}
	return nil
}

type contextRemainingSegment struct{}

func (s *contextRemainingSegment) Name() string { return "context.remaining" }
func (s *contextRemainingSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ContextData); ok && data != nil && data.Remaining != nil {
		v := fmt.Sprintf("%d%%", *data.Remaining)
		return &v
	}
	return nil
}

// --- Model ---

type modelNameSegment struct{}

func (s *modelNameSegment) Name() string { return "model.name" }
func (s *modelNameSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ModelData); ok && data != nil {
		return data.Name
	}
	return nil
}

type modelIDSegment struct{}

func (s *modelIDSegment) Name() string { return "model.id" }
func (s *modelIDSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ModelData); ok && data != nil {
		return data.ID
	}
	return nil
}

// --- Cost ---

type costUSDSegment struct{}

func (s *costUSDSegment) Name() string { return "cost.usd" }
func (s *costUSDSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.CostData); ok && data != nil {
		return data.USD
	}
	return nil
}

// --- Speed ---

type speedInputSegment struct{}

func (s *speedInputSegment) Name() string { return "speed.input" }
func (s *speedInputSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.SpeedData); ok && data != nil {
		return data.Input
	}
	return nil
}

type speedOutputSegment struct{}

func (s *speedOutputSegment) Name() string { return "speed.output" }
func (s *speedOutputSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.SpeedData); ok && data != nil {
		return data.Output
	}
	return nil
}

type speedTotalSegment struct{}

func (s *speedTotalSegment) Name() string { return "speed.total" }
func (s *speedTotalSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.SpeedData); ok && data != nil {
		return data.Total
	}
	return nil
}

// --- Session ---

type sessionDurationSegment struct{}

func (s *sessionDurationSegment) Name() string { return "session.duration" }
func (s *sessionDurationSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.SessionData); ok && data != nil {
		return data.Duration
	}
	return nil
}

type sessionLinesAddedSegment struct{}

func (s *sessionLinesAddedSegment) Name() string { return "session.lines-added" }
func (s *sessionLinesAddedSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.SessionData); ok && data != nil && data.LinesAdded != nil {
		v := fmt.Sprintf("%d", *data.LinesAdded)
		return &v
	}
	return nil
}

type sessionLinesRemovedSegment struct{}

func (s *sessionLinesRemovedSegment) Name() string { return "session.lines-removed" }
func (s *sessionLinesRemovedSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.SessionData); ok && data != nil && data.LinesRemoved != nil {
		v := fmt.Sprintf("%d", *data.LinesRemoved)
		return &v
	}
	return nil
}

// --- Claude ---

type claudeVersionSegment struct{}

func (s *claudeVersionSegment) Name() string { return "claude.version" }
func (s *claudeVersionSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ClaudeData); ok && data != nil {
		return data.Version
	}
	return nil
}

type claudeStyleSegment struct{}

func (s *claudeStyleSegment) Name() string { return "claude.style" }
func (s *claudeStyleSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.ClaudeData); ok && data != nil {
		return data.Style
	}
	return nil
}
