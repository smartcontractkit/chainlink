package bulletprooftxmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	// ErrExternalWalletUsedNonce is the error string we save if we come to the conclusion that the transaction nonce was used by an external account
	ErrExternalWalletUsedNonce = "external wallet used nonce"

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
	store             *store.Store
	gethClientWrapper store.GethClientWrapper
	config            orm.ConfigReader
}

func NewEthConfirmer(store *store.Store, config orm.ConfigReader) *ethConfirmer {
	return &ethConfirmer{
		store:             store,
		gethClientWrapper: store.GethClientWrapper,
		config:            config,
	}
}

// Do nothing on connect, simply wait for the next head
func (ec *ethConfirmer) Connect(*models.Head) error {
	return nil
}

func (ec *ethConfirmer) Disconnect() {
	// pass
}

func (ec *ethConfirmer) OnNewLongestChain(head models.Head) {
	if ec.config.EnableBulletproofTxManager() {
		if err := ec.ProcessHead(head); err != nil {
			logger.Error("EthConfirmer: ", err)
		}
	}
}

// ProcessHead takes all required transactions for the confirmer on a new head
func (ec *ethConfirmer) ProcessHead(head models.Head) error {
	return withAdvisoryLock(ec.store, ethConfirmerAdvisoryLockClassID, ethConfirmerAdvisoryLockObjectID, func() error {
		return ec.processHead(head)
	})
}

// NOTE: This SHOULD NOT be run concurrently or it could behave badly
func (ec *ethConfirmer) processHead(head models.Head) error {
	logger.Debugw("EthConfirmer: running SetBroadcastBeforeBlockNum", "headNum", head.Number)
	if err := ec.SetBroadcastBeforeBlockNum(head.Number); err != nil {
		return errors.Wrap(err, "SetBroadcastBeforeBlockNum failed")
	}

	logger.Debugw("EthConfirmer: running CheckForReceipts", "headNum", head.Number)
	if err := ec.CheckForReceipts(); err != nil {
		return errors.Wrap(err, "CheckForReceipts failed")
	}

	logger.Debugw("EthConfirmer: running BumpGasWhereNecessary", "headNum", head.Number)
	if err := ec.BumpGasWhereNecessary(head.Number); err != nil {
		return errors.Wrap(err, "BumpGasWhereNecessary failed")
	}

	logger.Debugw("EthConfirmer: running EnsureConfirmedTransactionsInLongestChain", "headNum", head.Number)
	return errors.Wrap(ec.EnsureConfirmedTransactionsInLongestChain(head), "EnsureConfirmedTransactionsInLongestChain failed")
}

func (ec *ethConfirmer) SetBroadcastBeforeBlockNum(blockNum int64) error {
	return ec.store.GetRawDB().Exec(
		`UPDATE eth_tx_attempts SET broadcast_before_block_num = ? WHERE broadcast_before_block_num IS NULL AND state = 'broadcast'`,
		blockNum,
	).Error
}

func (ec *ethConfirmer) CheckForReceipts() error {
	unconfirmedEtxs, err := ec.findUnconfirmedEthTxs()
	if err != nil {
		return errors.Wrap(err, "findUnconfirmedEthTxs failed")
	}
	if len(unconfirmedEtxs) > 0 {
		logger.Debugf("EthConfirmer: %v unconfirmed transactions", len(unconfirmedEtxs))
	}
	for _, etx := range unconfirmedEtxs {
		for _, attempt := range etx.EthTxAttempts {
			// NOTE: If this becomes a performance bottleneck due to eth node requests,
			// it may be possible to use goroutines here to speed it up by
			// issuing `fetchReceipt` requests in parallel
			receipt, err := ec.fetchReceipt(attempt.Hash)
			if isParityQueriedReceiptTooEarly(err) {
				logger.Debugw("EthConfirmer: got receipt for transaction but it's still in the mempool and not included in a block yet", "txHash", attempt.Hash.Hex())
				break
			} else if err != nil {
				return errors.Wrap(err, fmt.Sprintf("fetchReceipt failed for transaction %s", attempt.Hash.Hex()))
			}
			if receipt != nil {
				logger.Debugw("EthConfirmer: got receipt for transaction", "txHash", attempt.Hash.Hex(), "blockNumber", receipt.BlockNumber)
				if receipt.TxHash != attempt.Hash {
					return errors.Errorf("invariant violation: expected receipt with hash %s to have same hash as attempt with hash %s", receipt.TxHash.Hex(), attempt.Hash.Hex())
				}
				if err := ec.saveReceipt(*receipt, etx.ID); err != nil {
					return err
				}
				break
			} else {
				logger.Debugw("EthConfirmer: still waiting for receipt", "txHash", attempt.Hash.Hex())
			}
		}
	}
	return nil
}

