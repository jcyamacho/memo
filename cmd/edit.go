package cmd

import (
	"fmt"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newEditCommand(service *memory.Service) *cobra.Command {
	var contentFlag string
	var workspace string
	var global bool

	command := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit an existing memory",
		Long: `Replace the content of an existing memory by id.

Provide the new content with --content or piped stdin. The memory must belong to
the current workspace. Use --global for a global memory, or --workspace for a
memory in another workspace.

Edit cannot change a memory's scope. To move a memory, add it in the new scope,
then delete the old one.

Use edit when a memory is still useful but needs correction. The command prints
the updated memory as XML.`,
		Example: `  memo edit 550e8400-e29b-41d4-a716-446655440000 --content "corrected fact"
  echo "corrected fact" | memo edit 550e8400-e29b-41d4-a716-446655440000
  memo edit --global prefer-bun-3f9a2c --content "corrected preference"`,
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

			content := contentFlag
			if !cmd.Flags().Changed("content") {
				stdin := stdinReader(cmd.InOrStdin())
				if stdin == nil {
					return fmt.Errorf("provide --content or stdin")
				}
				content, err = readContentFrom(nil, stdin)
				if err != nil {
					return err
				}
			}

			m, err := service.Edit(cmd.Context(), memory.EditParams{
				ID:        id,
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
	command.Flags().StringVar(&contentFlag, "content", "", "new memory content")
	command.Flags().StringVar(&workspace, "workspace", "", "workspace path (default: cwd)")
	command.Flags().BoolVar(&global, "global", false, "edit a global memory")

	return command
}
