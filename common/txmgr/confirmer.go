package txmgr

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	iutils "github.com/smartcontractkit/chainlink/v2/common/internal/utils"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// processHeadTimeout represents a sanity limit on how long ProcessHead
// should take to complete
const processHeadTimeout = 10 * time.Minute

var (
	promNumConfirmedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_confirmed_transactions",
		Help: "Total number of confirmed transactions. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
	}, []string{"chainID"})
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
	promTimeUntilTxConfirmed = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tx_manager_time_until_tx_confirmed",
		Help: "The amount of time elapsed from a transaction being broadcast to being included in a block.",
		Buckets: []float64{
			float64(500 * time.Millisecond),
			float64(time.Second),
			float64(5 * time.Second),
			float64(15 * time.Second),
			float64(30 * time.Second),
			float64(time.Minute),
			float64(2 * time.Minute),
			float64(5 * time.Minute),
			float64(10 * time.Minute),
		},
	}, []string{"chainID"})
	promBlocksUntilTxConfirmed = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tx_manager_blocks_until_tx_confirmed",
		Help: "The amount of blocks that have been mined from a transaction being broadcast to being included in a block.",
		Buckets: []float64{
			float64(1),
			float64(5),
			float64(10),
			float64(20),
			float64(50),
			float64(100),
		},
	}, []string{"chainID"})
)

type Confirmer[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	services.StateMachine
	txStore        txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	lggr           logger.SugaredLogger
	client         txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	resumeCallback ResumeCallback
	chainConfig    txmgrtypes.ConfirmerChainConfig
	txConfig       txmgrtypes.ConfirmerTransactionsConfig
	chainID        CHAIN_ID

	ks               txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	enabledAddresses []ADDR

	mb        *mailbox.Mailbox[types.Head[BLOCK_HASH]]
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup

	isReceiptNil func(R) bool
}

func NewConfirmer[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	client txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	chainConfig txmgrtypes.ConfirmerChainConfig,
	txConfig txmgrtypes.ConfirmerTransactionsConfig,
	keystore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	lggr logger.Logger,
	isReceiptNil func(R) bool,
) *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	lggr = logger.Named(lggr, "Confirmer")
	return &Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		txStore:        txStore,
		lggr:           logger.Sugared(lggr),
		client:         client,
		resumeCallback: nil,
		chainConfig:    chainConfig,
		txConfig:       txConfig,
		chainID:        client.ConfiguredChainID(),
		ks:             keystore,
		mb:             mailbox.NewSingle[types.Head[BLOCK_HASH]](),
		isReceiptNil:   isReceiptNil,
	}
}

func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Start(_ context.Context) error {
	return ec.StartOnce("Confirmer", func() (err error) {
		ec.enabledAddresses, err = ec.ks.EnabledAddressesForChain(ec.chainID)
		if err != nil {
			return fmt.Errorf("Confirmer: failed to load EnabledAddressesForChain: %w", err)
		}

		ec.ctx, ec.ctxCancel = context.WithCancel(context.Background())
		ec.wg = sync.WaitGroup{}
		ec.wg.Add(1)
		go ec.runLoop()
		return nil
	})
}

func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() error {
	return ec.StopOnce("Confirmer", func() (err error) {
		ec.ctxCancel()
		ec.wg.Wait()
		return nil
	})
}

func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SetResumeCallback(callback ResumeCallback) {
	ec.resumeCallback = callback
}

func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Name() string {
	return ec.lggr.Name()
}

func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HealthReport() map[string]error {
	return map[string]error{ec.Name(): ec.Healthy()}
}

func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) runLoop() {
	defer ec.wg.Done()
	keysChanged, unsub := ec.ks.SubscribeToKeyChanges()
	defer unsub()

	for {
		select {
		case <-ec.mb.Notify():
			for {
				if ec.ctx.Err() != nil {
					return
				}
				head, exists := ec.mb.Retrieve()
				if !exists {
					break
				}
				if err := ec.ProcessHead(ec.ctx, head); err != nil {
					ec.lggr.Errorw("Error processing head", "err", err)
					continue
				}
			}
		case <-keysChanged:
			var err error
			ec.enabledAddresses, err = ec.ks.EnabledAddressesForChain(ec.chainID)
			if err != nil {
				ec.lggr.Critical("Failed to reload key states after key change")
				continue
			}
		case <-ec.ctx.Done():
			return
		}
	}
}

