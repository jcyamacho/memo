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

	getByID    map[string]Memory
	getErrByID map[string]error

	listItems          []Memory
	listErr            error
	listWorkspace      string
	listIncludeGlobals bool

	updated   []Memory
	updateErr error

	deleteErrByID map[string]error
	deletedIDs    []string

	workspaces        []string
	listWorkspacesErr error
}

func newMockStore(t *testing.T) *mockStore {
	t.Helper()

	return &mockStore{
		t:             t,
		getByID:       make(map[string]Memory),
		getErrByID:    make(map[string]error),
		deleteErrByID: make(map[string]error),
	}
}

func (m *mockStore) Insert(_ context.Context, memory Memory) error {
	m.inserted = append(m.inserted, memory)
	return m.insertErr
}

func (m *mockStore) Get(_ context.Context, id string) (Memory, error) {
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

func (m *mockStore) List(_ context.Context, workspace string, options ...ListOption) ([]Memory, error) {
	m.listWorkspace = workspace
	m.listIncludeGlobals = NewListOptions(options...).IncludeGlobals
	return m.listItems, m.listErr
}

func (m *mockStore) Update(_ context.Context, memory Memory) error {
	m.updated = append(m.updated, memory)
	return m.updateErr
}

func (m *mockStore) Delete(_ context.Context, id string) error {
	m.deletedIDs = append(m.deletedIDs, id)
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
