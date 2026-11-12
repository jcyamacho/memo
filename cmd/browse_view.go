package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/jcyamacho/memo/internal/memory"
)

var (
	browseAccent         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Cyan)
	browseLabel          = lipgloss.NewStyle().Bold(true)
	browseID             = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.BrightBlue)
	browseWorkspaceLabel = browseLabel.Foreground(lipgloss.Green)
	browseGlobalLabel    = browseLabel.Foreground(lipgloss.Magenta)
	browseSelection      = lipgloss.NewStyle().Bold(true).Reverse(true)
	browseHighlight      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Yellow)
	browseDelete         = browseLabel.Foreground(lipgloss.Red)
	browseDivider        = lipgloss.NewStyle().Faint(true)
)

func (m *browseModel) tooSmall() bool {
	return m.width < 34 || m.height < 11
}

func (m *browseModel) splitView() bool {
	return m.width >= 76
}

func (m *browseModel) listWidth() int {
	return min(52, m.width*2/5)
}

func (m *browseModel) bodyHeight() int {
	if m.height < 14 {
		return m.height - 5
	}
	return m.height - 8
}

func (m *browseModel) resize() {
	if m.tooSmall() {
		return
	}

	bodyWidth := m.width - 2
	detailWidth := bodyWidth
	if m.splitView() {
		detailWidth -= m.listWidth() + 3
	}

	m.search.SetWidth(bodyWidth - ansi.StringWidth(m.search.Prompt))
	m.detail.SetWidth(detailWidth)
	offset := m.detail.YOffset()
	m.showDetail()
	if offset > 0 {
		m.detail.SetYOffset(offset)
	}
}

func (m *browseModel) showDetail() {
	if len(m.matches) == 0 {
		m.detail.SetContent("")
		return
	}

	item := m.selectedItem()
	content := browseText(item.Content)
	if !m.tooSmall() {
		m.detail.SetHeight(m.bodyHeight() - len(m.detailHeader()))
		content = ansi.Wordwrap(content, m.detail.Width(), "")
	}
	m.detail.SetContent(content)
	m.detail.GotoTop()

	if m.query != nil {
		// A space in a search result can become a line break in the reader.
		displayQuery := strings.ReplaceAll(m.query.String(), " ", "[ \n]")
		pattern := regexp.MustCompile(displayQuery)
		m.detail.SetHighlights(pattern.FindAllStringIndex(content, -1))
		m.detail.SetXOffset(0)
	}
}

func (m *browseModel) View() tea.View {
	content := ""
	if m.tooSmall() {
		if m.width > 0 {
			content = ansi.Truncate("Resize terminal. Ctrl+C quits.", m.width, "")
		}
	} else {
		content = m.render()
	}

	view := tea.NewView(content)
	view.AltScreen = true
	if m.search.Focused() && !m.tooSmall() {
		if cursor := m.search.Cursor(); cursor != nil {
			cursor.X++
			cursor.Y = 2
			view.Cursor = cursor
		}
	}
	return view
}

func (m *browseModel) render() string {
	width := m.width - 2
	path := strings.Join(strings.Fields(browseText(m.workspace)), " ")
	title := browseAccent.Render("memo browse")
	header := title + "\n" + ansi.Truncate(path, width, "...")

	if m.confirming {
		item := m.selectedItem()
		metadata := browseMetadata(item, width)
		firstLine, _, _ := strings.Cut(strings.TrimSpace(browseText(item.Content)), "\n")
		excerpt := ansi.Truncate(firstLine, width, "...")
		confirm := browseDelete.Render("Delete this memory?") + "\n\n" +
			strings.Join(metadata, "\n") + "\n\n" + excerpt + "\n\ny Delete   n / Esc Cancel"
		if m.height < 14 {
			return browseInset(confirm)
		}
		return browseInset(header + "\n\n" + confirm)
	}

	body := m.renderBody(width)
	help := "↑↓ Select  Enter Read  / Find\ne Edit  d Delete  q Quit"
	switch {
	case m.search.Focused():
		help = "↑↓ Select  Enter Read  Tab List\nEsc Clear search"
	case m.reading:
		help = "↑↓ Scroll  Tab / Esc Results\ne Edit  d Delete  / Find  q Quit"
	case len(m.matches) == 0:
		help = "/ Search  Esc Clear search\nq Quit"
	}

	status := ansi.Truncate(strings.Join(strings.Fields(browseText(m.status)), " "), width, "...")
	if m.height < 14 {
		if status != "" {
			firstLine, _, _ := strings.Cut(help, "\n")
			help = firstLine + "\n" + browseLabel.Render(status)
		}
		return browseInset(header + "\n" + m.search.View() + "\n" + body + "\n" + help)
	}
	return browseInset(header + "\n" + m.search.View() + "\n\n" + body + "\n\n" +
		help + "\n" + browseLabel.Render(status))
}

