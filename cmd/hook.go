package cmd

import (
	"github.com/jcyamacho/memo/internal/claude"
	"github.com/jcyamacho/memo/internal/codex"
	"github.com/jcyamacho/memo/internal/cursor"
	"github.com/spf13/cobra"
)

func newHookCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "hook",
		Short: "Configure session-start hooks",
		Long: `Configure session-start hooks for coding-agent harnesses.

The configured hook runs memo context when a new agent session starts. Existing
memo context hooks are replaced, and unrelated hooks are preserved.`,
		Example: `  memo hook codex
  memo hook claude
  memo hook cursor`,
	}
	command.AddCommand(newHookCodexCommand(), newHookClaudeCommand(), newHookCursorCommand())

	return command
}

func newHookCodexCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "codex",
		Short: "Configure the Codex session-start hook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := codex.ConfigureHook()
			if err != nil {
				return err
			}

			return printOutput(cmd, "configured memo hook in "+target+"\nopen /hooks in Codex to review and trust the hook")
		},
	}
}

func newHookClaudeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "claude",
		Short: "Configure the Claude session-start hook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := claude.ConfigureHook()
			if err != nil {
				return err
			}

			return printOutput(cmd, "configured memo hook in "+target)
		},
	}
}

func newHookCursorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cursor",
		Short: "Configure the Cursor session-start hook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := cursor.ConfigureHook()
			if err != nil {
				return err
			}

			return printOutput(cmd, "configured memo hook in "+target+"\nreview the hook in Cursor Settings > Hooks; memo must be on PATH")
		},
	}
}
