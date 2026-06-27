package memory

import "context"

func (s *Service) List(ctx context.Context, queryWorkspace string) ([]Memory, error) {
	resolved, err := s.workspaceResolver.Resolve(ctx, queryWorkspace)
	if err != nil {
		return nil, err
	}

	items, err := s.store.List(ctx, ListFilter{
		Workspace: resolved,
		Global:    true,
	})
	if err != nil {
		return nil, err
	}

	for i := range items {
		if items[i].Workspace == resolved {
			items[i].Workspace = queryWorkspace
		}
	}
	return items, nil
}
