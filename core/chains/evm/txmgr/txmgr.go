package txmgr

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/label"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Config encompasses config used by txmgr package
// Unless otherwise specified, these should support changing at runtime
//go:generate mockery --recursive --name Config --output ./mocks/ --case=underscore --structname Config --filename config.go
type Config interface {
	gas.Config
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EthTxResendAfterThreshold() time.Duration
	EvmGasBumpThreshold() uint64
	EvmGasBumpTxDepth() uint16
	EvmGasLimitDefault() uint64
	EvmMaxInFlightTransactions() uint32
	EvmMaxQueuedTransactions() uint64
	EvmNonceAutoSync() bool
	EvmUseForwarders() bool
	EvmRPCDefaultBatchSize() uint32
	KeySpecificMaxGasPriceWei(addr common.Address) *big.Int
	TriggerFallbackDBPollInterval() time.Duration
	LogSQL() bool
}

// KeyStore encompasses the subset of keystore used by txmgr
type KeyStore interface {
	GetStatesForChain(chainID *big.Int) ([]ethkey.State, error)
	SignTx(fromAddress common.Address, tx *gethTypes.Transaction, chainID *big.Int) (*gethTypes.Transaction, error)
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())
}

// For more information about the Txm architecture, see the design doc:
// https://www.notion.so/chainlink/Txm-Architecture-Overview-9dc62450cd7a443ba9e7dceffa1a8d6b

var _ TxManager = &Txm{}

// ResumeCallback is assumed to be idempotent
type ResumeCallback func(id uuid.UUID, result interface{}, err error) error

//go:generate mockery --recursive --name TxManager --output ./mocks/ --case=underscore --structname TxManager --filename tx_manager.go
type TxManager interface {
	httypes.HeadTrackable
	services.ServiceCtx
	Trigger(addr common.Address)
	CreateEthTransaction(newTx NewTx, qopts ...pg.QOpt) (etx EthTx, err error)
	GetGasEstimator() gas.Estimator
	RegisterResumeCallback(fn ResumeCallback)
	SendEther(chainID *big.Int, from, to common.Address, value assets.Eth, gasLimit uint64) (etx EthTx, err error)
}

type Txm struct {
	utils.StartStopOnce

	logger           logger.Logger
	db               *sqlx.DB
	q                pg.Q
	ethClient        evmclient.Client
	config           Config
	keyStore         KeyStore
	eventBroadcaster pg.EventBroadcaster
	gasEstimator     gas.Estimator
	chainID          big.Int
	checkerFactory   TransmitCheckerFactory

	chHeads        chan *evmtypes.Head
	trigger        chan common.Address
	resumeCallback ResumeCallback

	chStop   chan struct{}
	chSubbed chan struct{}
	wg       sync.WaitGroup

	reaper      *Reaper
	ethResender *EthResender
	fwdMgr      *forwarders.FwdMgr
}

func (b *Txm) RegisterResumeCallback(fn ResumeCallback) {
	b.resumeCallback = fn
}

// NewTxm creates a new Txm with the given configuration.
func NewTxm(db *sqlx.DB, ethClient evmclient.Client, cfg Config, keyStore KeyStore, eventBroadcaster pg.EventBroadcaster, lggr logger.Logger, checkerFactory TransmitCheckerFactory, logPoller *logpoller.LogPoller) *Txm {
	lggr = lggr.Named("Txm")
	lggr.Infow("Initializing EVM transaction manager",
		"gasBumpTxDepth", cfg.EvmGasBumpTxDepth(),
		"maxInFlightTransactions", cfg.EvmMaxInFlightTransactions(),
		"maxQueuedTransactions", cfg.EvmMaxQueuedTransactions(),
		"nonceAutoSync", cfg.EvmNonceAutoSync(),
		"gasLimitDefault", cfg.EvmGasLimitDefault(),
	)
	b := Txm{
		StartStopOnce:    utils.StartStopOnce{},
		logger:           lggr,
		db:               db,
		q:                pg.NewQ(db, lggr, cfg),
		ethClient:        ethClient,
		config:           cfg,
		keyStore:         keyStore,
		eventBroadcaster: eventBroadcaster,
		gasEstimator:     gas.NewEstimator(lggr, ethClient, cfg),
		chainID:          *ethClient.ChainID(),
		checkerFactory:   checkerFactory,
		chHeads:          make(chan *evmtypes.Head),
		trigger:          make(chan common.Address),
		chStop:           make(chan struct{}),
		chSubbed:         make(chan struct{}),
	}
	if cfg.EthTxResendAfterThreshold() > 0 {
		b.ethResender = NewEthResender(lggr, db, ethClient, defaultResenderPollInterval, cfg)
	} else {
		b.logger.Info("EthResender: Disabled")
	}
	if cfg.EthTxReaperThreshold() > 0 && cfg.EthTxReaperInterval() > 0 {
		b.reaper = NewReaper(lggr, db, cfg, *ethClient.ChainID())
	} else {
		b.logger.Info("EthTxReaper: Disabled")
	}
	if cfg.EvmUseForwarders() {
		b.fwdMgr = forwarders.NewFwdMgr(db, ethClient, logPoller, lggr, cfg)
	} else {
		b.logger.Info("EvmForwardManager: Disabled")
	}

	return &b
}

