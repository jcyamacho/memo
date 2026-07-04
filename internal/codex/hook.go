package codex

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcyamacho/memo/internal/pkg/jsonutil"
)

func DefaultHookPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, ".codex", "hooks.json"), nil
}

func ConfigureHook() (string, error) {
	path, err := DefaultHookPath()
	if err != nil {
		return "", err
	}

	if err := configureHookAt(path); err != nil {
		return "", err
	}
	return path, nil
}

func configureHookAt(path string) error {
	config, err := jsonutil.Read[config](path)
	if err != nil {
		return fmt.Errorf("read Codex hook config: %w", err)
	}

	config.configureSessionStartHook()
	if err := jsonutil.Write(path, config); err != nil {
		return fmt.Errorf("write Codex hook config: %w", err)
	}
	return nil
}
