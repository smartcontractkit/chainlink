package txm

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

const (
	broadcastInterval       time.Duration = 15 * time.Second
	maxInFlightTransactions uint64        = 16
	maxAllowedAttempts      uint16        = 10
)

type Client interface {
	PendingNonceAt(context.Context, common.Address) (uint64, error)
	NonceAt(context.Context, common.Address, *big.Int) (uint64, error)
	SendTransaction(context.Context, *evmtypes.Transaction) error
}

type TxStore interface {
	AbandonPendingTransactions(context.Context, common.Address) error
	AppendAttemptToTransaction(context.Context, uint64, *types.Attempt) error
	CreateEmptyUnconfirmedTransaction(context.Context, common.Address, *big.Int, uint64, uint64) (*types.Transaction, error)
	CreateTransaction(context.Context, *types.TxRequest) (*types.Transaction, error)
	FetchUnconfirmedTransactionAtNonceWithCount(context.Context, uint64, common.Address) (*types.Transaction, int, error)
	MarkTransactionsConfirmed(context.Context, uint64, common.Address) ([]uint64, []uint64, error)
	MarkUnconfirmedTransactionPurgeable(context.Context, uint64) error
	UpdateTransactionBroadcast(context.Context, uint64, uint64, common.Hash) error
	UpdateUnstartedTransactionWithNonce(context.Context, common.Address, uint64) (*types.Transaction, error)

	// ErrorHandler
	DeleteAttemptForUnconfirmedTx(context.Context, uint64, *types.Attempt) error
	MarkTxFatal(context.Context, *types.Transaction) error
}

type AttemptBuilder interface {
	NewAttempt(context.Context, logger.Logger, *types.Transaction, bool) (*types.Attempt, error)
	NewBumpAttempt(context.Context, logger.Logger, *types.Transaction, types.Attempt) (*types.Attempt, error)
}

type ErrorHandler interface {
	HandleError(tx *types.Transaction, message error, attemptBuilder AttemptBuilder, client Client, txStore TxStore) (err error)
}

type StuckTxDetector interface {
	DetectStuckTransactions(tx *types.Transaction) (bool, error)
}

type Config struct {
	EIP1559             bool
	BlockTime           time.Duration
	RetryBlockThreshold uint16
	EmptyTxLimitDefault uint64
}

type Txm struct {
	services.StateMachine
	lggr            logger.SugaredLogger
	address         common.Address
	chainID         *big.Int
	client          Client
	attemptBuilder  AttemptBuilder
	errorHandler    ErrorHandler
	stuckTxDetector StuckTxDetector
	txStore         TxStore
	config          Config
	nonce           atomic.Uint64

	triggerCh       chan struct{}
	broadcastStopCh services.StopChan
	backfillStopCh  services.StopChan
	wg              sync.WaitGroup
}

func NewTxm(lggr logger.Logger, chainID *big.Int, client Client, attemptBuilder AttemptBuilder, txStore TxStore, config Config, address common.Address) *Txm {
	return &Txm{
		lggr:            logger.Sugared(logger.Named(lggr, "Txm")),
		address:         address,
		chainID:         chainID,
		client:          client,
		attemptBuilder:  attemptBuilder,
		txStore:         txStore,
		config:          config,
		triggerCh:       make(chan struct{}),
		broadcastStopCh: make(chan struct{}),
		backfillStopCh:  make(chan struct{}),
	}
}

func (t *Txm) Start(context.Context) error {
	return t.StartOnce("Txm", func() error {
		pendingNonce, err := t.client.PendingNonceAt(context.TODO(), t.address)
		if err != nil {
			return err
		}
		t.nonce.Store(pendingNonce)
		t.wg.Add(2)
		go t.broadcastLoop()
		go t.backfillLoop()

		return nil
	})
}

func (t *Txm) Close() error {
	return t.StopOnce("Txm", func() error {
		close(t.broadcastStopCh)
		close(t.backfillStopCh)
		t.wg.Wait()
		return nil
	})
}

