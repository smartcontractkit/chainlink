package adapters

import (
	"github.com/smartcontractkit/chainlink-go/models"
)

type NoOp struct {
}

func (self *NoOp) Perform(input models.RunResult) models.RunResult {
	return models.RunResult{}
}