// Start starts Txm service.
// The provided context can be used to terminate Start sequence.
func (b *Txm) Start(ctx context.Context) (merr error) {
	return b.StartOnce("Txm", func() error {
		keyStates, err := b.keyStore.GetStatesForChain(&b.chainID)
		if err != nil {
			return errors.Wrap(err, "Txm: failed to load key states")
		}

		if len(keyStates) > 0 {
			b.logger.Debugw(fmt.Sprintf("Booting with %d keys", len(keyStates)), "keys", keyStates)
		} else {
			b.logger.Warnf("Chain %s does not have any eth keys, no transactions will be sent on this chain", b.chainID.String())
		}

		eb := NewEthBroadcaster(b.db, b.ethClient, b.config, b.keyStore, b.eventBroadcaster, keyStates, b.gasEstimator, b.resumeCallback, b.logger, b.checkerFactory)
		ec := NewEthConfirmer(b.db, b.ethClient, b.config, b.keyStore, keyStates, b.gasEstimator, b.resumeCallback, b.logger)
		if err := eb.Start(ctx); err != nil {
			return errors.Wrap(err, "Txm: EthBroadcaster failed to start")
		}
		if err := ec.Start(); err != nil {
			return errors.Wrap(err, "Txm: EthConfirmer failed to start")
		}

		if err := b.gasEstimator.Start(ctx); err != nil {
			return errors.Wrap(err, "Txm: Estimator failed to start")
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

		if b.fwdMgr != nil {
			if err = b.fwdMgr.Start(); err != nil {
				return errors.Wrap(err, "Txm: EVMForwarderManager failed to start")
			}
		}

		return nil
	})
}

func (b *Txm) Close() (merr error) {
	return b.StopOnce("Txm", func() error {
		close(b.chStop)

		if b.reaper != nil {
			b.reaper.Stop()
		}
		if b.ethResender != nil {
			b.ethResender.Stop()
		}
		if b.fwdMgr != nil {
			if err := b.fwdMgr.Stop(); err != nil {
				return errors.Wrap(err, "Txm: failed to stop EVMForwarderManager")
			}
		}

		b.wg.Wait()

		b.gasEstimator.Close()

		return nil
	})
}

func (b *Txm) runLoop(eb *EthBroadcaster, ec *EthConfirmer) {
	defer b.wg.Done()
	keysChanged, unsub := b.keyStore.SubscribeToKeyChanges()
	defer unsub()

	close(b.chSubbed)

	ctx, cancel := utils.ContextFromChan(b.chStop)
	defer cancel()

	for {
		select {
		case address := <-b.trigger:
			eb.Trigger(address)
		case head := <-b.chHeads:
			ec.mb.Deliver(head)
		case <-b.chStop:
			b.logger.ErrorIfClosing(eb, "EthBroadcaster")
			b.logger.ErrorIfClosing(ec, "EthConfirmer")
			return
		case <-keysChanged:
			keyStates, err := b.keyStore.GetStatesForChain(&b.chainID)
			if err != nil {
				b.logger.Errorf("Failed to reload key states after key change")
				continue
			}
			b.logger.Debugw("Keys changed, reloading", "keyStates", keyStates)

			b.logger.ErrorIfClosing(eb, "EthBroadcaster")
			b.logger.ErrorIfClosing(ec, "EthConfirmer")

			eb = NewEthBroadcaster(b.db, b.ethClient, b.config, b.keyStore, b.eventBroadcaster, keyStates, b.gasEstimator, b.resumeCallback, b.logger, b.checkerFactory)
			ec = NewEthConfirmer(b.db, b.ethClient, b.config, b.keyStore, keyStates, b.gasEstimator, b.resumeCallback, b.logger)

			if err := eb.Start(ctx); err != nil {
				b.logger.Criticalw("Failed to start EthBroadcaster", "error", err)
			}
			if err := ec.Start(); err != nil {
				b.logger.Criticalw("Failed to start EthConfirmer", "error", err)
			}
		}
	}
}

// OnNewLongestChain conforms to HeadTrackable
func (b *Txm) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	ok := b.IfStarted(func() {
		if b.reaper != nil {
			b.reaper.SetLatestBlockNum(head.Number)
		}
		b.gasEstimator.OnNewLongestChain(ctx, head)
		select {
		case b.chHeads <- head:
		case <-ctx.Done():
			b.logger.Errorw("Timed out handling head", "blockNum", head.Number, "ctxErr", ctx.Err())
		}
	})
	if !ok {
		b.logger.Debugw("Not started; ignoring head", "head", head, "state", b.State())
	}
}

