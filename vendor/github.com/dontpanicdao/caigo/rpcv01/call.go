package rpcv01

import (
	"context"

	ctypes "github.com/dontpanicdao/caigo/types"
)

// Call a starknet function without creating a StarkNet transaction.
func (provider *Provider) Call(ctx context.Context, call ctypes.FunctionCall, block BlockID) ([]string, error) {
	call.EntryPointSelector = ctypes.BigToHex(ctypes.GetSelectorFromName(call.EntryPointSelector))
	if len(call.Calldata) == 0 {
		call.Calldata = make([]string, 0)
	}
	var result []string
	if err := do(ctx, provider.c, "starknet_call", &result, call, block); err != nil {
		return nil, err
	}
	return result, nil
}
