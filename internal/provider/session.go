package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

type sessionProvider struct{}

func (p *sessionProvider) Name() string { return "session" }

func (p *sessionProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	result := &types.ProviderResult{
		Values: map[string]any{
			"session.duration.total": "",
			"session.duration.api":   "",
			"session.lines-added":    0,
			"session.lines-removed":  0,
			"session.id":             "",
		},
	}

	if session.SessionID != "" {
		result.Values["session.id"] = session.SessionID
	}

	if session.Cost == nil {
		return result, nil
	}

	result.Values["session.duration.total"] = FormatDuration(session.Cost.TotalDurationMS)
	result.Values["session.duration.api"] = FormatDuration(session.Cost.TotalAPIDurationMS)

	if session.Cost.TotalLinesAdded > 0 {
		result.Values["session.lines-added"] = session.Cost.TotalLinesAdded
	}
	if session.Cost.TotalLinesRemoved > 0 {
		result.Values["session.lines-removed"] = session.Cost.TotalLinesRemoved
	}

	return result, nil
}

// FormatDuration formats milliseconds into a human-readable duration.
func FormatDuration(ms float64) string {
	totalMinutes := int(ms / 60_000)
	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
