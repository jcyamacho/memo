package cmd

import (
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

const contextInstructions = "Durable memories are managed by the `memo` CLI; use the memories below as durable context for this session. " +
	"If the harness truncated this output, read the full output file it points to, or run `memo list`, before acting on the first request. " +
	"When the user corrects you, states a lasting preference, or asks you to remember something, or when you confirm a non-obvious project fact, consider saving a memory. " +
	"Before adding, editing, or deleting memories, run `memo skill` unless its rules are already in context."

func newContextCommand(service *memory.Service) *cobra.Command {
	var workspace string

	command := &cobra.Command{
		Use:     "context",
		Aliases: []string{"ctx"},
		Short:   "Print session-start context for coding agents",
		Long: `Print durable memory context for a coding agent session.

The output wraps an <instructions> block, a <workspace path="..."> block, and a
<global> block in <memo>. The memory blocks match memo list: newest first, one
"- <id> (<updated date>): <content>" line per memory. Pass --global to get,
edit, or delete a memory from <global>.
If a harness truncates hook output, the instructions tell the agent to read the
full output file or run memo list before acting on the first request. They also
tell the agent when to consider saving a memory.
By default, the workspace is the current directory resolved to the Git repository
root when possible.`,
		Example: `  memo context
  memo context --workspace /path/to/project`,
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

			return printOutput(cmd, contextOutput(result))
		},
	}
	command.Flags().StringVar(&workspace, "workspace", "", "workspace path (default: cwd)")

	return command
}

func contextOutput(result memory.ListResult) string {
	return "<memo>\n<instructions>" + contextInstructions + "</instructions>\n" +
		memoryBlocks(result) + "</memo>"
}
