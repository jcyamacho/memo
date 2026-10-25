package cmd

import "github.com/spf13/cobra"

func newPathCommand(storeDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the memory store directory",
		Long: `Print the absolute path of the directory where memo stores memories.

The directory is MEMO_CONFIG_DIR when it is set, otherwise ~/.config/memo. The
output is plain text with one path.`,
		Example: `  memo path
  cd "$(memo path)"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return printOutput(cmd, storeDir)
		},
	}
}
