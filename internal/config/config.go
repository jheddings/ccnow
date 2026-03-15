package config

import (
	"encoding/json"
	"strings"

	"github.com/jheddings/ccglow/internal/types"
)

var noProviderSegments = map[string]bool{
	"literal": true,
	"newline": true,
	"group":   true,
}

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
		if node.Type == "" {
			continue
		}
		nodes = append(nodes, node)
	}

	InferProviders(nodes)
	return nodes, nil
}

// InferProviders sets the Provider field on nodes that don't have one,
// based on the segment type prefix (e.g. "git.branch" → provider "git").
func InferProviders(nodes []types.SegmentNode) {
	for i := range nodes {
		if nodes[i].Provider == "" && nodes[i].Type != "" && !noProviderSegments[nodes[i].Type] {
			parts := strings.SplitN(nodes[i].Type, ".", 2)
			if len(parts) > 0 {
				nodes[i].Provider = parts[0]
			}
		}
		if len(nodes[i].Children) > 0 {
			InferProviders(nodes[i].Children)
		}
	}
}
