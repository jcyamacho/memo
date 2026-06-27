package memory

import (
	"context"
)

type AddParams struct {
	Content string
	Global  bool
}

func (s *Service) Add(ctx context.Context, workspacePath string, params AddParams) (Memory, error) {
	scope := ""
	if !params.Global {
		resolved, err := s.workspaceResolver.Resolve(ctx, workspacePath)
		if err != nil {
			return Memory{}, err
		}
		scope = resolved
	}

	m, err := New(params.Content, scope)
	if err != nil {
		return Memory{}, err
	}
	if err := s.store.Insert(ctx, m); err != nil {
		return Memory{}, err
	}
	return m, nil
}
