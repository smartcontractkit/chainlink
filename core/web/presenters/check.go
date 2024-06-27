package presenters

import (
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

type Check struct {
	JAID
	Name   string          `json:"name"`
	Status services.Status `json:"status"`
	Output string          `json:"output"`
}

func (c Check) GetName() string {
	return "checks"
}
