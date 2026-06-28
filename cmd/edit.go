package cmd

import (
	"fmt"
	"os"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

var (
	editContent string
	editGlobal  bool
)

func init() {
	editCmd.Flags().StringVar(&editContent, "content", "", "new memory content")
	editCmd.Flags().BoolVar(&editGlobal, "global", false, "promote workspace memory to global scope")
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing memory",
	Long: `Edit an existing memory by id.

Provide --content or piped stdin to update text. Use --global to promote a
workspace memory to global scope. At least one content source or --global is
required.

Use edit when a memory is still useful but needs correction. The command prints
the updated memory as XML.`,
	Example: `  memo edit abc123 --content "corrected fact"
  echo "corrected fact" | memo edit abc123
  memo edit abc123 --global
  echo "applies everywhere" | memo edit abc123 --global`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contentSet := cmd.Flags().Changed("content")
		stdin := stdinReader(os.Stdin)
		if !contentSet && stdin == nil && !editGlobal {
			return fmt.Errorf("provide --content, stdin, and/or --global")
		}

		params := memory.EditParams{
			ID:            args[0],
			PromoteGlobal: editGlobal,
		}
		if contentSet {
			params.Content = &editContent
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
		return printOutput(memoryXMLOutput(m))
	},
}
