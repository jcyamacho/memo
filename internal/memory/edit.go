package memory

import "context"

type EditParams struct {
	ID        ID
	Content   string
	Workspace string
	Global    bool
}

func (s *Service) Edit(ctx context.Context, params EditParams) (MemoryDTO, error) {
	workspace, err := s.resolveScope(ctx, params.Workspace, params.Global)
	if err != nil {
		return MemoryDTO{}, err
	}

	m, err := s.store.Get(ctx, params.ID, workspace)
	if err != nil {
		return MemoryDTO{}, scopeNotFound(err, workspace)
	}

	changed, err := m.SetContent(params.Content)
	if err != nil {
		return MemoryDTO{}, err
	}
	if !changed {
		return MemoryDTO{}, ErrInvalidInput
	}

	if err := s.store.Update(ctx, m); err != nil {
		return MemoryDTO{}, err
	}
	return m.DTO(), nil
}
