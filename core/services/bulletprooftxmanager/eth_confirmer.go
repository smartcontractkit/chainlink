package bulletprooftxmanager

import (
	"context"
	"encoding/json"
	"math/big"
	"sync"
	"time"

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
			logger.Error("EthConfirmer: ", err)
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

	if err := ec.CheckForReceipts(ctx); err != nil {
		return errors.Wrap(err, "CheckForReceipts failed")
	}

	logger.Debugw("EthConfirmer: finished CheckForReceipts", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	if err := ec.BumpGasWhereNecessary(ctx, head.Number); err != nil {
		return errors.Wrap(err, "BumpGasWhereNecessary failed")
	}

	logger.Debugw("EthConfirmer: finished BumpGasWhereNecessary", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	mark = time.Now()

	defer func() {
		logger.Debugw("EthConfirmer: finished EnsureConfirmedTransactionsInLongestChain", "headNum", head.Number, "time", time.Since(mark), "id", "eth_confirmer")
	}()

	return errors.Wrap(ec.EnsureConfirmedTransactionsInLongestChain(ctx, head), "EnsureConfirmedTransactionsInLongestChain failed")
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

func (ec *ethConfirmer) CheckForReceipts(ctx context.Context) error {
	unconfirmedEtxs, err := ec.findUnconfirmedEthTxs()
	if err != nil {
		return errors.Wrap(err, "findUnconfirmedEthTxs failed")
	}
	if len(unconfirmedEtxs) > 0 {
		logger.Debugf("EthConfirmer: %v unconfirmed transactions", len(unconfirmedEtxs))
	}
	wg := sync.WaitGroup{}
	wg.Add(receiptFetcherWorkerCount)
	chEthTxes := make(chan models.EthTx)
	for i := 0; i < receiptFetcherWorkerCount; i++ {
		go ec.fetchReceipts(ctx, chEthTxes, &wg)
	}
	for _, etx := range unconfirmedEtxs {
		chEthTxes <- etx
	}
	close(chEthTxes)
	wg.Wait()
	return nil
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
			if isParityQueriedReceiptTooEarly(err) || (receipt != nil && receipt.BlockNumber == nil) {
				logger.Debugw("EthConfirmer#fetchReceipts: got receipt for transaction but it's still in the mempool and not included in a block yet", "txHash", attempt.Hash.Hex())
				break
			} else if err != nil {
				logger.Errorf("EthConfirmer#fetchReceipts: fetchReceipt failed for transaction %s", attempt.Hash.Hex())
				break
			}
			if receipt != nil {
				logger.Debugw("EthConfirmer#fetchReceipts: got receipt for transaction", "txHash", attempt.Hash.Hex(), "blockNumber", receipt.BlockNumber)
				if receipt.TxHash != attempt.Hash {
					logger.Errorf("EthConfirmer#fetchReceipts: invariant violation, expected receipt with hash %s to have same hash as attempt with hash %s", receipt.TxHash.Hex(), attempt.Hash.Hex())
					break
				}
				if err := ec.saveReceipt(*receipt, etx.ID); err != nil {
					logger.Errorf("EthConfirmer#fetchReceipts: saveReceipt failed")
					break
				}
				break
			} else {
				logger.Debugw("EthConfirmer#fetchReceipts: still waiting for receipt", "txHash", attempt.Hash.Hex(), "ethTxAttemptID", attempt.ID, "ethTxID", etx.ID)
			}
		}
	}
}

func (ec *ethConfirmer) findUnconfirmedEthTxs() ([]models.EthTx, error) {
	var etxs []models.EthTx
	err := ec.store.DB.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC")
		}).
		Order("nonce ASC").
		Find(&etxs, "eth_txes.state = 'unconfirmed'").Error
	return etxs, err
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

func (ec *ethConfirmer) BumpGasWhereNecessary(ctx context.Context, blockHeight int64) error {
	if err := ec.handleAnyInProgressAttempts(ctx, blockHeight); err != nil {
		return errors.Wrap(err, "handleAnyInProgressAttempts failed")
	}

	etxs, err := FindEthTxsRequiringNewAttempt(ec.store.DB, blockHeight, int64(ec.config.EthGasBumpThreshold()))
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

		if err := ec.handleInProgressAttempt(ctx, etx, attempt, blockHeight, true); err != nil {
			return errors.Wrap(err, "handleInProgressAttempt failed")
		}
	}
	return nil
}

