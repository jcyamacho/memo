package cmd

import (
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func New(service *memory.Service, storeDir string) *cobra.Command {
	root := &cobra.Command{
		Use:           "memo",
		Short:         "Durable memory for humans and coding agents",
		Version:       currentVersion(),
		SilenceErrors: true,
		SilenceUsage:  true,
		Long: `memo stores durable facts, preferences, and project context.

Scope:
  Workspace memories apply to one project. By default, memo resolves the current
  directory to the Git repository root when possible.
  Global memories apply across projects. Pass --global to add, get, edit, or
  delete them.

Output:
  context, add, get, list, edit, and delete print XML intended to be easy for
  humans to inspect and stable for LLM tools to parse.

Agent rules:
  Use loaded memory context before deciding what to add, edit, or delete.
  Store durable facts, preferences, and project context. Do not store secrets.
  Use edit for corrections to an existing memory. Use delete only for obsolete
  or incorrect memories.

Store:
  Set MEMO_CONFIG_DIR to choose the store directory. The default is
  ~/.config/memo. Run memo path to print the directory in use.`,
		Example: `  memo list
  memo context
  memo add "prefer small Go interfaces"
  some-command | memo add --global
  memo edit 550e8400-e29b-41d4-a716-446655440000 --content "corrected memory"
  memo delete 550e8400-e29b-41d4-a716-446655440000`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.SetVersionTemplate("memo {{.Version}}\n")

	root.AddCommand(
		newAddCommand(service),
		newContextCommand(service),
		newDeleteCommand(service),
		newEditCommand(service),
		newGetCommand(service),
		newHookCommand(),
		newListCommand(service),
		newPathCommand(storeDir),
		newSkillCommand(),
		newVersionCommand(),
		newWorkspacesCommand(service),
	)

	return root
}
