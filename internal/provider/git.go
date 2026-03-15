package provider

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
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

	if mod, stg, unt, err := parseGitStatus(cwd); err == nil {
		data.Modified = intPtr(mod)
		data.Staged = intPtr(stg)
		data.Untracked = intPtr(unt)
	}

	data.Worktree = detectWorktree(cwd)

	return data, nil
}

func parseGitStatus(cwd string) (modified, staged, untracked int, err error) {
	// Use gitExecRaw to preserve leading whitespace in porcelain output,
	// since the first column position is significant.
	out, err := gitExecRaw(cwd, "status", "--porcelain")
	if err != nil {
		return 0, 0, 0, err
	}
	if out == "" {
		return 0, 0, 0, nil
	}
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 2 {
			continue
		}
		if strings.HasPrefix(line, "??") {
			untracked++
			continue
		}
		x, y := line[0], line[1]
		// Column 1: staged changes
		if x == 'M' || x == 'A' || x == 'D' || x == 'R' || x == 'C' {
			staged++
		}
		// Column 2: unstaged changes
		if y == 'M' || y == 'D' || y == 'T' {
			modified++
		}
	}
	return modified, staged, untracked, nil
}

func intPtr(n int) *int { return &n }

// detectWorktree returns the worktree name if cwd is inside a linked
// worktree, or nil if it is the main working copy.
func detectWorktree(cwd string) *string {
	gitDir, err := gitExec(cwd, "rev-parse", "--git-dir")
	if err != nil {
		return nil
	}
	commonDir, err := gitExec(cwd, "rev-parse", "--git-common-dir")
	if err != nil {
		return nil
	}
	// Normalize to absolute paths for comparison
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(cwd, gitDir)
	}
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(cwd, commonDir)
	}
	gitDir = filepath.Clean(gitDir)
	commonDir = filepath.Clean(commonDir)

	if gitDir == commonDir {
		return nil // main working copy
	}

	// In a linked worktree — get the worktree root name
	toplevel, err := gitExec(cwd, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil
	}
	name := filepath.Base(toplevel)
	return &name
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

// gitExecRaw runs a git command and returns output with only trailing
// whitespace trimmed, preserving leading whitespace which is significant
// in commands like "status --porcelain".
func gitExecRaw(cwd string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", cwd}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), " \t\n\r"), nil
}
