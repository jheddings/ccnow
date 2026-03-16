package render

import (
	"strings"
	"sync"

	"github.com/jheddings/ccglow/internal/condition"
	"github.com/jheddings/ccglow/internal/segment"
	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
	"github.com/rs/zerolog/log"
)

var conditionCache = make(map[string]*condition.Condition)
var conditionMu sync.Mutex

func getCondition(expr string) *condition.Condition {
	if expr == "" {
		return nil
	}
	conditionMu.Lock()
	defer conditionMu.Unlock()

	if c, ok := conditionCache[expr]; ok {
		return c
	}

	c, err := condition.Compile(expr)
	if err != nil {
		log.Warn().Err(err).Str("expr", expr).Msg("invalid when expression")
		conditionCache[expr] = nil
		return nil
	}
	conditionCache[expr] = c
	return c
}

func isEnabled(node *types.SegmentNode, session *types.SessionData) bool {
	if node.EnabledFn != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Warn().Str("type", node.Type).Interface("panic", r).Msg("enabledFn panicked")
			}
		}()
		return node.EnabledFn(session)
	}
	if node.Enabled != nil {
		return *node.Enabled
	}
	return true
}

func renderNode(
	node *types.SegmentNode,
	segments *segment.Registry,
	session *types.SessionData,
	providerData map[string]any,
	segmentValues map[string]any,
	tagIdx TagIndex,
) *string {
	if !isEnabled(node, session) {
		return nil
	}

	// SegmentGroup: evaluate when, then render children
	if len(node.Children) > 0 {
		if node.When != "" {
			c := getCondition(node.When)
			if c == nil {
				return nil // compilation failed
			}
			var pd any
			if node.Provider != "" {
				pd = providerData[node.Provider]
			}
			env := condition.BuildEnv(pd, nil, "")
			if !c.Evaluate(env) {
				return nil
			}
		}

		var parts []string
		for i := range node.Children {
			rendered := renderNode(&node.Children[i], segments, session, providerData, segmentValues, tagIdx)
			if rendered != nil {
				parts = append(parts, *rendered)
			}
		}
		if len(parts) == 0 {
			return nil
		}
		joined := strings.Join(parts, "")
		styled := style.Apply(joined, node.Style)
		return &styled
	}

	// Built-in segment: delegate to registered segment (literal, newline)
	seg := segments.Get(node.Type)
	if seg != nil {
		ctx := &types.SegmentContext{
			Session: session,
			Props:   node.Props,
		}

		value := seg.Render(ctx)
		if value == nil {
			return nil
		}
		styled := style.Apply(*value, node.Style)
		return &styled
	}

	// DataSegment: resolve from segment values
	value, ok := segmentValues[node.Type]
	if !ok || value == nil {
		return nil
	}

	// Resolve format: config override > tag default > none
	format := node.Format
	if format == "" {
		if accessor, exists := tagIdx[node.Type]; exists {
			format = accessor.DefaultFormat
		}
	}

	text := FormatValue(value, format)
	if text == "" {
		return nil
	}

	// Evaluate when expression
	if node.When != "" {
		c := getCondition(node.When)
		if c == nil {
			return nil // compilation failed
		}
		var pd any
		if accessor, exists := tagIdx[node.Type]; exists {
			pd = providerData[accessor.Provider]
		}
		env := condition.BuildEnv(pd, value, text)
		if !c.Evaluate(env) {
			return nil
		}
	}

	styled := style.Apply(text, node.Style)
	return &styled
}

// Tree performs a depth-first traversal of the segment tree,
// resolving each node against the registries and provider data.
func Tree(
	tree []types.SegmentNode,
	segments *segment.Registry,
	session *types.SessionData,
	providerData map[string]any,
	segmentValues map[string]any,
	tagIdx TagIndex,
) string {
	var parts []string
	for i := range tree {
		rendered := renderNode(&tree[i], segments, session, providerData, segmentValues, tagIdx)
		if rendered != nil {
			parts = append(parts, *rendered)
		}
	}
	return strings.Join(parts, "")
}

// CollectProviderNames walks the tree and returns the set of provider
// names needed for rendering (skipping disabled nodes).
func CollectProviderNames(tree []types.SegmentNode, tagIdx TagIndex) map[string]bool {
	names := make(map[string]bool)
	collectNames(tree, names, tagIdx)
	return names
}

func collectNames(nodes []types.SegmentNode, names map[string]bool, idx TagIndex) {
	for _, node := range nodes {
		if node.Enabled != nil && !*node.Enabled {
			continue
		}
		if accessor, ok := idx[node.Type]; ok {
			names[accessor.Provider] = true
		}
		if node.Provider != "" {
			names[node.Provider] = true
		}
		if len(node.Children) > 0 {
			collectNames(node.Children, names, idx)
		}
	}
}

// ResolveProviders resolves all named providers concurrently and returns
// a map of provider name → resolved data.
func ResolveProviders(
	names map[string]bool,
	providers map[string]types.DataProvider,
	session *types.SessionData,
) map[string]any {
	results := make(map[string]any)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for name := range names {
		p, ok := providers[name]
		if !ok {
			log.Warn().Str("provider", name).Msg("provider not registered")
			continue
		}
		wg.Add(1)
		go func(prov types.DataProvider) {
			defer wg.Done()
			data, err := prov.Resolve(session)
			if err != nil {
				log.Warn().Err(err).Str("provider", prov.Name()).Msg("provider resolve failed")
				return
			}
			mu.Lock()
			results[prov.Name()] = data
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	return results
}
