package preset

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/jheddings/ccglow/internal/config"
	"github.com/jheddings/ccglow/internal/types"
)

//go:embed *.json
var presetFS embed.FS

// Get returns the segment tree for a named preset, or nil.
func Get(name string) []types.SegmentNode {
	data, err := presetFS.ReadFile(name + ".json")
	if err != nil {
		return nil
	}

	nodes, err := config.Parse(data)
	if err != nil {
		return nil
	}

	return nodes
}

// List returns all available preset names.
func List() []string {
	entries, err := presetFS.ReadDir(".")
	if err != nil {
		return nil
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			names = append(names, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	sort.Strings(names)
	return names
}

// Dump returns the raw JSON for a named preset.
func Dump(name string) ([]byte, error) {
	data, err := presetFS.ReadFile(name + ".json")
	if err != nil {
		return nil, fmt.Errorf("unknown preset: %s", name)
	}
	return data, nil
}
