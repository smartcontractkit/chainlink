package txmgr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

var _ Finalizer = (*evmFinalizer)(nil)

var (
	promNumSuccessfulTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_successful_transactions",
		Help: "Total number of successful transactions. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
	}, []string{"chainID"})
	promRevertedTxCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_tx_reverted",
		Help: "Number of times a transaction reverted on-chain. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
	}, []string{"chainID"})
	promFwdTxCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_fwd_tx_count",
		Help: "The number of forwarded transaction attempts labeled by status",
	}, []string{"chainID", "successful"})
	promTxAttemptCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tx_manager_tx_attempt_count",
		Help: "The number of transaction attempts that are currently being processed by the transaction manager",
	}, []string{"chainID"})
	promNumFinalizedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_finalized_transactions",
		Help: "Total number of finalized transactions",
	}, []string{"chainID"})
)

var (
	// ErrCouldNotGetReceipt is the error string we save if we reach our LatestFinalizedBlockNum for a confirmed transaction
	// without ever getting a receipt. This most likely happened because an external wallet used the account for this nonce
	ErrCouldNotGetReceipt = "could not get receipt"
)

// processHeadTimeout represents a sanity limit on how long ProcessHead should take to complete
const processHeadTimeout = 10 * time.Minute

type finalizerTxStore interface {
	FindAttemptsRequiringReceiptFetch(ctx context.Context, chainID *big.Int) (hashes []TxAttempt, err error)
	FindConfirmedTxesReceipts(ctx context.Context, finalizedBlockNum int64, chainID *big.Int) (receipts []*evmtypes.Receipt, err error)
	FindTxesPendingCallback(ctx context.Context, latest, finalized int64, chainID *big.Int) (receiptsPlus []ReceiptPlus, err error)
	FindTxesByIDs(ctx context.Context, etxIDs []int64, chainID *big.Int) (etxs []*Tx, err error)
	PreloadTxes(ctx context.Context, attempts []TxAttempt) error
	SaveFetchedReceipts(ctx context.Context, r []*evmtypes.Receipt) (err error)
	UpdateTxCallbackCompleted(ctx context.Context, pipelineTaskRunID uuid.UUID, chainID *big.Int) error
	UpdateTxFatalErrorAndDeleteAttempts(ctx context.Context, etx *Tx) error
	UpdateTxStatesToFinalizedUsingTxHashes(ctx context.Context, txHashes []common.Hash, chainID *big.Int) error
}

type finalizerChainClient interface {
	BatchCallContext(ctx context.Context, elems []rpc.BatchElem) error
	BatchGetReceipts(ctx context.Context, attempts []TxAttempt) (txReceipt []*evmtypes.Receipt, txErr []error, funcErr error)
	CallContract(ctx context.Context, a TxAttempt, blockNumber *big.Int) (rpcErr fmt.Stringer, extractErr error)
}

type finalizerHeadTracker interface {
	LatestAndFinalizedBlock(ctx context.Context) (latest, finalized *evmtypes.Head, err error)
}

type resumeCallback = func(context.Context, uuid.UUID, interface{}, error) error

// Finalizer handles processing new finalized blocks and marking transactions as finalized accordingly in the TXM DB
type evmFinalizer struct {
	services.StateMachine
	lggr              logger.SugaredLogger
	chainID           *big.Int
	rpcBatchSize      int
	forwardersEnabled bool

	txStore     finalizerTxStore
	client      finalizerChainClient
	headTracker finalizerHeadTracker

	mb     *mailbox.Mailbox[*evmtypes.Head]
	stopCh services.StopChan
	wg     sync.WaitGroup

	lastProcessedFinalizedBlockNum int64
	resumeCallback                 resumeCallback
}

