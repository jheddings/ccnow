package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestCostProvider(t *testing.T) {
	p := &costProvider{}
	sess := &types.SessionData{
		CWD:  "/tmp",
		Cost: &types.CostInfo{TotalCostUSD: 12.5},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["cost.usd"] != "$12.50" {
		t.Errorf("expected $12.50, got %s", result.Values["cost.usd"])
	}
}
