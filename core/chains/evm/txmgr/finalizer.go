package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

var _ Finalizer = (*evmFinalizer)(nil)

// processHeadTimeout represents a sanity limit on how long ProcessHead should take to complete
const processHeadTimeout = 10 * time.Minute

type finalizerTxStore interface {
	FindConfirmedTxesReceipts(ctx context.Context, finalizedBlockNum int64, chainID *big.Int) ([]Receipt, error)
	UpdateTxStatesToFinalizedUsingReceiptIds(ctx context.Context, txs []int64, chainId *big.Int) error
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

	lastProcessedFinalizedBlockNum int64
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
	// Only continue processing if the latestFinalizedHead has not already been processed
	// Helps avoid unnecessary processing on every head if blocks are finalized in batches
	if latestFinalizedHead.BlockNumber() == f.lastProcessedFinalizedBlockNum {
		return nil
	}
	if latestFinalizedHead.BlockNumber() < f.lastProcessedFinalizedBlockNum {
		f.lggr.Errorw("Received finalized block older than one already processed. This should never happen and could be an issue with RPCs.", "lastProcessedFinalizedBlockNum", f.lastProcessedFinalizedBlockNum, "retrievedFinalizedBlockNum", latestFinalizedHead.BlockNumber())
		return nil
	}

	earliestBlockNumInChain := latestFinalizedHead.EarliestHeadInChain().BlockNumber()
	f.lggr.Debugw("processing latest finalized head", "blockNum", latestFinalizedHead.BlockNumber(), "blockHash", latestFinalizedHead.BlockHash(), "earliestBlockNumInChain", earliestBlockNumInChain)

	// Retrieve all confirmed transactions with receipts older than or equal to the finalized block, loaded with attempts and receipts
	unfinalizedReceipts, err := f.txStore.FindConfirmedTxesReceipts(ctx, latestFinalizedHead.BlockNumber(), f.chainId)
	if err != nil {
		return fmt.Errorf("failed to retrieve receipts for confirmed, unfinalized transactions: %w", err)
	}

	var finalizedReceipts []Receipt
	// Group by block hash transactions whose receipts cannot be validated using the cached heads
	blockNumToReceiptsMap := make(map[int64][]Receipt)
	// Find transactions with receipt block nums older than the latest finalized block num and block hashes still in chain
	for _, receipt := range unfinalizedReceipts {
		// The tx store query ensures transactions have receipts but leaving this check here for a belts and braces approach
		if receipt.TxHash == utils.EmptyHash || receipt.BlockHash == utils.EmptyHash {
			f.lggr.AssumptionViolationw("invalid receipt found for confirmed transaction", "receipt", receipt)
			continue
		}
		// The tx store query only returns transactions with receipts older than or equal to the finalized block but leaving this check here for a belts and braces approach
		if receipt.BlockNumber > latestFinalizedHead.BlockNumber() {
			continue
		}
		// Receipt block num older than earliest head in chain. Validate hash using RPC call later
		if receipt.BlockNumber < earliestBlockNumInChain {
			blockNumToReceiptsMap[receipt.BlockNumber] = append(blockNumToReceiptsMap[receipt.BlockNumber], receipt)
			continue
		}
		blockHashInChain := latestFinalizedHead.HashAtHeight(receipt.BlockNumber)
		// Receipt block hash does not match the block hash in chain. Transaction has been re-org'd out but DB state has not been updated yet
		if blockHashInChain.String() != receipt.BlockHash.String() {
			// Log error if a transaction is marked as confirmed with a receipt older than the finalized block
			// This scenario could potentially point to a re-org'd transaction the Confirmer has lost track of
			f.lggr.Errorw("found confirmed transaction with re-org'd receipt older than finalized block", "receipt", receipt, "onchainBlockHash", blockHashInChain.String())
			continue
		}
		finalizedReceipts = append(finalizedReceipts, receipt)
	}

	// Check if block hashes exist for receipts on-chain older than the earliest cached head
	// Transactions are grouped by their receipt block hash to avoid repeat requests on the same hash in case transactions were confirmed in the same block
	validatedReceipts := f.batchCheckReceiptHashesOnchain(ctx, blockNumToReceiptsMap)
	finalizedReceipts = append(finalizedReceipts, validatedReceipts...)

	receiptIDs := f.buildReceiptIdList(finalizedReceipts)

	err = f.txStore.UpdateTxStatesToFinalizedUsingReceiptIds(ctx, receiptIDs, f.chainId)
	if err != nil {
		return fmt.Errorf("failed to update transactions as finalized: %w", err)
	}
	// Update lastProcessedFinalizedBlockNum after processing has completed to allow failed processing to retry on subsequent heads
	// Does not need to be protected with mutex lock because the Finalizer only runs in a single loop
	f.lastProcessedFinalizedBlockNum = latestFinalizedHead.BlockNumber()
	return nil
}

