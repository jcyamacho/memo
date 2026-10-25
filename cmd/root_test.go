package cmd

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"strings"
	"testing"

	"github.com/jcyamacho/memo/internal/fsstore"
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func TestCommandExecutesMemoryLifecycleWithInjectedIO(t *testing.T) {
	t.Parallel()

	service := newTestService(t)
	addedOutput, err := executeCommand(t, service, nil, "add", " durable fact ")
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	added := decodeMemoryXML(t, addedOutput)
	if added.Workspace != "/repo" || added.Content != "durable fact" {
		t.Fatalf("added memory = %#v", added)
	}

	editedOutput, err := executeCommand(t, service, strings.NewReader("corrected fact\n"), "edit", added.ID)
	if err != nil {
		t.Fatalf("edit: %v", err)
	}
	if edited := decodeMemoryXML(t, editedOutput); edited.Content != "corrected fact" {
		t.Fatalf("edited content = %q, want corrected fact", edited.Content)
	}

	listOutput, err := executeCommand(t, service, nil, "list", "--workspace", ".")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var listed memoriesXML
	if err := xml.Unmarshal([]byte(listOutput), &listed); err != nil {
		t.Fatalf("unmarshal list output: %v", err)
	}
	if listed.Workspace != "/repo" || len(listed.Items) != 1 || listed.Items[0].Content != "corrected fact" {
		t.Fatalf("listed memories = %#v", listed)
	}

	if _, err := executeCommand(t, service, nil, "delete", added.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := executeCommand(t, service, nil, "get", added.ID); !errors.Is(err, memory.ErrNotFound) {
		t.Fatalf("get deleted memory = %v, want ErrNotFound", err)
	}
}

func TestAddCommandUsesNameInID(t *testing.T) {
	t.Parallel()

	service := newTestService(t)
	addedOutput, err := executeCommand(t, service, nil, "add", "--name", "Prefer Bun", "durable fact")
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	added := decodeMemoryXML(t, addedOutput)
	if !strings.HasPrefix(added.ID, "prefer-bun-") {
		t.Fatalf("added id = %q, want prefer-bun- prefix", added.ID)
	}

	gotOutput, err := executeCommand(t, service, nil, "get", added.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got := decodeMemoryXML(t, gotOutput); got.ID != added.ID || got.Content != "durable fact" {
		t.Fatalf("got memory = %#v", got)
	}
}

func TestCommandsRequireGlobalFlagForGlobalMemories(t *testing.T) {
	t.Parallel()

	service := newTestService(t)
	addedOutput, err := executeCommand(t, service, nil, "add", "--global", "global fact")
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	added := decodeMemoryXML(t, addedOutput)

	if _, err := executeCommand(t, service, nil, "get", added.ID); !errors.Is(err, memory.ErrNotFound) {
		t.Fatalf("get without --global = %v, want ErrNotFound", err)
	}
	if _, err := executeCommand(t, service, nil, "edit", "--global", added.ID, "--content", "corrected global fact"); err != nil {
		t.Fatalf("edit --global: %v", err)
	}
	gotOutput, err := executeCommand(t, service, nil, "get", "--global", added.ID)
	if err != nil {
		t.Fatalf("get --global: %v", err)
	}
	if got := decodeMemoryXML(t, gotOutput); got.Content != "corrected global fact" {
		t.Fatalf("global content = %q, want corrected global fact", got.Content)
	}
	if _, err := executeCommand(t, service, nil, "delete", "--global", added.ID); err != nil {
		t.Fatalf("delete --global: %v", err)
	}
}

type fixedWorkspaceResolver string

func (r fixedWorkspaceResolver) Resolve(context.Context, string) (string, error) {
	return string(r), nil
}

func newTestService(t *testing.T) *memory.Service {
	t.Helper()

	store, err := fsstore.New(t.TempDir())
	if err != nil {
		t.Fatalf("fsstore.New: %v", err)
	}
	return memory.NewService(store, fixedWorkspaceResolver("/repo"))
}

func executeCommand(t *testing.T, service *memory.Service, stdin *strings.Reader, args ...string) (string, error) {
	t.Helper()

	return runCommand(t, New(service, ""), stdin, args...)
}

func runCommand(t *testing.T, command *cobra.Command, stdin *strings.Reader, args ...string) (string, error) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if stdin != nil {
		command.SetIn(stdin)
	}
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs(args)

	err := command.ExecuteContext(t.Context())
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
	return stdout.String(), err
}

func decodeMemoryXML(t *testing.T, output string) memoryXML {
	t.Helper()

	var decoded memoryXML
	if err := xml.Unmarshal([]byte(output), &decoded); err != nil {
		t.Fatalf("unmarshal memory output: %v", err)
	}
	return decoded
}
