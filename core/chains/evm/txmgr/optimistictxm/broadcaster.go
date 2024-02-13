package txm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jpillora/backoff"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/chains/label"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// InFlightTransactionRecheckInterval controls how often the Broadcaster
// will poll the unconfirmed queue to see if it is allowed to send another
// transaction
const InFlightTransactionRecheckInterval = 1 * time.Second

var ErrTxRemoved = errors.New("tx removed")

type TxAttemptBuilder interface {
	NewAttempt(context.Context, txmgr.Tx, logger.Logger) (txmgr.TxAttempt, error)
}

type BroadcasterTxStore interface {
	CountUnconfirmedTransactions(context.Context, common.Address, *big.Int) (uint32, error)
	CountUnstartedTransactions(context.Context, common.Address, *big.Int) (uint32, error)
	FindNextUnstartedTransactionFromAddress(context.Context, *txmgr.Tx, common.Address, *big.Int) error
	BroadcasterUpdateTxUnstartedToInProgress(context.Context, *txmgr.Tx) error
	BroadcasterGetTxInProgress(context.Context, common.Address) (*txmgr.Tx, error)
	UpdateTxInProgressToUnconfirmed(context.Context, *txmgr.Tx) error
}

type BroadcasterClient interface {
	ConfiguredChainID() *big.Int
	PendingNonceAt(context.Context, common.Address) (uint64, error)
	SendTransaction(context.Context, *types.Transaction) error
}

type BroadcasterConfig struct {
	FallbackPollInterval time.Duration
	MaxInFlight          uint32
	NonceAutoSync        bool
}

type KeyStore interface {
	EnabledAddressesForChain(*big.Int) ([]common.Address, error)
}

type SequenceSyncer interface {
	LoadNextSequenceMap(context.Context, []common.Address)
	GetNextSequence(context.Context, common.Address) (evmtypes.Nonce, error)
	IncrementNextSequence(common.Address)
	SyncSequence(context.Context, common.Address, services.StopChan)
}

type Broadcaster struct {
	services.StateMachine
	txAttemptBuilder TxAttemptBuilder
	lggr             logger.SugaredLogger
	txStore          BroadcasterTxStore
	client           BroadcasterClient
	chainID          *big.Int
	config           BroadcasterConfig

	ks               KeyStore
	sequenceSyncer   SequenceSyncer
	enabledAddresses []common.Address

	triggers map[common.Address]chan struct{}

	chStop services.StopChan
	wg     sync.WaitGroup
}

func NewBroadcaster(
	txAttemptBuilder TxAttemptBuilder,
	lggr logger.Logger,
	txStore BroadcasterTxStore,
	client BroadcasterClient,
	config BroadcasterConfig,
	keystore KeyStore,
	sequenceSyncer SequenceSyncer,
) *Broadcaster {
	lggr = logger.Named(lggr, "Broadcaster")
	return &Broadcaster{
		txAttemptBuilder: txAttemptBuilder,
		lggr:             logger.Sugared(lggr),
		txStore:          txStore,
		client:           client,
		chainID:          client.ConfiguredChainID(),
		config:           config,
		ks:               keystore,
		sequenceSyncer:   sequenceSyncer,
	}
}

func (b *Broadcaster) Start(ctx context.Context) error {
	return b.StartOnce("Broadcaster", func() (err error) {
		// TODO: handle subscription to new addresses properly
		b.enabledAddresses, err = b.ks.EnabledAddressesForChain(b.chainID)
		if err != nil {
			return fmt.Errorf("Broadcaster: failed to load EnabledAddressesForChain: %w", err)
		}
		if len(b.enabledAddresses) > 0 {
			b.lggr.Debugw(fmt.Sprintf("Booting with %d keys", len(b.enabledAddresses)), "keys", b.enabledAddresses)
		} else {
			b.lggr.Warnf("Chain %s does not have any keys, no transactions will be sent on this chain", b.chainID.String())
		}

		b.chStop = make(chan struct{})
		b.wg = sync.WaitGroup{}
		b.wg.Add(len(b.enabledAddresses))
		b.triggers = make(map[common.Address]chan struct{})
		b.sequenceSyncer.LoadNextSequenceMap(ctx, b.enabledAddresses)

		for _, addr := range b.enabledAddresses {
			triggerCh := make(chan struct{}, 1)
			b.triggers[addr] = triggerCh
			go b.monitorTxs(addr, triggerCh)
		}
		return
	})
}

