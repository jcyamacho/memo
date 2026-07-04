package memory

import (
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestServiceListSortsAndMapsResolvedWorkspace(t *testing.T) {
	t.Parallel()

	older := time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC)
	newer := older.Add(time.Hour)
	store := newMockStore(t)
	resolver := newMockWorkspaceResolver("/resolved/repo")
	service := NewService(store, resolver)

	store.listItems = []Memory{
		mustLoadMemory(t, "workspace-b", "workspace b", "/resolved/repo", older),
		mustLoadMemory(t, "global", "global", "", newer),
		mustLoadMemory(t, "workspace-a", "workspace a", "/resolved/repo", older),
	}

	got, err := service.List(t.Context(), ".")
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	want := []MemoryDTO{
		{ID: "global", Content: "global", UpdatedAt: newer},
		{ID: "workspace-a", Content: "workspace a", Workspace: ".", UpdatedAt: older},
		{ID: "workspace-b", Content: "workspace b", Workspace: ".", UpdatedAt: older},
	}
	testutil.AssertEqual(t, want, got)
	testutil.AssertEqual(t, []string{"."}, resolver.paths)
	testutil.AssertEqual(t, "/resolved/repo", store.listWorkspace)
	testutil.AssertEqual(t, true, store.listIncludeGlobals)
}
