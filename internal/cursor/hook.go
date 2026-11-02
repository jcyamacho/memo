package cursor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcyamacho/memo/internal/hookconfig"
)

const memoContextCommand = "memo context --format cursor-hook"

func DefaultHookPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, ".cursor", "hooks.json"), nil
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
	if err := hookconfig.Update(path, configureSessionStart); err != nil {
		return fmt.Errorf("configure Cursor hook: %w", err)
	}

	return nil
}
