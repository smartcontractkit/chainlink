package bulletprooftxmanager

import (
	"bytes"
	"context"
	"encoding/hex"
	"math/big"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Config encompasses config used by bulletprooftxmanager package
// Unless otherwise specified, these should support changing at runtime
//go:generate mockery --recursive --name Config --output ./mocks/ --case=underscore --structname Config --filename config.go
type Config interface {
	ChainID() *big.Int
	EthFinalityDepth() uint
	EthGasBumpPercent() uint16
	EthGasBumpThreshold() uint64
	EthGasBumpTxDepth() uint16
	EthGasBumpWei() *big.Int
	EthGasLimitDefault() uint64
	EthGasLimitMultiplier() float32
	EthGasPriceDefault() *big.Int
	EthMaxGasPriceWei() *big.Int
	EthNonceAutoSync() bool
	EthRPCDefaultBatchSize() uint32
	EthTxReaperThreshold() time.Duration
	EthTxReaperInterval() time.Duration
	EthTxResendAfterThreshold() time.Duration
	TriggerFallbackDBPollInterval() time.Duration
	EthMaxInFlightTransactions() uint32
	OptimismGasFees() bool
}

// KeyStore encompasses the subset of keystore used by bulletprooftxmanager
type KeyStore interface {
	AllKeys() (keys []models.Key, err error)
	SignTx(fromAddress common.Address, tx *gethTypes.Transaction, chainID *big.Int) (*gethTypes.Transaction, error)
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())
}

// For more information about the BulletproofTxManager architecture, see the design doc:
// https://www.notion.so/chainlink/BulletproofTxManager-Architecture-Overview-9dc62450cd7a443ba9e7dceffa1a8d6b

var (
	promNumGasBumps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_num_gas_bumps",
		Help: "Number of gas bumps",
	})

	promGasBumpExceedsLimit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_gas_bump_exceeds_limit",
		Help: "Number of times gas bumping failed from exceeding the configured limit. Any counts of this type indicate a serious problem.",
	})

	promRevertedTxCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_num_tx_reverted",
		Help: "Number of times a transaction reverted on-chain",
	})
)

var _ models.HeadTrackable = &BulletproofTxManager{}

type BulletproofTxManager struct {
	utils.StartStopOnce

	db               *gorm.DB
	ethClient        eth.Client
	config           Config
	keyStore         KeyStore
	advisoryLocker   postgres.AdvisoryLocker
	eventBroadcaster postgres.EventBroadcaster

	chHeads chan models.Head
	trigger chan common.Address

	chStop chan struct{}
	wg     sync.WaitGroup

	reaper      *Reaper
	ethResender *EthResender
}

func NewBulletproofTxManager(db *gorm.DB, ethClient eth.Client, config Config, keyStore KeyStore, advisoryLocker postgres.AdvisoryLocker, eventBroadcaster postgres.EventBroadcaster) *BulletproofTxManager {
	b := BulletproofTxManager{
		StartStopOnce:    utils.StartStopOnce{},
		db:               db,
		ethClient:        ethClient,
		config:           config,
		keyStore:         keyStore,
		advisoryLocker:   advisoryLocker,
		eventBroadcaster: eventBroadcaster,
		chHeads:          make(chan models.Head),
		trigger:          make(chan common.Address),
		chStop:           make(chan struct{}),
	}
	if config.EthTxResendAfterThreshold() > 0 {
		b.ethResender = NewEthResender(db, ethClient, defaultResenderPollInterval, config)
	} else {
		logger.Info("EthResender: Disabled")
	}
	if config.EthTxReaperThreshold() > 0 {
		b.reaper = NewReaper(db, config)
	} else {
		logger.Info("EthTxReaper: Disabled")
	}
	return &b
}

func (b *BulletproofTxManager) Start() (merr error) {
	return b.StartOnce("BulletproofTxManager", func() error {
		keys, err := b.keyStore.AllKeys()
		if err != nil {
			return errors.Wrap(err, "BulletproofTxManager: failed to load keys")
		}

		logger.Debugw("BulletproofTxManager: booting", "keys", keys)

		eb := NewEthBroadcaster(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, b.eventBroadcaster, keys)
		ec := NewEthConfirmer(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, keys)
		if err := eb.Start(); err != nil {
			return errors.Wrap(err, "BulletproofTxManager: EthBroadcaster failed to start")
		}
		if err := ec.Start(); err != nil {
			return errors.Wrap(err, "BulletproofTxManager: EthConfirmer failed to start")
		}

		b.wg.Add(1)
		go b.runLoop(eb, ec)

		if b.reaper != nil {
			b.reaper.Start()
		}

		if b.ethResender != nil {
			b.ethResender.Start()
		}

		return nil
	})
}

