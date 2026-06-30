package cmd

import "testing"

func TestVersionOutput(t *testing.T) {
	t.Parallel()

	got, err := versionOutput("1.2.3")
	if err != nil {
		t.Fatalf("versionOutput: %v", err)
	}

	want := "memo 1.2.3"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
