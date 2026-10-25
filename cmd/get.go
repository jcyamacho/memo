package cmd

import (
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newGetCommand(service *memory.Service) *cobra.Command {
	var workspace string
	var global bool

	command := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a memory by id",
		Long: `Fetch a single memory by id and print it as XML.

The memory must belong to the current workspace. Use --global for a global
memory, or --workspace for a memory in another workspace.

Use get when a caller already has an id and needs the full memory content.`,
		Example: `  memo get 550e8400-e29b-41d4-a716-446655440000
  memo get --global prefer-bun-3f9a2c
  memo get --workspace /path/to/project deploy-branch-1a2b3c`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspacePath, err := scopeWorkspacePath(global, workspace)
			if err != nil {
				return err
			}
			id, err := parseMemoryID(args[0])
			if err != nil {
				return err
			}
			m, err := service.Get(cmd.Context(), memory.GetParams{
				ID:        id,
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
	command.Flags().StringVar(&workspace, "workspace", "", "workspace path (default: cwd)")
	command.Flags().BoolVar(&global, "global", false, "get a global memory")

	return command
}
