package memory

import (
	"context"
	"errors"
)

type DeleteOutcome struct {
	Deleted bool
	ID      string
	Code    string
}

func (s *Service) Delete(ctx context.Context, ids []string) ([]DeleteOutcome, error) {
	seen := make(map[string]struct{}, len(ids))
	outcomes := make([]DeleteOutcome, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		err := s.store.Delete(ctx, id)
		switch {
		case err == nil:
			outcomes = append(outcomes, DeleteOutcome{Deleted: true, ID: id})
		case errors.Is(err, ErrNotFound):
			outcomes = append(outcomes, DeleteOutcome{Deleted: false, ID: id, Code: "not_found"})
		default:
			outcomes = append(outcomes, DeleteOutcome{Deleted: false, ID: id, Code: "internal_error"})
		}
	}
	return outcomes, nil
}