func (b *BulletproofTxManager) Close() (merr error) {
	return b.StopOnce("BulletproofTxManager", func() error {
		close(b.chStop)

		if b.reaper != nil {
			b.reaper.Stop()
		}
		if b.ethResender != nil {
			b.ethResender.Stop()
		}

		b.wg.Wait()

		return nil
	})
}

func (b *BulletproofTxManager) runLoop(eb *EthBroadcaster, ec *EthConfirmer) {
	defer b.wg.Done()
	keysChanged, unsub := b.keyStore.SubscribeToKeyChanges()
	defer unsub()

	for {
		select {
		case address := <-b.trigger:
			eb.Trigger(address)
		case head := <-b.chHeads:
			ec.mb.Deliver(head)
		case <-b.chStop:
			logger.ErrorIfCalling(eb.Close)
			logger.ErrorIfCalling(ec.Close)
			return
		case <-keysChanged:
			keys, err := b.keyStore.AllKeys()
			if err != nil {
				logger.Fatalf("BulletproofTxManager: expected keystore to be unlocked: %s", err.Error())
			}

			logger.Debugw("BulletproofTxManager: keys changed, reloading", "keys", keys)

			logger.ErrorIfCalling(eb.Close)
			logger.ErrorIfCalling(ec.Close)

			eb = NewEthBroadcaster(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, b.eventBroadcaster, keys)
			ec = NewEthConfirmer(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, keys)

			logger.ErrorIfCalling(eb.Start)
			logger.ErrorIfCalling(ec.Start)
		}
	}
}

// OnNewLongestChain conforms to HeadTrackable
func (b *BulletproofTxManager) OnNewLongestChain(ctx context.Context, head models.Head) {
	ok := b.IfStarted(func() {
		if b.reaper != nil {
			b.reaper.SetLatestBlockNum(head.Number)
		}
		select {
		case b.chHeads <- head:
		case <-ctx.Done():
			logger.Errorw("BulletproofTxManager: timed out handling head", "blockNum", head.Number, "ctxErr", ctx.Err())
		}
	})
	if !ok {
		logger.Debugw("BulletproofTxManager: not started; ignoring head", "head", head, "state", b.State())
	}
}

// Trigger forces the EthBroadcaster to check early for the given address
func (b *BulletproofTxManager) Trigger(addr common.Address) {
	select {
	case b.trigger <- addr:
	default:
	}
}

// Connect solely exists to conform to HeadTrackable
func (b *BulletproofTxManager) Connect(*models.Head) error {
	return nil
}

// Disconnect solely exists to conform to HeadTrackable
func (b *BulletproofTxManager) Disconnect() {}

const optimismGasPrice int64 = 1e9 // 1 GWei

// SendEther creates a transaction that transfers the given value of ether
func SendEther(db *gorm.DB, from, to common.Address, value assets.Eth, gasLimit uint64) (etx models.EthTx, err error) {
	if to == utils.ZeroAddress {
		return etx, errors.New("cannot send ether to zero address")
	}
	etx = models.EthTx{
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: []byte{},
		Value:          value,
		GasLimit:       gasLimit,
		State:          models.EthTxUnstarted,
	}
	err = db.Create(&etx).Error
	return etx, err
}

func newAttempt(ctx context.Context, ethClient eth.Client, ks KeyStore, config Config, etx models.EthTx, suggestedGasPrice *big.Int) (models.EthTxAttempt, error) {
	attempt := models.EthTxAttempt{}

	gasPrice := config.EthGasPriceDefault()
	if suggestedGasPrice != nil {
		gasPrice = suggestedGasPrice
	}

	if config.OptimismGasFees() {
		// Optimism requires special handling, it assumes that clients always call EstimateGas
		callMsg := ethereum.CallMsg{To: &etx.ToAddress, From: etx.FromAddress, Data: etx.EncodedPayload}
		gasLimit, estimateGasErr := ethClient.EstimateGas(ctx, callMsg)
		if estimateGasErr != nil {
			return attempt, errors.Wrapf(estimateGasErr, "error getting gas price for new transaction %v", etx.ID)
		}
		etx.GasLimit = gasLimit
		gasPrice = big.NewInt(optimismGasPrice)
	}

	gasLimit := decimal.NewFromBigInt(big.NewInt(0).SetUint64(etx.GasLimit), 0).Mul(decimal.NewFromFloat32(config.EthGasLimitMultiplier())).IntPart()
	etx.GasLimit = (uint64)(gasLimit)

	transaction := gethTypes.NewTransaction(uint64(*etx.Nonce), etx.ToAddress, etx.Value.ToInt(), etx.GasLimit, gasPrice, etx.EncodedPayload)
	hash, signedTxBytes, err := signTx(ks, etx.FromAddress, transaction, config.ChainID())
	if err != nil {
		return attempt, errors.Wrapf(err, "error using account %s to sign transaction %v", etx.FromAddress.String(), etx.ID)
	}

	attempt.State = models.EthTxAttemptInProgress
	attempt.SignedRawTx = signedTxBytes
	attempt.EthTxID = etx.ID
	attempt.GasPrice = *utils.NewBig(gasPrice)
	attempt.Hash = hash

	return attempt, nil
}

