package txutil

import (
	"context"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment"
)

// GeneratedTx represents a transaction that was generated but not sent to the chain. Can extend to include metadata.
type GeneratedTx struct {
	Tx                 *gethtypes.Transaction
	ChainSelector      uint64
	DestinationAddress string
	ContractType       string
}

type TxExecuteResult struct {
	Tx          GeneratedTx
	BlockNumber uint64
}

// ExecuteTransactions executes transactions directly on the chain with the given deployer address
// the transactions should not be already sent to the chain
func ExecuteTransactions(e deployment.Environment, generatedTxs []GeneratedTx) ([]TxExecuteResult, error) {
	var txExecuteResults []TxExecuteResult
	for _, tx := range generatedTxs {
		chain := e.Chains[tx.ChainSelector]
		err := chain.Client.SendTransaction(context.Background(), tx.Tx)
		if err != nil {
			return nil, err
		}
		blockNumber, err := chain.Confirm(tx.Tx)
		if err != nil {
			return nil, err
		}
		txExecuteResults = append(txExecuteResults, TxExecuteResult{Tx: tx, BlockNumber: blockNumber})
	}

	return txExecuteResults, nil
}
