package memory

import "context"

type WorkspaceResolver interface {
	Resolve(ctx context.Context, path string) (string, error)
}

type Service struct {
	store             Store
	workspaceResolver WorkspaceResolver
}

func NewService(store Store, workspaceResolver WorkspaceResolver) *Service {
	return &Service{
		store:             store,
		workspaceResolver: workspaceResolver,
	}
}
