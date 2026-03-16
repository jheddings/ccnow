package condition

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// dotFieldRe matches dot-prefixed field references like .name, .count
var dotFieldRe = regexp.MustCompile(`\.([a-zA-Z_][a-zA-Z0-9_]*)`)

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

	// Rewrite .field references to __field for the expr engine
	rewritten := dotFieldRe.ReplaceAllString(expression, "__$1")

	program, err := expr.Compile(rewritten, expr.AllowUndefinedVariables())
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

	// Rewrite .field keys to __field to match compiled expression
	translated := make(map[string]any, len(env))
	for k, v := range env {
		if strings.HasPrefix(k, ".") {
			translated["__"+k[1:]] = v
		} else {
			translated[k] = v
		}
	}

	result, err := expr.Run(c.program, translated)
	if err != nil {
		return false
	}

	b, ok := result.(bool)
	return ok && b
}

// BuildEnv builds the variable environment for expression evaluation.
func BuildEnv(providerData any, value any, text string) map[string]any {
	env := make(map[string]any)

	if providerData != nil {
		v := reflect.ValueOf(providerData)
		if v.Kind() == reflect.Ptr {
			if !v.IsNil() {
				v = v.Elem()
			} else {
				v = reflect.Value{}
			}
		}

		if v.IsValid() && v.Kind() == reflect.Struct {
			t := v.Type()
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				if !field.IsExported() {
					continue
				}

				key := "." + strings.ToLower(field.Name)
				fv := v.Field(i)

				if fv.Kind() == reflect.Ptr {
					if fv.IsNil() {
						env[key] = coerceNilPointer(field.Type)
						continue
					}
					fv = fv.Elem()
				}

				env[key] = fv.Interface()
			}
		}
	}

	env["value"] = value
	env["text"] = text

	return env
}

func coerceNilPointer(t reflect.Type) any {
	elem := t.Elem()
	switch elem.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0
	case reflect.Float32, reflect.Float64:
		return 0.0
	case reflect.String:
		return ""
	case reflect.Bool:
		return false
	default:
		return nil
	}
}
