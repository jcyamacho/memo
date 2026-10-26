package memory

import (
	"context"
	"slices"
	"strings"
)

type ListResult struct {
	Workspace      string
	WorkspaceItems []MemoryDTO
	GlobalItems    []MemoryDTO
}

func (s *Service) List(ctx context.Context, queryWorkspace string) (ListResult, error) {
	resolved, err := s.workspaceResolver.Resolve(ctx, queryWorkspace)
	if err != nil {
		return ListResult{}, err
	}

	workspaceItems, err := s.listScope(ctx, resolved)
	if err != nil {
		return ListResult{}, err
	}

	globalItems, err := s.listScope(ctx, "")
	if err != nil {
		return ListResult{}, err
	}

	return ListResult{Workspace: resolved, WorkspaceItems: workspaceItems, GlobalItems: globalItems}, nil
}

func (s *Service) listScope(ctx context.Context, workspace string) ([]MemoryDTO, error) {
	items, err := s.store.List(ctx, workspace)
	if err != nil {
		return nil, err
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

	return results, nil
}
