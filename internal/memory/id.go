package memory

import "github.com/google/uuid"

type ID struct {
	value uuid.UUID
}

func NewID() (ID, error) {
	value, err := uuid.NewRandom()
	if err != nil {
		return ID{}, err
	}
	return ID{value: value}, nil
}

func ParseID(value string) (ID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil || parsed == uuid.Nil {
		return ID{}, ErrInvalidInput
	}
	return ID{value: parsed}, nil
}

func (id ID) String() string {
	return id.value.String()
}
