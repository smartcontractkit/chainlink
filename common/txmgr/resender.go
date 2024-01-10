package txmgr

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink-common/pkg/chains/label"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	commonhex "github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/common/client"
	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

const (
	// pollInterval is the maximum amount of time in addition to
	// TxResendAfterThreshold that we will wait before resending an attempt
	DefaultResenderPollInterval = 5 * time.Second

	// Alert interval for unconfirmed transaction attempts
	unconfirmedTxAlertLogFrequency = 2 * time.Minute

	// timeout value for batchSendTransactions
	batchSendTransactionTimeout = 30 * time.Second
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
)

type Resender[
	CHAIN_ID types.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	services.StateMachine
	txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	txStore             txmgrtypes.TransactionStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, SEQ, FEE]
	client              txmgrtypes.TransactionClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	ks                  txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	chainID             CHAIN_ID
	interval            time.Duration
	config              txmgrtypes.ResenderChainConfig
	feeConfig           txmgrtypes.ResenderFeeConfig
	txConfig            txmgrtypes.ResenderTransactionsConfig
	dbConfig            txmgrtypes.ResenderDatabaseConfig
	logger              logger.SugaredLogger
	lastAlertTimestamps map[string]time.Time
	enabledAddresses    []ADDR

	mb        *mailbox.Mailbox[int64]
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup
}

func NewResender[
	CHAIN_ID types.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	lggr logger.Logger,
	txStore txmgrtypes.TransactionStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, SEQ, FEE],
	client txmgrtypes.TransactionClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	ks txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	pollInterval time.Duration,
	config txmgrtypes.ResenderChainConfig,
	feeConfig txmgrtypes.ResenderFeeConfig,
	txConfig txmgrtypes.ResenderTransactionsConfig,
	dbConfig txmgrtypes.ResenderDatabaseConfig,
) *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	if txConfig.ResendAfterThreshold() == 0 {
		panic("Resender requires a non-zero threshold")
	}
	return &Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		txStore:             txStore,
		client:              client,
		ks:                  ks,
		TxAttemptBuilder:    txAttemptBuilder,
		chainID:             client.ConfiguredChainID(),
		interval:            pollInterval,
		config:              config,
		feeConfig:           feeConfig,
		txConfig:            txConfig,
		dbConfig:            dbConfig,
		logger:              logger.Sugared(logger.Named(lggr, "Resender")),
		lastAlertTimestamps: make(map[string]time.Time),
		mb:                  mailbox.NewSingle[int64](),
	}
}

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Start(_ context.Context) error {
	return er.StartOnce("Resender", func() (err error) {
		er.enabledAddresses, err = er.ks.EnabledAddressesForChain(er.chainID)
		if err != nil {
			return fmt.Errorf("Resender: failed to load EnabledAddressesForChain: %w", err)
		}

		er.ctx, er.ctxCancel = context.WithCancel(context.Background())
		er.wg = sync.WaitGroup{}
		er.wg.Add(1)
		go er.runLoop()
		return nil
	})
}

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Close() error {
	return er.StopOnce("Resender", func() (err error) {
		er.ctxCancel()
		er.wg.Wait()
		return nil
	})
}

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) runLoop() {
	defer er.wg.Done()
	keysChanged, unsub := er.ks.SubscribeToKeyChanges()
	defer unsub()

	ticker := time.NewTicker(utils.WithJitter(er.interval))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := er.ResendUnconfirmed(er.ctx); err != nil {
				er.logger.Warnw("Failed to resend unconfirmed transactions", "err", err)
			}
		case <-er.mb.Notify():
			for {
				if er.ctx.Err() != nil {
					return
				}
				blockHeight, exists := er.mb.Retrieve()
				if !exists {
					break
				}
				if err := er.RebroadcastWhereNecessary(er.ctx, blockHeight); err != nil {
					er.logger.Errorw("Error processing head", "err", err)
					continue
				}
			}
		case <-keysChanged:
			var err error
			er.enabledAddresses, err = er.ks.EnabledAddressesForChain(er.chainID)
			if err != nil {
				er.logger.Critical("Failed to reload key states after key change")
				continue
			}
		case <-er.ctx.Done():
			return
		}
	}
}

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) ResendUnconfirmed(ctx context.Context) error {

	ageThreshold := er.txConfig.ResendAfterThreshold()
	maxInFlightTransactions := er.txConfig.MaxInFlight()
	olderThan := time.Now().Add(-ageThreshold)
	var allAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]

	for _, k := range er.enabledAddresses {
		attempts, err := er.txStore.FindTxAttemptsRequiringResend(ctx, olderThan, maxInFlightTransactions, er.chainID, k)
		if err != nil {
			return fmt.Errorf("failed to FindTxAttemptsRequiringResend: %w", err)
		}
		er.logStuckAttempts(attempts, k)

		allAttempts = append(allAttempts, attempts...)
	}

	if len(allAttempts) == 0 {
		for k := range er.lastAlertTimestamps {
			er.lastAlertTimestamps[k] = time.Now()
		}
		return nil
	}
	er.logger.Infow(fmt.Sprintf("Re-sending %d unconfirmed transactions that were last sent over %s ago. These transactions are taking longer than usual to be mined. %s", len(allAttempts), ageThreshold, label.NodeConnectivityProblemWarning), "n", len(allAttempts))

	batchSize := int(er.config.RPCDefaultBatchSize())
	batchCtx, cancel := context.WithTimeout(ctx, batchSendTransactionTimeout)
	defer cancel()
	txErrTypes, _, broadcastTime, txIDs, err := er.client.BatchSendTransactions(batchCtx, allAttempts, batchSize, er.logger)

	// update broadcast times before checking additional errors
	if len(txIDs) > 0 {
		if updateErr := er.txStore.UpdateBroadcastAts(batchCtx, broadcastTime, txIDs); updateErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to update broadcast time: %w", updateErr))
		}
	}
	if err != nil {
		return fmt.Errorf("failed to re-send transactions: %w", err)
	}
	logResendResult(er.logger, txErrTypes)

	return nil
}

