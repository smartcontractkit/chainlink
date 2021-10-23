package bulletprooftxmanager

import (
	"context"
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
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/multierr"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/gas"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/static"
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

	logger    logger.Logger
	db        *gorm.DB
	ethClient eth.Client
	ChainKeyStore
	estimator      gas.Estimator
	resumeCallback ResumeCallback

	keyStates []ethkey.State

	mb        *utils.Mailbox
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup

	nConsecutiveBlocksChainTooShort int
}

// NewEthConfirmer instantiates a new eth confirmer
func NewEthConfirmer(db *gorm.DB, ethClient eth.Client, config Config, keystore KeyStore,
	keyStates []ethkey.State, estimator gas.Estimator, resumeCallback ResumeCallback, lggr logger.Logger) *EthConfirmer {

	context, cancel := context.WithCancel(context.Background())
	return &EthConfirmer{
		utils.StartStopOnce{},
		lggr,
		db,
		ethClient,
		ChainKeyStore{
			*ethClient.ChainID(),
			config,
			keystore,
		},
		estimator,
		resumeCallback,
		keyStates,
		utils.NewMailbox(1),
		context,
		cancel,
		sync.WaitGroup{},
		0,
	}
}

func (ec *EthConfirmer) Start() error {
	return ec.StartOnce("EthConfirmer", func() error {
		if ec.config.EvmGasBumpThreshold() == 0 {
			ec.logger.Infow("EthConfirmer: Gas bumping is disabled (ETH_GAS_BUMP_THRESHOLD set to 0)", "ethGasBumpThreshold", 0)
		} else {
			ec.logger.Infow(fmt.Sprintf("EthConfirmer: Gas bumping is enabled, unconfirmed transactions will have their gas price bumped every %d blocks", ec.config.EvmGasBumpThreshold()), "ethGasBumpThreshold", ec.config.EvmGasBumpThreshold())
		}

		ec.wg.Add(1)
		go ec.runLoop()

		return nil
	})
}

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
				h, is := head.(eth.Head)
				if !is {
					ec.logger.Errorf("EthConfirmer: invariant violation, expected %T but got %T", eth.Head{}, head)
					continue
				}
				if err := ec.ProcessHead(ec.ctx, h); err != nil {
					ec.logger.Errorw("EthConfirmer error", "err", err)
					continue
				}
			}
		case <-ec.ctx.Done():
			return
		}
	}
}

// ProcessHead takes all required transactions for the confirmer on a new head
func (ec *EthConfirmer) ProcessHead(ctx context.Context, head eth.Head) error {
	ctx, cancel := context.WithTimeout(ctx, processHeadTimeout)
	defer cancel()

	return ec.processHead(ctx, head)
}

