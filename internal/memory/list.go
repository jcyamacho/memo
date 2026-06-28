package memory

import (
	"context"
	"sort"
)

func (s *Service) List(ctx context.Context, queryWorkspace string) ([]MemoryDTO, error) {
	resolved, err := s.workspaceResolver.Resolve(ctx, queryWorkspace)
	if err != nil {
		return nil, err
	}

	items, err := s.store.List(ctx, resolved, WithGlobals())
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].UpdatedAt().Equal(items[j].UpdatedAt()) {
			return items[i].ID() < items[j].ID()
		}
		return items[i].UpdatedAt().After(items[j].UpdatedAt())
	})

	results := make([]MemoryDTO, 0, len(items))
	for i := range items {
		item := items[i].DTO()
		if item.Workspace == resolved {
			item.Workspace = queryWorkspace
		}
		results = append(results, item)
	}
	return results, nil
}
