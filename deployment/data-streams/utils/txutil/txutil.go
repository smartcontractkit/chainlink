package txutil

import (
	"context"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

// PreparedTx represents a transaction that was prepared but not sent to the chain. This is intended to be
// either signed and executed directly or bundled into an MCMS operation. The internal Tx should be unsigned.
type PreparedTx struct {
	Tx                 *gethtypes.Transaction // unsigned transaction
	ChainSelector      uint64
	DestinationAddress string
	ContractType       string
	Tags               []string
}

type ExecuteTxResult struct {
	Tx          PreparedTx
	BlockNumber uint64
}

// SignAndExecute signs and then executes transactions directly on the chain with the given deployer key configured
// for the chain. The transactions should not be already sent to the chain.
func SignAndExecute(e deployment.Environment, preparedTxs []PreparedTx) ([]ExecuteTxResult, error) {
	for _, tx := range preparedTxs {
		chain := e.Chains[tx.ChainSelector]
		signedTx, err := chain.DeployerKey.Signer(chain.DeployerKey.From, tx.Tx)
		if err != nil {
			return nil, err
		}
		tx.Tx = signedTx
	}
	return Execute(e, preparedTxs)
}

// Execute executes the prepared transactions directly on the chain
// the transactions should not be already sent to the chain
func Execute(e deployment.Environment, preparedTxs []PreparedTx) ([]ExecuteTxResult, error) {
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
