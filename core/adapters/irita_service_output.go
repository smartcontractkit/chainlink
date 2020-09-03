package adapters

import (
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type IritaServiceOutput struct {
	Output Output
}

type Output struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

func (iso *IritaServiceOutput) TaskType() models.TaskType {
	return TaskTypeIritaServiceOutput
}

func (iso *IritaServiceOutput) Perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	iso.Output = Output{
		Header: "",
		Body:   input.Result().String(),
	}

	return models.NewRunOutputCompleteWithResult(iso.Output.Body)
}
