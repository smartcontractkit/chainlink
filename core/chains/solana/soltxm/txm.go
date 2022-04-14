package soltxm

import (
	"context"
	"fmt"
	"sync"
	"time"

	solanaGo "github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	MaxQueueLen    = 1000
	MaxRetryTimeMs = 500
)

var (
	_ services.ServiceCtx = (*Txm)(nil)
	_ solana.TxManager    = (*Txm)(nil)
)

// Txm manages transactions for the solana blockchain.
// simple implementation with no persistently stored txs
type Txm struct {
	starter  utils.StartStopOnce
	lggr     logger.Logger
	tc       func() (solanaClient.ReaderWriter, error)
	queue    chan queueMsg
	simQueue chan queueMsg
	stop     chan struct{}
	done     sync.WaitGroup
	cfg      config.Config
	txCache  *TxCache
	client   solanaClient.ReaderWriter
}

type queueMsg struct {
	tx        *solanaGo.Transaction
	timeout   time.Duration
	signature solanaGo.Signature
}

// NewTxm creates a txm. Uses simulation so should only be used to send txes to trusted contracts i.e. OCR.
func NewTxm(tc func() (solanaClient.ReaderWriter, error), cfg config.Config, lggr logger.Logger) *Txm {
	lggr = lggr.Named("Txm")
	return &Txm{
		starter:  utils.StartStopOnce{},
		tc:       tc,
		lggr:     lggr,
		queue:    make(chan queueMsg, MaxQueueLen), // queue can support 1000 pending txs
		simQueue: make(chan queueMsg, MaxQueueLen), // queue can support 1000 pending txs
		stop:     make(chan struct{}),
		cfg:      cfg,
		txCache:  NewTxCache(),
	}
}

// Start subscribes to queuing channel and processes them.
func (txm *Txm) Start(context.Context) error {
	return txm.starter.StartOnce("solanatxm", func() error {
		go txm.run()
		return nil
	})
}

func (txm *Txm) run() {
	txm.done.Add(1)
	defer txm.done.Done()
	ctx, cancel := utils.ContextFromChan(txm.stop)
	defer cancel()

	// start confirmer + simulator
	go txm.confirm(ctx)
	go txm.simulate(ctx)

	for {
		select {
		case msg := <-txm.queue:
			// fetch client (only if it doesn't exist, reduce need for db read each time)
			if txm.client == nil {
				newClient, err := txm.tc()
				if err != nil {
					txm.lggr.Errorw("failed to get client", "err", err)
					continue
				}
				txm.client = newClient
			}
			// process tx
			sig, err := txm.sendWithRetry(ctx, msg.tx, msg.timeout)
			if err != nil {
				txm.lggr.Criticalw("failed to send transaction", "err", err)
				txm.client = nil // clear client if tx fails immediately (potentially bad RPC)
				continue         // skip remainining
			}

			// send tx + signature to simulation queue
			msg.signature = sig
			txm.simQueue <- msg

			txm.lggr.Debugw("sent transaction", "signature", sig.String())
		case <-txm.stop:
			return
		}
	}
}

func (txm *Txm) sendWithRetry(chanCtx context.Context, tx *solanaGo.Transaction, timeout time.Duration) (solanaGo.Signature, error) {
	if txm.client == nil {
		return solanaGo.Signature{}, errors.New("transaction manager client is nil")
	}

	// create timeout context
	ctx, cancel := context.WithTimeout(chanCtx, timeout)

	// send initial tx (do not retry and exit early if fails)
	sig, err := txm.client.SendTx(ctx, tx)
	if err != nil {
		cancel() // cancel context when exiting early
		return solanaGo.Signature{}, errors.Wrap(err, "tx failed initial transmit")
	}

	// store tx signature + cancel function
	if err := txm.txCache.Insert(sig, cancel); err != nil {
		return solanaGo.Signature{}, errors.Wrap(err, "failed to save tx signature to cache")
	}

	// retry with exponential backoff
	// until context cancelled by timeout or called externally
	go func() {
		deltaT := 1
		tick := time.After(0)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick:
				retrySig, err := txm.client.SendTx(ctx, tx)
				// this could occur if endpoint goes down
				if err != nil {
					txm.lggr.Criticalw("failed to send retry transaction", "err", err, "signature", retrySig)
				}
				// this should never happen
				if retrySig != sig {
					txm.lggr.Criticalw("original signature does not match retry signature", "expectedSignature", sig, "receivedSignature", retrySig)
				}
			}

			// exponential increase in wait time, capped at 500ms
			deltaT *= 2
			if deltaT > MaxRetryTimeMs {
				deltaT = MaxRetryTimeMs
			}
			tick = time.After(time.Duration(deltaT) * time.Millisecond)
		}
	}()

	// return signature for use in simulation
	return sig, nil
}

// goroutine that polls to confirm implementation
// cancels the exponential retry once confirmed
func (txm *Txm) confirm(ctx context.Context) {
	txm.done.Add(1)
	defer txm.done.Done()

	tick := time.After(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			// TODO: implement confirmation
			fmt.Println("PLACEHOLDER: cache list", txm.txCache.List())
		}
		// TODO: defaults to 1s - should change to 0.5s
		tick = time.After(utils.WithJitter(txm.cfg.ConfirmPollPeriod()))
	}
}

// goroutine that simulates tx (use a bounded number of goroutines to pick from queue?)
func (txm *Txm) simulate(ctx context.Context) {
	txm.done.Add(1)
	defer txm.done.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-txm.simQueue:
			// TODO: simulate tx that is passed from queue
			fmt.Println("PLACEHOLDER: simQueue", msg)
		}
	}
}

// Enqueue enqueue a msg destined for the solana chain.
func (txm *Txm) Enqueue(accountID string, tx *solanaGo.Transaction) error {
	msg := queueMsg{
		tx:      tx,
		timeout: 5 * time.Second, // TODO: make this configurable via argument,
	}

	select {
	case txm.queue <- msg:
	default:
		txm.lggr.Errorw("failed to enqeue tx", "queueLength", len(txm.queue), "tx", msg)
		return errors.Errorf("failed to enqueue transaction for %s", accountID)
	}
	return nil
}

// Close close service
func (txm *Txm) Close() error {
	return txm.starter.StopOnce("solanatxm", func() error {
		close(txm.stop)
		txm.done.Wait()
		return nil
	})
}

// Healthy service is healthy
func (txm *Txm) Healthy() error {
	return nil
}

// Ready service is ready
func (txm *Txm) Ready() error {
	return nil
}
