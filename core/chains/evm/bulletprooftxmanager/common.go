package bulletprooftxmanager

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// Tries to send transactions in batches. Even if some batch(es) fail to get sent, it tries all remaining batches,
// before returning with error for the latest batch send. If a batch send fails, this sets the error on all
// elements in that batch.
// TODO: The batch send is just sending the batch to 1 RPC node. Change it to send to all nodes,
// similar to how Pool.SendTransaction sends a transaction to all nodes.
func batchSendTransactions(ctx context.Context, attempts []EthTxAttempt, batchSize int, logger logger.Logger,
	ethClient evmclient.Client) (reqs []rpc.BatchElem, err error) {
	if len(attempts) == 0 {
		return
	}

	reqs = make([]rpc.BatchElem, len(attempts))
	ethTxIDs := make([]int64, len(attempts))
	for i, attempt := range attempts {
		ethTxIDs[i] = attempt.EthTxID
		req := rpc.BatchElem{
			Method: "eth_sendRawTransaction",
			Args:   []interface{}{hexutil.Encode(attempt.SignedRawTx)},
			Result: &common.Hash{},
		}
		reqs[i] = req
	}

	if batchSize == 0 {
		batchSize = len(reqs)
	}
	for i := 0; i < len(reqs); i += batchSize {
		j := i + batchSize
		if j > len(reqs) {
			j = len(reqs)
		}

		logger.Debugw(fmt.Sprintf("Batch sending transactions %v thru %v", i, j))

		if errInternal := ethClient.BatchCallContext(ctx, reqs[i:j]); errInternal != nil {
			logger.Errorw(fmt.Sprintf("Failed to batch send transactions %v thru %v", i, j),
				"error", errInternal)
			// Set this error on call
			for idx := i; idx < j; idx++ {
				reqs[idx].Error = errInternal
			}
			err = errors.Wrap(errInternal, "Failed to batch send transactions")
		}
	}
	return
}