// Trigger forces the EthBroadcaster to check early for the given address
func (b *Txm) Trigger(addr common.Address) {
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

	// Checker defines the check that should be run before a transaction is submitted on chain.
	Checker TransmitCheckerSpec
}

// CreateEthTransaction inserts a new transaction
func (b *Txm) CreateEthTransaction(newTx NewTx, qs ...pg.QOpt) (etx EthTx, err error) {
	q := b.q.WithOpts(qs...)
	if b.config.EvmUseForwarders() {
		fwdAddr, fwdPayload, fwdErr := b.fwdMgr.MaybeForwardTransaction(newTx.FromAddress, newTx.ToAddress, newTx.EncodedPayload)
		if fwdErr == nil {
			newTx.ToAddress = fwdAddr
			newTx.EncodedPayload = fwdPayload
		} else {
			b.logger.Infof("Skipping using forwarders: %s", fwdErr.Error())
		}
	}

	err = CheckEthTxQueueCapacity(q, newTx.FromAddress, b.config.EvmMaxQueuedTransactions(), b.chainID)
	if err != nil {
		return etx, errors.Wrap(err, "Txm#CreateEthTransaction")
	}

	value := 0
	err = q.Transaction(func(tx pg.Queryer) error {
		if newTx.PipelineTaskRunID != nil {
			err = tx.Get(&etx, `SELECT * FROM eth_txes WHERE pipeline_task_run_id = $1 AND evm_chain_id = $2`, newTx.PipelineTaskRunID, b.chainID.String())
			// If no eth_tx matches (the common case) then continue
			if !errors.Is(err, sql.ErrNoRows) {
				if err != nil {
					return errors.Wrap(err, "Txm#CreateEthTransaction")
				}
				// if a previous transaction for this task run exists, immediately return it
				return nil
			}
		}
		if err = b.checkStateExists(tx, newTx.FromAddress); err != nil {
			return err
		}
		err := tx.Get(&etx, `
INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, transmit_checker)
VALUES (
$1,$2,$3,$4,$5,'unstarted',NOW(),$6,$7,$8,$9,$10,$11
)
RETURNING "eth_txes".*
`, newTx.FromAddress, newTx.ToAddress, newTx.EncodedPayload, value, newTx.GasLimit, newTx.Meta, newTx.Strategy.Subject(), b.chainID.String(), newTx.MinConfirmations, newTx.PipelineTaskRunID, newTx.Checker)
		if err != nil {
			return errors.Wrap(err, "Txm#CreateEthTransaction failed to insert eth_tx")
		}

		pruned, err := newTx.Strategy.PruneQueue(tx)
		if err != nil {
			return errors.Wrap(err, "Txm#CreateEthTransaction failed to prune eth_txes")
		}
		if pruned > 0 {
			b.logger.Warnw(fmt.Sprintf("Dropped %d old transactions from transaction queue", pruned), "fromAddress", newTx.FromAddress, "toAddress", newTx.ToAddress, "meta", newTx.Meta, "subject", newTx.Strategy.Subject(), "replacementID", etx.ID)
		}
		return nil
	})
	return
}

