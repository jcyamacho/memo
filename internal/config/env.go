package config

import (
	"os"
	"path/filepath"
)

const EnvDir = "MEMO_CONFIG_DIR"

func Dir() (string, error) {
	if dir := os.Getenv(EnvDir); dir != "" {
		return dir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "memo"), nil
}
