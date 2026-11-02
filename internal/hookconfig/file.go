package hookconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func ConfigureFile(path string, spec SessionStartHook) error {
	return Update(path, func(config *Config) error {
		if err := config.ConfigureSessionStart(spec); err != nil {
			return fmt.Errorf("configure session-start hook: %w", err)
		}
		return nil
	})
}

// Update reads a hook config, applies update, and writes it back.
// A missing file starts empty. Invalid JSON is returned unread.
func Update(path string, update func(*Config) error) error {
	config, err := readFile(path)
	if err != nil {
		return fmt.Errorf("read hook config: %w", err)
	}

	if err := update(&config); err != nil {
		return err
	}

	if err := writeFile(path, config); err != nil {
		return fmt.Errorf("write hook config: %w", err)
	}

	return nil
}

func readFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func writeFile(path string, config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
