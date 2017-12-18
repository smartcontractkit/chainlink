package adapters

import (
	"github.com/smartcontractkit/chainlink-go/models"
)

type NoOp struct {
	AdapterBase
}

func (self *NoOp) Perform(input models.RunResult) models.RunResult {
	return models.RunResult{}
}
