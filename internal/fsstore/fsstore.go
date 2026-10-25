package fsstore

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jcyamacho/memo/internal/memory"
)

const (
	directoryPerm = 0o755
	filePerm      = 0o644
	markdownExt   = ".md"
)

type FS struct {
	globalsDir    string
	workspacesDir string
}

func New(dir string) (*FS, error) {
	globalsDir := filepath.Join(dir, "globals")
	workspacesDir := filepath.Join(dir, "workspaces")
	if err := os.MkdirAll(globalsDir, directoryPerm); err != nil {
		return nil, fmt.Errorf("create globals directory: %w", err)
	}
	if err := os.MkdirAll(workspacesDir, directoryPerm); err != nil {
		return nil, fmt.Errorf("create workspaces directory: %w", err)
	}
	return &FS{
		globalsDir:    globalsDir,
		workspacesDir: workspacesDir,
	}, nil
}

func (s *FS) Insert(_ context.Context, m memory.Memory) error {
	path := s.memoryPath(m.ID(), m.Workspace())

	_, err := os.Stat(path)
	if err == nil {
		return memory.ErrExists
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("check memory file: %w", err)
	}

	return writeMemoryFile(path, m.Content(), m.UpdatedAt())
}

func (s *FS) Get(_ context.Context, id memory.ID, workspace string) (memory.Memory, error) {
	path := s.memoryPath(id, workspace)
	if err := requireMemoryFile(path); err != nil {
		return memory.Memory{}, err
	}

	return s.readRecord(path, workspace)
}

func (s *FS) List(_ context.Context, workspace string) ([]memory.Memory, error) {
	targets := []memoryTarget{
		{path: s.globalsDir},
		{path: s.workspaceDirectory(workspace), workspace: workspace},
	}
	var items []memory.Memory
	for _, target := range targets {
		records, err := s.readDirectoryRecords(target.path, target.workspace)
		if err != nil {
			return nil, err
		}
		items = append(items, records...)
	}

	return items, nil
}

func (s *FS) Update(_ context.Context, m memory.Memory) error {
	path := s.memoryPath(m.ID(), m.Workspace())
	if err := requireMemoryFile(path); err != nil {
		return err
	}

	return writeMemoryFile(path, m.Content(), m.UpdatedAt())
}

func (s *FS) Delete(_ context.Context, id memory.ID, workspace string) error {
	err := os.Remove(s.memoryPath(id, workspace))
	if errors.Is(err, fs.ErrNotExist) {
		return memory.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("remove memory file: %w", err)
	}

	return s.cleanupWorkspaceDirectory(workspace)
}

func (s *FS) ListWorkspaces(_ context.Context) ([]string, error) {
	entries, err := os.ReadDir(s.workspacesDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("list workspaces: %w", err)
	}

	var workspaces []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		workspace, err := url.PathUnescape(entry.Name())
		if err != nil {
			return nil, fmt.Errorf("decode workspace %q: %w", entry.Name(), err)
		}
		workspaces = append(workspaces, workspace)
	}
	sort.Strings(workspaces)
	return workspaces, nil
}

type memoryTarget struct {
	path      string
	workspace string
}

func (s *FS) readDirectoryRecords(dirPath, workspace string) ([]memory.Memory, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read directory %s: %w", dirPath, err)
	}

	var records []memory.Memory
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), markdownExt) {
			continue
		}
		record, err := s.readRecord(filepath.Join(dirPath, entry.Name()), workspace)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (s *FS) readRecord(path, workspace string) (memory.Memory, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return memory.Memory{}, fmt.Errorf("read memory file: %w", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return memory.Memory{}, fmt.Errorf("stat memory file: %w", err)
	}
	record, err := memory.Load(
		strings.TrimSuffix(filepath.Base(path), markdownExt),
		string(content),
		workspace,
		info.ModTime(),
	)
	if err != nil {
		return memory.Memory{}, err
	}
	return record, nil
}

func requireMemoryFile(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return memory.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("check memory file: %w", err)
	}
	return nil
}

func (s *FS) memoryPath(id memory.ID, workspace string) string {
	if workspace == "" {
		return filepath.Join(s.globalsDir, id.String()+markdownExt)
	}
	return filepath.Join(s.workspaceDirectory(workspace), id.String()+markdownExt)
}

func (s *FS) workspaceDirectory(workspace string) string {
	return filepath.Join(s.workspacesDir, url.PathEscape(workspace))
}

func (s *FS) cleanupWorkspaceDirectory(workspace string) error {
	if workspace == "" {
		return nil
	}

	dir := s.workspaceDirectory(workspace)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read workspace directory: %w", err)
	}
	if len(entries) == 0 {
		if err := os.Remove(dir); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove empty workspace directory: %w", err)
		}
	}

	return nil
}
