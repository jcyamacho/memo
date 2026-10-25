package memory

import "context"

type GetParams struct {
	ID        ID
	Workspace string
	Global    bool
}

func (s *Service) Get(ctx context.Context, params GetParams) (MemoryDTO, error) {
	workspace, err := s.resolveScope(ctx, params.Workspace, params.Global)
	if err != nil {
		return MemoryDTO{}, err
	}

	m, err := s.store.Get(ctx, params.ID, workspace)
	if err != nil {
		return MemoryDTO{}, scopeNotFound(err, workspace)
	}
	return m.DTO(), nil
}
