package render

import (
	"strings"
	"sync"
	"time"

	"github.com/jheddings/ccglow/internal/eval"
	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
	"github.com/rs/zerolog/log"
)

func isEnabled(node *types.SegmentNode, session *types.SessionData) bool {
	if node.EnabledFn != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Warn().Interface("panic", r).Msg("enabledFn panicked")
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
	session *types.SessionData,
	env map[string]any,
	defaultFormats map[string]string,
) *string {
	if !isEnabled(node, session) {
		return nil
	}

	// Composite: evaluate when, then render children
	if len(node.Children) > 0 {
		if node.When != "" {
			c := eval.CompileCached(node.When)
			if c == nil {
				return nil // compilation failed
			}
			segEnv := eval.BuildSegmentEnv(env, nil, "")
			if !c.Evaluate(segEnv) {
				return nil
			}
		}

		var parts []string
		for i := range node.Children {
			rendered := renderNode(&node.Children[i], session, env, defaultFormats)
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

	// Resolve raw value
	var raw any
	var hasValue bool

	if node.Value != nil {
		raw = node.Value
		hasValue = true
	} else if node.Expr != "" {
		result, err := eval.Eval(node.Expr, env)
		if err != nil {
			log.Warn().Err(err).Str("expr", node.Expr).Msg("expr eval failed")
			return nil
		}
		raw = result
		hasValue = true
	}

	if !hasValue {
		return nil
	}

	// Resolve format: config override > provider default > none
	format := node.Format
	if format == "" && node.Expr != "" {
		format = defaultFormats[node.Expr]
	}

	text := FormatValue(raw, format)
	if text == "" {
		return nil
	}

	// Evaluate when expression
	if node.When != "" {
		c := eval.CompileCached(node.When)
		if c == nil {
			return nil // compilation failed
		}
		segEnv := eval.BuildSegmentEnv(env, raw, text)
		if !c.Evaluate(segEnv) {
			return nil
		}
	}

	styled := style.Apply(text, node.Style)
	return &styled
}

// Tree performs a depth-first traversal of the segment tree,
// resolving each node against the environment and default formats.
func Tree(
	tree []types.SegmentNode,
	session *types.SessionData,
	env map[string]any,
	defaultFormats map[string]string,
) string {
	var parts []string
	for i := range tree {
		rendered := renderNode(&tree[i], session, env, defaultFormats)
		if rendered != nil {
			parts = append(parts, *rendered)
		}
	}
	return strings.Join(parts, "")
}

// BuildEnv resolves all providers concurrently and merges their nested
// results into a single environment map. Returns the merged env and
// flat default format map.
func BuildEnv(
	providers map[string]types.DataProvider,
	session *types.SessionData,
) (env map[string]any, defaultFormats map[string]string) {
	env = make(map[string]any)
	defaultFormats = make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, p := range providers {
		wg.Add(1)
		go func(prov types.DataProvider) {
			defer wg.Done()
			start := time.Now()
			result, err := prov.Resolve(session)
			elapsed := time.Since(start)
			if err != nil {
				log.Warn().Err(err).Str("provider", prov.Name()).Msg("provider resolve failed")
				return
			}
			mu.Lock()
			for k, v := range result.Values {
				// inject __metrics__ into the provider's value subtree
				if m, ok := v.(map[string]any); ok {
					m["__metrics__"] = map[string]any{
						"duration_ms": elapsed.Seconds() * 1000,
					}
				}
				env[k] = v
			}
			for k, v := range result.Formats {
				defaultFormats[k] = v
			}
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	return env, defaultFormats
}