func logResendResult(lggr logger.Logger, codes []client.SendTxReturnCode) {
	var nNew int
	var nFatal int
	for _, c := range codes {
		if c == client.Successful {
			nNew++
		} else if c == client.Fatal {
			nFatal++
		}
	}
	lggr.Debugw("Completed", "n", len(codes), "nNew", nNew, "nFatal", nFatal)
}

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) logStuckAttempts(attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], fromAddress ADDR) {
	if time.Since(er.lastAlertTimestamps[fromAddress.String()]) >= unconfirmedTxAlertLogFrequency {
		oldestAttempt, exists := findOldestUnconfirmedAttempt(attempts)
		if exists {
			// Wait at least 2 times the TxResendAfterThreshold to log critical with an unconfirmedTxAlertDelay
			if time.Since(oldestAttempt.CreatedAt) > er.txConfig.ResendAfterThreshold()*2 {
				er.lastAlertTimestamps[fromAddress.String()] = time.Now()
				er.logger.Errorw("TxAttempt has been unconfirmed for more than max duration", "maxDuration", er.txConfig.ResendAfterThreshold()*2,
					"txID", oldestAttempt.TxID, "txFee", oldestAttempt.TxFee,
					"BroadcastBeforeBlockNum", oldestAttempt.BroadcastBeforeBlockNum, "Hash", oldestAttempt.Hash, "fromAddress", fromAddress)
			}
		}
	}
}

func findOldestUnconfirmedAttempt[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
](attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) (txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], bool) {
	var oldestAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	if len(attempts) < 1 {
		return oldestAttempt, false
	}
	oldestAttempt = attempts[0]
	for i := 1; i < len(attempts); i++ {
		if oldestAttempt.CreatedAt.Sub(attempts[i].CreatedAt) <= 0 {
			oldestAttempt = attempts[i]
		}
	}
	return oldestAttempt, true
}

