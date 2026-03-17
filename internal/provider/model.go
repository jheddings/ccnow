package provider

import "github.com/jheddings/ccglow/internal/types"

type modelProvider struct{}

func (p *modelProvider) Name() string { return "model" }

func (p *modelProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	model := map[string]any{
		"name": "",
		"id":   "",
	}
	if session.Model != nil && session.Model.DisplayName != "" {
		model["name"] = session.Model.DisplayName
	}
	if session.Model != nil && session.Model.ID != "" {
		model["id"] = session.Model.ID
	}
	return &types.ProviderResult{
		Values: map[string]any{"model": model},
	}, nil
}
