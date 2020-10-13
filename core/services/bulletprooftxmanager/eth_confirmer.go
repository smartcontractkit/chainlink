package bulletprooftxmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

var (
	// ErrCouldNotGetReceipt is the error string we save if we reach our finality depth for a confirmed transaction without ever getting a receipt
	// This most likely happened because an external wallet used the account for this nonce
	ErrCouldNotGetReceipt = "could not get receipt"

	ethConfirmerAdvisoryLockClassID  = int32(1)
	ethConfirmerAdvisoryLockObjectID = int32(0)
)

// EthConfirmer is a broad service which performs four different tasks in sequence on every new longest chain
// Step 1: Mark that all currently pending transaction attempts were broadcast before this block
// Step 2: Check pending transactions for receipts
// Step 3: See if any transactions have exceeded the gas bumping block threshold and, if so, bump them
// Step 4: Check confirmed transactions to make sure they are still in the longest chain (reorg protection)
type EthConfirmer interface {
	store.HeadTrackable
}

type ethConfirmer struct {
	store     *store.Store
	ethClient eth.Client
	config    orm.ConfigReader
}

func NewEthConfirmer(store *store.Store, config orm.ConfigReader) *ethConfirmer {
	return &ethConfirmer{
		store:     store,
		ethClient: store.EthClient,
		config:    config,
	}
}

// Do nothing on connect, simply wait for the next head
func (ec *ethConfirmer) Connect(*models.Head) error {
	return nil
}

func (ec *ethConfirmer) Disconnect() {
	// pass
}

func (ec *ethConfirmer) OnNewLongestChain(ctx context.Context, head models.Head) {
	if ec.config.EnableBulletproofTxManager() {
		if err := ec.ProcessHead(ctx, head); err != nil {
			logger.Errorw("EthConfirmer error", "err", err)
		}
	}
}

// ProcessHead takes all required transactions for the confirmer on a new head
func (ec *ethConfirmer) ProcessHead(ctx context.Context, head models.Head) error {
	return withAdvisoryLock(ec.store, ethConfirmerAdvisoryLockClassID, ethConfirmerAdvisoryLockObjectID, func() error {
		return ec.processHead(ctx, head)
	})
}

