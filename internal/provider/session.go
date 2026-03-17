package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

type sessionProvider struct{}

func (p *sessionProvider) Name() string { return "session" }

func (p *sessionProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	sess := map[string]any{
		"duration": map[string]any{
			"total":     "",
			"api":       "",
			"total_min": 0,
			"api_min":   0,
		},
		"lines-added":   0,
		"lines-removed": 0,
		"id":            "",
	}

	result := &types.ProviderResult{
		Values: map[string]any{"session": sess},
	}

	if session.SessionID != "" {
		sess["id"] = session.SessionID
	}

	if session.Cost == nil {
		return result, nil
	}

	dur := sess["duration"].(map[string]any)
	dur["total"] = FormatDuration(session.Cost.TotalDurationMS)
	dur["api"] = FormatDuration(session.Cost.TotalAPIDurationMS)
	dur["total_min"] = int(session.Cost.TotalDurationMS / 60_000)
	dur["api_min"] = int(session.Cost.TotalAPIDurationMS / 60_000)

	if session.Cost.TotalLinesAdded > 0 {
		sess["lines-added"] = session.Cost.TotalLinesAdded
	}
	if session.Cost.TotalLinesRemoved > 0 {
		sess["lines-removed"] = session.Cost.TotalLinesRemoved
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
