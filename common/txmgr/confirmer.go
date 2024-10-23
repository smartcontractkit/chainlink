package txmgr

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"

	commonhex "github.com/smartcontractkit/chainlink-common/pkg/utils/hex"

	"github.com/smartcontractkit/chainlink-common/pkg/chains/label"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/common/client"
	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	iutils "github.com/smartcontractkit/chainlink/v2/common/internal/utils"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

const (
	// processHeadTimeout represents a sanity limit on how long ProcessHead
	// should take to complete
	processHeadTimeout = 10 * time.Minute
)

var (
	promNumGasBumps = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_gas_bumps",
		Help: "Number of gas bumps",
	}, []string{"chainID"})

	promGasBumpExceedsLimit = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_gas_bump_exceeds_limit",
		Help: "Number of times gas bumping failed from exceeding the configured limit. Any counts of this type indicate a serious problem.",
	}, []string{"chainID"})
	promNumConfirmedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_confirmed_transactions",
		Help: "Total number of confirmed transactions. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
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

// Confirmer is a broad service which performs four different tasks in sequence on every new longest chain
// Step 1: Mark that all currently pending transaction attempts were broadcast before this block
// Step 2: Check pending transactions for confirmation and confirmed transactions for re-org
// Step 3: Check if any pending transaction is stuck in the mempool. If so, mark for purge.
// Step 4: See if any transactions have exceeded the gas bumping block threshold and, if so, bump them
type Confirmer[
	CHAIN_ID types.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	services.StateMachine
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	lggr    logger.SugaredLogger
	client  txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	stuckTxDetector txmgrtypes.StuckTxDetector[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	resumeCallback  ResumeCallback
	feeConfig       txmgrtypes.ConfirmerFeeConfig
	txConfig        txmgrtypes.ConfirmerTransactionsConfig
	dbConfig        txmgrtypes.ConfirmerDatabaseConfig
	chainID         CHAIN_ID

	ks               txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	enabledAddresses []ADDR

	mb           *mailbox.Mailbox[HEAD]
	stopCh       services.StopChan
	wg           sync.WaitGroup
	initSync     sync.Mutex
	isStarted    bool
	isReceiptNil func(R) bool
}

func NewConfirmer[
	CHAIN_ID types.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	client txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	feeConfig txmgrtypes.ConfirmerFeeConfig,
	txConfig txmgrtypes.ConfirmerTransactionsConfig,
	dbConfig txmgrtypes.ConfirmerDatabaseConfig,
	keystore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	lggr logger.Logger,
	isReceiptNil func(R) bool,
	stuckTxDetector txmgrtypes.StuckTxDetector[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	lggr = logger.Named(lggr, "Confirmer")
	return &Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		txStore:          txStore,
		lggr:             logger.Sugared(lggr),
		client:           client,
		TxAttemptBuilder: txAttemptBuilder,
		resumeCallback:   nil,
		feeConfig:        feeConfig,
		txConfig:         txConfig,
		dbConfig:         dbConfig,
		chainID:          client.ConfiguredChainID(),
		ks:               keystore,
		mb:               mailbox.NewSingle[HEAD](),
		isReceiptNil:     isReceiptNil,
		stuckTxDetector:  stuckTxDetector,
	}
}

// Start is a comment to appease the linter
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Start(ctx context.Context) error {
	return ec.StartOnce("Confirmer", func() error {
		if ec.feeConfig.BumpThreshold() == 0 {
			ec.lggr.Infow("Gas bumping is disabled (FeeEstimator.BumpThreshold set to 0)", "feeBumpThreshold", 0)
		} else {
			ec.lggr.Infow(fmt.Sprintf("Fee bumping is enabled, unconfirmed transactions will have their fee bumped every %d blocks", ec.feeConfig.BumpThreshold()), "feeBumpThreshold", ec.feeConfig.BumpThreshold())
		}

		return ec.startInternal(ctx)
	})
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) startInternal(ctx context.Context) error {
	ec.initSync.Lock()
	defer ec.initSync.Unlock()
	if ec.isStarted {
		return errors.New("Confirmer is already started")
	}
	var err error
	ec.enabledAddresses, err = ec.ks.EnabledAddressesForChain(ctx, ec.chainID)
	if err != nil {
		return fmt.Errorf("Confirmer: failed to load EnabledAddressesForChain: %w", err)
	}
	if err = ec.stuckTxDetector.LoadPurgeBlockNumMap(ctx, ec.enabledAddresses); err != nil {
		ec.lggr.Debugf("Confirmer: failed to load the last purged block num for enabled addresses. Process can continue as normal but purge rate limiting may be affected.")
	}

	ec.stopCh = make(chan struct{})
	ec.wg = sync.WaitGroup{}
	ec.wg.Add(1)
	go ec.runLoop()
	ec.isStarted = true
	return nil
}

// Close is a comment to appease the linter
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() error {
	return ec.StopOnce("Confirmer", func() error {
		return ec.closeInternal()
	})
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) closeInternal() error {
	ec.initSync.Lock()
	defer ec.initSync.Unlock()
	if !ec.isStarted {
		return fmt.Errorf("Confirmer is not started: %w", services.ErrAlreadyStopped)
	}
	close(ec.stopCh)
	ec.wg.Wait()
	ec.isStarted = false
	return nil
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SetResumeCallback(callback ResumeCallback) {
	ec.resumeCallback = callback
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Name() string {
	return ec.lggr.Name()
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HealthReport() map[string]error {
	return map[string]error{ec.Name(): ec.Healthy()}
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) runLoop() {
	defer ec.wg.Done()
	ctx, cancel := ec.stopCh.NewCtx()
	defer cancel()
	for {
		select {
		case <-ec.mb.Notify():
			for {
				if ctx.Err() != nil {
					return
				}
				head, exists := ec.mb.Retrieve()
				if !exists {
					break
				}
				if err := ec.ProcessHead(ctx, head); err != nil {
					ec.lggr.Errorw("Error processing head", "err", err)
					continue
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// ProcessHead takes all required transactions for the confirmer on a new head
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ProcessHead(ctx context.Context, head types.Head[BLOCK_HASH]) error {
	ctx, cancel := context.WithTimeout(ctx, processHeadTimeout)
	defer cancel()
	return ec.processHead(ctx, head)
}

// NOTE: This SHOULD NOT be run concurrently or it could behave badly
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) processHead(ctx context.Context, head types.Head[BLOCK_HASH]) error {
	ec.lggr.Debugw("processHead start", "headNum", head.BlockNumber(), "id", "confirmer")

	mark := time.Now()
	if err := ec.txStore.SetBroadcastBeforeBlockNum(ctx, head.BlockNumber(), ec.chainID); err != nil {
		return err
	}
	ec.lggr.Debugw("Finished SetBroadcastBeforeBlockNum", "headNum", head.BlockNumber(), "time", time.Since(mark), "id", "confirmer")

	mark = time.Now()
	if err := ec.CheckForConfirmation(ctx, head); err != nil {
		return err
	}
	ec.lggr.Debugw("Finished CheckForConfirmation", "headNum", head.BlockNumber(), "time", time.Since(mark), "id", "confirmer")

	mark = time.Now()
	if err := ec.ProcessStuckTransactions(ctx, head.BlockNumber()); err != nil {
		return err
	}
	ec.lggr.Debugw("Finished ProcessStuckTransactions", "headNum", head.BlockNumber(), "time", time.Since(mark), "id", "confirmer")

	mark = time.Now()
	if err := ec.RebroadcastWhereNecessary(ctx, head.BlockNumber()); err != nil {
		return err
	}
	ec.lggr.Debugw("Finished RebroadcastWhereNecessary", "headNum", head.BlockNumber(), "time", time.Since(mark), "id", "confirmer")
	ec.lggr.Debugw("processHead finish", "headNum", head.BlockNumber(), "id", "confirmer")

	return nil
}

// CheckForConfirmation fetches the mined transaction count for each enabled address and marks transactions with a lower sequence as confirmed and ones with equal or higher sequence as unconfirmed
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CheckForConfirmation(ctx context.Context, head types.Head[BLOCK_HASH]) error {
	var errorList []error
	for _, fromAddress := range ec.enabledAddresses {
		minedTxCount, err := ec.client.SequenceAt(ctx, fromAddress, nil)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("unable to fetch mined transaction count for address %s: %w", fromAddress.String(), err))
			continue
		}
		reorgTxs, includedTxs, err := ec.txStore.FindReorgOrIncludedTxs(ctx, fromAddress, minedTxCount, ec.chainID)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("failed to find re-org'd or included transactions based on the mined transaction count %d: %w", minedTxCount.Int64(), err))
			continue
		}
		// If re-org'd transactions are identified, process them and mark them for rebroadcast
		err = ec.ProcessReorgTxs(ctx, reorgTxs, head)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("failed to process re-org'd transactions: %w", err))
			continue
		}
		// If unconfirmed transactions are identified as included, process them and mark them as confirmed or terminally stuck (if purge attempt exists)
		err = ec.ProcessIncludedTxs(ctx, includedTxs, head)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("failed to process confirmed transactions: %w", err))
			continue
		}
	}
	if len(errorList) > 0 {
		return errors.Join(errorList...)
	}
	return nil
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ProcessReorgTxs(ctx context.Context, reorgTxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], head types.Head[BLOCK_HASH]) error {
	if len(reorgTxs) == 0 {
		return nil
	}
	etxIDs := make([]int64, 0, len(reorgTxs))
	attemptIDs := make([]int64, 0, len(reorgTxs))
	for _, etx := range reorgTxs {
		if len(etx.TxAttempts) == 0 {
			return fmt.Errorf("invariant violation: expected tx %v to have at least one attempt", etx.ID)
		}

		// Rebroadcast the one with the highest gas price
		attempt := etx.TxAttempts[0]

		logValues := []interface{}{
			"txhash", attempt.Hash.String(),
			"currentBlockNum", head.BlockNumber(),
			"currentBlockHash", head.BlockHash().String(),
			"txID", etx.ID,
			"attemptID", attempt.ID,
			"nReceipts", len(attempt.Receipts),
			"attemptState", attempt.State,
			"id", "confirmer",
		}

		if len(attempt.Receipts) > 0 && attempt.Receipts[0] != nil {
			receipt := attempt.Receipts[0]
			logValues = append(logValues,
				"replacementBlockHashAtConfirmedHeight", head.HashAtHeight(receipt.GetBlockNumber().Int64()),
				"confirmedInBlockNum", receipt.GetBlockNumber(),
				"confirmedInBlockHash", receipt.GetBlockHash(),
				"confirmedInTxIndex", receipt.GetTransactionIndex(),
			)
		}

		if etx.State == TxFinalized {
			ec.lggr.AssumptionViolationw(fmt.Sprintf("Re-org detected for finalized transaction. This should never happen. Rebroadcasting transaction %s which may have been re-org'd out of the main chain", attempt.Hash.String()), logValues...)
		} else {
			ec.lggr.Infow(fmt.Sprintf("Re-org detected. Rebroadcasting transaction %s which may have been re-org'd out of the main chain", attempt.Hash.String()), logValues...)
		}

		etxIDs = append(etxIDs, etx.ID)
		attemptIDs = append(attemptIDs, attempt.ID)
	}

	// Mark transactions as unconfirmed, mark attempts as in-progress, and delete receipts since they do not apply to the new chain
	// This may revert some fatal error transactions to unconfirmed if terminally stuck transactions purge attempts get re-org'd
	return ec.txStore.UpdateTxsForRebroadcast(ctx, etxIDs, attemptIDs)
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ProcessIncludedTxs(ctx context.Context, includedTxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], head types.Head[BLOCK_HASH]) error {
	if len(includedTxs) == 0 {
		return nil
	}
	// Add newly confirmed transactions to the prom metric
	promNumConfirmedTxs.WithLabelValues(ec.chainID.String()).Add(float64(len(includedTxs)))

	purgeTxIDs := make([]int64, 0, len(includedTxs))
	confirmedTxIDs := make([]int64, 0, len(includedTxs))
	for _, tx := range includedTxs {
		// If any attempt in the transaction is marked for purge, the transaction was terminally stuck and should be marked as fatal error
		if tx.HasPurgeAttempt() {
			// Setting the purged block num here is ok since we have confirmation the tx has been included
			ec.stuckTxDetector.SetPurgeBlockNum(tx.FromAddress, head.BlockNumber())
			purgeTxIDs = append(purgeTxIDs, tx.ID)
			continue
		}
		confirmedTxIDs = append(confirmedTxIDs, tx.ID)
		observeUntilTxConfirmed(ec.chainID, tx.TxAttempts, head)
	}
	// Mark the transactions included on-chain with a purge attempt as fatal error with the terminally stuck error message
	if err := ec.txStore.UpdateTxFatalError(ctx, purgeTxIDs, ec.stuckTxDetector.StuckTxFatalError()); err != nil {
		return fmt.Errorf("failed to update terminally stuck transactions: %w", err)
	}
	// Mark the transactions included on-chain as confirmed
	if err := ec.txStore.UpdateTxConfirmed(ctx, confirmedTxIDs); err != nil {
		return fmt.Errorf("failed to update confirmed transactions: %w", err)
	}
	return nil
}

