package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workspacesCmd)
}

var workspacesCmd = &cobra.Command{
	Use:     "workspaces",
	Aliases: []string{"ws"},
	Short:   "List workspace paths that have memories",
	Long: `List all workspace paths that have workspace-scoped memories.

Global memories are not included. The output is plain text with one absolute
	workspace path per line.`,
	Example: `  memo workspaces`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaces, err := service.Workspaces(cmd.Context())
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(workspaces, "\n"))
		return nil
	},
}
