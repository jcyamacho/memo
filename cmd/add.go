package cmd

import (
	"fmt"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newAddCommand(service *memory.Service) *cobra.Command {
	var workspace string
	var global bool

	command := &cobra.Command{
		Use:   "add [content]",
		Short: "Add a new memory",
		Long: `Add a new durable memory.

Provide content as a positional argument or through stdin. By default the memory
is scoped to the current workspace. Use --global for facts that apply across projects.

The command prints the created memory as XML. Empty or whitespace-only content is
rejected after trimming.`,
		Example: `  memo add "this project uses Cobra for the CLI"
  memo add --global "prefer concise final answers"
  git diff --stat | memo add
  memo add --workspace /path/to/project "project-specific fact"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateAddFlags(global, workspace); err != nil {
				return err
			}
			content, err := readContent(cmd, args)
			if err != nil {
				return err
			}

			var workspacePath string
			if !global {
				workspacePath, err = resolveWorkspacePath(workspace)
				if err != nil {
					return err
				}
			}
			m, err := service.Add(cmd.Context(), memory.AddParams{
				Content:   content,
				Workspace: workspacePath,
				Global:    global,
			})
			if err != nil {
				return err
			}

			output, err := memoryXMLOutput(m)
			if err != nil {
				return err
			}
			return printOutput(cmd, output)
		},
	}
	command.Flags().StringVar(&workspace, "workspace", "", "project-scoped workspace path (default: cwd)")
	command.Flags().BoolVar(&global, "global", false, "save as global memory (not project-scoped)")

	return command
}

func validateAddFlags(global bool, workspaceFlag string) error {
	if global && workspaceFlag != "" {
		return fmt.Errorf("--global and --workspace cannot be used together")
	}
	return nil
}
