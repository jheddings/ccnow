package command

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// DefaultTimeout is the maximum time a command is allowed to run.
const DefaultTimeout = 2 * time.Second

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Interpolate replaces ${dotted.path} references in s by walking the nested
// env map. Unresolved references are replaced with an empty string.
func Interpolate(s string, env map[string]any) string {
	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		// Extract the dotted path from ${...}
		path := varPattern.FindStringSubmatch(match)[1]
		parts := strings.Split(path, ".")

		var current any = env
		for _, part := range parts {
			m, ok := current.(map[string]any)
			if !ok {
				return ""
			}
			current, ok = m[part]
			if !ok {
				return ""
			}
		}

		return fmt.Sprintf("%v", current)
	})
}

// Run interpolates variables in cmdStr, executes it via sh -c, and returns
// trimmed stdout. Empty output or non-zero exit returns an empty string.
func Run(cmdStr string, env map[string]any, cwd string, timeout time.Duration) string {
	if cmdStr == "" {
		return ""
	}

	interpolated := Interpolate(cmdStr, env)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", interpolated)
	if cwd != "" {
		cmd.Dir = cwd
	}

	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
