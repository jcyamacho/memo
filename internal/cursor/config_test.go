package cursor

import (
	"path/filepath"
	"testing"

	"github.com/jcyamacho/memo/internal/hookconfig"
	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestConfigureSessionStart(t *testing.T) {
	t.Parallel()

	config := testutil.DecodeJSON[hookconfig.Config](t, `{
		"theme": "dark",
		"hooks": {
			"afterFileEdit": [{"command": "echo keep"}],
			"sessionStart": [
				{"command": "echo keep", "timeout": 2},
				{"command": "memo context --format cursor-hook", "timeout": 1},
				{"command": "memo context --format cursor-hook"}
			]
		}
	}`)

	for range 2 {
		if err := configureSessionStart(&config); err != nil {
			t.Fatalf("configureSessionStart: %v", err)
		}
	}

	got := testutil.EncodeJSON[map[string]any](t, config)
	want := map[string]any{
		"theme":   "dark",
		"version": float64(1),
		"hooks": map[string]any{
			"afterFileEdit": []any{map[string]any{"command": "echo keep"}},
			"sessionStart": []any{
				map[string]any{"command": "echo keep", "timeout": float64(2)},
				map[string]any{
					"command": "memo context --format cursor-hook",
					"timeout": float64(5),
				},
			},
		},
	}
	testutil.AssertEqual(t, want, got)
}

func TestConfigureSessionStartRejectsInvalidHooks(t *testing.T) {
	t.Parallel()

	config := testutil.DecodeJSON[hookconfig.Config](t, `{"hooks": []}`)
	if err := configureSessionStart(&config); err == nil {
		t.Fatal("expected invalid hooks error")
	}
}

func TestConfigureHookAtPreservesUnrelatedConfig(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), "nested", "hooks.json")
	testutil.WriteFile(t, target, `{"theme":"dark"}`)

	if err := configureHookAt(target); err != nil {
		t.Fatalf("configureHookAt: %v", err)
	}

	got := testutil.ReadJSON[map[string]any](t, target)
	want := map[string]any{
		"theme":   "dark",
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
