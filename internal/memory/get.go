package memory

import "context"

func (s *Service) Get(ctx context.Context, id string) (MemoryDTO, error) {
	m, err := s.store.Get(ctx, id)
	if err != nil {
		return MemoryDTO{}, err
	}
	return m.DTO(), nil
}