func (ec *ethConfirmer) findUnconfirmedEthTxs() ([]models.EthTx, error) {
	var etxs []models.EthTx
	err := ec.store.GetRawDB().
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC")
		}).
		Order("nonce ASC").
		Find(&etxs, "eth_txes.state = 'unconfirmed'").Error
	return etxs, err
}

func (ec *ethConfirmer) fetchReceipt(hash gethCommon.Hash) (*gethTypes.Receipt, error) {
	var receipt *gethTypes.Receipt
	err := ec.gethClientWrapper.GethClient(func(gethClient eth.GethClient) error {
		ctx, cancel := context.WithTimeout(context.Background(), maxEthNodeRequestTime)
		defer cancel()
		var err error
		receipt, err = gethClient.TransactionReceipt(ctx, hash)
		return err
	})
	if err != nil && err.Error() == "not found" {
		return nil, nil
	}
	return receipt, err

}

func (ec *ethConfirmer) saveReceipt(receipt gethTypes.Receipt, ethTxID int64) error {
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

func (ec *ethConfirmer) BumpGasWhereNecessary(blockHeight int64) error {
	if err := ec.handleAnyInProgressAttempts(blockHeight); err != nil {
		return err
	}

	etxs, err := FindEthTxsRequiringNewAttempt(ec.store.GetRawDB(), blockHeight, int64(ec.config.EthGasBumpThreshold()))
	if err != nil {
		return err
	}
	if len(etxs) > 0 {
		logger.Debugf("EthConfirmer: Bumping gas for %v transactions", len(etxs))
	}
	for _, etx := range etxs {
		attempt, err := ec.newAttemptWithGasBump(etx)
		if err != nil {
			return err
		}

		if err := ec.saveInProgressAttempt(&attempt); err != nil {
			return err
		}

		if err := ec.handleInProgressAttempt(etx, attempt, blockHeight, true); err != nil {
			return err
		}
	}
	return nil
}

// "in_progress" attempts were left behind after a crash/restart and may or may not have been sent
// We should try to ensure they get on-chain so we can fetch a receipt for them
func (ec *ethConfirmer) handleAnyInProgressAttempts(blockHeight int64) error {
	attempts, err := getInProgressEthTxAttempts(ec.store)
	if err != nil {
		return err
	}
	for _, a := range attempts {
		if err := ec.handleInProgressAttempt(a.EthTx, a, blockHeight, false); err != nil {
			return err
		}
	}
	return nil
}

func getInProgressEthTxAttempts(s *store.Store) ([]models.EthTxAttempt, error) {
	var attempts []models.EthTxAttempt
	err := s.GetRawDB().
		Preload("EthTx").
		Joins("INNER JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state in ('confirmed', 'unconfirmed')").
		Where("eth_tx_attempts.state = 'in_progress'").
		Find(&attempts).Error
	return attempts, errors.Wrap(err, "getInProgressEthTxAttempts failed")
}

// FindEthTxsRequiringNewAttempt returns transactions that have all
// attempts which are unconfirmed for at least gasBumpThreshold blocks
func FindEthTxsRequiringNewAttempt(db *gorm.DB, blockNum int64, gasBumpThreshold int64) ([]models.EthTx, error) {
	var etxs []models.EthTx
	err := db.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC")
		}).
		Joins("LEFT JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id "+
			"AND (broadcast_before_block_num > ? OR broadcast_before_block_num IS NULL OR eth_tx_attempts.state != 'broadcast')", blockNum-gasBumpThreshold).
		Order("nonce ASC").
		Where("eth_txes.state = 'unconfirmed' AND eth_tx_attempts.id IS NULL").
		Find(&etxs).Error

	return etxs, errors.Wrap(err, "FindEthTxsRequiringNewAttempt failed")
}

