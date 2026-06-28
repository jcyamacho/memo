package fsstore

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/jcyamacho/memo/internal/memory"
)

func TestFSRoundTrip(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	store, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx := context.Background()
	workspace := "/repo/project"
	m, err := memory.New("workspace memory", workspace)
	if err != nil {
		t.Fatalf("memory.New: %v", err)
	}
	if err := store.Insert(ctx, m); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	global, err := memory.New("global memory", "")
	if err != nil {
		t.Fatalf("memory.New global: %v", err)
	}
	if err := store.Insert(ctx, global); err != nil {
		t.Fatalf("Insert global: %v", err)
	}

	got, err := store.Get(ctx, m.ID())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Content() != m.Content() || got.Workspace() != workspace {
		t.Fatalf("Get returned %+v, want content %q workspace %q", got, m.Content(), workspace)
	}

	items, err := store.List(ctx, workspace, memory.WithGlobals())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("List len = %d, want 2", len(items))
	}

	if _, err := m.SetContent("updated content"); err != nil {
		t.Fatalf("SetContent: %v", err)
	}
	if err := store.Update(ctx, m); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err = store.Get(ctx, m.ID())
	if err != nil {
		t.Fatalf("Get after update: %v", err)
	}
	if got.Content() != "updated content" {
		t.Fatalf("content = %q, want updated", got.Content())
	}

	if err := store.Delete(ctx, m.ID()); err != nil {
		t.Fatalf("Delete workspace memory: %v", err)
	}
	if err := store.Delete(ctx, "missing"); !errors.Is(err, memory.ErrNotFound) {
		t.Fatalf("Delete missing = %v, want ErrNotFound", err)
	}
	if err := store.Delete(ctx, global.ID()); err != nil {
		t.Fatalf("Delete global memory: %v", err)
	}

	workspaces, err := store.ListWorkspaces(ctx)
	if err != nil {
		t.Fatalf("ListWorkspaces: %v", err)
	}
	if len(workspaces) != 0 {
		t.Fatalf("workspaces = %v, want empty after delete", workspaces)
	}

	encoded := url.PathEscape(workspace)
	if _, err := os.Stat(filepath.Join(dir, "workspaces", encoded)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected encoded workspace dir to be removed, stat err = %v", err)
	}
}

func TestFSWorkspacePathWithSeparators(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	store, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx := context.Background()
	workspace := "/repo/feature branch\\nested"
	m, err := memory.New("workspace memory", workspace)
	if err != nil {
		t.Fatalf("memory.New: %v", err)
	}
	if err := store.Insert(ctx, m); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	encoded := url.PathEscape(workspace)
	if encoded == workspace {
		t.Fatal("expected workspace path to be escaped")
	}
	if _, err := os.Stat(filepath.Join(dir, "workspaces", encoded)); err != nil {
		t.Fatalf("stat encoded workspace directory: %v", err)
	}

	workspaces, err := store.ListWorkspaces(ctx)
	if err != nil {
		t.Fatalf("ListWorkspaces: %v", err)
	}
	if len(workspaces) != 1 || workspaces[0] != workspace {
		t.Fatalf("workspaces = %v, want [%q]", workspaces, workspace)
	}

	items, err := store.List(ctx, workspace)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 1 || items[0].Workspace() != workspace {
		t.Fatalf("items = %+v, want one item scoped to %q", items, workspace)
	}
}

func TestFSListTargets(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	store, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx := context.Background()
	insertMemory := func(content, workspace string) {
		t.Helper()

		m, err := memory.New(content, workspace)
		if err != nil {
			t.Fatalf("memory.New: %v", err)
		}
		if err := store.Insert(ctx, m); err != nil {
			t.Fatalf("Insert: %v", err)
		}
	}

	insertMemory("global", "")
	insertMemory("workspace one", "/repo/one")
	insertMemory("workspace two", "/repo/two")

	tests := []struct {
		name      string
		workspace string
		options   []memory.ListOption
		contents  []string
	}{
		{
			name:     "all workspaces",
			contents: []string{"workspace one", "workspace two"},
		},
		{
			name:     "all workspaces with globals",
			options:  []memory.ListOption{memory.WithGlobals()},
			contents: []string{"global", "workspace one", "workspace two"},
		},
		{
			name:      "one workspace",
			workspace: "/repo/one",
			contents:  []string{"workspace one"},
		},
		{
			name:      "one workspace with globals",
			workspace: "/repo/one",
			options:   []memory.ListOption{memory.WithGlobals()},
			contents:  []string{"global", "workspace one"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			items, err := store.List(ctx, test.workspace, test.options...)
			if err != nil {
				t.Fatalf("List: %v", err)
			}

			got := make(map[string]struct{}, len(items))
			for _, item := range items {
				got[item.Content()] = struct{}{}
			}
			if len(got) != len(test.contents) {
				t.Fatalf("got contents %v, want %v", got, test.contents)
			}
			for _, content := range test.contents {
				if _, ok := got[content]; !ok {
					t.Fatalf("missing content %q in %v", content, got)
				}
			}
		})
	}
}

func TestFSInsertDuplicate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	store, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx := context.Background()
	m, err := memory.New("content", "")
	if err != nil {
		t.Fatalf("memory.New: %v", err)
	}
	if err := store.Insert(ctx, m); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if err := store.Insert(ctx, m); !errors.Is(err, memory.ErrExists) {
		t.Fatalf("Insert duplicate = %v, want ErrExists", err)
	}
}
