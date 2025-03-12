package txutil

import (
	"context"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

// PreparedTx represents a transaction that was prepared but not sent to the chain. This is intended to be
// either executed directly or bundled into an MCMS operation
type PreparedTx struct {
	Tx                 *gethtypes.Transaction
	ChainSelector      uint64
	DestinationAddress string
	ContractType       string
	Tags               []string
}

type ExecuteTxResult struct {
	Tx          PreparedTx
	BlockNumber uint64
}

// ExecuteTransactions executes transactions directly on the chain with the given deployer address
// the transactions should not be already sent to the chain
func ExecuteTransactions(e deployment.Environment, preparedTxs []PreparedTx) ([]ExecuteTxResult, error) {
	var executeTxResults []ExecuteTxResult
	for _, tx := range preparedTxs {
		chain := e.Chains[tx.ChainSelector]
		err := chain.Client.SendTransaction(context.Background(), tx.Tx)
		if err != nil {
			return nil, err
		}
		blockNumber, err := chain.Confirm(tx.Tx)
		if err != nil {
			return nil, err
		}
		executeTxResults = append(executeTxResults, ExecuteTxResult{Tx: tx, BlockNumber: blockNumber})
	}

	return executeTxResults, nil
}
