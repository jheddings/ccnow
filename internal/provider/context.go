package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

type contextProvider struct{}

func (p *contextProvider) Name() string { return "context" }

func (p *contextProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	result := &types.ProviderResult{
		Values: map[string]any{
			"context.tokens":            "",
			"context.size":              "",
			"context.percent.used":      0,
			"context.percent.remaining": 0,
			"context.input":             "",
			"context.output":            "",
		},
		Formats: map[string]string{
			"context.percent.used":      "%d%%",
			"context.percent.remaining": "%d%%",
		},
	}

	cw := session.ContextWindow
	if cw == nil {
		return result, nil
	}

	totalTokens := 0
	if cw.CurrentUsage != nil {
		totalTokens = cw.CurrentUsage.InputTokens +
			cw.CurrentUsage.CacheCreationInputTokens +
			cw.CurrentUsage.CacheReadInputTokens
	}

	result.Values["context.tokens"] = FormatTokens(totalTokens)

	if cw.ContextWindowSize > 0 {
		result.Values["context.size"] = FormatTokens(cw.ContextWindowSize)
	}

	if cw.UsedPercentage > 0 || cw.CurrentUsage != nil {
		result.Values["context.percent.used"] = cw.UsedPercentage
	}

	if cw.RemainingPercentage > 0 || cw.CurrentUsage != nil {
		result.Values["context.percent.remaining"] = cw.RemainingPercentage
	}

	if cw.TotalInputTokens != nil {
		result.Values["context.input"] = FormatTokens(*cw.TotalInputTokens)
	} else if totalTokens > 0 {
		result.Values["context.input"] = FormatTokens(totalTokens)
	}

	if cw.TotalOutputTokens != nil {
		result.Values["context.output"] = FormatTokens(*cw.TotalOutputTokens)
	}

	return result, nil
}

// FormatTokens formats a token count for display (e.g. 1500000 -> "1.5M").
func FormatTokens(total int) string {
	if total >= 1_000_000 {
		m := float64(total) / 1_000_000.0
		if m == float64(int(m)) {
			return fmt.Sprintf("%dM", int(m))
		}
		return fmt.Sprintf("%.1fM", m)
	}
	if total >= 1_000 {
		return fmt.Sprintf("%dK", total/1_000)
	}
	return fmt.Sprintf("%d", total)
}
