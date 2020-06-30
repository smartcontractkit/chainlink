// +build !sgx_enclave

package adapters

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Perform returns the input's "result" field, multiplied times the adapter's
// "times" field.
//
// For example, if input value is "99.994" and the adapter's "times" is
// set to "100", the result's value will be "9999.4".
func (ma *Multiply) Perform(input models.RunInput, store *store.Store) models.RunOutput {
	var result string
	val := input.Result()
	dec, err := decimal.NewFromString(val.String())
	if err != nil {
		return models.NewRunOutputError(errors.Wrapf(err, "cannot parse into big.Float: %v", val.String()))
	}
	if ma.Times != nil {
		dec = dec.Mul(*ma.Times)
		if store != nil && store.Config.LeetMode() {
			dec = dec.Round(0)
			result = fmt.Sprintf("%v.1337", dec.String())
		} else {
			result = dec.String()
		}
	}
	return models.NewRunOutputCompleteWithResult(result)
}