func signTx(keyStore KeyStore, address common.Address, tx *gethTypes.Transaction, chainID *big.Int) (common.Hash, []byte, error) {
	signedTx, err := keyStore.SignTx(address, tx, chainID)
	if err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "signTx failed")
	}
	rlp := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(rlp); err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "signTx failed")
	}
	return signedTx.Hash(), rlp.Bytes(), nil

}

// send broadcasts the transaction to the ethereum network, writes any relevant
// data onto the attempt and returns an error (or nil) depending on the status
func sendTransaction(ctx context.Context, ethClient eth.Client, a models.EthTxAttempt, e models.EthTx) *eth.SendError {
	signedTx, err := a.GetSignedTx()
	if err != nil {
		return eth.NewFatalSendError(err)
	}

	ctx, cancel := eth.DefaultQueryCtx(ctx)
	defer cancel()

	err = ethClient.SendTransaction(ctx, signedTx)
	err = errors.WithStack(err)

	logger.Debugw("BulletproofTxManager: Sending transaction", "ethTxAttemptID", a.ID, "txHash", signedTx.Hash(), "gasPriceWei", a.GasPrice.ToInt().Int64(), "err", err, "meta", e.Meta)
	sendErr := eth.NewSendError(err)
	if sendErr.IsTransactionAlreadyInMempool() {
		logger.Debugw("transaction already in mempool", "txHash", signedTx.Hash(), "nodeErr", sendErr.Error())
		return nil
	}
	return eth.NewSendError(err)
}

