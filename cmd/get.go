package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a memory by id",
	Long: `Fetch a single memory by id and print it as XML.

Use get when a caller already has an id and needs the full memory content.`,
	Example: `  memo get abc123`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := service.Get(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return printOutput(memoryXMLOutput(m))
	},
}