func NewEvmFinalizer(
	lggr logger.Logger,
	chainID *big.Int,
	rpcBatchSize uint32,
	forwardersEnabled bool,
	txStore finalizerTxStore,
	client finalizerChainClient,
	headTracker finalizerHeadTracker,
) *evmFinalizer {
	lggr = logger.Named(lggr, "Finalizer")
	return &evmFinalizer{
		lggr:              logger.Sugared(lggr),
		chainID:           chainID,
		rpcBatchSize:      int(rpcBatchSize),
		forwardersEnabled: forwardersEnabled,
		txStore:           txStore,
		client:            client,
		headTracker:       headTracker,
		mb:                mailbox.NewSingle[*evmtypes.Head](),
		resumeCallback:    nil,
	}
}

func (f *evmFinalizer) SetResumeCallback(callback resumeCallback) {
	f.resumeCallback = callback
}

// Start the finalizer
func (f *evmFinalizer) Start(ctx context.Context) error {
	return f.StartOnce("Finalizer", func() error {
		f.lggr.Debugw("started Finalizer", "rpcBatchSize", f.rpcBatchSize, "forwardersEnabled", f.forwardersEnabled)
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
	// Fetch and store receipts for confirmed transactions that do not have locally stored receipts
	err = f.FetchAndStoreReceipts(ctx, head, latestFinalizedHead)
	// Do not return on error since other functions are not dependent on results
	if err != nil {
		f.lggr.Errorf("failed to fetch and store receipts for confirmed transactions: %s", err.Error())
	}
	// Resume pending task runs if any receipts match the min confirmation criteria
	err = f.ResumePendingTaskRuns(ctx, head.BlockNumber(), latestFinalizedHead.BlockNumber())
	// Do not return on error since other functions are not dependent on results
	if err != nil {
		f.lggr.Errorf("failed to resume pending task runs: %s", err.Error())
	}
	return f.processFinalizedHead(ctx, latestFinalizedHead)
}

// processFinalizedHead determines if any confirmed transactions can be marked as finalized by comparing their receipts against the latest finalized block
// Fetches receipts directly from on-chain so re-org detection is not needed during finalization
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

	mark := time.Now()
	// Retrieve all confirmed transactions with receipts older than or equal to the finalized block
	unfinalizedReceipts, err := f.txStore.FindConfirmedTxesReceipts(ctx, latestFinalizedHead.BlockNumber(), f.chainID)
	if err != nil {
		return fmt.Errorf("failed to retrieve receipts for confirmed, unfinalized transactions: %w", err)
	}
	f.lggr.Debugw(fmt.Sprintf("found %d receipts for potential finalized transactions", len(unfinalizedReceipts)), "timeElapsed", time.Since(mark))
	mark = time.Now()

	finalizedReceipts := make([]*evmtypes.Receipt, 0, len(unfinalizedReceipts))
	// Group by block hash transactions whose receipts cannot be validated using the cached heads
	blockNumToReceiptsMap := make(map[int64][]*evmtypes.Receipt)
	// Find receipts with block nums older than or equal to the latest finalized block num
	for _, receipt := range unfinalizedReceipts {
		// The tx store query ensures transactions have receipts but leaving this check here for a belts and braces approach
		if receipt.TxHash == utils.EmptyHash || receipt.BlockHash == utils.EmptyHash {
			f.lggr.AssumptionViolationw("invalid receipt found for confirmed transaction", "receipt", receipt)
			continue
		}
		// The tx store query only returns transactions with receipts older than or equal to the finalized block but leaving this check here for a belts and braces approach
		if receipt.BlockNumber.Int64() > latestFinalizedHead.BlockNumber() {
			continue
		}
		// Receipt block num older than earliest head in chain. Validate hash using RPC call later
		if receipt.BlockNumber.Int64() < earliestBlockNumInChain {
			blockNumToReceiptsMap[receipt.BlockNumber.Int64()] = append(blockNumToReceiptsMap[receipt.BlockNumber.Int64()], receipt)
			continue
		}
		blockHashInChain := latestFinalizedHead.HashAtHeight(receipt.BlockNumber.Int64())
		// Receipt block hash does not match the block hash in chain. Transaction has been re-org'd out but DB state has not been updated yet
		if blockHashInChain.String() != receipt.BlockHash.String() {
			// Log error if a transaction is marked as confirmed with a receipt older than the finalized block
			// This scenario could potentially point to a re-org'd transaction the Confirmer has lost track of
			f.lggr.Errorw("found confirmed transaction with re-org'd receipt older than finalized block", "receipt", receipt, "onchainBlockHash", blockHashInChain.String())
			continue
		}
		finalizedReceipts = append(finalizedReceipts, receipt)
	}
	f.lggr.Debugw(fmt.Sprintf("found %d finalized transactions using local block history", len(finalizedReceipts)), "latestFinalizedBlockNum", latestFinalizedHead.BlockNumber(), "timeElapsed", time.Since(mark))
	mark = time.Now()

	// Check if block hashes exist for receipts on-chain older than the earliest cached head
	// Transactions are grouped by their receipt block hash to avoid repeat requests on the same hash in case transactions were confirmed in the same block
	validatedReceipts := f.batchCheckReceiptHashesOnchain(ctx, blockNumToReceiptsMap)
	finalizedReceipts = append(finalizedReceipts, validatedReceipts...)
	f.lggr.Debugw(fmt.Sprintf("found %d finalized transactions validated against RPC", len(validatedReceipts)), "latestFinalizedBlockNum", latestFinalizedHead.BlockNumber(), "timeElapsed", time.Since(mark))

	txHashes := f.buildTxHashList(finalizedReceipts)

	err = f.txStore.UpdateTxStatesToFinalizedUsingTxHashes(ctx, txHashes, f.chainID)
	if err != nil {
		return fmt.Errorf("failed to update transactions as finalized: %w", err)
	}

	// add newly finalized transactions to the prom metric
	promNumFinalizedTxs.WithLabelValues(f.chainID.String()).Add(float64(len(txHashes)))

	return nil
}

func (f *evmFinalizer) batchCheckReceiptHashesOnchain(ctx context.Context, blockNumToReceiptsMap map[int64][]*evmtypes.Receipt) []*evmtypes.Receipt {
	if len(blockNumToReceiptsMap) == 0 {
		return nil
	}
	// Group the RPC batch calls in groups of rpcBatchSize
	rpcBatchGroups := make([][]rpc.BatchElem, 0, len(blockNumToReceiptsMap))
	rpcBatch := make([]rpc.BatchElem, 0, f.rpcBatchSize)
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
			rpcBatch = make([]rpc.BatchElem, 0, f.rpcBatchSize)
		}
	}
	if len(rpcBatch) > 0 {
		rpcBatchGroups = append(rpcBatchGroups, rpcBatch)
	}

	finalizedReceipts := make([]*evmtypes.Receipt, 0, len(blockNumToReceiptsMap))
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

