package provider

import "github.com/jheddings/ccglow/internal/types"

type claudeProvider struct{}

func (p *claudeProvider) Name() string { return "claude" }

func (p *claudeProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	result := &types.ProviderResult{
		Values: map[string]any{
			"claude.version": "",
			"claude.style":   "",
		},
	}
	if session.Version != "" {
		result.Values["claude.version"] = session.Version
	}
	if session.OutputStyle != nil && session.OutputStyle.Name != "" {
		result.Values["claude.style"] = session.OutputStyle.Name
	}
	return result, nil
}