// ProcessHead takes all required transactions for the confirmer on a new head.
// Handles all addresses, so there is no point of running it concurrently.
func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ProcessHead(ctx context.Context, head types.Head[BLOCK_HASH]) error {
	ctx, cancel := context.WithTimeout(ctx, processHeadTimeout)
	defer cancel()

	ec.lggr.Debugw("ProcessHead start", "headNum", head.BlockNumber())

	if err := ec.txStore.SetBroadcastBeforeBlockNum(ctx, head.BlockNumber(), ec.chainID); err != nil {
		return fmt.Errorf("SetBroadcastBeforeBlockNum failed: %w", err)
	}

	// TODO: Add addresses that are not enabled but we still have unconfirmed transactions for
	for _, from := range ec.enabledAddresses {
		total := time.Now()
		mark := time.Now()

		minedSequence, err := ec.client.SequenceAt(ctx, from, nil)
		if err != nil {
			return fmt.Errorf("unable to fetch pending sequence for address: %v: %w", from, err)
		}

		if err := ec.ConfirmUnconfirmedTransactions(ctx, head.BlockNumber(), from, minedSequence); err != nil {
			return fmt.Errorf("ConfirmUnconfirmedTransactions failed: %w", err)
		}

		ec.lggr.Debugw("Finished ConfirmUnconfirmedTransactions", "headNum", head.BlockNumber(), "time", time.Since(mark))
		mark = time.Now()

		if err := ec.HandleReorgedTransactions(ctx, head, from, minedSequence); err != nil {
			return fmt.Errorf("EnsureConfirmedTransactionsOnChain failed: %w", err)
		}

		//if err := ec.EnsureConfirmedTransactionReceiptsInLongestChain(ctx, head); err != nil {
		//	return fmt.Errorf("EnsureConfirmedTransactionReceiptsInLongestChain failed: %w", err)
		//}
		ec.lggr.Debugw("Finished Re-org detection", "headNum", head.BlockNumber(), "time", time.Since(mark))

		ec.lggr.Debugw("Finished transaction tracking.", "fromAddress", from, "headNum", head.BlockNumber(), "time", time.Since(total))
	}

	if ec.resumeCallback != nil {
		mark := time.Now()
		if err := ec.ResumePendingTaskRuns(ctx, head); err != nil {
			return fmt.Errorf("ResumePendingTaskRuns failed: %w", err)
		}

		ec.lggr.Debugw("Finished ResumePendingTaskRuns", "headNum", head.BlockNumber(), "time", time.Since(mark), "id", "confirmer")
	}

	ec.lggr.Debugw("processHead finish", "headNum", head.BlockNumber(), "id", "confirmer")

	return nil
}

func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ConfirmUnconfirmedTransactions(ctx context.Context, blockNum int64, from ADDR, minedSequence SEQ) error {
	attempts, err := ec.txStore.FindLikelyConfirmedTxAttemptsRequiringReceipt(ctx, ec.chainID, minedSequence)
	if err != nil {
		return fmt.Errorf("FindUnconfirmedTxAttemptsRequiringReceipt failed: %w", err)
	}
	if len(attempts) == 0 {
		return nil
	}
	promTxAttemptCount.WithLabelValues(ec.chainID.String()).Set(float64(len(attempts)))

	start := time.Now()
	allReceipts, err := ec.batchFetchReceipts(ctx, attempts, blockNum)
	if err != nil {
		return fmt.Errorf("batchFetchReceipts failed: %w", err)
	}

	if err := ec.txStore.SaveFetchedReceipts(ctx, allReceipts, ec.chainID); err != nil {
		return fmt.Errorf("saveFetchedReceipts failed: %w", err)
	}
	promNumConfirmedTxs.WithLabelValues(ec.chainID.String()).Add(float64(len(allReceipts)))

	observeUntilTxConfirmed(ec.chainID, attempts, allReceipts)

	// Any unconfirmed tx before the latest mined sequence that doesn't have a receipt is
	// marked as confirmed missing receipt. This is to avoid nonce gaps.
	if err := ec.txStore.MarkAllConfirmedMissingReceipt(ctx, ec.chainID); err != nil {
		return fmt.Errorf("unable to mark txes as 'confirmed_missing_receipt': %w", err)
	}

	ec.lggr.Debugw(fmt.Sprintf("Fetching and saving %v likely confirmed receipts done", len(attempts)),
		"time", time.Since(start))
	return nil
}

