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

type finalizerHeadTracker interface {
	LatestAndFinalizedBlock(ctx context.Context) (latest, finalized *evmtypes.Head, err error)
}

// Finalizer handles processing new finalized blocks and marking transactions as finalized accordingly in the TXM DB
type evmFinalizer struct {
	services.StateMachine
	lggr         logger.SugaredLogger
	chainId      *big.Int
	rpcBatchSize int

	txStore     finalizerTxStore
	client      finalizerChainClient
	headTracker finalizerHeadTracker

	mb     *mailbox.Mailbox[*evmtypes.Head]
	stopCh services.StopChan
	wg     sync.WaitGroup
}

func NewEvmFinalizer(
	lggr logger.Logger,
	chainId *big.Int,
	rpcBatchSize uint32,
	txStore finalizerTxStore,
	client finalizerChainClient,
	headTracker finalizerHeadTracker,
) *evmFinalizer {
	lggr = logger.Named(lggr, "Finalizer")
	return &evmFinalizer{
		lggr:         logger.Sugared(lggr),
		chainId:      chainId,
		rpcBatchSize: int(rpcBatchSize),
		txStore:      txStore,
		client:       client,
		headTracker:  headTracker,
		mb:           mailbox.NewSingle[*evmtypes.Head](),
	}
}

// Start the finalizer
func (f *evmFinalizer) Start(ctx context.Context) error {
	return f.StartOnce("Finalizer", func() error {
		f.lggr.Debugf("started Finalizer with RPC batch size limit: %d", f.rpcBatchSize)
		f.stopCh = make(chan struct{})
		f.wg.Add(1)
		go f.runLoop()
		return nil
	})
}

// Close the finalizer
func (f *evmFinalizer) Close() error {
	return f.StopOnce("Finalizer", func() error {
		f.lggr.Debug("closing Finalizer")
		close(f.stopCh)
		f.wg.Wait()
		return nil
	})
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

func (f *evmFinalizer) DeliverLatestHead(head *evmtypes.Head) bool {
	return f.mb.Deliver(head)
}

func (f *evmFinalizer) ProcessHead(ctx context.Context, head *evmtypes.Head) error {
	ctx, cancel := context.WithTimeout(ctx, processHeadTimeout)
	defer cancel()
	_, latestFinalizedHead, err := f.headTracker.LatestAndFinalizedBlock(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve latest finalized head: %w", err)
	}
	return f.processFinalizedHead(ctx, latestFinalizedHead)
}

// Determines if any confirmed transactions can be marked as finalized by comparing their receipts against the latest finalized block
func (f *evmFinalizer) processFinalizedHead(ctx context.Context, latestFinalizedHead *evmtypes.Head) error {
	// Cannot determine finality without a finalized head for comparison
	if latestFinalizedHead == nil || !latestFinalizedHead.IsValid() {
		return fmt.Errorf("invalid latestFinalizedHead")
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
	validatedReceiptTxs, err := f.batchCheckReceiptHashesOnchain(ctx, receiptBlockHashToTx, latestFinalizedHead.BlockNumber())
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

func (f *evmFinalizer) batchCheckReceiptHashesOnchain(ctx context.Context, receiptMap map[common.Hash][]*Tx, latestFinalizedBlockNum int64) ([]*Tx, error) {
	if len(receiptMap) == 0 {
		return nil, nil
	}
	// Group the RPC batch calls in groups of rpcBatchSize
	var rpcBatchGroups [][]rpc.BatchElem
	var rpcBatch []rpc.BatchElem
	for hash := range receiptMap {
		elem := rpc.BatchElem{
			Method: "eth_getBlockByHash",
			Args: []any{
				hash,
				false,
			},
			Result: new(evmtypes.Head),
		}
		rpcBatch = append(rpcBatch, elem)
		if len(rpcBatch) >= f.rpcBatchSize {
			rpcBatchGroups = append(rpcBatchGroups, rpcBatch)
			rpcBatch = []rpc.BatchElem{}
		}
	}

	var finalizedTxs []*Tx
	for _, rpcBatch := range rpcBatchGroups {
		err := f.client.BatchCallContext(ctx, rpcBatch)
		if err != nil {
			// Continue if batch RPC call failed so other batches can still be considered for finalization
			f.lggr.Debugw("failed to find blocks due to batch call failure")
			continue
		}
		for _, req := range rpcBatch {
			if req.Error != nil {
				// Continue if particular RPC call failed so other txs can still be considered for finalization
				f.lggr.Debugw("failed to find block by hash", "hash", req.Args[0])
				continue
			}
			head := req.Result.(*evmtypes.Head)
			if head == nil {
				// Continue if particular RPC call yielded a nil head so other txs can still be considered for finalization
				f.lggr.Debugw("failed to find block by hash", "hash", req.Args[0])
				continue
			}
			// Confirmed receipt's block hash exists on-chain
			// Add to finalized list if block num less than or equal to the latest finalized head block num
			if head.BlockNumber() <= latestFinalizedBlockNum {
				txs := receiptMap[head.BlockHash()]
				finalizedTxs = append(finalizedTxs, txs...)
			}
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
