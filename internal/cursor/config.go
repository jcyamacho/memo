package cursor

import (
	"encoding/json"
	"fmt"

	"github.com/jcyamacho/memo/internal/hookconfig"
)

func configureSessionStart(config *hookconfig.Config) error {
	if *config == nil {
		*config = hookconfig.Config{}
	}

	version, err := json.Marshal(1)
	if err != nil {
		return err
	}
	(*config)["version"] = version

	hooks := map[string]json.RawMessage{}
	if raw, ok := (*config)["hooks"]; ok {
		if err := json.Unmarshal(raw, &hooks); err != nil {
			return fmt.Errorf("decode hooks: %w", err)
		}
		if hooks == nil {
			hooks = map[string]json.RawMessage{}
		}
	}

	var entries []json.RawMessage
	if raw, ok := hooks["sessionStart"]; ok {
		if err := json.Unmarshal(raw, &entries); err != nil {
			return fmt.Errorf("decode sessionStart hooks: %w", err)
		}
	}

	configured, err := configureEntries(entries)
	if err != nil {
		return err
	}
	hooks["sessionStart"], err = json.Marshal(configured)
	if err != nil {
		return err
	}
	(*config)["hooks"], err = json.Marshal(hooks)
	return err
}

func configureEntries(entries []json.RawMessage) ([]json.RawMessage, error) {
	configured := make([]json.RawMessage, 0, len(entries)+1)
	for _, entry := range entries {
		var hook commandHook
		if err := json.Unmarshal(entry, &hook); err != nil {
			return nil, fmt.Errorf("decode sessionStart entry: %w", err)
		}
		if hook.Command == memoContextCommand {
			continue
		}
		configured = append(configured, entry)
	}

	entry, err := json.Marshal(commandHook{
		Command: memoContextCommand,
		Timeout: 5,
	})
	if err != nil {
		return nil, err
	}

	return append(configured, entry), nil
}

type commandHook struct {
	Command string `json:"command"`
	Timeout int    `json:"timeout,omitempty"`
}
