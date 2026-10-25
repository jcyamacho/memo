package memory

import (
	"context"
)

type AddParams struct {
	Content   string
	Name      string
	Workspace string
	Global    bool
}

func (s *Service) Add(ctx context.Context, params AddParams) (MemoryDTO, error) {
	workspace, err := s.resolveScope(ctx, params.Workspace, params.Global)
	if err != nil {
		return MemoryDTO{}, err
	}

	m, err := New(params.Content, params.Name, workspace)
	if err != nil {
		return MemoryDTO{}, err
	}
	if err := s.store.Insert(ctx, m); err != nil {
		return MemoryDTO{}, err
	}
	return m.DTO(), nil
}
