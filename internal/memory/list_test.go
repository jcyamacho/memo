package memory

import (
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestServiceListSortsAndReturnsResolvedWorkspace(t *testing.T) {
	t.Parallel()

	older := time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC)
	newer := older.Add(time.Hour)
	store := newMockStore(t)
	resolver := newMockWorkspaceResolver("/resolved/repo")
	service := NewService(store, resolver)

	store.listItems = []Memory{
		mustLoadMemory(t, testIDTwo, "workspace b", "/resolved/repo", older),
		mustLoadMemory(t, testIDThree, "global", "", newer),
		mustLoadMemory(t, testIDOne, "workspace a", "/resolved/repo", older),
	}

	result, err := service.List(t.Context(), ".")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	testutil.AssertEqual(t, "/resolved/repo", result.Workspace)

	type comparableMemoryDTO struct {
		ID        string
		Content   string
		Workspace string
		UpdatedAt time.Time
	}

	comparable := make([]comparableMemoryDTO, 0, len(result.Items))
	for _, item := range result.Items {
		comparable = append(comparable, comparableMemoryDTO{
			ID:        item.ID.String(),
			Content:   item.Content,
			Workspace: item.Workspace,
			UpdatedAt: item.UpdatedAt,
		})
	}
	want := []comparableMemoryDTO{
		{ID: testIDThree, Content: "global", UpdatedAt: newer},
		{ID: testIDOne, Content: "workspace a", Workspace: "/resolved/repo", UpdatedAt: older},
		{ID: testIDTwo, Content: "workspace b", Workspace: "/resolved/repo", UpdatedAt: older},
	}
	testutil.AssertEqual(t, want, comparable)
	testutil.AssertEqual(t, []string{"."}, resolver.paths)
	testutil.AssertEqual(t, "/resolved/repo", store.listWorkspace)
}
