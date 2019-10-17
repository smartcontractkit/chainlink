// +build sgx_enclave

package adapters

/*
#cgo LDFLAGS: -L../sgx/target/ -ladapters
#include <stdlib.h>
#include "../../sgx/libadapters/adapters.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Multiplier represents the number to multiply by in Multiply adapter.
type Multiplier float64

// UnmarshalJSON implements json.Unmarshaler.
func (m *Multiplier) UnmarshalJSON(input []byte) error {
	input = utils.RemoveQuotes(input)
	times, err := strconv.ParseFloat(string(input), 64)
	if err != nil {
		return fmt.Errorf("cannot parse into float: %s", input)
	}

	*m = Multiplier(times)

	return nil
}

// Multiply holds the a number to multiply the given value by.
type Multiply struct {
	Times Multiplier `json:"times"`
}

// Perform returns the input's "result" field, multiplied times the adapter's
// "times" field.
//
// For example, if input value is "99.994" and the adapter's "times" is
// set to "100", the result's value will be "9999.4".
func (ma *Multiply) Perform(input models.RunInput, _ *store.Store) models.RunOutput {
	adapterJSON, err := json.Marshal(ma)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return models.NewRunOutputError(err)
	}

	cAdapter := C.CString(string(adapterJSON))
	defer C.free(unsafe.Pointer(cAdapter))
	cInput := C.CString(string(inputJSON))
	defer C.free(unsafe.Pointer(cInput))

	buffer := make([]byte, 8192)
	output := (*C.char)(unsafe.Pointer(&buffer[0]))
	bufferCapacity := C.int(len(buffer))
	outputLen := C.int(0)
	outputLenPtr := (*C.int)(unsafe.Pointer(&outputLen))

	if _, err = C.multiply(cAdapter, cInput, output, bufferCapacity, outputLenPtr); err != nil {
		return models.NewRunOutputError(fmt.Errorf("SGX multiply: %v", err))
	}

	sgxResult := C.GoStringN(output, outputLen)
	var result models.RunResult
	if err := json.Unmarshal([]byte(sgxResult), &result); err != nil {
		return models.NewRunOutputError(fmt.Errorf("unmarshaling SGX result: %v", err))
	}

	return result
}
