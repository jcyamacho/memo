package memory

import (
	"context"
	"testing"
	"time"
)

type mockStore struct {
	t *testing.T

	inserted  []Memory
	insertErr error

	getByID    map[ID]Memory
	getErrByID map[ID]error

	listItems     []Memory
	listErr       error
	listWorkspace string

	updated   []Memory
	updateErr error

	deleteErrByID map[ID]error
	deletedIDs    []string

	workspaces        []string
	listWorkspacesErr error
}

func newMockStore(t *testing.T) *mockStore {
	t.Helper()

	return &mockStore{
		t:             t,
		getByID:       make(map[ID]Memory),
		getErrByID:    make(map[ID]error),
		deleteErrByID: make(map[ID]error),
	}
}

func (m *mockStore) Insert(_ context.Context, memory Memory) error {
	m.inserted = append(m.inserted, memory)
	return m.insertErr
}

func (m *mockStore) Get(_ context.Context, id ID) (Memory, error) {
	m.t.Helper()

	if err, ok := m.getErrByID[id]; ok {
		return Memory{}, err
	}
	if memory, ok := m.getByID[id]; ok {
		return memory, nil
	}
	m.t.Fatalf("Get called with unstubbed id %q", id)
	return Memory{}, nil
}

func (m *mockStore) List(_ context.Context, workspace string) ([]Memory, error) {
	m.listWorkspace = workspace
	return m.listItems, m.listErr
}

func (m *mockStore) Update(_ context.Context, memory Memory) error {
	m.updated = append(m.updated, memory)
	return m.updateErr
}

func (m *mockStore) Delete(_ context.Context, id ID) error {
	m.deletedIDs = append(m.deletedIDs, id.String())
	return m.deleteErrByID[id]
}

func (m *mockStore) ListWorkspaces(context.Context) ([]string, error) {
	return m.workspaces, m.listWorkspacesErr
}

type mockWorkspaceResolver struct {
	resolved string
	err      error
	paths    []string
}

func newMockWorkspaceResolver(resolved string) *mockWorkspaceResolver {
	return &mockWorkspaceResolver{resolved: resolved}
}

func (m *mockWorkspaceResolver) Resolve(_ context.Context, path string) (string, error) {
	m.paths = append(m.paths, path)
	return m.resolved, m.err
}

func mustLoadMemory(t *testing.T, id, content, workspace string, updatedAt time.Time) Memory {
	t.Helper()

	m, err := Load(id, content, workspace, updatedAt)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return m
}

func mustParseID(t *testing.T, value string) ID {
	t.Helper()

	id, err := ParseID(value)
	if err != nil {
		t.Fatalf("ParseID: %v", err)
	}
	return id
}