func (ec *ethConfirmer) newAttemptWithGasBump(etx models.EthTx) (models.EthTxAttempt, error) {
	var bumpedGasPrice *big.Int
	if len(etx.EthTxAttempts) > 0 {
		previousGasPrice := etx.EthTxAttempts[0].GasPrice
		bumpedGasPrice = BumpGas(ec.config, previousGasPrice.ToInt())
	} else {
		logger.Error("invariant violation: EthTx %v was unconfirmed but didn't have any attempts. "+
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
	return errors.Wrap(ec.store.GetRawDB().Save(attempt).Error, "saveInProgressAttempt failed")
}

func (ec *ethConfirmer) handleInProgressAttempt(etx models.EthTx, a models.EthTxAttempt, blockHeight int64, isVirginAttempt bool) error {
	if a.State != models.EthTxAttemptInProgress {
		return errors.Errorf("invariant violation: expected eth_tx_attempt %v to be in_progress, it was %s", a.ID, a.State)
	}

	sendError := sendTransaction(ec.gethClientWrapper, a)

	if sendError.IsTerminallyUnderpriced() {
		// This should really not ever happen in normal operation since we
		// already bumped the required amount in ethConfirmer.
		//
		// It may happen if the eth node changed it's configuration, or it is a
		// parity node that is rejecting transactions due to mempool pressure
		bumpedGasPrice := BumpGas(ec.config, a.GasPrice.ToInt())
		logger.Warnf("gas price %v wei was rejected by the eth node for being too low. "+
			"Eth node returned: '%s'. "+
			"Bumping to %v wei and retrying. "+
			"You should consider increasing ETH_GAS_PRICE_DEFAULT", a.GasPrice, sendError.Error(), bumpedGasPrice)
		replacementAttempt, err := newAttempt(ec.store, etx, bumpedGasPrice)
		if err != nil {
			return err
		}

		if err := saveReplacementInProgressAttempt(ec.store, a, &replacementAttempt); err != nil {
			return err
		}
		return ec.handleInProgressAttempt(etx, replacementAttempt, blockHeight, isVirginAttempt)
	}

	if sendError.Fatal() {
		// WARNING: This should never happen!
		// Should NEVER be fatal this is an invariant violation. The
		// EthBroadcaster can never create an EthTxAttempt that will
		// fatally error or be terminally underpriced.
		logger.Errorf("invariant violation: fatal error while reattempting transaction %v: '%s'. "+
			"SignedRawTx: %s\n"+
			"BlockHeight: %v\n"+
			"IsVirginAttempt: %v\n"+
			"ACTION REQUIRED: Your node is BROKEN - this error should never happen in normal operation. "+
			"Please consider raising an issue here: https://github.com/smartcontractkit/chainlink/issues", etx.ID, sendError, hexutil.Encode(a.SignedRawTx), blockHeight, isVirginAttempt)
		// This will loop continously on every new head so it must be handled manually by the node operator!
		return deleteInProgressAttempt(ec.store.GetRawDB(), a)
	}

	if sendError.IsNonceTooLowError() {
		// Nonce too low indicated that it was confirmed already. Success!
		// This attempt is unnecessary, so we can wait for one of the previous
		// attempts to catch a receipt on the next loop

		// NOTE: It is possible we will hit this over and over again forever if
		// another wallet used the chainlink keys. This is because without the
		// attempt in our database, we can never fetch the receipt and thus
		// never mark the transaction as confirmed.
		//
		// The simplest and safest thing to do in that case is to keep
		// requesting receipts on each head until one of two things happens:
		//
		// 1. We get a receipt for one of our transactions
		// 2. We hit such a high block height that whichever transaction is
		// there won't get re-org'd out so it doesn't matter anyway (for
		// gapless nonce purposes).
		//
		// NOTE: It may be possible to introduce some optimisation here since
		// we know that if a later nonce is confirmed, earlier nonces are also
		// automatically confirmed.
		if ec.IsSafeToAbandon(etx, blockHeight) {
			logger.Errorf("nonce %v for transaction %v was already used but we never got a receipt from the eth node for any of our attempts. "+
				"Current block height is %v. This transaction has not been sent and will be marked as fatally errored. "+
				"This can happen if an external wallet has been used to send a transaction from account %s with nonce %v."+
				" Please note that using the chainlink keys with an external wallet is NOT SUPPORTED and can lead to missed transactions",
				*etx.Nonce, etx.ID, blockHeight, etx.FromAddress.Hex(), *etx.Nonce)

			return saveExternalWalletUsedNonce(ec.store, &etx, a)
		}

		if isVirginAttempt {
			// If we get this error, and this attempt has never been sent before
			// it is extremely likely that either:
			//
			// 1. One of our previous attempts was successful
			// 2. An external wallet messed with our nonce
			//
			// Either way, there is no point keeping the current attempt around
			// since it will never confirm.
			//
			// On the extremely minute chance this is due to a network double
			// send or something bizarre, we will fail to get a receipt for one
			// of the other transactions and simply enter this loop again.
			return deleteInProgressAttempt(ec.store.GetRawDB(), a)
		}
		// If we already sent the attempt, we have to assume the one who was
		// confirmed was this one, so simply mark it as broadcast and wait for
		// a receipt.
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
		logger.Errorf("replacement transaction underpriced at %v wei for eth_tx %v. "+
			"Eth node returned error: '%s'. "+
			"Either you have set ETH_GAS_BUMP_PERCENT (currently %v%%) too low or an external wallet used this account. "+
			"Please note that using your node's private keys outside of the chainlink node is NOT SUPPORTED and can lead to missed transactions.",
			a.GasPrice.ToInt().Int64(), etx.ID, sendError.Error(), ec.store.Config.EthGasBumpPercent())

		sendError = nil
	}

	if sendError == nil {
		return saveSentAttempt(ec.store.GetRawDB(), &a)
	}

	// Any other type of error is considered temporary or resolvable by the
	// node operator. The node may have it in the mempool so we must keep the
	// attempt (leave it in_progress). Safest thing to do is bail out and wait
	// for the next head.
	return errors.Wrapf(sendError, "unexpected error sending eth_tx %v with hash %s", etx.ID, a.Hash.Hex())
}

func deleteInProgressAttempt(db *gorm.DB, a models.EthTxAttempt) error {
	if a.State != models.EthTxAttemptInProgress {
		return errors.New("deleteInProgressAttempt: expected attempt state to be in_progress")
	}
	if a.ID == 0 {
		return errors.New("deleteInProgressAttempt: expected attempt to have an id")
	}
	return errors.Wrap(db.Exec(`DELETE FROM eth_tx_attempts WHERE id = ?`, a.ID).Error, "deleteInProgressAttempt failed")
}

// IsSafeToAbandon determines whether the transaction has an attempt that was
// broadcast long enough ago that we consider it to be "final". Note that this
// is only used in the case the nonce has been used but we cannot get a
// receipt, because we do not have the right attempt in our database.
//
// This should only ever happen if an external wallet has used the account.
func (ec *ethConfirmer) IsSafeToAbandon(etx models.EthTx, blockHeight int64) bool {
	min := int64(0)
	for _, attempt := range etx.EthTxAttempts {
		if attempt.BroadcastBeforeBlockNum != nil && (min == 0 || *attempt.BroadcastBeforeBlockNum < min) {
			min = *attempt.BroadcastBeforeBlockNum
		}
	}
	return min != 0 && min < (blockHeight-int64(ec.config.EthFinalityDepth()))
}

func saveExternalWalletUsedNonce(s *store.Store, etx *models.EthTx, a models.EthTxAttempt) error {
	if etx.State != models.EthTxUnconfirmed {
		return errors.Errorf("can only set external wallet used nonce if unconfirmed, transaction is currently %s", etx.State)
	}
	etx.Nonce = nil
	etx.State = models.EthTxFatalError
	etx.Error = &ErrExternalWalletUsedNonce
	etx.BroadcastAt = nil
	return s.Transaction(func(tx *gorm.DB) error {
		if err := deleteInProgressAttempt(tx, a); err != nil {
			return err
		}
		return errors.Wrap(tx.Save(etx).Error, "saveExternalWalletUsedNonce failed")
	})
}

func saveSentAttempt(db *gorm.DB, a *models.EthTxAttempt) error {
	if a.State != models.EthTxAttemptInProgress {
		return errors.New("expected state to be in_progress")
	}
	a.State = models.EthTxAttemptBroadcast
	return errors.Wrap(db.Save(a).Error, "saveSentAttempt failed")
}

// EnsureConfirmedTransactionsInLongestChain finds all confirmed eth_txes up to the depth
// of the given chain and ensures that every one has a receipt with a block hash that is
// in the given chain.
//
// If any of the confirmed transactions does not have a receipt in the chain, it has been
// re-org'd out and will be rebroadcast.
func (ec *ethConfirmer) EnsureConfirmedTransactionsInLongestChain(head models.Head) error {
	etxs, err := findTransactionsConfirmedAtOrAboveBlockHeight(ec.store.GetRawDB(), head.EarliestInChain().Number)
	if err != nil {
		return err
	}

	for _, etx := range etxs {
		if !hasReceiptInLongestChain(etx, head) {
			if err := ec.markForRebroadcast(etx); err != nil {
				return err
			}
		}
	}

	// Send all the attempts we may have marked for rebroadcast (in_progress state)
	return ec.handleAnyInProgressAttempts(head.Number)
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
		Where("eth_txes.state = 'confirmed' AND block_number >= ?", blockNumber).
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
	a := etx.EthTxAttempts[0]

	// Put it back in progress and delete the receipt
	err := ec.store.Transaction(func(tx *gorm.DB) error {
		if err := deleteAllReceipts(tx, etx.ID); err != nil {
			return err
		}
		if err := unconfirmEthTx(tx, etx); err != nil {
			return err
		}
		return unbroadcastAttempt(tx, a)
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

func unbroadcastAttempt(db *gorm.DB, a models.EthTxAttempt) error {
	if a.State != models.EthTxAttemptBroadcast {
		return errors.New("expected eth_tx_attempt to be broadcast")
	}
	return errors.Wrap(db.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = NULL, state = 'in_progress' WHERE id = ?`, a.ID).Error, "unbroadcastAttempt failed")
}
