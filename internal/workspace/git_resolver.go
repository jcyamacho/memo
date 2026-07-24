package workspace

import (
	"context"
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

	cmd := exec.CommandContext(
		ctx,
		"git",
		"rev-parse",
		"--path-format=absolute",
		"--show-toplevel",
		"--git-dir",
		"--git-common-dir",
	)
	cmd.Dir = trimmed
	output, err := cmd.Output()
	if err != nil {
		return trimmed, nil
	}

	paths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(paths) != 3 {
		return trimmed, nil
	}
	if paths[1] != paths[2] {
		return filepath.Dir(paths[2]), nil
	}
	return paths[0], nil
}
