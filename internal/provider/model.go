package provider

import "github.com/jheddings/ccglow/internal/types"

// ModelData holds resolved model information.
type ModelData struct {
	Name *string `segment:"model.name"`
}

func (p *modelProvider) Fields() any { return &ModelData{} }

type modelProvider struct{}

func (p *modelProvider) Name() string { return "model" }

func (p *modelProvider) Resolve(session *types.SessionData) (any, error) {
	data := &ModelData{}
	if session.Model != nil && session.Model.DisplayName != "" {
		data.Name = &session.Model.DisplayName
	}
	return data, nil
}
