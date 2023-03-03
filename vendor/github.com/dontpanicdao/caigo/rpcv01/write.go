package rpcv01

import (
	"context"
	"fmt"
	"math/big"

	ctypes "github.com/dontpanicdao/caigo/types"
)

// AddInvokeTransaction estimates the fee for a given StarkNet transaction.
func (provider *Provider) AddInvokeTransaction(ctx context.Context, call ctypes.FunctionCall, signature []string, maxFee string, version string, nonce *big.Int) (*ctypes.AddInvokeTransactionOutput, error) {
	if call.EntryPointSelector != "" {
		call.EntryPointSelector = fmt.Sprintf("0x%s", ctypes.GetSelectorFromName(call.EntryPointSelector).Text(16))
	}
	var output ctypes.AddInvokeTransactionOutput
	if nonce == nil {
		if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, call, signature, maxFee, version); err != nil {
			return nil, err
		}
		return &output, nil
	}
	if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, call, signature, maxFee, version, fmt.Sprintf("0x%s", nonce.Text(16))); err != nil {
		return nil, err
	}
	return &output, nil
}

// AddDeclareTransaction submits a new class declaration transaction.
func (provider *Provider) AddDeclareTransaction(ctx context.Context, contractClass ctypes.ContractClass, version string) (*AddDeclareTransactionOutput, error) {
	var result AddDeclareTransactionOutput
	if err := do(ctx, provider.c, "starknet_addDeclareTransaction", &result, contractClass, version); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddDeployTransaction allows to declare a class and instantiate the
// associated contract in one command. This function will be deprecated and
// replaced by AddDeclareTransaction to declare a class, followed by
// AddInvokeTransaction to instantiate the contract. For now, it remains the only
// way to deploy an account without being charged for it.
func (provider *Provider) AddDeployTransaction(ctx context.Context, salt string, inputs []string, contractClass ctypes.ContractClass) (*AddDeployTransactionOutput, error) {
	var result AddDeployTransactionOutput
	if err := do(ctx, provider.c, "starknet_addDeployTransaction", &result, salt, inputs, contractClass); err != nil {
		return nil, err
	}
	return &result, nil
}
