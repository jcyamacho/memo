package cmd

import (
	"fmt"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

var (
	addWorkspace string
	addGlobal    bool
)

func init() {
	addCmd.Flags().StringVar(&addWorkspace, "workspace", "", "project-scoped workspace path (default: cwd)")
	addCmd.Flags().BoolVar(&addGlobal, "global", false, "save as global memory (not project-scoped)")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
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
		if err := validateAddFlags(addGlobal, addWorkspace); err != nil {
			return err
		}
		content, err := readContent(args)
		if err != nil {
			return err
		}

		var workspacePath string
		if !addGlobal {
			workspacePath, err = resolveWorkspacePath(addWorkspace)
			if err != nil {
				return err
			}
		}

		m, err := service.Add(cmd.Context(), memory.AddParams{
			Content:   content,
			Workspace: workspacePath,
			Global:    addGlobal,
		})
		if err != nil {
			return err
		}
		return printOutput(memoryXMLOutput(m))
	},
}

func validateAddFlags(global bool, workspaceFlag string) error {
	if global && workspaceFlag != "" {
		return fmt.Errorf("--global and --workspace cannot be used together")
	}
	return nil
}
