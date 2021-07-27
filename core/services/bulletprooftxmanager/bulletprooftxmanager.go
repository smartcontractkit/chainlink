package bulletprooftxmanager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/gas"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Config encompasses config used by bulletprooftxmanager package
// Unless otherwise specified, these should support changing at runtime
//go:generate mockery --recursive --name Config --output ./mocks/ --case=underscore --structname Config --filename config.go
type Config interface {
	BlockHistoryEstimatorBatchSize() uint32
	BlockHistoryEstimatorBlockDelay() uint16
	BlockHistoryEstimatorBlockHistorySize() uint16
	BlockHistoryEstimatorTransactionPercentile() uint16
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
	EthMaxInFlightTransactions() uint32
	EthMaxQueuedTransactions() uint64
	EthMinGasPriceWei() *big.Int
	EthNonceAutoSync() bool
	EthRPCDefaultBatchSize() uint32
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EthTxResendAfterThreshold() time.Duration
	GasEstimatorMode() string
	TriggerFallbackDBPollInterval() time.Duration
}

// KeyStore encompasses the subset of keystore used by bulletprooftxmanager
type KeyStore interface {
	AllKeys() (keys []ethkey.Key, err error)
	SignTx(fromAddress common.Address, tx *gethTypes.Transaction, chainID *big.Int) (*gethTypes.Transaction, error)
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())
}

// For more information about the BulletproofTxManager architecture, see the design doc:
// https://www.notion.so/chainlink/BulletproofTxManager-Architecture-Overview-9dc62450cd7a443ba9e7dceffa1a8d6b

var (
	promRevertedTxCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_num_tx_reverted",
		Help: "Number of times a transaction reverted on-chain",
	})
)

var _ TxManager = &BulletproofTxManager{}

//go:generate mockery --recursive --name TxManager --output ./mocks/ --case=underscore --structname TxManager --filename tx_manager.go
type TxManager interface {
	httypes.HeadTrackable
	service.Service
	Trigger(addr common.Address)
	CreateEthTransaction(db *gorm.DB, fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, meta interface{}, strategy TxStrategy) (etx EthTx, err error)
	GetGasEstimator() gas.Estimator
}

type BulletproofTxManager struct {
	utils.StartStopOnce

	db               *gorm.DB
	ethClient        eth.Client
	config           Config
	keyStore         KeyStore
	advisoryLocker   postgres.AdvisoryLocker
	eventBroadcaster postgres.EventBroadcaster
	gasEstimator     gas.Estimator

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
	b.gasEstimator = gas.NewEstimator(ethClient, config)

	return &b
}

