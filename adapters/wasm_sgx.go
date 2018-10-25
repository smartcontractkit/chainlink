// +build sgx_enclave

package adapters

/*
#cgo LDFLAGS: -L../sgx/target/ -ladapters
#include "../sgx/libadapters/adapters.h"
#include "stdlib.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Wasm represents a wasm binary encoded as base64 or wasm encoded as text (a lisp like language).
type Wasm struct {
	Wasm string `json:"wasm"`
}

// Perform ships the wasm representation to the SGX enclave where it is evaluated.
func (wasm *Wasm) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	adapterJSON, err := json.Marshal(wasm)
	if err != nil {
		return input.WithError(err)
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return input.WithError(err)
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

	_, err = C.wasm(cAdapter, cInput, output, bufferCapacity, outputLenPtr)
	if err != nil {
		return input.WithError(fmt.Errorf("SGX wasm: %v", err))
	}

	sgxResult := C.GoStringN(output, outputLen)
	var result models.RunResult
	if err := json.Unmarshal([]byte(sgxResult), &result); err != nil {
		return input.WithError(fmt.Errorf("unmarshaling SGX result: %v", err))
	}

	return result
}
