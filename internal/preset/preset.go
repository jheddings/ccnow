package preset

import "github.com/jheddings/ccnow/internal/types"

// Get returns the segment tree for a named preset, or nil.
func Get(name string) []types.SegmentNode {
	switch name {
	case "default":
		return defaultPreset()
	case "minimal":
		return minimalPreset()
	case "full":
		return fullPreset()
	default:
		return nil
	}
}

// List returns all available preset names.
func List() []string {
	return []string{"default", "minimal", "full"}
}
