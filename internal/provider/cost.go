package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

// CostData holds resolved session cost information.
type CostData struct {
	USD *string `segment:"cost.usd"`
}

func (p *costProvider) Fields() any { return &CostData{} }

type costProvider struct{}

func (p *costProvider) Name() string { return "cost" }

func (p *costProvider) Resolve(session *types.SessionData) (any, error) {
	data := &CostData{}
	if session.Cost != nil {
		usd := fmt.Sprintf("$%.2f", session.Cost.TotalCostUSD)
		data.USD = &usd
	}
	return data, nil
}
