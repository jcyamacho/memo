package memory

import (
	"context"
)

type AddParams struct {
	Content   string
	Workspace string
	Global    bool
}

func (s *Service) Add(ctx context.Context, params AddParams) (MemoryDTO, error) {
	resolvedWorkspace := ""
	if !params.Global {
		resolved, err := s.workspaceResolver.Resolve(ctx, params.Workspace)
		if err != nil {
			return MemoryDTO{}, err
		}
		resolvedWorkspace = resolved
	}

	m, err := New(params.Content, resolvedWorkspace)
	if err != nil {
		return MemoryDTO{}, err
	}
	if err := s.store.Insert(ctx, m); err != nil {
		return MemoryDTO{}, err
	}
	return m.DTO(), nil
}
