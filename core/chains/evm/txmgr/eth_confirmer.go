package txmgr

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
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
	promTxAttemptCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tx_manager_tx_attempt_count",
		Help: "The number of transaction attempts that are currently being processed by the transaction manager",
	}, []string{"evmChainID"})
)

// EthConfirmer is a broad service which performs four different tasks in sequence on every new longest chain
// Step 1: Mark that all currently pending transaction attempts were broadcast before this block
// Step 2: Check pending transactions for receipts
// Step 3: See if any transactions have exceeded the gas bumping block threshold and, if so, bump them
// Step 4: Check confirmed transactions to make sure they are still in the longest chain (reorg protection)
type EthConfirmer struct {
	utils.StartStopOnce

	lggr      logger.Logger
	db        *sqlx.DB
	q         pg.Q
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
func NewEthConfirmer(db *sqlx.DB, ethClient evmclient.Client, config Config, keystore KeyStore,
	keyStates []ethkey.State, estimator gas.Estimator, resumeCallback ResumeCallback, lggr logger.Logger) *EthConfirmer {

	context, cancel := context.WithCancel(context.Background())
	lggr = lggr.Named("EthConfirmer")
	q := pg.NewQ(db, lggr, config)

	return &EthConfirmer{
		utils.StartStopOnce{},
		lggr,
		db,
		q,
		ethClient,
		ChainKeyStore{
			*ethClient.ChainID(),
			config,
			keystore,
		},
		estimator,
		resumeCallback,
		keyStates,
		utils.NewMailbox[*evmtypes.Head](1),
		context,
		cancel,
		sync.WaitGroup{},
		0,
	}
}

// Start is a comment to appease the linter
func (ec *EthConfirmer) Start() error {
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

	ec.lggr.Debugw("processHead", "headNum", head.Number, "id", "eth_confirmer")

	if err := ec.SetBroadcastBeforeBlockNum(head.Number); err != nil {
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

	return nil
}

// SetBroadcastBeforeBlockNum updates already broadcast attempts with the
// current block number. This is safe no matter how old the head is because if
// the attempt is already broadcast it _must_ have been before this head.
func (ec *EthConfirmer) SetBroadcastBeforeBlockNum(blockNum int64) error {
	_, err := ec.q.Exec(
		`UPDATE eth_tx_attempts
SET broadcast_before_block_num = $1 
FROM eth_txes
WHERE eth_tx_attempts.broadcast_before_block_num IS NULL AND eth_tx_attempts.state = 'broadcast'
AND eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.evm_chain_id = $2`,
		blockNum, ec.chainID.String(),
	)
	return errors.Wrap(err, "SetBroadcastBeforeBlockNum failed")
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
	var attempts []EthTxAttempt
	err = ec.q.Select(&attempts,
		`SELECT DISTINCT ON (eth_tx_id) eth_tx_attempts.*
		FROM eth_tx_attempts
		JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state = 'confirmed_missing_receipt'
		WHERE evm_chain_id = $1
		ORDER BY eth_tx_attempts.eth_tx_id ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC`,
		ec.chainID.String())
	if err != nil {
		return errors.Wrap(err, "CheckConfirmedMissingReceipt failed to query")
	}
	if len(attempts) == 0 {
		return nil
	}
	ec.lggr.Infof("Found %d transactions confirmed_missing_receipt. The RPC node did not give us a receipt for these transactions even though it should have been mined. This could be due to using the wallet with an external account, or if the primary node is not synced or not propagating transactions properly", len(attempts))
	reqs, err := batchSendTransactions(ec.ctx, ec.db, attempts, int(ec.config.EvmRPCDefaultBatchSize()), ec.lggr, ec.ethClient)
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
	_, err = ec.q.Exec(`UPDATE eth_txes SET state='unconfirmed' WHERE id = ANY($1)`, pq.Array(ethTxIDsToUnconfirm))

	if err != nil {
		return errors.Wrap(err, "CheckConfirmedMissingReceipt: marking as unconfirmed failed")
	}
	return
}

// CheckForReceipts finds attempts that are still pending and checks to see if a receipt is present for the given block number
func (ec *EthConfirmer) CheckForReceipts(ctx context.Context, blockNum int64) error {
	attempts, err := ec.findEthTxAttemptsRequiringReceiptFetch()
	if err != nil {
		return errors.Wrap(err, "findEthTxAttemptsRequiringReceiptFetch failed")
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

	if err := ec.markAllConfirmedMissingReceipt(); err != nil {
		return errors.Wrap(err, "unable to mark eth_txes as 'confirmed_missing_receipt'")
	}

	if err := ec.markOldTxesMissingReceiptAsErrored(blockNum); err != nil {
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
	for i := 0; i < len(attempts); i += batchSize {
		j := i + batchSize
		if j > len(attempts) {
			j = len(attempts)
		}

		ec.lggr.Debugw(fmt.Sprintf("Batch fetching receipts at indexes %v until (excluded) %v", i, j), "blockNum", blockNum)

		batch := attempts[i:j]

		receipts, err := ec.batchFetchReceipts(ctx, batch)
		if err != nil {
			return errors.Wrap(err, "batchFetchReceipts failed")
		}
		if err := ec.saveFetchedReceipts(receipts); err != nil {
			return errors.Wrap(err, "saveFetchedReceipts failed")
		}
		promNumConfirmedTxs.WithLabelValues(ec.chainID.String()).Add(float64(len(receipts)))
	}
	return nil
}

func (ec *EthConfirmer) findEthTxAttemptsRequiringReceiptFetch() (attempts []EthTxAttempt, err error) {
	err = ec.q.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&attempts, `
SELECT eth_tx_attempts.* FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt') AND eth_txes.evm_chain_id = $1
WHERE eth_tx_attempts.state != 'insufficient_eth'
ORDER BY eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC
`, ec.chainID.String())
		if err != nil {
			return errors.Wrap(err, "findEthTxAttemptsRequiringReceiptFetch failed to load eth_tx_attempts")
		}
		err = loadEthTxes(tx, attempts)
		return errors.Wrap(err, "findEthTxAttemptsRequiringReceiptFetch failed to load eth_txes")
	}, pg.OptReadOnlyTx())
	return
}

func (ec *EthConfirmer) getMinedTransactionCount(ctx context.Context, from gethCommon.Address) (nonce uint64, err error) {
	return ec.ethClient.NonceAt(ctx, from, nil)
}

// Note this function will increment promRevertedTxCount upon receiving
// a reverted transaction receipt. Should only be called with unconfirmed attempts.
func (ec *EthConfirmer) batchFetchReceipts(ctx context.Context, attempts []EthTxAttempt) (receipts []evmtypes.Receipt, err error) {
	var reqs []rpc.BatchElem
	for _, attempt := range attempts {
		req := rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{attempt.Hash},
			Result: &evmtypes.Receipt{},
		}
		reqs = append(reqs, req)
	}

	lggr := ec.lggr.Named("batchFetchReceipts")

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
			l.Warnf("transaction %s reverted on-chain", receipt.TxHash)
			// This might increment more than once e.g. in case of re-orgs going back and forth we might re-fetch the same receipt
			promRevertedTxCount.WithLabelValues(ec.chainID.String()).Add(1)
		} else {
			promNumSuccessfulTxs.WithLabelValues(ec.chainID.String()).Add(1)
		}

		receipts = append(receipts, *receipt)
	}

	return
}

