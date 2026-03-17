package provider

import "github.com/jheddings/ccglow/internal/types"

type claudeProvider struct{}

func (p *claudeProvider) Name() string { return "claude" }

func (p *claudeProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	claude := map[string]any{
		"version": "",
		"style":   "",
	}
	if session.Version != "" {
		claude["version"] = session.Version
	}
	if session.OutputStyle != nil && session.OutputStyle.Name != "" {
		claude["style"] = session.OutputStyle.Name
	}
	return &types.ProviderResult{
		Values: map[string]any{"claude": claude},
	}, nil
}
