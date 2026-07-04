package codex

import (
	"encoding/json"
	"maps"
)

const memoContextCommand = "memo context"

type config struct {
	Hooks hooks `json:"hooks"`
	Extra map[string]json.RawMessage
}

func (c *config) UnmarshalJSON(data []byte) error {
	type configJSON struct {
		Hooks hooks `json:"hooks"`
	}

	var known configJSON
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}

	var extra map[string]json.RawMessage
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}
	delete(extra, "hooks")

	c.Hooks = known.Hooks
	c.Extra = extra
	return nil
}

func (c config) MarshalJSON() ([]byte, error) {
	fields := map[string]json.RawMessage{}
	maps.Copy(fields, c.Extra)
	if len(c.Hooks.SessionStart) > 0 {
		data, err := json.Marshal(c.Hooks)
		if err != nil {
			return nil, err
		}
		fields["hooks"] = data
	}
	return json.Marshal(fields)
}

func (c *config) configureSessionStartHook() {
	merged := make([]sessionStartEntry, 0, len(c.Hooks.SessionStart)+1)
	for _, entry := range c.Hooks.SessionStart {
		if !entry.hasHooks {
			merged = append(merged, entry)
			continue
		}

		entry.Hooks = removeMemoContextHooks(entry.Hooks)
		if len(entry.Hooks) > 0 {
			merged = append(merged, entry)
		}
	}

	c.Hooks.SessionStart = append(merged, hookEntry())
}

type hooks struct {
	SessionStart []sessionStartEntry `json:"SessionStart"`
	Extra        map[string]json.RawMessage
}

func (h *hooks) UnmarshalJSON(data []byte) error {
	type hooksJSON struct {
		SessionStart []sessionStartEntry `json:"SessionStart"`
	}

	var known hooksJSON
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}

	var extra map[string]json.RawMessage
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}
	delete(extra, "SessionStart")

	h.SessionStart = known.SessionStart
	h.Extra = extra
	return nil
}

func (h hooks) MarshalJSON() ([]byte, error) {
	fields := map[string]json.RawMessage{}
	maps.Copy(fields, h.Extra)
	if len(h.SessionStart) > 0 {
		data, err := json.Marshal(h.SessionStart)
		if err != nil {
			return nil, err
		}
		fields["SessionStart"] = data
	}
	return json.Marshal(fields)
}

type sessionStartEntry struct {
	Matcher  string `json:"matcher"`
	Hooks    []hook `json:"hooks"`
	Extra    map[string]json.RawMessage
	hasHooks bool
}

func (e *sessionStartEntry) UnmarshalJSON(data []byte) error {
	type entryJSON struct {
		Matcher string `json:"matcher"`
		Hooks   []hook `json:"hooks"`
	}

	var known entryJSON
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}

	var extra map[string]json.RawMessage
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}
	_, hasHooks := extra["hooks"]
	delete(extra, "matcher")
	delete(extra, "hooks")

	e.Matcher = known.Matcher
	e.Hooks = known.Hooks
	e.Extra = extra
	e.hasHooks = hasHooks
	return nil
}

func (e sessionStartEntry) MarshalJSON() ([]byte, error) {
	fields := map[string]json.RawMessage{}
	maps.Copy(fields, e.Extra)
	if e.Matcher != "" {
		data, err := json.Marshal(e.Matcher)
		if err != nil {
			return nil, err
		}
		fields["matcher"] = data
	}
	if len(e.Hooks) > 0 {
		data, err := json.Marshal(e.Hooks)
		if err != nil {
			return nil, err
		}
		fields["hooks"] = data
	}
	return json.Marshal(fields)
}

type hook struct {
	Type          string `json:"type"`
	Command       string `json:"command"`
	StatusMessage string `json:"statusMessage"`
	Timeout       *int   `json:"timeout"`
	Extra         map[string]json.RawMessage
}

func (h *hook) UnmarshalJSON(data []byte) error {
	type hookJSON struct {
		Type          string `json:"type"`
		Command       string `json:"command"`
		StatusMessage string `json:"statusMessage"`
		Timeout       *int   `json:"timeout"`
	}

	var known hookJSON
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}

	var extra map[string]json.RawMessage
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}
	delete(extra, "type")
	delete(extra, "command")
	delete(extra, "statusMessage")
	delete(extra, "timeout")

	h.Type = known.Type
	h.Command = known.Command
	h.StatusMessage = known.StatusMessage
	h.Timeout = known.Timeout
	h.Extra = extra
	return nil
}

func (h hook) MarshalJSON() ([]byte, error) {
	fields := map[string]json.RawMessage{}
	maps.Copy(fields, h.Extra)
	if h.Type != "" {
		data, err := json.Marshal(h.Type)
		if err != nil {
			return nil, err
		}
		fields["type"] = data
	}
	if h.Command != "" {
		data, err := json.Marshal(h.Command)
		if err != nil {
			return nil, err
		}
		fields["command"] = data
	}
	if h.StatusMessage != "" {
		data, err := json.Marshal(h.StatusMessage)
		if err != nil {
			return nil, err
		}
		fields["statusMessage"] = data
	}
	if h.Timeout != nil {
		data, err := json.Marshal(*h.Timeout)
		if err != nil {
			return nil, err
		}
		fields["timeout"] = data
	}
	return json.Marshal(fields)
}

func hookEntry() sessionStartEntry {
	return sessionStartEntry{
		Matcher: "startup|resume|clear|compact",
		Hooks: []hook{
			{
				Type:          "command",
				Command:       memoContextCommand,
				StatusMessage: "Loading memories",
				Timeout:       new(5),
			},
		},
		hasHooks: true,
	}
}

func removeMemoContextHooks(hooks []hook) []hook {
	filtered := make([]hook, 0, len(hooks))
	for _, hook := range hooks {
		if hook.Type == "command" && hook.Command == memoContextCommand {
			continue
		}
		filtered = append(filtered, hook)
	}
	return filtered
}
