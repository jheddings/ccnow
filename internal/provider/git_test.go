package provider

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func skipWithoutGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
}

// initTempRepo creates a temp dir with git init, an initial commit,
// and returns the path.
func initTempRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v failed: %s", args, out)
		}
	}

	// Create and commit a seed file so HEAD exists
	seed := filepath.Join(dir, "seed.txt")
	if err := os.WriteFile(seed, []byte("seed"), 0644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"git", "add", "seed.txt"},
		{"git", "commit", "-m", "initial"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v failed: %s", args, out)
		}
	}

	return dir
}

func TestGitStatusCounts(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	// Modified: edit an existing tracked file (unstaged)
	os.WriteFile(filepath.Join(dir, "seed.txt"), []byte("changed"), 0644)

	// Staged: create and stage a new file
	os.WriteFile(filepath.Join(dir, "staged.txt"), []byte("new"), 0644)
	cmd := exec.Command("git", "add", "staged.txt")
	cmd.Dir = dir
	cmd.Run()

	// Untracked: create a file without adding
	os.WriteFile(filepath.Join(dir, "untracked.txt"), []byte("extra"), 0644)

	p := &gitProvider{}
	sess := &types.SessionData{CWD: dir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["git.modified"] != 1 {
		t.Errorf("expected 1 modified, got %v", result.Values["git.modified"])
	}
	if result.Values["git.staged"] != 1 {
		t.Errorf("expected 1 staged, got %v", result.Values["git.staged"])
	}
	if result.Values["git.untracked"] != 1 {
		t.Errorf("expected 1 untracked, got %v", result.Values["git.untracked"])
	}
}

func TestGitStatusClean(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	p := &gitProvider{}
	sess := &types.SessionData{CWD: dir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["git.modified"] != 0 {
		t.Errorf("expected 0 modified, got %v", result.Values["git.modified"])
	}
	if result.Values["git.staged"] != 0 {
		t.Errorf("expected 0 staged, got %v", result.Values["git.staged"])
	}
	if result.Values["git.untracked"] != 0 {
		t.Errorf("expected 0 untracked, got %v", result.Values["git.untracked"])
	}
}

func TestGitWorktreeDetection(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	// Create a linked worktree
	wtDir := filepath.Join(t.TempDir(), "my-worktree")
	cmd := exec.Command("git", "worktree", "add", wtDir, "-b", "test-branch")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("worktree add failed: %s", out)
	}

	p := &gitProvider{}
	sess := &types.SessionData{CWD: wtDir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["git.worktree"] != "my-worktree" {
		t.Errorf("expected worktree name 'my-worktree', got %q", result.Values["git.worktree"])
	}
}

func TestGitRemoteOwnerRepo(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	// Add an HTTPS remote
	cmd := exec.Command("git", "remote", "add", "origin", "https://github.com/myowner/myrepo.git")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("remote add failed: %s", out)
	}

	p := &gitProvider{}
	sess := &types.SessionData{CWD: dir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["git.owner"] != "myowner" {
		t.Errorf("expected owner 'myowner', got %q", result.Values["git.owner"])
	}
	if result.Values["git.repo"] != "myrepo" {
		t.Errorf("expected repo 'myrepo', got %q", result.Values["git.repo"])
	}
}

func TestGitRemoteSSH(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	// Add an SSH remote
	cmd := exec.Command("git", "remote", "add", "origin", "git@github.com:sshowner/sshrepo.git")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("remote add failed: %s", out)
	}

	p := &gitProvider{}
	sess := &types.SessionData{CWD: dir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["git.owner"] != "sshowner" {
		t.Errorf("expected owner 'sshowner', got %q", result.Values["git.owner"])
	}
	if result.Values["git.repo"] != "sshrepo" {
		t.Errorf("expected repo 'sshrepo', got %q", result.Values["git.repo"])
	}
}

func TestGitRemoteNoOrigin(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	// No remote added

	p := &gitProvider{}
	sess := &types.SessionData{CWD: dir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["git.owner"] != "" {
		t.Errorf("expected empty Owner, got %q", result.Values["git.owner"])
	}
	if result.Values["git.repo"] != "" {
		t.Errorf("expected empty Repo, got %q", result.Values["git.repo"])
	}
}

func TestGitWorktreeMainCopy(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	p := &gitProvider{}
	sess := &types.SessionData{CWD: dir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if result.Values["git.worktree"] != "" {
		t.Errorf("expected empty Worktree in main copy, got %q", result.Values["git.worktree"])
	}
}
