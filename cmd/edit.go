package cmd

import (
	"fmt"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newEditCommand(service *memory.Service) *cobra.Command {
	var contentFlag string
	var promoteGlobal bool

	command := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit an existing memory",
		Long: `Edit an existing memory by id.

Provide --content or piped stdin to update text. Use --global to promote a
workspace memory to global scope. At least one content source or --global is
required.

Use edit when a memory is still useful but needs correction. The command prints
the updated memory as XML.`,
		Example: `  memo edit 550e8400-e29b-41d4-a716-446655440000 --content "corrected fact"
  echo "corrected fact" | memo edit 550e8400-e29b-41d4-a716-446655440000
  memo edit 550e8400-e29b-41d4-a716-446655440000 --global
  echo "applies everywhere" | memo edit 550e8400-e29b-41d4-a716-446655440000 --global`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseMemoryID(args[0])
			if err != nil {
				return err
			}

			contentSet := cmd.Flags().Changed("content")
			stdin := stdinReader(cmd.InOrStdin())
			if !contentSet && stdin == nil && !promoteGlobal {
				return fmt.Errorf("provide --content, stdin, and/or --global")
			}

			params := memory.EditParams{
				ID:            id,
				PromoteGlobal: promoteGlobal,
			}
			if contentSet {
				params.Content = &contentFlag
			} else if stdin != nil {
				content, err := readContentFrom(nil, stdin)
				if err != nil {
					return err
				}
				params.Content = &content
			}
			m, err := service.Edit(cmd.Context(), params)
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
	command.Flags().StringVar(&contentFlag, "content", "", "new memory content")
	command.Flags().BoolVar(&promoteGlobal, "global", false, "promote workspace memory to global scope")

	return command
}
