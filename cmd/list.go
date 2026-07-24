package cmd

import (
	"encoding/xml"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newListCommand(service *memory.Service) *cobra.Command {
	var workspace string

	command := &cobra.Command{
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

			output, err := memoriesXMLOutput(result.Workspace, result.Items)
			if err != nil {
				return err
			}
			return printOutput(cmd, output)
		},
	}
	command.Flags().StringVar(&workspace, "workspace", "", "workspace path (default: cwd)")

	return command
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
