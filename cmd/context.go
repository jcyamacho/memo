package cmd

import (
	"encoding/xml"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

const contextInstructions = "Durable memories are managed by the `memo` CLI; use the memories below as durable context for this session. " +
	"Before adding, editing, or deleting memories, run `memo skill` unless its rules are already in context."

var contextWorkspace string

func init() {
	contextCmd.Flags().StringVar(&contextWorkspace, "workspace", "", "workspace path (default: cwd)")
	rootCmd.AddCommand(contextCmd)
}

var contextCmd = &cobra.Command{
	Use:     "context",
	Aliases: []string{"ctx"},
	Short:   "Print session-start context for coding agents",
	Long: `Print durable memory context for a coding agent session.

The output includes a short instruction block plus workspace and global memories.
By default, the workspace is the current directory resolved to the Git repository
root when possible.`,
	Example: `  memo context
  memo context --workspace /path/to/project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workspacePath, err := resolveWorkspacePath(contextWorkspace)
		if err != nil {
			return err
		}

		items, err := service.List(cmd.Context(), workspacePath)
		if err != nil {
			return err
		}
		return printOutput(contextXMLOutput(workspacePath, items))
	},
}

type memoContextXML struct {
	XMLName      xml.Name           `xml:"memo_context"`
	Workspace    string             `xml:"workspace,attr"`
	Instructions string             `xml:"instructions"`
	Memories     contextMemoriesXML `xml:"memories"`
}

type contextMemoriesXML struct {
	Items []memoryXML `xml:"memory"`
}

func contextXMLOutput(workspace string, items []memory.MemoryDTO) (string, error) {
	dto := memoContextXML{
		Workspace:    workspace,
		Instructions: contextInstructions,
		Memories:     contextMemoriesXML{Items: make([]memoryXML, 0, len(items))},
	}

	for _, item := range items {
		dto.Memories.Items = append(dto.Memories.Items, newMemoryXML(item, workspace))
	}

	return marshalXML(dto)
}
