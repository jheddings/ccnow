package condition

import (
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// Condition is a compiled when expression.
type Condition struct {
	program *vm.Program
}

// Compile compiles an expression string into a reusable Condition.
// Returns nil for empty expressions (always true).
func Compile(expression string) (*Condition, error) {
	if expression == "" {
		return nil, nil
	}

	program, err := expr.Compile(expression, expr.AllowUndefinedVariables())
	if err != nil {
		return nil, err
	}

	return &Condition{program: program}, nil
}

// Evaluate runs the compiled expression against the environment.
// Returns true only if the result is boolean true.
// Nil receiver (empty expression) returns true.
func (c *Condition) Evaluate(env map[string]any) bool {
	if c == nil {
		return true
	}

	result, err := expr.Run(c.program, env)
	if err != nil {
		return false
	}

	b, ok := result.(bool)
	return ok && b
}

// BuildNestedEnv converts a flat segment values map into nested maps
// for expr-lang member access. "git.repo" becomes env["git"]["repo"].
func BuildNestedEnv(segmentValues map[string]any) map[string]any {
	env := make(map[string]any)

	for key, value := range segmentValues {
		parts := strings.SplitN(key, ".", 2)
		if len(parts) != 2 {
			env[key] = value
			continue
		}

		ns := parts[0]
		field := parts[1]

		sub, ok := env[ns].(map[string]any)
		if !ok {
			sub = make(map[string]any)
			env[ns] = sub
		}

		// Handle dotted field names like "percent.used" -> nested further
		fieldParts := strings.SplitN(field, ".", 2)
		if len(fieldParts) == 2 {
			inner, ok := sub[fieldParts[0]].(map[string]any)
			if !ok {
				inner = make(map[string]any)
				sub[fieldParts[0]] = inner
			}
			inner[fieldParts[1]] = value
		} else {
			sub[field] = value
		}
	}

	return env
}

// BuildSegmentEnv creates the evaluation environment for a single segment's
// when expression. It shallow-copies the nested env and adds value/text keys.
func BuildSegmentEnv(nestedEnv map[string]any, value any, text string) map[string]any {
	env := make(map[string]any, len(nestedEnv)+2)
	for k, v := range nestedEnv {
		env[k] = v
	}
	env["value"] = value
	env["text"] = text
	return env
}
