package cmd

import (
	"github.com/jcyamacho/memo/internal/claude"
	"github.com/jcyamacho/memo/internal/codex"
	"github.com/spf13/cobra"
)

func init() {
	hookCmd.AddCommand(hookCodexCmd, hookClaudeCmd)
	rootCmd.AddCommand(hookCmd)
}

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Configure session-start hooks",
	Long: `Configure session-start hooks for coding-agent harnesses.

The configured hook runs memo context when a new agent session starts. Existing
memo context hooks are replaced, and unrelated hooks are preserved.`,
	Example: `  memo hook codex
  memo hook claude`,
}

var hookCodexCmd = &cobra.Command{
	Use:   "codex",
	Short: "Configure the Codex session-start hook",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := codex.ConfigureHook()
		if err != nil {
			return err
		}
		return printOutput("configured memo hook in "+target, nil)
	},
}

var hookClaudeCmd = &cobra.Command{
	Use:   "claude",
	Short: "Configure the Claude session-start hook",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := claude.ConfigureHook()
		if err != nil {
			return err
		}
		return printOutput("configured memo hook in "+target, nil)
	},
}
