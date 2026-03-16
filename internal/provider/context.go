package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

// ContextData holds resolved context window information.
type ContextData struct {
	Tokens    string `segment:"context.tokens"`
	Size      string `segment:"context.size"`
	Percent   *int   `segment:"context.percent.used,format:%d%%"`
	Remaining *int   `segment:"context.percent.remaining,format:%d%%"`
	Input     string `segment:"context.input"`
	Output    string `segment:"context.output"`
}

func (p *contextProvider) Fields() any { return &ContextData{} }

type contextProvider struct{}

func (p *contextProvider) Name() string { return "context" }

func (p *contextProvider) Resolve(session *types.SessionData) (any, error) {
	cw := session.ContextWindow
	if cw == nil {
		return &ContextData{}, nil
	}

	data := &ContextData{}

	totalTokens := 0
	if cw.CurrentUsage != nil {
		totalTokens = cw.CurrentUsage.InputTokens +
			cw.CurrentUsage.CacheCreationInputTokens +
			cw.CurrentUsage.CacheReadInputTokens
	}

	data.Tokens = FormatTokens(totalTokens)

	if cw.ContextWindowSize > 0 {
		data.Size = FormatTokens(cw.ContextWindowSize)
	}

	if cw.UsedPercentage > 0 || cw.CurrentUsage != nil {
		pct := cw.UsedPercentage
		data.Percent = &pct
	}

	if cw.RemainingPercentage > 0 || cw.CurrentUsage != nil {
		rem := cw.RemainingPercentage
		data.Remaining = &rem
	}

	if cw.TotalInputTokens != nil {
		data.Input = FormatTokens(*cw.TotalInputTokens)
	} else if totalTokens > 0 {
		data.Input = FormatTokens(totalTokens)
	}

	if cw.TotalOutputTokens != nil {
		data.Output = FormatTokens(*cw.TotalOutputTokens)
	}

	return data, nil
}

// FormatTokens formats a token count for display (e.g. 1500000 → "1.5M").
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