func (t *Txm) CreateTransaction(ctx context.Context, txRequest *types.TxRequest) (tx *types.Transaction, err error) {
	tx, err = t.txStore.CreateTransaction(ctx, txRequest)
	if err == nil {
		t.lggr.Infow("Created transaction", "tx", tx)
	}
	return
}

func (t *Txm) Trigger() error {
	if !t.IfStarted(func() {
		t.triggerCh <- struct{}{}
	}) {
		return fmt.Errorf("Txm unstarted")
	}
	return nil
}

func (t *Txm) Abandon() error {
	return t.txStore.AbandonPendingTransactions(context.TODO(), t.address)
}

func (t *Txm) broadcastLoop() {
	defer t.wg.Done()
	broadcasterTicker := time.NewTicker(utils.WithJitter(broadcastInterval))
	defer broadcasterTicker.Stop()

	for {
		select {
		case <-t.broadcastStopCh:
			return
		case <-t.triggerCh:
			start := time.Now()
			if err := t.broadcastTransaction(); err != nil {
				t.lggr.Errorf("Error during triggered transaction broadcasting: %v", err)
			} else {
				t.lggr.Debug("Triggered transaction broadcasting time elapsed: ", time.Since(start))
			}
			broadcasterTicker.Reset(utils.WithJitter(broadcastInterval))
		case <-broadcasterTicker.C:
			start := time.Now()
			if err := t.broadcastTransaction(); err != nil {
				t.lggr.Errorf("Error during transaction broadcasting: %v", err)
			} else {
				t.lggr.Debug("Transaction broadcasting time elapsed: ", time.Since(start))
			}
		}
	}
}

func (t *Txm) backfillLoop() {
	defer t.wg.Done()
	backfillTicker := time.NewTicker(utils.WithJitter(t.config.BlockTime))
	defer backfillTicker.Stop()

	for {
		select {
		case <-t.backfillStopCh:
			return
		case <-backfillTicker.C:
			start := time.Now()
			if err := t.backfillTransactions(); err != nil {
				t.lggr.Errorf("Error during backfill: %v", err)
			} else {
				t.lggr.Debug("Backfill time elapsed: ", time.Since(start))
			}
		}
	}
}

func (t *Txm) broadcastTransaction() error {
	_, unconfirmedCount, err := t.txStore.FetchUnconfirmedTransactionAtNonceWithCount(context.TODO(), 0, t.address)
	if err != nil {
		return err
	}

	// Optimistically send up to 1/3 of the maxInFlightTransactions. After that threshold, broadcast more cautiously
	// by checking the pending nonce so no more than maxInFlightTransactions/3 can get stuck simultaneously i.e. due
	// to insufficient balance. We're making this trade-off to avoid storing stuck transactions and making unnecessary
	// RPC calls. The upper limit is always maxInFlightTransactions regardless of the pending nonce.
	if unconfirmedCount >= int(maxInFlightTransactions)/3 {
		if unconfirmedCount > int(maxInFlightTransactions) {
			t.lggr.Warnf("Reached transaction limit: %d for unconfirmed transactions", maxInFlightTransactions)
			return nil
		}
		pendingNonce, err := t.client.PendingNonceAt(context.TODO(), t.address)
		if err != nil {
			return err
		}
		if t.nonce.Load() > pendingNonce {
			t.lggr.Warnf("Reached transaction limit. LocalNonce: %d, PendingNonce %d, unconfirmedCount: %d",
				t.nonce.Load(), pendingNonce, unconfirmedCount)
				return nil
		}
	}

	tx, err := t.txStore.UpdateUnstartedTransactionWithNonce(context.TODO(), t.address, t.nonce.Load())
	if err != nil {
		return err
	}
	if tx == nil {
		return err
	}
	tx.Nonce = t.nonce.Load()
	tx.State = types.TxUnconfirmed
	t.nonce.Add(1)

	return t.createAndSendAttempt(tx)
}

func (t *Txm) createAndSendAttempt(tx *types.Transaction) error {
	attempt, err := t.attemptBuilder.NewAttempt(context.TODO(), t.lggr, tx, t.config.EIP1559)
	if err != nil {
		return err
	}

	if err = t.txStore.AppendAttemptToTransaction(context.TODO(), tx.Nonce, attempt); err != nil {
		return err
	}

	return t.sendTransactionWithError(tx, attempt)
}

