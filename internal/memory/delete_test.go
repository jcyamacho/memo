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

	missingID := mustParseID(t, testIDTwo)
	failedID := mustParseID(t, testIDThree)
	store.deleteErrByID[missingID] = ErrNotFound
	store.deleteErrByID[failedID] = errInternal

	okID := mustParseID(t, testIDOne)
	got, err := service.Delete(t.Context(), []ID{okID, missingID, okID, failedID})
	if !errors.Is(err, errInternal) {
		t.Fatalf("Delete error = %v, want store error", err)
	}
	if len(got) != 3 {
		t.Fatalf("Delete outcomes = %d, want 3", len(got))
	}
	if !got[0].Deleted || got[0].ID != okID {
		t.Fatalf("successful outcome = %#v", got[0])
	}
	if got[1].ID != missingID || got[1].Code != DeleteCodeNotFound {
		t.Fatalf("missing outcome = %#v", got[1])
	}
	if got[2].ID != failedID || got[2].Code != DeleteCodeInternalError {
		t.Fatalf("failed outcome = %#v", got[2])
	}
	testutil.AssertEqual(t, []string{testIDOne, testIDTwo, testIDThree}, store.deletedIDs)
}
