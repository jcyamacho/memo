package memory

import (
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestServiceListGroupsScopesAndSortsNewestFirst(t *testing.T) {
	t.Parallel()

	older := time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC)
	newer := older.Add(time.Hour)
	store := newMockStore(t)
	resolver := newMockWorkspaceResolver("/resolved/repo")
	service := NewService(store, resolver)

	store.listItemsByWorkspace = map[string][]Memory{
		"/resolved/repo": {
			mustLoadMemory(t, testIDTwo, "workspace b", "/resolved/repo", older),
			mustLoadMemory(t, testIDOne, "workspace a", "/resolved/repo", older),
		},
		"": {
			mustLoadMemory(t, testIDOne, "older global", "", older),
			mustLoadMemory(t, testIDThree, "newer global", "", newer),
		},
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

	toComparable := func(items []MemoryDTO) []comparableMemoryDTO {
		comparable := make([]comparableMemoryDTO, 0, len(items))
		for _, item := range items {
			comparable = append(comparable, comparableMemoryDTO{
				ID:        item.ID.String(),
				Content:   item.Content,
				Workspace: item.Workspace,
				UpdatedAt: item.UpdatedAt,
			})
		}
		return comparable
	}

	testutil.AssertEqual(t, []comparableMemoryDTO{
		{ID: testIDOne, Content: "workspace a", Workspace: "/resolved/repo", UpdatedAt: older},
		{ID: testIDTwo, Content: "workspace b", Workspace: "/resolved/repo", UpdatedAt: older},
	}, toComparable(result.WorkspaceItems))
	testutil.AssertEqual(t, []comparableMemoryDTO{
		{ID: testIDThree, Content: "newer global", UpdatedAt: newer},
		{ID: testIDOne, Content: "older global", UpdatedAt: older},
	}, toComparable(result.GlobalItems))
	testutil.AssertEqual(t, []string{"."}, resolver.paths)
	testutil.AssertEqual(t, []string{"/resolved/repo", ""}, store.listedWorkspaces)
}
