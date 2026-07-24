package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/memory"
)

const (
	testMemoryIDOne = "550e8400-e29b-41d4-a716-446655440000"
	testMemoryIDTwo = "550e8400-e29b-41d4-a716-446655440001"
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
		testMemoryIDOne,
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

func TestContextXML(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	workspaceMemory, err := memory.Load(testMemoryIDOne, "workspace fact", "/repo", updatedAt)
	if err != nil {
		t.Fatalf("memory.Load workspace: %v", err)
	}
	globalMemory, err := memory.Load(testMemoryIDTwo, "global fact", "", updatedAt)
	if err != nil {
		t.Fatalf("memory.Load global: %v", err)
	}

	got, err := contextXMLOutput("/repo", []memory.MemoryDTO{
		workspaceMemory.DTO(),
		globalMemory.DTO(),
	})
	if err != nil {
		t.Fatalf("contextXMLOutput: %v", err)
	}

	if !strings.Contains(got, `<memo_context workspace="/repo">`) {
		t.Fatalf("missing context workspace: %s", got)
	}
	if !strings.Contains(got, "Durable memories are managed by the `memo` CLI") {
		t.Fatalf("missing durable memory instruction: %s", got)
	}
	if !strings.Contains(got, "run `memo skill` unless its rules are already in context") {
		t.Fatalf("missing memo skill instruction: %s", got)
	}
	if strings.Contains(got, "<instruction>") {
		t.Fatalf("instructions should be a single text node: %s", got)
	}
	if count := strings.Count(got, `workspace="/repo"`); count != 1 {
		t.Fatalf("workspace attribute count = %d, want 1: %s", count, got)
	}
	if !strings.Contains(got, `<memory id="550e8400-e29b-41d4-a716-446655440000" updated_at="2026-03-01T12:00:00.000Z">workspace fact</memory>`) {
		t.Fatalf("missing workspace memory: %s", got)
	}
	if !strings.Contains(got, `<memory id="550e8400-e29b-41d4-a716-446655440001" updated_at="2026-03-01T12:00:00.000Z" global="true">global fact</memory>`) {
		t.Fatalf("missing global memory: %s", got)
	}
}

func TestDeleteResultsXML(t *testing.T) {
	t.Parallel()

	deletedID, err := memory.ParseID(testMemoryIDOne)
	if err != nil {
		t.Fatalf("ParseID deleted: %v", err)
	}
	missingID, err := memory.ParseID(testMemoryIDTwo)
	if err != nil {
		t.Fatalf("ParseID missing: %v", err)
	}

	outcomes := []memory.DeleteOutcome{
		{Deleted: true, ID: deletedID},
		{Deleted: false, ID: missingID, Code: memory.DeleteCodeNotFound},
	}

	got, err := deleteResultsXMLOutput(outcomes)
	if err != nil {
		t.Fatalf("deleteResultsXMLOutput: %v", err)
	}
	if !strings.Contains(got, `<delete_results deleted="1" failed="1">`) {
		t.Fatalf("unexpected wrapper: %s", got)
	}
	if !strings.Contains(got, `<deleted id="550e8400-e29b-41d4-a716-446655440000">`) {
		t.Fatalf("missing deleted id: %s", got)
	}
	if !strings.Contains(got, `<failure id="550e8400-e29b-41d4-a716-446655440001" status="not_found">`) {
		t.Fatalf("missing failure: %s", got)
	}
}