// NOTE: This SHOULD NOT be run concurrently or it could behave badly
func (ec *ethConfirmer) processHead(ctx context.Context, head models.Head) error {
	if err := ec.SetBroadcastBeforeBlockNum(head.Number); err != nil {
		return errors.Wrap(err, "SetBroadcastBeforeBlockNum failed")
	}

	mark := time.Now()

	if err := ec.CheckForReceipts(ctx, head.Number); err != nil {
		return errors.Wrap(err, "CheckForReceipts failed")
	}

	logger.Debugw("EthConfirmer: finished CheckForReceipts", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	keys, err := ec.store.SendKeys()
	if err != nil {
		return errors.Wrap(err, "could not fetch keys")
	}
	if err := ec.BumpGasWhereNecessary(ctx, keys, head.Number); err != nil {
		return errors.Wrap(err, "BumpGasWhereNecessary failed")
	}

	logger.Debugw("EthConfirmer: finished BumpGasWhereNecessary", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	defer func() {
		logger.Debugw("EthConfirmer: finished EnsureConfirmedTransactionsInLongestChain", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	}()

	return errors.Wrap(ec.EnsureConfirmedTransactionsInLongestChain(ctx, keys, head), "EnsureConfirmedTransactionsInLongestChain failed")
}

func (ec *ethConfirmer) SetBroadcastBeforeBlockNum(blockNum int64) error {
	return ec.store.DB.Exec(
		`UPDATE eth_tx_attempts SET broadcast_before_block_num = ? WHERE broadcast_before_block_num IS NULL AND state = 'broadcast'`,
		blockNum,
	).Error
}

// receiptFetcherWorkerCount is the max number of concurrently executing
// workers that will fetch receipts for eth transactions
const receiptFetcherWorkerCount = 10

func (ec *ethConfirmer) CheckForReceipts(ctx context.Context, blockNum int64) error {
	etxs, err := ec.findEthTxsRequiringReceiptFetch()
	if err != nil {
		return errors.Wrap(err, "findEthTxsRequiringReceiptFetch failed")
	}
	if len(etxs) == 0 {
		return nil
	}

	logger.Debugf("EthConfirmer: fetching receipt for %v transactions", len(etxs))

	ec.concurrentlyFetchReceipts(ctx, etxs)

	if err := ec.markConfirmedMissingReceipt(ctx); err != nil {
		return errors.Wrap(err, "unable to mark eth_txes as 'confirmed_missing_receipt'")
	}

	if err := ec.markOldTxesMissingReceiptAsErrored(ctx, blockNum); err != nil {
		return errors.Wrap(err, "unable to confirm buried unconfirmed eth_txes")
	}

	return nil
}

func (ec *ethConfirmer) findEthTxsRequiringReceiptFetch() (etxs []models.EthTx, err error) {
	err = ec.store.DB.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC")
		}).
		Order("nonce ASC").
		Find(&etxs, "state IN ('unconfirmed', 'confirmed_missing_receipt')").Error

	return
}

func (ec *ethConfirmer) concurrentlyFetchReceipts(ctx context.Context, etxs []models.EthTx) {
	var wg sync.WaitGroup
	wg.Add(receiptFetcherWorkerCount)
	chEthTxes := make(chan models.EthTx)
	for i := 0; i < receiptFetcherWorkerCount; i++ {
		go ec.fetchReceipts(ctx, chEthTxes, &wg)
	}
	for _, etx := range etxs {
		chEthTxes <- etx
	}
	close(chEthTxes)
	wg.Wait()

}

func (ec *ethConfirmer) fetchReceipts(ctx context.Context, chEthTxes <-chan models.EthTx, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		etx, ok := <-chEthTxes
		if !ok {
			return
		}
		for _, attempt := range etx.EthTxAttempts {
			// NOTE: This could conceivably be optimised even further at the
			// expense of slightly higher load for the remote eth node, by
			// batch requesting all receipts at once
			receipt, err := ec.fetchReceipt(ctx, attempt.Hash)
			if eth.IsParityQueriedReceiptTooEarly(err) || (receipt != nil && receipt.BlockNumber == nil) {
				logger.Debugw("EthConfirmer#fetchReceipts: got receipt for transaction but it's still in the mempool and not included in a block yet", "txHash", attempt.Hash.Hex())
				break
			} else if err != nil {
				logger.Errorw("EthConfirmer#fetchReceipts: fetchReceipt failed", "txHash", attempt.Hash.Hex(), "err", err)
				break
			}
			if receipt != nil {
				logger.Debugw("EthConfirmer#fetchReceipts: got receipt for transaction", "txHash", attempt.Hash.Hex(), "blockNumber", receipt.BlockNumber)
				if receipt.TxHash != attempt.Hash {
					logger.Errorf("EthConfirmer#fetchReceipts: invariant violation, expected receipt with hash %s to have same hash as attempt with hash %s", receipt.TxHash.Hex(), attempt.Hash.Hex())
					break
				}
				if err := ec.saveReceipt(*receipt, etx.ID); err != nil {
					logger.Errorw("EthConfirmer#fetchReceipts: saveReceipt failed", "err", err)
					break
				}
				break
			} else {
				logger.Debugw("EthConfirmer#fetchReceipts: still waiting for receipt", "txHash", attempt.Hash.Hex(), "ethTxAttemptID", attempt.ID, "ethTxID", etx.ID)
			}
		}
	}
}

func (ec *ethConfirmer) fetchReceipt(ctx context.Context, hash gethCommon.Hash) (*gethTypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, maxEthNodeRequestTime)
	defer cancel()
	receipt, err := ec.ethClient.TransactionReceipt(ctx, hash)
	if err != nil && err.Error() == "not found" {
		return nil, nil
	}
	return receipt, err
}

