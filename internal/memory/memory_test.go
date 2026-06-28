package memory

import (
	"testing"
	"time"
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
	if changed {
		t.Fatal("SetContent unchanged reported a change")
	}
	if !m.UpdatedAt().Equal(updatedAt) {
		t.Fatalf("UpdatedAt changed to %s, want %s", m.UpdatedAt(), updatedAt)
	}

	changed, err = m.SetContent("new content")
	if err != nil {
		t.Fatalf("SetContent changed: %v", err)
	}
	if !changed {
		t.Fatal("SetContent changed reported no change")
	}
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

	if !m.PromoteToGlobal() {
		t.Fatal("PromoteToGlobal reported no change")
	}
	if !m.IsGlobal() {
		t.Fatal("memory is not global")
	}
	if !m.UpdatedAt().After(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want after %s", m.UpdatedAt(), updatedAt)
	}

	promotedAt := m.UpdatedAt()
	if m.PromoteToGlobal() {
		t.Fatal("PromoteToGlobal on global memory reported a change")
	}
	if !m.UpdatedAt().Equal(promotedAt) {
		t.Fatalf("UpdatedAt changed to %s, want %s", m.UpdatedAt(), promotedAt)
	}
}
