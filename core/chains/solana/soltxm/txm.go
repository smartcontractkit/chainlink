package soltxm

import (
	"context"
	"sync"
	"time"

	solanaGo "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
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
	queue    chan queueMsg
	simQueue chan queueMsg
	stop     chan struct{}
	done     sync.WaitGroup
	cfg      config.Config
	txCache  *TxCache
	client   *ValidClient
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
		lggr:     lggr,
		queue:    make(chan queueMsg, MaxQueueLen), // queue can support 1000 pending txs
		simQueue: make(chan queueMsg, MaxQueueLen), // queue can support 1000 pending txs
		stop:     make(chan struct{}),
		cfg:      cfg,
		txCache:  NewTxCache(),
		client:   NewValidClient(tc),
	}
}

// Start subscribes to queuing channel and processes them.
func (txm *Txm) Start(context.Context) error {
	return txm.starter.StartOnce("solanatxm", func() error {
		txm.done.Add(3) // waitgroup: tx retry, confirmer, simulator
		go txm.run()
		return nil
	})
}

func (txm *Txm) run() {
	defer txm.done.Done()
	ctx, cancel := utils.ContextFromChan(txm.stop)
	defer cancel()

	// start confirmer + simulator
	go txm.confirm(ctx)
	go txm.simulate(ctx)

	for {
		select {
		case msg := <-txm.queue:
			// process tx
			sig, err := txm.sendWithRetry(ctx, msg.tx, msg.timeout)
			if err != nil {
				txm.lggr.Errorw("failed to send transaction", "error", err)
				txm.client.Clear() // clear client if tx fails immediately (potentially bad RPC)
				continue           // skip remainining
			}

			// send tx + signature to simulation queue
			msg.signature = sig
			select {
			case txm.simQueue <- msg:
			default:
				txm.lggr.Warnw("failed to enqeue tx for simulation", "queueFull", len(txm.queue) == MaxQueueLen, "tx", msg)
			}

			txm.lggr.Debugw("transaction sent", "signature", sig.String())
		case <-txm.stop:
			return
		}
	}
}

func (txm *Txm) sendWithRetry(chanCtx context.Context, tx *solanaGo.Transaction, timeout time.Duration) (solanaGo.Signature, error) {
	// fetch client
	client, err := txm.client.Get()
	if err != nil {
		return solanaGo.Signature{}, errors.Wrap(err, "failed to get client in soltxm.sendWithRetry")
	}

	// create timeout context
	ctx, cancel := context.WithTimeout(chanCtx, timeout)

	// send initial tx (do not retry and exit early if fails)
	sig, err := client.SendTx(ctx, tx)
	if err != nil {
		cancel() // cancel context when exiting early
		return solanaGo.Signature{}, errors.Wrap(err, "tx failed initial transmit")
	}

	// store tx signature + cancel function
	if err := txm.txCache.Insert(sig, cancel); err != nil {
		cancel() // cancel context when exiting early
		return solanaGo.Signature{}, errors.Wrap(err, "failed to save tx signature to cache")
	}

	// retry with exponential backoff
	// until context cancelled by timeout or called externally
	go func() {
		// remove sig from cache if context is cancelled
		defer txm.txCache.Cancel(sig)

		deltaT := 1
		tick := time.After(0)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick:
				retrySig, err := client.SendTx(ctx, tx)
				// this could occur if endpoint goes down
				if err != nil {
					txm.lggr.Errorw("failed to send retry transaction", "error", err, "signature", retrySig)
					break // exit switch
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
	defer txm.done.Done()

	tick := time.After(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			// get list of tx signatures to confirm
			sigs := txm.txCache.List()

			// skip loop if not txs to confirm
			if len(sigs) == 0 {
				break
			}

			// get client
			client, err := txm.client.Get()
			if err != nil {
				txm.lggr.Errorw("failed to get client in soltxm.confirm", "error", err)
				break // exit switch
			}

			// fetch signature statuses
			statuses, err := client.SignatureStatuses(ctx, sigs)
			if err != nil {
				txm.lggr.Errorw("failed to get signature statuses in soltxm.confirm", "error", err)
				break // exit switch
			}

			// process signatures
			for i := 0; i < len(statuses); i++ {
				// if status is nil (sig not found), continue polling
				// sig not found could mean invalid tx or not picked up yet
				if statuses[i] == nil {
					txm.lggr.Debugw("transaction not found",
						"signature", sigs[i],
					)
					continue
				}

				// if signature has an error, end polling
				if statuses[i].Err != nil {
					txm.lggr.Errorw("transaction failed",
						"signature", sigs[i],
						"error", err,
						"status", statuses[i].ConfirmationStatus,
					)
					txm.txCache.Cancel(sigs[i])
					continue
				}

				// if signature is processed, keep polling
				if statuses[i].ConfirmationStatus == rpc.ConfirmationStatusProcessed {
					txm.lggr.Tracew("transaction processed",
						"signature", sigs[i],
					)
					continue
				}

				// if signature is confirmed, end polling
				if statuses[i].ConfirmationStatus == rpc.ConfirmationStatusConfirmed {
					txm.lggr.Debugw("transaction confirmed",
						"signature", sigs[i],
					)
					txm.txCache.Cancel(sigs[i])
					continue
				}
			}
		}
		// TODO: defaults to 1s - should change to 0.5s
		tick = time.After(utils.WithJitter(txm.cfg.ConfirmPollPeriod()))
	}
}

// goroutine that simulates tx (use a bounded number of goroutines to pick from queue?)
func (txm *Txm) simulate(ctx context.Context) {
	defer txm.done.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-txm.simQueue:
			// TODO: consider bounded number of workers for simulating
			// current: single worker
			// compared to tx retrying (unbounded goroutines)

			// get client
			client, err := txm.client.Get()
			if err != nil {
				txm.lggr.Errorw("failed to get client in soltxm.simulate", "error", err)
				break // exit switch
			}

			res, err := client.SimulateTx(ctx, msg.tx, nil) // use default options
			if err != nil {
				txm.lggr.Errorw("failed to simulate tx", "signature", msg.signature, "error", err)
				break // exit switch
			}

			// stop tx retrying if simulate returns error
			if res.Err != nil {
				txm.txCache.Cancel(msg.signature) // cancel retry
				txm.lggr.Errorw("simulate error", "signature", msg.signature, "error", res)
				break
			}
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
		txm.lggr.Errorw("failed to enqeue tx", "queueFull", len(txm.queue) == MaxQueueLen, "tx", msg)
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
