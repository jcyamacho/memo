package memory

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type Memory struct {
	ID        string
	Content   string
	Workspace string
	UpdatedAt time.Time
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
		ID:        id,
		Content:   trimmed,
		Workspace: workspace,
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (m *Memory) SetContent(content string) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ErrInvalidInput
	}
	m.Content = trimmed
	return nil
}

func (m *Memory) PromoteToGlobal() {
	m.Workspace = ""
}

func (m Memory) IsGlobal() bool {
	return m.Workspace == ""
}

func newID() (string, error) {
	var buf [12]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf[:]), nil
}