func (ec *EthConfirmer) saveFetchedReceipts(receipts []evmtypes.Receipt) (err error) {
	if len(receipts) == 0 {
		return nil
	}
	// Notes on this query:
	//
	// # Receipts insert
	// Conflict on (tx_hash, block_hash) shouldn't be possible because there
	// should only ever be one receipt for an eth_tx.
	//
	// ASIDE: This is because we mark confirmed atomically with receipt insert
	// in this query, and delete receipts upon marking unconfirmed - see
	// markForRebroadcast.
	//
	// If a receipt with the same (tx_hash, block_hash) exists then the
	// transaction is marked confirmed which means we _should_ never get here.
	// However, even so, it still shouldn't be an error to upsert a receipt we
	// already have.
	//
	// # EthTxAttempts update
	// It should always be safe to mark the attempt as broadcast here because
	// if it were not successfully broadcast how could it possibly have a
	// receipt?
	//
	// This state is reachable for example if the eth node errors so the
	// attempt was left in_progress but the transaction was actually accepted
	// and mined.
	//
	// # EthTxes update
	// Should be self-explanatory. If we got a receipt, the eth_tx is confirmed.
	//
	var valueStrs []string
	var valueArgs []interface{}
	for _, r := range receipts {
		var receiptJSON []byte
		receiptJSON, err = json.Marshal(r)
		if err != nil {
			return errors.Wrap(err, "saveFetchedReceipts failed to marshal JSON")
		}
		valueStrs = append(valueStrs, "(?,?,?,?,?,NOW())")
		valueArgs = append(valueArgs, r.TxHash, r.BlockHash, r.BlockNumber.Int64(), r.TransactionIndex, receiptJSON)
	}
	valueArgs = append(valueArgs, ec.chainID.String())

	/* #nosec G201 */
	sql := `
	WITH inserted_receipts AS (
		INSERT INTO eth_receipts (tx_hash, block_hash, block_number, transaction_index, receipt, created_at)
		VALUES %s
		ON CONFLICT (tx_hash, block_hash) DO UPDATE SET
			block_number = EXCLUDED.block_number,
			transaction_index = EXCLUDED.transaction_index,
			receipt = EXCLUDED.receipt
		RETURNING eth_receipts.tx_hash, eth_receipts.block_number
	),
	updated_eth_tx_attempts AS (
		UPDATE eth_tx_attempts
		SET
			state = 'broadcast',
			broadcast_before_block_num = COALESCE(eth_tx_attempts.broadcast_before_block_num, inserted_receipts.block_number)
		FROM inserted_receipts
		WHERE inserted_receipts.tx_hash = eth_tx_attempts.hash
		RETURNING eth_tx_attempts.eth_tx_id
	)
	UPDATE eth_txes
	SET state = 'confirmed'
	FROM updated_eth_tx_attempts
	WHERE updated_eth_tx_attempts.eth_tx_id = eth_txes.id
	AND evm_chain_id = ?
	`

	stmt := fmt.Sprintf(sql, strings.Join(valueStrs, ","))

	stmt = sqlx.Rebind(sqlx.DOLLAR, stmt)

	err = ec.q.ExecQ(stmt, valueArgs...)
	return errors.Wrap(err, "saveFetchedReceipts failed to save receipts")
}