func (ec *ethConfirmer) saveReceipt(receipt gethTypes.Receipt, ethTxID int64) error {
	if receipt.BlockNumber == nil {
		return errors.Errorf("receipt was missing block number: %#v", receipt)
	}

	return ec.store.Transaction(func(tx *gorm.DB) error {
		receiptJSON, err := json.Marshal(receipt)
		if err != nil {
			return errors.Wrap(err, "saveReceipt failed")
		}
		// Conflict here shouldn't be possible because there should only ever
		// be one receipt for an eth_tx, and if it exists then the transaction
		// is marked confirmed which means we can never get here.
		// However, even so, it still shouldn't be an error to re-insert a receipt we already have.
		err = tx.Set("gorm:insert_option", "ON CONFLICT (tx_hash, block_hash) DO NOTHING").
			Create(&models.EthReceipt{
				Receipt:          receiptJSON,
				TxHash:           receipt.TxHash,
				BlockHash:        receipt.BlockHash,
				BlockNumber:      receipt.BlockNumber.Int64(),
				TransactionIndex: receipt.TransactionIndex,
			}).Error
		if err == nil || err.Error() == "sql: no rows in result set" {
			return errors.Wrap(tx.Exec(`UPDATE eth_txes SET state = 'confirmed' WHERE id = ?`, ethTxID).Error, "saveReceipt failed to update eth_txes")
		}

		return errors.Wrap(err, "saveReceipt failed to save receipt")
	})
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
// We will continue to try to fetch a receipt for these attempts until all
// attempts are below the finality depth from current head.
func (ec *ethConfirmer) markConfirmedMissingReceipt(ctx context.Context) (err error) {
	_, err = ec.store.DB.DB().ExecContext(ctx, `
UPDATE eth_txes
SET state = 'confirmed_missing_receipt'
WHERE state = 'unconfirmed'
AND nonce < (
	SELECT MAX(nonce) FROM eth_txes
	WHERE state = 'confirmed'
)
	`)
	return
}

// markOldTxesMissingReceiptAsErrored
//
// Once eth_tx has all of its attempts broadcast before some cutoff threshold,
// we mark it as fatally errored (never sent).
//
// The job run will also be marked as errored in this case since we never got a
// receipt and thus cannot pass on any transaction hash
func (ec *ethConfirmer) markOldTxesMissingReceiptAsErrored(ctx context.Context, blockNum int64) error {
	// cutoff is a block height
	// Any 'confirmed_missing_receipt' eth_tx with all attempts older than this block height will be marked as errored
	// We will not try to query for receipts for this transaction any more
	cutoff := blockNum - int64(ec.config.EthFinalityDepth())
	if cutoff <= 0 {
		return nil
	}
	rows, err := ec.store.DB.DB().QueryContext(ctx, `
UPDATE eth_txes
SET state='fatal_error', nonce=NULL, error=$1, broadcast_at=NULL
WHERE id IN (
	SELECT eth_txes.id FROM eth_txes
	INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id
	WHERE eth_txes.state = 'confirmed_missing_receipt'
	GROUP BY eth_txes.id
	HAVING max(eth_tx_attempts.broadcast_before_block_num) < $2
)
RETURNING id, nonce, from_address`, ErrCouldNotGetReceipt, cutoff)

	if err != nil {
		return errors.Wrap(err, "markOldTxesMissingReceiptAsErrored failed to query")
	}

	for rows.Next() {
		var ethTxID int64
		var nonce null.Int64
		var fromAddress gethCommon.Address
		if err = rows.Scan(&ethTxID, &nonce, &fromAddress); err != nil {
			return errors.Wrap(err, "error scanning row")
		}

		logger.Errorf("EthConfirmer: eth_tx with ID %v expired without ever getting a receipt for any of our attempts. "+
			"Current block height is %v. This transaction has not been sent and will be marked as fatally errored. "+
			"This can happen if an external wallet has been used to send a transaction from account %s with nonce %v."+
			" Please note that using the chainlink keys with an external wallet is NOT SUPPORTED and WILL lead to missed transactions",
			ethTxID, blockNum, fromAddress.Hex(), nonce.Int64)
	}

	return errors.Wrap(rows.Close(), "markOldTxesMissingReceiptAsErrored failed to close rows")
}

func (ec *ethConfirmer) BumpGasWhereNecessary(ctx context.Context, keys []models.Key, blockHeight int64) error {
	var wg sync.WaitGroup

	// It is safe to process separate keys concurrently
	// NOTE: This design will block one key if another takes a really long time to execute
	wg.Add(len(keys))
	errors := []error{}
	var errMu sync.Mutex
	for _, key := range keys {
		go func(fromAddress gethCommon.Address) {
			if err := ec.bumpGasWhereNecessary(ctx, fromAddress, blockHeight); err != nil {
				errMu.Lock()
				errors = append(errors, err)
				errMu.Unlock()
				logger.Errorw("Error in BumpGasWhereNecessary", "error", err, "fromAddress", fromAddress)
			}

			wg.Done()
		}(key.Address.Address())
	}

	wg.Wait()

	return multierr.Combine(errors...)
}

func (ec *ethConfirmer) bumpGasWhereNecessary(ctx context.Context, address gethCommon.Address, blockHeight int64) error {
	if err := ec.handleAnyInProgressAttempts(ctx, address, blockHeight); err != nil {
		return errors.Wrap(err, "handleAnyInProgressAttempts failed")
	}

	threshold := int64(ec.config.EthGasBumpThreshold())
	depth := int64(ec.config.EthGasBumpTxDepth())
	etxs, err := FindEthTxsRequiringNewAttempt(ec.store.DB, address, blockHeight, threshold, depth)
	if err != nil {
		return errors.Wrap(err, "FindEthTxsRequiringNewAttempt failed")
	}
	if len(etxs) > 0 {
		logger.Debugf("EthConfirmer: Bumping gas for %v transactions", len(etxs))
	}
	for _, etx := range etxs {
		attempt, err := ec.newAttemptWithGasBump(etx)
		if err != nil {
			return errors.Wrap(err, "newAttemptWithGasBump failed")
		}

		if err := ec.saveInProgressAttempt(&attempt); err != nil {
			return errors.Wrap(err, "saveInProgressAttempt failed")
		}

		if err := ec.handleInProgressAttempt(ctx, etx, attempt, blockHeight); err != nil {
			return errors.Wrap(err, "handleInProgressAttempt failed")
		}
	}
	return nil
}

// "in_progress" attempts were left behind after a crash/restart and may or may not have been sent
// We should try to ensure they get on-chain so we can fetch a receipt for them
func (ec *ethConfirmer) handleAnyInProgressAttempts(ctx context.Context, address gethCommon.Address, blockHeight int64) error {
	attempts, err := getInProgressEthTxAttempts(ec.store, address)
	if err != nil {
		return errors.Wrap(err, "getInProgressEthTxAttempts failed")
	}
	for _, a := range attempts {
		if err := ec.handleInProgressAttempt(ctx, a.EthTx, a, blockHeight); err != nil {
			return errors.Wrap(err, "handleInProgressAttempt failed")
		}
	}
	return nil
}

func getInProgressEthTxAttempts(s *store.Store, address gethCommon.Address) ([]models.EthTxAttempt, error) {
	var attempts []models.EthTxAttempt
	err := s.DB.
		Preload("EthTx").
		Joins("INNER JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state in ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')").
		Where("eth_tx_attempts.state = 'in_progress'").
		Where("eth_txes.from_address = ?", address).
		Find(&attempts).Error
	return attempts, errors.Wrap(err, "getInProgressEthTxAttempts failed")
}

// FindEthTxsRequiringNewAttempt returns transactions that have all
// attempts which are unconfirmed for at least gasBumpThreshold blocks,
// limited by limit pending transactions
func FindEthTxsRequiringNewAttempt(db *gorm.DB, address gethCommon.Address, blockNum, gasBumpThreshold, depth int64) (etxs []models.EthTx, err error) {
	q := db.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC")
		}).
		Joins("LEFT JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id "+
			"AND eth_tx_attempts.state != 'insufficient_eth' "+
			"AND (broadcast_before_block_num > ? OR broadcast_before_block_num IS NULL OR eth_tx_attempts.state != 'broadcast')", blockNum-gasBumpThreshold).
		Where("eth_txes.state = 'unconfirmed' AND eth_tx_attempts.id IS NULL")

	if depth > 0 {
		q = q.Where("eth_txes.id IN (SELECT id FROM eth_txes WHERE state = 'unconfirmed' AND from_address = ? ORDER BY nonce ASC LIMIT ?)", address, depth)
	}

	err = q.Order("nonce ASC").Find(&etxs).Error

	err = errors.Wrap(err, "FindEthTxsRequiringNewAttempt failed")

	return
}