func (b *Broadcaster) Close() error {
	return b.StopOnce("Broadcaster", func() error {
		close(b.chStop)
		b.wg.Wait()
		return nil
	})
}

func (b *Broadcaster) HealthReport() map[string]error {
	return map[string]error{b.lggr.Name(): b.Healthy()}
}

func (b *Broadcaster) Trigger(addr common.Address) {
	if !b.IfStarted(func() {
		triggerCh, exists := b.triggers[addr]
		if !exists {
			return
		}
		select {
		case triggerCh <- struct{}{}:
		default:
		}
	}) {
		b.lggr.Debugf("Unstarted; ignoring trigger for %s", addr)
	}
}

func (b *Broadcaster) newResendBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}
}

func (b *Broadcaster) monitorTxs(addr common.Address, triggerCh chan struct{}) {
	defer b.wg.Done()

	ctx, cancel := b.chStop.NewCtx()
	defer cancel()

	if b.config.NonceAutoSync {
		b.lggr.Debugw("Auto-syncing sequence", "address", addr.String())
		b.sequenceSyncer.SyncSequence(ctx, addr, b.chStop)
		if ctx.Err() != nil {
			return
		}
	} else {
		b.lggr.Debugw("Skipping sequence auto-sync", "address", addr.String())
	}

	var errorRetryCh <-chan time.Time
	bf := b.newResendBackoff()

	for {
		pollDBTimer := time.NewTimer(utils.WithJitter(b.config.FallbackPollInterval))

		err := b.ProcessUnstartedTxs(ctx, addr)
		if err != nil {
			// On errors we implement exponential backoff retries. This
			// handles intermittent connectivity, remote RPC races, timing issues etc
			b.lggr.Errorw("Error occurred while handling tx queue in ProcessUnstartedTxs", "err", err)
			pollDBTimer.Reset(utils.WithJitter(b.config.FallbackPollInterval))
			errorRetryCh = time.After(bf.Duration())
		} else {
			bf = b.newResendBackoff()
			errorRetryCh = nil
		}

		select {
		case <-ctx.Done():
			// NOTE: See: https://godoc.org/time#Timer.Stop for an explanation of this pattern
			if !pollDBTimer.Stop() {
				<-pollDBTimer.C
			}
			return
		case <-triggerCh:
			// tx was inserted
			if !pollDBTimer.Stop() {
				<-pollDBTimer.C
			}
			continue
		case <-pollDBTimer.C:
			// DB poller timed out
			continue
		case <-errorRetryCh:
			// Error backoff period reached
			continue
		}
	}
}

func (b *Broadcaster) ProcessUnstartedTxs(ctx context.Context, fromAddress common.Address) (err error) {
	err = b.handleAnyInProgressTx(ctx, fromAddress)
	if err != nil {
		return fmt.Errorf("ProcessUnstartedTxs failed on handleAnyInProgressTx: %w", err)
	}
	for {
		maxInFlightTransactions := b.config.MaxInFlight
		if maxInFlightTransactions > 0 {
			nUnconfirmed, err := b.txStore.CountUnconfirmedTransactions(ctx, fromAddress, b.chainID)
			if err != nil {
				return fmt.Errorf("CountUnconfirmedTransactions failed: %w", err)
			}
			if nUnconfirmed >= maxInFlightTransactions {
				nUnstarted, err := b.txStore.CountUnstartedTransactions(ctx, fromAddress, b.chainID)
				if err != nil {
					return fmt.Errorf("CountUnstartedTransactions failed: %w", err)
				}
				b.lggr.Warnw(fmt.Sprintf(`Transaction throttling; %d transactions in-flight and %d unstarted transactions pending (maximum number of in-flight transactions is %d per key). %s`, nUnconfirmed, nUnstarted, maxInFlightTransactions, label.MaxInFlightTransactionsWarning), "maxInFlightTransactions", maxInFlightTransactions, "nUnconfirmed", nUnconfirmed, "nUnstarted", nUnstarted)
				select {
				case <-time.After(InFlightTransactionRecheckInterval):
				case <-ctx.Done():
					return context.Cause(ctx)
				}
				continue
			}
		}
		etx, err := b.nextUnstartedTransactionWithSequence(fromAddress)
		if err != nil {
			return fmt.Errorf("processUnstartedTxs failed on nextUnstartedTransactionWithSequence: %w", err)
		}
		if etx == nil {
			return nil
		}

		if err := b.handleInProgressTx(ctx, *etx); err != nil {
			return fmt.Errorf("processUnstartedTxs failed on handleUnstartedTx: %w", err)
		}
	}
}

