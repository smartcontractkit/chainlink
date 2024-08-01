package rpc

import (
	"context"
)

// AddInvokeTransaction adds an invoke transaction to the provider.
//
// Parameters:
// - ctx: The context for the function.
// - invokeTxn: The invoke transaction to be added.
// Returns:
// - *AddInvokeTransactionResponse: the response of adding the invoke transaction
// - error: an error if any
func (provider *Provider) AddInvokeTransaction(ctx context.Context, invokeTxn BroadcastInvokeTxnType) (*AddInvokeTransactionResponse, error) {
	var output AddInvokeTransactionResponse
	if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, invokeTxn); err != nil {
		return nil, tryUnwrapToRPCErr(
			err,
			ErrInsufficientAccountBalance,
			ErrInsufficientMaxFee,
			ErrInvalidTransactionNonce,
			ErrValidationFailure,
			ErrNonAccount,
			ErrDuplicateTx,
			ErrUnsupportedTxVersion,
			ErrUnexpectedError,
		)
	}
	return &output, nil
}

// AddDeclareTransaction submits a declare transaction to the StarkNet contract.
//
// Parameters:
// - ctx: The context.Context object for the request.
// - declareTransaction: The input for the declare transaction.
// Returns:
// - *AddDeclareTransactionResponse: The response of submitting the declare transaction
// - error: an error if any
func (provider *Provider) AddDeclareTransaction(ctx context.Context, declareTransaction BroadcastDeclareTxnType) (*AddDeclareTransactionResponse, error) {

	switch txn := declareTransaction.(type) {
	case DeclareTxnV2:
		// DeclareTxnV2 should not have a populated class hash field. It is only needed for signing.
		txn.ClassHash = nil
		declareTransaction = txn
	}

	var result AddDeclareTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeclareTransaction", &result, declareTransaction); err != nil {
		return nil, tryUnwrapToRPCErr(
			err,
			ErrClassAlreadyDeclared,
			ErrCompilationFailed,
			ErrCompiledClassHashMismatch,
			ErrInsufficientAccountBalance,
			ErrInsufficientMaxFee,
			ErrInvalidTransactionNonce,
			ErrValidationFailure,
			ErrNonAccount,
			ErrDuplicateTx,
			ErrContractClassSizeTooLarge,
			ErrUnsupportedTxVersion,
			ErrUnsupportedContractClassVersion,
		)
	}
	return &result, nil
}

// AddDeployAccountTransaction adds a DEPLOY_ACCOUNT transaction to the provider.
//
// Parameters:
// - ctx: The context of the function
// - deployAccountTransaction: The deploy account transaction to be added
// Returns:
// - *AddDeployAccountTransactionResponse: the response of adding the deploy account transaction or an error
func (provider *Provider) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction BroadcastAddDeployTxnType) (*AddDeployAccountTransactionResponse, error) {
	var result AddDeployAccountTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeployAccountTransaction", &result, deployAccountTransaction); err != nil {
		return nil, tryUnwrapToRPCErr(
			err,
			ErrInsufficientAccountBalance,
			ErrInsufficientMaxFee,
			ErrInvalidTransactionNonce,
			ErrValidationFailure,
			ErrNonAccount,
			ErrClassHashNotFound,
			ErrDuplicateTx,
			ErrUnsupportedTxVersion,
		)
	}
	return &result, nil
}
