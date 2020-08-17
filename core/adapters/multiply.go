package adapters

import (
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/shopspring/decimal"
)

// Multiply holds the a number to multiply the given value by.
type Multiply struct {
	Times *decimal.Decimal `json:"times,omitempty"`
}

// TaskType returns the type of Adapter.
func (m *Multiply) TaskType() models.TaskType {
	return TaskTypeMultiply
}
