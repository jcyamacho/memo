package jsonutil

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func Read[T any](path string) (T, error) {
	var value T

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return value, nil
	}
	if err != nil {
		return value, err
	}

	if err := json.Unmarshal(data, &value); err != nil {
		return value, err
	}
	return value, nil
}

func Write[T any](path string, value T) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
