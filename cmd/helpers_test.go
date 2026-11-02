package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestReadContentUsesPositionalArgument(t *testing.T) {
	t.Parallel()

	got, err := readContentFrom([]string{"content"}, strings.NewReader("stdin"))
	if err != nil {
		t.Fatalf("readContentFrom: %v", err)
	}
	if got != "content" {
		t.Fatalf("got %q, want positional content", got)
	}
}

func TestReadContentUsesStdin(t *testing.T) {
	t.Parallel()

	got, err := readContentFrom(nil, strings.NewReader("stdin content\n"))
	if err != nil {
		t.Fatalf("readContentFrom: %v", err)
	}
	if got != "stdin content\n" {
		t.Fatalf("got %q, want stdin content", got)
	}
}

func TestReadContentRequiresContentWithoutStdin(t *testing.T) {
	t.Parallel()

	if _, err := readContentFrom(nil, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveContextWorkspacePathPrefersCursorProjectDir(t *testing.T) {
	project := filepath.Join(t.TempDir(), "project")
	t.Setenv("CURSOR_PROJECT_DIR", project)

	got, err := resolveContextWorkspacePath("")
	if err != nil {
		t.Fatalf("resolveContextWorkspacePath: %v", err)
	}
	testutil.AssertEqual(t, project, got)
}

func TestResolveContextWorkspacePathFlagOverridesCursorProjectDir(t *testing.T) {
	t.Setenv("CURSOR_PROJECT_DIR", filepath.Join(t.TempDir(), "cursor-project"))

	got, err := resolveContextWorkspacePath("/explicit")
	if err != nil {
		t.Fatalf("resolveContextWorkspacePath: %v", err)
	}
	testutil.AssertEqual(t, "/explicit", got)
}

func TestResolveWorkspacePathIgnoresCursorProjectDir(t *testing.T) {
	t.Setenv("CURSOR_PROJECT_DIR", filepath.Join(t.TempDir(), "cursor-project"))

	got, err := resolveWorkspacePath("")
	if err != nil {
		t.Fatalf("resolveWorkspacePath: %v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	testutil.AssertEqual(t, cwd, got)
}
