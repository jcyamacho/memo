package cmd

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func newBrowseCommand(service *memory.Service) *cobra.Command {
	var workspace string

	command := &cobra.Command{
		Use:   "browse",
		Short: "Search, read, and manage memories in the terminal",
		Long: `Browse memories for the current workspace and the global scope.
Use --workspace to open a different project.

Results show memory IDs grouped by scope. The detail view shows the scope,
update date, and full ID above the memory text.

Press / to search IDs and full memory text. Use the arrow keys or j and k
to select a memory, and g or G to jump to the first or last one.

In a wide terminal the details stay next to the results. Press PgUp or PgDn,
or Ctrl+U or Ctrl+D, to scroll them. In a narrow terminal, press Enter or Tab
to open the details and Esc to return to the results.

Press e to edit in the terminal editor set by VISUAL or EDITOR. Save and
exit the editor to apply the changes. Press d to review the ID and scope
before deleting. Press y to confirm or n or Esc to cancel. Press q to quit.`,
		Example: "  memo browse\n  memo browse --workspace /path/to/project",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			in, inputOK := cmd.InOrStdin().(*os.File)
			out, outputOK := cmd.OutOrStdout().(*os.File)
			if !inputOK || !outputOK || !term.IsTerminal(in.Fd()) || !term.IsTerminal(out.Fd()) {
				return errors.New("memo browse requires an interactive terminal")
			}

			workspacePath, err := resolveWorkspacePath(workspace)
			if err != nil {
				return err
			}
			ctx := cmd.Context()
			result, err := service.List(ctx, workspacePath)
			if err != nil {
				return err
			}

			model := newBrowseModel(ctx, service, result)
			defer model.removeEditFile()

			program := tea.NewProgram(model, tea.WithContext(ctx), tea.WithInput(in), tea.WithOutput(out))
			_, err = program.Run()
			return err
		},
	}
	command.Flags().StringVar(&workspace, "workspace", "", "workspace path (default: cwd)")

	return command
}

type browseModel struct {
	ctx       context.Context
	service   *memory.Service
	workspace string
	items     []memory.MemoryDTO
	matches   []int
	selected  int
	search    textinput.Model
	query     *regexp.Regexp
	detail    viewport.Model

	width, height int
	reading       bool
	confirming    bool
	status        browseStatus
	editFile      string
}

type browseStatus struct {
	text   string
	failed bool
}

func browseNotice(text string) browseStatus {
	return browseStatus{text: text}
}

func browseFailure(format string, args ...any) browseStatus {
	return browseStatus{text: fmt.Sprintf(format, args...), failed: true}
}

func newBrowseModel(ctx context.Context, service *memory.Service, result memory.ListResult) *browseModel {
	search := textinput.New()
	search.Prompt = "Search: "
	search.Placeholder = "ID or content"
	search.SetVirtualCursor(false)

	detail := viewport.New()
	detail.SoftWrap = true
	detail.HighlightStyle = browseHighlight
	detail.SelectedHighlightStyle = browseHighlight

	items := slices.Concat(result.WorkspaceItems, result.GlobalItems)
	model := &browseModel{
		ctx: ctx, service: service, workspace: result.Workspace, items: items,
		search: search, detail: detail,
	}
	model.filter()
	return model
}

func (m *browseModel) Init() tea.Cmd { return nil }

func (m *browseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.resize()
	case browseEditorDone:
		m.finishEdit(msg.err)
	case tea.KeyPressMsg:
		return m, m.handleKey(msg)
	case tea.PasteMsg:
		if m.search.Focused() {
			return m, m.updateSearch(msg)
		}
	}

	return m, nil
}

func (m *browseModel) handleKey(key tea.KeyPressMsg) tea.Cmd {
	stroke := key.Keystroke()
	if stroke == "ctrl+c" {
		return tea.Quit
	}
	if m.editFile != "" {
		return nil
	}
	if m.tooSmall() {
		if stroke == "q" && !m.search.Focused() {
			return tea.Quit
		}
		return nil
	}
	if m.confirming {
		switch stroke {
		case "y":
			m.confirming = false
			m.deleteSelected()
		case "n", "esc":
			m.confirming = false
		}
		return nil
	}
	if m.search.Focused() {
		switch stroke {
		case "enter":
			m.search.Blur()
			m.reading = len(m.matches) > 0 && !m.splitView()
		case "tab":
			m.search.Blur()
		case "esc":
			m.search.Blur()
			m.clearSearch()
		case "up", "down":
			m.moveSelection(stroke)
		default:
			return m.updateSearch(key)
		}
		return nil
	}

	m.status = browseStatus{}
	switch stroke {
	case "q":
		return tea.Quit
	case "/":
		m.reading = false
		return m.search.Focus()
	case "esc":
		if m.readerOpen() {
			m.reading = false
		} else {
			m.clearSearch()
		}
	case "enter", "tab":
		if len(m.matches) > 0 && !m.splitView() {
			m.reading = !m.reading
		}
	case "e":
		if len(m.matches) > 0 {
			return m.openEditor()
		}
	case "d":
		m.confirming = len(m.matches) > 0
	default:
		if m.readerOpen() || m.splitView() && browseDetailScrollKey(stroke) {
			var cmd tea.Cmd
			m.detail, cmd = m.detail.Update(key)
			return cmd
		}
		m.moveSelection(stroke)
	}

	return nil
}