// Determines if any of the unconfirmed transactions are terminally stuck for each enabled address
// If any transaction is found to be terminally stuck, this method sends an empty attempt with bumped gas in an attempt to purge the stuck transaction
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ProcessStuckTransactions(ctx context.Context, blockNum int64) error {
	// Use the detector to find a stuck tx for each enabled address
	stuckTxs, err := ec.stuckTxDetector.DetectStuckTransactions(ctx, ec.enabledAddresses, blockNum)
	if err != nil {
		return fmt.Errorf("failed to detect stuck transactions: %w", err)
	}
	if len(stuckTxs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(stuckTxs))
	errorList := []error{}
	var errMu sync.Mutex
	for _, tx := range stuckTxs {
		// All stuck transactions will have unique from addresses. It is safe to process separate keys concurrently
		// NOTE: This design will block one key if another takes a really long time to execute
		go func(tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
			defer wg.Done()
			lggr := tx.GetLogger(ec.lggr)
			// Create a purge attempt for tx
			purgeAttempt, err := ec.TxAttemptBuilder.NewPurgeTxAttempt(ctx, tx, lggr)
			if err != nil {
				errMu.Lock()
				errorList = append(errorList, fmt.Errorf("failed to create a purge attempt: %w", err))
				errMu.Unlock()
				return
			}
			// Save purge attempt
			if err := ec.txStore.SaveInProgressAttempt(ctx, &purgeAttempt); err != nil {
				errMu.Lock()
				errorList = append(errorList, fmt.Errorf("failed to save purge attempt: %w", err))
				errMu.Unlock()
				return
			}
			lggr.Warnw("marked transaction as terminally stuck", "etx", tx)
			// Send purge attempt
			if err := ec.handleInProgressAttempt(ctx, lggr, tx, purgeAttempt, blockNum); err != nil {
				errMu.Lock()
				errorList = append(errorList, fmt.Errorf("failed to send purge attempt: %w", err))
				errMu.Unlock()
				return
			}
			// Resume pending task runs with failure for stuck transactions
			if err := ec.resumeFailedTaskRuns(ctx, tx); err != nil {
				errMu.Lock()
				errorList = append(errorList, fmt.Errorf("failed to resume pending task run for transaction: %w", err))
				errMu.Unlock()
				return
			}
		}(tx)
	}
	wg.Wait()
	return errors.Join(errorList...)
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) resumeFailedTaskRuns(ctx context.Context, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	if !etx.PipelineTaskRunID.Valid || ec.resumeCallback == nil || !etx.SignalCallback || etx.CallbackCompleted {
		return nil
	}
	err := ec.resumeCallback(ctx, etx.PipelineTaskRunID.UUID, nil, errors.New(ec.stuckTxDetector.StuckTxFatalError()))
	if errors.Is(err, sql.ErrNoRows) {
		ec.lggr.Debugw("callback missing or already resumed", "etxID", etx.ID)
	} else if err != nil {
		return fmt.Errorf("failed to resume pipeline: %w", err)
	} else {
		// Mark tx as having completed callback
		if err := ec.txStore.UpdateTxCallbackCompleted(ctx, etx.PipelineTaskRunID.UUID, ec.chainID); err != nil {
			return err
		}
	}
	return nil
}

