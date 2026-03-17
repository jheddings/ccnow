package eval

import (
	"sync"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/rs/zerolog/log"
)

// Condition is a compiled boolean expression (used for when guards).
type Condition struct {
	program *vm.Program
}

var (
	cache   = make(map[string]*vm.Program)
	cacheMu sync.Mutex
)

// compile returns a cached compiled program, or compiles and caches it.
func compile(expression string) (*vm.Program, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if p, ok := cache[expression]; ok {
		return p, nil
	}

	p, err := expr.Compile(expression, expr.AllowUndefinedVariables())
	if err != nil {
		return nil, err
	}

	cache[expression] = p
	return p, nil
}

// Compile compiles an expression string into a reusable Condition.
// Returns nil for empty expressions (always true).
func Compile(expression string) (*Condition, error) {
	if expression == "" {
		return nil, nil
	}

	p, err := compile(expression)
	if err != nil {
		return nil, err
	}

	return &Condition{program: p}, nil
}

// CompileCached compiles a when expression with caching and warning on failure.
// Returns nil for empty or invalid expressions.
func CompileCached(expression string) *Condition {
	if expression == "" {
		return nil
	}

	c, err := Compile(expression)
	if err != nil {
		log.Warn().Err(err).Str("expr", expression).Msg("invalid when expression")
		return nil
	}
	return c
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

// Eval compiles (with caching) and runs an expression against the environment,
// returning the result as any value. Used to resolve expr fields in nodes.
func Eval(expression string, env map[string]any) (any, error) {
	if expression == "" {
		return nil, nil
	}

	p, err := compile(expression)
	if err != nil {
		return nil, err
	}

	return expr.Run(p, env)
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
