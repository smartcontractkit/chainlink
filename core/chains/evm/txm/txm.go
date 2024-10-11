package txm

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

const (
	broadcastInterval       time.Duration = 30 * time.Second
	maxInFlightTransactions uint64        = 16
	maxAllowedAttempts      uint16        = 10
)

type Client interface {
	PendingNonceAt(context.Context, common.Address) (uint64, error)
	NonceAt(context.Context, common.Address, *big.Int) (uint64, error)
	SendTransaction(context.Context, *evmtypes.Transaction) error
	BatchCallContext(context.Context, []rpc.BatchElem) error
}

type Storage interface {
	AbandonPendingTransactions(context.Context, common.Address) error
	AppendAttemptToTransaction(context.Context, uint64, *types.Attempt) error
	CountUnstartedTransactions(context.Context, common.Address) (int, error)
	CreateEmptyUnconfirmedTransaction(context.Context, common.Address, *big.Int, uint64, uint64) (*types.Transaction, error)
	CreateTransaction(context.Context, *types.Transaction) (uint64, error)
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
	HandleError(tx *types.Transaction, message error, attemptBuilder AttemptBuilder, client Client, storage Storage) (err error)
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
	storage         Storage
	config          Config
	nonce           atomic.Uint64

	triggerCh       chan struct{}
	broadcastStopCh services.StopChan
	backfillStopCh  services.StopChan
	wg              *sync.WaitGroup
}

