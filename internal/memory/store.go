package memory

import "context"

type Store interface {
	Insert(ctx context.Context, m Memory) error
	Get(ctx context.Context, id ID) (Memory, error)
	List(ctx context.Context, workspace string) ([]Memory, error)
	Update(ctx context.Context, m Memory) error
	Delete(ctx context.Context, id ID) error
	ListWorkspaces(ctx context.Context) ([]string, error)
}
