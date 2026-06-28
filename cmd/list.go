package cmd

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

var listWorkspace string

func init() {
	listCmd.Flags().StringVar(&listWorkspace, "workspace", "", "workspace path (default: cwd)")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List memories for a workspace",
	Long: `List workspace-scoped and global memories, newest first.

By default, the workspace is the current directory resolved to the Git repository
root when possible. The result includes memories for that workspace plus global
memories. The command returns the full result set with no pagination as XML.`,
	Example: `  memo list
  memo ls
  memo list --workspace /path/to/project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workspacePath := listWorkspace
		if workspacePath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve working directory: %w", err)
			}
			workspacePath = cwd
		}

		items, err := service.List(cmd.Context(), workspacePath)
		if err != nil {
			return err
		}
		return printOutput(memoriesXMLOutput(workspacePath, items))
	},
}

type memoriesXML struct {
	XMLName   xml.Name    `xml:"memories"`
	Workspace string      `xml:"workspace,attr"`
	Items     []memoryXML `xml:"memory"`
}

func memoriesXMLOutput(workspace string, items []memory.MemoryDTO) (string, error) {
	dto := memoriesXML{
		Workspace: workspace,
		Items:     make([]memoryXML, 0, len(items)),
	}

	for _, item := range items {
		dto.Items = append(dto.Items, newMemoryXML(item, workspace))
	}

	return marshalXML(dto)
}