// sendEmptyTransaction sends a transaction with 0 Eth and an empty payload to the burn address
// May be useful for clearing stuck nonces
func sendEmptyTransaction(
	ethClient eth.Client,
	keyStore KeyStore,
	nonce uint64,
	gasLimit uint64,
	gasPriceWei *big.Int,
	fromAddress common.Address,
	chainID *big.Int,
) (_ *gethTypes.Transaction, err error) {
	defer utils.WrapIfError(&err, "sendEmptyTransaction failed")

	signedTx, err := makeEmptyTransaction(keyStore, nonce, gasLimit, gasPriceWei, fromAddress, chainID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := eth.DefaultQueryCtx()
	defer cancel()
	err = ethClient.SendTransaction(ctx, signedTx)
	return signedTx, err
}

// makes a transaction that sends 0 eth to self
func makeEmptyTransaction(keyStore KeyStore, nonce uint64, gasLimit uint64, gasPriceWei *big.Int, fromAddress common.Address, chainID *big.Int) (*gethTypes.Transaction, error) {
	value := big.NewInt(0)
	payload := []byte{}
	tx := gethTypes.NewTransaction(nonce, fromAddress, value, gasLimit, gasPriceWei, payload)
	return keyStore.SignTx(fromAddress, tx, chainID)
}

func saveReplacementInProgressAttempt(db *gorm.DB, oldAttempt models.EthTxAttempt, replacementAttempt *models.EthTxAttempt) error {
	if oldAttempt.State != models.EthTxAttemptInProgress || replacementAttempt.State != models.EthTxAttemptInProgress {
		return errors.New("expected attempts to be in_progress")
	}
	if oldAttempt.ID == 0 {
		return errors.New("expected oldAttempt to have an ID")
	}
	return postgres.GormTransactionWithDefaultContext(db, func(tx *gorm.DB) error {
		if err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE id = ? `, oldAttempt.ID).Error; err != nil {
			return errors.Wrap(err, "saveReplacementInProgressAttempt failed")
		}
		return errors.Wrap(tx.Create(replacementAttempt).Error, "saveReplacementInProgressAttempt failed")
	})
}

// BumpGas computes the next gas price to attempt as the largest of:
// - A configured percentage bump (ETH_GAS_BUMP_PERCENT) on top of the baseline price.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func BumpGas(config Config, originalGasPrice *big.Int) (*big.Int, error) {
	baselinePrice := max(originalGasPrice, config.EthGasPriceDefault())

	var priceByPercentage = new(big.Int)
	priceByPercentage.Mul(baselinePrice, big.NewInt(int64(100+config.EthGasBumpPercent())))
	priceByPercentage.Div(priceByPercentage, big.NewInt(100))

	var priceByIncrement = new(big.Int)
	priceByIncrement.Add(baselinePrice, config.EthGasBumpWei())

	bumpedGasPrice := max(priceByPercentage, priceByIncrement)
	if bumpedGasPrice.Cmp(config.EthMaxGasPriceWei()) > 0 {
		promGasBumpExceedsLimit.Inc()
		return config.EthMaxGasPriceWei(), errors.Errorf("bumped gas price of %s would exceed configured max gas price of %s (original price was %s). %s",
			bumpedGasPrice.String(), config.EthMaxGasPriceWei(), originalGasPrice.String(), EthNodeConnectivityProblemLabel)
	} else if bumpedGasPrice.Cmp(originalGasPrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// ETH_GAS_BUMP_PERCENT and ETH_GAS_BUMP_WEI in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedGasPrice, errors.Errorf("bumped gas price of %s is equal to original gas price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI", bumpedGasPrice.String(), originalGasPrice.String())
	}
	promNumGasBumps.Inc()
	return bumpedGasPrice, nil
}

func max(a, b *big.Int) *big.Int {
	if a.Cmp(b) >= 0 {
		return a
	}
	return b
}

// CountUnconfirmedTransactions returns the number of unconfirmed transactions
func CountUnconfirmedTransactions(db *gorm.DB, fromAddress common.Address) (count uint32, err error) {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err = db.WithContext(ctx).Raw(`SELECT count(*) FROM eth_txes WHERE from_address = ? AND state = 'unconfirmed'`, fromAddress).Scan(&count).Error
	return
}

// CreateEthTransaction inserts a new transaction
func CreateEthTransaction(db *gorm.DB, fromAddress, toAddress common.Address, payload []byte, gasLimit, maxUnconfirmedTransactions uint64) (etx models.EthTx, err error) {
	err = CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
	if err != nil {
		return etx, errors.Wrap(err, "transmitter#CreateEthTransaction")
	}

	value := 0
	// NOTE: It is important to remember that eth_tx_attempts with state
	// insufficient_eth can actually hang around long after the node has been
	// refunded and started sending transactions again.
	// This is because they are not ever deleted if attached to an eth_tx that
	// is moved into confirmed/fatal_error state
	res := db.Raw(`
INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at)
(
SELECT ?,?,?,?,?,'unstarted',NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM eth_tx_attempts
	JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id
	WHERE eth_txes.from_address = ?
		AND eth_txes.state = 'unconfirmed'
		AND eth_tx_attempts.state = 'insufficient_eth'
)
)
RETURNING "eth_txes".*
`, fromAddress, toAddress, payload, value, gasLimit, fromAddress).Scan(&etx)
	err = res.Error
	if err != nil {
		return etx, errors.Wrap(err, "transmitter failed to insert eth_tx")
	}

	if res.RowsAffected == 0 {
		err = errors.Errorf("wallet is out of eth: %s", fromAddress.Hex())
		logger.Warnw(err.Error(),
			"fromAddress", fromAddress,
			"toAddress", toAddress,
			"payload", "0x"+hex.EncodeToString(payload),
			"value", value,
			"gasLimit", gasLimit,
		)
	}
	return
}

const EthMaxInFlightTransactionsWarningLabel = `WARNING: You may need to increase ETH_MAX_IN_FLIGHT_TRANSACTIONS to boost your node's transaction throughput, however you do this at your own risk. You MUST first ensure your ethereum node is configured not to ever evict local transactions that exceed this number otherwise the node can get permanently stuck.`

const EthMaxQueuedTransactionsLabel = `WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse.`

const EthNodeConnectivityProblemLabel = `WARNING: If this keeps happening it may be a sign that your eth node has a connectivity problem, and your transactions are not making it to any miners.`

// CheckEthTxQueueCapacity returns an error if inserting this transaction would
// exceed the maximum queue size.
func CheckEthTxQueueCapacity(db *gorm.DB, fromAddress common.Address, maxQueuedTransactions uint64) (err error) {
	if maxQueuedTransactions == 0 {
		return nil
	}
	var count uint64
	err = db.Raw(`SELECT count(*) FROM eth_txes WHERE from_address = ? AND state = 'unstarted'`, fromAddress).Scan(&count).Error
	if err != nil {
		err = errors.Wrap(err, "bulletprooftxmanager.CheckEthTxQueueCapacity query failed")
		return
	}

	if count >= maxQueuedTransactions {
		err = errors.Errorf("cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, EthMaxQueuedTransactionsLabel)
	}
	return
}