func (b *Txm) checkStateExists(q pg.Queryer, addr common.Address) error {
	var state ethkey.State
	err := q.Get(&state, `SELECT * FROM eth_key_states WHERE address = $1`, addr)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.Errorf("no eth key exists with address %s", addr.Hex())
	} else if err != nil {
		return errors.Wrap(err, "failed to query state")
	}
	if state.EVMChainID.Cmp(utils.NewBig(&b.chainID)) != 0 {
		return errors.Errorf("cannot send transaction on chain ID %s; eth key with address %s is pegged to chain ID %s", b.chainID.String(), addr.Hex(), state.EVMChainID.String())
	}
	return nil
}

// GetGasEstimator returns the gas estimator, mostly useful for tests
func (b *Txm) GetGasEstimator() gas.Estimator {
	return b.gasEstimator
}

// SendEther creates a transaction that transfers the given value of ether
func (b *Txm) SendEther(chainID *big.Int, from, to common.Address, value assets.Eth, gasLimit uint64) (etx EthTx, err error) {
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
		EVMChainID:     *utils.NewBig(chainID),
	}
	query := `INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, evm_chain_id, created_at) VALUES (
:from_address, :to_address, :encoded_payload, :value, :gas_limit, :state, :evm_chain_id, NOW()
) RETURNING eth_txes.*`
	err = b.q.GetNamed(query, &etx, etx)
	return etx, errors.Wrap(err, "SendEther failed to insert eth_tx")
}

type ChainKeyStore struct {
	chainID  big.Int
	config   Config
	keystore KeyStore
}

func NewChainKeyStore(chainID big.Int, config Config, keystore KeyStore) ChainKeyStore {
	return ChainKeyStore{chainID, config, keystore}
}

func (c *ChainKeyStore) SignTx(address common.Address, tx *gethTypes.Transaction) (common.Hash, []byte, error) {
	signedTx, err := c.keystore.SignTx(address, tx, &c.chainID)
	if err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "SignTx failed")
	}
	rlp := new(bytes.Buffer)
	if err = signedTx.EncodeRLP(rlp); err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "SignTx failed")
	}
	return signedTx.Hash(), rlp.Bytes(), nil
}

// send broadcasts the transaction to the ethereum network, writes any relevant
// data onto the attempt and returns an error (or nil) depending on the status
func sendTransaction(ctx context.Context, ethClient evmclient.Client, a EthTxAttempt, e EthTx, logger logger.Logger) *evmclient.SendError {
	signedTx, err := a.GetSignedTx()
	if err != nil {
		return evmclient.NewFatalSendError(err)
	}

	err = ethClient.SendTransaction(ctx, signedTx)
	err = errors.WithStack(err)

	a.EthTx = e // for logging
	logger.Debugw("Sent transaction", "ethTxAttemptID", a.ID, "txHash", a.Hash, "err", err, "meta", e.Meta, "gasLimit", e.GasLimit, "attempt", a)
	sendErr := evmclient.NewSendError(err)
	if sendErr.IsTransactionAlreadyInMempool() {
		logger.Debugw("Transaction already in mempool", "txHash", a.Hash, "nodeErr", sendErr.Error())
		return nil
	}
	return sendErr
}

