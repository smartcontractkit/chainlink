package bulletprooftxmanager

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
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
	exchainutils "github.com/okex/exchain-ethereum-compatible/utils"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
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
	EvmFinalityDepth() uint
	EvmGasBumpPercent() uint16
	EvmGasBumpThreshold() uint64
	EvmGasBumpTxDepth() uint16
	EvmGasBumpWei() *big.Int
	EvmGasLimitDefault() uint64
	EvmGasLimitMultiplier() float32
	EvmGasPriceDefault() *big.Int
	EvmMaxGasPriceWei() *big.Int
	EvmMaxInFlightTransactions() uint32
	EvmMaxQueuedTransactions() uint64
	EvmMinGasPriceWei() *big.Int
	EvmNonceAutoSync() bool
	EvmRPCDefaultBatchSize() uint32
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EthTxResendAfterThreshold() time.Duration
	GasEstimatorMode() string
	TriggerFallbackDBPollInterval() time.Duration
}

// KeyStore encompasses the subset of keystore used by bulletprooftxmanager
type KeyStore interface {
	GetAll() (keys []ethkey.KeyV2, err error)
	SignTx(fromAddress common.Address, tx *gethTypes.Transaction, chainID *big.Int) (*gethTypes.Transaction, error)
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())
	GetState(id string) (ethkey.State, error)
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
	CreateEthTransaction(db *gorm.DB, newTx NewTx) (etx EthTx, err error)
	GetGasEstimator() gas.Estimator
	RegisterResumeCallback(fn func(id uuid.UUID, value interface{}) error)
}

type BulletproofTxManager struct {
	utils.StartStopOnce

	logger           *logger.Logger
	db               *gorm.DB
	ethClient        eth.Client
	config           Config
	keyStore         KeyStore
	advisoryLocker   postgres.AdvisoryLocker
	eventBroadcaster postgres.EventBroadcaster
	gasEstimator     gas.Estimator

	chHeads        chan models.Head
	trigger        chan common.Address
	resumeCallback func(id uuid.UUID, value interface{}) error

	chStop   chan struct{}
	chSubbed chan struct{}
	wg       sync.WaitGroup

	reaper      *Reaper
	ethResender *EthResender
}

func (b *BulletproofTxManager) RegisterResumeCallback(fn func(id uuid.UUID, value interface{}) error) {
	b.resumeCallback = fn
}

func NewBulletproofTxManager(db *gorm.DB, ethClient eth.Client, config Config, keyStore KeyStore,
	advisoryLocker postgres.AdvisoryLocker, eventBroadcaster postgres.EventBroadcaster, logger *logger.Logger) *BulletproofTxManager {
	b := BulletproofTxManager{
		StartStopOnce:    utils.StartStopOnce{},
		logger:           logger,
		db:               db,
		ethClient:        ethClient,
		config:           config,
		keyStore:         keyStore,
		advisoryLocker:   advisoryLocker,
		eventBroadcaster: eventBroadcaster,
		chHeads:          make(chan models.Head),
		trigger:          make(chan common.Address),
		chStop:           make(chan struct{}),
		chSubbed:         make(chan struct{}),
	}
	if config.EthTxResendAfterThreshold() > 0 {
		b.ethResender = NewEthResender(db, ethClient, defaultResenderPollInterval, config)
	} else {
		b.logger.Info("EthResender: Disabled")
	}
	if config.EthTxReaperThreshold() > 0 {
		b.reaper = NewReaper(db, config)
	} else {
		b.logger.Info("EthTxReaper: Disabled")
	}
	b.gasEstimator = gas.NewEstimator(ethClient, config)

	return &b
}

