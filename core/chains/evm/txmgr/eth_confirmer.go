package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	// ErrCouldNotGetReceipt is the error string we save if we reach our finality depth for a confirmed transaction without ever getting a receipt
	// This most likely happened because an external wallet used the account for this nonce
	ErrCouldNotGetReceipt = "could not get receipt"

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

// EthConfirmer is a broad service which performs four different tasks in sequence on every new longest chain
// Step 1: Mark that all currently pending transaction attempts were broadcast before this block
// Step 2: Check pending transactions for receipts
// Step 3: See if any transactions have exceeded the gas bumping block threshold and, if so, bump them
// Step 4: Check confirmed transactions to make sure they are still in the longest chain (reorg protection)
type EthConfirmer struct {
	utils.StartStopOnce
	orm       ORM
	lggr      logger.Logger
	ethClient evmclient.Client
	ChainKeyStore
	estimator      gas.Estimator
	resumeCallback ResumeCallback

	keyStates []ethkey.State

	mb        *utils.Mailbox[*evmtypes.Head]
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup

	nConsecutiveBlocksChainTooShort int
}

// NewEthConfirmer instantiates a new eth confirmer
func NewEthConfirmer(orm ORM, ethClient evmclient.Client, config Config, keystore KeyStore,
	keyStates []ethkey.State, estimator gas.Estimator, resumeCallback ResumeCallback, lggr logger.Logger) *EthConfirmer {

	ctx, cancel := context.WithCancel(context.Background())
	lggr = lggr.Named("EthConfirmer")

	return &EthConfirmer{
		utils.StartStopOnce{},
		orm,
		lggr,
		ethClient,
		ChainKeyStore{
			*ethClient.ChainID(),
			config,
			keystore,
		},
		estimator,
		resumeCallback,
		keyStates,
		utils.NewSingleMailbox[*evmtypes.Head](),
		ctx,
		cancel,
		sync.WaitGroup{},
		0,
	}
}

// Start is a comment to appease the linter
func (ec *EthConfirmer) Start(_ context.Context) error {
	return ec.StartOnce("EthConfirmer", func() error {
		if ec.config.EvmGasBumpThreshold() == 0 {
			ec.lggr.Infow("Gas bumping is disabled (ETH_GAS_BUMP_THRESHOLD set to 0)", "ethGasBumpThreshold", 0)
		} else {
			ec.lggr.Infow(fmt.Sprintf("Gas bumping is enabled, unconfirmed transactions will have their gas price bumped every %d blocks", ec.config.EvmGasBumpThreshold()), "ethGasBumpThreshold", ec.config.EvmGasBumpThreshold())
		}

		ec.wg.Add(1)
		go ec.runLoop()

		return nil
	})
}

// Close is a comment to appease the linter
func (ec *EthConfirmer) Close() error {
	return ec.StopOnce("EthConfirmer", func() error {
		ec.ctxCancel()
		ec.wg.Wait()
		return nil
	})
}

func (ec *EthConfirmer) Name() string {
	return ec.lggr.Name()
}

func (ec *EthConfirmer) HealthReport() map[string]error {
	return map[string]error{ec.Name(): ec.Healthy()}
}

