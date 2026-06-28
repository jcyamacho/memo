package memory

import "context"

type EditParams struct {
	ID            string
	Content       *string
	PromoteGlobal bool
}

func (s *Service) Edit(ctx context.Context, params EditParams) (MemoryDTO, error) {
	m, err := s.store.Get(ctx, params.ID)
	if err != nil {
		return MemoryDTO{}, err
	}

	changed := false
	if params.Content != nil {
		contentChanged, err := m.SetContent(*params.Content)
		if err != nil {
			return MemoryDTO{}, err
		}
		changed = changed || contentChanged
	}
	if params.PromoteGlobal {
		changed = m.PromoteToGlobal() || changed
	}
	if !changed {
		return MemoryDTO{}, ErrInvalidInput
	}

	if err := s.store.Update(ctx, m); err != nil {
		return MemoryDTO{}, err
	}
	return m.DTO(), nil
}
