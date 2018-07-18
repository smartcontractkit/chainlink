// +build sgx_enclave

package adapters

/*
#cgo LDFLAGS: -L../sgx/target/ -ladapters
#include "../sgx/libadapters/adapters.h"
#include "stdlib.h"
*/
import "C"

import (
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
	buffer := make([]byte, 8192)
	output := (*C.char)(unsafe.Pointer(&buffer[0]))
	bufferCapacity := C.int(len(buffer))
	outputLen := C.int(0)
	outputLenPtr := (*C.int)(unsafe.Pointer(&outputLen))

	fmt.Println("Input", input.JobRunID, input.Data, input.Amount)

	wasmCStr := C.CString(wasm.Wasm)
	defer C.free(unsafe.Pointer(wasmCStr))
	valueCStr := C.CString(input.Get("value").String())
	defer C.free(unsafe.Pointer(valueCStr))

	fmt.Println("value", input.Get("value"), "valueCstr", valueCStr)

	_, err := C.wasm(wasmCStr, valueCStr, output, bufferCapacity, outputLenPtr)
	if err != nil {
		return input.WithError(err)
	}

	outputStr := C.GoStringN(output, outputLen)
	fmt.Println("Wasm output", outputLen, outputStr)
	return input.WithValue(outputStr)
}
