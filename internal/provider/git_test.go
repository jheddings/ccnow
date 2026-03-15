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

	data := result.(*GitData)
	if data.Modified == nil {
		t.Error("expected non-nil Modified")
	} else if *data.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", *data.Modified)
	}
	if data.Staged == nil {
		t.Error("expected non-nil Staged")
	} else if *data.Staged != 1 {
		t.Errorf("expected 1 staged, got %d", *data.Staged)
	}
	if data.Untracked == nil {
		t.Error("expected non-nil Untracked")
	} else if *data.Untracked != 1 {
		t.Errorf("expected 1 untracked, got %d", *data.Untracked)
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

	data := result.(*GitData)
	if data.Modified == nil {
		t.Fatal("expected non-nil Modified for clean repo")
	}
	if *data.Modified != 0 {
		t.Errorf("expected 0 modified, got %d", *data.Modified)
	}
	if data.Staged == nil || *data.Staged != 0 {
		t.Errorf("expected 0 staged, got %v", data.Staged)
	}
	if data.Untracked == nil || *data.Untracked != 0 {
		t.Errorf("expected 0 untracked, got %v", data.Untracked)
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

	data := result.(*GitData)
	if data.Worktree == nil {
		t.Fatal("expected non-nil Worktree in linked worktree")
	}
	if *data.Worktree != "my-worktree" {
		t.Errorf("expected worktree name 'my-worktree', got %q", *data.Worktree)
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

	data := result.(*GitData)
	if data.Worktree != nil {
		t.Errorf("expected nil Worktree in main copy, got %q", *data.Worktree)
	}
}