// markAllConfirmedMissingReceipt
// It is possible that we can fail to get a receipt for all eth_tx_attempts
// even though a transaction with this nonce has long since been confirmed (we
// know this because transactions with higher nonces HAVE returned a receipt).
//
// This can probably only happen if an external wallet used the account (or
// conceivably because of some bug in the remote eth node that prevents it
// from returning a receipt for a valid transaction).
//
// In this case we mark these transactions as 'confirmed_missing_receipt' to
// prevent gas bumping.
//
// NOTE: We continue to attempt to resend eth_txes in this state on
// every head to guard against the extremely rare scenario of nonce gap due to
// reorg that excludes the transaction (from another wallet) that had this
// nonce (until finality depth is reached, after which we make the explicit
// decision to give up). This is done in the EthResender.
//
// We will continue to try to fetch a receipt for these attempts until all
// attempts are below the finality depth from current head.
func (ec *EthConfirmer) markAllConfirmedMissingReceipt() (err error) {
	res, err := ec.q.Exec(`
UPDATE eth_txes
SET state = 'confirmed_missing_receipt'
WHERE state = 'unconfirmed'
AND nonce < (
	SELECT MAX(nonce) FROM eth_txes
	WHERE state = 'confirmed'
)
AND evm_chain_id = $1
	`, ec.chainID.String())
	if err != nil {
		return errors.Wrap(err, "markAllConfirmedMissingReceipt failed")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "markAllConfirmedMissingReceipt RowsAffected failed")
	}
	if rowsAffected > 0 {
		ec.lggr.Infow(fmt.Sprintf("%d transactions missing receipt", rowsAffected), "n", rowsAffected)
	}
	return
}

// markOldTxesMissingReceiptAsErrored
//
// Once eth_tx has all of its attempts broadcast before some cutoff threshold
// without receiving any receipts, we mark it as fatally errored (never sent).
//
// The job run will also be marked as errored in this case since we never got a
// receipt and thus cannot pass on any transaction hash
func (ec *EthConfirmer) markOldTxesMissingReceiptAsErrored(blockNum int64) error {
	// cutoff is a block height
	// Any 'confirmed_missing_receipt' eth_tx with all attempts older than this block height will be marked as errored
	// We will not try to query for receipts for this transaction any more
	cutoff := blockNum - int64(ec.config.EvmFinalityDepth())
	if cutoff <= 0 {
		return nil
	}

	rows, err := ec.q.Query(`
UPDATE eth_txes
SET state='fatal_error', nonce=NULL, error=$1, broadcast_at=NULL, initial_broadcast_at=NULL
FROM (
	SELECT e1.id, e1.nonce, e1.from_address FROM eth_txes AS e1 WHERE id IN (
		SELECT e2.id FROM eth_txes AS e2
		INNER JOIN eth_tx_attempts ON e2.id = eth_tx_attempts.eth_tx_id
		WHERE e2.state = 'confirmed_missing_receipt'
		AND e2.evm_chain_id = $3
		GROUP BY e2.id
		HAVING max(eth_tx_attempts.broadcast_before_block_num) < $2
	)
	FOR UPDATE OF e1
) e0
WHERE e0.id = eth_txes.id
RETURNING e0.id, e0.nonce, e0.from_address`, ErrCouldNotGetReceipt, cutoff, ec.chainID.String())

	if err != nil {
		return errors.Wrap(err, "markOldTxesMissingReceiptAsErrored failed to query")
	}
	defer ec.lggr.ErrorIfClosing(rows, "eth_txes rows")

	for rows.Next() {
		var ethTxID int64
		var nonce null.Int64
		var fromAddress gethCommon.Address
		if err = rows.Scan(&ethTxID, &nonce, &fromAddress); err != nil {
			return errors.Wrap(err, "error scanning row")
		}

		ec.lggr.Criticalw(fmt.Sprintf("eth_tx with ID %v expired without ever getting a receipt for any of our attempts. "+
			"Current block height is %v. This transaction may not have not been sent and will be marked as fatally errored. "+
			"This can happen if there is another instance of chainlink running that is using the same private key, or if "+
			"an external wallet has been used to send a transaction from account %s with nonce %v."+
			" Please note that Chainlink requires exclusive ownership of it's private keys and sharing keys across multiple"+
			" chainlink instances, or using the chainlink keys with an external wallet is NOT SUPPORTED and WILL lead to missed transactions",
			ethTxID, blockNum, fromAddress.Hex(), nonce.Int64), "ethTxID", ethTxID, "nonce", nonce, "fromAddress", fromAddress)
	}

	return rows.Err()
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
	etxs, err := FindEthTxsRequiringRebroadcast(ctx, ec.q, ec.lggr, address, blockHeight, threshold, bumpDepth, maxInFlightTransactions, ec.chainID)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "FindEthTxsRequiringRebroadcast failed")
	}
	for _, etx := range etxs {
		lggr := etx.GetLogger(ec.lggr)

		attempt, err := ec.attemptForRebroadcast(ctx, lggr, *etx)
		if err != nil {
			return errors.Wrap(err, "attemptForRebroadcast failed")
		}

		lggr.Debugw("Rebroadcasting transaction", "nPreviousAttempts", len(etx.EthTxAttempts), "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)

		if err := ec.saveInProgressAttempt(&attempt); err != nil {
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
	attempts, err := getInProgressEthTxAttempts(ctx, ec.q, ec.lggr, address, ec.chainID)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "getInProgressEthTxAttempts failed")
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

func getInProgressEthTxAttempts(ctx context.Context, q pg.Q, lggr logger.Logger, address gethCommon.Address, chainID big.Int) (attempts []EthTxAttempt, err error) {
	qq := q.WithOpts(pg.WithParentCtx(ctx))
	err = qq.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&attempts, `
SELECT eth_tx_attempts.* FROM eth_tx_attempts
INNER JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state in ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')
WHERE eth_tx_attempts.state = 'in_progress' AND eth_txes.from_address = $1 AND eth_txes.evm_chain_id = $2
`, address, chainID.String())
		if err != nil {
			return errors.Wrap(err, "getInProgressEthTxAttempts failed to load eth_tx_attempts")
		}
		err = loadEthTxes(q, attempts)
		return errors.Wrap(err, "getInProgressEthTxAttempts failed to load eth_txes")
	}, pg.OptReadOnlyTx())
	return attempts, errors.Wrap(err, "getInProgressEthTxAttempts failed")
}

