package services

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type BlockWithReceipts struct {
	block    *types.Block
	receipts []bulletprooftxmanager.Receipt
}

type BlockFetcher struct {
	store  *strpkg.Store
	logger *logger.Logger
}

func NewBlockFetcher(store *strpkg.Store) *BlockFetcher {
	return &BlockFetcher{
		store: store,
		logger: logger.CreateLogger(logger.Default.With(
			"id", "block_fetcher",
		)),
	}
}

func (ht *BlockFetcher) fetchBlock(ctx context.Context, head models.Head) (*BlockWithReceipts, error) {

	start := time.Now()
	block, err := ht.store.EthClient.FastBlockByHash(ctx, head.Hash)
	if err != nil {
		return nil, errors.Wrap(err, "HeadTracker#handleNewHighestHead failed fetching BlockByHash")
	}
	elapsed := time.Since(start)

	ht.logger.Debugw("========================= HeadTracker: getting whole block took", "elapsedMs", elapsed.Milliseconds())

	logger.Warnf("rpcBlock.Hash: %v", block.Hash())
	logger.Warnf("rpcBlock.num Transactions: %v", block.Transactions().Len())

	start2 := time.Now()

	receipts, err := ht.batchFetchReceipts(ctx, block.Transactions())

	receiptsSizeBytes := 0
	for _, receipt := range receipts {
		for _, log := range receipt.Logs {
			receiptsSizeBytes += len(log.Data)
			receiptsSizeBytes += (len(log.Topics) + 2) * common.HashLength
			receiptsSizeBytes += common.AddressLength
		}
	}

	elapsed2 := time.Since(start2)
	ht.logger.Debugw("========================= HeadTracker: getting tx receipts took",
		"elapsedMs", elapsed2.Milliseconds(), "receiptsSizeBytes", receiptsSizeBytes)

	if err != nil {
		return nil, errors.Wrap(err, "HeadTracker#batchFetchReceipts failed ")
	}

	return &BlockWithReceipts{
		block,
		receipts,
	}, nil
}

func (ht *BlockFetcher) batchFetchReceipts(ctx context.Context, txs []*types.Transaction) (receipts []bulletprooftxmanager.Receipt, err error) {
	var reqs []rpc.BatchElem
	for _, tx := range txs { // TODO: how many is too many?
		req := rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{tx.Hash()},
			Result: &bulletprooftxmanager.Receipt{},
		}
		reqs = append(reqs, req)
	}

	ctx, cancel := eth.DefaultQueryCtx(ctx)
	defer cancel()

	err = ht.store.EthClient.BatchCallContext(ctx, reqs)
	if err != nil {
		return nil, errors.Wrap(err, "EthConfirmer#batchFetchReceipts error fetching receipts with BatchCallContext")
	}

	for i, req := range reqs {
		tx := txs[i]
		result, err := req.Result, req.Error

		receipt, is := result.(*bulletprooftxmanager.Receipt)
		if !is {
			return nil, errors.Errorf("expected result to be a %T, got %T", (*bulletprooftxmanager.Receipt)(nil), result)
		}

		l := ht.logger.With(
			"txHash", tx.Hash().Hex(), "err", err,
		)

		if err != nil {
			l.Errorw("EthConfirmer#batchFetchReceipts: fetchReceipt failed")
			continue
		}

		if receipt == nil {
			// NOTE: This should never possibly happen, but it seems safer to
			// check regardless to avoid a potential panic
			l.Errorw("EthConfirmer#batchFetchReceipts: invariant violation, got nil receipt")
			continue
		}

		if receipt.IsZero() {
			l.Debugw("EthConfirmer#batchFetchReceipts: still waiting for receipt")
			continue
		}

		l = l.With("receipt", receipt)

		if receipt.IsUnmined() {
			l.Debugw("EthConfirmer#batchFetchReceipts: got receipt for transaction but it's still in the mempool and not included in a block yet")
			continue
		}

		//l.Debugw("EthConfirmer#batchFetchReceipts: got receipt for transaction", "blockNumber", receipt.BlockNumber)

		if receipt.TxHash != tx.Hash() {
			l.Errorf("EthConfirmer#batchFetchReceipts: invariant violation, expected receipt with hash %s to have same hash as attempt with hash %s", receipt.TxHash.Hex(), tx.Hash().Hex())
			continue
		}

		if receipt.BlockNumber == nil {
			l.Error("EthConfirmer#batchFetchReceipts: invariant violation, receipt was missing block number")
			continue
		}

		if receipt.Status == 0 {
			l.Warnf("transaction %s reverted on-chain", receipt.TxHash)
			// This is safe to increment here because we save the receipt immediately after
			// and once its saved we do not fetch it again.
			//promRevertedTxCount.Add(1)
		}

		receipts = append(receipts, *receipt)
	}

	return
}
