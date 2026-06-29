package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print release version information",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return printOutput(versionOutput(version, commit, date))
	},
}

func versionOutput(version, commit, date string) (string, error) {
	return fmt.Sprintf("memo %s (commit %s, built %s)", version, commit, date), nil
}
