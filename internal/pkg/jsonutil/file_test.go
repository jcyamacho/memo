package jsonutil

import (
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	Enabled bool   `json:"enabled,omitempty"`
	Name    string `json:"name,omitempty"`
}

func TestReadMissingFileReturnsZeroValue(t *testing.T) {
	t.Parallel()

	got, err := Read[testConfig](filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got != (testConfig{}) {
		t.Fatalf("got %#v, want zero value", got)
	}
}

func TestWriteAndRead(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), "nested", "config.json")
	want := testConfig{Enabled: true, Name: "memo"}
	if err := Write(target, want); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := Read[testConfig](target)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got != want {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestReadRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(target, []byte("{invalid"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := Read[testConfig](target); err == nil {
		t.Fatal("expected error")
	}
}
