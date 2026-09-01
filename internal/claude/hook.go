package claude

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcyamacho/memo/internal/hookconfig"
)

const memoContextCommand = "memo context"

func DefaultHookPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, ".claude", "settings.json"), nil
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
	if err := hookconfig.ConfigureFile(path, hookconfig.SessionStartHook{
		Command:       memoContextCommand,
		StatusMessage: "Loading memories",
		Timeout:       5,
	}); err != nil {
		return fmt.Errorf("configure Claude hook: %w", err)
	}

	return nil
}
