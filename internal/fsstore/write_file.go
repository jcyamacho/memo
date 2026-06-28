package fsstore

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func writeMemoryFile(path, content string, updatedAt time.Time) error {
	if err := os.MkdirAll(filepath.Dir(path), directoryPerm); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), filePerm); err != nil {
		return fmt.Errorf("write memory file: %w", err)
	}

	if err := os.Chtimes(path, updatedAt, updatedAt); err != nil {
		return fmt.Errorf("set file times: %w", err)
	}

	return nil
}
