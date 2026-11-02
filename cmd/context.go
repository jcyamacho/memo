package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

const (
	contextFormatXML        = "xml"
	contextFormatCursorHook = "cursor-hook"
)

const contextInstructions = "Durable memories are managed by the `memo` CLI; use the memories below as durable context for this session. " +
	"If the harness truncated this output, read the full output file it points to, or run `memo list`, before acting on the first request. " +
	"When the user corrects you, states a lasting preference, or asks you to remember something, or when you confirm a non-obvious project fact, consider saving a memory. " +
	"Before adding, editing, or deleting memories, run `memo skill` unless its rules are already in context."

func newContextCommand(service *memory.Service) *cobra.Command {
	var workspace string
	var format string

	command := &cobra.Command{
		Use:     "context",
		Aliases: []string{"ctx"},
		Short:   "Print session-start context for coding agents",
		Long: `Print durable memory context for a coding agent session.

The default output wraps an <instructions> block, a <workspace path="..."> block,
and a <global> block in <memo>. The memory blocks match memo list: newest first,
one "- <id> (<updated date>): <content>" line per memory. Pass --global to get,
edit, or delete a memory from <global>.
If a harness truncates hook output, the instructions tell the agent to read the
full output file or run memo list before acting on the first request. They also
tell the agent when to consider saving a memory.
Use --format cursor-hook to wrap the same XML in Cursor's sessionStart JSON
envelope ({"additional_context":"..."}).
By default, the workspace is the current directory resolved to the Git repository
root when possible. When CURSOR_PROJECT_DIR is set and --workspace is omitted,
that directory is used instead (Cursor user hooks run from ~/.cursor).`,
		Example: `  memo context
  memo context --workspace /path/to/project
  memo context --format cursor-hook`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspacePath, err := resolveContextWorkspacePath(workspace)
			if err != nil {
				return err
			}
			result, err := service.List(cmd.Context(), workspacePath)
			if err != nil {
				return err
			}

			output, err := formatContextOutput(contextOutput(result), format)
			if err != nil {
				return err
			}

			return printOutput(cmd, output)
		},
	}
	command.Flags().StringVar(&workspace, "workspace", "", "workspace path (default: cwd)")
	command.Flags().StringVar(&format, "format", contextFormatXML, "output format: xml or cursor-hook")

	return command
}

func contextOutput(result memory.ListResult) string {
	return "<memo>\n<instructions>" + contextInstructions + "</instructions>\n" +
		memoryBlocks(result) + "</memo>"
}

func formatContextOutput(xml string, format string) (string, error) {
	switch format {
	case "", contextFormatXML:
		return xml, nil
	case contextFormatCursorHook:
		payload, err := json.Marshal(struct {
			AdditionalContext string `json:"additional_context"`
		}{AdditionalContext: xml})
		if err != nil {
			return "", fmt.Errorf("encode cursor-hook output: %w", err)
		}
		return string(payload), nil
	default:
		return "", fmt.Errorf("unsupported format %q (want xml or cursor-hook)", format)
	}
}
