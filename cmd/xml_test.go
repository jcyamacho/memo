package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/memory"
)

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