func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HandleReorgedTransactions(ctx context.Context, head types.Head[BLOCK_HASH], from ADDR, minedSequence SEQ) error {
	txs, err := ec.txStore.FindConfirmedTxsWithSequenceHigherThanMined(ctx, ec.chainID, minedSequence)
	if err != nil {
		return fmt.Errorf("FindUnconfirmedTxAttemptsRequiringReceipt failed: %w", err)
	}
	if len(txs) == 0 {
		return nil
	}

	for _, tx := range txs {
		ec.markForRebroadcast(*tx, head)
	}

	ec.lggr.Debugw(fmt.Sprintf("Marked %v confirmed and confirmed_missing_receipt txs as unconfirmed.", len(txs)))

	return nil
}

func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) batchFetchReceipts(ctx context.Context, attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], blockNum int64) (allReceipts []R, err error) {
	batchSize := int(ec.chainConfig.RPCDefaultBatchSize())
	if batchSize == 0 {
		batchSize = len(attempts)
	}

	for i := 0; i < len(attempts); i += batchSize {
		j := i + batchSize
		if j > len(attempts) {
			j = len(attempts)
		}

		ec.lggr.Debugw(fmt.Sprintf("Batch fetching receipts at indexes %v until (excluded) %v", i, j), "blockNum", blockNum)

		batch := attempts[i:j]

		// Metadata is required to determine whether a tx is forwarded or not.
		if ec.txConfig.ForwardersEnabled() {
			err = ec.txStore.PreloadTxes(ctx, batch)
			if err != nil {
				return nil, fmt.Errorf("Confirmer#batchFetchReceipts error loading txs for attempts: %w", err)
			}
		}

		lggr := ec.lggr.Named("BatchFetchReceipts").With("blockNum", blockNum)

		txReceipts, txErrs, err := ec.client.BatchGetReceipts(ctx, batch)
		if err != nil {
			return nil, err
		}

		for i := range txReceipts {
			attempt := batch[i]
			receipt := txReceipts[i]
			err := txErrs[i]

			l := attempt.Tx.GetLogger(lggr).With("txHash", attempt.Hash.String(), "txAttemptID", attempt.ID,
				"txID", attempt.TxID, "err", err, "sequence", attempt.Tx.Sequence,
			)

			if err != nil {
				l.Error("FetchReceipt failed")
				continue
			}

			if ec.isReceiptNil(receipt) {
				// NOTE: This should never happen, but it seems safer to check
				// regardless to avoid a potential panic
				l.AssumptionViolation("got nil receipt")
				continue
			}

			if receipt.IsZero() {
				l.Debug("Still waiting for receipt")
				continue
			}

			l = l.With("blockHash", receipt.GetBlockHash().String(), "status", receipt.GetStatus(), "transactionIndex", receipt.GetTransactionIndex())

			if receipt.IsUnmined() {
				l.Debug("Got receipt for transaction but it's still in the mempool and not included in a block yet")
				continue
			}

			l.Debugw("Got receipt for transaction", "blockNumber", receipt.GetBlockNumber(), "feeUsed", receipt.GetFeeUsed())

			if receipt.GetTxHash().String() != attempt.Hash.String() {
				l.Errorf("Invariant violation, expected receipt with hash %s to have same hash as attempt with hash %s", receipt.GetTxHash().String(), attempt.Hash.String())
				continue
			}

			if receipt.GetBlockNumber() == nil {
				l.Error("Invariant violation, receipt was missing block number")
				continue
			}

			if receipt.GetStatus() == 0 {
				rpcError, errExtract := ec.client.CallContract(ctx, attempt, receipt.GetBlockNumber())
				if errExtract == nil {
					l.Warnw("transaction reverted on-chain", "hash", receipt.GetTxHash(), "rpcError", rpcError.String())
				} else {
					l.Warnw("transaction reverted on-chain unable to extract revert reason", "hash", receipt.GetTxHash(), "err", err)
				}
				// This might increment more than once e.g. in case of re-orgs going back and forth we might re-fetch the same receipt
				promRevertedTxCount.WithLabelValues(ec.chainID.String()).Add(1)
			} else {
				promNumSuccessfulTxs.WithLabelValues(ec.chainID.String()).Add(1)
			}

			// This is only recording forwarded tx that were mined and have a status.
			// Counters are prone to being inaccurate due to re-orgs.
			if ec.txConfig.ForwardersEnabled() {
				meta, metaErr := attempt.Tx.GetMeta()
				if metaErr == nil && meta != nil && meta.FwdrDestAddress != nil {
					// promFwdTxCount takes two labels, chainId and a boolean of whether a tx was successful or not.
					promFwdTxCount.WithLabelValues(ec.chainID.String(), strconv.FormatBool(receipt.GetStatus() != 0)).Add(1)
				}
			}

			allReceipts = append(allReceipts, receipt)
		}

	}

	return
}

