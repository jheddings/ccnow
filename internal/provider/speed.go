package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

type speedProvider struct{}

func (p *speedProvider) Name() string { return "speed" }

func (p *speedProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	speed := map[string]any{
		"input":  "",
		"output": "",
		"total":  "",
	}

	result := &types.ProviderResult{
		Values: map[string]any{"speed": speed},
	}

	cw := session.ContextWindow
	cost := session.Cost
	if cw == nil || cost == nil || cost.TotalAPIDurationMS == 0 {
		return result, nil
	}

	durationSec := cost.TotalAPIDurationMS / 1000.0

	if cw.TotalInputTokens != nil {
		s := float64(*cw.TotalInputTokens) / durationSec
		speed["input"] = FormatSpeed(s)
	}

	if cw.TotalOutputTokens != nil {
		s := float64(*cw.TotalOutputTokens) / durationSec
		speed["output"] = FormatSpeed(s)
	}

	if cw.TotalInputTokens != nil && cw.TotalOutputTokens != nil {
		s := float64(*cw.TotalInputTokens+*cw.TotalOutputTokens) / durationSec
		speed["total"] = FormatSpeed(s)
	}

	return result, nil
}

// FormatSpeed formats a tokens-per-second value for display.
func FormatSpeed(tokensPerSec float64) string {
	if tokensPerSec >= 1000 {
		return fmt.Sprintf("%.1fK t/s", tokensPerSec/1000)
	}
	return fmt.Sprintf("%d t/s", int(tokensPerSec))
}
