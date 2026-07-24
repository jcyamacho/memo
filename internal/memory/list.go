package memory

import (
	"context"
	"slices"
	"strings"
)

type ListResult struct {
	Workspace string
	Items     []MemoryDTO
}

func (s *Service) List(ctx context.Context, queryWorkspace string) (ListResult, error) {
	resolved, err := s.workspaceResolver.Resolve(ctx, queryWorkspace)
	if err != nil {
		return ListResult{}, err
	}

	items, err := s.store.List(ctx, resolved)
	if err != nil {
		return ListResult{}, err
	}

	slices.SortFunc(items, func(a, b Memory) int {
		if order := b.UpdatedAt().Compare(a.UpdatedAt()); order != 0 {
			return order
		}
		return strings.Compare(a.ID().String(), b.ID().String())
	})

	results := make([]MemoryDTO, 0, len(items))
	for i := range items {
		results = append(results, items[i].DTO())
	}

	return ListResult{Workspace: resolved, Items: results}, nil
}