func (ec *EthConfirmer) runLoop() {
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
func (ec *EthConfirmer) ProcessHead(ctx context.Context, head *evmtypes.Head) error {
	ctx, cancel := context.WithTimeout(ctx, processHeadTimeout)
	defer cancel()

	return ec.processHead(ctx, head)
}

// NOTE: This SHOULD NOT be run concurrently or it could behave badly
func (ec *EthConfirmer) processHead(ctx context.Context, head *evmtypes.Head) error {
	mark := time.Now()

	ec.lggr.Debugw("processHead start", "headNum", head.Number, "id", "eth_confirmer")

	if err := ec.orm.SetBroadcastBeforeBlockNum(head.Number, ec.chainID); err != nil {
		return errors.Wrap(err, "SetBroadcastBeforeBlockNum failed")
	}
	if err := ec.CheckConfirmedMissingReceipt(ctx); err != nil {
		return errors.Wrap(err, "CheckConfirmedMissingReceipt failed")
	}

	if err := ec.CheckForReceipts(ctx, head.Number); err != nil {
		return errors.Wrap(err, "CheckForReceipts failed")
	}

	ec.lggr.Debugw("Finished CheckForReceipts", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	if err := ec.RebroadcastWhereNecessary(ctx, head.Number); err != nil {
		return errors.Wrap(err, "RebroadcastWhereNecessary failed")
	}

	ec.lggr.Debugw("Finished RebroadcastWhereNecessary", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	if err := ec.EnsureConfirmedTransactionsInLongestChain(ctx, head); err != nil {
		return errors.Wrap(err, "EnsureConfirmedTransactionsInLongestChain failed")
	}

	ec.lggr.Debugw("Finished EnsureConfirmedTransactionsInLongestChain", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")

	if ec.resumeCallback != nil {
		mark = time.Now()
		if err := ec.ResumePendingTaskRuns(ctx, head); err != nil {
			return errors.Wrap(err, "ResumePendingTaskRuns failed")
		}

		ec.lggr.Debugw("Finished ResumePendingTaskRuns", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	}

	ec.lggr.Debugw("processHead finish", "headNum", head.Number, "id", "eth_confirmer")

	return nil
}

// CheckConfirmedMissingReceipt will attempt to re-send any transaction in the
// state of "confirmed_missing_receipt". If we get back any type of senderror
// other than "nonce too low" it means that this transaction isn't actually
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
// 4. RPC node 2 still has an old view of the chain, so it returns us "nonce too low" indicating "no problem this transaction is already mined"
// 5. Now the transaction is marked "confirmed_missing_receipt" but the latest chain does not actually include it
// 6. Now we are reliant on the EthResender to propagate it, and this transaction will not be gas bumped, so in the event of gas spikes it could languish or even be evicted from the mempool and hold up the queue
// 7. Even if/when RPC node 2 catches up, the transaction is still stuck in state "confirmed_missing_receipt"
//
// This scenario might sound unlikely but has been observed to happen multiple times in the wild on Polygon.
func (ec *EthConfirmer) CheckConfirmedMissingReceipt(ctx context.Context) (err error) {
	attempts, err := ec.orm.FindEtxAttemptsConfirmedMissingReceipt(ec.chainID)
	if err != nil {
		return err
	}
	if len(attempts) == 0 {
		return nil
	}
	ec.lggr.Infow(fmt.Sprintf("Found %d transactions confirmed_missing_receipt. The RPC node did not give us a receipt for these transactions even though it should have been mined. This could be due to using the wallet with an external account, or if the primary node is not synced or not propagating transactions properly", len(attempts)), "attempts", attempts)
	reqs, err := batchSendTransactions(ctx, ec.orm, attempts, int(ec.config.EvmRPCDefaultBatchSize()), ec.lggr, ec.ethClient)
	if err != nil {
		ec.lggr.Debugw("Batch sending transactions failed", err)
	}
	var ethTxIDsToUnconfirm []int64
	for idx, req := range reqs {
		// Add to Unconfirm array, all tx where error wasn't NonceTooLow.
		if req.Error != nil {
			err := evmclient.NewSendError(req.Error)
			if err.IsNonceTooLowError() || err.IsTransactionAlreadyMined() {
				continue
			}
		}

		ethTxIDsToUnconfirm = append(ethTxIDsToUnconfirm, attempts[idx].EthTxID)
	}
	err = ec.orm.UpdateEthTxsUnconfirmed(ethTxIDsToUnconfirm)

	if err != nil {
		return err
	}
	return
}

// CheckForReceipts finds attempts that are still pending and checks to see if a receipt is present for the given block number
func (ec *EthConfirmer) CheckForReceipts(ctx context.Context, blockNum int64) error {
	attempts, err := ec.orm.FindEthTxAttemptsRequiringReceiptFetch(ec.chainID)
	if err != nil {
		return errors.Wrap(err, "FindEthTxAttemptsRequiringReceiptFetch failed")
	}
	if len(attempts) == 0 {
		return nil
	}

	ec.lggr.Debugw(fmt.Sprintf("Fetching receipts for %v transaction attempts", len(attempts)), "blockNum", blockNum)

	attemptsByAddress := make(map[gethCommon.Address][]EthTxAttempt)
	for _, att := range attempts {
		attemptsByAddress[att.EthTx.FromAddress] = append(attemptsByAddress[att.EthTx.FromAddress], att)
	}

	for from, attempts := range attemptsByAddress {
		minedTransactionCount, err := ec.getMinedTransactionCount(ctx, from)
		if err != nil {
			return errors.Wrapf(err, "unable to fetch pending nonce for address: %v", from)
		}

		// separateLikelyConfirmedAttempts is used as an optimisation: there is
		// no point trying to fetch receipts for attempts with a nonce higher
		// than the highest nonce the RPC node thinks it has seen
		likelyConfirmed := ec.separateLikelyConfirmedAttempts(from, attempts, minedTransactionCount)
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

	if err := ec.orm.MarkAllConfirmedMissingReceipt(ec.chainID); err != nil {
		return errors.Wrap(err, "unable to mark eth_txes as 'confirmed_missing_receipt'")
	}

	if err := ec.orm.MarkOldTxesMissingReceiptAsErrored(blockNum, ec.config.EvmFinalityDepth(), ec.chainID); err != nil {
		return errors.Wrap(err, "unable to confirm buried unconfirmed eth_txes")
	}
	return nil
}

func (ec *EthConfirmer) separateLikelyConfirmedAttempts(from gethCommon.Address, attempts []EthTxAttempt, minedTransactionCount uint64) []EthTxAttempt {
	if len(attempts) == 0 {
		return attempts
	}

	firstAttemptNonce := *attempts[len(attempts)-1].EthTx.Nonce
	lastAttemptNonce := *attempts[0].EthTx.Nonce
	latestMinedNonce := int64(minedTransactionCount) - 1 // this can be -1 if a transaction has never been mined on this account
	ec.lggr.Debugw(fmt.Sprintf("There are %d attempts from address %s, mined transaction count is %d (latest mined nonce is %d) and for the attempts' nonces: first = %d, last = %d",
		len(attempts), from.Hex(), minedTransactionCount, latestMinedNonce, firstAttemptNonce, lastAttemptNonce), "nAttempts", len(attempts), "fromAddress", from, "minedTransactionCount", minedTransactionCount, "latestMinedNonce", latestMinedNonce, "firstAttemptNonce", firstAttemptNonce, "lastAttemptNonce", lastAttemptNonce)

	likelyConfirmed := attempts
	// attempts are ordered by nonce ASC
	for i := 0; i < len(attempts); i++ {
		// If the attempt nonce is lower or equal to the latestBlockNonce
		// it must have been confirmed, we just didn't get a receipt yet
		//
		// Examples:
		// 3 transactions confirmed, highest has nonce 2
		// 5 total attempts, highest has nonce 4
		// minedTransactionCount=3
		// likelyConfirmed will be attempts[0:3] which gives the first 3 transactions, as expected
		if *attempts[i].EthTx.Nonce > int64(minedTransactionCount) {
			ec.lggr.Debugf("Marking attempts as likely confirmed just before index %v, at nonce: %v", i, *attempts[i].EthTx.Nonce)
			likelyConfirmed = attempts[0:i]
			break
		}
	}

	if len(likelyConfirmed) == 0 {
		ec.lggr.Debug("There are no likely confirmed attempts - so will skip checking any")
	}

	return likelyConfirmed
}

func (ec *EthConfirmer) fetchAndSaveReceipts(ctx context.Context, attempts []EthTxAttempt, blockNum int64) error {
	promTxAttemptCount.WithLabelValues(ec.chainID.String()).Set(float64(len(attempts)))

	batchSize := int(ec.config.EvmRPCDefaultBatchSize())
	if batchSize == 0 {
		batchSize = len(attempts)
	}
	var allReceipts []evmtypes.Receipt
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
		if err := ec.orm.SaveFetchedReceipts(receipts, ec.chainID); err != nil {
			return errors.Wrap(err, "saveFetchedReceipts failed")
		}
		promNumConfirmedTxs.WithLabelValues(ec.chainID.String()).Add(float64(len(receipts)))

		allReceipts = append(allReceipts, receipts...)
	}

	observeUntilTxConfirmed(ec.chainID, attempts, allReceipts)

	return nil
}

func (ec *EthConfirmer) getMinedTransactionCount(ctx context.Context, from gethCommon.Address) (nonce uint64, err error) {
	return ec.ethClient.NonceAt(ctx, from, nil)
}

// Note this function will increment promRevertedTxCount upon receiving
// a reverted transaction receipt. Should only be called with unconfirmed attempts.
func (ec *EthConfirmer) batchFetchReceipts(ctx context.Context, attempts []EthTxAttempt, blockNum int64) (receipts []evmtypes.Receipt, err error) {
	var reqs []rpc.BatchElem

	// Metadata is required to determine whether a tx is forwarded or not.
	if ec.config.EvmUseForwarders() {
		err = ec.orm.PreloadEthTxes(attempts)
		if err != nil {
			return nil, errors.Wrap(err, "EthConfirmer#batchFetchReceipts error loading txs for attempts")
		}
	}

	for _, attempt := range attempts {
		req := rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{attempt.Hash},
			Result: &evmtypes.Receipt{},
		}
		reqs = append(reqs, req)
	}

	lggr := ec.lggr.Named("batchFetchReceipts").With("blockNum", blockNum)

	err = ec.ethClient.BatchCallContext(ctx, reqs)
	if err != nil {
		return nil, errors.Wrap(err, "EthConfirmer#batchFetchReceipts error fetching receipts with BatchCallContext")
	}

	for i, req := range reqs {
		attempt := attempts[i]
		result, err := req.Result, req.Error

		receipt, is := result.(*evmtypes.Receipt)
		if !is {
			return nil, errors.Errorf("expected result to be a %T, got %T", (*evmtypes.Receipt)(nil), result)
		}

		l := logger.Sugared(attempt.EthTx.GetLogger(lggr).With(
			"txHash", attempt.Hash.Hex(), "ethTxAttemptID", attempt.ID,
			"ethTxID", attempt.EthTxID, "err", err, "nonce", attempt.EthTx.Nonce,
		))

		if err != nil {
			l.Error("FetchReceipt failed")
			continue
		}

		if receipt == nil {
			// NOTE: This should never happen, but it seems safer to check
			// regardless to avoid a potential panic
			l.AssumptionViolation("got nil receipt")
			continue
		}

		if receipt.IsZero() {
			l.Debug("Still waiting for receipt")
			continue
		}

		l = logger.Sugared(l.With("blockHash", receipt.BlockHash.Hex(), "status", receipt.Status, "transactionIndex", receipt.TransactionIndex))

		if receipt.IsUnmined() {
			l.Debug("Got receipt for transaction but it's still in the mempool and not included in a block yet")
			continue
		}

		l.Debugw("Got receipt for transaction", "blockNumber", receipt.BlockNumber, "gasUsed", receipt.GasUsed)

		if receipt.TxHash != attempt.Hash {
			l.Errorf("Invariant violation, expected receipt with hash %s to have same hash as attempt with hash %s", receipt.TxHash.Hex(), attempt.Hash.Hex())
			continue
		}

		if receipt.BlockNumber == nil {
			l.Error("Invariant violation, receipt was missing block number")
			continue
		}

		if receipt.Status == 0 {
			// Do an eth call to obtain the revert reason.
			_, errCall := ec.ethClient.CallContract(ctx, ethereum.CallMsg{
				From:       attempt.EthTx.FromAddress,
				To:         &attempt.EthTx.ToAddress,
				Gas:        uint64(attempt.EthTx.GasLimit),
				GasPrice:   attempt.GasPrice.ToInt(),
				GasFeeCap:  attempt.GasFeeCap.ToInt(),
				GasTipCap:  attempt.GasTipCap.ToInt(),
				Value:      nil,
				Data:       attempt.EthTx.EncodedPayload,
				AccessList: nil,
			}, receipt.BlockNumber)
			rpcError, errExtract := evmclient.ExtractRPCError(errCall)
			if errExtract == nil {
				l.Warnw("transaction reverted on-chain", "hash", receipt.TxHash, "rpcError", rpcError.String())
			} else {
				l.Warnw("transaction reverted on-chain unable to extract revert reason", "hash", receipt.TxHash, "err", err)
			}
			// This might increment more than once e.g. in case of re-orgs going back and forth we might re-fetch the same receipt
			promRevertedTxCount.WithLabelValues(ec.chainID.String()).Add(1)
		} else {
			promNumSuccessfulTxs.WithLabelValues(ec.chainID.String()).Add(1)
		}

		// This is only recording forwarded tx that were mined and have a status.
		// Counters are prone to being inaccurate due to re-orgs.
		if ec.config.EvmUseForwarders() {
			meta, err := attempt.EthTx.GetMeta()
			if err == nil && meta != nil && meta.FwdrDestAddress != nil {
				// promFwdTxCount takes two labels, chainId and a boolean of whether a tx was successful or not.
				promFwdTxCount.WithLabelValues(ec.chainID.String(), strconv.FormatBool(receipt.Status != 0)).Add(1)
			}
		}

		receipts = append(receipts, *receipt)
	}

	return
}

// RebroadcastWhereNecessary bumps gas or resends transactions that were previously out-of-eth
func (ec *EthConfirmer) RebroadcastWhereNecessary(ctx context.Context, blockHeight int64) error {
	var wg sync.WaitGroup

	// It is safe to process separate keys concurrently
	// NOTE: This design will block one key if another takes a really long time to execute
	wg.Add(len(ec.keyStates))
	errors := []error{}
	var errMu sync.Mutex
	for _, key := range ec.keyStates {
		go func(fromAddress gethCommon.Address) {
			if err := ec.rebroadcastWhereNecessary(ctx, fromAddress, blockHeight); err != nil {
				errMu.Lock()
				errors = append(errors, err)
				errMu.Unlock()
				ec.lggr.Errorw("Error in RebroadcastWhereNecessary", "error", err, "fromAddress", fromAddress)
			}

			wg.Done()
		}(key.Address.Address())
	}

	wg.Wait()

	return multierr.Combine(errors...)
}

func (ec *EthConfirmer) rebroadcastWhereNecessary(ctx context.Context, address gethCommon.Address, blockHeight int64) error {
	if err := ec.handleAnyInProgressAttempts(ctx, address, blockHeight); err != nil {
		return errors.Wrap(err, "handleAnyInProgressAttempts failed")
	}

	threshold := int64(ec.config.EvmGasBumpThreshold())
	bumpDepth := int64(ec.config.EvmGasBumpTxDepth())
	maxInFlightTransactions := ec.config.EvmMaxInFlightTransactions()
	etxs, err := ec.FindEthTxsRequiringRebroadcast(ctx, ec.lggr, address, blockHeight, threshold, bumpDepth, maxInFlightTransactions, ec.chainID)
	if err != nil {
		return errors.Wrap(err, "FindEthTxsRequiringRebroadcast failed")
	}
	for _, etx := range etxs {
		lggr := etx.GetLogger(ec.lggr)

		attempt, err := ec.attemptForRebroadcast(ctx, lggr, *etx)
		if err != nil {
			return errors.Wrap(err, "attemptForRebroadcast failed")
		}

		lggr.Debugw("Rebroadcasting transaction", "nPreviousAttempts", len(etx.EthTxAttempts), "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)

		if err := ec.orm.SaveInProgressAttempt(&attempt); err != nil {
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
func (ec *EthConfirmer) handleAnyInProgressAttempts(ctx context.Context, address gethCommon.Address, blockHeight int64) error {
	attempts, err := ec.orm.GetInProgressEthTxAttempts(ctx, address, ec.chainID)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "GetInProgressEthTxAttempts failed")
	}
	for _, a := range attempts {
		err := ec.handleInProgressAttempt(ctx, a.EthTx.GetLogger(ec.lggr), a.EthTx, a, blockHeight)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			return errors.Wrap(err, "handleInProgressAttempt failed")
		}
	}
	return nil
}

// FindEthTxsRequiringRebroadcast returns attempts that hit insufficient eth,
// and attempts that need bumping, in nonce ASC order
func (ec *EthConfirmer) FindEthTxsRequiringRebroadcast(ctx context.Context, lggr logger.Logger, address gethCommon.Address, blockNum, gasBumpThreshold, bumpDepth int64, maxInFlightTransactions uint32, chainID big.Int) (etxs []*EthTx, err error) {
	// NOTE: These two queries could be combined into one using union but it
	// becomes harder to read and difficult to test in isolation. KISS principle
	etxInsufficientEths, err := ec.orm.FindEthTxsRequiringResubmissionDueToInsufficientEth(address, chainID, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}

	if len(etxInsufficientEths) > 0 {
		lggr.Infow(fmt.Sprintf("Found %d transactions to be re-sent that were previously rejected due to insufficient eth balance", len(etxInsufficientEths)), "blockNum", blockNum, "address", address)
	}

	// TODO: Just pass the Q through everything
	etxBumps, err := ec.orm.FindEthTxsRequiringGasBump(ctx, address, blockNum, gasBumpThreshold, bumpDepth, chainID)
	if ctx.Err() != nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if len(etxBumps) > 0 {
		// txes are ordered by nonce asc so the first will always be the oldest
		etx := etxBumps[0]
		// attempts are ordered by time sent asc so first will always be the oldest
		var oldestBlocksBehind int64 = -1 // It should never happen that the oldest attempt has no BroadcastBeforeBlockNum set, but in case it does, we shouldn't crash - log this sentinel value instead
		if len(etx.EthTxAttempts) > 0 {
			oldestBlockNum := etx.EthTxAttempts[0].BroadcastBeforeBlockNum
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
		return *(etxs[i].Nonce) < *(etxs[j].Nonce)
	})

	if maxInFlightTransactions > 0 && len(etxs) > int(maxInFlightTransactions) {
		lggr.Warnf("%d transactions to rebroadcast which exceeds limit of %d. %s", len(etxs), maxInFlightTransactions, label.MaxInFlightTransactionsWarning)
		etxs = etxs[:maxInFlightTransactions]
	}

	return
}

func (ec *EthConfirmer) attemptForRebroadcast(ctx context.Context, lggr logger.Logger, etx EthTx) (attempt EthTxAttempt, err error) {
	if len(etx.EthTxAttempts) > 0 {
		etx.EthTxAttempts[0].EthTx = etx
		previousAttempt := etx.EthTxAttempts[0]
		logFields := ec.logFieldsPreviousAttempt(previousAttempt)
		if previousAttempt.State == EthTxAttemptInsufficientEth {
			// Do not create a new attempt if we ran out of eth last time since bumping gas is pointless
			// Instead try to resubmit the same attempt at the same price, in the hope that the wallet was funded since our last attempt
			lggr.Debugw("Rebroadcast InsufficientEth", logFields...)
			previousAttempt.State = EthTxAttemptInProgress
			return previousAttempt, nil
		}
		attempt, err = ec.bumpGas(ctx, etx, etx.EthTxAttempts)

		if gas.IsBumpErr(err) {
			lggr.Errorw("Failed to bump gas", append(logFields, "err", err)...)
			// Do not create a new attempt if bumping gas would put us over the limit or cause some other problem
			// Instead try to resubmit the previous attempt, and keep resubmitting until its accepted
			previousAttempt.BroadcastBeforeBlockNum = nil
			previousAttempt.State = EthTxAttemptInProgress
			return previousAttempt, nil
		}
		return attempt, err
	}
	return attempt, errors.Errorf("invariant violation: EthTx %v was unconfirmed but didn't have any attempts. "+
		"Falling back to default gas price instead."+
		"This is a bug! Please report to https://github.com/smartcontractkit/chainlink/issues", etx.ID)
}

func (ec *EthConfirmer) logFieldsPreviousAttempt(attempt EthTxAttempt) []interface{} {
	etx := attempt.EthTx
	return []interface{}{
		"etxID", etx.ID,
		"txHash", attempt.Hash,
		"previousAttempt", attempt,
		"gasLimit", etx.GasLimit,
		"maxGasPrice", ec.config.EvmMaxGasPriceWei(),
		"nonce", etx.Nonce,
	}
}

func (ec *EthConfirmer) bumpGas(ctx context.Context, etx EthTx, previousAttempts []EthTxAttempt) (bumpedAttempt EthTxAttempt, err error) {
	priorAttempts := make([]gas.PriorAttempt, len(previousAttempts))
	// This feels a bit useless but until we get iterators there is no other
	// way to cast an array of structs to an array of interfaces
	for i, attempt := range previousAttempts {
		priorAttempts[i] = attempt
	}
	previousAttempt := previousAttempts[0]
	logFields := ec.logFieldsPreviousAttempt(previousAttempt)
	keySpecificMaxGasPriceWei := ec.config.KeySpecificMaxGasPriceWei(etx.FromAddress)
	switch previousAttempt.TxType {
	case 0x0: // Legacy
		var bumpedGasPrice *assets.Wei
		var bumpedGasLimit uint32
		bumpedGasPrice, bumpedGasLimit, err = ec.estimator.BumpLegacyGas(ctx, previousAttempt.GasPrice, etx.GasLimit, keySpecificMaxGasPriceWei, priorAttempts)
		if err == nil {
			promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
			ec.lggr.Debugw("Rebroadcast bumping gas for Legacy tx", append(logFields, "bumpedGasPrice", bumpedGasPrice.String())...)
			return ec.NewLegacyAttempt(etx, bumpedGasPrice, bumpedGasLimit)
		}
	case 0x2: // EIP1559
		var bumpedFee gas.DynamicFee
		var bumpedGasLimit uint32
		original := previousAttempt.DynamicFee()
		bumpedFee, bumpedGasLimit, err = ec.estimator.BumpDynamicFee(ctx, original, etx.GasLimit, keySpecificMaxGasPriceWei, priorAttempts)
		if err == nil {
			promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
			ec.lggr.Debugw("Rebroadcast bumping gas for DynamicFee tx", append(logFields, "bumpedTipCap", bumpedFee.TipCap.String(), "bumpedFeeCap", bumpedFee.FeeCap.String())...)
			return ec.NewDynamicFeeAttempt(etx, bumpedFee, bumpedGasLimit)
		}
	default:
		err = errors.Errorf("invariant violation: Attempt %v had unrecognised transaction type %v"+
			"This is a bug! Please report to https://github.com/smartcontractkit/chainlink/issues", previousAttempt.ID, previousAttempt.TxType)
	}

	if errors.Is(errors.Cause(err), gas.ErrBumpGasExceedsLimit) {
		promGasBumpExceedsLimit.WithLabelValues(ec.chainID.String()).Inc()
	}

	return bumpedAttempt, errors.Wrap(err, "error bumping gas")
}

func (ec *EthConfirmer) handleInProgressAttempt(ctx context.Context, lggr logger.Logger, etx EthTx, attempt EthTxAttempt, blockHeight int64) error {
	if attempt.State != EthTxAttemptInProgress {
		return errors.Errorf("invariant violation: expected eth_tx_attempt %v to be in_progress, it was %s", attempt.ID, attempt.State)
	}

	now := time.Now()
	sendError := sendTransaction(ctx, ec.ethClient, attempt, etx, lggr)

	if sendError.IsTerminallyUnderpriced() {
		// This should really not ever happen in normal operation since we
		// already bumped above the required minimum in ethBroadcaster.
		ec.lggr.Warnw("Got terminally underpriced error for gas bump, this should never happen unless the remote RPC node changed its configuration on the fly, or you are using multiple RPC nodes with different minimum gas price requirements. This is not recommended", "err", sendError, "attempt", attempt)
		// "Lazily" load attempts here since the overwhelmingly common case is
		// that we don't need them unless we enter this path
		if err := ec.orm.LoadEthTxAttempts(&etx, pg.WithParentCtx(ctx)); err != nil {
			return errors.Wrap(err, "failed to load EthTxAttempts while bumping on terminally underpriced error")
		}
		if len(etx.EthTxAttempts) == 0 {
			err := errors.New("expected to find at least 1 attempt")
			logger.Sugared(ec.lggr).AssumptionViolationw(err.Error(), "err", err, "attempt", attempt)
			return err
		}
		if attempt.ID != etx.EthTxAttempts[0].ID {
			err := errors.New("expected highest priced attempt to be the current in_progress attempt")
			logger.Sugared(ec.lggr).AssumptionViolationw(err.Error(), "err", err, "attempt", attempt, "ethTxAttempts", etx.EthTxAttempts)
			return err
		}
		replacementAttempt, err := ec.bumpGas(ctx, etx, etx.EthTxAttempts)
		if err != nil {
			return errors.Wrap(err, "could not bump gas for terminally underpriced transaction")
		}
		promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
		lggr.With(
			"sendError", sendError,
			"maxGasPriceConfig", ec.config.EvmMaxGasPriceWei(),
			"previousAttempt", attempt,
			"replacementAttempt", replacementAttempt,
		).Errorf("gas price was rejected by the eth node for being too low. Eth node returned: '%s'", sendError.Error())

		if err := ec.orm.SaveReplacementInProgressAttempt(attempt, &replacementAttempt); err != nil {
			return errors.Wrap(err, "saveReplacementInProgressAttempt failed")
		}
		return ec.handleInProgressAttempt(ctx, lggr, etx, replacementAttempt, blockHeight)
	}

	if sendError.IsTemporarilyUnderpriced() {
		// Most likely scenario here is a parity node that is rejecting
		// low-priced transactions due to mempool pressure
		//
		// In that case, the safest thing to do is to pretend the transaction
		// was accepted and continue the normal gas bumping cycle until we can
		// get it into the mempool
		lggr.Debugw("Transaction temporarily underpriced", "attemptID", attempt.ID, "err", sendError.Error(), "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)
		sendError = nil
	}

	if sendError.IsTxFeeExceedsCap() {
		// The gas price was bumped too high. This transaction attempt cannot be accepted.
		//
		// Best thing we can do is to re-send the previous attempt at the old
		// price and discard this bumped version.
		lggr.Errorw(fmt.Sprintf("Transaction gas bump failed; %s", label.RPCTxFeeCapConfiguredIncorrectlyWarning),
			"err", sendError,
			"gasPrice", attempt.GasPrice,
			"gasLimit", etx.GasLimit,
			"signedRawTx", hexutil.Encode(attempt.SignedRawTx),
			"blockHeight", blockHeight,
			"id", "RPCTxFeeCapExceeded",
		)
		return ec.orm.DeleteInProgressAttempt(ctx, attempt)
	}

	if sendError.Fatal() {
		// WARNING: This should never happen!
		// Should NEVER be fatal this is an invariant violation. The
		// EthBroadcaster can never create an EthTxAttempt that will
		// fatally error.
		//
		// The only scenario imaginable where this might take place is if
		// geth/parity have been updated between broadcasting and confirming steps.
		lggr.Criticalw("Invariant violation: fatal error while re-attempting transaction",
			"err", sendError,
			"signedRawTx", hexutil.Encode(attempt.SignedRawTx),
			"blockHeight", blockHeight,
		)
		// This will loop continuously on every new head so it must be handled manually by the node operator!
		return ec.orm.DeleteInProgressAttempt(ctx, attempt)
	}

	if sendError.IsNonceTooLowError() || sendError.IsTransactionAlreadyMined() {
		// Nonce too low indicated that a transaction at this nonce was confirmed already.
		// Mark confirmed_missing_receipt and wait for the next cycle to try to get a receipt
		sendError = nil
		lggr.Debugw("Nonce already used", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash.Hex(), "err", sendError)
		timeout := ec.config.DatabaseDefaultQueryTimeout()
		return ec.orm.SaveConfirmedMissingReceiptAttempt(ctx, timeout, &attempt, now)
	}

	if sendError.IsReplacementUnderpriced() {
		// Our system constraints guarantee that the attempt referenced in this
		// function has the highest gas price of all attempts.
		//
		// Thus, there are only two possible scenarios where this can happen.
		//
		// 1. Our gas bump was insufficient compared to our previous attempt
		// 2. An external wallet used the account to manually send a transaction
		// at a higher gas price
		//
		// In this case the simplest and most robust way to recover is to ignore
		// this attempt and wait until the next bump threshold is reached in
		// order to bump again.
		lggr.Errorw(fmt.Sprintf("Replacement transaction underpriced for eth_tx %v. "+
			"Eth node returned error: '%s'. "+
			"Either you have set ETH_GAS_BUMP_PERCENT (currently %v%%) too low or an external wallet used this account. "+
			"Please note that using your node's private keys outside of the chainlink node is NOT SUPPORTED and can lead to missed transactions.",
			etx.ID, sendError.Error(), ec.config.EvmGasBumpPercent()), "err", sendError, "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)

		// Assume success and hand off to the next cycle.
		sendError = nil
	}

	if sendError.IsInsufficientEth() {
		lggr.Criticalw(fmt.Sprintf("Tx 0x%x with type 0x%d was rejected due to insufficient eth: %s\n"+
			"ACTION REQUIRED: Chainlink wallet with address 0x%x is OUT OF FUNDS",
			attempt.ID, attempt.Hash, sendError.Error(), etx.FromAddress,
		), "err", sendError, "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)
		timeout := ec.config.DatabaseDefaultQueryTimeout()
		return ec.orm.SaveInsufficientEthAttempt(timeout, &attempt, now)
	}

	if sendError == nil {
		lggr.Debugw("Successfully broadcast transaction", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash.Hex())
		timeout := ec.config.DatabaseDefaultQueryTimeout()
		return ec.orm.SaveSentAttempt(timeout, &attempt, now)
	}

	// Any other type of error is considered temporary or resolvable by the
	// node operator. The node may have it in the mempool so we must keep the
	// attempt (leave it in_progress). Safest thing to do is bail out and wait
	// for the next head.
	return errors.Wrapf(sendError, "unexpected error sending eth_tx %v with hash %s", etx.ID, attempt.Hash.Hex())
}

// EnsureConfirmedTransactionsInLongestChain finds all confirmed eth_txes up to the depth
// of the given chain and ensures that every one has a receipt with a block hash that is
// in the given chain.
//
// If any of the confirmed transactions does not have a receipt in the chain, it has been
// re-org'd out and will be rebroadcast.
func (ec *EthConfirmer) EnsureConfirmedTransactionsInLongestChain(ctx context.Context, head *evmtypes.Head) error {
	if head.ChainLength() < ec.config.EvmFinalityDepth() {
		logArgs := []interface{}{
			"chainLength", head.ChainLength(), "evmFinalityDepth", ec.config.EvmFinalityDepth(),
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
	etxs, err := ec.orm.FindTransactionsConfirmedInBlockRange(head.Number, head.EarliestInChain().Number, ec.chainID)
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
	wg.Add(len(ec.keyStates))
	for _, key := range ec.keyStates {
		go func(fromAddress gethCommon.Address) {
			if err := ec.handleAnyInProgressAttempts(ctx, fromAddress, head.Number); err != nil {
				errMu.Lock()
				errors = append(errors, err)
				errMu.Unlock()
				ec.lggr.Errorw("Error in handleAnyInProgressAttempts", "err", err, "fromAddress", fromAddress)
			}

			wg.Done()
		}(key.Address.Address())
	}

	wg.Wait()

	return multierr.Combine(errors...)
}

func hasReceiptInLongestChain(etx EthTx, head *evmtypes.Head) bool {
	for {
		for _, attempt := range etx.EthTxAttempts {
			for _, receipt := range attempt.EthReceipts {
				if receipt.BlockHash == head.Hash && receipt.BlockNumber == head.Number {
					return true
				}
			}
		}
		if head.Parent == nil {
			return false
		}
		head = head.Parent
	}
}

func (ec *EthConfirmer) markForRebroadcast(etx EthTx, head *evmtypes.Head) error {
	if len(etx.EthTxAttempts) == 0 {
		return errors.Errorf("invariant violation: expected eth_tx %v to have at least one attempt", etx.ID)
	}

	// Rebroadcast the one with the highest gas price
	attempt := etx.EthTxAttempts[0]
	var receipt EthReceipt
	if len(attempt.EthReceipts) > 0 {
		receipt = attempt.EthReceipts[0]
	}

	ec.lggr.Infow(fmt.Sprintf("Re-org detected. Rebroadcasting transaction %s which may have been re-org'd out of the main chain", attempt.Hash.Hex()),
		"txhash", attempt.Hash.Hex(),
		"currentBlockNum", head.Number,
		"currentBlockHash", head.Hash.Hex(),
		"replacementBlockHashAtConfirmedHeight", head.HashAtHeight(receipt.BlockNumber),
		"confirmedInBlockNum", receipt.BlockNumber,
		"confirmedInBlockHash", receipt.BlockHash,
		"confirmedInTxIndex", receipt.TransactionIndex,
		"ethTxID", etx.ID,
		"attemptID", attempt.ID,
		"receiptID", receipt.ID,
		"nReceipts", len(attempt.EthReceipts),
		"id", "eth_confirmer")

	// Put it back in progress and delete all receipts (they do not apply to the new chain)
	err := ec.orm.UpdateEthTxForRebroadcast(etx, attempt)
	return errors.Wrap(err, "markForRebroadcast failed")
}

// ForceRebroadcast sends a transaction for every nonce in the given nonce range at the given gas price.
// If an eth_tx exists for this nonce, we re-send the existing eth_tx with the supplied parameters.
// If an eth_tx doesn't exist for this nonce, we send a zero transaction.
// This operates completely orthogonal to the normal EthConfirmer and can result in untracked attempts!
// Only for emergency usage.
// This is in case of some unforeseen scenario where the node is refusing to release the lock. KISS.
func (ec *EthConfirmer) ForceRebroadcast(beginningNonce uint, endingNonce uint, gasPriceWei uint64, address gethCommon.Address, overrideGasLimit uint32) error {
	ec.lggr.Infof("ForceRebroadcast: will rebroadcast transactions for all nonces between %v and %v", beginningNonce, endingNonce)

	for n := beginningNonce; n <= endingNonce; n++ {
		etx, err := ec.orm.FindEthTxWithNonce(address, n)
		if err != nil {
			return errors.Wrap(err, "ForceRebroadcast failed")
		}
		if etx == nil {
			ec.lggr.Debugf("ForceRebroadcast: no eth_tx found with nonce %v, will rebroadcast empty transaction", n)
			hash, err := ec.sendEmptyTransaction(context.TODO(), address, n, overrideGasLimit, gasPriceWei)
			if err != nil {
				ec.lggr.Errorw("ForceRebroadcast: failed to send empty transaction", "nonce", n, "err", err)
				continue
			}
			ec.lggr.Infow("ForceRebroadcast: successfully rebroadcast empty transaction", "nonce", n, "hash", hash.String())
		} else {
			ec.lggr.Debugf("ForceRebroadcast: got eth_tx %v with nonce %v, will rebroadcast this transaction", etx.ID, *etx.Nonce)
			if overrideGasLimit != 0 {
				etx.GasLimit = overrideGasLimit
			}
			attempt, err := ec.NewLegacyAttempt(*etx, assets.NewWeiI(int64(gasPriceWei)), etx.GasLimit)
			if err != nil {
				ec.lggr.Errorw("ForceRebroadcast: failed to create new attempt", "ethTxID", etx.ID, "err", err)
				continue
			}
			if err := sendTransaction(context.TODO(), ec.ethClient, attempt, *etx, ec.lggr); err != nil {
				ec.lggr.Errorw(fmt.Sprintf("ForceRebroadcast: failed to rebroadcast eth_tx %v with nonce %v and gas limit %v: %s", etx.ID, *etx.Nonce, etx.GasLimit, err.Error()), "err", err, "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)
				continue
			}
			ec.lggr.Infof("ForceRebroadcast: successfully rebroadcast eth_tx %v with hash: 0x%x", etx.ID, attempt.Hash)
		}
	}
	return nil
}

func (ec *EthConfirmer) sendEmptyTransaction(ctx context.Context, fromAddress gethCommon.Address, nonce uint, overrideGasLimit uint32, gasPriceWei uint64) (gethCommon.Hash, error) {
	gasLimit := overrideGasLimit
	if gasLimit == 0 {
		gasLimit = ec.config.EvmGasLimitDefault()
	}
	tx, err := sendEmptyTransaction(ctx, ec.ethClient, ec.keystore, uint64(nonce), gasLimit, big.NewInt(int64(gasPriceWei)), fromAddress, &ec.chainID)
	if err != nil {
		return gethCommon.Hash{}, errors.Wrap(err, "(EthConfirmer).sendEmptyTransaction failed")
	}
	return tx.Hash(), nil
}

// ResumePendingTaskRuns issues callbacks to task runs that are pending waiting for receipts
func (ec *EthConfirmer) ResumePendingTaskRuns(ctx context.Context, head *evmtypes.Head) error {

	receiptsPlus, err := ec.orm.FindEthReceiptsPendingConfirmation(ctx, head.Number, ec.chainID)

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
		if data.FailOnRevert && data.Receipt.Status == 0 {
			taskErr = errors.Errorf("transaction %s reverted on-chain", data.Receipt.TxHash)
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
func observeUntilTxConfirmed(chainID big.Int, attempts []EthTxAttempt, receipts []evmtypes.Receipt) {
	for _, attempt := range attempts {
		for _, r := range receipts {
			if attempt.Hash != r.TxHash {
				continue
			}

			// We estimate the time until confirmation by subtracting from the time the eth tx (not the attempt)
			// was created. We want to measure the amount of time taken from when a transaction is created
			// via e.g Txm.CreateTransaction to when it is confirmed on-chain, regardless of how many attempts
			// were needed to achieve this.
			duration := time.Since(attempt.EthTx.CreatedAt)
			promTimeUntilTxConfirmed.
				WithLabelValues(chainID.String()).
				Observe(float64(duration))

			// Since a eth tx can have many attempts, we take the number of blocks to confirm as the block number
			// of the receipt minus the block number of the first ever broadcast for this transaction.
			broadcastBefore := utils.MinKey(attempt.EthTx.EthTxAttempts, func(attempt EthTxAttempt) int64 {
				if attempt.BroadcastBeforeBlockNum != nil {
					return *attempt.BroadcastBeforeBlockNum
				}
				return 0
			})
			if broadcastBefore > 0 {
				blocksElapsed := r.BlockNumber.Int64() - broadcastBefore
				promBlocksUntilTxConfirmed.
					WithLabelValues(chainID.String()).
					Observe(float64(blocksElapsed))
			}
		}
	}
}
