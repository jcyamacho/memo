package cmd

import (
	"encoding/xml"
	"fmt"

	"github.com/jcyamacho/memo/internal/memory"
)

type memoryXML struct {
	XMLName   xml.Name `xml:"memory"`
	ID        string   `xml:"id,attr"`
	UpdatedAt string   `xml:"updated_at,attr"`
	Global    *bool    `xml:"global,attr,omitempty"`
	Workspace string   `xml:"workspace,attr,omitempty"`
	Content   string   `xml:",chardata"`
}

func memoryXMLOutput(m memory.MemoryDTO) (string, error) {
	return marshalXML(newMemoryXML(m, ""))
}

func newMemoryXML(m memory.MemoryDTO, skipWorkspaceIfEquals string) memoryXML {
	dto := memoryXML{
		ID:        m.ID,
		UpdatedAt: m.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		Content:   m.Content,
	}

	if m.IsGlobal() {
		value := true
		dto.Global = &value
	} else if m.Workspace != skipWorkspaceIfEquals {
		dto.Workspace = m.Workspace
	}

	return dto
}

func marshalXML(value any) (string, error) {
	output, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal XML: %w", err)
	}
	return string(output), nil
}