func (m *browseModel) renderBody(width int) string {
	height := m.bodyHeight()
	if len(m.matches) == 0 {
		empty := "No memories yet."
		if len(m.items) > 0 {
			empty = "No memories match this view."
		}
		return lipgloss.NewStyle().Height(height).Render(empty)
	}

	if m.splitView() {
		list := m.renderList(m.listWidth(), height)
		divider := " " + browseDivider.Render("│") + " "
		gutter := strings.TrimSuffix(strings.Repeat(divider+"\n", height), "\n")
		return lipgloss.JoinHorizontal(lipgloss.Top, list, gutter, m.renderDetail())
	}
	if m.reading {
		return m.renderDetail()
	}
	return m.renderList(width, height)
}

func (m *browseModel) renderList(width, height int) string {
	header := fmt.Sprintf("RESULTS  %d of %d", m.selected+1, len(m.matches))
	if !m.reading && !m.search.Focused() {
		header = browseAccent.Render(header)
	}
	window := m.listWindow()
	rows := []string{header}
	lastScope := ""
	for i, index := range m.matches[window.start:window.end] {
		item := m.items[index]
		scope := browseScope(item)
		if scope != lastScope {
			if lastScope != "" {
				rows = append(rows, "")
			}
			rows = append(rows, browseScopeStyle(item).Render(scope))
			lastScope = scope
		}

		id := browseListID(item.ID.String(), width-2)
		row := "  " + id
		if window.start+i == m.selected {
			row = browseSelection.Width(width).Render("> " + id)
		}
		rows = append(rows, row)
	}

	style := lipgloss.NewStyle().Width(width).Height(height)
	return style.Render(strings.Join(rows, "\n"))
}

type browseListWindow struct{ start, end int }

func (w browseListWindow) size() int { return w.end - w.start }

func (m *browseModel) listWindow() browseListWindow {
	if len(m.matches) == 0 {
		return browseListWindow{}
	}

	capacity := m.bodyHeight() - 1
	start := max(0, m.selected-capacity+2)
	for m.listRows(start, m.selected+1) > capacity {
		start++
	}
	end := m.selected + 1
	for end < len(m.matches) && m.listRows(start, end+1) <= capacity {
		end++
	}
	for start > 0 && m.listRows(start-1, end) <= capacity {
		start--
	}
	return browseListWindow{start: start, end: end}
}

func (m *browseModel) listRows(start, end int) int {
	rows := end - start + 1
	if m.items[m.matches[start]].IsGlobal() != m.items[m.matches[end-1]].IsGlobal() {
		rows += 2
	}
	return rows
}

func (m *browseModel) renderDetail() string {
	header := m.detailHeader()
	contentLabel := len(header) - 1
	if m.reading {
		header[contentLabel] = browseAccent.Render(header[contentLabel])
	} else {
		header[contentLabel] = browseLabel.Render(header[contentLabel])
	}
	return strings.Join(header, "\n") + "\n" + m.detail.View()
}

func (m *browseModel) detailHeader() []string {
	header := browseMetadata(m.selectedItem(), m.detail.Width())
	if m.bodyHeight() >= 9 {
		header = append(header, "")
	}
	return append(header, "CONTENT")
}

func browseMetadata(item memory.MemoryDTO, width int) []string {
	scope := browseScopeStyle(item).Render(strings.ToUpper(browseScope(item)) + " MEMORY")
	date := item.UpdatedAt.Local().Format("2 Jan 2006")
	wrappedID := ansi.Hardwrap(item.ID.String(), width-4, true)
	lines := []string{scope + "  " + date}
	label := "ID  "
	for line := range strings.SplitSeq(wrappedID, "\n") {
		lines = append(lines, label+browseID.Render(line))
		label = "    "
	}
	return lines
}

func browseListID(id string, width int) string {
	if len(id) <= width {
		return id
	}

	const tailWidth = 7
	const ellipsis = "..."
	return id[:width-tailWidth-len(ellipsis)] + ellipsis + id[len(id)-tailWidth:]
}

func browseScope(item memory.MemoryDTO) string {
	if item.IsGlobal() {
		return "Global"
	}
	return "Workspace"
}

func browseScopeStyle(item memory.MemoryDTO) lipgloss.Style {
	if item.IsGlobal() {
		return browseGlobalLabel
	}
	return browseWorkspaceLabel
}

func browseInset(content string) string {
	return " " + strings.ReplaceAll(content, "\n", "\n ")
}

func browseText(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\t", "    ")
	return strings.Map(func(r rune) rune {
		if r != '\n' && unicode.IsControl(r) {
			return -1
		}
		return r
	}, ansi.Strip(content))
}