// NOTE: This SHOULD NOT be run concurrently or it could behave badly
func (ec *EthConfirmer) processHead(ctx context.Context, head eth.Head) error {
	mark := time.Now()

	ec.logger.Debugw("EthConfirmer: processHead", "headNum", head.Number, "id", "eth_confirmer")

	if err := ec.SetBroadcastBeforeBlockNum(head.Number); err != nil {
		return errors.Wrap(err, "SetBroadcastBeforeBlockNum failed")
	}

	if err := ec.CheckForReceipts(ctx, head.Number); err != nil {
		return errors.Wrap(err, "CheckForReceipts failed")
	}

	ec.logger.Debugw("EthConfirmer: finished CheckForReceipts", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	if err := ec.RebroadcastWhereNecessary(ctx, head.Number); err != nil {
		return errors.Wrap(err, "RebroadcastWhereNecessary failed")
	}

	ec.logger.Debugw("EthConfirmer: finished RebroadcastWhereNecessary", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	if err := ec.EnsureConfirmedTransactionsInLongestChain(ctx, head); err != nil {
		return errors.Wrap(err, "EnsureConfirmedTransactionsInLongestChain failed")
	}

	ec.logger.Debugw("EthConfirmer: finished EnsureConfirmedTransactionsInLongestChain", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")

	if ec.resumeCallback != nil {
		mark = time.Now()
		if err := ec.ResumePendingTaskRuns(ctx, head); err != nil {
			return errors.Wrap(err, "ResumePendingTaskRuns failed")
		}

		ec.logger.Debugw("EthConfirmer: finished ResumePendingTaskRuns", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	}

	return nil
}

// SetBroadcastBeforeBlockNum updates already broadcast attempts with the
// current block number. This is safe no matter how old the head is because if
// the attempt is already broadcast it _must_ have been before this head.
func (ec *EthConfirmer) SetBroadcastBeforeBlockNum(blockNum int64) error {
	return ec.db.Exec(
		`UPDATE eth_tx_attempts SET broadcast_before_block_num = ? WHERE broadcast_before_block_num IS NULL AND state = 'broadcast'`,
		blockNum,
	).Error
}

func (ec *EthConfirmer) CheckForReceipts(ctx context.Context, blockNum int64) error {
	attempts, err := ec.findEthTxAttemptsRequiringReceiptFetch()
	if err != nil {
		return errors.Wrap(err, "findEthTxAttemptsRequiringReceiptFetch failed")
	}
	if len(attempts) == 0 {
		return nil
	}

	ec.logger.Debugw(fmt.Sprintf("EthConfirmer: fetching receipts for %v transaction attempts", len(attempts)), "blockNum", blockNum)

	attemptsByAddress := make(map[gethCommon.Address][]EthTxAttempt)
	for _, att := range attempts {
		attemptsByAddress[att.EthTx.FromAddress] = append(attemptsByAddress[att.EthTx.FromAddress], att)
	}

	for from, attempts := range attemptsByAddress {
		ctxInner, cancel := eth.DefaultQueryCtx(ctx)
		latestBlockNonce, err := ec.getNonceForLatestBlock(ctxInner, from)

		if ctxInner.Err() != nil { // timeout
			return errors.Wrapf(ctxInner.Err(), "unable to fetch pending nonce for address: %v - timeout or interrupt", from)
		}
		cancel()
		if err != nil {
			return errors.Wrapf(err, "unable to fetch pending nonce for address: %v", from)
		}

		likelyConfirmed := ec.separateLikelyConfirmedAttempts(from, attempts, latestBlockNonce)
		likelyConfirmedCount := len(likelyConfirmed)
		if likelyConfirmedCount > 0 {
			likelyUnconfirmedCount := len(attempts) - likelyConfirmedCount

			ec.logger.Debugf("EthConfirmer: Fetching and saving %v likely confirmed receipts. Skipping checking the others (%v)",
				likelyConfirmedCount, likelyUnconfirmedCount)

			start := time.Now()
			err = ec.fetchAndSaveReceipts(ctx, likelyConfirmed, blockNum)
			if err != nil {
				return errors.Wrapf(err, "unable to fetch and save receipts for likely confirmed txs, for address: %v", from)
			}
			ec.logger.Debugw(fmt.Sprintf("EthConfirmer: Fetching and saving %v likely confirmed receipts done", likelyConfirmedCount),
				"time", time.Since(start))
		}
	}

	if err := ec.markConfirmedMissingReceipt(); err != nil {
		return errors.Wrap(err, "unable to mark eth_txes as 'confirmed_missing_receipt'")
	}

	if err := ec.markOldTxesMissingReceiptAsErrored(blockNum); err != nil {
		return errors.Wrap(err, "unable to confirm buried unconfirmed eth_txes")
	}
	return nil
}

func (ec *EthConfirmer) separateLikelyConfirmedAttempts(from gethCommon.Address, attempts []EthTxAttempt, latestBlockNonce uint64) []EthTxAttempt {
	if len(attempts) == 0 {
		return attempts
	}

	firstAttemptNonce := fmt.Sprintf("%v", *attempts[len(attempts)-1].EthTx.Nonce)
	lastAttemptNonce := fmt.Sprintf("%v", *attempts[0].EthTx.Nonce)
	ec.logger.Debugw(fmt.Sprintf("EthConfirmer: There are %v attempts from address %v, latest nonce for it is %v and for the attempts' nonces: first = %v, last = %v",
		len(attempts), from, latestBlockNonce, firstAttemptNonce, lastAttemptNonce), "nAttempts", len(attempts), "fromAddress", from, "latestBlockNonce", latestBlockNonce, "firstAttemptNonce", firstAttemptNonce, "lastAttemptNonce", lastAttemptNonce)

	likelyConfirmed := attempts
	for i := 0; i < len(attempts); i++ {
		if attempts[i].EthTx.Nonce != nil && *attempts[i].EthTx.Nonce > int64(latestBlockNonce) {
			ec.logger.Debugf("EthConfirmer: Marking attempts as likely confirmed just before index %v, at nonce: %v", i, *attempts[i].EthTx.Nonce)
			likelyConfirmed = attempts[0:i]
			break
		}
	}

	if len(likelyConfirmed) == 0 {
		ec.logger.Debug("EthConfirmer: There are no likely confirmed attempts - so will skip checking any")
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

		ec.logger.Debugw(fmt.Sprintf("EthConfirmer: batch fetching receipts at indexes %v until (excluded) %v", i, j), "blockNum", blockNum)

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
	err = ec.db.
		Joins("EthTx"). // Joins("EthTx") is needed for the query to actually return data from eth_txes table as well.
		Joins("JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt') AND eth_txes.evm_chain_id = ?", ec.chainID.String()).
		Order("eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC").
		Where("eth_tx_attempts.state != 'insufficient_eth'").
		Find(&attempts).Error

	return
}

func (ec *EthConfirmer) getNonceForLatestBlock(ctx context.Context, from gethCommon.Address) (nonce uint64, err error) {
	return ec.ethClient.NonceAt(ctx, from, nil)
}

// Note this function will increment promRevertedTxCount upon receiving
// a reverted transaction receipt. Should only be called with unconfirmed attempts.
func (ec *EthConfirmer) batchFetchReceipts(ctx context.Context, attempts []EthTxAttempt) (receipts []Receipt, err error) {
	var reqs []rpc.BatchElem
	for _, attempt := range attempts {
		req := rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{attempt.Hash},
			Result: &Receipt{},
		}
		reqs = append(reqs, req)
	}

	ctx, cancel := eth.DefaultQueryCtx(ctx)
	defer cancel()

	err = ec.ethClient.BatchCallContext(ctx, reqs)
	if err != nil {
		return nil, errors.Wrap(err, "EthConfirmer#batchFetchReceipts error fetching receipts with BatchCallContext")
	}

	for i, req := range reqs {
		attempt := attempts[i]
		result, err := req.Result, req.Error

		receipt, is := result.(*Receipt)
		if !is {
			return nil, errors.Errorf("expected result to be a %T, got %T", (*Receipt)(nil), result)
		}

		l := ec.logger.With(
			"txHash", attempt.Hash.Hex(), "ethTxAttemptID", attempt.ID, "ethTxID", attempt.EthTxID, "err", err,
		)

		if err != nil {
			l.Errorw("EthConfirmer#batchFetchReceipts: fetchReceipt failed")
			continue
		}

		if receipt == nil {
			// NOTE: This should never possibly happen, but it seems safer to
			// check regardless to avoid a potential panic
			l.Errorw("EthConfirmer#batchFetchReceipts: invariant violation, got nil receipt")
			continue
		}

		if receipt.IsZero() {
			l.Debugw("EthConfirmer#batchFetchReceipts: still waiting for receipt")
			continue
		}

		l = l.With("blockHash", receipt.BlockHash.Hex(), "status", receipt.Status, "transactionIndex", receipt.TransactionIndex)

		if receipt.IsUnmined() {
			l.Debugw("EthConfirmer#batchFetchReceipts: got receipt for transaction but it's still in the mempool and not included in a block yet")
			continue
		}

		l.Debugw("EthConfirmer#batchFetchReceipts: got receipt for transaction", "blockNumber", receipt.BlockNumber, "gasUsed", receipt.GasUsed)

		if receipt.TxHash != attempt.Hash {
			l.Errorf("EthConfirmer#batchFetchReceipts: invariant violation, expected receipt with hash %s to have same hash as attempt with hash %s", receipt.TxHash.Hex(), attempt.Hash.Hex())
			continue
		}

		if receipt.BlockNumber == nil {
			l.Error("EthConfirmer#batchFetchReceipts: invariant violation, receipt was missing block number")
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

func (ec *EthConfirmer) saveFetchedReceipts(receipts []Receipt) (err error) {
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

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()

	err = ec.db.WithContext(ctx).Exec(stmt, valueArgs...).Error
	return errors.Wrap(err, "saveFetchedReceipts failed to save receipts")
}

// markConfirmedMissingReceipt
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
func (ec *EthConfirmer) markConfirmedMissingReceipt() (err error) {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()

	res := ec.db.WithContext(ctx).Exec(`
UPDATE eth_txes
SET state = 'confirmed_missing_receipt'
WHERE state = 'unconfirmed'
AND nonce < (
	SELECT MAX(nonce) FROM eth_txes
	WHERE state = 'confirmed'
)
AND evm_chain_id = ?
	`, ec.chainID.String())
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		ec.logger.Infow(fmt.Sprintf("EthConfirmer: %d transactions missing receipt", res.RowsAffected), "n", res.RowsAffected)
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
	d, err := ec.db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()

	rows, err := d.QueryContext(ctx, `
UPDATE eth_txes
SET state='fatal_error', nonce=NULL, error=$1, broadcast_at=NULL
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
	defer ec.logger.ErrorIfCalling(rows.Close)

	for rows.Next() {
		var ethTxID int64
		var nonce null.Int64
		var fromAddress gethCommon.Address
		if err = rows.Scan(&ethTxID, &nonce, &fromAddress); err != nil {
			return errors.Wrap(err, "error scanning row")
		}

		ec.logger.Errorw(fmt.Sprintf("EthConfirmer: eth_tx with ID %v expired without ever getting a receipt for any of our attempts. "+
			"Current block height is %v. This transaction has not been sent and will be marked as fatally errored. "+
			"This can happen if there is another instance of chainlink running that is using the same private key, or if "+
			"an external wallet has been used to send a transaction from account %s with nonce %v."+
			" Please note that Chainlink requires exclusive ownership of it's private keys and sharing keys across multiple"+
			" chainlink instances, or using the chainlink keys with an external wallet is NOT SUPPORTED and WILL lead to missed transactions",
			ethTxID, blockNum, fromAddress.Hex(), nonce.Int64), "ethTxID", ethTxID, "nonce", nonce, "fromAddress", fromAddress)
	}

	return nil
}

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
				ec.logger.Errorw("Error in RebroadcastWhereNecessary", "error", err, "fromAddress", fromAddress)
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
	etxs, err := FindEthTxsRequiringRebroadcast(ctx, ec.db, address, blockHeight, threshold, bumpDepth, maxInFlightTransactions, ec.logger, ec.chainID)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "FindEthTxsRequiringRebroadcast failed")
	}
	for _, etx := range etxs {
		attempt, err := ec.attemptForRebroadcast(ctx, etx)
		if err != nil {
			return errors.Wrap(err, "attemptForRebroadcast failed")
		}

		ec.logger.Debugw("EthConfirmer: Rebroadcasting transaction", "ethTxID", etx.ID, "nonce", etx.Nonce, "nPreviousAttempts", len(etx.EthTxAttempts), "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)

		if err := ec.saveInProgressAttempt(&attempt); err != nil {
			return errors.Wrap(err, "saveInProgressAttempt failed")
		}

		if err := ec.handleInProgressAttempt(ctx, etx, attempt, blockHeight); err != nil {
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
	attempts, err := getInProgressEthTxAttempts(ctx, ec.db, address, ec.chainID)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "getInProgressEthTxAttempts failed")
	}
	for _, a := range attempts {
		err := ec.handleInProgressAttempt(ctx, a.EthTx, a, blockHeight)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			return errors.Wrap(err, "handleInProgressAttempt failed")
		}
	}
	return nil
}

func getInProgressEthTxAttempts(ctx context.Context, db *gorm.DB, address gethCommon.Address, chainID big.Int) ([]EthTxAttempt, error) {
	var attempts []EthTxAttempt
	err := db.
		WithContext(ctx).
		Preload("EthTx").
		Joins("INNER JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state in ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')").
		Where("eth_tx_attempts.state = 'in_progress' AND eth_txes.from_address = ? AND eth_txes.evm_chain_id = ?", address, chainID.String()).
		Find(&attempts).Error
	return attempts, errors.Wrap(err, "getInProgressEthTxAttempts failed")
}

// FindEthTxsRequiringRebroadcast returns attempts that hit insufficient eth,
// and attempts that need bumping, in nonce ASC order
func FindEthTxsRequiringRebroadcast(ctx context.Context, db *gorm.DB, address gethCommon.Address, blockNum, gasBumpThreshold, bumpDepth int64, maxInFlightTransactions uint32, l logger.Logger, chainID big.Int) (etxs []EthTx, err error) {
	// NOTE: These two queries could be combined into one using union but it
	// becomes harder to read and difficult to test in isolation. KISS principle
	etxInsufficientEths, err := FindEthTxsRequiringResubmissionDueToInsufficientEth(ctx, db, address, chainID)
	if ctx.Err() != nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if len(etxInsufficientEths) > 0 {
		l.Infow(fmt.Sprintf("EthConfirmer: Found %d transactions to be re-sent that were previously rejected due to insufficient eth balance", len(etxInsufficientEths)), "blockNum", blockNum, "address", address)
	}

	etxBumps, err := FindEthTxsRequiringGasBump(ctx, db, address, blockNum, gasBumpThreshold, bumpDepth, chainID)
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
			l.Warnw("EthConfirmer: expected eth_tx for gas bump to have at least one attempt", "etxID", etx.ID, "blockNum", blockNum, "address", address)
		}
		l.Infow(fmt.Sprintf("EthConfirmer: Found %d transactions to re-sent that have still not been confirmed after at least %d blocks. The oldest of these has not still not been confirmed after %d blocks. These transactions will have their gas price bumped. %s", len(etxBumps), gasBumpThreshold, oldestBlocksBehind, static.EthNodeConnectivityProblemLabel), "blockNum", blockNum, "address", address, "gasBumpThreshold", gasBumpThreshold)
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
		l.Warnf("EthConfirmer: %d transactions to rebroadcast which exceeds limit of %d. %s", len(etxs), maxInFlightTransactions, static.EvmMaxInFlightTransactionsWarningLabel)
		etxs = etxs[:maxInFlightTransactions]
	}

	return
}

// FindEthTxsRequiringResubmissionDueToInsufficientEth returns transactions
// that need to be re-sent because they hit an out-of-eth error on a previous
// block
func FindEthTxsRequiringResubmissionDueToInsufficientEth(ctx context.Context, db *gorm.DB, address gethCommon.Address, chainID big.Int) (etxs []EthTx, err error) {
	err = db.
		WithContext(ctx).
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC")
		}).
		Joins("INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_tx_attempts.state = 'insufficient_eth'").
		Where("eth_txes.from_address = ? AND eth_txes.state = 'unconfirmed' AND eth_txes.evm_chain_id = ?", address, chainID.String()).
		Order("nonce ASC").
		Find(&etxs).Error

	err = errors.Wrap(err, "FindEthTxsRequiringResubmissionDueToInsufficientEth failed to load eth_txes having insufficient eth")

	return
}

// FindEthTxsRequiringGasBump returns transactions that have all
// attempts which are unconfirmed for at least gasBumpThreshold blocks,
// limited by limit pending transactions
//
// It also returns eth_txes that are unconfirmed with no eth_tx_attempts
func FindEthTxsRequiringGasBump(ctx context.Context, db *gorm.DB, address gethCommon.Address, blockNum, gasBumpThreshold, depth int64, chainID big.Int) (etxs []EthTx, err error) {
	if gasBumpThreshold == 0 {
		return
	}
	q := db.
		WithContext(ctx).
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC")
		}).
		Joins("LEFT JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id "+
			"AND (broadcast_before_block_num > ? OR broadcast_before_block_num IS NULL OR eth_tx_attempts.state != 'broadcast')", blockNum-gasBumpThreshold).
		Where("eth_txes.state = 'unconfirmed' AND eth_tx_attempts.id IS NULL AND eth_txes.from_address = ? AND eth_txes.evm_chain_id = ?", address, chainID.String())

	if depth > 0 {
		q = q.Where("eth_txes.id IN (SELECT id FROM eth_txes WHERE state = 'unconfirmed' AND from_address = ? ORDER BY nonce ASC LIMIT ?)", address, depth)
	}

	err = q.Order("nonce ASC").Find(&etxs).Error
	err = errors.Wrap(err, "FindEthTxsRequiringGasBump failed to load eth_txes requiring gas bump")

	return
}

func (ec *EthConfirmer) attemptForRebroadcast(ctx context.Context, etx EthTx) (attempt EthTxAttempt, err error) {
	if len(etx.EthTxAttempts) > 0 {
		previousAttempt := etx.EthTxAttempts[0]
		previousAttempt.EthTx = etx
		logFields := ec.logFieldsPreviousAttempt(previousAttempt)
		if previousAttempt.State == EthTxAttemptInsufficientEth {
			// Do not create a new attempt if we ran out of eth last time since bumping gas is pointless
			// Instead try to resubmit the same attempt at the same price, in the hope that the wallet was funded since our last attempt
			ec.logger.Debugw("EthConfirmer: rebroadcast InsufficientEth", logFields...)
			previousAttempt.State = EthTxAttemptInProgress
			return previousAttempt, nil
		}
		attempt, err = ec.bumpGas(previousAttempt)

		if gas.IsBumpErr(err) {
			ec.logger.Errorw("Failed to bump gas", append(logFields, "err", err)...)
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
	case 0x0:
		var bumpedGasPrice *big.Int
		var bumpedGasLimit uint64
		bumpedGasPrice, bumpedGasLimit, err = ec.estimator.BumpLegacyGas(previousAttempt.GasPrice.ToInt(), previousAttempt.EthTx.GasLimit)
		if err == nil {
			promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
			ec.logger.Debugw("EthConfirmer: rebroadcast bumping gas for Legacy tx", append(logFields, "bumpedGasPrice", bumpedGasPrice.String())...)
			return ec.NewLegacyAttempt(previousAttempt.EthTx, bumpedGasPrice, bumpedGasLimit)
		}
	case 0x2:
		// BumpDynamicFee(original DynamicFee, gasLimit uint64) (bumped DynamicFee, chainSpecificGasLimit uint64, err error)
		var bumpedFee gas.DynamicFee
		var bumpedGasLimit uint64
		original := previousAttempt.DynamicFee()
		bumpedFee, bumpedGasLimit, err = ec.estimator.BumpDynamicFee(original, previousAttempt.EthTx.GasLimit)
		if err == nil {
			promNumGasBumps.WithLabelValues(ec.chainID.String()).Inc()
			ec.logger.Debugw("EthConfirmer: rebroadcast bumping gas for DynamicFee tx", append(logFields, "bumpedTipCap", bumpedFee.TipCap.String(), "bumpedFeeCap", bumpedFee.FeeCap.String())...)
			return ec.NewDynamicFeeAttempt(previousAttempt.EthTx, bumpedFee, bumpedGasLimit)
		}
	default:
		err = errors.Errorf("invariant violation: Attempt %v had unrecognised transaction type %v"+
			"This is a bug! Please report to https://github.com/smartcontractkit/chainlink/issues", previousAttempt.ID, previousAttempt.TxType)
	}

	if errors.Cause(err) == gas.ErrBumpGasExceedsLimit {
		promGasBumpExceedsLimit.WithLabelValues(ec.chainID.String()).Inc()
	}

	return bumpedAttempt, errors.Wrap(err, "error bumping gas")
}

func (ec *EthConfirmer) saveInProgressAttempt(attempt *EthTxAttempt) error {
	if attempt.State != EthTxAttemptInProgress {
		return errors.New("saveInProgressAttempt failed: attempt state must be in_progress")
	}
	return errors.Wrap(ec.db.Save(attempt).Error, "saveInProgressAttempt failed")
}

func (ec *EthConfirmer) handleInProgressAttempt(ctx context.Context, etx EthTx, attempt EthTxAttempt, blockHeight int64) error {
	if attempt.State != EthTxAttemptInProgress {
		return errors.Errorf("invariant violation: expected eth_tx_attempt %v to be in_progress, it was %s", attempt.ID, attempt.State)
	}

	now := time.Now()
	sendError := sendTransaction(ctx, ec.ethClient, attempt, etx, ec.logger)

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
		ec.logger.Errorw(
			fmt.Sprintf("gas price %v wei was rejected by the eth node for being too low. Eth node returned: '%s'",
				attempt.GasPrice.String(), sendError.Error()),
			"previousAttempt", attempt, "replacementAttempt", replacementAttempt)

		if err := saveReplacementInProgressAttempt(ec.db, attempt, &replacementAttempt); err != nil {
			return errors.Wrap(err, "saveReplacementInProgressAttempt failed")
		}
		return ec.handleInProgressAttempt(ctx, etx, replacementAttempt, blockHeight)
	}

	if sendError.IsTemporarilyUnderpriced() {
		// Most likely scenario here is a parity node that is rejecting
		// low-priced transactions due to mempool pressure
		//
		// In that case, the safest thing to do is to pretend the transaction
		// was accepted and continue the normal gas bumping cycle until we can
		// get it into the mempool
		ec.logger.Infow("EthConfirmer: Transaction temporarily underpriced", "ethTxID", etx.ID, "attemptID", attempt.ID, "err", sendError.Error(), "gasPrice", attempt.GasPrice, "gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)
		sendError = nil
	}

	if sendError.IsTooExpensive() {
		// The gas price was bumped too high. This transaction attempt cannot be accepted.
		//
		// Best thing we can do is to re-send the previous attempt at the old
		// price and discard this bumped version.
		ec.logger.Errorw("EthConfirmer: bumped transaction gas price was rejected by the eth node for being too high. Consider increasing your eth node's RPCTxFeeCap (it is suggested to run geth with no cap i.e. --rpc.gascap=0 --rpc.txfeecap=0)",
			"ethTxID", etx.ID,
			"err", sendError,
			"gasPrice", attempt.GasPrice,
			"gasLimit", etx.GasLimit,
			"signedRawTx", hexutil.Encode(attempt.SignedRawTx),
			"blockHeight", blockHeight,
			"id", "RPCTxFeeCapExceeded",
		)
		return deleteInProgressAttempt(ec.db, attempt)
	}

	if sendError.Fatal() {
		// WARNING: This should never happen!
		// Should NEVER be fatal this is an invariant violation. The
		// EthBroadcaster can never create an EthTxAttempt that will
		// fatally error.
		//
		// The only scenario imaginable where this might take place is if
		// geth/parity have been updated between broadcasting and confirming steps.
		ec.logger.Errorw("invariant violation: fatal error while re-attempting transaction",
			"ethTxID", etx.ID,
			"err", sendError,
			"signedRawTx", hexutil.Encode(attempt.SignedRawTx),
			"blockHeight", blockHeight,
		)
		// This will loop continuously on every new head so it must be handled manually by the node operator!
		return deleteInProgressAttempt(ec.db, attempt)
	}

	if sendError.IsNonceTooLowError() {
		// Nonce too low indicated that a transaction at this nonce was confirmed already.
		// Assume success and hand off to the next cycle to fetch a receipt and mark confirmed.
		sendError = nil
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
		ec.logger.Errorw(fmt.Sprintf("EthConfirmer: replacement transaction underpriced at %v wei for eth_tx %v. "+
			"Eth node returned error: '%s'. "+
			"Either you have set ETH_GAS_BUMP_PERCENT (currently %v%%) too low or an external wallet used this account. "+
			"Please note that using your node's private keys outside of the chainlink node is NOT SUPPORTED and can lead to missed transactions.",
			attempt.GasPrice.ToInt().Int64(), etx.ID, sendError.Error(), ec.config.EvmGasBumpPercent()), "err", sendError)

		// Assume success and hand off to the next cycle.
		sendError = nil
	}

	if sendError.IsInsufficientEth() {
		ec.logger.Errorw(fmt.Sprintf("EthConfirmer: EthTxAttempt %v (hash 0x%x) at gas price (%s Wei) was rejected due to insufficient eth. "+
			"The eth node returned %s. "+
			"ACTION REQUIRED: Chainlink wallet with address 0x%x is OUT OF FUNDS",
			attempt.ID, attempt.Hash, attempt.GasPrice.String(), sendError.Error(), etx.FromAddress,
		), "err", sendError)
		return saveInsufficientEthAttempt(ec.db, &attempt, now)
	}

	if sendError == nil {
		ec.logger.Debugw("EthConfirmer: successfully broadcast transaction", "ethTxID", etx.ID, "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash.Hex())
		return saveSentAttempt(ec.db, &attempt, now)
	}

	// Any other type of error is considered temporary or resolvable by the
	// node operator. The node may have it in the mempool so we must keep the
	// attempt (leave it in_progress). Safest thing to do is bail out and wait
	// for the next head.
	return errors.Wrapf(sendError, "unexpected error sending eth_tx %v with hash %s", etx.ID, attempt.Hash.Hex())
}

func deleteInProgressAttempt(db *gorm.DB, attempt EthTxAttempt) error {
	if attempt.State != EthTxAttemptInProgress {
		return errors.New("deleteInProgressAttempt: expected attempt state to be in_progress")
	}
	if attempt.ID == 0 {
		return errors.New("deleteInProgressAttempt: expected attempt to have an id")
	}
	return errors.Wrap(db.Exec(`DELETE FROM eth_tx_attempts WHERE id = ?`, attempt.ID).Error, "deleteInProgressAttempt failed")
}

func saveSentAttempt(db *gorm.DB, attempt *EthTxAttempt, broadcastAt time.Time) error {
	if attempt.State != EthTxAttemptInProgress {
		return errors.New("expected state to be in_progress")
	}
	attempt.State = EthTxAttemptBroadcast
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	return postgres.GormTransaction(ctx, db, func(tx *gorm.DB) error {
		// In case of null broadcast_at (shouldn't happen) we don't want to
		// update anyway because it indicates a state where broadcast_at makes
		// no sense e.g. fatal_error
		if err := tx.Exec(`UPDATE eth_txes SET broadcast_at = ? WHERE id = ? AND broadcast_at < ?`, broadcastAt, attempt.EthTxID, broadcastAt).Error; err != nil {
			return errors.Wrap(err, "saveSentAttempt failed")
		}
		return errors.Wrap(db.Save(attempt).Error, "saveSentAttempt failed")
	})
}

func saveInsufficientEthAttempt(db *gorm.DB, attempt *EthTxAttempt, broadcastAt time.Time) error {
	if !(attempt.State == EthTxAttemptInProgress || attempt.State == EthTxAttemptInsufficientEth) {
		return errors.New("expected state to be either in_progress or insufficient_eth")
	}
	attempt.State = EthTxAttemptInsufficientEth
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	return postgres.GormTransaction(ctx, db, func(tx *gorm.DB) error {
		// In case of null broadcast_at (shouldn't happen) we don't want to
		// update anyway because it indicates a state where broadcast_at makes
		// no sense e.g. fatal_error
		if err := tx.Exec(`UPDATE eth_txes SET broadcast_at = ? WHERE id = ? AND broadcast_at < ?`, broadcastAt, attempt.EthTxID, broadcastAt).Error; err != nil {
			return errors.Wrap(err, "saveInsufficientEthAttempt failed")
		}
		return errors.Wrap(db.Save(attempt).Error, "saveInsufficientEthAttempt failed")
	})

}

// EnsureConfirmedTransactionsInLongestChain finds all confirmed eth_txes up to the depth
// of the given chain and ensures that every one has a receipt with a block hash that is
// in the given chain.
//
// If any of the confirmed transactions does not have a receipt in the chain, it has been
// re-org'd out and will be rebroadcast.
func (ec *EthConfirmer) EnsureConfirmedTransactionsInLongestChain(ctx context.Context, head eth.Head) error {
	if head.ChainLength() < ec.config.EvmFinalityDepth() {
		logArgs := []interface{}{
			"evmChainID", ec.chainID.String(), "chainLength", head.ChainLength(), "evmFinalityDepth", ec.config.EvmFinalityDepth(),
		}
		if ec.nConsecutiveBlocksChainTooShort > logAfterNConsecutiveBlocksChainTooShort {
			warnMsg := "EthConfirmer: chain length supplied for re-org detection was shorter than EvmFinalityDepth. Re-org protection is not working properly. This could indicate a problem with the remote RPC endpoint, a compatibility issue with a particular blockchain, a bug with this particular blockchain, heads table being truncated too early, remote node out of sync, or something else. If this happens a lot please raise a bug with the Chainlink team including a log output sample and details of the chain and RPC endpoint you are using."
			ec.logger.Warnw(warnMsg, append(logArgs, "nConsecutiveBlocksChainTooShort", ec.nConsecutiveBlocksChainTooShort)...)
		} else {
			logMsg := "EthConfirmer: chain length supplied for re-org detection was shorter than EvmFinalityDepth"
			ec.logger.Debugw(logMsg, append(logArgs, "nConsecutiveBlocksChainTooShort", ec.nConsecutiveBlocksChainTooShort)...)
		}
		ec.nConsecutiveBlocksChainTooShort++
	} else {
		ec.nConsecutiveBlocksChainTooShort = 0
	}
	etxs, err := findTransactionsConfirmedInBlockRange(ec.db, head.Number, head.EarliestInChain().Number, ec.chainID)
	if err != nil {
		return errors.Wrap(err, "findTransactionsConfirmedInBlockRange failed")
	}

	for _, etx := range etxs {
		if !hasReceiptInLongestChain(etx, head) {
			if err := ec.markForRebroadcast(etx, head); err != nil {
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
				ec.logger.Errorw("Error in handleAnyInProgressAttempts", "err", err, "fromAddress", fromAddress)
			}

			wg.Done()
		}(key.Address.Address())
	}

	wg.Wait()

	return multierr.Combine(errors...)
}

func findTransactionsConfirmedInBlockRange(db *gorm.DB, highBlockNumber, lowBlockNumber int64, chainID big.Int) ([]EthTx, error) {
	var etxs []EthTx
	err := db.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC")
		}).
		Preload("EthTxAttempts.EthReceipts").
		Joins("INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_tx_attempts.state = 'broadcast'").
		Joins("INNER JOIN eth_receipts ON eth_receipts.tx_hash = eth_tx_attempts.hash").
		Order("nonce ASC").
		Where("eth_txes.state IN ('confirmed', 'confirmed_missing_receipt') AND block_number BETWEEN ? AND ? AND evm_chain_id = ?", lowBlockNumber, highBlockNumber, chainID.String()).
		Find(&etxs).Error
	return etxs, errors.Wrap(err, "findTransactionsConfirmedInBlockRange failed")
}

func hasReceiptInLongestChain(etx EthTx, head eth.Head) bool {
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
		head = *head.Parent
	}
}

func (ec *EthConfirmer) markForRebroadcast(etx EthTx, head eth.Head) error {
	if len(etx.EthTxAttempts) == 0 {
		return errors.Errorf("invariant violation: expected eth_tx %v to have at least one attempt", etx.ID)
	}

	// Rebroadcast the one with the highest gas price
	attempt := etx.EthTxAttempts[0]
	var receipt EthReceipt
	if len(attempt.EthReceipts) > 0 {
		receipt = attempt.EthReceipts[0]
	}

	ec.logger.Infow(fmt.Sprintf("EthConfirmer: re-org detected. Rebroadcasting transaction %s which may have been re-org'd out of the main chain", attempt.Hash.Hex()),
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
	err := postgres.GormTransactionWithDefaultContext(ec.db, func(tx *gorm.DB) error {
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

func deleteAllReceipts(db *gorm.DB, etxID int64) error {
	return db.Exec(`
		DELETE FROM eth_receipts
		USING eth_tx_attempts
		WHERE eth_receipts.tx_hash = eth_tx_attempts.hash
		AND eth_tx_attempts.eth_tx_id = ?
	`, etxID).Error
}

func unconfirmEthTx(db *gorm.DB, etx EthTx) error {
	if etx.State != EthTxConfirmed {
		return errors.New("expected eth_tx state to be confirmed")
	}
	return errors.Wrap(db.Exec(`UPDATE eth_txes SET state = 'unconfirmed' WHERE id = ?`, etx.ID).Error, "unconfirmEthTx failed")
}

func unbroadcastAttempt(db *gorm.DB, attempt EthTxAttempt) error {
	if attempt.State != EthTxAttemptBroadcast {
		return errors.New("expected eth_tx_attempt to be broadcast")
	}
	return errors.Wrap(db.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = NULL, state = 'in_progress' WHERE id = ?`, attempt.ID).Error, "unbroadcastAttempt failed")
}

// ForceRebroadcast sends a transaction for every nonce in the given nonce range at the given gas price.
// If an eth_tx exists for this nonce, we re-send the existing eth_tx with the supplied parameters.
// If an eth_tx doesn't exist for this nonce, we send a zero transaction.
// This operates completely orthogonal to the normal EthConfirmer and can result in untracked attempts!
// Only for emergency usage.
// This is in case of some unforeseen scenario where the node is refusing to release the lock. KISS.
func (ec *EthConfirmer) ForceRebroadcast(beginningNonce uint, endingNonce uint, gasPriceWei uint64, address gethCommon.Address, overrideGasLimit uint64) error {
	ec.logger.Infof("ForceRebroadcast: will rebroadcast transactions for all nonces between %v and %v", beginningNonce, endingNonce)

	for n := beginningNonce; n <= endingNonce; n++ {
		etx, err := findEthTxWithNonce(ec.db, address, n)
		if err != nil {
			return errors.Wrap(err, "ForceRebroadcast failed")
		}
		if etx == nil {
			ec.logger.Debugf("ForceRebroadcast: no eth_tx found with nonce %v, will rebroadcast empty transaction", n)
			hash, err := ec.sendEmptyTransaction(context.TODO(), address, n, overrideGasLimit, gasPriceWei)
			if err != nil {
				ec.logger.Errorw("ForceRebroadcast: failed to send empty transaction", "nonce", n, "err", err)
				continue
			}
			ec.logger.Infow("ForceRebroadcast: successfully rebroadcast empty transaction", "nonce", n, "hash", hash.String())
		} else {
			ec.logger.Debugf("ForceRebroadcast: got eth_tx %v with nonce %v, will rebroadcast this transaction", etx.ID, *etx.Nonce)
			if overrideGasLimit != 0 {
				etx.GasLimit = overrideGasLimit
			}
			attempt, err := ec.NewLegacyAttempt(*etx, big.NewInt(int64(gasPriceWei)), etx.GasLimit)
			if err != nil {
				ec.logger.Errorw("ForceRebroadcast: failed to create new attempt", "ethTxID", etx.ID, "err", err)
				continue
			}
			if err := sendTransaction(context.TODO(), ec.ethClient, attempt, *etx, ec.logger); err != nil {
				ec.logger.Errorw(fmt.Sprintf("ForceRebroadcast: failed to rebroadcast eth_tx %v with nonce %v at gas price %s wei and gas limit %v: %s", etx.ID, *etx.Nonce, attempt.GasPrice.String(), etx.GasLimit, err.Error()), "err", err)
				continue
			}
			ec.logger.Infof("ForceRebroadcast: successfully rebroadcast eth_tx %v with hash: 0x%x", etx.ID, attempt.Hash)
		}
	}
	return nil
}

func (ec *EthConfirmer) sendEmptyTransaction(ctx context.Context, fromAddress gethCommon.Address, nonce uint, overrideGasLimit uint64, gasPriceWei uint64) (gethCommon.Hash, error) {
	gasLimit := overrideGasLimit
	if gasLimit == 0 {
		gasLimit = ec.config.EvmGasLimitDefault()
	}
	tx, err := sendEmptyTransaction(ec.ethClient, ec.keystore, uint64(nonce), gasLimit, big.NewInt(int64(gasPriceWei)), fromAddress, &ec.chainID)
	if err != nil {
		return gethCommon.Hash{}, errors.Wrap(err, "(EthConfirmer).sendEmptyTransaction failed")
	}
	hash, err := signedTxHash(tx, ec.config.ChainType())
	if err != nil {
		return hash, err
	}
	return hash, nil
}

// findEthTxWithNonce returns any broadcast ethtx with the given nonce
func findEthTxWithNonce(db *gorm.DB, fromAddress gethCommon.Address, nonce uint) (*EthTx, error) {
	etx := EthTx{}
	err := db.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC")
		}).
		First(&etx, "from_address = ? AND nonce = ? AND state IN ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')", fromAddress, nonce).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &etx, errors.Wrap(err, "findEthTxsWithNonce failed")
}

func (ec *EthConfirmer) ResumePendingTaskRuns(ctx context.Context, head eth.Head) error {
	sqlxDB := postgres.UnwrapGormDB(ec.db)

	type x struct {
		ID      uuid.UUID
		Receipt []byte
	}
	var receipts []x
	// NOTE: we don't filter on eth_txes.state = 'confirmed', because a transaction with an attached receipt
	// is guaranteed to be confirmed. This results in a slightly better query plan.
	if err := sqlxDB.Select(&receipts, `
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
