package txmgr

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// processHeadTimeout represents a sanity limit on how long ProcessHead
	// should take to complete
	processHeadTimeout = 10 * time.Minute

	// logAfterNConsecutiveBlocksChainTooShort logs a warning if we go at least
	// this many consecutive blocks with a re-org protection chain that is too
	// short
	//
	// we don't log every time because on startup it can be lower, only if it
	// persists does it indicate a serious problem
	logAfterNConsecutiveBlocksChainTooShort = 10
)

var (
	promNumGasBumps = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_gas_bumps",
		Help: "Number of gas bumps",
	}, []string{"evmChainID"})

	promGasBumpExceedsLimit = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_gas_bump_exceeds_limit",
		Help: "Number of times gas bumping failed from exceeding the configured limit. Any counts of this type indicate a serious problem.",
	}, []string{"evmChainID"})
	promNumConfirmedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_confirmed_transactions",
		Help: "Total number of confirmed transactions. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
	}, []string{"evmChainID"})
	promNumSuccessfulTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_successful_transactions",
		Help: "Total number of successful transactions. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
	}, []string{"evmChainID"})
	promRevertedTxCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_tx_reverted",
		Help: "Number of times a transaction reverted on-chain. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
	}, []string{"evmChainID"})
	promFwdTxCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_fwd_tx_count",
		Help: "The number of forwarded transaction attempts labeled by status",
	}, []string{"evmChainID", "successful"})
	promTxAttemptCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tx_manager_tx_attempt_count",
		Help: "The number of transaction attempts that are currently being processed by the transaction manager",
	}, []string{"evmChainID"})
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
	}, []string{"evmChainID"})
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
	}, []string{"evmChainID"})
)

// Confirmer is a broad service which performs four different tasks in sequence on every new longest chain
// Step 1: Mark that all currently pending transaction attempts were broadcast before this block
// Step 2: Check pending transactions for receipts
// Step 3: See if any transactions have exceeded the gas bumping block threshold and, if so, bump them
// Step 4: Check confirmed transactions to make sure they are still in the longest chain (reorg protection)
type Confirmer[
	CHAIN_ID txmgrtypes.ID,
	HEAD commontypes.Head[BLOCK_HASH],
	ADDR commontypes.Hashable,
	TX_HASH commontypes.Hashable,
	BLOCK_HASH commontypes.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ txmgrtypes.Sequence,
	FEE feetypes.Fee,
] struct {
	utils.StartStopOnce
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	lggr    logger.Logger
	client  txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	resumeCallback ResumeCallback
	config         txmgrtypes.ConfirmerConfig
	txConfig       txmgrtypes.ConfirmerTransactionsConfig
	dbConfig       txmgrtypes.ConfirmerDatabaseConfig
	chainID        CHAIN_ID

	ks               txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	enabledAddresses []ADDR

	mb        *utils.Mailbox[HEAD]
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup
	initSync  sync.Mutex
	isStarted bool

	nConsecutiveBlocksChainTooShort int
	isReceiptNil                    func(R) bool
}

func NewConfirmer[
	CHAIN_ID txmgrtypes.ID,
	HEAD commontypes.Head[BLOCK_HASH],
	ADDR commontypes.Hashable,
	TX_HASH commontypes.Hashable,
	BLOCK_HASH commontypes.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ txmgrtypes.Sequence,
	FEE feetypes.Fee,
](
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	client txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	config txmgrtypes.ConfirmerConfig,
	txConfig txmgrtypes.ConfirmerTransactionsConfig,
	dbConfig txmgrtypes.ConfirmerDatabaseConfig,
	keystore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	lggr logger.Logger,
	isReceiptNil func(R) bool,
) *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	lggr = lggr.Named("Confirmer")
	return &Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		txStore:          txStore,
		lggr:             lggr,
		client:           client,
		TxAttemptBuilder: txAttemptBuilder,
		resumeCallback:   nil,
		config:           config,
		txConfig:         txConfig,
		dbConfig:         dbConfig,
		chainID:          client.ConfiguredChainID(),
		ks:               keystore,
		mb:               utils.NewSingleMailbox[HEAD](),
		isReceiptNil:     isReceiptNil,
	}
}

// Start is a comment to appease the linter
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Start(_ context.Context) error {
	return ec.StartOnce("Confirmer", func() error {
		if ec.config.FeeBumpThreshold() == 0 {
			ec.lggr.Infow("Gas bumping is disabled (EVM.GasEstimator.BumpThreshold set to 0)", "ethGasBumpThreshold", 0)
		} else {
			ec.lggr.Infow(fmt.Sprintf("Gas bumping is enabled, unconfirmed transactions will have their gas price bumped every %d blocks", ec.config.FeeBumpThreshold()), "ethGasBumpThreshold", ec.config.FeeBumpThreshold())
		}

		return ec.startInternal()
	})
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) startInternal() error {
	ec.initSync.Lock()
	defer ec.initSync.Unlock()
	if ec.isStarted {
		return errors.New("Confirmer is already started")
	}
	var err error
	ec.enabledAddresses, err = ec.ks.EnabledAddressesForChain(ec.chainID)
	if err != nil {
		return errors.Wrap(err, "Confirmer: failed to load EnabledAddressesForChain")
	}

	ec.ctx, ec.ctxCancel = context.WithCancel(context.Background())
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
		return errors.Wrap(utils.ErrAlreadyStopped, "Confirmer is not started")
	}
	ec.ctxCancel()
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
	return map[string]error{ec.Name(): ec.StartStopOnce.Healthy()}
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) runLoop() {
	defer ec.wg.Done()
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
		case <-ec.ctx.Done():
			return
		}
	}
}

