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
	service := NewService(store, nil)
	content := "new content"

	original := mustLoadMemory(t, "abc123", "old content", "/repo", updatedAt)
	store.getByID["abc123"] = original

	got, err := service.Edit(t.Context(), EditParams{
		ID:            "abc123",
		Content:       &content,
		PromoteGlobal: true,
	})
	if err != nil {
		t.Fatalf("Edit: %v", err)
	}

	testutil.AssertEqual(t, 1, len(store.updated))
	testutil.AssertEqual(t, "abc123", store.updated[0].ID())
	testutil.AssertEqual(t, "new content", store.updated[0].Content())
	testutil.AssertEqual(t, "", store.updated[0].Workspace())
	if !store.updated[0].UpdatedAt().After(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want after %s", store.updated[0].UpdatedAt(), updatedAt)
	}
	testutil.AssertEqual(t, "abc123", got.ID)
	testutil.AssertEqual(t, "new content", got.Content)
	testutil.AssertEqual(t, "", got.Workspace)
	if !got.UpdatedAt.After(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want after %s", got.UpdatedAt, updatedAt)
	}
}

func TestServiceEditRejectsNoop(t *testing.T) {
	t.Parallel()

	store := newMockStore(t)
	service := NewService(store, nil)
	content := "content"

	original := mustLoadMemory(t, "abc123", "content", "", time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC))
	store.getByID["abc123"] = original

	_, err := service.Edit(t.Context(), EditParams{
		ID:            "abc123",
		Content:       &content,
		PromoteGlobal: true,
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("Edit error = %v, want ErrInvalidInput", err)
	}
	testutil.AssertEqual(t, 0, len(store.updated))
}
