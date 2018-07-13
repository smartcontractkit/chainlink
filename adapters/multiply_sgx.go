// +build sgx_enclave

package adapters

/*
#cgo LDFLAGS: -L../sgx/target/release/ -ladapters
#include "../sgx/libadapters/adapters.h"
*/
import "C"

import (
	"fmt"
	"strconv"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Multiplier represents the number to multiply by in Multiply adapter.
type Multiplier float64

// UnmarshalJSON implements json.Unmarshaler.
func (m *Multiplier) UnmarshalJSON(input []byte) error {
	if isString(input) {
		input = input[1 : len(input)-1]
	}

	times, err := strconv.ParseFloat(string(input), 64)
	if err != nil {
		return fmt.Errorf("cannot parse into float: %s", input)
	}

	*m = Multiplier(times)

	return nil
}

func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// Multiply holds the a number to multiply the given value by.
type Multiply struct {
	Times Multiplier `json:"times"`
}

// Perform returns the input's "value" field, multiplied times the adapter's
// "times" field.
//
// For example, if input value is "99.994" and the adapter's "times" is
// set to "100", the result's value will be "9999.4".
func (ma *Multiply) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	multiplicand := C.CString(strconv.FormatUint(uint64(ma.Times), 10))
	multiplier := C.CString(input.Data.Get("value").String())
	body, err := C.multiply(multiplicand, multiplier)
	if err != nil {
		return input.WithError(fmt.Errorf(C.GoString(body)))
	}
	return input.WithValue(C.GoString(body))
}
