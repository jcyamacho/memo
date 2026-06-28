package cmd

import (
	"strings"
	"testing"
)

func TestReadContentUsesPositionalArgument(t *testing.T) {
	t.Parallel()

	got, err := readContentFrom([]string{"content"}, strings.NewReader("stdin"))
	if err != nil {
		t.Fatalf("readContentFrom: %v", err)
	}
	if got != "content" {
		t.Fatalf("got %q, want positional content", got)
	}
}

func TestReadContentUsesStdin(t *testing.T) {
	t.Parallel()

	got, err := readContentFrom(nil, strings.NewReader("stdin content\n"))
	if err != nil {
		t.Fatalf("readContentFrom: %v", err)
	}
	if got != "stdin content\n" {
		t.Fatalf("got %q, want stdin content", got)
	}
}

func TestReadContentRequiresContentWithoutStdin(t *testing.T) {
	t.Parallel()

	if _, err := readContentFrom(nil, nil); err == nil {
		t.Fatal("expected error")
	}
}
