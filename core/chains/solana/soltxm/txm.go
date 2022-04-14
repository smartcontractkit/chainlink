package soltxm

import (
	"context"
	"fmt"
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
	starter    utils.StartStopOnce
	lggr       logger.Logger
	tc         func() (solanaClient.ReaderWriter, error)
	queue      chan *solanaGo.Transaction
	stop, done chan struct{}
	cfg        config.Config
	txCache    *TxCache
	client     solanaClient.ReaderWriter
}

// NewTxm creates a txm. Uses simulation so should only be used to send txes to trusted contracts i.e. OCR.
func NewTxm(tc func() (solanaClient.ReaderWriter, error), cfg config.Config, lggr logger.Logger) *Txm {
	lggr = lggr.Named("Txm")
	return &Txm{
		starter: utils.StartStopOnce{},
		tc:      tc,
		lggr:    lggr,
		queue:   make(chan *solanaGo.Transaction, MaxQueueLen), // queue can support 1000 pending txs
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		cfg:     cfg,
		txCache: NewTxCache(),
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
	defer close(txm.done)
	ctx, cancel := utils.ContextFromChan(txm.stop)
	defer cancel()

	for {
		select {
		case tx := <-txm.queue:
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
			sig, err := txm.sendWithRetry(ctx, tx)
			if err != nil {
				txm.lggr.Criticalw("failed to send transaction", "err", err)
				txm.client = nil // clear client if tx fails immediately (potentially bad RPC)
				continue         // skip remainining
			}

			// send tx + signature to simulation queue

			txm.lggr.Debugw("successfully sent transaction", "signature", sig.String())
		case <-txm.stop:
			return
		}
	}
}

func (txm *Txm) sendWithRetry(chanCtx context.Context, tx *solanaGo.Transaction) (solanaGo.Signature, error) {
	if txm.client == nil {
		return solanaGo.Signature{}, errors.New("transaction manager client is nil")
	}

	// create timeout context
	// TODO: pass in timeout parameter
	timeout := 5 * time.Second
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
			fmt.Println(deltaT)
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

// TODO: goroutine that polls to confirm implementation
// cancels the exponential retry once confirmed
func (txm *Txm) confirm() {}

// TODO: goroutine that simulates tx (use a bounded number of goroutines to pick from queue?)
func (txm *Txm) simulate() {}

// Enqueue enqueue a msg destined for the solana chain.
func (txm *Txm) Enqueue(accountID string, msg *solanaGo.Transaction) error {
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
		<-txm.done
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
