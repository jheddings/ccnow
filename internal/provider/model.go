package provider

import "github.com/jheddings/ccglow/internal/types"

type modelProvider struct{}

func (p *modelProvider) Name() string { return "model" }

func (p *modelProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	result := &types.ProviderResult{
		Values: map[string]any{
			"model.name": "",
			"model.id":   "",
		},
	}
	if session.Model != nil && session.Model.DisplayName != "" {
		result.Values["model.name"] = session.Model.DisplayName
	}
	if session.Model != nil && session.Model.ID != "" {
		result.Values["model.id"] = session.Model.ID
	}
	return result, nil
}
