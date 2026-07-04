package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertEqual[T any](t *testing.T, want T, got T) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}
}
