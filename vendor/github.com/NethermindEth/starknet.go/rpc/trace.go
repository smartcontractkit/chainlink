package rpc

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
)

// TraceTransaction returns the transaction trace for the given transaction hash.
//
// Parameters:
//   - ctx: the context.Context object for the request
//   - transactionHash: the transaction hash to trace
//
// Returns:
//   - TxnTrace: the transaction trace
//   - error: an error if the transaction trace cannot be retrieved
func (provider *Provider) TraceTransaction(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error) {
	var rawTxnTrace map[string]any
	if err := do(ctx, provider.c, "starknet_traceTransaction", &rawTxnTrace, transactionHash); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrHashNotFound, ErrNoTraceAvailable)
	}

	rawTraceByte, err := json.Marshal(rawTxnTrace)
	if err != nil {
		return nil, Err(InternalError, err)
	}

	switch rawTxnTrace["type"] {
	case string(TransactionType_Invoke):
		var trace InvokeTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, Err(InternalError, err)
		}
		return trace, nil
	case string(TransactionType_Declare):
		var trace DeclareTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, Err(InternalError, err)
		}
		return trace, nil
	case string(TransactionType_DeployAccount):
		var trace DeployAccountTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, Err(InternalError, err)
		}
		return trace, nil
	case string(TransactionType_L1Handler):
		var trace L1HandlerTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, Err(InternalError, err)
		}
		return trace, nil
	}
	return nil, Err(InternalError, "Unknown transaction type")

}

// TraceBlockTransactions retrieves the traces of transactions in a given block.
//
// Parameters:
// - ctx: the context.Context object for controlling the request
// - blockHash: the hash of the block to retrieve the traces from
// Returns:
// - []Trace: a slice of Trace objects representing the traces of transactions in the block
// - error: an error if there was a problem retrieving the traces.
func (provider *Provider) TraceBlockTransactions(ctx context.Context, blockID BlockID) ([]Trace, error) {
	var output []Trace
	if err := do(ctx, provider.c, "starknet_traceBlockTransactions", &output, blockID); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrBlockNotFound)
	}
	return output, nil

}

// SimulateTransactions simulates transactions on the blockchain.
// Simulate a given sequence of transactions on the requested state, and generate the execution traces.
// Note that some of the transactions may revert, in which case no error is thrown, but revert details can be seen on the returned trace object.
// Note that some of the transactions may revert, this will be reflected by the revert_error property in the trace. Other types of failures (e.g. unexpected error or failure in the validation phase) will result in TRANSACTION_EXECUTION_ERROR.
func (provider *Provider) SimulateTransactions(ctx context.Context, blockID BlockID, txns []Transaction, simulationFlags []SimulationFlag) ([]SimulatedTransaction, error) {

	var output []SimulatedTransaction
	if err := do(ctx, provider.c, "starknet_simulateTransactions", &output, blockID, txns, simulationFlags); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTxnExec, ErrBlockNotFound)
	}

	return output, nil

}