// RebroadcastWhereNecessary bumps gas or resends transactions that were previously out-of-funds
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RebroadcastWhereNecessary(ctx context.Context, blockHeight int64) error {
	var wg sync.WaitGroup

	// It is safe to process separate keys concurrently
	// NOTE: This design will block one key if another takes a really long time to execute
	wg.Add(len(ec.enabledAddresses))
	errors := []error{}
	var errMu sync.Mutex
	for _, address := range ec.enabledAddresses {
		go func(fromAddress ADDR) {
			if err := ec.rebroadcastWhereNecessary(ctx, fromAddress, blockHeight); err != nil {
				errMu.Lock()
				errors = append(errors, err)
				errMu.Unlock()
				ec.lggr.Errorw("Error in RebroadcastWhereNecessary", "err", err, "fromAddress", fromAddress)
			}

			wg.Done()
		}(address)
	}

	wg.Wait()

	return multierr.Combine(errors...)
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) rebroadcastWhereNecessary(ctx context.Context, address ADDR, blockHeight int64) error {
	if err := ec.handleAnyInProgressAttempts(ctx, address, blockHeight); err != nil {
		return fmt.Errorf("handleAnyInProgressAttempts failed: %w", err)
	}

	threshold := int64(ec.feeConfig.BumpThreshold())
	bumpDepth := int64(ec.feeConfig.BumpTxDepth())
	maxInFlightTransactions := ec.txConfig.MaxInFlight()
	etxs, err := ec.FindTxsRequiringRebroadcast(ctx, ec.lggr, address, blockHeight, threshold, bumpDepth, maxInFlightTransactions, ec.chainID)
	if err != nil {
		return fmt.Errorf("FindTxsRequiringRebroadcast failed: %w", err)
	}
	for _, etx := range etxs {
		lggr := etx.GetLogger(ec.lggr)

		attempt, err := ec.attemptForRebroadcast(ctx, lggr, *etx)
		if err != nil {
			return fmt.Errorf("attemptForRebroadcast failed: %w", err)
		}

		lggr.Debugw("Rebroadcasting transaction", "nPreviousAttempts", len(etx.TxAttempts), "fee", attempt.TxFee)

		if err := ec.txStore.SaveInProgressAttempt(ctx, &attempt); err != nil {
			return fmt.Errorf("saveInProgressAttempt failed: %w", err)
		}

		if err := ec.handleInProgressAttempt(ctx, lggr, *etx, attempt, blockHeight); err != nil {
			return fmt.Errorf("handleInProgressAttempt failed: %w", err)
		}
	}
	return nil
}

