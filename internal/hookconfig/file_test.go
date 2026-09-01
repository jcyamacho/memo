package hookconfig

import (
	"path/filepath"
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestConfigureFilePreservesUnrelatedConfig(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), "nested", "hooks.json")
	testutil.WriteFile(t, target, `{"theme":"dark"}`)

	err := ConfigureFile(target, SessionStartHook{
		Command:       "memo context",
		StatusMessage: "Loading memories",
		Timeout:       5,
	})
	if err != nil {
		t.Fatalf("ConfigureFile: %v", err)
	}

	config := testutil.ReadJSON[map[string]any](t, target)
	if config["theme"] != "dark" {
		t.Fatalf("theme = %v, want dark", config["theme"])
	}
	if _, ok := config["hooks"]; !ok {
		t.Fatal("hooks not configured")
	}
}
