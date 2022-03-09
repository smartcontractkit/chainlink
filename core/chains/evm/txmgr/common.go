package txmgr

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// Tries to send transactions in batches. Even if some batch(es) fail to get sent, it tries all remaining batches,
// before returning with error for the latest batch send. If a batch send fails, this sets the error on all
// elements in that batch.
func batchSendTransactions(
	ctx context.Context,
	db *sqlx.DB,
	attempts []EthTxAttempt,
	batchSize int,
	logger logger.Logger,
	ethClient evmclient.Client) ([]rpc.BatchElem, error) {
	if len(attempts) == 0 {
		return nil, nil
	}

	reqs := make([]rpc.BatchElem, len(attempts))
	ethTxIDs := make([]int64, len(attempts))
	hashes := make([]common.Hash, len(attempts))
	for i, attempt := range attempts {
		ethTxIDs[i] = attempt.EthTxID
		hashes[i] = attempt.Hash
		req := rpc.BatchElem{
			Method: "eth_sendRawTransaction",
			Args:   []interface{}{hexutil.Encode(attempt.SignedRawTx)},
			Result: &common.Hash{},
		}
		reqs[i] = req
	}

	logger.Debugw(fmt.Sprintf("Batch sending %d unconfirmed transactions.", len(attempts)), "n", len(attempts), "ethTxIDs", ethTxIDs, "hashes", hashes)

	now := time.Now()
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
			return reqs, errors.Wrap(err, "failed to batch send transactions")
		}

		if err := updateBroadcastAts(db, now, ethTxIDs[i:j]); err != nil {
			return reqs, errors.Wrap(err, "failed to update last succeeded on attempts")
		}
	}
	return reqs, nil
}

func updateBroadcastAts(db *sqlx.DB, now time.Time, etxIDs []int64) error {
	// Deliberately do nothing on NULL broadcast_at because that indicates the
	// tx has been moved into a state where broadcast_at is not relevant, e.g.
	// fatally errored.
	//
	// Since EthConfirmer/EthResender can race (totally OK since highest
	// priced transaction always wins) we only want to update broadcast_at if
	// our version is later.
	_, err := db.Exec(`UPDATE eth_txes SET broadcast_at = $1 WHERE id = ANY($2) AND broadcast_at < $1`, now, pq.Array(etxIDs))
	return errors.Wrap(err, "updateBroadcastAts failed to update eth_txes")
}
