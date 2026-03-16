package segment

import "github.com/jheddings/ccglow/internal/types"

// RegisterBuiltin adds all built-in segment implementations to the registry.
func RegisterBuiltin(registry *Registry) {
	registry.Register(&literalSegment{})
	registry.Register(&newlineSegment{})
}

// Registry maps segment type names to their implementations.
type Registry struct {
	segments map[string]types.Segment
}

// NewRegistry creates an empty segment registry.
func NewRegistry() *Registry {
	return &Registry{segments: make(map[string]types.Segment)}
}

// Register adds a segment implementation.
func (r *Registry) Register(seg types.Segment) {
	r.segments[seg.Name()] = seg
}

// Get returns the segment for the given type name, or nil.
func (r *Registry) Get(name string) types.Segment {
	return r.segments[name]
}

// --- Literal ---

type literalSegment struct{}

func (s *literalSegment) Name() string { return "literal" }
func (s *literalSegment) Render(ctx *types.SegmentContext) *string {
	if ctx.Props == nil {
		return nil
	}
	if text, ok := ctx.Props["text"].(string); ok {
		return &text
	}
	return nil
}

// --- Newline ---

type newlineSegment struct{}

func (s *newlineSegment) Name() string { return "newline" }
func (s *newlineSegment) Render(ctx *types.SegmentContext) *string {
	v := "\n"
	return &v
}