// TODO: best approach would be to clear all attempts and receipts and mark tx as in progress and let the broadcaster handle it from the beginning, otherwise last attempt might not be accurate at all.
func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) markForRebroadcast(tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], head types.Head[BLOCK_HASH]) error {
	if len(tx.TxAttempts) == 0 {
		return fmt.Errorf("invariant violation: expected tx %v to have at least one attempt", tx.ID)
	}

	// Rebroadcast the one with the highest gas price
	attempt := tx.TxAttempts[0]
	var receipt txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH]
	if len(attempt.Receipts) > 0 {
		receipt = attempt.Receipts[0]
	}

	logValues := []interface{}{
		"txhash", attempt.Hash.String(),
		"currentBlockNum", head.BlockNumber(),
		"currentBlockHash", head.BlockHash().String(),
		"txID", tx.ID,
		"attemptID", attempt.ID,
		"nReceipts", len(attempt.Receipts),
		"id", "confirmer",
	}

	// nil check on receipt interface
	if receipt != nil {
		logValues = append(logValues,
			"replacementBlockHashAtConfirmedHeight", head.HashAtHeight(receipt.GetBlockNumber().Int64()),
			"confirmedInBlockNum", receipt.GetBlockNumber(),
			"confirmedInBlockHash", receipt.GetBlockHash(),
			"confirmedInTxIndex", receipt.GetTransactionIndex(),
		)
	}

	ec.lggr.Infow(fmt.Sprintf("Re-org detected. Rebroadcasting transaction %s which may have been re-org'd out of the main chain", attempt.Hash.String()), logValues...)

	// Put it back in progress and delete all receipts (they do not apply to the new chain)
	if err := ec.txStore.UpdateTxForRebroadcast(ec.ctx, tx, attempt); err != nil {
		return fmt.Errorf("markForRebroadcast failed: %w", err)
	}

	return nil
}

// ResumePendingTaskRuns issues callbacks to task runs that are pending waiting for receipts
func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ResumePendingTaskRuns(ctx context.Context, head types.Head[BLOCK_HASH]) error {

	receiptsPlus, err := ec.txStore.FindTxesPendingCallback(ctx, head.BlockNumber(), ec.chainID)

	if err != nil {
		return err
	}

	if len(receiptsPlus) > 0 {
		ec.lggr.Debugf("Resuming %d task runs pending receipt", len(receiptsPlus))
	} else {
		ec.lggr.Debug("No task runs to resume")
	}
	for _, data := range receiptsPlus {
		var taskErr error
		var output interface{}
		if data.FailOnRevert && data.Receipt.GetStatus() == 0 {
			taskErr = fmt.Errorf("transaction %s reverted on-chain", data.Receipt.GetTxHash())
		} else {
			output = data.Receipt
		}

		ec.lggr.Debugw("Callback: resuming tx with receipt", "output", output, "taskErr", taskErr, "pipelineTaskRunID", data.ID)
		if err := ec.resumeCallback(data.ID, output, taskErr); err != nil {
			return fmt.Errorf("failed to resume suspended pipeline run: %w", err)
		}
		// Mark tx as having completed callback
		if err := ec.txStore.UpdateTxCallbackCompleted(ctx, data.ID, ec.chainID); err != nil {
			return err
		}
	}

	return nil
}

