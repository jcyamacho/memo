package memory

import (
	"context"
	"errors"
)

type DeleteCode string

const (
	DeleteCodeNotFound      DeleteCode = "not_found"
	DeleteCodeInternalError DeleteCode = "internal_error"
)

type DeleteOutcome struct {
	Deleted bool
	ID      ID
	Code    DeleteCode
}

type DeleteParams struct {
	IDs       []ID
	Workspace string
	Global    bool
}

func (s *Service) Delete(ctx context.Context, params DeleteParams) ([]DeleteOutcome, error) {
	workspace, err := s.resolveScope(ctx, params.Workspace, params.Global)
	if err != nil {
		return nil, err
	}

	seen := make(map[ID]struct{}, len(params.IDs))
	outcomes := make([]DeleteOutcome, 0, len(params.IDs))
	var deleteErr error
	for _, id := range params.IDs {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		err := s.store.Delete(ctx, id, workspace)
		switch {
		case err == nil:
			outcomes = append(outcomes, DeleteOutcome{Deleted: true, ID: id})
		case errors.Is(err, ErrNotFound):
			outcomes = append(outcomes, DeleteOutcome{Deleted: false, ID: id, Code: DeleteCodeNotFound})
		default:
			outcomes = append(outcomes, DeleteOutcome{Deleted: false, ID: id, Code: DeleteCodeInternalError})
			if deleteErr == nil {
				deleteErr = err
			}
		}
	}
	return outcomes, deleteErr
}
