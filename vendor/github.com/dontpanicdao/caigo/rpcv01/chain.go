package rpcv01

import (
	"context"

	ctypes "github.com/dontpanicdao/caigo/types"
)

// ChainID retrieves the current chain ID for transaction replay protection.
func (provider *Provider) ChainID(ctx context.Context) (string, error) {
	var result string
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := provider.c.CallContext(ctx, &result, "starknet_chainId", []interface{}{}...); err != nil {
		return "", err
	}
	return ctypes.HexToShortStr(result), nil
}

// Syncing checks the syncing status of the node.
func (provider *Provider) Syncing(ctx context.Context) (*SyncResponse, error) {
	var result SyncResponse
	// Note: []interface{}{}...force an empty `params[]` in the jsonrpc request
	if err := provider.c.CallContext(ctx, &result, "starknet_syncing", []interface{}{}...); err != nil {
		return nil, err
	}
	return &result, nil
}

// StateUpdate gets the information about the result of executing the requested block.
func (provider *Provider) StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, error) {
	var state StateUpdateOutput
	if err := do(ctx, provider.c, "starknet_getStateUpdate", &state, blockID); err != nil {
		return nil, err
	}
	return &state, nil
}
