package adapters

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Multiply holds the a number to multiply the given value by.
type Multiply struct {
	Times int64 `json:"times"`
}

// Perform returns the input's "value" field, multiplied times the adapter's
// "times" field.
//
// For example, if input value is "99.994" and the adapter's "times" is
// set to "100", the result's value will be "9999.4".
func (ma *Multiply) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	val, err := input.Get("value")
	if err != nil {
		return models.RunResultWithError(err)
	}

	i, ok := (&big.Float{}).SetString(val.String())
	if !ok {
		return models.RunResultWithError(fmt.Errorf("cannot parse into big.Float: %v", val.String()))
	}
	res := i.Mul(i, big.NewFloat(float64(ma.Times)))

	return input.WithValue(res.String())
}
