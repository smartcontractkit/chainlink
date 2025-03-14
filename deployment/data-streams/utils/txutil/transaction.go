package txutil

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

// PreparedTx represents a transaction that was prepared but not sent to the chain. This is intended to be
// either signed and executed directly or bundled into an MCMS operation.
type PreparedTx struct {
	Tx            *gethtypes.Transaction
	ChainSelector uint64
	ContractType  string
	Tags          []string
}

type ExecuteTxResult struct {
	Tx          *PreparedTx
	BlockNumber uint64
}

// SignAndExecute signs and then executes transactions directly on the chain with the given deployer key configured
// for the chain. The transactions should not be already sent to the chain.
func SignAndExecute(e deployment.Environment, preparedTxs []*PreparedTx) ([]ExecuteTxResult, error) {
	var executeTxResults []ExecuteTxResult
	// To execute the txs in parallel this would need to batch up the txs by chain to avoid nonce issues
	for _, tx := range preparedTxs {
		chain, exists := e.Chains[tx.ChainSelector]
		if !exists {
			return executeTxResults, fmt.Errorf("chain not found in env %d", tx.ChainSelector)
		}
		// The environment level context `e.GetContext()` is also an option to use here
		ctx := chain.DeployerKey.Context
		reconfiguredTx, err := reconfigureTx(ctx, chain, tx)
		if err != nil {
			return executeTxResults, fmt.Errorf("chain %d: failed to reconfigure transaction: %w", chain.Selector, err)
		}
		signedTx, err := chain.DeployerKey.Signer(chain.DeployerKey.From, reconfiguredTx)
		if err != nil {
			return executeTxResults, fmt.Errorf("chain %d: failed to sign transaction: %w", chain.Selector, err)
		}
		tx.Tx = signedTx
		err = chain.Client.SendTransaction(ctx, tx.Tx)
		if err != nil {
			return executeTxResults, fmt.Errorf("chain %d: failed to send transaction: %w", chain.Selector, err)
		}
		blockNumber, err := chain.Confirm(tx.Tx)
		if err != nil {
			return executeTxResults, fmt.Errorf("chain %d: failed to confirm transaction: %w", chain.Selector, err)
		}
		e.Logger.Infow("Transaction confirmed", "blockNumber", blockNumber, "tx", tx)
		executeTxResults = append(executeTxResults, ExecuteTxResult{Tx: tx, BlockNumber: blockNumber})
	}
	return executeTxResults, nil
}

// reconfigureTx takes the tx `call data` and reconfigures the transaction to use valid nonce, gas price and gas limit
func reconfigureTx(ctx context.Context, chain deployment.Chain, preparedTx *PreparedTx) (*gethtypes.Transaction, error) {
	nonce, err := chain.Client.NonceAt(ctx, chain.DeployerKey.From, nil)
	if err != nil {
		return nil, fmt.Errorf("chain %d: failed to get nonce: %w", chain.Selector, err)
	}
	gasPrice, err := chain.Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain %d: failed to get gas price: %w", chain.Selector, err)
	}
	estimate, err := chain.Client.EstimateGas(ctx, ethereum.CallMsg{
		From: chain.DeployerKey.From,
		To:   preparedTx.Tx.To(),
		Data: preparedTx.Tx.Data(),
	})
	if err != nil {
		return nil, fmt.Errorf("chain %d: failed to estimate gas: %w", chain.Selector, err)
	}

	rawTx := gethtypes.NewTx(&gethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      estimate,
		To:       preparedTx.Tx.To(),
		Value:    preparedTx.Tx.Value(),
		Data:     preparedTx.Tx.Data(),
	})
	return rawTx, nil
}
