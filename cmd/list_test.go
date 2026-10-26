package cmd

import (
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestMemoriesOutputEmpty(t *testing.T) {
	t.Parallel()

	got := memoriesOutput(memory.ListResult{Workspace: "/repo"})

	testutil.AssertEqual(t, `<memories>
<workspace path="/repo">
</workspace>
<global>
</global>
</memories>`, got)
}

func TestMemoriesOutputGroupsScopesAsRawLines(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	got := memoriesOutput(memory.ListResult{
		Workspace: "/repo",
		WorkspaceItems: []memory.MemoryDTO{
			mustLoadMemoryDTO(t, testMemoryIDOne, "workspace fact\nsecond line", "/repo", updatedAt),
		},
		GlobalItems: []memory.MemoryDTO{
			mustLoadMemoryDTO(t, testMemoryIDTwo, `global "fact" & it's <raw>`, "", updatedAt),
		},
	})

	testutil.AssertEqual(t, `<memories>
<workspace path="/repo">
- 550e8400-e29b-41d4-a716-446655440000 (2026-03-01): workspace fact
  second line
</workspace>
<global>
- 550e8400-e29b-41d4-a716-446655440001 (2026-03-01): global "fact" & it's <raw>
</global>
</memories>`, got)
}
