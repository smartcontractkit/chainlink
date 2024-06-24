package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

var _ Finalizer = (*evmFinalizer)(nil)

// processHeadTimeout represents a sanity limit on how long ProcessHead
// should take to complete
const processHeadTimeout = 10 * time.Minute

type finalizerTxStore interface {
	FindConfirmedTxesAwaitingFinalization(ctx context.Context, chainID *big.Int) ([]*Tx, error)
	UpdateTxesFinalized(ctx context.Context, txs []int64, chainId *big.Int) error
}

type finalizerChainClient interface {
	BatchCallContext(ctx context.Context, elems []rpc.BatchElem) error
}

// Finalizer handles processing new finalized blocks and marking transactions as finalized accordingly in the TXM DB
type evmFinalizer struct {
	services.StateMachine
	lggr    logger.SugaredLogger
	chainId *big.Int
	txStore finalizerTxStore
	client  finalizerChainClient
	mb      *mailbox.Mailbox[*evmtypes.Head]
	stopCh  services.StopChan
	wg      sync.WaitGroup
}

func NewEvmFinalizer(
	lggr logger.Logger,
	chainId *big.Int,
	txStore finalizerTxStore,
	client finalizerChainClient,
) *evmFinalizer {
	lggr = logger.Named(lggr, "Finalizer")
	return &evmFinalizer{
		txStore: txStore,
		lggr:    logger.Sugared(lggr),
		chainId: chainId,
		client:  client,
		mb:      mailbox.NewSingle[*evmtypes.Head](),
	}
}

// Start is a comment to appease the linter
func (f *evmFinalizer) Start(ctx context.Context) error {
	return f.StartOnce("Finalizer", func() error {
		return f.startInternal(ctx)
	})
}

func (f *evmFinalizer) startInternal(_ context.Context) error {
	f.stopCh = make(chan struct{})
	f.wg = sync.WaitGroup{}
	f.wg.Add(1)
	go f.runLoop()
	return nil
}

// Close is a comment to appease the linter
func (f *evmFinalizer) Close() error {
	return f.StopOnce("Finalizer", func() error {
		return f.closeInternal()
	})
}

func (f *evmFinalizer) closeInternal() error {
	close(f.stopCh)
	f.wg.Wait()
	return nil
}

func (f *evmFinalizer) Name() string {
	return f.lggr.Name()
}

func (f *evmFinalizer) HealthReport() map[string]error {
	return map[string]error{f.Name(): f.Healthy()}
}

func (f *evmFinalizer) runLoop() {
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

func (f *evmFinalizer) DeliverHead(head *evmtypes.Head) bool {
	return f.mb.Deliver(head)
}

func (f *evmFinalizer) ProcessHead(ctx context.Context, head *evmtypes.Head) error {
	ctx, cancel := context.WithTimeout(ctx, processHeadTimeout)
	defer cancel()
	return f.processHead(ctx, head)
}

// Determines if any confirmed transactions can be marked as finalized by comparing their receipts against the latest finalized block
func (f *evmFinalizer) processHead(ctx context.Context, head *evmtypes.Head) error {
	latestFinalizedHead := head.LatestFinalizedHead()
	// Cannot determine finality without a finalized head for comparison
	if latestFinalizedHead == nil || !latestFinalizedHead.IsValid() {
		return fmt.Errorf("failed to find latest finalized head in chain")
	}
	earliestBlockNumInChain := latestFinalizedHead.EarliestHeadInChain().BlockNumber()
	f.lggr.Debugw("processing latest finalized head", "block num", latestFinalizedHead.BlockNumber(), "block hash", latestFinalizedHead.BlockHash(), "earliest block num in chain", earliestBlockNumInChain)

	// Retrieve all confirmed transactions, loaded with attempts and receipts
	unfinalizedTxs, err := f.txStore.FindConfirmedTxesAwaitingFinalization(ctx, f.chainId)
	if err != nil {
		return fmt.Errorf("failed to retrieve confirmed transactions: %w", err)
	}

	var finalizedTxs []*Tx
	// Group by block hash transactions whose receipts cannot be validated using the cached heads
	receiptBlockHashToTx := make(map[common.Hash][]*Tx)
	// Find transactions with receipt block nums older than the latest finalized block num and block hashes still in chain
	for _, tx := range unfinalizedTxs {
		// Only consider transactions not already marked as finalized
		if tx.Finalized {
			continue
		}
		receipt := tx.GetReceipt()
		if receipt == nil || receipt.IsZero() || receipt.IsUnmined() {
			f.lggr.AssumptionViolationw("invalid receipt found for confirmed transaction", "tx", tx, "receipt", receipt)
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
	// Transactions are grouped by their receipt block hash to avoid repeat requests on the same hash in case transactions were confirmed in the same block
	validatedReceiptTxs, err := f.batchCheckReceiptHashes(ctx, receiptBlockHashToTx, latestFinalizedHead.BlockNumber())
	if err != nil {
		// Do not error out to allow transactions that did not need RPC validation to still be marked as finalized
		// The transactions failed to be validated will be checked again in the next round
		f.lggr.Errorf("failed to validate receipt block hashes over RPC: %v", err)
	}
	finalizedTxs = append(finalizedTxs, validatedReceiptTxs...)

	etxIDs := f.buildTxIdList(finalizedTxs)

	err = f.txStore.UpdateTxesFinalized(ctx, etxIDs, f.chainId)
	if err != nil {
		return fmt.Errorf("failed to update transactions as finalized: %w", err)
	}
	return nil
}

func (f *evmFinalizer) batchCheckReceiptHashes(ctx context.Context, receiptMap map[common.Hash][]*Tx, latestFinalizedBlockNum int64) ([]*Tx, error) {
	if len(receiptMap) == 0 {
		return nil, nil
	}
	var rpcBatchCalls []rpc.BatchElem
	for hash := range receiptMap {
		elem := rpc.BatchElem{
			Method: "eth_getBlockByHash",
			Args: []any{
				hash,
				false,
			},
			Result: new(evmtypes.Head),
		}
		rpcBatchCalls = append(rpcBatchCalls, elem)
	}

	err := f.client.BatchCallContext(ctx, rpcBatchCalls)
	if err != nil {
		return nil, fmt.Errorf("get block hash batch call failed: %w", err)
	}
	var finalizedTxs []*Tx
	for _, req := range rpcBatchCalls {
		if req.Error != nil {
			f.lggr.Debugw("failed to find block by hash", "hash", req.Args[0])
			continue
		}
		head := req.Result.(*evmtypes.Head)
		if head == nil {
			f.lggr.Debugw("failed to find block by hash", "hash", req.Args[0])
			continue
		}
		// Confirmed receipt's block hash exists on-chain still
		// Add to finalized list if block num less than or equal to the latest finalized head block num
		if head.BlockNumber() <= latestFinalizedBlockNum {
			txs := receiptMap[head.BlockHash()]
			finalizedTxs = append(finalizedTxs, txs...)
		}
	}
	return finalizedTxs, nil
}

// Build list of transaction IDs
func (f *evmFinalizer) buildTxIdList(finalizedTxs []*Tx) []int64 {
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
