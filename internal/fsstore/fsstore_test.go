package fsstore

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/memory"
)

func TestFSRoundTrip(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	store := newStoreAt(t, dir)
	ctx := t.Context()
	workspace := "/repo/project"
	m := mustMemory(t, "workspace memory", workspace)
	mustInsert(t, store, m)
	global := mustMemory(t, "global memory", "")
	mustInsert(t, store, global)

	got, err := store.Get(ctx, m.ID())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Content() != m.Content() || got.Workspace() != workspace {
		t.Fatalf("Get returned %+v, want content %q workspace %q", got, m.Content(), workspace)
	}

	items, err := store.List(ctx, workspace)
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
	missingID, err := memory.ParseID("550e8400-e29b-41d4-a716-446655440099")
	if err != nil {
		t.Fatalf("ParseID missing: %v", err)
	}
	if err := store.Delete(ctx, missingID); !errors.Is(err, memory.ErrNotFound) {
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
	store := newStoreAt(t, dir)
	ctx := t.Context()
	workspace := "/repo/feature branch\\nested"
	mustInsert(t, store, mustMemory(t, "workspace memory", workspace))

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

func TestFSListIncludesWorkspaceAndGlobals(t *testing.T) {
	t.Parallel()

	store := newStore(t)
	ctx := t.Context()
	mustInsert(t, store, mustMemory(t, "global", ""))
	mustInsert(t, store, mustMemory(t, "workspace one", "/repo/one"))
	mustInsert(t, store, mustMemory(t, "workspace two", "/repo/two"))

	items, err := store.List(ctx, "/repo/one")
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	got := make(map[string]struct{}, len(items))
	for _, item := range items {
		got[item.Content()] = struct{}{}
	}
	want := []string{"global", "workspace one"}
	if len(got) != len(want) {
		t.Fatalf("got contents %v, want %v", got, want)
	}
	for _, content := range want {
		if _, ok := got[content]; !ok {
			t.Fatalf("missing content %q in %v", content, got)
		}
	}
}

func TestFSInsertDuplicate(t *testing.T) {
	t.Parallel()

	store := newStore(t)
	ctx := t.Context()
	m := mustMemory(t, "content", "")
	mustInsert(t, store, m)
	if err := store.Insert(ctx, m); !errors.Is(err, memory.ErrExists) {
		t.Fatalf("Insert duplicate = %v, want ErrExists", err)
	}
}

func TestFSInsertRejectsIDUsedInAnotherScope(t *testing.T) {
	t.Parallel()

	store := newStore(t)
	ctx := t.Context()
	m := mustMemory(t, "content", "/repo/one")
	mustInsert(t, store, m)

	for _, workspace := range []string{"", "/repo/two"} {
		duplicate, err := memory.Load(m.ID().String(), "other content", workspace, time.Now())
		if err != nil {
			t.Fatalf("Load duplicate: %v", err)
		}
		if err := store.Insert(ctx, duplicate); !errors.Is(err, memory.ErrExists) {
			t.Fatalf("Insert duplicate in %q = %v, want ErrExists", workspace, err)
		}
	}
}

func TestFSNamedIDRoundTrip(t *testing.T) {
	t.Parallel()

	store := newStore(t)
	ctx := t.Context()
	workspace := "/repo/project"
	m, err := memory.New("named memory", "Prefer Bun", workspace)
	if err != nil {
		t.Fatalf("memory.New: %v", err)
	}
	mustInsert(t, store, m)

	path := store.memoryPath(m.ID(), workspace)
	if !strings.HasPrefix(filepath.Base(path), "prefer-bun-") {
		t.Fatalf("memory file = %q, want prefer-bun- prefix", path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stat named memory: %v", err)
	}

	items, err := store.List(ctx, workspace)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 1 || items[0].ID() != m.ID() {
		t.Fatalf("items = %+v, want one item with id %q", items, m.ID())
	}

	if err := store.Delete(ctx, m.ID()); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestFSUpdatePromotesWithoutLeavingDuplicate(t *testing.T) {
	t.Parallel()

	store := newStore(t)
	m := mustMemory(t, "workspace memory", "/repo/project")
	mustInsert(t, store, m)

	m.PromoteToGlobal()
	if err := store.Update(t.Context(), m); err != nil {
		t.Fatalf("Update: %v", err)
	}

	globalPath := store.memoryPath(m.ID(), "")
	workspacePath := store.memoryPath(m.ID(), "/repo/project")
	if _, err := os.Stat(globalPath); err != nil {
		t.Fatalf("stat global memory: %v", err)
	}
	if _, err := os.Stat(workspacePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("workspace memory still exists: %v", err)
	}
}

func newStore(t *testing.T) *FS {
	t.Helper()

	return newStoreAt(t, t.TempDir())
}

func newStoreAt(t *testing.T, dir string) *FS {
	t.Helper()

	store, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return store
}

func mustMemory(t *testing.T, content, workspace string) memory.Memory {
	t.Helper()

	m, err := memory.New(content, "", workspace)
	if err != nil {
		t.Fatalf("memory.New: %v", err)
	}
	return m
}

func mustInsert(t *testing.T, store *FS, m memory.Memory) {
	t.Helper()

	if err := store.Insert(t.Context(), m); err != nil {
		t.Fatalf("Insert: %v", err)
	}
}
