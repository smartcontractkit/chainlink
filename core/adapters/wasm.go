// +build !sgx_enclave

package adapters

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// Wasm represents a wasm binary encoded as base64 or wasm encoded as text (a lisp like language).
type Wasm struct {
	WasmT string `json:"wasmt"`
}

// Perform ships the wasm representation to the SGX enclave where it is evaluated.
func (wasm *Wasm) Perform(_ models.JSON, result models.RunResult, _ *store.Store) models.RunResult {
	result.SetError(fmt.Errorf("Wasm is not supported without SGX"))
	return result
}
