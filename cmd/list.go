package cmd

import (
	"strings"
	"time"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newListCommand(service *memory.Service) *cobra.Command {
	var workspace string

	command := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List memories for a workspace",
		Long: `List workspace-scoped and global memories.

The output wraps a <workspace path="..."> block and a <global> block in
<memories>, newest first within each block. Each memory is one line,
"- <id> (<updated date>): <content>", and continuation lines of multi-line
content are indented by two spaces. Pass --global to get, edit, or delete a
memory from <global>.

By default, the workspace is the current directory resolved to the Git repository
root when possible. The result includes memories for that workspace plus global
memories. The command returns the full result set with no pagination.`,
		Example: `  memo list
  memo ls
  memo list --workspace /path/to/project`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspacePath, err := resolveWorkspacePath(workspace)
			if err != nil {
				return err
			}
			result, err := service.List(cmd.Context(), workspacePath)
			if err != nil {
				return err
			}

			return printOutput(cmd, memoriesOutput(result))
		},
	}
	command.Flags().StringVar(&workspace, "workspace", "", "workspace path (default: cwd)")

	return command
}

func memoriesOutput(result memory.ListResult) string {
	return "<memories>\n" + memoryBlocks(result) + "</memories>"
}

// memoryBlocks puts workspace memories before globals so they survive harnesses
// that keep only the start of long output.
func memoryBlocks(result memory.ListResult) string {
	return `<workspace path="` + result.Workspace + "\">\n" +
		memoryLines(result.WorkspaceItems) +
		"</workspace>\n<global>\n" +
		memoryLines(result.GlobalItems) +
		"</global>\n"
}

func memoryLines(items []memory.MemoryDTO) string {
	var lines strings.Builder
	for _, item := range items {
		lines.WriteString(memoryLine(item))
		lines.WriteByte('\n')
	}

	return lines.String()
}

// memoryLine indents continuation lines so multi-line content stays in one list item.
func memoryLine(m memory.MemoryDTO) string {
	return "- " + m.ID.String() + " (" + m.UpdatedAt.UTC().Format(time.DateOnly) + "): " +
		strings.ReplaceAll(m.Content, "\n", "\n  ")
}