// sendEmptyTransaction sends a transaction with 0 Eth and an empty payload to the burn address
// May be useful for clearing stuck nonces
func sendEmptyTransaction(
	ctx context.Context,
	ethClient evmclient.Client,
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

const insertIntoEthTxAttemptsQuery = `
INSERT INTO eth_tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, broadcast_before_block_num, state, created_at, chain_specific_gas_limit, tx_type, gas_tip_cap, gas_fee_cap)
VALUES (:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :broadcast_before_block_num, :state, NOW(), :chain_specific_gas_limit, :tx_type, :gas_tip_cap, :gas_fee_cap)
RETURNING *;
`

func saveReplacementInProgressAttempt(q pg.Q, oldAttempt EthTxAttempt, replacementAttempt *EthTxAttempt) error {
	if oldAttempt.State != EthTxAttemptInProgress || replacementAttempt.State != EthTxAttemptInProgress {
		return errors.New("expected attempts to be in_progress")
	}
	if oldAttempt.ID == 0 {
		return errors.New("expected oldAttempt to have an ID")
	}
	return q.Transaction(func(tx pg.Queryer) error {
		if _, err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE id=$1`, oldAttempt.ID); err != nil {
			return errors.Wrap(err, "saveReplacementInProgressAttempt failed to delete from eth_tx_attempts")
		}
		query, args, e := tx.BindNamed(insertIntoEthTxAttemptsQuery, replacementAttempt)
		if e != nil {
			return errors.Wrap(e, "saveReplacementInProgressAttempt failed to BindNamed")
		}
		return errors.Wrap(tx.Get(replacementAttempt, query, args...), "saveReplacementInProgressAttempt failed to insert replacement attempt")
	})
}

// CountUnconfirmedTransactions returns the number of unconfirmed transactions
func CountUnconfirmedTransactions(q pg.Q, fromAddress common.Address, chainID big.Int) (count uint32, err error) {
	return countTransactionsWithState(q, fromAddress, EthTxUnconfirmed, chainID)
}

// CountUnstartedTransactions returns the number of unconfirmed transactions
func CountUnstartedTransactions(q pg.Q, fromAddress common.Address, chainID big.Int) (count uint32, err error) {
	return countTransactionsWithState(q, fromAddress, EthTxUnstarted, chainID)
}

func countTransactionsWithState(q pg.Q, fromAddress common.Address, state EthTxState, chainID big.Int) (count uint32, err error) {
	err = q.Get(&count, `SELECT count(*) FROM eth_txes WHERE from_address = $1 AND state = $2 AND evm_chain_id = $3`,
		fromAddress, state, chainID.String())
	return count, errors.Wrap(err, "failed to countTransactionsWithState")
}

// CheckEthTxQueueCapacity returns an error if inserting this transaction would
// exceed the maximum queue size.
func CheckEthTxQueueCapacity(q pg.Queryer, fromAddress common.Address, maxQueuedTransactions uint64, chainID big.Int) (err error) {
	if maxQueuedTransactions == 0 {
		return nil
	}
	var count uint64
	err = q.Get(&count, `SELECT count(*) FROM eth_txes WHERE from_address = $1 AND state = 'unstarted' AND evm_chain_id = $2`, fromAddress, chainID.String())
	if err != nil {
		err = errors.Wrap(err, "txmgr.CheckEthTxQueueCapacity query failed")
		return
	}

	if count >= maxQueuedTransactions {
		err = errors.Errorf("cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, label.MaxQueuedTransactionsWarning)
	}
	return
}

var _ TxManager = &NullTxManager{}

type NullTxManager struct {
	ErrMsg string
}

func (n *NullTxManager) OnNewLongestChain(context.Context, *evmtypes.Head) {}

// Start does noop for NullTxManager.
func (n *NullTxManager) Start(context.Context) error { return nil }

// Close does noop for NullTxManager.
func (n *NullTxManager) Close() error { return nil }

// Trigger does noop for NullTxManager.
func (n *NullTxManager) Trigger(common.Address) { panic(n.ErrMsg) }
func (n *NullTxManager) CreateEthTransaction(NewTx, ...pg.QOpt) (etx EthTx, err error) {
	return etx, errors.New(n.ErrMsg)
}

// SendEther does nothing, null functionality
func (n *NullTxManager) SendEther(chainID *big.Int, from, to common.Address, value assets.Eth, gasLimit uint64) (etx EthTx, err error) {
	return etx, errors.New(n.ErrMsg)
}
func (n *NullTxManager) Healthy() error                           { return nil }
func (n *NullTxManager) Ready() error                             { return nil }
func (n *NullTxManager) GetGasEstimator() gas.Estimator           { return nil }
func (n *NullTxManager) RegisterResumeCallback(fn ResumeCallback) {}
