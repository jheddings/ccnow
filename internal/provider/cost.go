package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

type costProvider struct{}

func (p *costProvider) Name() string { return "cost" }

func (p *costProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	result := &types.ProviderResult{
		Values: map[string]any{
			"cost.usd": "",
		},
	}
	if session.Cost != nil {
		result.Values["cost.usd"] = fmt.Sprintf("$%.2f", session.Cost.TotalCostUSD)
	}
	return result, nil
}