// ProcessHead takes all required transactions for the confirmer on a new head
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ProcessHead(ctx context.Context, head commontypes.Head[BLOCK_HASH]) error {
	ctx, cancel := context.WithTimeout(ctx, processHeadTimeout)
	defer cancel()
	return ec.processHead(ctx, head)
}

// NOTE: This SHOULD NOT be run concurrently or it could behave badly
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) processHead(ctx context.Context, head commontypes.Head[BLOCK_HASH]) error {
	mark := time.Now()

	ec.lggr.Debugw("processHead start", "headNum", head.BlockNumber(), "id", "eth_confirmer")

	if err := ec.txStore.SetBroadcastBeforeBlockNum(head.BlockNumber(), ec.chainID); err != nil {
		return errors.Wrap(err, "SetBroadcastBeforeBlockNum failed")
	}
	if err := ec.CheckConfirmedMissingReceipt(ctx); err != nil {
		return errors.Wrap(err, "CheckConfirmedMissingReceipt failed")
	}

	if err := ec.CheckForReceipts(ctx, head.BlockNumber()); err != nil {
		return errors.Wrap(err, "CheckForReceipts failed")
	}

	ec.lggr.Debugw("Finished CheckForReceipts", "headNum", head.BlockNumber(), "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	if err := ec.RebroadcastWhereNecessary(ctx, head.BlockNumber()); err != nil {
		return errors.Wrap(err, "RebroadcastWhereNecessary failed")
	}

	ec.lggr.Debugw("Finished RebroadcastWhereNecessary", "headNum", head.BlockNumber(), "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	if err := ec.EnsureConfirmedTransactionsInLongestChain(ctx, head); err != nil {
		return errors.Wrap(err, "EnsureConfirmedTransactionsInLongestChain failed")
	}

	ec.lggr.Debugw("Finished EnsureConfirmedTransactionsInLongestChain", "headNum", head.BlockNumber(), "time", time.Since(mark), "id", "eth_confirmer")

	if ec.resumeCallback != nil {
		mark = time.Now()
		if err := ec.ResumePendingTaskRuns(ctx, head); err != nil {
			return errors.Wrap(err, "ResumePendingTaskRuns failed")
		}

		ec.lggr.Debugw("Finished ResumePendingTaskRuns", "headNum", head.BlockNumber(), "time", time.Since(mark), "id", "eth_confirmer")
	}

	ec.lggr.Debugw("processHead finish", "headNum", head.BlockNumber(), "id", "eth_confirmer")

	return nil
}

// CheckConfirmedMissingReceipt will attempt to re-send any transaction in the
// state of "confirmed_missing_receipt". If we get back any type of senderror
// other than "sequence too low" it means that this transaction isn't actually
// confirmed and needs to be put back into "unconfirmed" state, so that it can enter
// the gas bumping cycle. This is necessary in rare cases (e.g. Polygon) where
// network conditions are extremely hostile.
//
// For example, assume the following scenario:
//
// 0. We are connected to multiple primary nodes via load balancer
// 1. We send a transaction, it is confirmed and, we get a receipt
// 2. A new head comes in from RPC node 1 indicating that this transaction was re-org'd, so we put it back into unconfirmed state
// 3. We re-send that transaction to a RPC node 2 **which hasn't caught up to this re-org yet**
// 4. RPC node 2 still has an old view of the chain, so it returns us "sequence too low" indicating "no problem this transaction is already mined"
// 5. Now the transaction is marked "confirmed_missing_receipt" but the latest chain does not actually include it
// 6. Now we are reliant on the EthResender to propagate it, and this transaction will not be gas bumped, so in the event of gas spikes it could languish or even be evicted from the mempool and hold up the queue
// 7. Even if/when RPC node 2 catches up, the transaction is still stuck in state "confirmed_missing_receipt"
//
// This scenario might sound unlikely but has been observed to happen multiple times in the wild on Polygon.
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CheckConfirmedMissingReceipt(ctx context.Context) (err error) {
	attempts, err := ec.txStore.FindTxAttemptsConfirmedMissingReceipt(ec.chainID)
	if err != nil {
		return err
	}
	if len(attempts) == 0 {
		return nil
	}
	ec.lggr.Infow(fmt.Sprintf("Found %d transactions confirmed_missing_receipt. The RPC node did not give us a receipt for these transactions even though it should have been mined. This could be due to using the wallet with an external account, or if the primary node is not synced or not propagating transactions properly", len(attempts)), "attempts", attempts)
	txCodes, txErrs, err := ec.client.BatchSendTransactions(ctx, ec.txStore.UpdateBroadcastAts, attempts, int(ec.config.RPCDefaultBatchSize()), ec.lggr)
	if err != nil {
		ec.lggr.Debugw("Batch sending transactions failed", err)
	}
	var ethTxIDsToUnconfirm []int64
	for idx, txErr := range txErrs {
		// Add to Unconfirm array, all tx where error wasn't TransactionAlreadyKnown.
		if txErr != nil {
			if txCodes[idx] == clienttypes.TransactionAlreadyKnown {
				continue
			}
		}

		ethTxIDsToUnconfirm = append(ethTxIDsToUnconfirm, attempts[idx].TxID)
	}
	err = ec.txStore.UpdateTxsUnconfirmed(ethTxIDsToUnconfirm)

	if err != nil {
		return err
	}
	return
}

