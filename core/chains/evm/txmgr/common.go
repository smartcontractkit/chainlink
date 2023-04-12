package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// timeout value for batchSendTransactions
const batchSendTransactionTimeout = 30 * time.Second

// Tries to send transactions in batches. Even if some batch(es) fail to get sent, it tries all remaining batches,
// before returning with error for the latest batch send. If a batch send fails, this sets the error on all
// elements in that batch.
func batchSendTransactions[ADDR types.Hashable[ADDR], TX_HASH types.Hashable[TX_HASH], BLOCK_HASH types.Hashable[BLOCK_HASH]](
	ctx context.Context,
	txStorageService txmgrtypes.TxStore[ADDR, big.Int, TX_HASH, BLOCK_HASH, NewTx[ADDR], *evmtypes.Receipt, EthTx[ADDR, TX_HASH], EthTxAttempt[ADDR, TX_HASH], int64, int64],
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

		if err := txStorageService.UpdateBroadcastAts(now, ethTxIDs[i:j]); err != nil {
			return reqs, errors.Wrap(err, "failed to update last succeeded on attempts")
		}
	}
	return reqs, nil
}

func getGethAddressFromADDR[ADDR commontypes.Hashable[ADDR]](addr ADDR) (common.Address, error) {
	addrHex, err := addr.MarshalText()
	if err != nil {
		return common.Address{}, errors.Wrapf(err, "failed to serialize address to text: %s", addr.String())
	}
	var gethAddr common.Address
	err = gethAddr.UnmarshalText(addrHex)
	if err != nil {
		return common.Address{}, errors.Wrapf(err, "failed to deserialize address from text: %s. Original address: %s", addrHex, addr.String())
	}
	return gethAddr, nil
}

func getGethHashFromHash[HASH commontypes.Hashable[HASH]](hash HASH) (common.Hash, error) {
	hashHex, err := hash.MarshalText()
	if err != nil {
		return common.Hash{}, errors.Wrapf(err, "failed to serialize hash to text: %s", hash.String())
	}
	var gethHash common.Hash
	err = gethHash.UnmarshalText(hashHex)
	if err != nil {
		return common.Hash{}, errors.Wrapf(err, "failed to deserialize hash from text: %s. Original hash: %s", hashHex, hash.String())
	}
	return gethHash, nil
}