func (f *evmFinalizer) FetchAndStoreReceipts(ctx context.Context, head, latestFinalizedHead *evmtypes.Head) error {
	attempts, err := f.txStore.FindAttemptsRequiringReceiptFetch(ctx, f.chainID)
	if err != nil {
		return fmt.Errorf("failed to fetch broadcasted attempts for confirmed transactions: %w", err)
	}
	if len(attempts) == 0 {
		return nil
	}
	promTxAttemptCount.WithLabelValues(f.chainID.String()).Set(float64(len(attempts)))

	f.lggr.Debugw(fmt.Sprintf("Fetching receipts for %v transaction attempts", len(attempts)))

	batchSize := f.rpcBatchSize
	if batchSize == 0 {
		batchSize = len(attempts)
	}
	allReceipts := make([]*evmtypes.Receipt, 0, len(attempts))
	errorList := make([]error, 0, len(attempts))
	for i := 0; i < len(attempts); i += batchSize {
		j := i + batchSize
		if j > len(attempts) {
			j = len(attempts)
		}
		batch := attempts[i:j]

		receipts, fetchErr := f.batchFetchReceipts(ctx, batch)
		if fetchErr != nil {
			errorList = append(errorList, fetchErr)
			continue
		}

		allReceipts = append(allReceipts, receipts...)

		if err = f.txStore.SaveFetchedReceipts(ctx, receipts); err != nil {
			errorList = append(errorList, err)
			continue
		}
	}
	if len(errorList) > 0 {
		return errors.Join(errorList...)
	}

	oldTxIDs := findOldTxIDsWithoutReceipts(attempts, allReceipts, latestFinalizedHead)
	// Process old transactions that never received receipts and need to be marked as fatal
	err = f.ProcessOldTxsWithoutReceipts(ctx, oldTxIDs, head, latestFinalizedHead)
	if err != nil {
		return err
	}

	return nil
}