func loadEthTxes(q pg.Queryer, attempts []EthTxAttempt) error {
	ethTxM := make(map[int64]EthTx)
	for _, attempt := range attempts {
		ethTxM[attempt.EthTxID] = EthTx{}
	}
	ethTxIDs := make([]int64, len(ethTxM))
	var i int
	for id := range ethTxM {
		ethTxIDs[i] = id
		i++
	}
	ethTxs := make([]EthTx, len(ethTxIDs))
	if err := q.Select(&ethTxs, `SELECT * FROM eth_txes WHERE id = ANY($1)`, pq.Array(ethTxIDs)); err != nil {
		return errors.Wrap(err, "loadEthTxes failed")
	}
	for _, etx := range ethTxs {
		ethTxM[etx.ID] = etx
	}
	for i, attempt := range attempts {
		attempts[i].EthTx = ethTxM[attempt.EthTxID]
	}
	return nil
}

// FindEthTxsRequiringRebroadcast returns attempts that hit insufficient eth,
// and attempts that need bumping, in nonce ASC order
func FindEthTxsRequiringRebroadcast(ctx context.Context, q pg.Q, lggr logger.Logger, address gethCommon.Address, blockNum, gasBumpThreshold, bumpDepth int64, maxInFlightTransactions uint32, chainID big.Int) (etxs []*EthTx, err error) {
	// NOTE: These two queries could be combined into one using union but it
	// becomes harder to read and difficult to test in isolation. KISS principle
	etxInsufficientEths, err := FindEthTxsRequiringResubmissionDueToInsufficientEth(ctx, q, lggr, address, chainID)
	if ctx.Err() != nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if len(etxInsufficientEths) > 0 {
		lggr.Infow(fmt.Sprintf("Found %d transactions to be re-sent that were previously rejected due to insufficient eth balance", len(etxInsufficientEths)), "blockNum", blockNum, "address", address)
	}

	// TODO: Just pass the Q through everything
	etxBumps, err := FindEthTxsRequiringGasBump(ctx, q, lggr, address, blockNum, gasBumpThreshold, bumpDepth, chainID)
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
			lggr.Warnw("Expected eth_tx for gas bump to have at least one attempt", "etxID", etx.ID, "blockNum", blockNum, "address", address)
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

// FindEthTxsRequiringResubmissionDueToInsufficientEth returns transactions
// that need to be re-sent because they hit an out-of-eth error on a previous
// block
func FindEthTxsRequiringResubmissionDueToInsufficientEth(ctx context.Context, q pg.Q, lggr logger.Logger, address gethCommon.Address, chainID big.Int) (etxs []*EthTx, err error) {
	qq := q.WithOpts(pg.WithParentCtx(ctx))
	err = qq.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&etxs, `
SELECT DISTINCT eth_txes.* FROM eth_txes
INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_tx_attempts.state = 'insufficient_eth'
WHERE eth_txes.from_address = $1 AND eth_txes.state = 'unconfirmed' AND eth_txes.evm_chain_id = $2
ORDER BY nonce ASC
`, address, chainID.String())
		if err != nil {
			return errors.Wrap(err, "FindEthTxsRequiringResubmissionDueToInsufficientEth failed to load eth_txes")
		}

		err = loadEthTxesAttempts(tx, etxs)
		return errors.Wrap(err, "FindEthTxsRequiringResubmissionDueToInsufficientEth failed to load eth_tx_attempts")
	}, pg.OptReadOnlyTx())
	return
}