// CheckForReceipts finds attempts that are still pending and checks to see if a receipt is present for the given block number
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CheckForReceipts(ctx context.Context, blockNum int64) error {
	attempts, err := ec.txStore.FindTxAttemptsRequiringReceiptFetch(ec.chainID)
	if err != nil {
		return errors.Wrap(err, "FindTxAttemptsRequiringReceiptFetch failed")
	}
	if len(attempts) == 0 {
		return nil
	}

	ec.lggr.Debugw(fmt.Sprintf("Fetching receipts for %v transaction attempts", len(attempts)), "blockNum", blockNum)

	attemptsByAddress := make(map[ADDR][]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
	for _, att := range attempts {
		attemptsByAddress[att.Tx.FromAddress] = append(attemptsByAddress[att.Tx.FromAddress], att)
	}

	for from, attempts := range attemptsByAddress {
		minedSequence, err := ec.getMinedSequenceForAddress(ctx, from)
		if err != nil {
			return errors.Wrapf(err, "unable to fetch pending sequence for address: %v", from)
		}

		// separateLikelyConfirmedAttempts is used as an optimisation: there is
		// no point trying to fetch receipts for attempts with a sequence higher
		// than the highest sequence the RPC node thinks it has seen
		likelyConfirmed := ec.separateLikelyConfirmedAttempts(from, attempts, minedSequence)
		likelyConfirmedCount := len(likelyConfirmed)
		if likelyConfirmedCount > 0 {
			likelyUnconfirmedCount := len(attempts) - likelyConfirmedCount

			ec.lggr.Debugf("Fetching and saving %v likely confirmed receipts. Skipping checking the others (%v)",
				likelyConfirmedCount, likelyUnconfirmedCount)

			start := time.Now()
			err = ec.fetchAndSaveReceipts(ctx, likelyConfirmed, blockNum)
			if err != nil {
				return errors.Wrapf(err, "unable to fetch and save receipts for likely confirmed txs, for address: %v", from)
			}
			ec.lggr.Debugw(fmt.Sprintf("Fetching and saving %v likely confirmed receipts done", likelyConfirmedCount),
				"time", time.Since(start))
		}
	}

	if err := ec.txStore.MarkAllConfirmedMissingReceipt(ec.chainID); err != nil {
		return errors.Wrap(err, "unable to mark eth_txes as 'confirmed_missing_receipt'")
	}

	if err := ec.txStore.MarkOldTxesMissingReceiptAsErrored(blockNum, ec.config.FinalityDepth(), ec.chainID); err != nil {
		return errors.Wrap(err, "unable to confirm buried unconfirmed eth_txes")
	}
	return nil
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) separateLikelyConfirmedAttempts(from ADDR, attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], minedSequence SEQ) []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	if len(attempts) == 0 {
		return attempts
	}

	firstAttemptSequence := *attempts[len(attempts)-1].Tx.Sequence
	lastAttemptSequence := *attempts[0].Tx.Sequence
	latestMinedSequence := minedSequence.Int64() - 1 // this can be -1 if a transaction has never been mined on this account
	ec.lggr.Debugw(fmt.Sprintf("There are %d attempts from address %s, mined transaction count is %d (latest mined sequence is %d) and for the attempts' sequences: first = %d, last = %d",
		len(attempts), from, minedSequence.Int64(), latestMinedSequence, firstAttemptSequence.Int64(), lastAttemptSequence.Int64()), "nAttempts", len(attempts), "fromAddress", from, "minedSequence", minedSequence, "latestMinedSequence", latestMinedSequence, "firstAttemptSequence", firstAttemptSequence, "lastAttemptSequence", lastAttemptSequence)

	likelyConfirmed := attempts
	// attempts are ordered by sequence ASC
	for i := 0; i < len(attempts); i++ {
		// If the attempt sequence is lower or equal to the latestBlockSequence
		// it must have been confirmed, we just didn't get a receipt yet
		//
		// Examples:
		// 3 transactions confirmed, highest has sequence 2
		// 5 total attempts, highest has sequence 4
		// minedSequence=3
		// likelyConfirmed will be attempts[0:3] which gives the first 3 transactions, as expected
		if (*attempts[i].Tx.Sequence).Int64() > minedSequence.Int64() {
			ec.lggr.Debugf("Marking attempts as likely confirmed just before index %v, at sequence: %v", i, *attempts[i].Tx.Sequence)
			likelyConfirmed = attempts[0:i]
			break
		}
	}

	if len(likelyConfirmed) == 0 {
		ec.lggr.Debug("There are no likely confirmed attempts - so will skip checking any")
	}

	return likelyConfirmed
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) fetchAndSaveReceipts(ctx context.Context, attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], blockNum int64) error {
	promTxAttemptCount.WithLabelValues(ec.chainID.String()).Set(float64(len(attempts)))

	batchSize := int(ec.config.RPCDefaultBatchSize())
	if batchSize == 0 {
		batchSize = len(attempts)
	}
	var allReceipts []R
	for i := 0; i < len(attempts); i += batchSize {
		j := i + batchSize
		if j > len(attempts) {
			j = len(attempts)
		}

		ec.lggr.Debugw(fmt.Sprintf("Batch fetching receipts at indexes %v until (excluded) %v", i, j), "blockNum", blockNum)

		batch := attempts[i:j]

		receipts, err := ec.batchFetchReceipts(ctx, batch, blockNum)
		if err != nil {
			return errors.Wrap(err, "batchFetchReceipts failed")
		}
		if err := ec.txStore.SaveFetchedReceipts(receipts, ec.chainID); err != nil {
			return errors.Wrap(err, "saveFetchedReceipts failed")
		}
		promNumConfirmedTxs.WithLabelValues(ec.chainID.String()).Add(float64(len(receipts)))

		allReceipts = append(allReceipts, receipts...)
	}

	observeUntilTxConfirmed(ec.chainID, attempts, allReceipts)

	return nil
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) getMinedSequenceForAddress(ctx context.Context, from ADDR) (SEQ, error) {
	return ec.client.SequenceAt(ctx, from, nil)
}

