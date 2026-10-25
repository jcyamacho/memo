package memory

import (
	"errors"
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestServiceGetReturnsMemoryDTO(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC)
	store := newMockStore(t)
	resolver := newMockWorkspaceResolver("/repo")
	service := NewService(store, resolver)

	m := mustLoadMemory(t, testIDOne, "content", "/repo", updatedAt)
	store.getByID[m.ID()] = m

	got, err := service.Get(t.Context(), GetParams{ID: m.ID(), Workspace: "."})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	testutil.AssertEqual(t, []string{"."}, resolver.paths)
	testutil.AssertEqual(t, []string{"/repo"}, store.getWorkspaces)
	testutil.AssertEqual(t, testIDOne, got.ID.String())
	testutil.AssertEqual(t, "content", got.Content)
	testutil.AssertEqual(t, "/repo", got.Workspace)
	testutil.AssertEqual(t, updatedAt, got.UpdatedAt)
}

func TestServiceGetGlobalSkipsWorkspaceResolver(t *testing.T) {
	t.Parallel()

	store := newMockStore(t)
	resolver := newMockWorkspaceResolver("/repo")
	service := NewService(store, resolver)

	m := mustLoadMemory(t, testIDOne, "content", "", time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC))
	store.getByID[m.ID()] = m

	if _, err := service.Get(t.Context(), GetParams{ID: m.ID(), Global: true}); err != nil {
		t.Fatalf("Get: %v", err)
	}

	testutil.AssertEqual(t, []string(nil), resolver.paths)
	testutil.AssertEqual(t, []string{""}, store.getWorkspaces)
}

func TestServiceGetNotFoundNamesScope(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params GetParams
		want   string
	}{
		{name: "workspace", params: GetParams{Workspace: "."}, want: "memory not found in workspace /repo"},
		{name: "global", params: GetParams{Global: true}, want: "memory not found in globals"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := newMockStore(t)
			service := NewService(store, newMockWorkspaceResolver("/repo"))
			id := mustParseID(t, testIDOne)
			store.getErrByID[id] = ErrNotFound

			test.params.ID = id
			_, err := service.Get(t.Context(), test.params)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("Get error = %v, want ErrNotFound", err)
			}
			testutil.AssertEqual(t, test.want, err.Error())
		})
	}
}
