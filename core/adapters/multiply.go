package adapters

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// Multiply holds the a number to multiply the given value by.
type Multiply struct {
	Times *decimal.Decimal `json:"times,omitempty"`
}

// TaskType returns the type of Adapter.
func (m *Multiply) TaskType() models.TaskType {
	return TaskTypeMultiply
}

// Perform returns the input's "result" field, multiplied times the adapter's
// "times" field.
//
// For example, if input value is "99.994" and the adapter's "times" is
// set to "100", the result's value will be "9999.4".
func (m *Multiply) Perform(input models.RunInput, _ *store.Store, _ *keystore.Master) models.RunOutput {
	val := input.Result()
	dec, err := decimal.NewFromString(val.String())
	if err != nil {
		return models.NewRunOutputError(errors.Wrapf(err, "cannot parse into big.Float: %s", val.String()))
	}
	if m.Times != nil {
		dec = dec.Mul(*m.Times)
	}
	return models.NewRunOutputCompleteWithResult(dec.String(), input.ResultCollection())
}
