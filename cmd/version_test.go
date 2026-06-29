package cmd

import "testing"

func TestVersionOutput(t *testing.T) {
	t.Parallel()

	got, err := versionOutput("1.2.3", "abc123", "2026-06-29T10:00:00Z")
	if err != nil {
		t.Fatalf("versionOutput: %v", err)
	}
	want := "memo 1.2.3 (commit abc123, built 2026-06-29T10:00:00Z)"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