func (b *BulletproofTxManager) Start() (merr error) {
	return b.StartOnce("BulletproofTxManager", func() error {
		keys, err := b.keyStore.GetAll()
		if err != nil {
			return errors.Wrap(err, "BulletproofTxManager: failed to load keys")
		}

		b.logger.Debugw("BulletproofTxManager: booting", "keys", keys)

		eb := NewEthBroadcaster(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, b.eventBroadcaster, keys, b.gasEstimator, b.logger)
		ec := NewEthConfirmer(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, keys, b.gasEstimator, b.resumeCallback, b.logger)
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
		<-b.chSubbed

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

	close(b.chSubbed)

	for {
		select {
		case address := <-b.trigger:
			eb.Trigger(address)
		case head := <-b.chHeads:
			ec.mb.Deliver(head)
		case <-b.chStop:
			b.logger.ErrorIfCalling(eb.Close)
			b.logger.ErrorIfCalling(ec.Close)
			return
		case <-keysChanged:
			keys, err := b.keyStore.GetAll()
			if err != nil {
				b.logger.Fatalf("BulletproofTxManager: expected keystore to be unlocked: %s", err.Error())
			}

			b.logger.Debugw("BulletproofTxManager: keys changed, reloading", "keys", keys)

			b.logger.ErrorIfCalling(eb.Close)
			b.logger.ErrorIfCalling(ec.Close)

			eb = NewEthBroadcaster(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, b.eventBroadcaster, keys, b.gasEstimator, b.logger)
			ec = NewEthConfirmer(b.db, b.ethClient, b.config, b.keyStore, b.advisoryLocker, keys, b.gasEstimator, b.resumeCallback, b.logger)

			b.logger.ErrorIfCalling(eb.Start)
			b.logger.ErrorIfCalling(ec.Start)
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
			b.logger.Errorw("BulletproofTxManager: timed out handling head", "blockNum", head.Number, "ctxErr", ctx.Err())
		}
	})
	if !ok {
		b.logger.Debugw("BulletproofTxManager: not started; ignoring head", "head", head, "state", b.State())
	}
}

// Trigger forces the EthBroadcaster to check early for the given address
func (b *BulletproofTxManager) Trigger(addr common.Address) {
	select {
	case b.trigger <- addr:
	default:
	}
}

type NewTx struct {
	FromAddress    common.Address
	ToAddress      common.Address
	EncodedPayload []byte
	GasLimit       uint64
	Meta           *EthTxMeta

	MinConfirmations  null.Uint32
	PipelineTaskRunID *uuid.UUID

	Strategy TxStrategy
}

// CreateEthTransaction inserts a new transaction
func (b *BulletproofTxManager) CreateEthTransaction(db *gorm.DB, newTx NewTx) (etx EthTx, err error) {
	err = CheckEthTxQueueCapacity(db, newTx.FromAddress, b.config.EvmMaxQueuedTransactions())
	if err != nil {
		return etx, errors.Wrap(err, "BulletproofTxManager#CreateEthTransaction")
	}

	value := 0
	err = postgres.GormTransactionWithDefaultContext(db, func(tx *gorm.DB) error {
		if newTx.PipelineTaskRunID != nil {
			err = tx.Raw(`SELECT * FROM eth_txes WHERE pipeline_task_run_id = ?`, newTx.PipelineTaskRunID).Scan(&etx).Error
			if err != nil {
				return errors.Wrap(err, "BulletproofTxManager#CreateEthTransaction")
			}

			// if a previous transaction for this task run exists, immediately return it
			if etx.ID != 0 {
				return nil
			}
		}
		res := tx.Raw(`
INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, min_confirmations, pipeline_task_run_id)
VALUES (
?,?,?,?,?,'unstarted',NOW(),?,?,?,?
)
RETURNING "eth_txes".*
`, newTx.FromAddress, newTx.ToAddress, newTx.EncodedPayload, value, newTx.GasLimit, newTx.Meta, newTx.Strategy.Subject(), newTx.MinConfirmations, newTx.PipelineTaskRunID).Scan(&etx)
		err = res.Error
		if err != nil {
			return errors.Wrap(err, "BulletproofTxManager#CreateEthTransaction failed to insert eth_tx")
		}

		pruned, err := newTx.Strategy.PruneQueue(tx)
		if err != nil {
			return errors.Wrap(err, "BulletproofTxManager#CreateEthTransaction failed to prune eth_txes")
		}
		if pruned > 0 {
			b.logger.Warnw(fmt.Sprintf("BulletproofTxManager: dropped %d old transactions from transaction queue", pruned), "fromAddress", newTx.FromAddress, "toAddress", newTx.ToAddress, "meta", newTx.Meta, "subject", newTx.Strategy.Subject(), "replacementID", etx.ID)
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
	hash, signedTxBytes, err := SignTx(ks, etx.FromAddress, transaction, chainID)
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

func SignTx(keyStore KeyStore, address common.Address, tx *gethTypes.Transaction, chainID *big.Int) (common.Hash, []byte, error) {
	signedTx, err := keyStore.SignTx(address, tx, chainID)
	if err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "SignTx failed")
	}
	rlp := new(bytes.Buffer)
	if err = signedTx.EncodeRLP(rlp); err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "SignTx failed")
	}
	var hash common.Hash
	hash, err = signedTxHash(signedTx, chainID)
	if err != nil {
		return hash, nil, err
	}
	return hash, rlp.Bytes(), nil
}

func signedTxHash(signedTx *gethTypes.Transaction, chainID *big.Int) (hash common.Hash, err error) {
	if evmtypes.IsExChain(chainID) {
		hash, err = exchainutils.Hash(signedTx)
		if err != nil {
			return hash, errors.Wrapf(err, "error getting signed tx hash from exchain (chain ID %s)", chainID.String())
		}
	} else {
		hash = signedTx.Hash()
	}
	return hash, nil
}

// send broadcasts the transaction to the ethereum network, writes any relevant
// data onto the attempt and returns an error (or nil) depending on the status
func sendTransaction(ctx context.Context, ethClient eth.Client, a EthTxAttempt, e EthTx, logger *logger.Logger) *eth.SendError {
	signedTx, err := a.GetSignedTx()
	if err != nil {
		return eth.NewFatalSendError(err)
	}

	ctx, cancel := eth.DefaultQueryCtx(ctx)
	defer cancel()

	err = ethClient.SendTransaction(ctx, signedTx)
	err = errors.WithStack(err)

	logger.Debugw("BulletproofTxManager: Sent transaction", "ethTxAttemptID", a.ID, "txHash", a.Hash, "gasPriceWei", a.GasPrice.ToInt().Int64(), "err", err, "meta", e.Meta, "gasLimit", e.GasLimit)
	sendErr := eth.NewSendError(err)
	if sendErr.IsTransactionAlreadyInMempool() {
		logger.Debugw("transaction already in mempool", "txHash", a.Hash, "nodeErr", sendErr.Error())
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
		err = errors.Errorf("cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, static.EvmMaxQueuedTransactionsLabel)
	}
	return
}

var _ TxManager = &NullTxManager{}

type NullTxManager struct {
	ErrMsg string
}

func (n *NullTxManager) OnNewLongestChain(context.Context, models.Head) {}
func (n *NullTxManager) Start() error                                   { return errors.New(n.ErrMsg) }
func (n *NullTxManager) Close() error                                   { return errors.New(n.ErrMsg) }
func (n *NullTxManager) Trigger(common.Address)                         { panic(n.ErrMsg) }
func (n *NullTxManager) CreateEthTransaction(*gorm.DB, NewTx) (etx EthTx, err error) {
	return etx, errors.New(n.ErrMsg)
}
func (n *NullTxManager) Healthy() error                                                        { return nil }
func (n *NullTxManager) Ready() error                                                          { return nil }
func (n *NullTxManager) GetGasEstimator() gas.Estimator                                        { return nil }
func (n *NullTxManager) RegisterResumeCallback(fn func(id uuid.UUID, value interface{}) error) {}