func (f *evmFinalizer) batchCheckReceiptHashesOnchain(ctx context.Context, blockNumToReceiptsMap map[int64][]Receipt) []Receipt {
	if len(blockNumToReceiptsMap) == 0 {
		return nil
	}
	// Group the RPC batch calls in groups of rpcBatchSize
	var rpcBatchGroups [][]rpc.BatchElem
	var rpcBatch []rpc.BatchElem
	for blockNum := range blockNumToReceiptsMap {
		elem := rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args: []any{
				hexutil.EncodeBig(big.NewInt(blockNum)),
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
	if len(rpcBatch) > 0 {
		rpcBatchGroups = append(rpcBatchGroups, rpcBatch)
	}

	var finalizedReceipts []Receipt
	for _, rpcBatch := range rpcBatchGroups {
		err := f.client.BatchCallContext(ctx, rpcBatch)
		if err != nil {
			// Continue if batch RPC call failed so other batches can still be considered for finalization
			f.lggr.Errorw("failed to find blocks due to batch call failure", "error", err)
			continue
		}
		for _, req := range rpcBatch {
			if req.Error != nil {
				// Continue if particular RPC call failed so other txs can still be considered for finalization
				f.lggr.Errorw("failed to find block by number", "blockNum", req.Args[0], "error", req.Error)
				continue
			}
			head, ok := req.Result.(*evmtypes.Head)
			if !ok || !head.IsValid() {
				// Continue if particular RPC call yielded a nil block so other txs can still be considered for finalization
				f.lggr.Errorw("retrieved nil head for block number", "blockNum", req.Args[0])
				continue
			}
			receipts := blockNumToReceiptsMap[head.BlockNumber()]
			// Check if transaction receipts match the block hash at the given block num
			// If they do not, the transactions may have been re-org'd out
			// The expectation is for the Confirmer to pick up on these re-orgs and get the transaction included
			for _, receipt := range receipts {
				if receipt.BlockHash.String() == head.BlockHash().String() {
					finalizedReceipts = append(finalizedReceipts, receipt)
				} else {
					// Log error if a transaction is marked as confirmed with a receipt older than the finalized block
					// This scenario could potentially point to a re-org'd transaction the Confirmer has lost track of
					f.lggr.Errorw("found confirmed transaction with re-org'd receipt older than finalized block", "receipt", receipt, "onchainBlockHash", head.BlockHash().String())
				}
			}
		}
	}
	return finalizedReceipts
}

// Build list of transaction IDs
func (f *evmFinalizer) buildReceiptIdList(finalizedReceipts []Receipt) []int64 {
	receiptIds := make([]int64, len(finalizedReceipts))
	for i, receipt := range finalizedReceipts {
		f.lggr.Debugw("transaction considered finalized",
			"txHash", receipt.TxHash.String(),
			"receiptBlockNum", receipt.BlockNumber,
			"receiptBlockHash", receipt.BlockHash.String(),
		)
		receiptIds[i] = receipt.ID
	}
	return receiptIds
}
