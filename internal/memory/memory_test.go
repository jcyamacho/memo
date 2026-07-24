package memory

import (
	"errors"
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestMemoryRejectsEmptyContent(t *testing.T) {
	t.Parallel()

	if _, err := New(" \n\t", "/repo"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("New error = %v, want ErrInvalidInput", err)
	}

	m, err := Load(testIDOne, "content", "/repo", time.Now())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if _, err := m.SetContent(" \n\t"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("SetContent error = %v, want ErrInvalidInput", err)
	}
	testutil.AssertEqual(t, "content", m.Content())
}

func TestSetContentTouchesOnlyWhenContentChanges(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	m, err := Load(testIDOne, "old content", "/repo", updatedAt)
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
	m, err := Load(testIDOne, "content", "/repo", updatedAt)
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
