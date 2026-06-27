package memory

import "context"

func (s *Service) Get(ctx context.Context, id string) (Memory, error) {
	return s.store.Get(ctx, id)
}
