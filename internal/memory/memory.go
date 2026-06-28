package memory

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Memory struct {
	id        string
	content   string
	workspace string
	updatedAt time.Time
}

func New(content, workspace string) (Memory, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return Memory{}, ErrInvalidInput
	}

	id, err := newID()
	if err != nil {
		return Memory{}, fmt.Errorf("generate memory id: %w", err)
	}

	return Memory{
		id:        id,
		content:   trimmed,
		workspace: workspace,
		updatedAt: time.Now().UTC(),
	}, nil
}

func Load(id, content, workspace string, updatedAt time.Time) (Memory, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return Memory{}, ErrInvalidInput
	}
	return Memory{
		id:        id,
		content:   trimmed,
		workspace: workspace,
		updatedAt: updatedAt.UTC(),
	}, nil
}

type MemoryDTO struct {
	ID        string
	Content   string
	Workspace string
	UpdatedAt time.Time
}

func (m *Memory) DTO() MemoryDTO {
	return MemoryDTO{
		ID:        m.id,
		Content:   m.content,
		Workspace: m.workspace,
		UpdatedAt: m.updatedAt,
	}
}

func (m MemoryDTO) IsGlobal() bool {
	return m.Workspace == ""
}

func (m *Memory) ID() string {
	return m.id
}

func (m *Memory) Content() string {
	return m.content
}

func (m *Memory) Workspace() string {
	return m.workspace
}

func (m *Memory) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *Memory) SetContent(content string) (bool, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return false, ErrInvalidInput
	}
	if m.content == trimmed {
		return false, nil
	}
	m.content = trimmed
	m.touch()
	return true, nil
}

func (m *Memory) PromoteToGlobal() bool {
	if m.IsGlobal() {
		return false
	}
	m.workspace = ""
	m.touch()
	return true
}

func (m *Memory) IsGlobal() bool {
	return m.workspace == ""
}

func (m *Memory) touch() {
	m.updatedAt = time.Now().UTC()
}

func newID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