// Note this function will increment promRevertedTxCount upon receiving
// a reverted transaction receipt. Should only be called with unconfirmed attempts.
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) batchFetchReceipts(ctx context.Context, attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], blockNum int64) (receipts []R, err error) {
	// Metadata is required to determine whether a tx is forwarded or not.
	if ec.txConfig.ForwardersEnabled() {
		err = ec.txStore.PreloadTxes(attempts)
		if err != nil {
			return nil, errors.Wrap(err, "Confirmer#batchFetchReceipts error loading txs for attempts")
		}
	}

	lggr := ec.lggr.Named("batchFetchReceipts").With("blockNum", blockNum)

	txReceipts, txErrs, err := ec.client.BatchGetReceipts(ctx, attempts)
	if err != nil {
		return nil, err
	}

	for i := range txReceipts {
		attempt := attempts[i]
		receipt := txReceipts[i]
		err := txErrs[i]

		l := logger.Sugared(attempt.Tx.GetLogger(lggr).With(
			"txHash", attempt.Hash.String(), "ethTxAttemptID", attempt.ID,
			"ethTxID", attempt.TxID, "err", err, "sequence", attempt.Tx.Sequence,
		))

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

		l = logger.Sugared(l.With("blockHash", receipt.GetBlockHash().String(), "status", receipt.GetStatus(), "transactionIndex", receipt.GetTransactionIndex()))

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

		receipts = append(receipts, receipt)
	}

	return
}

// RebroadcastWhereNecessary bumps gas or resends transactions that were previously out-of-eth
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
				ec.lggr.Errorw("Error in RebroadcastWhereNecessary", "error", err, "fromAddress", fromAddress)
			}

			wg.Done()
		}(address)
	}

	wg.Wait()

	return multierr.Combine(errors...)
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) rebroadcastWhereNecessary(ctx context.Context, address ADDR, blockHeight int64) error {
	if err := ec.handleAnyInProgressAttempts(ctx, address, blockHeight); err != nil {
		return errors.Wrap(err, "handleAnyInProgressAttempts failed")
	}

	threshold := int64(ec.config.FeeBumpThreshold())
	bumpDepth := int64(ec.config.FeeBumpTxDepth())
	maxInFlightTransactions := ec.txConfig.MaxInFlight()
	etxs, err := ec.FindTxsRequiringRebroadcast(ctx, ec.lggr, address, blockHeight, threshold, bumpDepth, maxInFlightTransactions, ec.chainID)
	if err != nil {
		return errors.Wrap(err, "FindTxsRequiringRebroadcast failed")
	}
	for _, etx := range etxs {
		lggr := etx.GetLogger(ec.lggr)

		attempt, err := ec.attemptForRebroadcast(ctx, lggr, *etx)
		if err != nil {
			return errors.Wrap(err, "attemptForRebroadcast failed")
		}

		lggr.Debugw("Rebroadcasting transaction", "nPreviousAttempts", len(etx.TxAttempts), "fee", attempt.TxFee)

		if err := ec.txStore.SaveInProgressAttempt(&attempt); err != nil {
			return errors.Wrap(err, "saveInProgressAttempt failed")
		}

		if err := ec.handleInProgressAttempt(ctx, lggr, *etx, attempt, blockHeight); err != nil {
			return errors.Wrap(err, "handleInProgressAttempt failed")
		}
	}
	return nil
}

