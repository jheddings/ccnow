package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

type contextProvider struct{}

func (p *contextProvider) Name() string { return "context" }

func (p *contextProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	ctx := map[string]any{
		"tokens": "",
		"size":   "",
		"percent": map[string]any{
			"used":      0,
			"remaining": 0,
		},
		"input":  "",
		"output": "",
	}

	result := &types.ProviderResult{
		Values: map[string]any{"context": ctx},
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

	ctx["tokens"] = FormatTokens(totalTokens)

	if cw.ContextWindowSize > 0 {
		ctx["size"] = FormatTokens(cw.ContextWindowSize)
	}

	pct := ctx["percent"].(map[string]any)
	if cw.UsedPercentage > 0 || cw.CurrentUsage != nil {
		pct["used"] = cw.UsedPercentage
	}

	if cw.RemainingPercentage > 0 || cw.CurrentUsage != nil {
		pct["remaining"] = cw.RemainingPercentage
	}

	if cw.TotalInputTokens != nil {
		ctx["input"] = FormatTokens(*cw.TotalInputTokens)
	} else if totalTokens > 0 {
		ctx["input"] = FormatTokens(totalTokens)
	}

	if cw.TotalOutputTokens != nil {
		ctx["output"] = FormatTokens(*cw.TotalOutputTokens)
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