func NewTxm(lggr logger.Logger, chainID *big.Int, client Client, attemptBuilder AttemptBuilder, storage Storage, config Config, address common.Address) *Txm {
	return &Txm{
		lggr:            logger.Sugared(logger.Named(lggr, "Txm")),
		address:         address,
		chainID:         chainID,
		client:          client,
		attemptBuilder:  attemptBuilder,
		storage:         storage,
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

func (t *Txm) Trigger() error {
	if !t.IfStarted(func() {
		t.triggerCh <- struct{}{}
	}) {
		return fmt.Errorf("Txm unstarted")
	}
	return nil
}

func (t *Txm) Abandon() error {
	return t.storage.AbandonPendingTransactions(context.TODO(), t.address)
}

func (t *Txm) broadcastLoop() {
	defer t.wg.Done()
	broadcasterTicker := time.NewTicker(utils.WithJitter(broadcastInterval) * time.Second)
	defer broadcasterTicker.Stop()

	for {
		select {
		case <-t.broadcastStopCh:
			return
		case <-t.triggerCh:
			start := time.Now()
			if err := t.broadcastTransaction(); err != nil {
				t.lggr.Errorf("Error during triggered transaction broadcasting %w", err)
			} else {
				t.lggr.Debug("Triggered transaction broadcasting time elapsed: ", time.Since(start))
			}
			broadcasterTicker.Reset(utils.WithJitter(broadcastInterval) * time.Second)
		case <-broadcasterTicker.C:
			start := time.Now()
			if err := t.broadcastTransaction(); err != nil {
				t.lggr.Errorf("Error during transaction broadcasting: %w", err)
			} else {
				t.lggr.Debug("Transaction broadcasting time elapsed: ", time.Since(start))
			}
		}
	}
}

func (t *Txm) backfillLoop() {
	defer t.wg.Done()
	backfillTicker := time.NewTicker(utils.WithJitter(t.config.BlockTime) * time.Second)
	defer backfillTicker.Stop()

	for {
		select {
		case <-t.backfillStopCh:
			return
		case <-backfillTicker.C:
			start := time.Now()
			if err := t.backfillTransactions(); err != nil {
				t.lggr.Errorf("Error during backfill: %w", err)
			} else {
				t.lggr.Debug("Backfill time elapsed: ", time.Since(start))
			}
		}
	}
}

func (t *Txm) broadcastTransaction() (err error) {
	pendingNonce, latestNonce, err := t.pendingAndLatestNonce(context.TODO(), t.address)
	if err != nil {
		return
	}

	// Some clients allow out-of-order nonce filling, but it's safer to disable it.
	if pendingNonce-latestNonce > maxInFlightTransactions || t.nonce.Load() > pendingNonce {
		t.lggr.Warnf("Reached transaction limit. LocalNonce: %d, PendingNonce %d, LatestNonce: %d, maxInFlightTransactions: %d",
			t.nonce.Load(), pendingNonce, latestNonce, maxInFlightTransactions)
		return
	}

	tx, err := t.storage.UpdateUnstartedTransactionWithNonce(context.TODO(), t.address, t.nonce.Load())
	if err != nil {
		return
	}
	if tx == nil {
		return
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

	if err = t.storage.AppendAttemptToTransaction(context.TODO(), tx.Nonce, attempt); err != nil {
		return err
	}

	return t.sendTransactionWithError(tx, attempt)
}

func (t *Txm) sendTransactionWithError(tx *types.Transaction, attempt *types.Attempt) (err error) {
	txErr := t.client.SendTransaction(context.TODO(), attempt.SignedTransaction)
	tx.AttemptCount++
	t.lggr.Infof("Broadcasted attempt", "tx", tx, "attempt", attempt, "txErr: ", txErr)
	if txErr != nil && t.errorHandler != nil {
		if err = t.errorHandler.HandleError(tx, txErr, t.attemptBuilder, t.client, t.storage); err != nil {
			return
		}
	} else if txErr != nil {
		pendingNonce, err := t.client.PendingNonceAt(context.TODO(), t.address)
		if err != nil {
			return err
		}
		if pendingNonce > tx.Nonce {
			return nil
		}
		t.lggr.Debugf("Pending nonce for txID: %v didn't increase. PendingNonce: %d, TxNonce: %d", tx.ID, pendingNonce, tx.Nonce)
		return nil
	}

	return t.storage.UpdateTransactionBroadcast(context.TODO(), attempt.TxID, tx.Nonce, attempt.Hash)
}

func (t *Txm) backfillTransactions() error {
	latestNonce, err := t.client.NonceAt(context.TODO(), t.address, nil)
	if err != nil {
		return err
	}

	// TODO: Update LastBroadcast(?)
	confirmedTransactionIDs, unconfirmedTransactionIDs, err := t.storage.MarkTransactionsConfirmed(context.TODO(), latestNonce, t.address)
	if err != nil {
		return err
	}
	t.lggr.Infof("Confirmed transactions: %v . Re-orged transactions: %v", confirmedTransactionIDs, unconfirmedTransactionIDs)

	tx, unconfirmedCount, err := t.storage.FetchUnconfirmedTransactionAtNonceWithCount(context.TODO(), latestNonce, t.address)
	if err != nil {
		return err
	}
	if unconfirmedCount == 0 {
		pendingNonce, err := t.client.PendingNonceAt(context.TODO(), t.address)
		if err != nil {
			return err
		}
		// if local nonce is incorrect, we need to fill the gap to start new transactions
		count, err := t.storage.CountUnstartedTransactions(context.TODO(), t.address)
		if err != nil {
			return err
		}
		if t.nonce.Load() <= pendingNonce && count == 0 {
			t.lggr.Debugf("All transactions confirmed for address: %v", t.address)
			return nil
		}
	}

	if tx == nil || tx.Nonce != latestNonce {
		t.lggr.Warn("Nonce gap at nonce: %d - address: %v. Creating a new transaction\n", latestNonce, t.address)
		return t.createAndSendEmptyTx(latestNonce)
	} else {
		if !tx.IsPurgeable && t.stuckTxDetector != nil {
			isStuck, err := t.stuckTxDetector.DetectStuckTransactions(tx)
			if err != nil {
				return err
			}
			if isStuck {
				tx.IsPurgeable = true
				t.storage.MarkUnconfirmedTransactionPurgeable(context.TODO(), tx.Nonce)
				t.lggr.Infof("Marked tx as purgeable. Sending purge attempt for tx: ", tx.ID, tx)
				return t.createAndSendAttempt(tx)
			}
		}
		if (time.Since(tx.LastBroadcastAt) > (t.config.BlockTime*time.Duration(t.config.RetryBlockThreshold)) || tx.LastBroadcastAt.IsZero()) &&
			tx.AttemptCount < maxAllowedAttempts {
			// TODO: add graceful bumping
			t.lggr.Infow("Rebroadcasting attempt for tx: ", tx)
			return t.createAndSendAttempt(tx)
		}
	}
	return nil
}

func (t *Txm) createAndSendEmptyTx(latestNonce uint64) error {
	tx, err := t.storage.CreateEmptyUnconfirmedTransaction(context.TODO(), t.address, t.chainID, latestNonce, t.config.EmptyTxLimitDefault)
	if err != nil {
		return err
	}
	return t.createAndSendAttempt(tx)
}

func (t *Txm) pendingAndLatestNonce(ctx context.Context, fromAddress common.Address) (pending uint64, latest uint64, err error) {
	pendingS, latestS := new(string), new(string)
	reqs := []rpc.BatchElem{
		{Method: "eth_getTransactionCount", Args: []interface{}{fromAddress, "pending"}, Result: &pendingS},
		{Method: "eth_getTransactionCount", Args: []interface{}{fromAddress, "latest"}, Result: &latestS},
	}

	if err = t.client.BatchCallContext(ctx, reqs); err != nil {
		return
	}

	for _, response := range reqs {
		if response.Error != nil {
			return 0, 0, response.Error
		}
	}

	if pending, err = hexutil.DecodeUint64(*pendingS); err != nil {
		return
	}
	if latest, err = hexutil.DecodeUint64(*latestS); err != nil {
		return
	}

	if pending < latest {
		return 0, 0, fmt.Errorf("RPC nonce state out of sync. Pending: %d, Latest: %d", pending, latest)
	}

	return pending, latest, err
}