// "in_progress" attempts were left behind after a crash/restart and may or may not have been sent.
// We should try to ensure they get on-chain so we can fetch a receipt for them.
// NOTE: We also use this to mark attempts for rebroadcast in event of a
// re-org, so multiple attempts are allowed to be in in_progress state (but
// only one per eth_tx).
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleAnyInProgressAttempts(ctx context.Context, address ADDR, blockHeight int64) error {
	attempts, err := ec.txStore.GetInProgressTxAttempts(ctx, address, ec.chainID)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "GetInProgressTxAttempts failed")
	}
	for _, a := range attempts {
		err := ec.handleInProgressAttempt(ctx, a.Tx.GetLogger(ec.lggr), a.Tx, a, blockHeight)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			return errors.Wrap(err, "handleInProgressAttempt failed")
		}
	}
	return nil
}

// FindTxsRequiringRebroadcast returns attempts that hit insufficient eth,
// and attempts that need bumping, in sequence ASC order
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxsRequiringRebroadcast(ctx context.Context, lggr logger.Logger, address ADDR, blockNum, gasBumpThreshold, bumpDepth int64, maxInFlightTransactions uint32, chainID CHAIN_ID) (etxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	// NOTE: These two queries could be combined into one using union but it
	// becomes harder to read and difficult to test in isolation. KISS principle
	etxInsufficientEths, err := ec.txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(address, chainID, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}

	if len(etxInsufficientEths) > 0 {
		lggr.Infow(fmt.Sprintf("Found %d transactions to be re-sent that were previously rejected due to insufficient eth balance", len(etxInsufficientEths)), "blockNum", blockNum, "address", address)
	}

	// TODO: Just pass the Q through everything
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
			logger.Sugared(lggr).AssumptionViolationf("Expected eth_tx for gas bump to have at least one attempt", "etxID", etx.ID, "blockNum", blockNum, "address", address)
		}
		lggr.Infow(fmt.Sprintf("Found %d transactions to re-sent that have still not been confirmed after at least %d blocks. The oldest of these has not still not been confirmed after %d blocks. These transactions will have their gas price bumped. %s", len(etxBumps), gasBumpThreshold, oldestBlocksBehind, label.NodeConnectivityProblemWarning), "blockNum", blockNum, "address", address, "gasBumpThreshold", gasBumpThreshold)
	}

	seen := make(map[int64]struct{})

	for _, etx := range etxInsufficientEths {
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
		if previousAttempt.State == txmgrtypes.TxAttemptInsufficientEth {
			// Do not create a new attempt if we ran out of eth last time since bumping gas is pointless
			// Instead try to resubmit the same attempt at the same price, in the hope that the wallet was funded since our last attempt
			lggr.Debugw("Rebroadcast InsufficientEth", logFields...)
			previousAttempt.State = txmgrtypes.TxAttemptInProgress
			return previousAttempt, nil
		}
		attempt, err = ec.bumpGas(ctx, etx, etx.TxAttempts)

		if gas.IsBumpErr(err) {
			lggr.Errorw("Failed to bump gas", append(logFields, "err", err)...)
			// Do not create a new attempt if bumping gas would put us over the limit or cause some other problem
			// Instead try to resubmit the previous attempt, and keep resubmitting until its accepted
			previousAttempt.BroadcastBeforeBlockNum = nil
			previousAttempt.State = txmgrtypes.TxAttemptInProgress
			return previousAttempt, nil
		}
		return attempt, err
	}
	return attempt, errors.Errorf("invariant violation: Tx %v was unconfirmed but didn't have any attempts. "+
		"Falling back to default gas price instead."+
		"This is a bug! Please report to https://github.com/smartcontractkit/chainlink/issues", etx.ID)
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) logFieldsPreviousAttempt(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) []interface{} {
	etx := attempt.Tx
	return []interface{}{
		"etxID", etx.ID,
		"txHash", attempt.Hash,
		"previousAttempt", attempt,
		"feeLimit", etx.FeeLimit,
		"maxGasPrice", ec.config.MaxFeePrice(),
		"sequence", etx.Sequence,
	}
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) bumpGas(ctx context.Context, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], previousAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) (bumpedAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	previousAttempt := previousAttempts[0]
	logFields := ec.logFieldsPreviousAttempt(previousAttempt)

	var bumpedFee FEE
	var bumpedFeeLimit uint32
	bumpedAttempt, bumpedFee, bumpedFeeLimit, _, err = ec.NewBumpTxAttempt(ctx, etx, previousAttempt, previousAttempts, ec.lggr)

	// if no error, return attempt
	// if err, continue below
	if err == nil {
		promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
		ec.lggr.Debugw("Rebroadcast bumping fee for tx", append(logFields, "bumpedFee", bumpedFee.String(), "bumpedFeeLimit", bumpedFeeLimit)...)
		return bumpedAttempt, err
	}

	if errors.Is(errors.Cause(err), gas.ErrBumpGasExceedsLimit) {
		promGasBumpExceedsLimit.WithLabelValues(ec.chainID.String()).Inc()
	}

	return bumpedAttempt, errors.Wrap(err, "error bumping gas")
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleInProgressAttempt(ctx context.Context, lggr logger.Logger, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], blockHeight int64) error {
	if attempt.State != txmgrtypes.TxAttemptInProgress {

		return errors.Errorf("invariant violation: expected eth_tx_attempt %v to be in_progress, it was %s", attempt.ID, attempt.State)
	}

	now := time.Now()
	lggr.Debugw("Sending transaction", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash, "meta", etx.Meta, "feeLimit", etx.FeeLimit, "attempt", attempt, "etx", etx)
	errType, sendError := ec.client.SendTransactionReturnCode(ctx, etx, attempt, lggr)

	switch errType {
	case clienttypes.Underpriced:
		// This should really not ever happen in normal operation since we
		// already bumped above the required minimum in ethBroadcaster.
		ec.lggr.Warnw("Got terminally underpriced error for gas bump, this should never happen unless the remote RPC node changed its configuration on the fly, or you are using multiple RPC nodes with different minimum gas price requirements. This is not recommended", "err", sendError, "attempt", attempt)
		// "Lazily" load attempts here since the overwhelmingly common case is
		// that we don't need them unless we enter this path
		if err := ec.txStore.LoadTxAttempts(&etx, pg.WithParentCtx(ctx)); err != nil {
			return errors.Wrap(err, "failed to load TxAttempts while bumping on terminally underpriced error")
		}
		if len(etx.TxAttempts) == 0 {
			err := errors.New("expected to find at least 1 attempt")
			logger.Sugared(ec.lggr).AssumptionViolationw(err.Error(), "err", err, "attempt", attempt)
			return err
		}
		if attempt.ID != etx.TxAttempts[0].ID {
			err := errors.New("expected highest priced attempt to be the current in_progress attempt")
			logger.Sugared(ec.lggr).AssumptionViolationw(err.Error(), "err", err, "attempt", attempt, "ethTxAttempts", etx.TxAttempts)
			return err
		}
		replacementAttempt, err := ec.bumpGas(ctx, etx, etx.TxAttempts)
		if err != nil {
			return errors.Wrap(err, "could not bump gas for terminally underpriced transaction")
		}
		promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
		lggr.With(
			"sendError", sendError,
			"maxGasPriceConfig", ec.config.MaxFeePrice(),
			"previousAttempt", attempt,
			"replacementAttempt", replacementAttempt,
		).Errorf("gas price was rejected by the eth node for being too low. Eth node returned: '%s'", sendError.Error())

		if err := ec.txStore.SaveReplacementInProgressAttempt(attempt, &replacementAttempt); err != nil {
			return errors.Wrap(err, "saveReplacementInProgressAttempt failed")
		}
		return ec.handleInProgressAttempt(ctx, lggr, etx, replacementAttempt, blockHeight)
	case clienttypes.ExceedsMaxFee:
		// Confirmer: The gas price was bumped too high. This transaction attempt cannot be accepted.
		// Best thing we can do is to re-send the previous attempt at the old
		// price and discard this bumped version.
		fallthrough
	case clienttypes.Fatal:
		// WARNING: This should never happen!
		// Should NEVER be fatal this is an invariant violation. The
		// Broadcaster can never create a TxAttempt that will
		// fatally error.
		lggr.Criticalw("Invariant violation: fatal error while re-attempting transaction",
			"err", sendError,
			"fee", attempt.TxFee,
			"feeLimit", etx.FeeLimit,
			"signedRawTx", hexutil.Encode(attempt.SignedRawTx),
			"blockHeight", blockHeight,
		)
		ec.SvcErrBuffer.Append(sendError)
		// This will loop continuously on every new head so it must be handled manually by the node operator!
		return ec.txStore.DeleteInProgressAttempt(ctx, attempt)
	case clienttypes.TransactionAlreadyKnown:
		// Sequence too low indicated that a transaction at this sequence was confirmed already.
		// Mark confirmed_missing_receipt and wait for the next cycle to try to get a receipt
		lggr.Debugw("Sequence already used", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash.String(), "err", sendError)
		timeout := ec.dbConfig.DefaultQueryTimeout()
		return ec.txStore.SaveConfirmedMissingReceiptAttempt(ctx, timeout, &attempt, now)
	case clienttypes.InsufficientFunds:
		timeout := ec.dbConfig.DefaultQueryTimeout()
		return ec.txStore.SaveInsufficientFundsAttempt(timeout, &attempt, now)
	case clienttypes.Successful:
		lggr.Debugw("Successfully broadcast transaction", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash.String())
		timeout := ec.dbConfig.DefaultQueryTimeout()
		return ec.txStore.SaveSentAttempt(timeout, &attempt, now)
	case clienttypes.Unknown:
		// Every error that doesn't fall under one of the above categories will be treated as Unknown.
		fallthrough
	default:
		// Any other type of error is considered temporary or resolvable by the
		// node operator. The node may have it in the mempool so we must keep the
		// attempt (leave it in_progress). Safest thing to do is bail out and wait
		// for the next head.
		return errors.Wrapf(sendError, "unexpected error sending eth_tx %v with hash %s", etx.ID, attempt.Hash.String())
	}
}

