package memory

import (
	"errors"
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestServiceEditUpdatesChangedMemory(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC)
	store := newMockStore(t)
	service := NewService(store, newMockWorkspaceResolver("/repo"))

	original := mustLoadMemory(t, testIDOne, "old content", "/repo", updatedAt)
	store.getByID[original.ID()] = original

	got, err := service.Edit(t.Context(), EditParams{
		ID:        original.ID(),
		Content:   "new content",
		Workspace: ".",
	})
	if err != nil {
		t.Fatalf("Edit: %v", err)
	}

	testutil.AssertEqual(t, []string{"/repo"}, store.getWorkspaces)
	testutil.AssertEqual(t, 1, len(store.updated))
	testutil.AssertEqual(t, original.ID().String(), store.updated[0].ID().String())
	testutil.AssertEqual(t, "new content", store.updated[0].Content())
	testutil.AssertEqual(t, "/repo", store.updated[0].Workspace())
	if !store.updated[0].UpdatedAt().After(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want after %s", store.updated[0].UpdatedAt(), updatedAt)
	}
	testutil.AssertEqual(t, original.ID().String(), got.ID.String())
	testutil.AssertEqual(t, "new content", got.Content)
	testutil.AssertEqual(t, "/repo", got.Workspace)
	if !got.UpdatedAt.After(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want after %s", got.UpdatedAt, updatedAt)
	}
}

func TestServiceEditRejectsNoop(t *testing.T) {
	t.Parallel()

	store := newMockStore(t)
	service := NewService(store, nil)

	original := mustLoadMemory(t, testIDOne, "content", "", time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC))
	store.getByID[original.ID()] = original

	_, err := service.Edit(t.Context(), EditParams{
		ID:      original.ID(),
		Content: "content",
		Global:  true,
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("Edit error = %v, want ErrInvalidInput", err)
	}
	testutil.AssertEqual(t, []string{""}, store.getWorkspaces)
	testutil.AssertEqual(t, 0, len(store.updated))
}
