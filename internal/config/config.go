package config

import (
	"encoding/json"

	"github.com/jheddings/ccglow/internal/types"
)

type configFile struct {
	Segments []json.RawMessage `json:"segments"`
}

// Parse parses a JSON config file into a segment tree.
func Parse(data []byte) ([]types.SegmentNode, error) {
	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	var nodes []types.SegmentNode
	for _, raw := range cfg.Segments {
		var node types.SegmentNode
		if err := json.Unmarshal(raw, &node); err != nil {
			continue
		}
		if node.Expr == "" && node.Value == nil && len(node.Children) == 0 {
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}