func (b *BulletproofTxManager) Start() (merr error) {
	return b.StartOnce("BulletproofTxManager", func() error {
		keys, err := b.keyStore.AllKeys()
		if err != nil {
			return errors.Wrap(err, "BulletproofTxManager: failed to load keys")
		}

		logger.Debugw("BulletproofTxManager: booting", "keys", keys)

		eb := NewEthBroadcaster(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, b.eventBroadcaster, keys, b.gasEstimator)
		ec := NewEthConfirmer(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, keys, b.gasEstimator)
		if err := eb.Start(); err != nil {
			return errors.Wrap(err, "BulletproofTxManager: EthBroadcaster failed to start")
		}
		if err := ec.Start(); err != nil {
			return errors.Wrap(err, "BulletproofTxManager: EthConfirmer failed to start")
		}

		if err := b.gasEstimator.Start(); err != nil {
			return errors.Wrap(err, "BulletproofTxManager: Estimator failed to start")
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

		b.gasEstimator.Close()

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

			eb = NewEthBroadcaster(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, b.eventBroadcaster, keys, b.gasEstimator)
			ec = NewEthConfirmer(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, keys, b.gasEstimator)

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
		b.gasEstimator.OnNewLongestChain(ctx, head)
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

// CreateEthTransaction inserts a new transaction
func (b *BulletproofTxManager) CreateEthTransaction(db *gorm.DB, fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, meta interface{}, strategy TxStrategy) (etx EthTx, err error) {
	err = CheckEthTxQueueCapacity(db, fromAddress, b.config.EthMaxQueuedTransactions())
	if err != nil {
		return etx, errors.Wrap(err, "BulletproofTxManager#CreateEthTransaction")
	}

	// meta can hold arbitrary data and is mostly useful for logging/debugging
	var metaBytes []byte
	if meta != nil {
		metaBytes, err = json.Marshal(meta)
		if err != nil {
			return etx, errors.Wrap(err, "BulletproofTxManager#CreateEthTransaction failed to marshal ethtx metadata")
		}
	}

	value := 0
	err = postgres.GormTransactionWithDefaultContext(db, func(tx *gorm.DB) error {
		res := tx.Raw(`
INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject)
VALUES (
?,?,?,?,?,'unstarted',NOW(),?,?
)
RETURNING "eth_txes".*
`, fromAddress, toAddress, payload, value, gasLimit, metaBytes, strategy.Subject()).Scan(&etx)
		err = res.Error
		if err != nil {
			return errors.Wrap(err, "BulletproofTxManager#CreateEthTransaction failed to insert eth_tx")
		}

		pruned, err := strategy.PruneQueue(tx)
		if err != nil {
			return errors.Wrap(err, "BulletproofTxManager#CreateEthTransaction failed to prune eth_txes")
		}
		if pruned > 0 {
			logger.Warnw(fmt.Sprintf("BulletproofTxManager: dropped %d old transactions from transaction queue", pruned), "fromAddress", fromAddress, "toAddress", toAddress, "meta", meta, "subject", strategy.Subject(), "replacementID", etx.ID)
		}
		return nil
	})
	return
}

// GetGasEstimator returns the gas estimator, mostly useful for tests
func (b *BulletproofTxManager) GetGasEstimator() gas.Estimator {
	return b.gasEstimator
}

// SendEther creates a transaction that transfers the given value of ether
func SendEther(db *gorm.DB, from, to common.Address, value assets.Eth, gasLimit uint64) (etx EthTx, err error) {
	if to == utils.ZeroAddress {
		return etx, errors.New("cannot send ether to zero address")
	}
	etx = EthTx{
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: []byte{},
		Value:          value,
		GasLimit:       gasLimit,
		State:          EthTxUnstarted,
	}
	err = db.Create(&etx).Error
	return etx, err
}

func newAttempt(ethClient eth.Client, ks KeyStore, chainID *big.Int, etx EthTx, gasPrice *big.Int, gasLimit uint64) (EthTxAttempt, error) {
	attempt := EthTxAttempt{}

	tx := newLegacyTransaction(
		uint64(*etx.Nonce),
		etx.ToAddress,
		etx.Value.ToInt(),
		gasLimit,
		gasPrice,
		etx.EncodedPayload,
	)

	transaction := gethTypes.NewTx(&tx)
	hash, signedTxBytes, err := signTx(ks, etx.FromAddress, transaction, chainID)
	if err != nil {
		return attempt, errors.Wrapf(err, "error using account %s to sign transaction %v", etx.FromAddress.String(), etx.ID)
	}

	attempt.State = EthTxAttemptInProgress
	attempt.SignedRawTx = signedTxBytes
	attempt.EthTxID = etx.ID
	attempt.GasPrice = *utils.NewBig(gasPrice)
	attempt.Hash = hash

	return attempt, nil
}

func newLegacyTransaction(nonce uint64, to common.Address, value *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) gethTypes.LegacyTx {
	return gethTypes.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}
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
func sendTransaction(ctx context.Context, ethClient eth.Client, a EthTxAttempt, e EthTx) *eth.SendError {
	signedTx, err := a.GetSignedTx()
	if err != nil {
		return eth.NewFatalSendError(err)
	}

	ctx, cancel := eth.DefaultQueryCtx(ctx)
	defer cancel()

	err = ethClient.SendTransaction(ctx, signedTx)
	err = errors.WithStack(err)

	logger.Debugw("BulletproofTxManager: Sent transaction", "ethTxAttemptID", a.ID, "txHash", signedTx.Hash(), "gasPriceWei", a.GasPrice.ToInt().Int64(), "err", err, "meta", e.Meta, "gasLimit", e.GasLimit)
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

func saveReplacementInProgressAttempt(db *gorm.DB, oldAttempt EthTxAttempt, replacementAttempt *EthTxAttempt) error {
	if oldAttempt.State != EthTxAttemptInProgress || replacementAttempt.State != EthTxAttemptInProgress {
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

// CountUnconfirmedTransactions returns the number of unconfirmed transactions
func CountUnconfirmedTransactions(db *gorm.DB, fromAddress common.Address) (count uint32, err error) {
	return countTransactionsWithState(db, fromAddress, EthTxUnconfirmed)
}

// CountUnstartedTransactions returns the number of unconfirmed transactions
func CountUnstartedTransactions(db *gorm.DB, fromAddress common.Address) (count uint32, err error) {
	return countTransactionsWithState(db, fromAddress, EthTxUnstarted)
}

func countTransactionsWithState(db *gorm.DB, fromAddress common.Address, state EthTxState) (count uint32, err error) {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err = db.WithContext(ctx).Raw(`SELECT count(*) FROM eth_txes WHERE from_address = ? AND state = ?`, fromAddress, state).Scan(&count).Error
	return
}

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
		err = errors.Errorf("cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, static.EthMaxQueuedTransactionsLabel)
	}
	return
}

var _ TxManager = &NullTxManager{}

type NullTxManager struct {
	ErrMsg string
}

func (n *NullTxManager) Connect(*models.Head) error                     { return errors.New(n.ErrMsg) }
func (n *NullTxManager) OnNewLongestChain(context.Context, models.Head) {}
func (n *NullTxManager) Start() error                                   { return errors.New(n.ErrMsg) }
func (n *NullTxManager) Close() error                                   { return errors.New(n.ErrMsg) }
func (n *NullTxManager) Trigger(common.Address)                         { panic(n.ErrMsg) }
func (n *NullTxManager) CreateEthTransaction(*gorm.DB, common.Address, common.Address, []byte, uint64, interface{}, TxStrategy) (etx EthTx, err error) {
	return etx, errors.New(n.ErrMsg)
}
func (n *NullTxManager) Healthy() error                 { return nil }
func (n *NullTxManager) Ready() error                   { return nil }
func (n *NullTxManager) GetGasEstimator() gas.Estimator { return nil }
