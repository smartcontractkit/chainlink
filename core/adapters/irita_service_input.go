package adapters

import (
	simplejson "github.com/bitly/go-simplejson"

	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type IritaServiceInput struct {
	PathKey string `json:"pathKey"`
}

func (isi *IritaServiceInput) TaskType() models.TaskType {
	return TaskTypeIritaServiceInput
}

func (isi *IritaServiceInput) Perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	// var httpExtendedPath ExtendedPath
	serviceInput := isi.GetServiceInput(input.JobRunID().String())

	js, err := simplejson.NewJson([]byte(serviceInput))
	if err != nil {
		return models.NewRunOutputError(err)
	}

	value, err := dig(js, []string{isi.PathKey})
	if err != nil {
		return models.NewRunOutputError(err)
	}

	return models.NewRunOutputCompleteWithResult(value.MustString())
}

func (isi *IritaServiceInput) GetServiceInput(jobRunID string) string {
	return strpkg.GetServiceMemory()[jobRunID].RequestResponse.Input
}
