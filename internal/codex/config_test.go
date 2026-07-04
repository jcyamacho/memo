package codex

import (
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestConfigureSessionStartHookWritesCodexHook(t *testing.T) {
	t.Parallel()

	config := config{}
	config.configureSessionStartHook()

	if len(config.Hooks.SessionStart) != 1 {
		t.Fatalf("SessionStart length = %d, want 1: %#v", len(config.Hooks.SessionStart), config.Hooks.SessionStart)
	}

	entry := config.Hooks.SessionStart[0]
	if entry.Matcher != "startup|resume|clear|compact" {
		t.Fatalf("matcher = %v, want startup|resume|clear|compact", entry.Matcher)
	}

	hooks := entry.Hooks
	if len(hooks) != 1 {
		t.Fatalf("hooks length = %d, want 1: %#v", len(hooks), hooks)
	}

	hook := hooks[0]
	if hook.Type != "command" || hook.Command != "memo context" {
		t.Fatalf("unexpected hook command: %#v", hook)
	}
	if hook.StatusMessage != "Loading memories" {
		t.Fatalf("statusMessage = %v, want Loading memories", hook.StatusMessage)
	}
	if hook.Timeout == nil || *hook.Timeout != 5 {
		t.Fatalf("timeout = %v, want 5", hook.Timeout)
	}
}

func TestConfigureSessionStartHookPreservesUnrelatedCodexSettings(t *testing.T) {
	t.Parallel()

	config := testutil.DecodeJSON[config](t, `{
		"theme": "dark",
		"hooks": {
			"Other": [{"name": "keep"}],
			"SessionStart": [
				{"matcher": "other", "hooks": [{"type": "command", "command": "echo keep"}]},
				{"label": "without hooks"}
			]
		}
	}`)
	config.configureSessionStartHook()
	got := testutil.EncodeJSON[map[string]any](t, config)
	want := map[string]any{
		"theme": "dark",
		"hooks": map[string]any{
			"Other": []any{
				map[string]any{"name": "keep"},
			},
			"SessionStart": []any{
				map[string]any{
					"matcher": "other",
					"hooks": []any{
						map[string]any{"type": "command", "command": "echo keep"},
					},
				},
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

func TestConfigureSessionStartHookPreservesExplicitZeroTimeout(t *testing.T) {
	t.Parallel()

	config := testutil.DecodeJSON[config](t, `{
		"hooks": {
			"SessionStart": [
				{"matcher": "", "hooks": [
					{"type": "command", "command": "echo keep", "statusMessage": "", "timeout": 0}
				]}
			]
		}
	}`)
	config.configureSessionStartHook()
	got := testutil.EncodeJSON[map[string]any](t, config)
	want := map[string]any{
		"hooks": map[string]any{
			"SessionStart": []any{
				map[string]any{
					"hooks": []any{
						map[string]any{"type": "command", "command": "echo keep", "timeout": float64(0)},
					},
				},
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

func TestConfigureSessionStartHookRemovesDuplicateCodexMemoHooks(t *testing.T) {
	t.Parallel()

	config := testutil.DecodeJSON[config](t, `{
		"hooks": {
			"SessionStart": [
				{"hooks": [
					{"type": "command", "command": "memo context"},
					{"type": "command", "command": "echo keep"}
				]},
				{"hooks": [{"type": "command", "command": "memo context"}]},
				{"matcher": "startup|resume|clear|compact", "hooks": [
					{"type": "command", "command": "memo context", "statusMessage": "old"}
				]}
			]
		}
	}`)
	config.configureSessionStartHook()

	sessionStart := config.Hooks.SessionStart
	if len(sessionStart) != 2 {
		t.Fatalf("SessionStart length = %d, want 2: %#v", len(sessionStart), sessionStart)
	}

	remaining := sessionStart[0].Hooks
	if len(remaining) != 1 {
		t.Fatalf("remaining hooks length = %d, want 1: %#v", len(remaining), remaining)
	}
	if remaining[0].Command != "echo keep" {
		t.Fatalf("unexpected remaining hook: %#v", remaining[0])
	}
}

func TestConfigureSessionStartHookIsIdempotent(t *testing.T) {
	t.Parallel()

	config := config{}
	for range 2 {
		config.configureSessionStartHook()
	}

	sessionStart := config.Hooks.SessionStart
	if len(sessionStart) != 1 {
		t.Fatalf("SessionStart length = %d, want 1: %#v", len(sessionStart), sessionStart)
	}
}
