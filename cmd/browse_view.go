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
	browseDanger         = browseLabel.Foreground(lipgloss.Red)
	browseMuted          = lipgloss.NewStyle().Faint(true)
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
		return m.height - 4
	}
	return m.height - 7
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
			cursor.Y = 1
			view.Cursor = cursor
		}
	}
	return view
}

func (m *browseModel) render() string {
	width := m.width - 2
	path := strings.Join(strings.Fields(browseText(m.workspace)), " ")
	const title = "memo browse"
	header := browseAccent.Render(title) + "  " + browsePathTail(path, width-len(title)-2)

	if m.confirming {
		item := m.selectedItem()
		metadata := browseMetadata(item, width)
		firstLine, _, _ := strings.Cut(strings.TrimSpace(browseText(item.Content)), "\n")
		excerpt := ansi.Truncate(firstLine, width, "...")
		confirm := browseDanger.Render("Delete this memory?") + "\n\n" +
			strings.Join(metadata, "\n") + "\n\n" + excerpt + "\n\ny Delete   n / Esc Cancel"
		if m.height < 14 {
			return browseInset(confirm)
		}
		return browseInset(header + "\n\n" + confirm)
	}

	search := m.search.View()
	if !m.search.Focused() && m.search.Value() == "" {
		search = browseMuted.Render("Press / to search")
	}

	body := m.renderBody(width)
	help := m.help()

	statusStyle := browseLabel
	if m.status.failed {
		statusStyle = browseDanger
	}
	statusText := ansi.Truncate(strings.Join(strings.Fields(browseText(m.status.text)), " "), width, "...")
	status := statusStyle.Render(statusText)

	if m.height < 14 {
		if statusText != "" {
			help[1] = status
		}
		return browseInset(header + "\n" + search + "\n" + body + "\n" + strings.Join(help, "\n"))
	}
	return browseInset(header + "\n" + search + "\n\n" + body + "\n\n" +
		strings.Join(help, "\n") + "\n" + status)
}

type browseHint struct{ keys, action string }

func (m *browseModel) help() []string {
	manage := []browseHint{{"e", "Edit"}, {"d", "Delete"}, {"q", "Quit"}}
	switch {
	case m.search.Focused() && m.splitView():
		return browseHelp([]browseHint{{"↑↓", "Select"}, {"Enter", "List"}}, []browseHint{{"Esc", "Clear search"}})
	case m.search.Focused():
		return browseHelp(
			[]browseHint{{"↑↓", "Select"}, {"Enter", "Read"}, {"Tab", "List"}},
			[]browseHint{{"Esc", "Clear search"}},
		)
	case len(m.matches) == 0:
		return browseHelp([]browseHint{{"/", "Search"}, {"Esc", "Clear search"}}, []browseHint{{"q", "Quit"}})
	case m.splitView():
		return browseHelp([]browseHint{{"↑↓", "Select"}, {"PgUp/PgDn", "Scroll"}, {"/", "Search"}}, manage)
	case m.reading:
		return browseHelp([]browseHint{{"↑↓", "Scroll"}, {"Esc", "List"}, {"/", "Search"}}, manage)
	}
	return browseHelp([]browseHint{{"↑↓", "Select"}, {"Enter", "Read"}, {"/", "Search"}}, manage)
}

func browseHelp(lines ...[]browseHint) []string {
	rendered := make([]string, 0, len(lines))
	for _, hints := range lines {
		parts := make([]string, 0, len(hints))
		for _, hint := range hints {
			parts = append(parts, browseLabel.Render(hint.keys)+" "+browseMuted.Render(hint.action))
		}
		rendered = append(rendered, strings.Join(parts, "  "))
	}
	return rendered
}

func browsePathTail(path string, width int) string {
	const ellipsis = "..."
	overflow := ansi.StringWidth(path) - width
	if overflow <= 0 {
		return path
	}
	return ellipsis + ansi.TruncateLeft(path, overflow+len(ellipsis), "")
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
		divider := " " + browseMuted.Render("│") + " "
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
	if !m.readerOpen() && !m.search.Focused() {
		header = browseAccent.Render(header)
	}

	globalCount := 0
	for _, index := range m.matches {
		if m.items[index].IsGlobal() {
			globalCount++
		}
	}
	scopeCounts := map[string]int{"Global": globalCount, "Workspace": len(m.matches) - globalCount}

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
			rows = append(rows, browseScopeStyle(item).Render(fmt.Sprintf("%s (%d)", scope, scopeCounts[scope])))
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
	labelStyle := browseLabel
	if m.readerOpen() {
		labelStyle = browseAccent
	}
	header[contentLabel] = labelStyle.Render(header[contentLabel]) + browseMuted.Render(m.scrollHint())

	return strings.Join(header, "\n") + "\n" + m.detail.View()
}

// scrollHint is empty when the whole memory fits in the detail pane.
func (m *browseModel) scrollHint() string {
	arrows := ""
	if !m.detail.AtTop() {
		arrows += "↑"
	}
	if !m.detail.AtBottom() {
		arrows += "↓"
	}
	if arrows == "" {
		return ""
	}
	return fmt.Sprintf("  %s %d%%", arrows, int(m.detail.ScrollPercent()*100))
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
