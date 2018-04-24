package presenters

import (
	"encoding/json"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/store/models"
)

type JobSpecsDocument struct {
	Data []models.JobSpec
	jsonapi.Links
}

func (js *JobSpecsDocument) UnmarshalJSON(input []byte) error {
	// First unmarshal using the jsonAPI into the JobSpec slice, as is api2go will discard the links
	err := jsonapi.Unmarshal(input, &js.Data)
	if err != nil {
		return err
	}

	// Unmarshal using the stdlib Unmarshal to extract the Links part of the document
	document := jsonapi.Document{}
	err = json.Unmarshal(input, &document)
	if err != nil {
		return err
	}
	js.Links = document.Links

	return nil
}