func (ec *ethConfirmer) newAttemptWithGasBump(etx models.EthTx) (attempt models.EthTxAttempt, err error) {
	var bumpedGasPrice *big.Int
	if len(etx.EthTxAttempts) > 0 {
		previousAttempt := etx.EthTxAttempts[0]
		if previousAttempt.State == models.EthTxAttemptInsufficientEth {
			// Do not create a new attempt if we ran out of eth last time since bumping gas is pointless
			// Instead try to resubmit the same attempt at the same price, in the hope that the wallet was funded since our last attempt
			previousAttempt.State = models.EthTxAttemptInProgress
			return previousAttempt, nil
		}
		previousGasPrice := previousAttempt.GasPrice
		bumpedGasPrice, err = BumpGas(ec.config, previousGasPrice.ToInt())
		if err != nil {
			logger.Errorw("Failed to bump gas", "err", err, "etxID", etx.ID, "txHash", attempt.Hash, "originalGasPrice", previousGasPrice.String(), "maxGasPrice", ec.config.EthMaxGasPriceWei())
			// Do not create a new attempt if bumping gas would put us over the limit or cause some other problem
			// Instead try to resubmit the previous attempt, and keep resubmitting until its accepted
			previousAttempt.BroadcastBeforeBlockNum = nil
			previousAttempt.State = models.EthTxAttemptInProgress
			return previousAttempt, nil
		}
	} else {
		logger.Errorf("invariant violation: EthTx %v was unconfirmed but didn't have any attempts. "+
			"Falling back to default gas price instead."+
			"This is a bug! Please report to https://github.com/smartcontractkit/chainlink/issues", etx.ID)
		bumpedGasPrice = ec.config.EthGasPriceDefault()
	}
	return newAttempt(ec.store, etx, bumpedGasPrice)
}