func (t *Txm) sendTransactionWithError(tx *types.Transaction, attempt *types.Attempt) (err error) {
	txErr := t.client.SendTransaction(context.TODO(), attempt.SignedTransaction)
	tx.AttemptCount++
	t.lggr.Infow("Broadcasted attempt", "tx", tx, "attempt", attempt, "txErr: ", txErr)
	if txErr != nil && t.errorHandler != nil {
		if err = t.errorHandler.HandleError(tx, txErr, t.attemptBuilder, t.client, t.txStore); err != nil {
			return
		}
	} else if txErr != nil {
		pendingNonce, err := t.client.PendingNonceAt(context.TODO(), t.address)
		if err != nil {
			return err
		}
		if pendingNonce <= tx.Nonce {
			t.lggr.Debugf("Pending nonce for txID: %v didn't increase. PendingNonce: %d, TxNonce: %d", tx.ID, pendingNonce, tx.Nonce)
			return nil
		}
	}

	return t.txStore.UpdateTransactionBroadcast(context.TODO(), attempt.TxID, tx.Nonce, attempt.Hash)
}

func (t *Txm) backfillTransactions() error {
	latestNonce, err := t.client.NonceAt(context.TODO(), t.address, nil)
	if err != nil {
		return err
	}

	confirmedTransactionIDs, unconfirmedTransactionIDs, err := t.txStore.MarkTransactionsConfirmed(context.TODO(), latestNonce, t.address)
	if err != nil {
		return err
	}
	if len(confirmedTransactionIDs) > 0 || len(unconfirmedTransactionIDs) > 0 {
		t.lggr.Infof("Confirmed transaction IDs: %v . Re-orged transaction IDs: %v", confirmedTransactionIDs, unconfirmedTransactionIDs)
	}

	tx, unconfirmedCount, err := t.txStore.FetchUnconfirmedTransactionAtNonceWithCount(context.TODO(), latestNonce, t.address)
	if err != nil {
		return err
	}
	if unconfirmedCount == 0 {
		t.lggr.Debugf("All transactions confirmed for address: %v", t.address)
		return nil
	}

	if tx == nil || tx.Nonce != latestNonce {
		t.lggr.Warnf("Nonce gap at nonce: %d - address: %v. Creating a new transaction\n", latestNonce, t.address)
		return t.createAndSendEmptyTx(latestNonce)
	} else {
		if !tx.IsPurgeable && t.stuckTxDetector != nil {
			isStuck, err := t.stuckTxDetector.DetectStuckTransactions(tx)
			if err != nil {
				return err
			}
			if isStuck {
				tx.IsPurgeable = true
				t.txStore.MarkUnconfirmedTransactionPurgeable(context.TODO(), tx.Nonce)
				t.lggr.Infof("Marked tx as purgeable. Sending purge attempt for txID: ", tx.ID)
				return t.createAndSendAttempt(tx)
			}
		}

		if tx.AttemptCount >= maxAllowedAttempts {
			return fmt.Errorf("reached max allowed attempts for txID: %d. TXM won't broadcast any more attempts."+
				"If this error persists, it means the transaction won't be confirmed and the TXM needs to be restarted."+
				"Look for any error messages from previous attempts that may indicate why this happened, i.e. wallet is out of funds. Tx: %v", tx.ID, tx)
		}

		if time.Since(tx.LastBroadcastAt) > (t.config.BlockTime*time.Duration(t.config.RetryBlockThreshold)) || tx.LastBroadcastAt.IsZero() {
			// TODO: add optional graceful bumping strategy
			t.lggr.Info("Rebroadcasting attempt for txID: ", tx.ID)
			return t.createAndSendAttempt(tx)
		}
	}
	return nil
}

func (t *Txm) createAndSendEmptyTx(latestNonce uint64) error {
	tx, err := t.txStore.CreateEmptyUnconfirmedTransaction(context.TODO(), t.address, t.chainID, latestNonce, t.config.EmptyTxLimitDefault)
	if err != nil {
		return err
	}
	return t.createAndSendAttempt(tx)
}
