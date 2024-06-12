package txmgr

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type finalizerTxStore[CHAIN_ID types.ID, ADDR types.Hashable, TX_HASH types.Hashable, BLOCK_HASH types.Hashable, SEQ types.Sequence, FEE feetypes.Fee] interface {
	FindTransactionsByState(ctx context.Context, state txmgrtypes.TxState, chainID CHAIN_ID) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error)
	UpdateTxesFinalized(ctx context.Context, txs []int64, chainId CHAIN_ID) error
}

type finalizerChainClient[BLOCK_HASH types.Hashable, HEAD types.Head[BLOCK_HASH]] interface {
	HeadByHash(ctx context.Context, hash BLOCK_HASH) (HEAD, error)
}

// Finalizer handles processing new finalized blocks and marking transactions as finalized accordingly in the TXM DB
type Finalizer[CHAIN_ID types.ID, ADDR types.Hashable, TX_HASH types.Hashable, BLOCK_HASH types.Hashable, SEQ types.Sequence, FEE feetypes.Fee, HEAD types.Head[BLOCK_HASH]] struct {
	services.StateMachine
	lggr      logger.SugaredLogger
	chainId   CHAIN_ID
	txStore   finalizerTxStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	client    finalizerChainClient[BLOCK_HASH, HEAD]
	mb        *mailbox.Mailbox[HEAD]
	stopCh    services.StopChan
	wg        sync.WaitGroup
	initSync  sync.Mutex
	isStarted bool
}

func NewFinalizer[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
](
	lggr logger.Logger,
	chainId CHAIN_ID,
	txStore finalizerTxStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	client finalizerChainClient[BLOCK_HASH, HEAD],
) *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD] {
	lggr = logger.Named(lggr, "Finalizer")
	return &Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]{
		txStore: txStore,
		lggr:    logger.Sugared(lggr),
		chainId: chainId,
		client:  client,
		mb:      mailbox.NewSingle[HEAD](),
	}
}

// Start is a comment to appease the linter
func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) Start(ctx context.Context) error {
	return f.StartOnce("Finalizer", func() error {
		return f.startInternal(ctx)
	})
}

func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) startInternal(_ context.Context) error {
	f.initSync.Lock()
	defer f.initSync.Unlock()
	if f.isStarted {
		return errors.New("Finalizer is already started")
	}

	f.stopCh = make(chan struct{})
	f.wg = sync.WaitGroup{}
	f.wg.Add(1)
	go f.runLoop()
	f.isStarted = true
	return nil
}

// Close is a comment to appease the linter
func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) Close() error {
	return f.StopOnce("Finalizer", func() error {
		return f.closeInternal()
	})
}

func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) closeInternal() error {
	f.initSync.Lock()
	defer f.initSync.Unlock()
	if !f.isStarted {
		return fmt.Errorf("Finalizer is not started: %w", services.ErrAlreadyStopped)
	}
	close(f.stopCh)
	f.wg.Wait()
	f.isStarted = false
	return nil
}

func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) Name() string {
	return f.lggr.Name()
}

func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) HealthReport() map[string]error {
	return map[string]error{f.Name(): f.Healthy()}
}

