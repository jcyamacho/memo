package cursor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestConfigureHookConfiguresDefaultCursorHook(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	target, err := ConfigureHook()
	if err != nil {
		t.Fatalf("ConfigureHook: %v", err)
	}

	wantTarget := filepath.Join(home, ".cursor", "hooks.json")
	testutil.AssertEqual(t, wantTarget, target)

	got := testutil.ReadJSON[map[string]any](t, target)
	want := map[string]any{
		"version": float64(1),
		"hooks": map[string]any{
			"sessionStart": []any{
				map[string]any{
					"command": "memo context --format cursor-hook",
					"timeout": float64(5),
				},
			},
		},
	}
	testutil.AssertEqual(t, want, got)
}

func TestConfigureHookLeavesInvalidCursorJSONUnchanged(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), ".cursor", "hooks.json")
	testutil.WriteFile(t, target, "{invalid")

	if err := configureHookAt(target); err == nil {
		t.Fatal("expected invalid JSON error")
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	testutil.AssertEqual(t, "{invalid", string(data))
}
