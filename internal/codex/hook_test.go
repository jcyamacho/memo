package codex

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestConfigureHookConfiguresDefaultCodexHook(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target, err := ConfigureHook()
	if err != nil {
		t.Fatalf("ConfigureHook: %v", err)
	}

	wantTarget := filepath.Join(home, ".codex", "hooks.json")
	if target != wantTarget {
		t.Fatalf("target = %q, want %q", target, wantTarget)
	}

	sessionStart := readSessionStart(t, target)
	if len(sessionStart) != 1 {
		t.Fatalf("SessionStart length = %d, want 1: %#v", len(sessionStart), sessionStart)
	}

	entry := sessionStart[0].(map[string]any)
	if entry["matcher"] != "startup|resume|clear|compact" {
		t.Fatalf("matcher = %v, want startup|resume|clear|compact", entry["matcher"])
	}

	hooks := entry["hooks"].([]any)
	if len(hooks) != 1 {
		t.Fatalf("hooks length = %d, want 1: %#v", len(hooks), hooks)
	}
	if hooks[0].(map[string]any)["command"] != "memo context" {
		t.Fatalf("unexpected hook: %#v", hooks[0])
	}
}

func TestConfigureHookLeavesInvalidCodexJSONUnchanged(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), ".codex", "hooks.json")
	testutil.WriteFile(t, target, "{invalid")

	if err := configureHookAt(target); err == nil {
		t.Fatal("expected invalid JSON error")
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "{invalid" {
		t.Fatalf("file changed after invalid JSON: %q", data)
	}
}

func readSessionStart(t *testing.T, path string) []any {
	t.Helper()

	config := testutil.ReadJSON[map[string]any](t, path)
	return config["hooks"].(map[string]any)["SessionStart"].([]any)
}
