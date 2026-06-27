package memory

import (
	"context"
	"time"
)

type EditParams struct {
	ID            string
	Content       *string
	PromoteGlobal bool
}

func (s *Service) Edit(ctx context.Context, params EditParams) (Memory, error) {
	m, err := s.store.Get(ctx, params.ID)
	if err != nil {
		return Memory{}, err
	}

	changed := false
	if params.Content != nil {
		if err := m.SetContent(*params.Content); err != nil {
			return Memory{}, err
		}
		changed = true
	}
	if params.PromoteGlobal {
		m.PromoteToGlobal()
		changed = true
	}
	if !changed {
		return Memory{}, ErrInvalidInput
	}

	m.UpdatedAt = time.Now().UTC()
	if err := s.store.Update(ctx, m); err != nil {
		return Memory{}, err
	}
	return m, nil
}