func (f *evmFinalizer) batchFetchReceipts(ctx context.Context, attempts []TxAttempt) (receipts []*evmtypes.Receipt, err error) {
	// Metadata is required to determine whether a tx is forwarded or not.
	if f.forwardersEnabled {
		err = f.txStore.PreloadTxes(ctx, attempts)
		if err != nil {
			return nil, fmt.Errorf("Confirmer#batchFetchReceipts error loading txs for attempts: %w", err)
		}
	}

	txReceipts, txErrs, err := f.client.BatchGetReceipts(ctx, attempts)
	if err != nil {
		return nil, err
	}

	for i, receipt := range txReceipts {
		attempt := attempts[i]
		err := txErrs[i]
		if err != nil {
			f.lggr.Error("FetchReceipts failed")
			continue
		}
		ok := f.validateReceipt(ctx, receipt, attempt)
		if !ok {
			continue
		}
		receipts = append(receipts, receipt)
	}

	return
}

// Note this function will increment promRevertedTxCount upon receiving a reverted transaction receipt
func (f *evmFinalizer) validateReceipt(ctx context.Context, receipt *evmtypes.Receipt, attempt TxAttempt) bool {
	l := attempt.Tx.GetLogger(f.lggr).With("txHash", attempt.Hash.String(), "txAttemptID", attempt.ID,
		"txID", attempt.TxID, "nonce", attempt.Tx.Sequence,
	)

	if receipt == nil {
		// NOTE: This should never happen, but it seems safer to check
		// regardless to avoid a potential panic
		l.AssumptionViolation("got nil receipt")
		return false
	}

	if receipt.IsZero() {
		l.Debug("Still waiting for receipt")
		return false
	}

	l = l.With("blockHash", receipt.GetBlockHash().String(), "status", receipt.GetStatus(), "transactionIndex", receipt.GetTransactionIndex())

	if receipt.IsUnmined() {
		l.Debug("Got receipt for transaction but it's still in the mempool and not included in a block yet")
		return false
	}

	l.Debugw("Got receipt for transaction", "blockNumber", receipt.GetBlockNumber(), "feeUsed", receipt.GetFeeUsed())

	if receipt.GetTxHash().String() != attempt.Hash.String() {
		l.Errorf("Invariant violation, expected receipt with hash %s to have same hash as attempt with hash %s", receipt.GetTxHash().String(), attempt.Hash.String())
		return false
	}

	if receipt.GetBlockNumber() == nil {
		l.Error("Invariant violation, receipt was missing block number")
		return false
	}

	if receipt.GetStatus() == 0 {
		if receipt.GetRevertReason() != nil {
			l.Warnw("transaction reverted on-chain", "hash", receipt.GetTxHash(), "revertReason", *receipt.GetRevertReason())
		} else {
			rpcError, errExtract := f.client.CallContract(ctx, attempt, receipt.GetBlockNumber())
			if errExtract == nil {
				l.Warnw("transaction reverted on-chain", "hash", receipt.GetTxHash(), "rpcError", rpcError.String())
			} else {
				l.Warnw("transaction reverted on-chain unable to extract revert reason", "hash", receipt.GetTxHash(), "err", errExtract)
			}
		}
		// This might increment more than once e.g. in case of re-orgs going back and forth we might re-fetch the same receipt
		promRevertedTxCount.WithLabelValues(f.chainID.String()).Add(1)
	} else {
		promNumSuccessfulTxs.WithLabelValues(f.chainID.String()).Add(1)
	}

	// This is only recording forwarded tx that were mined and have a status.
	// Counters are prone to being inaccurate due to re-orgs.
	if f.forwardersEnabled {
		meta, metaErr := attempt.Tx.GetMeta()
		if metaErr == nil && meta != nil && meta.FwdrDestAddress != nil {
			// promFwdTxCount takes two labels, chainID and a boolean of whether a tx was successful or not.
			promFwdTxCount.WithLabelValues(f.chainID.String(), strconv.FormatBool(receipt.GetStatus() != 0)).Add(1)
		}
	}
	return true
}

