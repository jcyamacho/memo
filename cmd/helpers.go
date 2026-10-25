package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func readContent(cmd *cobra.Command, args []string) (string, error) {
	return readContentFrom(args, stdinReader(cmd.InOrStdin()))
}

func readContentFrom(args []string, stdin io.Reader) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	if stdin == nil {
		return "", fmt.Errorf("content is required")
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", fmt.Errorf("read stdin: %w", err)
	}
	return string(data), nil
}

func stdinReader(stdin io.Reader) io.Reader {
	file, ok := stdin.(*os.File)
	if !ok {
		return stdin
	}

	info, err := file.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice != 0 {
		return nil
	}

	return stdin
}

func printOutput(cmd *cobra.Command, output string) error {
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), output); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func parseMemoryID(value string) (memory.ID, error) {
	id, err := memory.ParseID(value)
	if err != nil {
		return memory.ID{}, fmt.Errorf("memory id %q: %w", value, err)
	}

	return id, nil
}

// scopeWorkspacePath returns the workspace path for commands that accept
// --global and --workspace; it is empty for the global scope.
func scopeWorkspacePath(global bool, workspaceFlag string) (string, error) {
	if !global {
		return resolveWorkspacePath(workspaceFlag)
	}
	if workspaceFlag != "" {
		return "", fmt.Errorf("--global and --workspace cannot be used together")
	}
	return "", nil
}

func resolveWorkspacePath(workspaceFlag string) (string, error) {
	if workspaceFlag != "" {
		return workspaceFlag, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	return cwd, nil
}
