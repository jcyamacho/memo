package cmd

import (
	"testing"

	"github.com/jcyamacho/memo/internal/pkg/testutil"
)

func TestPathCommandPrintsStoreDir(t *testing.T) {
	t.Parallel()

	output, err := runCommand(t, New(newTestService(t), "/store/memo"), nil, "path")
	if err != nil {
		t.Fatalf("path: %v", err)
	}
	testutil.AssertEqual(t, "/store/memo\n", output)
}