// RebroadcastWhereNecessary bumps gas or resends transactions that were previously out-of-funds
func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) RebroadcastWhereNecessary(ctx context.Context, blockHeight int64) error {
	var wg sync.WaitGroup

	if err := er.txStore.SetBroadcastBeforeBlockNum(ctx, blockHeight, er.chainID); err != nil {
		return fmt.Errorf("SetBroadcastBeforeBlockNum failed: %w", err)
	}

	// It is safe to process separate keys concurrently
	// NOTE: This design will block one key if another takes a really long time to execute
	wg.Add(len(er.enabledAddresses))
	errors := []error{}
	var errMu sync.Mutex
	for _, address := range er.enabledAddresses {
		go func(fromAddress ADDR) {
			if err := er.rebroadcastWhereNecessary(ctx, fromAddress, blockHeight); err != nil {
				errMu.Lock()
				errors = append(errors, err)
				errMu.Unlock()
				er.logger.Errorw("Error in RebroadcastWhereNecessary", "err", err, "fromAddress", fromAddress)
			}

			wg.Done()
		}(address)
	}

	wg.Wait()

	return multierr.Combine(errors...)
}

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) rebroadcastWhereNecessary(ctx context.Context, address ADDR, blockHeight int64) error {
	if err := er.handleAnyInProgressAttempts(ctx, address, blockHeight); err != nil {
		return fmt.Errorf("handleAnyInProgressAttempts failed: %w", err)
	}

	threshold := int64(er.feeConfig.BumpThreshold())
	bumpDepth := int64(er.feeConfig.BumpTxDepth())
	maxInFlightTransactions := er.txConfig.MaxInFlight()
	etxs, err := er.FindTxsRequiringRebroadcast(ctx, er.logger, address, blockHeight, threshold, bumpDepth, maxInFlightTransactions, er.chainID)
	if err != nil {
		return fmt.Errorf("FindTxsRequiringRebroadcast failed: %w", err)
	}
	for _, etx := range etxs {
		lggr := etx.GetLogger(er.logger)

		attempt, err := er.attemptForRebroadcast(ctx, lggr, *etx)
		if err != nil {
			return fmt.Errorf("attemptForRebroadcast failed: %w", err)
		}

		lggr.Debugw("Rebroadcasting transaction", "nPreviousAttempts", len(etx.TxAttempts), "fee", attempt.TxFee)

		if err := er.txStore.SaveInProgressAttempt(ctx, &attempt); err != nil {
			return fmt.Errorf("saveInProgressAttempt failed: %w", err)
		}

		if err := er.handleInProgressAttempt(ctx, lggr, *etx, attempt, blockHeight); err != nil {
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
func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) handleAnyInProgressAttempts(ctx context.Context, address ADDR, blockHeight int64) error {
	attempts, err := er.txStore.GetInProgressTxAttempts(ctx, address, er.chainID)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return fmt.Errorf("GetInProgressTxAttempts failed: %w", err)
	}
	for _, a := range attempts {
		err := er.handleInProgressAttempt(ctx, a.Tx.GetLogger(er.logger), a.Tx, a, blockHeight)
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
func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) FindTxsRequiringRebroadcast(ctx context.Context, lggr logger.Logger, address ADDR, blockNum, gasBumpThreshold, bumpDepth int64, maxInFlightTransactions uint32, chainID CHAIN_ID) (etxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	// NOTE: These two queries could be combined into one using union but it
	// becomes harder to read and difficult to test in isolation. KISS principle
	etxInsufficientFunds, err := er.txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, address, chainID)
	if err != nil {
		return nil, err
	}

	if len(etxInsufficientFunds) > 0 {
		lggr.Infow(fmt.Sprintf("Found %d transactions to be re-sent that were previously rejected due to insufficient native token balance", len(etxInsufficientFunds)), "blockNum", blockNum, "address", address)
	}

	// TODO: Just pass the Q through everything
	etxBumps, err := er.txStore.FindTxsRequiringGasBump(ctx, address, blockNum, gasBumpThreshold, bumpDepth, chainID)
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
			logger.Sugared(lggr).AssumptionViolationf("Expected tx for gas bump to have at least one attempt", "etxID", etx.ID, "blockNum", blockNum, "address", address)
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

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) attemptForRebroadcast(ctx context.Context, lggr logger.Logger, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) (attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	if len(etx.TxAttempts) > 0 {
		etx.TxAttempts[0].Tx = etx
		previousAttempt := etx.TxAttempts[0]
		logFields := er.logFieldsPreviousAttempt(previousAttempt)
		if previousAttempt.State == txmgrtypes.TxAttemptInsufficientFunds {
			// Do not create a new attempt if we ran out of funds last time since bumping gas is pointless
			// Instead try to resubmit the same attempt at the same price, in the hope that the wallet was funded since our last attempt
			lggr.Debugw("Rebroadcast InsufficientFunds", logFields...)
			previousAttempt.State = txmgrtypes.TxAttemptInProgress
			return previousAttempt, nil
		}
		attempt, err = er.bumpGas(ctx, etx, etx.TxAttempts)

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

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) logFieldsPreviousAttempt(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) []interface{} {
	etx := attempt.Tx
	return []interface{}{
		"etxID", etx.ID,
		"txHash", attempt.Hash,
		"previousAttempt", attempt,
		"feeLimit", etx.FeeLimit,
		"maxGasPrice", er.feeConfig.MaxFeePrice(),
		"sequence", etx.Sequence,
	}
}

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bumpGas(ctx context.Context, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], previousAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) (bumpedAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	previousAttempt := previousAttempts[0]
	logFields := er.logFieldsPreviousAttempt(previousAttempt)

	var bumpedFee FEE
	var bumpedFeeLimit uint32
	bumpedAttempt, bumpedFee, bumpedFeeLimit, _, err = er.NewBumpTxAttempt(ctx, etx, previousAttempt, previousAttempts, er.logger)

	// if no error, return attempt
	// if err, continue below
	if err == nil {
		promNumGasBumps.WithLabelValues(er.chainID.String()).Inc()
		er.logger.Debugw("Rebroadcast bumping fee for tx", append(logFields, "bumpedFee", bumpedFee.String(), "bumpedFeeLimit", bumpedFeeLimit)...)
		return bumpedAttempt, err
	}

	if errors.Is(err, commonfee.ErrBumpFeeExceedsLimit) {
		promGasBumpExceedsLimit.WithLabelValues(er.chainID.String()).Inc()
	}

	return bumpedAttempt, fmt.Errorf("error bumping gas: %w", err)
}

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) handleInProgressAttempt(ctx context.Context, lggr logger.SugaredLogger, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], blockHeight int64) error {
	if attempt.State != txmgrtypes.TxAttemptInProgress {

		return fmt.Errorf("invariant violation: expected tx_attempt %v to be in_progress, it was %s", attempt.ID, attempt.State)
	}

	now := time.Now()
	lggr.Debugw("Sending transaction", "txAttemptID", attempt.ID, "txHash", attempt.Hash, "meta", etx.Meta, "feeLimit", etx.FeeLimit, "attempt", attempt, "etx", etx)
	errType, sendError := er.client.SendTransactionReturnCode(ctx, etx, attempt, lggr)

	switch errType {
	case client.Underpriced:
		// This should really not ever happen in normal operation since we
		// already bumped above the required minimum in broadcaster.
		er.logger.Warnw("Got terminally underpriced error for gas bump, this should never happen unless the remote RPC node changed its configuration on the fly, or you are using multiple RPC nodes with different minimum gas price requirements. This is not recommended", "attempt", attempt)
		// "Lazily" load attempts here since the overwhelmingly common case is
		// that we don't need them unless we enter this path
		if err := er.txStore.LoadTxAttempts(ctx, &etx); err != nil {
			return fmt.Errorf("failed to load TxAttempts while bumping on terminally underpriced error: %w", err)
		}
		if len(etx.TxAttempts) == 0 {
			err := errors.New("expected to find at least 1 attempt")
			er.logger.AssumptionViolationw(err.Error(), "err", err, "attempt", attempt)
			return err
		}
		if attempt.ID != etx.TxAttempts[0].ID {
			err := errors.New("expected highest priced attempt to be the current in_progress attempt")
			er.logger.AssumptionViolationw(err.Error(), "err", err, "attempt", attempt, "txAttempts", etx.TxAttempts)
			return err
		}
		replacementAttempt, err := er.bumpGas(ctx, etx, etx.TxAttempts)
		if err != nil {
			return fmt.Errorf("could not bump gas for terminally underpriced transaction: %w", err)
		}
		promNumGasBumps.WithLabelValues(er.chainID.String()).Inc()
		lggr.With(
			"sendError", sendError,
			"maxGasPriceConfig", er.feeConfig.MaxFeePrice(),
			"previousAttempt", attempt,
			"replacementAttempt", replacementAttempt,
		).Errorf("gas price was rejected by the node for being too low. Node returned: '%s'", sendError.Error())

		if err := er.txStore.SaveReplacementInProgressAttempt(ctx, attempt, &replacementAttempt); err != nil {
			return fmt.Errorf("saveReplacementInProgressAttempt failed: %w", err)
		}
		return er.handleInProgressAttempt(ctx, lggr, etx, replacementAttempt, blockHeight)
	case client.ExceedsMaxFee:
		// Confirmer: The gas price was bumped too high. This transaction attempt cannot be accepted.
		// Best thing we can do is to re-send the previous attempt at the old
		// price and discard this bumped version.
		fallthrough
	case client.Fatal:
		// WARNING: This should never happen!
		// Should NEVER be fatal this is an invariant violation. The
		// Broadcaster can never create a TxAttempt that will
		// fatally error.
		lggr.Criticalw("Invariant violation: fatal error while re-attempting transaction",
			"fee", attempt.TxFee,
			"feeLimit", etx.FeeLimit,
			"signedRawTx", commonhex.EnsurePrefix(hex.EncodeToString(attempt.SignedRawTx)),
			"blockHeight", blockHeight,
		)
		er.SvcErrBuffer.Append(sendError)
		// This will loop continuously on every new head so it must be handled manually by the node operator!
		return er.txStore.DeleteInProgressAttempt(ctx, attempt)
	case client.TransactionAlreadyKnown:
		// Sequence too low indicated that a transaction at this sequence was confirmed already.
		// Mark confirmed_missing_receipt and wait for the next cycle to try to get a receipt
		lggr.Debugw("Sequence already used", "txAttemptID", attempt.ID, "txHash", attempt.Hash.String())
		timeout := er.dbConfig.DefaultQueryTimeout()
		return er.txStore.SaveConfirmedMissingReceiptAttempt(ctx, timeout, &attempt, now)
	case client.InsufficientFunds:
		timeout := er.dbConfig.DefaultQueryTimeout()
		return er.txStore.SaveInsufficientFundsAttempt(ctx, timeout, &attempt, now)
	case client.Successful:
		lggr.Debugw("Successfully broadcast transaction", "txAttemptID", attempt.ID, "txHash", attempt.Hash.String())
		timeout := er.dbConfig.DefaultQueryTimeout()
		return er.txStore.SaveSentAttempt(ctx, timeout, &attempt, now)
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
func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) ForceRebroadcast(ctx context.Context, seqs []SEQ, fee FEE, address ADDR, overrideGasLimit uint32) error {
	if len(seqs) == 0 {
		er.logger.Infof("ForceRebroadcast: No sequences provided. Skipping")
		return nil
	}
	er.logger.Infof("ForceRebroadcast: will rebroadcast transactions for all sequences between %v and %v", seqs[0], seqs[len(seqs)-1])

	for _, seq := range seqs {

		etx, err := er.txStore.FindTxWithSequence(ctx, address, seq)
		if err != nil {
			return fmt.Errorf("ForceRebroadcast failed: %w", err)
		}
		if etx == nil {
			er.logger.Debugf("ForceRebroadcast: no tx found with sequence %s, will rebroadcast empty transaction", seq)
			hashStr, err := er.sendEmptyTransaction(ctx, address, seq, overrideGasLimit, fee)
			if err != nil {
				er.logger.Errorw("ForceRebroadcast: failed to send empty transaction", "sequence", seq, "err", err)
				continue
			}
			er.logger.Infow("ForceRebroadcast: successfully rebroadcast empty transaction", "sequence", seq, "hash", hashStr)
		} else {
			er.logger.Debugf("ForceRebroadcast: got tx %v with sequence %v, will rebroadcast this transaction", etx.ID, *etx.Sequence)
			if overrideGasLimit != 0 {
				etx.FeeLimit = overrideGasLimit
			}
			attempt, _, err := er.NewCustomTxAttempt(*etx, fee, etx.FeeLimit, 0x0, er.logger)
			if err != nil {
				er.logger.Errorw("ForceRebroadcast: failed to create new attempt", "txID", etx.ID, "err", err)
				continue
			}
			attempt.Tx = *etx // for logging
			er.logger.Debugw("Sending transaction", "txAttemptID", attempt.ID, "txHash", attempt.Hash, "err", err, "meta", etx.Meta, "feeLimit", etx.FeeLimit, "attempt", attempt)
			if errCode, err := er.client.SendTransactionReturnCode(ctx, *etx, attempt, er.logger); errCode != client.Successful && err != nil {
				er.logger.Errorw(fmt.Sprintf("ForceRebroadcast: failed to rebroadcast tx %v with sequence %v and gas limit %v: %s", etx.ID, *etx.Sequence, etx.FeeLimit, err.Error()), "err", err, "fee", attempt.TxFee)
				continue
			}
			er.logger.Infof("ForceRebroadcast: successfully rebroadcast tx %v with hash: 0x%x", etx.ID, attempt.Hash)
		}
	}
	return nil
}

func (er *Resender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) sendEmptyTransaction(ctx context.Context, fromAddress ADDR, seq SEQ, overrideGasLimit uint32, fee FEE) (string, error) {
	gasLimit := overrideGasLimit
	if gasLimit == 0 {
		gasLimit = er.feeConfig.LimitDefault()
	}
	txhash, err := er.client.SendEmptyTransaction(ctx, er.TxAttemptBuilder.NewEmptyTxAttempt, seq, gasLimit, fee, fromAddress)
	if err != nil {
		return "", fmt.Errorf("(Resender).sendEmptyTransaction failed: %w", err)
	}
	return txhash, nil
}