func (ec *ethConfirmer) saveInProgressAttempt(attempt *models.EthTxAttempt) error {
	if attempt.State != models.EthTxAttemptInProgress {
		return errors.New("saveInProgressAttempt failed: attempt state must be in_progress")
	}
	return errors.Wrap(ec.store.DB.Save(attempt).Error, "saveInProgressAttempt failed")
}

func (ec *ethConfirmer) handleInProgressAttempt(ctx context.Context, etx models.EthTx, attempt models.EthTxAttempt, blockHeight int64) error {
	if attempt.State != models.EthTxAttemptInProgress {
		return errors.Errorf("invariant violation: expected eth_tx_attempt %v to be in_progress, it was %s", attempt.ID, attempt.State)
	}

	sendError := sendTransaction(ctx, ec.ethClient, attempt)

	if sendError.IsTerminallyUnderpriced() {
		// This should really not ever happen in normal operation since we
		// already bumped above the required minimum in ethBroadcaster.
		//
		// It could conceivably happen if the remote eth node changed it's configuration.
		bumpedGasPrice, err := BumpGas(ec.config, attempt.GasPrice.ToInt())
		if err != nil {
			return errors.Wrap(err, "could not bump gas for terminally underpriced transaction")
		}
		logger.Errorf("gas price %v wei was rejected by the eth node for being too low. "+
			"Eth node returned: '%s'. "+
			"Bumping to %v wei and retrying. "+
			"ACTION REQUIRED: You should consider increasing ETH_GAS_PRICE_DEFAULT", attempt.GasPrice, sendError.Error(), bumpedGasPrice)
		replacementAttempt, err := newAttempt(ec.store, etx, bumpedGasPrice)
		if err != nil {
			return errors.Wrap(err, "newAttempt failed")
		}

		if err := saveReplacementInProgressAttempt(ec.store, attempt, &replacementAttempt); err != nil {
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
		logger.Infow("EthConfirmer: Transaction temporarily underpriced", "ethTxID", etx.ID, "attemptID", attempt.ID, "err", sendError.Error(), "gasPriceWei", attempt.GasPrice.String())
		sendError = nil
	}

	if sendError.Fatal() {
		// WARNING: This should never happen!
		// Should NEVER be fatal this is an invariant violation. The
		// EthBroadcaster can never create an EthTxAttempt that will
		// fatally error.
		//
		// The only scenario imaginable where this might take place is if
		// geth/parity have been updated between broadcasting and confirming steps.
		logger.Errorw("invariant violation: fatal error while re-attempting transaction",
			"ethTxID", etx.ID,
			"err", sendError,
			"signedRawTx", hexutil.Encode(attempt.SignedRawTx),
			"blockHeight", blockHeight,
		)
		// This will loop continuously on every new head so it must be handled manually by the node operator!
		return deleteInProgressAttempt(ec.store.DB, attempt)
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
		logger.Errorw(fmt.Sprintf("EthConfirmer: replacement transaction underpriced at %v wei for eth_tx %v. "+
			"Eth node returned error: '%s'. "+
			"Either you have set ETH_GAS_BUMP_PERCENT (currently %v%%) too low or an external wallet used this account. "+
			"Please note that using your node's private keys outside of the chainlink node is NOT SUPPORTED and can lead to missed transactions.",
			attempt.GasPrice.ToInt().Int64(), etx.ID, sendError.Error(), ec.store.Config.EthGasBumpPercent()), "err", sendError)

		// Assume success and hand off to the next cycle.
		sendError = nil
	}

	if sendError.IsInsufficientEth() {
		logger.Errorw(fmt.Sprintf("EthConfirmer: EthTxAttempt %v (hash 0x%x) at gas price (%s Wei) was rejected due to insufficient eth. "+
			"The eth node returned %s. "+
			"ACTION REQUIRED: Chainlink wallet with address 0x%x is OUT OF FUNDS",
			attempt.ID, attempt.Hash, attempt.GasPrice.String(), sendError.Error(), etx.FromAddress,
		), "err", sendError)
		return saveInsufficientEthAttempt(ec.store.DB, &attempt)
	}

	if sendError == nil {
		return saveSentAttempt(ec.store.DB, &attempt)
	}

	// Any other type of error is considered temporary or resolvable by the
	// node operator. The node may have it in the mempool so we must keep the
	// attempt (leave it in_progress). Safest thing to do is bail out and wait
	// for the next head.
	return errors.Wrapf(sendError, "unexpected error sending eth_tx %v with hash %s", etx.ID, attempt.Hash.Hex())
}

func deleteInProgressAttempt(db *gorm.DB, attempt models.EthTxAttempt) error {
	if attempt.State != models.EthTxAttemptInProgress {
		return errors.New("deleteInProgressAttempt: expected attempt state to be in_progress")
	}
	if attempt.ID == 0 {
		return errors.New("deleteInProgressAttempt: expected attempt to have an id")
	}
	return errors.Wrap(db.Exec(`DELETE FROM eth_tx_attempts WHERE id = ?`, attempt.ID).Error, "deleteInProgressAttempt failed")
}

func saveSentAttempt(db *gorm.DB, attempt *models.EthTxAttempt) error {
	if attempt.State != models.EthTxAttemptInProgress {
		return errors.New("expected state to be in_progress")
	}
	attempt.State = models.EthTxAttemptBroadcast
	return errors.Wrap(db.Save(attempt).Error, "saveSentAttempt failed")
}

func saveInsufficientEthAttempt(db *gorm.DB, attempt *models.EthTxAttempt) error {
	if !(attempt.State == models.EthTxAttemptInProgress || attempt.State == models.EthTxAttemptInsufficientEth) {
		return errors.New("expected state to be either in_progress or insufficient_eth")
	}
	attempt.State = models.EthTxAttemptInsufficientEth
	return errors.Wrap(db.Save(attempt).Error, "saveInsufficientEthAttempt failed")

}

// EnsureConfirmedTransactionsInLongestChain finds all confirmed eth_txes up to the depth
// of the given chain and ensures that every one has a receipt with a block hash that is
// in the given chain.
//
// If any of the confirmed transactions does not have a receipt in the chain, it has been
// re-org'd out and will be rebroadcast.
func (ec *ethConfirmer) EnsureConfirmedTransactionsInLongestChain(ctx context.Context, keys []models.Key, head models.Head) error {
	etxs, err := findTransactionsConfirmedAtOrAboveBlockHeight(ec.store.DB, head.EarliestInChain().Number)
	if err != nil {
		return errors.Wrap(err, "findTransactionsConfirmedAtOrAboveBlockHeight failed")
	}

	for _, etx := range etxs {
		if !hasReceiptInLongestChain(etx, head) {
			if err := ec.markForRebroadcast(etx); err != nil {
				return errors.Wrapf(err, "markForRebroadcast failed for etx %v", etx.ID)
			}
		}
	}

	// It is safe to process separate keys concurrently
	// NOTE: This design will block one key if another takes a really long time to execute
	var wg sync.WaitGroup
	errors := []error{}
	var errMu sync.Mutex
	wg.Add(len(keys))
	for _, key := range keys {
		go func(fromAddress gethCommon.Address) {
			if err := ec.handleAnyInProgressAttempts(ctx, fromAddress, head.Number); err != nil {
				errMu.Lock()
				errors = append(errors, err)
				errMu.Unlock()
				logger.Errorw("Error in BumpGasWhereNecessary", "error", err, "fromAddress", fromAddress)
			}

			wg.Done()
		}(key.Address.Address())
	}

	wg.Wait()

	return multierr.Combine(errors...)
}

func findTransactionsConfirmedAtOrAboveBlockHeight(db *gorm.DB, blockNumber int64) ([]models.EthTx, error) {
	var etxs []models.EthTx
	err := db.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC")
		}).
		Preload("EthTxAttempts.EthReceipts").
		Joins("INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_tx_attempts.state = 'broadcast'").
		Joins("INNER JOIN eth_receipts ON eth_receipts.tx_hash = eth_tx_attempts.hash").
		Order("nonce ASC").
		Where("eth_txes.state IN ('confirmed', 'confirmed_missing_receipt') AND block_number >= ?", blockNumber).
		Find(&etxs).Error
	return etxs, errors.Wrap(err, "findTransactionsConfirmedAtOrAboveBlockHeight failed")
}