// "in_progress" attempts were left behind after a crash/restart and may or may not have been sent
// We should try to ensure they get on-chain so we can fetch a receipt for them
func (ec *ethConfirmer) handleAnyInProgressAttempts(ctx context.Context, blockHeight int64) error {
	attempts, err := getInProgressEthTxAttempts(ec.store)
	if err != nil {
		return errors.Wrap(err, "getInProgressEthTxAttempts failed")
	}
	for _, a := range attempts {
		if err := ec.handleInProgressAttempt(ctx, a.EthTx, a, blockHeight, false); err != nil {
			return errors.Wrap(err, "handleInProgressAttempt failed")
		}
	}
	return nil
}

func getInProgressEthTxAttempts(s *store.Store) ([]models.EthTxAttempt, error) {
	var attempts []models.EthTxAttempt
	err := s.DB.
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
			"AND eth_tx_attempts.state != 'insufficient_eth' "+
			"AND (broadcast_before_block_num > ? OR broadcast_before_block_num IS NULL OR eth_tx_attempts.state != 'broadcast')", blockNum-gasBumpThreshold).
		Order("nonce ASC").
		Where("eth_txes.state = 'unconfirmed' AND eth_tx_attempts.id IS NULL").
		Find(&etxs).Error

	return etxs, errors.Wrap(err, "FindEthTxsRequiringNewAttempt failed")
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
			return attempt, errors.Wrapf(err, "could not create newAttemptWithGasBump")
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

