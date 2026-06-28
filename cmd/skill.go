package cmd

import (
	"fmt"

	"github.com/jcyamacho/memo/internal/skill"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(skillCmd)
}

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Print LLM instructions for operating memo",
	Long: `Print a self-contained Markdown guide for LLMs that need to use memo.

The guide explains command mapping, scope rules, output format, and conservative
memory maintenance policy. It does not read or write the memory store.`,
	Example: `  memo skill`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(skill.Guide())
		return nil
	},
}
