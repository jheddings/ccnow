package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

type speedProvider struct{}

func (p *speedProvider) Name() string { return "speed" }

func (p *speedProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	result := &types.ProviderResult{
		Values: map[string]any{
			"speed.input":  "",
			"speed.output": "",
			"speed.total":  "",
		},
	}

	cw := session.ContextWindow
	cost := session.Cost
	if cw == nil || cost == nil || cost.TotalAPIDurationMS == 0 {
		return result, nil
	}

	durationSec := cost.TotalAPIDurationMS / 1000.0

	if cw.TotalInputTokens != nil {
		speed := float64(*cw.TotalInputTokens) / durationSec
		result.Values["speed.input"] = FormatSpeed(speed)
	}

	if cw.TotalOutputTokens != nil {
		speed := float64(*cw.TotalOutputTokens) / durationSec
		result.Values["speed.output"] = FormatSpeed(speed)
	}

	if cw.TotalInputTokens != nil && cw.TotalOutputTokens != nil {
		speed := float64(*cw.TotalInputTokens+*cw.TotalOutputTokens) / durationSec
		result.Values["speed.total"] = FormatSpeed(speed)
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