func (ec *ethConfirmer) handleInProgressAttempt(ctx context.Context, etx models.EthTx, attempt models.EthTxAttempt, blockHeight int64, isVirginAttempt bool) error {
	if attempt.State != models.EthTxAttemptInProgress {
		return errors.Errorf("invariant violation: expected eth_tx_attempt %v to be in_progress, it was %s", attempt.ID, attempt.State)
	}

	sendError := sendTransaction(ctx, ec.ethClient, attempt)

	if sendError.IsTerminallyUnderpriced() {
		// This should really not ever happen in normal operation since we
		// already bumped above the required minimum in ethBroadcaster.
		//
		// It could concievably happen if the remote eth node changed it's configuration.
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
		return ec.handleInProgressAttempt(ctx, etx, replacementAttempt, blockHeight, isVirginAttempt)
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
		logger.Errorf("invariant violation: fatal error while reattempting transaction %v: '%s'. "+
			"SignedRawTx: %s\n"+
			"BlockHeight: %v\n"+
			"IsVirginAttempt: %v\n"+
			"ACTION REQUIRED: Your node is BROKEN - this error should never happen in normal operation. "+
			"Please consider raising an issue here: https://github.com/smartcontractkit/chainlink/issues", etx.ID, sendError, hexutil.Encode(attempt.SignedRawTx), blockHeight, isVirginAttempt)
		// This will loop continuously on every new head so it must be handled manually by the node operator!
		return deleteInProgressAttempt(ec.store.DB, attempt)
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

			return saveExternalWalletUsedNonce(ec.store, &etx, attempt)
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
			return deleteInProgressAttempt(ec.store.DB, attempt)
		}
		// If we already sent the attempt, we have to assume the one who was
		// confirmed was this one, so simply mark it as broadcast and wait for
		// a receipt.
		//
		// Assume success and hand off to the next cycle.
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
		logger.Errorf("EthConfirmer: replacement transaction underpriced at %v wei for eth_tx %v. "+
			"Eth node returned error: '%s'. "+
			"Either you have set ETH_GAS_BUMP_PERCENT (currently %v%%) too low or an external wallet used this account. "+
			"Please note that using your node's private keys outside of the chainlink node is NOT SUPPORTED and can lead to missed transactions.",
			attempt.GasPrice.ToInt().Int64(), etx.ID, sendError.Error(), ec.store.Config.EthGasBumpPercent())

		// Assume success and hand off to the next cycle.
		sendError = nil
	}

	if sendError.IsInsufficientEth() {
		logger.Errorf("EthConfirmer: EthTxAttempt %v (hash 0x%x) at gas price (%s Wei) was rejected due to insufficient eth. "+
			"The eth node returned %s. "+
			"ACTION REQUIRED: Chainlink wallet with address 0x%x is OUT OF FUNDS",
			attempt.ID, attempt.Hash, attempt.GasPrice.String(), sendError.Error(), etx.FromAddress,
		)
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

func saveExternalWalletUsedNonce(s *store.Store, etx *models.EthTx, attempt models.EthTxAttempt) error {
	if etx.State != models.EthTxUnconfirmed {
		return errors.Errorf("can only set external wallet used nonce if unconfirmed, transaction is currently %s", etx.State)
	}
	etx.Nonce = nil
	etx.State = models.EthTxFatalError
	etx.Error = &ErrExternalWalletUsedNonce
	etx.BroadcastAt = nil
	return s.Transaction(func(tx *gorm.DB) error {
		if err := deleteInProgressAttempt(tx, attempt); err != nil {
			return errors.Wrap(err, "deleteInProgressAttempt failed")
		}
		return errors.Wrap(tx.Save(etx).Error, "saveExternalWalletUsedNonce failed")
	})
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
func (ec *ethConfirmer) EnsureConfirmedTransactionsInLongestChain(ctx context.Context, head models.Head) error {
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

	// Send all the attempts we may have marked for rebroadcast (in_progress state)
	return ec.handleAnyInProgressAttempts(ctx, head.Number)
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
	attempt := etx.EthTxAttempts[0]

	// Put it back in progress and delete the receipt
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
	logger.Info("ForceRebroadcast: will rebroadcast transactions for all nonces between %v and %v", beginningNonce, endingNonce)

	for n := beginningNonce; n <= endingNonce; n++ {
		etx, err := findEthTxWithNonce(ec.store.DB, address, n)
		if err != nil {
			return errors.Wrap(err, "ForceRebroadcast failed")
		}
		if etx == nil {
			logger.Debugf("ForceRebroadcast: no eth_tx found with nonce %v, will rebroadcast empty transaction", n)
			hash, err := ec.sendEmptyTransaction(context.TODO(), address, n, overrideGasLimit, gasPriceWei)
			if err != nil {
				logger.Errorf("ForceRebroadcast: failed to send empty transaction with nonce %v: %s", n, err.Error())
				continue
			}
			logger.Infof("ForceRebroadcast: successfully rebroadcast empty transaction with nonce %v and hash %s", n, hash)
		} else {
			logger.Debugf("ForceRebroadcast: got eth_tx %v with nonce %v, will rebroadcast this transaction", etx.ID, etx.Nonce)
			if overrideGasLimit != 0 {
				etx.GasLimit = overrideGasLimit
			}
			attempt, err := newAttempt(ec.store, *etx, big.NewInt(int64(gasPriceWei)))
			if err != nil {
				logger.Errorf("ForceRebroadcast: failed to create new attempt for eth_tx %v: %s", etx.ID, err.Error())
				continue
			}
			if err := sendTransaction(context.TODO(), ec.ethClient, attempt); err != nil {
				logger.Errorf("ForceRebroadcast: failed to rebroadcast eth_tx %v with nonce %v at gas price %s wei and gas limit %v: %s", etx.ID, *etx.Nonce, attempt.GasPrice.String(), etx.GasLimit, err.Error())
				continue
			}
			logger.Infof("ForceRebroadcast: successfully rebroadcast eth_tx %v with hash: %s", etx.ID, attempt.Hash)
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

// findAllEthTxsInNonceRange returns an array of eth_txes for the given key
// matching the inclusive range between beginningNonce and endingNonce
func findEthTxWithNonce(db *gorm.DB, fromAddress gethCommon.Address, nonce uint) (*models.EthTx, error) {
	etx := models.EthTx{}
	err := db.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("eth_tx_attempts.gas_price DESC")
		}).
		First(&etx, "from_address = ? AND nonce = ? AND state IN ('confirmed','unconfirmed')", fromAddress, nonce).
		Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &etx, errors.Wrap(err, "findEthTxsWithNonce failed")
}
