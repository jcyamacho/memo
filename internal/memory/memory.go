package memory

import (
	"fmt"
	"strings"
	"time"
)

type Memory struct {
	id        ID
	content   string
	workspace string
	updatedAt time.Time
}

func New(content, workspace string) (Memory, error) {
	normalized, err := normalizeContent(content)
	if err != nil {
		return Memory{}, err
	}

	id, err := NewID()
	if err != nil {
		return Memory{}, fmt.Errorf("generate memory id: %w", err)
	}

	return Memory{
		id:        id,
		content:   normalized,
		workspace: workspace,
		updatedAt: time.Now().UTC(),
	}, nil
}

func Load(id, content, workspace string, updatedAt time.Time) (Memory, error) {
	parsedID, err := ParseID(id)
	if err != nil {
		return Memory{}, err
	}

	normalized, err := normalizeContent(content)
	if err != nil {
		return Memory{}, err
	}
	return Memory{
		id:        parsedID,
		content:   normalized,
		workspace: workspace,
		updatedAt: updatedAt.UTC(),
	}, nil
}

type MemoryDTO struct {
	ID        ID
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

func (m *Memory) ID() ID {
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
	normalized, err := normalizeContent(content)
	if err != nil {
		return false, err
	}
	if m.content == normalized {
		return false, nil
	}
	m.content = normalized
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

func normalizeContent(content string) (string, error) {
	normalized := strings.TrimSpace(content)
	if normalized == "" {
		return "", ErrInvalidInput
	}

	return normalized, nil
}
