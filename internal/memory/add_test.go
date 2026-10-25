package memory

import (
	"errors"
	"strings"
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestServiceAddResolvesWorkspaceAndInserts(t *testing.T) {
	t.Parallel()

	store := newMockStore(t)
	resolver := newMockWorkspaceResolver("/repo")
	service := NewService(store, resolver)

	got, err := service.Add(t.Context(), AddParams{Content: " content ", Workspace: "."})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	testutil.AssertEqual(t, []string{"."}, resolver.paths)
	testutil.AssertEqual(t, 1, len(store.inserted))
	testutil.AssertEqual(t, "content", store.inserted[0].Content())
	testutil.AssertEqual(t, "/repo", store.inserted[0].Workspace())
	testutil.AssertEqual(t, "content", got.Content)
	testutil.AssertEqual(t, "/repo", got.Workspace)
}

func TestServiceAddGlobalSkipsWorkspaceResolver(t *testing.T) {
	t.Parallel()

	store := newMockStore(t)
	resolver := newMockWorkspaceResolver("/repo")
	service := NewService(store, resolver)

	got, err := service.Add(t.Context(), AddParams{Content: "content", Global: true})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	testutil.AssertEqual(t, []string(nil), resolver.paths)
	testutil.AssertEqual(t, 1, len(store.inserted))
	testutil.AssertEqual(t, "content", store.inserted[0].Content())
	testutil.AssertEqual(t, "", store.inserted[0].Workspace())
	testutil.AssertEqual(t, "content", got.Content)
	testutil.AssertEqual(t, "", got.Workspace)
}

func TestServiceAddUsesNameInID(t *testing.T) {
	t.Parallel()

	store := newMockStore(t)
	service := NewService(store, newMockWorkspaceResolver("/repo"))

	got, err := service.Add(t.Context(), AddParams{Content: "content", Name: "Prefer Bun", Global: true})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if !strings.HasPrefix(got.ID.String(), "prefer-bun-") {
		t.Fatalf("ID = %q, want prefix prefer-bun-", got.ID)
	}
	testutil.AssertEqual(t, 1, len(store.inserted))
	testutil.AssertEqual(t, got.ID.String(), store.inserted[0].ID().String())
}

func TestServiceAddRejectsInvalidName(t *testing.T) {
	t.Parallel()

	store := newMockStore(t)
	service := NewService(store, newMockWorkspaceResolver("/repo"))

	_, err := service.Add(t.Context(), AddParams{Content: "content", Name: "!!!", Global: true})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("Add = %v, want ErrInvalidInput", err)
	}
	testutil.AssertEqual(t, 0, len(store.inserted))
}
