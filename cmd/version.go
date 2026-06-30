package cmd

import "github.com/spf13/cobra"

var version = "dev"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print release version information",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return printOutput(versionOutput(version))
	},
}

func versionOutput(version string) (string, error) {
	return "memo " + version, nil
}
