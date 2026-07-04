package memory

import (
	"errors"
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestServiceDeleteDeduplicatesAndClassifiesOutcomes(t *testing.T) {
	t.Parallel()

	store := newMockStore(t)
	service := NewService(store, nil)
	errInternal := errors.New("store failed")

	store.deleteErrByID["missing"] = ErrNotFound
	store.deleteErrByID["failed"] = errInternal

	got, err := service.Delete(t.Context(), []string{"ok", "missing", "ok", "failed"})
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	want := []DeleteOutcome{
		{Deleted: true, ID: "ok"},
		{ID: "missing", Code: DeleteCodeNotFound},
		{ID: "failed", Code: DeleteCodeInternalError},
	}
	testutil.AssertEqual(t, want, got)
	testutil.AssertEqual(t, []string{"ok", "missing", "failed"}, store.deletedIDs)
}
