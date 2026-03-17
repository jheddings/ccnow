package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

type costProvider struct{}

func (p *costProvider) Name() string { return "cost" }

func (p *costProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	cost := map[string]any{
		"usd": "",
	}
	if session.Cost != nil {
		cost["usd"] = fmt.Sprintf("$%.2f", session.Cost.TotalCostUSD)
	}
	return &types.ProviderResult{
		Values: map[string]any{"cost": cost},
	}, nil
}
