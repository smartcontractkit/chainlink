// +build sgx_enclave

package adapters

/*
#cgo LDFLAGS: -L../sgx/target/release/ -ladapters
#include "../sgx/libadapters/adapters.h"
*/
import "C"

import (
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Wasm represents a wasm binary encoded as base64 or wasm encoded as text (a lisp like language).
type Wasm struct {
	WasmT string `json:"wasmt"`
}

// Perform ships the wasm representation to the SGX enclave where it is evaluated.
func (wasm *Wasm) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	result, err := C.wasm(C.CString(wasm.WasmT))
	if err != nil {
		return input.WithError(err)
	}
	return input.WithValue(C.GoString(result))
}