func loadEthTxesAttempts(q pg.Queryer, etxs []*EthTx) error {
	ethTxIDs := make([]int64, len(etxs))
	ethTxesM := make(map[int64]*EthTx, len(etxs))
	for i, etx := range etxs {
		ethTxIDs[i] = etx.ID
		ethTxesM[etx.ID] = etxs[i]
	}
	var ethTxAttempts []EthTxAttempt
	if err := q.Select(&ethTxAttempts, `SELECT * FROM eth_tx_attempts WHERE eth_tx_id = ANY($1) ORDER BY eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC`, pq.Array(ethTxIDs)); err != nil {
		return errors.Wrap(err, "loadEthTxesAttempts failed to load eth_tx_attempts")
	}
	for _, attempt := range ethTxAttempts {
		etx := ethTxesM[attempt.EthTxID]
		etx.EthTxAttempts = append(etx.EthTxAttempts, attempt)
	}
	return nil
}

// FindEthTxsRequiringGasBump returns transactions that have all
// attempts which are unconfirmed for at least gasBumpThreshold blocks,
// limited by limit pending transactions
//
// It also returns eth_txes that are unconfirmed with no eth_tx_attempts
func FindEthTxsRequiringGasBump(ctx context.Context, q pg.Q, lggr logger.Logger, address gethCommon.Address, blockNum, gasBumpThreshold, depth int64, chainID big.Int) (etxs []*EthTx, err error) {
	if gasBumpThreshold == 0 {
		return
	}
	qq := q.WithOpts(pg.WithParentCtx(ctx))
	err = qq.Transaction(func(tx pg.Queryer) error {
		stmt := `
SELECT eth_txes.* FROM eth_txes
LEFT JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND (broadcast_before_block_num > $4 OR broadcast_before_block_num IS NULL OR eth_tx_attempts.state != 'broadcast')
WHERE eth_txes.state = 'unconfirmed' AND eth_tx_attempts.id IS NULL AND eth_txes.from_address = $1 AND eth_txes.evm_chain_id = $2
	AND (($3 = 0) OR (eth_txes.id IN (SELECT id FROM eth_txes WHERE state = 'unconfirmed' AND from_address = $1 ORDER BY nonce ASC LIMIT $3)))
ORDER BY nonce ASC
`
		if err = tx.Select(&etxs, stmt, address, chainID.String(), depth, blockNum-gasBumpThreshold); err != nil {
			return errors.Wrap(err, "FindEthTxsRequiringGasBump failed to load eth_txes")
		}
		err = loadEthTxesAttempts(tx, etxs)
		return errors.Wrap(err, "FindEthTxsRequiringGasBump failed to load eth_tx_attempts")
	}, pg.OptReadOnlyTx())
	return
}