// ResumePendingTaskRuns issues callbacks to task runs that are pending waiting for receipts
func (f *evmFinalizer) ResumePendingTaskRuns(ctx context.Context, latest, finalized int64) error {
	if f.resumeCallback == nil {
		return nil
	}
	receiptsPlus, err := f.txStore.FindTxesPendingCallback(ctx, latest, finalized, f.chainID)

	if err != nil {
		return err
	}

	if len(receiptsPlus) > 0 {
		f.lggr.Debugf("Resuming %d task runs pending receipt", len(receiptsPlus))
	} else {
		f.lggr.Debug("No task runs to resume")
	}
	for _, data := range receiptsPlus {
		var taskErr error
		var output interface{}
		if data.FailOnRevert && data.Receipt.GetStatus() == 0 {
			taskErr = fmt.Errorf("transaction %s reverted on-chain", data.Receipt.GetTxHash())
		} else {
			output = data.Receipt
		}

		f.lggr.Debugw("Callback: resuming tx with receipt", "output", output, "taskErr", taskErr, "pipelineTaskRunID", data.ID)
		if err := f.resumeCallback(ctx, data.ID, output, taskErr); err != nil {
			return fmt.Errorf("failed to resume suspended pipeline run: %w", err)
		}
		// Mark tx as having completed callback
		if err := f.txStore.UpdateTxCallbackCompleted(ctx, data.ID, f.chainID); err != nil {
			return err
		}
	}

	return nil
}

func (f *evmFinalizer) ProcessOldTxsWithoutReceipts(ctx context.Context, oldTxIDs []int64, head, latestFinalizedHead *evmtypes.Head) error {
	if len(oldTxIDs) == 0 {
		return nil
	}
	oldTxs, err := f.txStore.FindTxesByIDs(ctx, oldTxIDs, f.chainID)
	if err != nil {
		return fmt.Errorf("failed to find transactions with IDs: %w", err)
	}

	errorList := make([]error, 0, len(oldTxs))
	for _, oldTx := range oldTxs {
		f.lggr.Criticalw(fmt.Sprintf("transaction with ID %v expired without ever getting a receipt for any of our attempts. "+
			"Current block height is %d, transaction was broadcast before finalized block %d. This transaction may not have not been sent and will be marked as fatally errored. "+
			"This can happen if there is another instance of chainlink running that is using the same private key, or if "+
			"an external wallet has been used to send a transaction from account %s with nonce %s."+
			" Please note that Chainlink requires exclusive ownership of it's private keys and sharing keys across multiple"+
			" chainlink instances, or using the chainlink keys with an external wallet is NOT SUPPORTED and WILL lead to missed transactions",
			oldTx.ID, head.BlockNumber(), latestFinalizedHead.BlockNumber(), oldTx.FromAddress, oldTx.Sequence.String()), "txID", oldTx.ID, "sequence", oldTx.Sequence.String(), "fromAddress", oldTx.FromAddress)

		// Signal pending tasks for these transactions as failed
		// Store errors and continue to allow all transactions a chance to be signaled
		if f.resumeCallback != nil && oldTx.PipelineTaskRunID.Valid && oldTx.SignalCallback && !oldTx.CallbackCompleted {
			err = f.resumeCallback(ctx, oldTx.PipelineTaskRunID.UUID, nil, errors.New(ErrCouldNotGetReceipt))
			if errors.Is(err, sql.ErrNoRows) {
				f.lggr.Debugw("callback missing or already resumed", "etxID", oldTx.ID)
			} else if err != nil {
				errorList = append(errorList, fmt.Errorf("failed to resume pipeline for ID %s: %w", oldTx.PipelineTaskRunID.UUID.String(), err))
				continue
			} else {
				// Mark tx as having completed callback
				if err = f.txStore.UpdateTxCallbackCompleted(ctx, oldTx.PipelineTaskRunID.UUID, f.chainID); err != nil {
					errorList = append(errorList, fmt.Errorf("failed to update callback as complete for tx ID %d: %w", oldTx.ID, err))
					continue
				}
			}
		}

		// Mark transaction as fatal error and delete attempts to prevent further receipt fetching
		oldTx.Error = null.StringFrom(ErrCouldNotGetReceipt)
		if err = f.txStore.UpdateTxFatalErrorAndDeleteAttempts(ctx, oldTx); err != nil {
			errorList = append(errorList, fmt.Errorf("failed to mark tx with ID %d as fatal: %w", oldTx.ID, err))
		}
	}
	if len(errorList) > 0 {
		return errors.Join(errorList...)
	}

	return nil
}

