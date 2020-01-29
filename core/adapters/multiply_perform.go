// +build !sgx_enclave

package adapters

import (
	"chainlink/core/store"
	"chainlink/core/store/models"
	"fmt"
	"math/big"
)

// Perform returns the input's "result" field, multiplied times the adapter's
// "times" field.
//
// For example, if input value is "99.994" and the adapter's "times" is
// set to "100", the result's value will be "9999.4".
func (ma *Multiply) Perform(input models.RunInput, _ *store.Store) models.RunOutput {
	val := input.Result()
	i, ok := (&big.Float{}).SetString(val.String())
	if !ok {
		return models.NewRunOutputError(fmt.Errorf("cannot parse into big.Float: %v", val.String()))
	}

	if ma.Times != nil {
		i.Mul(i, big.NewFloat(float64(*ma.Times)))
	}
	return models.NewRunOutputCompleteWithResult(i.String())
}
