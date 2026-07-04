package memory

import (
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestServiceGetReturnsMemoryDTO(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC)
	store := newMockStore(t)
	service := NewService(store, nil)

	m := mustLoadMemory(t, "abc123", "content", "/repo", updatedAt)
	store.getByID["abc123"] = m

	got, err := service.Get(t.Context(), "abc123")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	want := MemoryDTO{
		ID:        "abc123",
		Content:   "content",
		Workspace: "/repo",
		UpdatedAt: updatedAt,
	}
	testutil.AssertEqual(t, want, got)
}