// findOldTxIDsWithoutReceipts finds IDs for transactions without receipts and attempts broadcasted at or before the finalized head
func findOldTxIDsWithoutReceipts(attempts []TxAttempt, receipts []*evmtypes.Receipt, latestFinalizedHead *evmtypes.Head) []int64 {
	if len(attempts) == 0 {
		return nil
	}
	txIDToAttemptsMap := make(map[int64][]TxAttempt)
	hashToReceiptMap := make(map[common.Hash]bool)
	// Store all receipts hashes in a map to easily access which attempt hash has a receipt
	for _, receipt := range receipts {
		hashToReceiptMap[receipt.TxHash] = true
	}
	// Store all attempts in a map of tx ID to attempts
	for _, attempt := range attempts {
		txIDToAttemptsMap[attempt.TxID] = append(txIDToAttemptsMap[attempt.TxID], attempt)
	}

	// Determine which transactions still do not have a receipt and if all of their attempts are older or equal to the latest finalized head
	oldTxIDs := make([]int64, 0, len(txIDToAttemptsMap))
	for txID, attempts := range txIDToAttemptsMap {
		hasReceipt := false
		hasAttemptAfterFinalizedHead := false
		for _, attempt := range attempts {
			if _, exists := hashToReceiptMap[attempt.Hash]; exists {
				hasReceipt = true
				break
			}
			if attempt.BroadcastBeforeBlockNum != nil && *attempt.BroadcastBeforeBlockNum > latestFinalizedHead.BlockNumber() {
				hasAttemptAfterFinalizedHead = true
				break
			}
		}
		if hasReceipt || hasAttemptAfterFinalizedHead {
			continue
		}
		oldTxIDs = append(oldTxIDs, txID)
	}
	return oldTxIDs
}

// buildTxHashList builds list of transaction hashes from receipts considered finalized
func (f *evmFinalizer) buildTxHashList(finalizedReceipts []*evmtypes.Receipt) []common.Hash {
	txHashes := make([]common.Hash, len(finalizedReceipts))
	for i, receipt := range finalizedReceipts {
		f.lggr.Debugw("transaction considered finalized",
			"txHash", receipt.TxHash.String(),
			"receiptBlockNum", receipt.BlockNumber.Int64(),
			"receiptBlockHash", receipt.BlockHash.String(),
		)
		txHashes[i] = receipt.TxHash
	}
	return txHashes
}