func (b *Broadcaster) nextUnstartedTransactionWithSequence(fromAddress common.Address) (*txmgr.Tx, error) {
	ctx, cancel := b.chStop.NewCtx()
	defer cancel()
	tx := &txmgr.Tx{}
	if err := b.txStore.FindNextUnstartedTransactionFromAddress(ctx, tx, fromAddress, b.chainID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("FindNextUnstartedTransactionFromAddress failed: %w", err)
	}

	sequence, err := b.sequenceSyncer.GetNextSequence(ctx, tx.FromAddress)
	if err != nil {
		return nil, err
	}
	tx.Sequence = &sequence

	if tx.State != commontxmgr.TxUnstarted {
		return nil, fmt.Errorf("invariant violation: expected transaction %v to be unstarted, it was %s", tx.ID, tx.State)
	}

	if err = b.txStore.BroadcasterUpdateTxUnstartedToInProgress(ctx, tx); errors.Is(err, ErrTxRemoved) {
		b.lggr.Debugw("tx removed", "txID", tx.ID, "subject", tx.Subject)
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("failed on BroadcasterUpdateTxUnstartedToInProgress: %w", err)
	}
	return tx, nil
}

func (b *Broadcaster) handleAnyInProgressTx(ctx context.Context, fromAddress common.Address) (err error) {
	tx, err := b.txStore.BroadcasterGetTxInProgress(ctx, fromAddress)
	if err != nil {
		return fmt.Errorf("handleAnyInProgressTx failed: %w", err)
	}
	if tx == nil {
		return
	}

	return b.handleInProgressTx(ctx, *tx)
}

func (b *Broadcaster) handleInProgressTx(ctx context.Context, tx txmgr.Tx) error {
	if tx.State != commontxmgr.TxInProgress {
		return fmt.Errorf("invariant violation: expected transaction %v to be in_progress, it was %s", tx.ID, tx.State)
	}

	attempt, err := b.txAttemptBuilder.NewAttempt(ctx, tx, b.lggr)
	if err != nil {
		return fmt.Errorf("failed on NewAttempt: %w", err)
	}

	lgr := tx.GetLogger(logger.With(b.lggr))
	signedTx, err := txmgr.GetGethSignedTx(attempt.SignedRawTx)
	if err != nil {
		b.lggr.Criticalw("Fatal error signing transaction", "err", err, "tx", tx)
		return fmt.Errorf("error while sending transaction %s (tx ID %d): %w", attempt.Hash.String(), tx.ID, err)
	}
	err = b.client.SendTransaction(ctx, signedTx)
	timeStamp := time.Now()
	tx.InitialBroadcastAt = &timeStamp
	tx.BroadcastAt = &timeStamp

	lgr.Infow("Sent transaction", "tx", tx.PrettyPrint(), "attempt", attempt.PrettyPrint(), "error", err)

	if err != nil {
		nextSequence, e := b.client.PendingNonceAt(ctx, tx.FromAddress)
		if e != nil {
			err = multierr.Combine(e, err)
			return fmt.Errorf("error while sending transaction %s (tx ID %d): %w", attempt.Hash.String(), tx.ID, err)
		}
		lgr.Warnw(err.Error(), "attempt", attempt)
		if nextSequence <= (*tx.Sequence).Uint64() {
			return fmt.Errorf("error while sending transaction %s (tx ID %d): %w", attempt.Hash.String(), tx.ID, err)
		}
	}

	err = b.txStore.UpdateTxInProgressToUnconfirmed(ctx, &tx)
	if err != nil {
		return err
	}

	b.sequenceSyncer.IncrementNextSequence(tx.FromAddress)
	return err
}
