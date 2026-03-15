package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

// SessionData holds resolved session timing and line-change data.
type SessionData struct {
	Duration     *string `segment:"session.duration"`
	LinesAdded   *int    `segment:"session.lines-added"`
	LinesRemoved *int    `segment:"session.lines-removed"`
}

func (p *sessionProvider) Fields() any { return &SessionData{} }

type sessionProvider struct{}

func (p *sessionProvider) Name() string { return "session" }

func (p *sessionProvider) Resolve(session *types.SessionData) (any, error) {
	data := &SessionData{}
	if session.Cost == nil {
		return data, nil
	}

	dur := FormatDuration(session.Cost.TotalDurationMS)
	data.Duration = &dur

	if session.Cost.TotalLinesAdded > 0 {
		n := session.Cost.TotalLinesAdded
		data.LinesAdded = &n
	}
	if session.Cost.TotalLinesRemoved > 0 {
		n := session.Cost.TotalLinesRemoved
		data.LinesRemoved = &n
	}

	return data, nil
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
