package txmgr

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// timeout value for batchSendTransactions
const batchSendTransactionTimeout = 30 * time.Second

// Tries to send transactions in batches. Even if some batch(es) fail to get sent, it tries all remaining batches,
// before returning with error for the latest batch send. If a batch send fails, this sets the error on all
// elements in that batch.
func batchSendTransactions[
	CHAIN_ID txmgrtypes.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R any,
	SEQ txmgrtypes.Sequence,
](
	ctx context.Context,
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, txmgrtypes.NewTx[ADDR, TX_HASH], R, EthTx[ADDR, TX_HASH], EthTxAttempt[ADDR, TX_HASH], SEQ],
	attempts []EthTxAttempt[ADDR, TX_HASH],
	batchSize int,
	logger logger.Logger,
	ethClient evmclient.Client) ([]rpc.BatchElem, error) {
	if len(attempts) == 0 {
		return nil, nil
	}

	reqs := make([]rpc.BatchElem, len(attempts))
	ethTxIDs := make([]int64, len(attempts))
	hashes := make([]string, len(attempts))
	for i, attempt := range attempts {
		ethTxIDs[i] = attempt.EthTxID
		hashes[i] = attempt.Hash.String()
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

		if err := txStore.UpdateBroadcastAts(now, ethTxIDs[i:j]); err != nil {
			return reqs, errors.Wrap(err, "failed to update last succeeded on attempts")
		}
	}
	return reqs, nil
}

func stringToGethAddress(s string) (common.Address, error) {
	if !common.IsHexAddress(s) {
		return common.Address{}, fmt.Errorf("invalid hex address: %s", s)
	}
	return common.HexToAddress(s), nil
}

func stringToGethHash(s string) (h common.Hash, err error) {
	err = h.UnmarshalText([]byte(s))
	return
}
