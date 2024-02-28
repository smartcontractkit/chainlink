package txmgr

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

// Tries to send transactions in batches. Even if some batch(es) fail to get sent, it tries all remaining batches,
// before returning with error for the latest batch send. If a batch send fails, this sets the error on all
// elements in that batch.
func batchSendTransactions(
	ctx context.Context,
	attempts []TxAttempt,
	batchSize int,
	logger logger.Logger,
	ethClient evmclient.Client,
) (
	[]rpc.BatchElem,
	time.Time, // batch broadcast time
	[]int64, // successfully broadcast tx IDs
	error) {
	if len(attempts) == 0 {
		return nil, time.Now(), nil, nil
	}

	reqs := make([]rpc.BatchElem, len(attempts))
	ethTxIDs := make([]int64, len(attempts))
	hashes := make([]string, len(attempts))
	now := time.Now()
	successfulBroadcast := []int64{}
	for i, attempt := range attempts {
		ethTxIDs[i] = attempt.TxID
		hashes[i] = attempt.Hash.String()
		// Decode the signed raw tx back into a Transaction object
		signedTx, decodeErr := GetGethSignedTx(attempt.SignedRawTx)
		if decodeErr != nil {
			return reqs, now, successfulBroadcast, fmt.Errorf("failed to decode signed raw tx into Transaction object: %w", decodeErr)
		}
		// Get the canonical encoding of the Transaction object needed for the eth_sendRawTransaction request
		// The signed raw tx cannot be used directly because it uses a different encoding
		txBytes, marshalErr := signedTx.MarshalBinary()
		if marshalErr != nil {
			return reqs, now, successfulBroadcast, fmt.Errorf("failed to marshal tx into canonical encoding: %w", marshalErr)
		}
		req := rpc.BatchElem{
			Method: "eth_sendRawTransaction",
			Args:   []interface{}{hexutil.Encode(txBytes)},
			Result: &common.Hash{},
		}
		reqs[i] = req
	}

	logger.Debugw(fmt.Sprintf("Batch sending %d unconfirmed transactions.", len(attempts)), "n", len(attempts), "ethTxIDs", ethTxIDs, "hashes", hashes)

	if batchSize == 0 {
		batchSize = len(reqs)
	}

	for i := 0; i < len(reqs); i += batchSize {
		j := i + batchSize
		if j > len(reqs) {
			j = len(reqs)
		}

		logger.Debugw(fmt.Sprintf("Batch sending transactions %v thru %v", i, j))

		if err := ethClient.BatchCallContextAll(ctx, reqs[i:j]); err != nil {
			return reqs, now, successfulBroadcast, pkgerrors.Wrap(err, "failed to batch send transactions")
		}
		successfulBroadcast = append(successfulBroadcast, ethTxIDs[i:j]...)
	}
	return reqs, now, successfulBroadcast, nil
}