// EnsureConfirmedTransactionsInLongestChain finds all confirmed eth_txes up to the depth
// of the given chain and ensures that every one has a receipt with a block hash that is
// in the given chain.
//
// If any of the confirmed transactions does not have a receipt in the chain, it has been
// re-org'd out and will be rebroadcast.
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) EnsureConfirmedTransactionsInLongestChain(ctx context.Context, head commontypes.Head[BLOCK_HASH]) error {
	if head.ChainLength() < ec.config.FinalityDepth() {
		logArgs := []interface{}{
			"chainLength", head.ChainLength(), "evmFinalityDepth", ec.config.FinalityDepth(),
		}
		if ec.nConsecutiveBlocksChainTooShort > logAfterNConsecutiveBlocksChainTooShort {
			warnMsg := "Chain length supplied for re-org detection was shorter than EvmFinalityDepth. Re-org protection is not working properly. This could indicate a problem with the remote RPC endpoint, a compatibility issue with a particular blockchain, a bug with this particular blockchain, heads table being truncated too early, remote node out of sync, or something else. If this happens a lot please raise a bug with the Chainlink team including a log output sample and details of the chain and RPC endpoint you are using."
			ec.lggr.Warnw(warnMsg, append(logArgs, "nConsecutiveBlocksChainTooShort", ec.nConsecutiveBlocksChainTooShort)...)
		} else {
			logMsg := "Chain length supplied for re-org detection was shorter than EvmFinalityDepth"
			ec.lggr.Debugw(logMsg, append(logArgs, "nConsecutiveBlocksChainTooShort", ec.nConsecutiveBlocksChainTooShort)...)
		}
		ec.nConsecutiveBlocksChainTooShort++
	} else {
		ec.nConsecutiveBlocksChainTooShort = 0
	}
	etxs, err := ec.txStore.FindTransactionsConfirmedInBlockRange(head.BlockNumber(), head.EarliestHeadInChain().BlockNumber(), ec.chainID)
	if err != nil {
		return errors.Wrap(err, "findTransactionsConfirmedInBlockRange failed")
	}

	for _, etx := range etxs {
		if !hasReceiptInLongestChain(*etx, head) {
			if err := ec.markForRebroadcast(*etx, head); err != nil {
				return errors.Wrapf(err, "markForRebroadcast failed for etx %v", etx.ID)
			}
		}
	}

	// It is safe to process separate keys concurrently
	// NOTE: This design will block one key if another takes a really long time to execute
	var wg sync.WaitGroup
	errors := []error{}
	var errMu sync.Mutex
	wg.Add(len(ec.enabledAddresses))
	for _, address := range ec.enabledAddresses {
		go func(fromAddress ADDR) {
			if err := ec.handleAnyInProgressAttempts(ctx, fromAddress, head.BlockNumber()); err != nil {
				errMu.Lock()
				errors = append(errors, err)
				errMu.Unlock()
				ec.lggr.Errorw("Error in handleAnyInProgressAttempts", "err", err, "fromAddress", fromAddress)
			}

			wg.Done()
		}(address)
	}

	wg.Wait()

	return multierr.Combine(errors...)
}

