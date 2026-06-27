package memory

import "context"

func (s *Service) Workspaces(ctx context.Context) ([]string, error) {
	return s.store.ListWorkspaces(ctx)
}
