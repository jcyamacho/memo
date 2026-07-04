package memory

import (
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestSetContentTouchesOnlyWhenContentChanges(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	m, err := Load("abc123", "old content", "/repo", updatedAt)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	changed, err := m.SetContent("old content")
	if err != nil {
		t.Fatalf("SetContent unchanged: %v", err)
	}
	testutil.AssertEqual(t, false, changed)
	testutil.AssertEqual(t, updatedAt, m.UpdatedAt())

	changed, err = m.SetContent("new content")
	if err != nil {
		t.Fatalf("SetContent changed: %v", err)
	}
	testutil.AssertEqual(t, true, changed)
	if !m.UpdatedAt().After(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want after %s", m.UpdatedAt(), updatedAt)
	}
}

func TestPromoteToGlobalTouchesOnlyWhenScopeChanges(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	m, err := Load("abc123", "content", "/repo", updatedAt)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	testutil.AssertEqual(t, true, m.PromoteToGlobal())
	testutil.AssertEqual(t, true, m.IsGlobal())
	if !m.UpdatedAt().After(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want after %s", m.UpdatedAt(), updatedAt)
	}

	promotedAt := m.UpdatedAt()
	testutil.AssertEqual(t, false, m.PromoteToGlobal())
	testutil.AssertEqual(t, promotedAt, m.UpdatedAt())
}