func hasReceiptInLongestChain[
	CHAIN_ID txmgrtypes.ID,
	ADDR commontypes.Hashable,
	TX_HASH, BLOCK_HASH commontypes.Hashable,
	SEQ txmgrtypes.Sequence,
	FEE feetypes.Fee,
](etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], head commontypes.Head[BLOCK_HASH]) bool {
	for {
		for _, attempt := range etx.TxAttempts {
			for _, receipt := range attempt.Receipts {
				if receipt.GetBlockHash().String() == head.BlockHash().String() && receipt.GetBlockNumber().Int64() == head.BlockNumber() {
					return true
				}
			}
		}
		if head.GetParent() == nil {
			return false
		}
		head = head.GetParent()
	}
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) markForRebroadcast(etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], head commontypes.Head[BLOCK_HASH]) error {
	if len(etx.TxAttempts) == 0 {
		return errors.Errorf("invariant violation: expected eth_tx %v to have at least one attempt", etx.ID)
	}

	// Rebroadcast the one with the highest gas price
	attempt := etx.TxAttempts[0]
	var receipt txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH]
	if len(attempt.Receipts) > 0 {
		receipt = attempt.Receipts[0]
	}

	logValues := []interface{}{
		"txhash", attempt.Hash.String(),
		"currentBlockNum", head.BlockNumber(),
		"currentBlockHash", head.BlockHash().String(),
		"ethTxID", etx.ID,
		"attemptID", attempt.ID,
		"nReceipts", len(attempt.Receipts),
		"id", "eth_confirmer",
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
	err := ec.txStore.UpdateTxForRebroadcast(etx, attempt)
	return errors.Wrap(err, "markForRebroadcast failed")
}

