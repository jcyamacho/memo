package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/memory"
)

func TestMemoriesXMLEmpty(t *testing.T) {
	t.Parallel()

	got, err := memoriesXMLOutput("/repo", nil)
	if err != nil {
		t.Fatalf("memoriesXMLOutput: %v", err)
	}
	want := `<memories workspace="/repo"></memories>`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestMemoryXMLGlobal(t *testing.T) {
	t.Parallel()

	m, err := memory.Load(
		"abc123",
		"hello & world",
		"",
		time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("memory.Load: %v", err)
	}

	got, err := memoryXMLOutput(m.DTO())
	if err != nil {
		t.Fatalf("memoryXMLOutput: %v", err)
	}
	if !strings.Contains(got, `global="true"`) {
		t.Fatalf("missing global attribute: %s", got)
	}
	if !strings.Contains(got, "hello &amp; world") {
		t.Fatalf("content not escaped: %s", got)
	}
}

func TestDeleteResultsXML(t *testing.T) {
	t.Parallel()

	outcomes := []memory.DeleteOutcome{
		{Deleted: true, ID: "a"},
		{Deleted: false, ID: "b", Code: memory.DeleteCodeNotFound},
	}

	got, err := deleteResultsXMLOutput(outcomes)
	if err != nil {
		t.Fatalf("deleteResultsXMLOutput: %v", err)
	}
	if !strings.Contains(got, `<delete_results deleted="1" failed="1">`) {
		t.Fatalf("unexpected wrapper: %s", got)
	}
	if !strings.Contains(got, `<deleted id="a">`) {
		t.Fatalf("missing deleted id: %s", got)
	}
	if !strings.Contains(got, `<failure id="b" status="not_found">`) {
		t.Fatalf("missing failure: %s", got)
	}
}
