package workspace

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitResolver struct{}

func (GitResolver) Resolve(ctx context.Context, path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", errRequired
	}
	absolute, err := filepath.Abs(trimmed)
	if err != nil {
		return "", fmt.Errorf("resolve absolute workspace path: %w", err)
	}

	cmd := exec.CommandContext(
		ctx,
		"git",
		"rev-parse",
		"--path-format=absolute",
		"--show-toplevel",
		"--git-dir",
		"--git-common-dir",
	)
	cmd.Dir = absolute
	output, err := cmd.Output()
	if err != nil {
		return absolute, nil
	}

	paths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(paths) != 3 {
		return absolute, nil
	}
	if paths[1] != paths[2] {
		return filepath.Dir(paths[2]), nil
	}
	return paths[0], nil
}
