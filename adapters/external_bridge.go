package adapters

import (
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

type ExternalBridge struct {
	*models.CustomTaskType
}

func (eb *ExternalBridge) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	return input
}
