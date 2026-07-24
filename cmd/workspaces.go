package cmd

import (
	"strings"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newWorkspacesCommand(service *memory.Service) *cobra.Command {
	return &cobra.Command{
		Use:     "workspaces",
		Aliases: []string{"ws"},
		Short:   "List workspace paths that have memories",
		Long: `List all workspace paths that have workspace-scoped memories.

Global memories are not included. The output is plain text with one absolute
workspace path per line.`,
		Example: `  memo workspaces`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaces, err := service.Workspaces(cmd.Context())
			if err != nil {
				return err
			}
			if len(workspaces) == 0 {
				return nil
			}
			return printOutput(cmd, strings.Join(workspaces, "\n"))
		},
	}
}
