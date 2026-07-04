package testutil

import (
	"encoding/json"
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/jsonutil"
)

func DecodeJSON[T any](t *testing.T, content string) T {
	t.Helper()

	var value T
	if err := json.Unmarshal([]byte(content), &value); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	return value
}

func EncodeJSON[T any](t *testing.T, value any) T {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var encoded T
	if err := json.Unmarshal(data, &encoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	return encoded
}

func ReadJSON[T any](t *testing.T, path string) T {
	t.Helper()

	value, err := jsonutil.Read[T](path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	return value
}
