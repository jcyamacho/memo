package memory

import "context"

type Store interface {
	Insert(ctx context.Context, m Memory) error
	Get(ctx context.Context, id string) (Memory, error)
	List(ctx context.Context, workspace string, options ...ListOption) ([]Memory, error)
	Update(ctx context.Context, m Memory) error
	Delete(ctx context.Context, id string) error
	ListWorkspaces(ctx context.Context) ([]string, error)
}

type ListOption func(*ListOptions)

type ListOptions struct {
	IncludeGlobals bool
}

func WithGlobals() ListOption {
	return func(options *ListOptions) {
		options.IncludeGlobals = true
	}
}

func NewListOptions(options ...ListOption) ListOptions {
	var result ListOptions
	for _, option := range options {
		option(&result)
	}
	return result
}
