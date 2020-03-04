package adapters

import (
	"github.com/shopspring/decimal"
)

// Multiply holds the a number to multiply the given value by.
type Multiply struct {
	Times *decimal.Decimal `json:"times,omitempty"`
}
