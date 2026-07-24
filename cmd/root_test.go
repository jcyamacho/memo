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

	command := New(service)
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
