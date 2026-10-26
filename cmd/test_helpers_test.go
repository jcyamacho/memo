package cmd

import (
	"testing"
	"time"

	"github.com/jcyamacho/memo/internal/memory"
)

const (
	testMemoryIDOne = "550e8400-e29b-41d4-a716-446655440000"
	testMemoryIDTwo = "550e8400-e29b-41d4-a716-446655440001"
)

func mustLoadMemoryDTO(t *testing.T, id, content, workspace string, updatedAt time.Time) memory.MemoryDTO {
	t.Helper()

	m, err := memory.Load(id, content, workspace, updatedAt)
	if err != nil {
		t.Fatalf("memory.Load: %v", err)
	}
	return m.DTO()
}
