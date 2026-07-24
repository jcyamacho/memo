package hookconfig

import (
	"encoding/json"
	"fmt"
)

type Config map[string]json.RawMessage

type SessionStartHook struct {
	Matcher       string
	Command       string
	StatusMessage string
	Timeout       int
}

func (c *Config) ConfigureSessionStart(spec SessionStartHook) error {
	if *c == nil {
		*c = Config{}
	}

	hooks := map[string]json.RawMessage{}
	if raw, ok := (*c)["hooks"]; ok {
		if err := json.Unmarshal(raw, &hooks); err != nil {
			return fmt.Errorf("decode hooks: %w", err)
		}
		if hooks == nil {
			hooks = map[string]json.RawMessage{}
		}
	}

	var entries []json.RawMessage
	if raw, ok := hooks["SessionStart"]; ok {
		if err := json.Unmarshal(raw, &entries); err != nil {
			return fmt.Errorf("decode SessionStart hooks: %w", err)
		}
	}

	configured, err := configureEntries(entries, spec)
	if err != nil {
		return err
	}
	hooks["SessionStart"], err = json.Marshal(configured)
	if err != nil {
		return err
	}
	(*c)["hooks"], err = json.Marshal(hooks)
	return err
}

func configureEntries(entries []json.RawMessage, spec SessionStartHook) ([]json.RawMessage, error) {
	configured := make([]json.RawMessage, 0, len(entries)+1)
	for _, entry := range entries {
		updated, keep, err := removeCommand(entry, spec.Command)
		if err != nil {
			return nil, err
		}
		if keep {
			configured = append(configured, updated)
		}
	}

	entry, err := json.Marshal(struct {
		Matcher string        `json:"matcher,omitempty"`
		Hooks   []commandHook `json:"hooks"`
	}{
		Matcher: spec.Matcher,
		Hooks: []commandHook{{
			Type:          "command",
			Command:       spec.Command,
			StatusMessage: spec.StatusMessage,
			Timeout:       spec.Timeout,
		}},
	})
	if err != nil {
		return nil, err
	}

	return append(configured, entry), nil
}

func removeCommand(raw json.RawMessage, command string) (json.RawMessage, bool, error) {
	var entry map[string]json.RawMessage
	if err := json.Unmarshal(raw, &entry); err != nil {
		return nil, false, fmt.Errorf("decode SessionStart entry: %w", err)
	}

	rawHooks, ok := entry["hooks"]
	if !ok {
		return raw, true, nil
	}

	var hooks []json.RawMessage
	if err := json.Unmarshal(rawHooks, &hooks); err != nil {
		return nil, false, fmt.Errorf("decode command hooks: %w", err)
	}

	filtered := hooks[:0]
	removed := false
	for _, rawHook := range hooks {
		var hook commandHook
		if err := json.Unmarshal(rawHook, &hook); err != nil {
			return nil, false, fmt.Errorf("decode command hook: %w", err)
		}
		if hook.Type == "command" && hook.Command == command {
			removed = true
			continue
		}
		filtered = append(filtered, rawHook)
	}

	if !removed {
		return raw, true, nil
	}
	if len(filtered) == 0 {
		return nil, false, nil
	}

	encodedHooks, err := json.Marshal(filtered)
	if err != nil {
		return nil, false, err
	}
	entry["hooks"] = encodedHooks
	updated, err := json.Marshal(entry)
	return updated, true, err
}

type commandHook struct {
	Type          string `json:"type"`
	Command       string `json:"command"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Timeout       int    `json:"timeout,omitempty"`
}
