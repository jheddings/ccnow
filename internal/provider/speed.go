package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

// SpeedData holds resolved token speed information.
type SpeedData struct {
	Input  *string `segment:"speed.input"`
	Output *string `segment:"speed.output"`
	Total  *string `segment:"speed.total"`
}

func (p *speedProvider) Fields() any { return &SpeedData{} }

type speedProvider struct{}

func (p *speedProvider) Name() string { return "speed" }

func (p *speedProvider) Resolve(session *types.SessionData) (any, error) {
	data := &SpeedData{}

	cw := session.ContextWindow
	cost := session.Cost
	if cw == nil || cost == nil || cost.TotalAPIDurationMS == 0 {
		return data, nil
	}

	durationSec := cost.TotalAPIDurationMS / 1000.0

	if cw.TotalInputTokens != nil {
		speed := float64(*cw.TotalInputTokens) / durationSec
		s := FormatSpeed(speed)
		data.Input = &s
	}

	if cw.TotalOutputTokens != nil {
		speed := float64(*cw.TotalOutputTokens) / durationSec
		s := FormatSpeed(speed)
		data.Output = &s
	}

	if cw.TotalInputTokens != nil && cw.TotalOutputTokens != nil {
		speed := float64(*cw.TotalInputTokens+*cw.TotalOutputTokens) / durationSec
		s := FormatSpeed(speed)
		data.Total = &s
	}

	return data, nil
}

// FormatSpeed formats a tokens-per-second value for display.
func FormatSpeed(tokensPerSec float64) string {
	if tokensPerSec >= 1000 {
		return fmt.Sprintf("%.1fK t/s", tokensPerSec/1000)
	}
	return fmt.Sprintf("%d t/s", int(tokensPerSec))
}
