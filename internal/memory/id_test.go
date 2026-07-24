package memory

import (
	"errors"
	"testing"
)

const (
	testIDOne   = "550e8400-e29b-41d4-a716-446655440000"
	testIDTwo   = "550e8400-e29b-41d4-a716-446655440001"
	testIDThree = "550e8400-e29b-41d4-a716-446655440002"
)

func TestNewIDCreatesValidID(t *testing.T) {
	t.Parallel()

	id, err := NewID()
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}

	parsed, err := ParseID(id.String())
	if err != nil {
		t.Fatalf("ParseID generated id: %v", err)
	}
	if parsed != id {
		t.Fatalf("ParseID generated id = %q, want %q", parsed, id)
	}
}

func TestParseIDAcceptsUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "version 1", value: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", want: "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
		{name: "version 4", value: testIDOne, want: testIDOne},
		{name: "version 7", value: "019f4b1f-f34a-7182-9051-652396fc3159", want: "019f4b1f-f34a-7182-9051-652396fc3159"},
		{name: "normalized", value: "550E8400-E29B-41D4-A716-446655440000", want: testIDOne},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, err := ParseID(test.value)
			if err != nil {
				t.Fatalf("ParseID: %v", err)
			}
			if id.String() != test.want {
				t.Fatalf("String = %q, want %q", id.String(), test.want)
			}
		})
	}
}

func TestParseIDRejectsInvalidID(t *testing.T) {
	t.Parallel()

	tests := []string{
		"",
		"abc123",
		"../outside",
		"00000000-0000-0000-0000-000000000000",
	}

	for _, value := range tests {
		t.Run(value, func(t *testing.T) {
			_, err := ParseID(value)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("ParseID(%q) = %v, want ErrInvalidInput", value, err)
			}
		})
	}
}
