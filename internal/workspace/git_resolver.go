package workspace

import (
	"context"
	"os/exec"
	"strings"
)

type GitResolver struct{}

func (GitResolver) Resolve(ctx context.Context, path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", errRequired
	}

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--path-format=absolute", "--show-toplevel")
	cmd.Dir = trimmed
	output, err := cmd.Output()
	if err != nil {
		return trimmed, nil
	}

	return strings.TrimSpace(string(output)), nil
}
