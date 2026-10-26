package cmd

import (
	"fmt"

	"github.com/jcyamacho/memo/internal/skill"
	"github.com/spf13/cobra"
)

func newSkillCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "skill",
		Short: "Print LLM instructions for operating memo",
		Long: `Print a self-contained Markdown guide for LLMs that need to use memo.

The guide explains command mapping, scope rules, how to write effective
memories, output format, and conservative memory maintenance policy. It does not read or write the memory store.`,
		Example: `  memo skill`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := fmt.Fprint(cmd.OutOrStdout(), skill.Guide()); err != nil {
				return fmt.Errorf("write output: %w", err)
			}

			return nil
		},
	}
}
