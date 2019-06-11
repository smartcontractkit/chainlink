package adapters

import (
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/tidwall/gjson"
)

var evmFalse = "0x0000000000000000000000000000000000000000000000000000000000000000"
var evmTrue = "0x0000000000000000000000000000000000000000000000000000000000000001"

// EthBool holds no fields
type EthBool struct{}

// Perform returns the abi encoding for a boolean
//
// For example, after converting the result false to hex encoded Ethereum
// ABI, it would be:
// "0x0000000000000000000000000000000000000000000000000000000000000000"
func (*EthBool) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	r := input.Result()
	if boolean(r.Type) {
		input.CompleteWithResult(evmTrue)
		return input
	}
	input.CompleteWithResult(evmFalse)
	return input
}

func boolean(t gjson.Type) bool {
	switch t {
	case gjson.False, gjson.Null:
		return false
	default:
		return true
	}
}
