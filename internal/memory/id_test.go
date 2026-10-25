package memory

import (
	"encoding/hex"
	"errors"
	"strings"
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

func TestNewNamedIDNormalizesName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		prefix string
	}{
		{name: "case and punctuation", input: "Prefer Bun for JS!", prefix: "prefer-bun-for-js-"},
		{name: "separator runs", input: "  deploy__branch / release ", prefix: "deploy-branch-release-"},
		{name: "non ascii", input: "configuración", prefix: "configuraci-n-"},
		{name: "truncated", input: strings.Repeat("a", 47) + " bcd", prefix: strings.Repeat("a", 47) + "-"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewNamedID(test.input)
			if err != nil {
				t.Fatalf("NewNamedID: %v", err)
			}

			suffix, ok := strings.CutPrefix(id.String(), test.prefix)
			if !ok {
				t.Fatalf("id = %q, want prefix %q", id, test.prefix)
			}
			if _, err := hex.DecodeString(suffix); err != nil || len(suffix) != 6 {
				t.Fatalf("suffix = %q, want six hex characters", suffix)
			}

			parsed, err := ParseID(id.String())
			if err != nil {
				t.Fatalf("ParseID named id: %v", err)
			}
			if parsed != id {
				t.Fatalf("ParseID named id = %q, want %q", parsed, id)
			}
		})
	}
}

func TestNewNamedIDRejectsEmptyName(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"", " - ", "ñ"} {
		t.Run(value, func(t *testing.T) {
			_, err := NewNamedID(value)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("NewNamedID(%q) = %v, want ErrInvalidInput", value, err)
			}
		})
	}
}

func TestParseIDAcceptsValidID(t *testing.T) {
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
		{name: "named", value: "prefer-bun-3f9a2c", want: "prefer-bun-3f9a2c"},
		{name: "single word", value: "abc123", want: "abc123"},
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
		"../outside",
		"a/b",
		"Upper-Case",
		"a--b",
		"-a",
		"a-",
		strings.Repeat("a", 65),
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