// ForceRebroadcast sends a transaction for every sequence in the given sequence range at the given gas price.
// If an eth_tx exists for this sequence, we re-send the existing eth_tx with the supplied parameters.
// If an eth_tx doesn't exist for this sequence, we send a zero transaction.
// This operates completely orthogonal to the normal Confirmer and can result in untracked attempts!
// Only for emergency usage.
// This is in case of some unforeseen scenario where the node is refusing to release the lock. KISS.
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ForceRebroadcast(seqs []SEQ, fee FEE, address ADDR, overrideGasLimit uint32) error {
	if len(seqs) == 0 {
		ec.lggr.Infof("ForceRebroadcast: No sequences provided. Skipping")
		return nil
	}
	ec.lggr.Infof("ForceRebroadcast: will rebroadcast transactions for all sequences between %v and %v", seqs[0], seqs[len(seqs)-1])

	for _, seq := range seqs {

		etx, err := ec.txStore.FindTxWithSequence(address, seq)
		if err != nil {
			return errors.Wrap(err, "ForceRebroadcast failed")
		}
		if etx == nil {
			ec.lggr.Debugf("ForceRebroadcast: no eth_tx found with sequence %s, will rebroadcast empty transaction", seq)
			hashStr, err := ec.sendEmptyTransaction(context.TODO(), address, seq, overrideGasLimit, fee)
			if err != nil {
				ec.lggr.Errorw("ForceRebroadcast: failed to send empty transaction", "sequence", seq, "err", err)
				continue
			}
			ec.lggr.Infow("ForceRebroadcast: successfully rebroadcast empty transaction", "sequence", seq, "hash", hashStr)
		} else {
			ec.lggr.Debugf("ForceRebroadcast: got eth_tx %v with sequence %v, will rebroadcast this transaction", etx.ID, *etx.Sequence)
			if overrideGasLimit != 0 {
				etx.FeeLimit = overrideGasLimit
			}
			attempt, _, err := ec.NewCustomTxAttempt(*etx, fee, etx.FeeLimit, 0x0, ec.lggr)
			if err != nil {
				ec.lggr.Errorw("ForceRebroadcast: failed to create new attempt", "ethTxID", etx.ID, "err", err)
				continue
			}
			attempt.Tx = *etx // for logging
			ec.lggr.Debugw("Sending transaction", "txAttemptID", attempt.ID, "txHash", attempt.Hash, "err", err, "meta", etx.Meta, "feeLimit", etx.FeeLimit, "attempt", attempt)
			if errCode, err := ec.client.SendTransactionReturnCode(context.TODO(), *etx, attempt, ec.lggr); errCode != clienttypes.Successful && err != nil {
				ec.lggr.Errorw(fmt.Sprintf("ForceRebroadcast: failed to rebroadcast eth_tx %v with sequence %v and gas limit %v: %s", etx.ID, *etx.Sequence, etx.FeeLimit, err.Error()), "err", err, "fee", attempt.TxFee)
				continue
			}
			ec.lggr.Infof("ForceRebroadcast: successfully rebroadcast eth_tx %v with hash: 0x%x", etx.ID, attempt.Hash)
		}
	}
	return nil
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) sendEmptyTransaction(ctx context.Context, fromAddress ADDR, seq SEQ, overrideGasLimit uint32, fee FEE) (string, error) {
	gasLimit := overrideGasLimit
	if gasLimit == 0 {
		gasLimit = ec.config.FeeLimitDefault()
	}
	txhash, err := ec.client.SendEmptyTransaction(ctx, ec.TxAttemptBuilder.NewEmptyTxAttempt, seq, gasLimit, fee, fromAddress)
	if err != nil {
		return "", errors.Wrap(err, "(Confirmer).sendEmptyTransaction failed")
	}
	return txhash, nil
}

// ResumePendingTaskRuns issues callbacks to task runs that are pending waiting for receipts
func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ResumePendingTaskRuns(ctx context.Context, head commontypes.Head[BLOCK_HASH]) error {

	receiptsPlus, err := ec.txStore.FindReceiptsPendingConfirmation(ctx, head.BlockNumber(), ec.chainID)

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
			taskErr = errors.Errorf("transaction %s reverted on-chain", data.Receipt.GetTxHash())
		} else {
			output = data.Receipt
		}

		ec.lggr.Debugw("Callback: resuming ethtx with receipt", "output", output, "taskErr", taskErr, "pipelineTaskRunID", data.ID)
		if err := ec.resumeCallback(data.ID, output, taskErr); err != nil {
			return err
		}
	}

	return nil
}

// observeUntilTxConfirmed observes the promBlocksUntilTxConfirmed metric for each confirmed
// transaction.
func observeUntilTxConfirmed[
	CHAIN_ID txmgrtypes.ID,
	ADDR commontypes.Hashable,
	TX_HASH, BLOCK_HASH commontypes.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ txmgrtypes.Sequence,
	FEE feetypes.Fee,
](chainID CHAIN_ID, attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], receipts []R) {
	for _, attempt := range attempts {
		for _, r := range receipts {
			if attempt.Hash.String() != r.GetTxHash().String() {
				continue
			}

			// We estimate the time until confirmation by subtracting from the time the eth tx (not the attempt)
			// was created. We want to measure the amount of time taken from when a transaction is created
			// via e.g Txm.CreateTransaction to when it is confirmed on-chain, regardless of how many attempts
			// were needed to achieve this.
			duration := time.Since(attempt.Tx.CreatedAt)
			promTimeUntilTxConfirmed.
				WithLabelValues(chainID.String()).
				Observe(float64(duration))

			// Since a eth tx can have many attempts, we take the number of blocks to confirm as the block number
			// of the receipt minus the block number of the first ever broadcast for this transaction.
			broadcastBefore := utils.MinKey(attempt.Tx.TxAttempts, func(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) int64 {
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
