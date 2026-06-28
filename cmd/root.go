package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "memo",
	Short: "Durable memory for humans and coding agents",
	Long: `memo stores durable facts, preferences, and project context.

Scope:
  Workspace memories apply to one project. By default, memo resolves the current
  directory to the Git repository root when possible.
  Global memories apply across projects and are created with --global.

Output:
  add, get, list, edit, and delete print XML intended to be easy for humans to
  inspect and stable for LLM tools to parse.

Agent rules:
  Use list before deciding what to add, edit, or delete.
  Store durable facts, preferences, and project context. Do not store secrets.
  Use edit for corrections to an existing memory. Use delete only for obsolete
  or incorrect memories.

Store:
  Set MEMO_CONFIG_DIR to choose the store directory. The default is
  ~/.config/memo.`,
	Example: `  memo list
  memo add "prefer small Go interfaces"
  some-command | memo add --global
  memo edit abc123 --content "corrected memory"
  memo delete abc123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
