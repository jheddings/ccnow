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
	segmentValues map[string]any,
	defaultFormats map[string]string,
	conditionEnv map[string]any,
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
			env := condition.BuildSegmentEnv(conditionEnv, nil, "")
			if !c.Evaluate(env) {
				return nil
			}
		}

		var parts []string
		for i := range node.Children {
			rendered := renderNode(&node.Children[i], segments, session, segmentValues, defaultFormats, conditionEnv)
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
	if !ok {
		return nil
	}

	// Resolve format: config override > provider default > none
	format := node.Format
	if format == "" {
		format = defaultFormats[node.Type]
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
		env := condition.BuildSegmentEnv(conditionEnv, value, text)
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
	segmentValues map[string]any,
	defaultFormats map[string]string,
	conditionEnv map[string]any,
) string {
	var parts []string
	for i := range tree {
		rendered := renderNode(&tree[i], segments, session, segmentValues, defaultFormats, conditionEnv)
		if rendered != nil {
			parts = append(parts, *rendered)
		}
	}
	return strings.Join(parts, "")
}

// ResolveAll resolves all providers and returns merged segment values
// and default formats.
func ResolveAll(
	providers map[string]types.DataProvider,
	session *types.SessionData,
) (segmentValues map[string]any, defaultFormats map[string]string) {
	segmentValues = make(map[string]any)
	defaultFormats = make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, p := range providers {
		wg.Add(1)
		go func(prov types.DataProvider) {
			defer wg.Done()
			result, err := prov.Resolve(session)
			if err != nil {
				log.Warn().Err(err).Str("provider", prov.Name()).Msg("provider resolve failed")
				return
			}
			mu.Lock()
			for k, v := range result.Values {
				segmentValues[k] = v
			}
			for k, v := range result.Formats {
				defaultFormats[k] = v
			}
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	return segmentValues, defaultFormats
}
