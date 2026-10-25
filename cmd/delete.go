package cmd

import (
	"encoding/xml"
	"errors"

	"github.com/jcyamacho/memo/internal/memory"
	"github.com/spf13/cobra"
)

var errDeleteIncomplete = errors.New("one or more memories could not be deleted")

func newDeleteCommand(service *memory.Service) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <id> [id...]",
		Aliases: []string{"remove", "rm"},
		Short:   "Delete one or more memories",
		Long: `Delete memories confirmed obsolete or incorrect by id.

Accepts 1 to 50 memory ids. Duplicate ids are ignored. The command prints XML
with one result per requested id that remains after deduplication. If any item
fails, it prints all results before returning a non-zero exit status.`,
		Example: `  memo delete 550e8400-e29b-41d4-a716-446655440000
  memo delete 550e8400-e29b-41d4-a716-446655440000 550e8400-e29b-41d4-a716-446655440001`,
		Args: cobra.RangeArgs(1, 50),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids := make([]memory.ID, 0, len(args))
			for _, value := range args {
				id, err := parseMemoryID(value)
				if err != nil {
					return err
				}
				ids = append(ids, id)
			}
			outcomes, deleteErr := service.Delete(cmd.Context(), ids)

			output, err := deleteResultsXMLOutput(outcomes)
			if err != nil {
				return err
			}
			if err := printOutput(cmd, output); err != nil {
				return err
			}
			if deleteErr != nil {
				return deleteErr
			}
			for _, outcome := range outcomes {
				if !outcome.Deleted {
					return errDeleteIncomplete
				}
			}

			return nil
		},
	}
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
			dto.Items = append(dto.Items, deleteResultXML{Deleted: &deletedXML{ID: outcome.ID.String()}})
			continue
		}

		dto.Items = append(dto.Items, deleteResultXML{Failure: &failureXML{
			ID:     outcome.ID.String(),
			Status: string(outcome.Code),
		}})
	}

	dto.Failed = len(outcomes) - dto.Deleted
	return marshalXML(dto)
}
