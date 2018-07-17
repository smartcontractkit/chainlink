// +build sgx_enclave

package adapters

/*
#cgo LDFLAGS: -L../sgx/target/release/ -ladapters
#include "../sgx/libadapters/adapters.h"
*/
import "C"

import (
	"encoding/json"
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
	adapter_json, err := json.Marshal(ma)
	if err != nil {
		return input.WithError(err)
	}
	input_json, err := json.Marshal(input)
	if err != nil {
		return input.WithError(err)
	}

	cAdapter := C.CString(string(adapter_json))
	defer C.free(unsafe.Pointer(cAdapter))
	cInput := C.CString(string(input_json))
	defer C.free(unsafe.Pointer(cInput))
	sgxResult := C.CString(stringOfLength(4096))
	defer C.free(unsafe.Pointer(sgxResult))

	_, err = C.multiply(cAdapter, cInput, sgxResult)
	if err != nil {
		return input.WithError(fmt.Errorf("CGO multiply: %v", err))
	}

	return input.WithValue(C.GoString(sgxResult))
}

func stringOfLength(l int) string {
	var s string
	for i := 0; i < l; i++ {
		s = s + " "
	}
	return s
}
