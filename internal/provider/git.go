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
	Branch     *string `segment:"git.branch"`
	Insertions *int    `segment:"git.insertions"`
	Deletions  *int    `segment:"git.deletions"`
	Modified   *int    `segment:"git.modified"`
	Staged     *int    `segment:"git.staged"`
	Untracked  *int    `segment:"git.untracked"`
	Worktree   *string `segment:"git.worktree"`
	Owner      *string `segment:"git.owner"`
	Repo       *string `segment:"git.repo"`
}

func (p *gitProvider) Fields() any { return &GitData{} }

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

	if owner, repo, err := parseRemoteOwnerRepo(cwd); err == nil {
		data.Owner = &owner
		data.Repo = &repo
	}

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

// parseRemoteOwnerRepo extracts the owner and repository name from the
// origin remote URL.  It handles both SSH (git@host:owner/repo.git) and
// HTTPS (https://host/owner/repo.git) formats.
func parseRemoteOwnerRepo(cwd string) (owner, repo string, err error) {
	url, err := gitExec(cwd, "remote", "get-url", "origin")
	if err != nil {
		return "", "", err
	}

	// Normalize SSH URLs: git@host:owner/repo -> host/owner/repo
	if strings.Contains(url, ":") && !strings.Contains(url, "://") {
		// SSH format — replace first ":" after the host with "/"
		url = url[strings.Index(url, ":")+1:]
	} else {
		// HTTPS format — strip scheme and host
		// e.g. https://github.com/owner/repo.git -> /owner/repo.git
		if idx := strings.Index(url, "://"); idx != -1 {
			url = url[idx+3:]
			// Remove host portion
			if slash := strings.Index(url, "/"); slash != -1 {
				url = url[slash+1:]
			}
		}
	}

	// Strip .git suffix
	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("cannot parse owner/repo from remote URL")
	}

	owner = parts[len(parts)-2]
	repo = parts[len(parts)-1]
	return owner, repo, nil
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
