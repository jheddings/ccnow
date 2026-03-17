package provider

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jheddings/ccglow/internal/types"
)

type pwdProvider struct{}

func (p *pwdProvider) Name() string { return "pwd" }

func (p *pwdProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	cwd := session.CWD
	name := filepath.Base(cwd)
	dir := filepath.Dir(cwd)
	if dir != "/" {
		dir += "/"
	}

	return &types.ProviderResult{
		Values: map[string]any{
			"pwd": map[string]any{
				"name":  name,
				"path":  dir,
				"smart": smartPrefix(cwd),
			},
		},
	}, nil
}

func smartPrefix(cwd string) string {
	if cwd == "/" {
		return ""
	}

	home, _ := os.UserHomeDir()

	display := cwd
	if home != "" && strings.HasPrefix(cwd, home) {
		display = "~" + cwd[len(home):]
	}

	dir := filepath.Dir(display)
	if dir == "." || dir == "/" {
		return ""
	}

	// Separate the root prefix from the relative path segments
	root := ""
	rel := dir
	if strings.HasPrefix(dir, "~/") {
		root = "~/"
		rel = dir[2:]
	} else if dir == "~" {
		return "~/"
	} else if strings.HasPrefix(dir, "/") {
		root = "/"
		rel = dir[1:]
	}

	if rel == "" {
		return root
	}

	segments := strings.Split(rel, "/")

	if len(segments) <= 2 {
		return root + strings.Join(segments, "/") + "/"
	}

	// Abbreviate: first char of leading parts, then ellipsis
	var abbrev []string
	for i := 0; i < len(segments)-1 && i < 2; i++ {
		if len(segments[i]) > 0 {
			abbrev = append(abbrev, string(segments[i][0]))
		}
	}
	abbrev = append(abbrev, "\u2026")

	return root + strings.Join(abbrev, "/") + "/"
}
