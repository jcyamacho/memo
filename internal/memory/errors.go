package memory

import "errors"

var (
	ErrNotFound     = errors.New("memory not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrExists       = errors.New("memory already exists")
)
