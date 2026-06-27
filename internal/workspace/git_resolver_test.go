package workspace

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGitResolverPassthroughWithoutGit(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	got, err := GitResolver{}.Resolve(t.Context(), dir)
	if err != nil {
		t.Fatalf("GitResolver.Resolve: %v", err)
	}
	if got != dir {
		t.Fatalf("got %q, want %q", got, dir)
	}
}

func TestGitResolverEmptyPath(t *testing.T) {
	t.Parallel()

	_, err := GitResolver{}.Resolve(t.Context(), "  ")
	if err == nil {
		t.Fatal("expected error for empty workspace")
	}
}

func TestGitResolverFindsGitRootFromNestedPath(t *testing.T) {
	t.Parallel()
	requireGit(t)

	repo := initGitRepo(t)
	nested := filepath.Join(repo, "internal", "workspace")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	got, err := GitResolver{}.Resolve(t.Context(), nested)
	if err != nil {
		t.Fatalf("GitResolver.Resolve: %v", err)
	}
	want := canonicalPath(t, repo)
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGitResolverFindsLinkedWorktreeRoot(t *testing.T) {
	t.Parallel()
	requireGit(t)

	repo := initGitRepo(t)
	worktree := filepath.Join(t.TempDir(), "linked-worktree")
	runGit(t, repo, "worktree", "add", "--detach", worktree, "HEAD")

	nested := filepath.Join(worktree, "internal", "workspace")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	got, err := GitResolver{}.Resolve(t.Context(), nested)
	if err != nil {
		t.Fatalf("GitResolver.Resolve: %v", err)
	}
	want := canonicalPath(t, worktree)
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func requireGit(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found")
	}
}

func initGitRepo(t *testing.T) string {
	t.Helper()

	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.com")
	runGit(t, repo, "config", "user.name", "Test User")

	readme := filepath.Join(repo, "README.md")
	if err := os.WriteFile(readme, []byte("test repo\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "-m", "initial commit")

	return repo
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func canonicalPath(t *testing.T, path string) string {
	t.Helper()

	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	return resolved
}
