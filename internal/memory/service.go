package memory

import (
	"context"
	"errors"
	"fmt"
)

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

// resolveScope returns the store workspace for a request, where "" is the
// global scope.
func (s *Service) resolveScope(ctx context.Context, workspace string, global bool) (string, error) {
	if global {
		return "", nil
	}
	return s.workspaceResolver.Resolve(ctx, workspace)
}

func scopeNotFound(err error, workspace string) error {
	if !errors.Is(err, ErrNotFound) {
		return err
	}
	if workspace == "" {
		return fmt.Errorf("%w in globals", err)
	}
	return fmt.Errorf("%w in workspace %s", err, workspace)
}
