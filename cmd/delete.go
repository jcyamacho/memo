package cmd

import (
	"encoding/xml"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id> [id...]",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete one or more memories",
	Long: `Delete memories confirmed obsolete or incorrect.

Accepts 1 to 50 memory ids. Duplicate ids are ignored. The command prints XML
with one result per requested id that remains after deduplication.`,
	Example: `  memo delete abc123
  memo delete abc123 def456 missing-id`,
	Args: cobra.RangeArgs(1, 50),
	RunE: func(cmd *cobra.Command, args []string) error {
		outcomes, err := service.Delete(cmd.Context(), args)
		if err != nil {
			return err
		}
		return printOutput(deleteResultsXMLOutput(outcomes))
	},
}

type deleteResultsXML struct {
	XMLName xml.Name          `xml:"delete_results"`
	Deleted int               `xml:"deleted,attr"`
	Failed  int               `xml:"failed,attr"`
	Items   []deleteResultXML `xml:",any"`
}

type deleteResultXML struct {
	Deleted *deletedXML
	Failure *failureXML
}

type deletedXML struct {
	XMLName xml.Name `xml:"deleted"`
	ID      string   `xml:"id,attr"`
}

type failureXML struct {
	XMLName xml.Name `xml:"failure"`
	ID      string   `xml:"id,attr"`
	Status  string   `xml:"status,attr"`
}

func (r deleteResultXML) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if r.Deleted != nil {
		return e.Encode(r.Deleted)
	}
	if r.Failure != nil {
		return e.Encode(r.Failure)
	}
	return nil
}

func deleteResultsXMLOutput(outcomes []memory.DeleteOutcome) (string, error) {
	dto := deleteResultsXML{
		Items: make([]deleteResultXML, 0, len(outcomes)),
	}

	for _, outcome := range outcomes {
		if outcome.Deleted {
			dto.Deleted++
			dto.Items = append(dto.Items, deleteResultXML{Deleted: &deletedXML{ID: outcome.ID}})
			continue
		}

		dto.Items = append(dto.Items, deleteResultXML{Failure: &failureXML{
			ID:     outcome.ID,
			Status: string(outcome.Code),
		}})
	}

	dto.Failed = len(outcomes) - dto.Deleted
	return marshalXML(dto)
}