func hasReceiptInLongestChain(etx models.EthTx, head models.Head) bool {
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

func (ec *ethConfirmer) markForRebroadcast(etx models.EthTx) error {
	if len(etx.EthTxAttempts) == 0 {
		return errors.Errorf("invariant violation: expected eth_tx %v to have at least one attempt", etx.ID)
	}

	// Rebroadcast the one with the highest gas price
	attempt := etx.EthTxAttempts[0]

	// Put it back in progress and delete all receipts (they do not apply to the new chain)
	err := ec.store.Transaction(func(tx *gorm.DB) error {
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

func unconfirmEthTx(db *gorm.DB, etx models.EthTx) error {
	if etx.State != models.EthTxConfirmed {
		return errors.New("expected eth_tx state to be confirmed")
	}
	return errors.Wrap(db.Exec(`UPDATE eth_txes SET state = 'unconfirmed' WHERE id = ?`, etx.ID).Error, "unconfirmEthTx failed")
}

func unbroadcastAttempt(db *gorm.DB, attempt models.EthTxAttempt) error {
	if attempt.State != models.EthTxAttemptBroadcast {
		return errors.New("expected eth_tx_attempt to be broadcast")
	}
	return errors.Wrap(db.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = NULL, state = 'in_progress' WHERE id = ?`, attempt.ID).Error, "unbroadcastAttempt failed")
}

// ForceRebroadcast sends a transaction for every nonce in the given nonce range at the given gas price.
// If an eth_tx exists for this nonce, we re-send the existing eth_tx with the supplied parameters.
// If an eth_tx doesn't exist for this nonce, we send a zero transaction.
// This operates completely orthogonal to the normal EthConfirmer and can result in untracked attempts!
// Only for emergency usage.
// Deliberately does not take the advisory lock (we don't write to the database so this is safe from a data integrity perspective).
// This is in case of some unforeseen scenario where the node is refusing to release the lock. KISS.
func (ec *ethConfirmer) ForceRebroadcast(beginningNonce uint, endingNonce uint, gasPriceWei uint64, address gethCommon.Address, overrideGasLimit uint64) error {
	logger.Infof("ForceRebroadcast: will rebroadcast transactions for all nonces between %v and %v", beginningNonce, endingNonce)

	for n := beginningNonce; n <= endingNonce; n++ {
		etx, err := findEthTxWithNonce(ec.store.DB, address, n)
		if err != nil {
			return errors.Wrap(err, "ForceRebroadcast failed")
		}
		if etx == nil {
			logger.Debugf("ForceRebroadcast: no eth_tx found with nonce %v, will rebroadcast empty transaction", n)
			hash, err := ec.sendEmptyTransaction(context.TODO(), address, n, overrideGasLimit, gasPriceWei)
			if err != nil {
				logger.Errorw("ForceRebroadcast: failed to send empty transaction", "nonce", n, "err", err)
				continue
			}
			logger.Infow("ForceRebroadcast: successfully rebroadcast empty transaction", "nonce", n, "hash", hash.String())
		} else {
			logger.Debugf("ForceRebroadcast: got eth_tx %v with nonce %v, will rebroadcast this transaction", etx.ID, *etx.Nonce)
			if overrideGasLimit != 0 {
				etx.GasLimit = overrideGasLimit
			}
			attempt, err := newAttempt(ec.store, *etx, big.NewInt(int64(gasPriceWei)))
			if err != nil {
				logger.Errorw("ForceRebroadcast: failed to create new attempt", "ethTxID", etx.ID, "err", err)
				continue
			}
			if err := sendTransaction(context.TODO(), ec.ethClient, attempt); err != nil {
				logger.Errorw(fmt.Sprintf("ForceRebroadcast: failed to rebroadcast eth_tx %v with nonce %v at gas price %s wei and gas limit %v: %s", etx.ID, *etx.Nonce, attempt.GasPrice.String(), etx.GasLimit, err.Error()), "err", err)
				continue
			}
			logger.Infof("ForceRebroadcast: successfully rebroadcast eth_tx %v with hash: 0x%x", etx.ID, attempt.Hash)
		}
	}
	return nil
}

func (ec *ethConfirmer) sendEmptyTransaction(ctx context.Context, fromAddress gethCommon.Address, nonce uint, overrideGasLimit uint64, gasPriceWei uint64) (gethCommon.Hash, error) {
	gasLimit := overrideGasLimit
	if gasLimit == 0 {
		gasLimit = ec.config.EthGasLimitDefault()
	}
	account, err := ec.store.KeyStore.GetAccountByAddress(fromAddress)
	if err != nil {
		return gethCommon.Hash{}, errors.Wrap(err, "(ethConfirmer).sendEmptyTransaction failed")
	}
	tx, err := sendEmptyTransaction(ec.ethClient, ec.store.KeyStore, uint64(nonce), gasLimit, big.NewInt(int64(gasPriceWei)), account, ec.config.ChainID())
	if err != nil {
		return gethCommon.Hash{}, errors.Wrap(err, "(ethConfirmer).sendEmptyTransaction failed")
	}
	return tx.Hash(), nil
}

// findEthTxWithNonce returns any broadcast ethtx with the given nonce
func findEthTxWithNonce(db *gorm.DB, fromAddress gethCommon.Address, nonce uint) (*models.EthTx, error) {
	etx := models.EthTx{}
	err := db.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC")
		}).
		First(&etx, "from_address = ? AND nonce = ? AND state IN ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')", fromAddress, nonce).
		Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &etx, errors.Wrap(err, "findEthTxsWithNonce failed")
}