func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) runLoop() {
	defer f.wg.Done()
	ctx, cancel := f.stopCh.NewCtx()
	defer cancel()
	for {
		select {
		case <-f.mb.Notify():
			for {
				if ctx.Err() != nil {
					return
				}
				head, exists := f.mb.Retrieve()
				if !exists {
					break
				}
				if err := f.ProcessHead(ctx, head); err != nil {
					f.lggr.Errorw("Error processing head", "err", err)
					f.SvcErrBuffer.Append(err)
					continue
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) ProcessHead(ctx context.Context, head types.Head[BLOCK_HASH]) error {
	ctx, cancel := context.WithTimeout(ctx, processHeadTimeout)
	defer cancel()
	return f.processHead(ctx, head)
}

// Determines if any confirmed transactions can be marked as finalized by comparing their receipts against the latest finalized block
func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) processHead(ctx context.Context, head types.Head[BLOCK_HASH]) error {
	latestFinalizedHead := head.LatestFinalizedHead()
	// Cannot determine finality without a finalized head for comparison
	if latestFinalizedHead == nil || !latestFinalizedHead.IsValid() {
		return fmt.Errorf("failed to find latest finalized head in chain")
	}
	earliestBlockNumInChain := latestFinalizedHead.EarliestHeadInChain().BlockNumber()
	f.lggr.Debugw("processing latest finalized head", "block num", latestFinalizedHead.BlockNumber(), "block hash", latestFinalizedHead.BlockHash(), "earliest block num in chain", earliestBlockNumInChain)

	// Retrieve all confirmed transactions, loaded with attempts and receipts
	confirmedTxs, err := f.txStore.FindTransactionsByState(ctx, TxConfirmed, f.chainId)
	if err != nil {
		return fmt.Errorf("failed to retrieve confirmed transactions: %w", err)
	}

	var finalizedTxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	// Group by block hash transactions whose receipts cannot be validated using the cached heads
	receiptBlockHashToTx := make(map[BLOCK_HASH][]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
	// Find transactions with receipt block nums older than the latest finalized block num and block hashes still in chain
	for _, tx := range confirmedTxs {
		// Only consider transactions not already marked as finalized
		if tx.Finalized {
			continue
		}
		receipt := tx.GetReceipt()
		if receipt == nil || receipt.IsZero() || receipt.IsUnmined() {
			continue
		}
		// Receipt newer than latest finalized head block num
		if receipt.GetBlockNumber().Cmp(big.NewInt(latestFinalizedHead.BlockNumber())) > 0 {
			continue
		}
		// Receipt block num older than earliest head in chain. Validate hash using RPC call later
		if receipt.GetBlockNumber().Int64() < earliestBlockNumInChain {
			receiptBlockHashToTx[receipt.GetBlockHash()] = append(receiptBlockHashToTx[receipt.GetBlockHash()], tx)
			continue
		}
		blockHashInChain := latestFinalizedHead.HashAtHeight(receipt.GetBlockNumber().Int64())
		// Receipt block hash does not match the block hash in chain. Transaction has been re-org'd out but DB state has not been updated yet
		if blockHashInChain.String() != receipt.GetBlockHash().String() {
			continue
		}
		finalizedTxs = append(finalizedTxs, tx)
	}

	// Check if block hashes exist for receipts on-chain older than the earliest cached head
	// Transactions are grouped by their receipt block hash to minimize the number of RPC calls in case transactions were confirmed in the same block
	// This check is only expected to be used in rare cases if there was an issue with the HeadTracker or if the node was down for significant time
	var wg sync.WaitGroup
	var txMu sync.RWMutex
	for receiptBlockHash, txs := range receiptBlockHashToTx {
		wg.Add(1)
		go func(hash BLOCK_HASH, txs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
			defer wg.Done()
			if head, rpcErr := f.client.HeadByHash(ctx, hash); rpcErr == nil && head.IsValid() {
				txMu.Lock()
				finalizedTxs = append(finalizedTxs, txs...)
				txMu.Unlock()
			}
		}(receiptBlockHash, txs)
	}
	wg.Wait()

	etxIDs := f.buildTxIdList(finalizedTxs)

	err = f.txStore.UpdateTxesFinalized(ctx, etxIDs, f.chainId)
	if err != nil {
		return fmt.Errorf("failed to update transactions as finalized: %w", err)
	}
	return nil
}

// Build list of transaction IDs
func (f *Finalizer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, HEAD]) buildTxIdList(finalizedTxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) []int64 {
	etxIDs := make([]int64, len(finalizedTxs))
	for i, tx := range finalizedTxs {
		receipt := tx.GetReceipt()
		f.lggr.Debugw("transaction considered finalized",
			"sequence", tx.Sequence,
			"fromAddress", tx.FromAddress.String(),
			"txHash", receipt.GetTxHash().String(),
			"receiptBlockNum", receipt.GetBlockNumber().Int64(),
			"receiptBlockHash", receipt.GetBlockHash().String(),
		)
		etxIDs[i] = tx.ID
	}
	return etxIDs
}
