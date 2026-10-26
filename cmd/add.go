package cmd

import (
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newAddCommand(service *memory.Service) *cobra.Command {
	var workspace string
	var name string
	var global bool

	command := &cobra.Command{
		Use:   "add [content]",
		Short: "Add a new memory",
		Long: `Add a new durable memory.

Provide content as a positional argument or through stdin. By default the memory
is scoped to the current workspace. Use --global for facts that apply across projects.

Agents should pass content through a quoted heredoc (<<'EOF'), which keeps every
character, including quotes, backticks, and $(...). Inside double quotes the shell
runs backticks and $(...) and would save their output in the memory. Without a
content argument, add reads stdin: it fails when stdin is a terminal and can wait
for input when a harness leaves stdin open.

Use --name with a short summary of two to five ASCII words. It is lowercased,
every run of characters other than a-z and 0-9 becomes one dash, and it is cut
to 48 characters. A random six-character hex suffix is appended, and the result
becomes the memory id, for example "prefer-bun-3f9a2c". Without --name the id is
a UUID.

The command prints the created memory as XML, including its id. Empty or
whitespace-only content is rejected after trimming.`,
		Example: `  memo add --name "deploy branch" <<'EOF'
  Deploys must run from the release branch, not from main.
  EOF
  memo add --global --name "concise answers" 'prefer concise final answers'
  memo add --workspace /path/to/project --name "project fact" 'project-specific fact'`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspacePath, err := scopeWorkspacePath(global, workspace)
			if err != nil {
				return err
			}
			content, err := readContent(cmd, args)
			if err != nil {
				return err
			}

			m, err := service.Add(cmd.Context(), memory.AddParams{
				Content:   content,
				Name:      name,
				Workspace: workspacePath,
				Global:    global,
			})
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
	command.Flags().StringVar(&workspace, "workspace", "", "project-scoped workspace path (default: cwd)")
	command.Flags().StringVar(&name, "name", "", "short name used as the readable part of the memory id")
	command.Flags().BoolVar(&global, "global", false, "save as global memory (not project-scoped)")

	return command
}