// "in_progress" attempts were left behind after a crash/restart and may or may not have been sent.
// We should try to ensure they get on-chain so we can fetch a receipt for them.
// NOTE: We also use this to mark attempts for rebroadcast in event of a
// re-org, so multiple attempts are allowed to be in in_progress state (but
// only one per tx).
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleAnyInProgressAttempts(ctx context.Context, address ADDR, blockHeight int64) error {
	attempts, err := ec.txStore.GetInProgressTxAttempts(ctx, address, ec.chainID)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return fmt.Errorf("GetInProgressTxAttempts failed: %w", err)
	}
	for _, a := range attempts {
		err := ec.handleInProgressAttempt(ctx, a.Tx.GetLogger(ec.lggr), a.Tx, a, blockHeight)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			return fmt.Errorf("handleInProgressAttempt failed: %w", err)
		}
	}
	return nil
}

// FindTxsRequiringRebroadcast returns attempts that hit insufficient native tokens,
// and attempts that need bumping, in sequence ASC order
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxsRequiringRebroadcast(ctx context.Context, lggr logger.Logger, address ADDR, blockNum, gasBumpThreshold, bumpDepth int64, maxInFlightTransactions uint32, chainID CHAIN_ID) (etxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	// NOTE: These two queries could be combined into one using union but it
	// becomes harder to read and difficult to test in isolation. KISS principle
	etxInsufficientFunds, err := ec.txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, address, chainID)
	if err != nil {
		return nil, err
	}

	if len(etxInsufficientFunds) > 0 {
		lggr.Infow(fmt.Sprintf("Found %d transactions to be re-sent that were previously rejected due to insufficient native token balance", len(etxInsufficientFunds)), "blockNum", blockNum, "address", address)
	}

	etxBumps, err := ec.txStore.FindTxsRequiringGasBump(ctx, address, blockNum, gasBumpThreshold, bumpDepth, chainID)
	if ctx.Err() != nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if len(etxBumps) > 0 {
		// txes are ordered by sequence asc so the first will always be the oldest
		etx := etxBumps[0]
		// attempts are ordered by time sent asc so first will always be the oldest
		var oldestBlocksBehind int64 = -1 // It should never happen that the oldest attempt has no BroadcastBeforeBlockNum set, but in case it does, we shouldn't crash - log this sentinel value instead
		if len(etx.TxAttempts) > 0 {
			oldestBlockNum := etx.TxAttempts[0].BroadcastBeforeBlockNum
			if oldestBlockNum != nil {
				oldestBlocksBehind = blockNum - *oldestBlockNum
			}
		} else {
			logger.Sugared(lggr).AssumptionViolationw("Expected tx for gas bump to have at least one attempt", "etxID", etx.ID, "blockNum", blockNum, "address", address)
		}
		lggr.Infow(fmt.Sprintf("Found %d transactions to re-sent that have still not been confirmed after at least %d blocks. The oldest of these has not still not been confirmed after %d blocks. These transactions will have their gas price bumped. %s", len(etxBumps), gasBumpThreshold, oldestBlocksBehind, label.NodeConnectivityProblemWarning), "blockNum", blockNum, "address", address, "gasBumpThreshold", gasBumpThreshold)
	}

	seen := make(map[int64]struct{})

	for _, etx := range etxInsufficientFunds {
		seen[etx.ID] = struct{}{}
		etxs = append(etxs, etx)
	}
	for _, etx := range etxBumps {
		if _, exists := seen[etx.ID]; !exists {
			etxs = append(etxs, etx)
		}
	}

	sort.Slice(etxs, func(i, j int) bool {
		return (*etxs[i].Sequence).Int64() < (*etxs[j].Sequence).Int64()
	})

	if maxInFlightTransactions > 0 && len(etxs) > int(maxInFlightTransactions) {
		lggr.Warnf("%d transactions to rebroadcast which exceeds limit of %d. %s", len(etxs), maxInFlightTransactions, label.MaxInFlightTransactionsWarning)
		etxs = etxs[:maxInFlightTransactions]
	}

	return
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) attemptForRebroadcast(ctx context.Context, lggr logger.Logger, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) (attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	if len(etx.TxAttempts) > 0 {
		etx.TxAttempts[0].Tx = etx
		previousAttempt := etx.TxAttempts[0]
		logFields := ec.logFieldsPreviousAttempt(previousAttempt)
		if previousAttempt.State == txmgrtypes.TxAttemptInsufficientFunds {
			// Do not create a new attempt if we ran out of funds last time since bumping gas is pointless
			// Instead try to resubmit the same attempt at the same price, in the hope that the wallet was funded since our last attempt
			lggr.Debugw("Rebroadcast InsufficientFunds", logFields...)
			previousAttempt.State = txmgrtypes.TxAttemptInProgress
			return previousAttempt, nil
		}
		attempt, err = ec.bumpGas(ctx, etx, etx.TxAttempts)

		if commonfee.IsBumpErr(err) {
			lggr.Errorw("Failed to bump gas", append(logFields, "err", err)...)
			// Do not create a new attempt if bumping gas would put us over the limit or cause some other problem
			// Instead try to resubmit the previous attempt, and keep resubmitting until its accepted
			previousAttempt.BroadcastBeforeBlockNum = nil
			previousAttempt.State = txmgrtypes.TxAttemptInProgress
			return previousAttempt, nil
		}
		return attempt, err
	}
	return attempt, fmt.Errorf("invariant violation: Tx %v was unconfirmed but didn't have any attempts. "+
		"Falling back to default gas price instead."+
		"This is a bug! Please report to https://github.com/smartcontractkit/chainlink/issues", etx.ID)
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) logFieldsPreviousAttempt(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) []interface{} {
	etx := attempt.Tx
	return []interface{}{
		"etxID", etx.ID,
		"txHash", attempt.Hash,
		"previousAttempt", attempt,
		"feeLimit", attempt.ChainSpecificFeeLimit,
		"callerProvidedFeeLimit", etx.FeeLimit,
		"maxGasPrice", ec.feeConfig.MaxFeePrice(),
		"sequence", etx.Sequence,
	}
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) bumpGas(ctx context.Context, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], previousAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) (bumpedAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	previousAttempt := previousAttempts[0]
	logFields := ec.logFieldsPreviousAttempt(previousAttempt)

	var bumpedFee FEE
	var bumpedFeeLimit uint64
	bumpedAttempt, bumpedFee, bumpedFeeLimit, _, err = ec.NewBumpTxAttempt(ctx, etx, previousAttempt, previousAttempts, ec.lggr)

	// if no error, return attempt
	// if err, continue below
	if err == nil {
		promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
		ec.lggr.Debugw("Rebroadcast bumping fee for tx", append(logFields, "bumpedFee", bumpedFee.String(), "bumpedFeeLimit", bumpedFeeLimit)...)
		return bumpedAttempt, err
	}

	if errors.Is(err, commonfee.ErrBumpFeeExceedsLimit) {
		promGasBumpExceedsLimit.WithLabelValues(ec.chainID.String()).Inc()
	}

	return bumpedAttempt, fmt.Errorf("error bumping gas: %w", err)
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleInProgressAttempt(ctx context.Context, lggr logger.SugaredLogger, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], blockHeight int64) error {
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return fmt.Errorf("invariant violation: expected tx_attempt %v to be in_progress, it was %s", attempt.ID, attempt.State)
	}

	now := time.Now()
	lggr.Debugw("Sending transaction", "txAttemptID", attempt.ID, "txHash", attempt.Hash, "meta", etx.Meta, "feeLimit", attempt.ChainSpecificFeeLimit, "callerProvidedFeeLimit", etx.FeeLimit, "attempt", attempt, "etx", etx)
	errType, sendError := ec.client.SendTransactionReturnCode(ctx, etx, attempt, lggr)

	switch errType {
	case client.Underpriced:
		// This should really not ever happen in normal operation since we
		// already bumped above the required minimum in broadcaster.
		ec.lggr.Warnw("Got terminally underpriced error for gas bump, this should never happen unless the remote RPC node changed its configuration on the fly, or you are using multiple RPC nodes with different minimum gas price requirements. This is not recommended", "attempt", attempt)
		// "Lazily" load attempts here since the overwhelmingly common case is
		// that we don't need them unless we enter this path
		if err := ec.txStore.LoadTxAttempts(ctx, &etx); err != nil {
			return fmt.Errorf("failed to load TxAttempts while bumping on terminally underpriced error: %w", err)
		}
		if len(etx.TxAttempts) == 0 {
			err := errors.New("expected to find at least 1 attempt")
			ec.lggr.AssumptionViolationw(err.Error(), "err", err, "attempt", attempt)
			return err
		}
		if attempt.ID != etx.TxAttempts[0].ID {
			err := errors.New("expected highest priced attempt to be the current in_progress attempt")
			ec.lggr.AssumptionViolationw(err.Error(), "err", err, "attempt", attempt, "txAttempts", etx.TxAttempts)
			return err
		}
		replacementAttempt, err := ec.bumpGas(ctx, etx, etx.TxAttempts)
		if err != nil {
			return fmt.Errorf("could not bump gas for terminally underpriced transaction: %w", err)
		}
		promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
		lggr.With(
			"sendError", sendError,
			"maxGasPriceConfig", ec.feeConfig.MaxFeePrice(),
			"previousAttempt", attempt,
			"replacementAttempt", replacementAttempt,
		).Errorf("gas price was rejected by the node for being too low. Node returned: '%s'", sendError.Error())

		if err := ec.txStore.SaveReplacementInProgressAttempt(ctx, attempt, &replacementAttempt); err != nil {
			return fmt.Errorf("saveReplacementInProgressAttempt failed: %w", err)
		}
		return ec.handleInProgressAttempt(ctx, lggr, etx, replacementAttempt, blockHeight)
	case client.ExceedsMaxFee:
		// Confirmer: Note it is not guaranteed that all nodes share the same tx fee cap.
		// So it is very likely that this attempt was successful on another node since
		// it was already successfully broadcasted. So we assume it is successful and
		// warn the operator that the RPC node is misconfigured.
		// This failure scenario is a strong indication that the RPC node
		// is misconfigured. This is a critical error and should be resolved by the
		// node operator.
		// If there is only one RPC node, or all RPC nodes have the same
		// configured cap, this transaction will get stuck and keep repeating
		// forever until the issue is resolved.
		lggr.Criticalw(`RPC node rejected this tx as outside Fee Cap but it may have been accepted by another Node`, "attempt", attempt)
		timeout := ec.dbConfig.DefaultQueryTimeout()
		return ec.txStore.SaveSentAttempt(ctx, timeout, &attempt, now)
	case client.Fatal:
		// WARNING: This should never happen!
		// Should NEVER be fatal this is an invariant violation. The
		// Broadcaster can never create a TxAttempt that will
		// fatally error.
		lggr.Criticalw("Invariant violation: fatal error while re-attempting transaction",
			"fee", attempt.TxFee,
			"feeLimit", attempt.ChainSpecificFeeLimit,
			"callerProvidedFeeLimit", etx.FeeLimit,
			"signedRawTx", commonhex.EnsurePrefix(hex.EncodeToString(attempt.SignedRawTx)),
			"blockHeight", blockHeight,
		)
		ec.SvcErrBuffer.Append(sendError)
		// This will loop continuously on every new head so it must be handled manually by the node operator!
		return ec.txStore.DeleteInProgressAttempt(ctx, attempt)
	case client.TerminallyStuck:
		// A transaction could broadcast successfully but then be considered terminally stuck on another attempt
		// Even though the transaction can succeed under different circumstances, we want to purge this transaction as soon as we get this error
		lggr.Warnw("terminally stuck transaction detected", "err", sendError.Error())
		ec.SvcErrBuffer.Append(sendError)
		// Create a purge attempt for tx
		purgeAttempt, err := ec.TxAttemptBuilder.NewPurgeTxAttempt(ctx, etx, lggr)
		if err != nil {
			return fmt.Errorf("NewPurgeTxAttempt failed: %w", err)
		}
		// Replace the in progress attempt with the purge attempt
		if err := ec.txStore.SaveReplacementInProgressAttempt(ctx, attempt, &purgeAttempt); err != nil {
			return fmt.Errorf("saveReplacementInProgressAttempt failed: %w", err)
		}
		return ec.handleInProgressAttempt(ctx, lggr, etx, purgeAttempt, blockHeight)
	case client.TransactionAlreadyKnown:
		// Sequence too low indicated that a transaction at this sequence was confirmed already.
		// Mark confirmed_missing_receipt and wait for the next cycle to try to get a receipt
		lggr.Debugw("Sequence already used", "txAttemptID", attempt.ID, "txHash", attempt.Hash.String())
		timeout := ec.dbConfig.DefaultQueryTimeout()
		return ec.txStore.SaveConfirmedAttempt(ctx, timeout, &attempt, now)
	case client.InsufficientFunds:
		timeout := ec.dbConfig.DefaultQueryTimeout()
		return ec.txStore.SaveInsufficientFundsAttempt(ctx, timeout, &attempt, now)
	case client.Successful:
		lggr.Debugw("Successfully broadcast transaction", "txAttemptID", attempt.ID, "txHash", attempt.Hash.String())
		timeout := ec.dbConfig.DefaultQueryTimeout()
		return ec.txStore.SaveSentAttempt(ctx, timeout, &attempt, now)
	case client.Unknown:
		// Every error that doesn't fall under one of the above categories will be treated as Unknown.
		fallthrough
	default:
		// Any other type of error is considered temporary or resolvable by the
		// node operator. The node may have it in the mempool so we must keep the
		// attempt (leave it in_progress). Safest thing to do is bail out and wait
		// for the next head.
		return fmt.Errorf("unexpected error sending tx %v with hash %s: %w", etx.ID, attempt.Hash.String(), sendError)
	}
}

