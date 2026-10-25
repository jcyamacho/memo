package cmd

import (
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newGetCommand(service *memory.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a memory by id",
		Long: `Fetch a single memory by id and print it as XML.

Use get when a caller already has an id and needs the full memory content.`,
		Example: `  memo get 550e8400-e29b-41d4-a716-446655440000`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseMemoryID(args[0])
			if err != nil {
				return err
			}
			m, err := service.Get(cmd.Context(), id)
			if err != nil {
				return err
			}

			output, err := memoryXMLOutput(m)
			if err != nil {
				return err
			}
			return printOutput(cmd, output)
		},
	}
}
