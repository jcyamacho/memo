package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/jcyamacho/memo/internal/fsstore"
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestBrowseFiltersFullContentAndConfirmsScopedDelete(t *testing.T) {
	t.Parallel()

	model := newTestBrowseModel(t)
	model.handleKey(browseTestKey('/'))
	model.Update(tea.PasteMsg{Content: "SECOND LINE"})
	testutil.AssertEqual(t, len(model.matches), 1)
	testutil.AssertEqual(t, model.selectedItem().IsGlobal(), true)
	view := ansi.Strip(model.View().Content)
	if !strings.Contains(view, "> "+testMemoryIDOne) || !strings.Contains(view, "Second line") {
		t.Fatal("content search should show the result ID and its matching text in the detail pane")
	}
	model.handleKey(browseTestKey(tea.KeyEnter))
	testutil.AssertEqual(t, model.reading, true)

	model.handleKey(browseTestKey('d'))
	view = model.View().Content
	if !strings.Contains(view, "Delete this memory?") || !strings.Contains(view, "Global") {
		t.Fatalf("confirmation does not identify the scope: %q", view)
	}
	model.handleKey(browseTestKey(tea.KeyEsc))
	params := memory.GetParams{ID: model.selectedItem().ID, Global: true}
	if _, err := model.service.Get(t.Context(), params); err != nil {
		t.Fatalf("get after canceling deletion: %v", err)
	}

	model.handleKey(browseTestKey('d'))
	model.handleKey(browseTestKey('y'))
	if _, err := model.service.Get(t.Context(), params); !errors.Is(err, memory.ErrNotFound) {
		t.Fatalf("get after confirming deletion = %v, want ErrNotFound", err)
	}
	workspace, err := model.service.Get(t.Context(), memory.GetParams{ID: params.ID, Workspace: "/repo"})
	if err != nil {
		t.Fatalf("get workspace memory with the same ID: %v", err)
	}
	testutil.AssertEqual(t, workspace.Content, "Workspace text")
}

func TestBrowseKeepsScopeAndContentSeparateOnNarrowScreens(t *testing.T) {
	t.Parallel()

	model := newTestBrowseModel(t)
	longID := strings.Repeat("a", 57) + "-abc123"
	content := "First paragraph.\n\n" + strings.Repeat("More text.\n", 8) + "Last paragraph."
	model.items[1] = mustLoadMemoryDTO(t, longID, content, "", time.Now())
	model.filter()
	model.Update(tea.WindowSizeMsg{Width: 34, Height: 11})

	list := ansi.Strip(model.View().Content)
	if !strings.Contains(list, "\n Workspace ") ||
		!strings.Contains(list, "\n > 550e8400-e29b-41d4-a...5440000") ||
		!strings.Contains(list, "\n Global ") ||
		!strings.Contains(list, "\n   "+strings.Repeat("a", 20)+"...-abc123") ||
		strings.Contains(list, "First paragraph.") ||
		strings.Contains(list, longID) {
		t.Fatalf("narrow list does not group shortened IDs by scope:\n%s", list)
	}

	model.handleKey(browseTestKey(tea.KeyDown))
	model.handleKey(browseTestKey(tea.KeyEnter))
	model.detail.GotoBottom()
	detail := ansi.Strip(model.View().Content)
	testutil.AssertEqual(t, strings.Count(detail, "\n")+1, 11)
	for _, expected := range []string{"GLOBAL MEMORY", "ID  " + strings.Repeat("a", 28), "CONTENT\n Last paragraph."} {
		if !strings.Contains(detail, expected) {
			t.Fatalf("scrolled detail is missing %q:\n%s", expected, detail)
		}
	}
	for line := range strings.SplitSeq(detail, "\n") {
		if ansi.StringWidth(line) > 34 {
			t.Fatalf("detail line exceeds the terminal width: %q", line)
		}
	}
}

func TestBrowseSavesEditToSelectedScope(t *testing.T) {
	t.Parallel()

	model := newTestBrowseModel(t)
	model.handleKey(browseTestKey(tea.KeyDown))
	file := filepath.Join(t.TempDir(), "memory.md")
	if err := os.WriteFile(file, []byte("Edited global text\n"), 0o600); err != nil {
		t.Fatalf("write edited memory: %v", err)
	}
	model.editFile = file
	model.handleKey(browseTestKey(tea.KeyUp))
	model.Update(browseEditorDone{})

	params := memory.GetParams{ID: model.selectedItem().ID, Global: true}
	updated, err := model.service.Get(t.Context(), params)
	if err != nil {
		t.Fatalf("get edited global memory: %v", err)
	}
	testutil.AssertEqual(t, updated.Content, "Edited global text")

	workspace, err := model.service.Get(t.Context(), memory.GetParams{ID: params.ID, Workspace: "/repo"})
	if err != nil {
		t.Fatalf("get workspace memory with the same ID: %v", err)
	}
	testutil.AssertEqual(t, workspace.Content, "Workspace text")

	if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat editor file = %v, want ErrNotExist", err)
	}
}

func newTestBrowseModel(t *testing.T) *browseModel {
	t.Helper()

	store, err := fsstore.New(t.TempDir())
	if err != nil {
		t.Fatalf("fsstore.New: %v", err)
	}
	for scope, content := range map[string]string{"/repo": "Workspace text", "": "Global text\nSecond line"} {
		item, err := memory.Load(testMemoryIDOne, content, scope, time.Now())
		if err != nil {
			t.Fatalf("load test memory: %v", err)
		}
		if err := store.Insert(t.Context(), item); err != nil {
			t.Fatalf("insert test memory: %v", err)
		}
	}
	service := memory.NewService(store, fixedWorkspaceResolver("/repo"))
	result, err := service.List(t.Context(), "/repo")
	if err != nil {
		t.Fatalf("list test memories: %v", err)
	}

	model := newBrowseModel(t.Context(), service, result)
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 24})
	return model
}

func browseTestKey(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code, Text: string(code)}
}