// observeUntilTxConfirmed observes the promBlocksUntilTxConfirmed metric for each confirmed
// transaction.
func observeUntilTxConfirmed[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](chainID CHAIN_ID, attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], receipts []R) {
	for _, attempt := range attempts {
		for _, r := range receipts {
			if attempt.Hash.String() != r.GetTxHash().String() {
				continue
			}

			// We estimate the time until confirmation by subtracting from the time the tx (not the attempt)
			// was created. We want to measure the amount of time taken from when a transaction is created
			// via e.g Txm.CreateTransaction to when it is confirmed on-chain, regardless of how many attempts
			// were needed to achieve this.
			duration := time.Since(attempt.Tx.CreatedAt)
			promTimeUntilTxConfirmed.
				WithLabelValues(chainID.String()).
				Observe(float64(duration))

			// Since a tx can have many attempts, we take the number of blocks to confirm as the block number
			// of the receipt minus the block number of the first ever broadcast for this transaction.
			broadcastBefore := iutils.MinFunc(attempt.Tx.TxAttempts, func(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) int64 {
				if attempt.BroadcastBeforeBlockNum != nil {
					return *attempt.BroadcastBeforeBlockNum
				}
				return 0
			})
			if broadcastBefore > 0 {
				blocksElapsed := r.GetBlockNumber().Int64() - broadcastBefore
				promBlocksUntilTxConfirmed.
					WithLabelValues(chainID.String()).
					Observe(float64(blocksElapsed))
			}
		}
	}
}

// There is a chance some transactions get reorged but they get included in a different block, leading to a confirmed transaction with a different receipt.
// The TXM cab still detect the nonce difference and keep functioning, but we will have a mismatch on the receipt for a given transaction.
// To fully protect against this we would had to fetch every receipt for every confirmed transaction after the finality depth.
// To mitigate this, we do a best effort search in the recent longest chain to try and detect such cases. If a receipt is not found, the tx is marked as
// unconfirmed and enters the next confirmer cycle.
//func (ec *Confirmer[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) EnsureConfirmedTransactionReceiptsInLongestChain(ctx context.Context, head types.Head[BLOCK_HASH]) error {
//	etxs, err := ec.txStore.FindTransactionsConfirmedInBlockRange(ctx, head.BlockNumber(), head.EarliestHeadInChain().BlockNumber(), ec.chainID)
//	if err != nil {
//		return fmt.Errorf("findTransactionsConfirmedInBlockRange failed: %w", err)
//	}
//
//	for _, etx := range etxs {
//		if !hasReceiptInLongestChain(*etx, head) {
//			if err := ec.markForRebroadcast(*etx, head); err != nil {
//				return fmt.Errorf("markForRebroadcast failed for etx %v: %w", etx.ID, err)
//			}
//		}
//	}
//
//	return nil
//}
//
//func hasReceiptInLongestChain[
//	CHAIN_ID types.ID,
//	ADDR types.Hashable,
//	TX_HASH, BLOCK_HASH types.Hashable,
//	SEQ types.Sequence,
//	FEE feetypes.Fee,
//](etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], head types.Head[BLOCK_HASH]) bool {
//	for {
//		for _, attempt := range etx.TxAttempts {
//			for _, receipt := range attempt.Receipts {
//				if receipt.GetBlockHash().String() == head.BlockHash().String() && receipt.GetBlockNumber().Int64() == head.BlockNumber() {
//					return true
//				}
//			}
//		}
//		if head.GetParent() == nil {
//			return false
//		}
//		head = head.GetParent()
//	}
//}
//