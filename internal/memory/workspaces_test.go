package memory

import (
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestServiceWorkspacesReturnsStoreWorkspaces(t *testing.T) {
	t.Parallel()

	store := newMockStore(t)
	service := NewService(store, nil)

	store.workspaces = []string{"/repo/one", "/repo/two"}

	got, err := service.Workspaces(t.Context())
	if err != nil {
		t.Fatalf("Workspaces: %v", err)
	}

	testutil.AssertEqual(t, []string{"/repo/one", "/repo/two"}, got)
}
