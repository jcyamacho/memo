package memory

import (
	"context"
	"sort"
)

func (s *Service) List(ctx context.Context, queryWorkspace string) ([]Memory, error) {
	resolved, err := s.workspaceResolver.Resolve(ctx, queryWorkspace)
	if err != nil {
		return nil, err
	}

	items, err := s.store.List(ctx, ListFilter{
		Workspace:      resolved,
		IncludeGlobals: true,
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	for i := range items {
		if items[i].Workspace == resolved {
			items[i].Workspace = queryWorkspace
		}
	}
	return items, nil
}
