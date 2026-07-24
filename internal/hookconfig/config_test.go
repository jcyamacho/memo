package hookconfig

import (
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestConfigureSessionStart(t *testing.T) {
	t.Parallel()

	config := testutil.DecodeJSON[Config](t, `{
		"theme": "dark",
		"hooks": {
			"Other": [{"name": "keep"}],
			"SessionStart": [
				{"matcher": "other", "hooks": [
					{"type": "command", "command": "echo keep", "timeout": 0},
					{"type": "command", "command": "memo context"}
				]},
				{"matcher": "external", "hooks": [], "label": "keep"},
				{"hooks": [{"type": "command", "command": "memo context"}]},
				{"label": "without hooks"}
			]
		}
	}`)
	spec := SessionStartHook{
		Matcher:       "startup|resume|clear|compact",
		Command:       "memo context",
		StatusMessage: "Loading memories",
		Timeout:       5,
	}

	for range 2 {
		if err := config.ConfigureSessionStart(spec); err != nil {
			t.Fatalf("ConfigureSessionStart: %v", err)
		}
	}

	got := testutil.EncodeJSON[map[string]any](t, config)
	want := map[string]any{
		"theme": "dark",
		"hooks": map[string]any{
			"Other": []any{map[string]any{"name": "keep"}},
			"SessionStart": []any{
				map[string]any{
					"matcher": "other",
					"hooks": []any{
						map[string]any{"type": "command", "command": "echo keep", "timeout": float64(0)},
					},
				},
				map[string]any{"matcher": "external", "hooks": []any{}, "label": "keep"},
				map[string]any{"label": "without hooks"},
				map[string]any{
					"matcher": "startup|resume|clear|compact",
					"hooks": []any{
						map[string]any{
							"type":          "command",
							"command":       "memo context",
							"statusMessage": "Loading memories",
							"timeout":       float64(5),
						},
					},
				},
			},
		},
	}
	testutil.AssertEqual(t, want, got)
}

func TestConfigureSessionStartRejectsInvalidHooks(t *testing.T) {
	t.Parallel()

	config := testutil.DecodeJSON[Config](t, `{"hooks": []}`)
	if err := config.ConfigureSessionStart(SessionStartHook{Command: "memo context"}); err == nil {
		t.Fatal("expected invalid hooks error")
	}
}
