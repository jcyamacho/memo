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

	m := mustLoadMemory(t, testIDOne, "content", "/repo", updatedAt)
	store.getByID[m.ID()] = m

	got, err := service.Get(t.Context(), m.ID())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	testutil.AssertEqual(t, testIDOne, got.ID.String())
	testutil.AssertEqual(t, "content", got.Content)
	testutil.AssertEqual(t, "/repo", got.Workspace)
	testutil.AssertEqual(t, updatedAt, got.UpdatedAt)
}
