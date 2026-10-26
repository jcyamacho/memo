package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/memory"
)

func TestContextOutput(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	got := contextOutput(memory.ListResult{
		Workspace:      "/repo",
		WorkspaceItems: []memory.MemoryDTO{mustLoadMemoryDTO(t, testMemoryIDOne, "workspace fact", "/repo", updatedAt)},
		GlobalItems:    []memory.MemoryDTO{mustLoadMemoryDTO(t, testMemoryIDTwo, "global fact", "", updatedAt)},
	})

	if !strings.HasPrefix(got, "<memo>\n<instructions>") {
		t.Fatalf("missing memo and instructions opening: %s", got)
	}
	for _, instruction := range []string{
		"Durable memories are managed by the `memo` CLI",
		"If the harness truncated this output, read the full output file it points to, or run `memo list`",
		"consider saving a memory",
		"run `memo skill` unless its rules are already in context",
	} {
		if !strings.Contains(got, instruction) {
			t.Fatalf("missing instruction %q: %s", instruction, got)
		}
	}

	wantMemories := `</instructions>
<workspace path="/repo">
- 550e8400-e29b-41d4-a716-446655440000 (2026-03-01): workspace fact
</workspace>
<global>
- 550e8400-e29b-41d4-a716-446655440001 (2026-03-01): global fact
</global>
</memo>`
	if !strings.HasSuffix(got, wantMemories) {
		t.Fatalf("memories = %s, want suffix %s", got, wantMemories)
	}
}