func (ec *EthConfirmer) attemptForRebroadcast(ctx context.Context, lggr logger.Logger, etx EthTx) (attempt EthTxAttempt, err error) {
	if len(etx.EthTxAttempts) > 0 {
		previousAttempt := etx.EthTxAttempts[0]
		previousAttempt.EthTx = etx
		logFields := ec.logFieldsPreviousAttempt(previousAttempt)
		if previousAttempt.State == EthTxAttemptInsufficientEth {
			// Do not create a new attempt if we ran out of eth last time since bumping gas is pointless
			// Instead try to resubmit the same attempt at the same price, in the hope that the wallet was funded since our last attempt
			lggr.Debugw("Rebroadcast InsufficientEth", logFields...)
			previousAttempt.State = EthTxAttemptInProgress
			return previousAttempt, nil
		}
		attempt, err = ec.bumpGas(previousAttempt)

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

func (ec *EthConfirmer) bumpGas(previousAttempt EthTxAttempt) (bumpedAttempt EthTxAttempt, err error) {
	logFields := ec.logFieldsPreviousAttempt(previousAttempt)
	switch previousAttempt.TxType {
	case 0x0: // Legacy
		var bumpedGasPrice *big.Int
		var bumpedGasLimit uint64
		bumpedGasPrice, bumpedGasLimit, err = ec.estimator.BumpLegacyGas(previousAttempt.GasPrice.ToInt(), previousAttempt.EthTx.GasLimit)
		if err == nil {
			promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
			ec.lggr.Debugw("Rebroadcast bumping gas for Legacy tx", append(logFields, "bumpedGasPrice", bumpedGasPrice.String())...)
			return ec.NewLegacyAttempt(previousAttempt.EthTx, bumpedGasPrice, bumpedGasLimit)
		}
	case 0x2: // EIP1559
		var bumpedFee gas.DynamicFee
		var bumpedGasLimit uint64
		original := previousAttempt.DynamicFee()
		bumpedFee, bumpedGasLimit, err = ec.estimator.BumpDynamicFee(original, previousAttempt.EthTx.GasLimit)
		if err == nil {
			promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
			ec.lggr.Debugw("Rebroadcast bumping gas for DynamicFee tx", append(logFields, "bumpedTipCap", bumpedFee.TipCap.String(), "bumpedFeeCap", bumpedFee.FeeCap.String())...)
			return ec.NewDynamicFeeAttempt(previousAttempt.EthTx, bumpedFee, bumpedGasLimit)
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

// saveInProgressAttempt inserts or updates an attempt
func (ec *EthConfirmer) saveInProgressAttempt(attempt *EthTxAttempt) error {
	if attempt.State != EthTxAttemptInProgress {
		return errors.New("saveInProgressAttempt failed: attempt state must be in_progress")
	}
	// Insert is the usual mode because the attempt is new
	if attempt.ID == 0 {
		query, args, e := ec.q.BindNamed(insertIntoEthTxAttemptsQuery, attempt)
		if e != nil {
			return errors.Wrap(e, "saveInProgressAttempt failed to BindNamed")
		}
		return errors.Wrap(ec.q.Get(attempt, query, args...), "saveInProgressAttempt failed to insert into eth_tx_attempts")
	}
	// Update only applies to case of insufficient eth and simply changes the state to in_progress
	res, err := ec.q.Exec(`UPDATE eth_tx_attempts SET state=$1, broadcast_before_block_num=$2 WHERE id=$3`, attempt.State, attempt.BroadcastBeforeBlockNum, attempt.ID)
	if err != nil {
		return errors.Wrap(err, "saveInProgressAttempt failed to update eth_tx_attempts")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "saveInProgressAttempt failed to get RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrapf(sql.ErrNoRows, "saveInProgressAttempt tried to update eth_tx_attempts but no rows matched id %d", attempt.ID)
	}
	return nil
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
		//
		// It could conceivably happen if the remote eth node changed its configuration.
		replacementAttempt, err := ec.bumpGas(attempt)
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

		if err := saveReplacementInProgressAttempt(ec.q, attempt, &replacementAttempt); err != nil {
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
		lggr.Infow("Transaction temporarily underpriced", "attemptID", attempt.ID, "err", sendError.Error(), "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)
		sendError = nil
	}

	if sendError.IsTooExpensive() {
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
		return deleteInProgressAttempt(ec.q.WithOpts(pg.WithParentCtx(ctx)), attempt)
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
		return deleteInProgressAttempt(ec.q.WithOpts(pg.WithParentCtx(ctx)), attempt)
	}

	if sendError.IsNonceTooLowError() || sendError.IsTransactionAlreadyMined() {
		// Nonce too low indicated that a transaction at this nonce was confirmed already.
		// Mark confirmed_missing_receipt and wait for the next cycle to try to get a receipt
		sendError = nil
		lggr.Debugw("Nonce already used", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash.Hex(), "err", sendError)
		return saveConfirmedMissingReceiptAttempt(ec.q.WithOpts(pg.WithParentCtx(ctx)), ec.lggr, &attempt, now)
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
		lggr.Errorw(fmt.Sprintf("EthTxAttempt %v (hash 0x%x) was rejected due to insufficient eth. "+
			"The eth node returned %s. "+
			"ACTION REQUIRED: Chainlink wallet with address 0x%x is OUT OF FUNDS",
			attempt.ID, attempt.Hash, sendError.Error(), etx.FromAddress,
		), "err", sendError, "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)
		return saveInsufficientEthAttempt(ec.q, ec.lggr, &attempt, now)
	}

	if sendError == nil {
		lggr.Debugw("Successfully broadcast transaction", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash.Hex())
		return saveSentAttempt(ec.db, lggr, &attempt, now)
	}

	// Any other type of error is considered temporary or resolvable by the
	// node operator. The node may have it in the mempool so we must keep the
	// attempt (leave it in_progress). Safest thing to do is bail out and wait
	// for the next head.
	return errors.Wrapf(sendError, "unexpected error sending eth_tx %v with hash %s", etx.ID, attempt.Hash.Hex())
}

func deleteInProgressAttempt(q pg.Q, attempt EthTxAttempt) error {
	if attempt.State != EthTxAttemptInProgress {
		return errors.New("deleteInProgressAttempt: expected attempt state to be in_progress")
	}
	if attempt.ID == 0 {
		return errors.New("deleteInProgressAttempt: expected attempt to have an id")
	}
	_, err := q.Exec(`DELETE FROM eth_tx_attempts WHERE id = $1`, attempt.ID)
	return errors.Wrap(err, "deleteInProgressAttempt failed")
}

func saveConfirmedMissingReceiptAttempt(q pg.Q, lggr logger.Logger, attempt *EthTxAttempt, broadcastAt time.Time) error {
	err := q.Transaction(func(tx pg.Queryer) error {
		if err := saveSentAttempt(tx, lggr, attempt, broadcastAt); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE eth_txes SET state = 'confirmed_missing_receipt' WHERE id = $1`, attempt.EthTxID); err != nil {
			return errors.Wrap(err, "failed to update eth_txes")
		}
		return nil
	})
	return errors.Wrap(err, "saveConfirmedMissingReceiptAttempt failed")
}

func saveSentAttempt(q pg.Queryer, lggr logger.Logger, attempt *EthTxAttempt, broadcastAt time.Time) error {
	if attempt.State != EthTxAttemptInProgress {
		return errors.New("expected state to be in_progress")
	}
	attempt.State = EthTxAttemptBroadcast
	return errors.Wrap(saveAttemptWithNewState(q, lggr, *attempt, broadcastAt), "saveSentAttempt failed")
}

func saveInsufficientEthAttempt(q pg.Queryer, lggr logger.Logger, attempt *EthTxAttempt, broadcastAt time.Time) error {
	if !(attempt.State == EthTxAttemptInProgress || attempt.State == EthTxAttemptInsufficientEth) {
		return errors.New("expected state to be either in_progress or insufficient_eth")
	}
	attempt.State = EthTxAttemptInsufficientEth
	return errors.Wrap(saveAttemptWithNewState(q, lggr, *attempt, broadcastAt), "saveInsufficientEthAttempt failed")

}

func saveAttemptWithNewState(q pg.Queryer, lggr logger.Logger, attempt EthTxAttempt, broadcastAt time.Time) error {
	return pg.SqlxTransactionWithDefaultCtx(q, lggr, func(tx pg.Queryer) error {
		// In case of null broadcast_at (shouldn't happen) we don't want to
		// update anyway because it indicates a state where broadcast_at makes
		// no sense e.g. fatal_error
		if _, err := tx.Exec(`UPDATE eth_txes SET broadcast_at = $1 WHERE id = $2 AND broadcast_at < $1`, broadcastAt, attempt.EthTxID); err != nil {
			return errors.Wrap(err, "saveAttemptWithNewState failed to update eth_txes")
		}
		_, err := tx.Exec(`UPDATE eth_tx_attempts SET state=$1 WHERE id=$2`, attempt.State, attempt.ID)
		return errors.Wrap(err, "saveAttemptWithNewState failed to update eth_tx_attempts")
	})
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
			"evmChainID", ec.chainID.String(), "chainLength", head.ChainLength(), "evmFinalityDepth", ec.config.EvmFinalityDepth(),
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
	etxs, err := findTransactionsConfirmedInBlockRange(ec.q, ec.lggr, head.Number, head.EarliestInChain().Number, ec.chainID)
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

func findTransactionsConfirmedInBlockRange(q pg.Q, lggr logger.Logger, highBlockNumber, lowBlockNumber int64, chainID big.Int) (etxs []*EthTx, err error) {
	err = q.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&etxs, `
SELECT DISTINCT eth_txes.* FROM eth_txes
INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_tx_attempts.state = 'broadcast'
INNER JOIN eth_receipts ON eth_receipts.tx_hash = eth_tx_attempts.hash
WHERE eth_txes.state IN ('confirmed', 'confirmed_missing_receipt') AND block_number BETWEEN $1 AND $2 AND evm_chain_id = $3
ORDER BY nonce ASC
`, lowBlockNumber, highBlockNumber, chainID.String())
		if err != nil {
			return errors.Wrap(err, "findTransactionsConfirmedInBlockRange failed to load eth_txes")
		}
		if err = loadEthTxesAttempts(tx, etxs); err != nil {
			return errors.Wrap(err, "findTransactionsConfirmedInBlockRange failed to load eth_tx_attempts")
		}
		err = loadEthTxesAttemptsReceipts(tx, etxs)
		return errors.Wrap(err, "findTransactionsConfirmedInBlockRange failed to load eth_receipts")
	}, pg.OptReadOnlyTx())
	return etxs, errors.Wrap(err, "findTransactionsConfirmedInBlockRange failed")
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
	err := ec.q.Transaction(func(tx pg.Queryer) error {
		if err := deleteAllReceipts(tx, etx.ID); err != nil {
			return errors.Wrapf(err, "deleteAllReceipts failed for etx %v", etx.ID)
		}
		if err := unconfirmEthTx(tx, etx); err != nil {
			return errors.Wrapf(err, "unconfirmEthTx failed for etx %v", etx.ID)
		}
		return unbroadcastAttempt(tx, attempt)
	})
	return errors.Wrap(err, "markForRebroadcast failed")
}

func deleteAllReceipts(q pg.Queryer, etxID int64) (err error) {
	_, err = q.Exec(`
DELETE FROM eth_receipts
USING eth_tx_attempts
WHERE eth_receipts.tx_hash = eth_tx_attempts.hash
AND eth_tx_attempts.eth_tx_id = $1
	`, etxID)
	return errors.Wrap(err, "deleteAllReceipts failed")
}

func unconfirmEthTx(q pg.Queryer, etx EthTx) error {
	if etx.State != EthTxConfirmed {
		return errors.New("expected eth_tx state to be confirmed")
	}
	_, err := q.Exec(`UPDATE eth_txes SET state = 'unconfirmed' WHERE id = $1`, etx.ID)
	return errors.Wrap(err, "unconfirmEthTx failed")
}

func unbroadcastAttempt(q pg.Queryer, attempt EthTxAttempt) error {
	if attempt.State != EthTxAttemptBroadcast {
		return errors.New("expected eth_tx_attempt to be broadcast")
	}
	_, err := q.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = NULL, state = 'in_progress' WHERE id = $1`, attempt.ID)
	return errors.Wrap(err, "unbroadcastAttempt failed")
}

// ForceRebroadcast sends a transaction for every nonce in the given nonce range at the given gas price.
// If an eth_tx exists for this nonce, we re-send the existing eth_tx with the supplied parameters.
// If an eth_tx doesn't exist for this nonce, we send a zero transaction.
// This operates completely orthogonal to the normal EthConfirmer and can result in untracked attempts!
// Only for emergency usage.
// This is in case of some unforeseen scenario where the node is refusing to release the lock. KISS.
func (ec *EthConfirmer) ForceRebroadcast(beginningNonce uint, endingNonce uint, gasPriceWei uint64, address gethCommon.Address, overrideGasLimit uint64) error {
	ec.lggr.Infof("ForceRebroadcast: will rebroadcast transactions for all nonces between %v and %v", beginningNonce, endingNonce)

	for n := beginningNonce; n <= endingNonce; n++ {
		etx, err := findEthTxWithNonce(ec.q, ec.lggr, address, n)
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
			attempt, err := ec.NewLegacyAttempt(*etx, big.NewInt(int64(gasPriceWei)), etx.GasLimit)
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

func (ec *EthConfirmer) sendEmptyTransaction(ctx context.Context, fromAddress gethCommon.Address, nonce uint, overrideGasLimit uint64, gasPriceWei uint64) (gethCommon.Hash, error) {
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

// findEthTxWithNonce returns any broadcast ethtx with the given nonce
func findEthTxWithNonce(q pg.Q, lggr logger.Logger, fromAddress gethCommon.Address, nonce uint) (etx *EthTx, err error) {
	etx = new(EthTx)
	err = q.Transaction(func(tx pg.Queryer) error {
		err = tx.Get(etx, `
SELECT * FROM eth_txes WHERE from_address = $1 AND nonce = $2 AND state IN ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')
`, fromAddress, nonce)
		if err != nil {
			return errors.Wrap(err, "findEthTxWithNonce failed to load eth_txes")
		}
		err = loadEthTxAttempts(tx, etx)
		return errors.Wrap(err, "findEthTxWithNonce failed to load eth_tx_attempts")
	}, pg.OptReadOnlyTx())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}

// ResumePendingTaskRuns issues callbacks to task runs that are pending waiting for receipts
func (ec *EthConfirmer) ResumePendingTaskRuns(ctx context.Context, head *evmtypes.Head) error {
	type x struct {
		ID      uuid.UUID
		Receipt []byte
	}
	var receipts []x
	// NOTE: we don't filter on eth_txes.state = 'confirmed', because a transaction with an attached receipt
	// is guaranteed to be confirmed. This results in a slightly better query plan.
	if err := ec.q.Select(&receipts, `
	SELECT pipeline_task_runs.id, eth_receipts.receipt FROM pipeline_task_runs
	INNER JOIN pipeline_runs ON pipeline_runs.id = pipeline_task_runs.pipeline_run_id
	INNER JOIN eth_txes ON eth_txes.pipeline_task_run_id = pipeline_task_runs.id
	INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id
	INNER JOIN eth_receipts ON eth_tx_attempts.hash = eth_receipts.tx_hash
	WHERE pipeline_runs.state = 'suspended' AND eth_receipts.block_number <= ($1 - eth_txes.min_confirmations) AND eth_txes.evm_chain_id = $2
	`, head.Number, ec.chainID.String()); err != nil {
		return err
	}

	for _, data := range receipts {
		if err := ec.resumeCallback(data.ID, data.Receipt, nil); err != nil {
			return err
		}
	}

	return nil
}