// ForceRebroadcast sends a transaction for every sequence in the given sequence range at the given gas price.
// If an tx exists for this sequence, we re-send the existing tx with the supplied parameters.
// If an tx doesn't exist for this sequence, we send a zero transaction.
// This operates completely orthogonal to the normal Confirmer and can result in untracked attempts!
// Only for emergency usage.
// This is in case of some unforeseen scenario where the node is refusing to release the lock. KISS.
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ForceRebroadcast(ctx context.Context, seqs []SEQ, fee FEE, address ADDR, overrideGasLimit uint64) error {
	if len(seqs) == 0 {
		ec.lggr.Infof("ForceRebroadcast: No sequences provided. Skipping")
		return nil
	}
	ec.lggr.Infof("ForceRebroadcast: will rebroadcast transactions for all sequences between %v and %v", seqs[0], seqs[len(seqs)-1])

	for _, seq := range seqs {
		etx, err := ec.txStore.FindTxWithSequence(ctx, address, seq)
		if err != nil {
			return fmt.Errorf("ForceRebroadcast failed: %w", err)
		}
		if etx == nil {
			ec.lggr.Debugf("ForceRebroadcast: no tx found with sequence %s, will rebroadcast empty transaction", seq)
			hashStr, err := ec.sendEmptyTransaction(ctx, address, seq, overrideGasLimit, fee)
			if err != nil {
				ec.lggr.Errorw("ForceRebroadcast: failed to send empty transaction", "sequence", seq, "err", err)
				continue
			}
			ec.lggr.Infow("ForceRebroadcast: successfully rebroadcast empty transaction", "sequence", seq, "hash", hashStr)
		} else {
			ec.lggr.Debugf("ForceRebroadcast: got tx %v with sequence %v, will rebroadcast this transaction", etx.ID, *etx.Sequence)
			if overrideGasLimit != 0 {
				etx.FeeLimit = overrideGasLimit
			}
			attempt, _, err := ec.NewCustomTxAttempt(ctx, *etx, fee, etx.FeeLimit, 0x0, ec.lggr)
			if err != nil {
				ec.lggr.Errorw("ForceRebroadcast: failed to create new attempt", "txID", etx.ID, "err", err)
				continue
			}
			attempt.Tx = *etx // for logging
			ec.lggr.Debugw("Sending transaction", "txAttemptID", attempt.ID, "txHash", attempt.Hash, "err", err, "meta", etx.Meta, "feeLimit", attempt.ChainSpecificFeeLimit, "callerProvidedFeeLimit", etx.FeeLimit, "attempt", attempt)
			if errCode, err := ec.client.SendTransactionReturnCode(ctx, *etx, attempt, ec.lggr); errCode != client.Successful && err != nil {
				ec.lggr.Errorw(fmt.Sprintf("ForceRebroadcast: failed to rebroadcast tx %v with sequence %v, gas limit %v, and caller provided fee Limit %v	: %s", etx.ID, *etx.Sequence, attempt.ChainSpecificFeeLimit, etx.FeeLimit, err.Error()), "err", err, "fee", attempt.TxFee)
				continue
			}
			ec.lggr.Infof("ForceRebroadcast: successfully rebroadcast tx %v with hash: 0x%x", etx.ID, attempt.Hash)
		}
	}
	return nil
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) sendEmptyTransaction(ctx context.Context, fromAddress ADDR, seq SEQ, overrideGasLimit uint64, fee FEE) (string, error) {
	gasLimit := overrideGasLimit
	if gasLimit == 0 {
		gasLimit = ec.feeConfig.LimitDefault()
	}
	txhash, err := ec.client.SendEmptyTransaction(ctx, ec.TxAttemptBuilder.NewEmptyTxAttempt, seq, gasLimit, fee, fromAddress)
	if err != nil {
		return "", fmt.Errorf("(Confirmer).sendEmptyTransaction failed: %w", err)
	}
	return txhash, nil
}

// observeUntilTxConfirmed observes the promBlocksUntilTxConfirmed metric for each confirmed
// transaction.
func observeUntilTxConfirmed[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
](chainID CHAIN_ID, attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], head types.Head[BLOCK_HASH]) {
	for _, attempt := range attempts {
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
			blocksElapsed := head.BlockNumber() - broadcastBefore
			promBlocksUntilTxConfirmed.
				WithLabelValues(chainID.String()).
				Observe(float64(blocksElapsed))
		}
	}
}
