package provider

import "github.com/jheddings/ccglow/internal/types"

// ClaudeData holds resolved Claude Code application metadata.
type ClaudeData struct {
	Version *string `segment:"claude.version"`
	Style   *string `segment:"claude.style"`
}

func (p *claudeProvider) Fields() any { return &ClaudeData{} }

type claudeProvider struct{}

func (p *claudeProvider) Name() string { return "claude" }

func (p *claudeProvider) Resolve(session *types.SessionData) (any, error) {
	data := &ClaudeData{}
	if session.Version != "" {
		data.Version = &session.Version
	}
	if session.OutputStyle != nil && session.OutputStyle.Name != "" {
		data.Style = &session.OutputStyle.Name
	}
	return data, nil
}
