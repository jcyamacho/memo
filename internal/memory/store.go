package memory

import "context"

type Store interface {
	Insert(ctx context.Context, m Memory) error
	Get(ctx context.Context, id string) (Memory, error)
	List(ctx context.Context, filter ListFilter) ([]Memory, error)
	Update(ctx context.Context, m Memory) error
	Delete(ctx context.Context, id string) error
	ListWorkspaces(ctx context.Context) ([]string, error)
}

type ListFilter struct {
	Workspace string
	Global    bool
}
