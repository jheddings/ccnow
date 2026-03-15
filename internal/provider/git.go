package provider

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/jheddings/ccglow/internal/types"
)

const gitTimeout = 5 * time.Second

// GitData holds resolved git repository information.
type GitData struct {
	Branch     *string
	Insertions *int
	Deletions  *int
	Modified   *int
	Staged     *int
	Untracked  *int
	Worktree   *string
}

type gitProvider struct{}

func (p *gitProvider) Name() string { return "git" }

func (p *gitProvider) Resolve(session *types.SessionData) (any, error) {
	cwd := session.CWD

	if !gitAvailable(cwd) {
		return &GitData{}, nil
	}

	data := &GitData{}

	if branch, err := gitExec(cwd, "branch", "--show-current"); err == nil && branch != "" {
		data.Branch = &branch
	}

	if diff, err := gitExec(cwd, "diff", "--shortstat", "HEAD"); err == nil && diff != "" {
		if m := insertionRe.FindStringSubmatch(diff); len(m) > 1 {
			var n int
			fmt.Sscanf(m[1], "%d", &n)
			data.Insertions = &n
		}
		if m := deletionRe.FindStringSubmatch(diff); len(m) > 1 {
			var n int
			fmt.Sscanf(m[1], "%d", &n)
			data.Deletions = &n
		}
	}

	return data, nil
}

var (
	insertionRe = regexp.MustCompile(`(\d+) insertion`)
	deletionRe  = regexp.MustCompile(`(\d+) deletion`)
)

func gitAvailable(cwd string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()
	return exec.CommandContext(ctx, "git", "-C", cwd, "rev-parse", "--git-dir").Run() == nil
}

func gitExec(cwd string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", cwd}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