// The reader replaces the list only when the terminal is too narrow to show both.
func (m *browseModel) readerOpen() bool {
	return m.reading && !m.splitView()
}

func browseDetailScrollKey(stroke string) bool {
	switch stroke {
	case "pgup", "pgdown", "ctrl+u", "ctrl+d":
		return true
	}
	return false
}

func (m *browseModel) updateSearch(msg tea.Msg) tea.Cmd {
	previous := m.search.Value()
	var cmd tea.Cmd
	m.search, cmd = m.search.Update(msg)
	if m.search.Value() != previous {
		m.selected = 0
		m.filter()
	}
	return cmd
}

func (m *browseModel) clearSearch() {
	m.search.Reset()
	m.selected = 0
	m.filter()
}

func (m *browseModel) filter() {
	var query *regexp.Regexp
	if value := strings.TrimSpace(m.search.Value()); value != "" {
		query = regexp.MustCompile("(?i)" + regexp.QuoteMeta(value))
	}

	m.query = query
	m.matches = m.matches[:0]
	for index, item := range m.items {
		if query == nil || query.MatchString(item.ID.String()) || query.MatchString(item.Content) {
			m.matches = append(m.matches, index)
		}
	}

	m.selected = min(m.selected, max(0, len(m.matches)-1))
	if len(m.matches) == 0 {
		m.reading = false
	}
	m.showDetail()
}

func (m *browseModel) moveSelection(key string) {
	next := m.selected
	switch key {
	case "up", "k":
		next--
	case "down", "j":
		next++
	case "pgup":
		next -= m.listWindow().size()
	case "pgdown":
		next += m.listWindow().size()
	case "home", "g":
		next = 0
	case "end", "G", "shift+g":
		next = len(m.matches) - 1
	}
	next = max(0, min(next, len(m.matches)-1))
	if next != m.selected {
		m.selected = next
		m.showDetail()
	}
}

func (m *browseModel) selectedItem() memory.MemoryDTO {
	return m.items[m.matches[m.selected]]
}

type browseEditorDone struct{ err error }

func (m *browseModel) openEditor() tea.Cmd {
	editor := cmp.Or(strings.TrimSpace(os.Getenv("VISUAL")), strings.TrimSpace(os.Getenv("EDITOR")))
	if editor == "" {
		m.status = browseFailure("Set VISUAL or EDITOR to a terminal editor.")
		return nil
	}

	file, err := os.CreateTemp("", "memo-*.md")
	if err != nil {
		m.status = browseFailure("Could not prepare editor: %v", err)
		return nil
	}
	m.editFile = file.Name()
	item := m.selectedItem()
	_, writeErr := fmt.Fprintln(file, item.Content)
	if err := errors.Join(writeErr, file.Close()); err != nil {
		m.removeEditFile()
		m.status = browseFailure("Could not prepare editor: %v", err)
		return nil
	}

	// The editor setting can include shell arguments. Pass the file as a separate argument.
	command := exec.CommandContext(m.ctx, "sh", "-c", editor+` "$1"`, "memo browse", m.editFile)
	return tea.ExecProcess(command, func(err error) tea.Msg {
		return browseEditorDone{err: err}
	})
}

func (m *browseModel) finishEdit(editorErr error) {
	defer m.removeEditFile()
	if editorErr != nil {
		m.status = browseFailure("Editor exited without saving to memo: %v", editorErr)
		return
	}

	data, err := os.ReadFile(m.editFile)
	if err != nil {
		m.status = browseFailure("Could not read edited memory: %v", err)
		return
	}

	index := m.matches[m.selected]
	item := m.items[index]
	content := strings.TrimSpace(string(data))
	switch content {
	case item.Content:
		m.status = browseNotice("No changes.")
		return
	case "":
		m.status = browseFailure("Memory text is empty. No changes saved.")
		return
	}

	current, err := m.service.Get(m.ctx, memory.GetParams{
		ID: item.ID, Workspace: item.Workspace, Global: item.IsGlobal(),
	})
	if err != nil {
		m.status = browseFailure("Could not read memory before saving: %v", err)
		return
	}
	if current.Content != item.Content {
		m.items[index] = current
		m.filter()
		m.status = browseFailure("Memory changed outside the browser. Reopen the editor to edit the latest version.")
		return
	}

	updated, err := m.service.Edit(m.ctx, memory.EditParams{
		ID: item.ID, Workspace: item.Workspace, Global: item.IsGlobal(), Content: content,
	})
	if err != nil {
		m.status = browseFailure("Could not save memory: %v", err)
		return
	}

	m.items[index] = updated
	m.filter()
	m.status = browseNotice("Memory saved.")
}

func (m *browseModel) removeEditFile() {
	if m.editFile != "" {
		_ = os.Remove(m.editFile)
		m.editFile = ""
	}
}

func (m *browseModel) deleteSelected() {
	index := m.matches[m.selected]
	item := m.items[index]
	outcomes, err := m.service.Delete(m.ctx, memory.DeleteParams{
		IDs: []memory.ID{item.ID}, Workspace: item.Workspace, Global: item.IsGlobal(),
	})
	if err != nil {
		m.status = browseFailure("Could not delete memory: %v", err)
		return
	}

	m.items = slices.Delete(m.items, index, index+1)
	m.filter()
	m.status = browseNotice("Memory deleted.")
	if !outcomes[0].Deleted {
		m.status = browseNotice("Memory was already deleted.")
	}
}
