package soltxm

import (
	"context"
	"fmt"
	"strings"
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
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	MaxQueueLen      = 1000
	MaxRetryTimeMs   = 250 // max tx retry time (exponential retry will taper to retry every 0.25s)
	MaxSigsToConfirm = 256 // max number of signatures in GetSignatureStatus call
)

var (
	_ services.ServiceCtx = (*Txm)(nil)
	_ solana.TxManager    = (*Txm)(nil)
)

// Txm manages transactions for the solana blockchain.
// simple implementation with no persistently stored txs
type Txm struct {
	starter utils.StartStopOnce
	lggr    logger.Logger
	chSend  chan pendingTx
	chSim   chan pendingTx
	chStop  chan struct{}
	done    sync.WaitGroup
	cfg     config.Config
	txs     PendingTxContext
	ks      keystore.Solana
	client  *utils.LazyLoad[solanaClient.ReaderWriter]
}

type pendingTx struct {
	tx        *solanaGo.Transaction
	timeout   time.Duration
	signature solanaGo.Signature
}

// NewTxm creates a txm. Uses simulation so should only be used to send txes to trusted contracts i.e. OCR.
func NewTxm(chainID string, tc func() (solanaClient.ReaderWriter, error), cfg config.Config, ks keystore.Solana, lggr logger.Logger) *Txm {
	lggr = lggr.Named("Txm")
	return &Txm{
		starter: utils.StartStopOnce{},
		lggr:    lggr,
		chSend:  make(chan pendingTx, MaxQueueLen), // queue can support 1000 pending txs
		chSim:   make(chan pendingTx, MaxQueueLen), // queue can support 1000 pending txs
		chStop:  make(chan struct{}),
		cfg:     cfg,
		txs:     newPendingTxContextWithProm(chainID),
		ks:      ks,
		client:  utils.NewLazyLoad(tc),
	}
}

// Start subscribes to queuing channel and processes them.
func (txm *Txm) Start(context.Context) error {
	return txm.starter.StartOnce("solana_txm", func() error {
		txm.done.Add(3) // waitgroup: tx retry, confirmer, simulator
		go txm.run()
		return nil
	})
}

func (txm *Txm) run() {
	defer txm.done.Done()
	ctx, cancel := utils.ContextFromChan(txm.chStop)
	defer cancel()

	// start confirmer + simulator
	go txm.confirm(ctx)
	go txm.simulate(ctx)

	for {
		select {
		case msg := <-txm.chSend:
			// process tx
			sig, err := txm.sendWithRetry(ctx, msg.tx, msg.timeout)
			if err != nil {
				txm.lggr.Errorw("failed to send transaction", "error", err)
				txm.client.Reset() // clear client if tx fails immediately (potentially bad RPC)
				continue           // skip remainining
			}

			// send tx + signature to simulation queue
			msg.signature = sig
			select {
			case txm.chSim <- msg:
			default:
				txm.lggr.Warnw("failed to enqeue tx for simulation", "queueFull", len(txm.chSend) == MaxQueueLen, "tx", msg)
			}

			txm.lggr.Debugw("transaction sent", "signature", sig.String())
		case <-txm.chStop:
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
		cancel()                           // cancel context when exiting early
		txm.txs.OnError(sig, TxFailReject) // increment failed metric
		return solanaGo.Signature{}, errors.Wrap(err, "tx failed initial transmit")
	}

	// store tx signature + cancel function
	if err := txm.txs.Add(sig, cancel); err != nil {
		cancel() // cancel context when exiting early
		return solanaGo.Signature{}, errors.Wrapf(err, "failed to save tx signature (%s) to inflight txs", sig)
	}

	// retry with exponential backoff
	// until context cancelled by timeout or called externally
	go func() {
		deltaT := 1 // ms
		tick := time.After(0)
		for {
			select {
			case <-ctx.Done():
				// stop sending tx after retry tx ctx times out (does not stop confirmation polling for tx)
				txm.lggr.Debugw("stopped tx retry", "signature", sig)
				return
			case <-tick:
				go func() {
					retrySig, err := client.SendTx(ctx, tx)
					// this could occur if endpoint goes down or if ctx cancelled
					if err != nil {
						if strings.Contains(err.Error(), "context canceled") || strings.Contains(err.Error(), "context deadline exceeded") {
							txm.lggr.Debugw("ctx error on send retry transaction", "error", err, "signature", sig)
						} else {
							txm.lggr.Warnw("failed to send retry transaction", "error", err, "signature", sig)
						}
						return
					}
					// this should never happen
					if retrySig != sig {
						txm.lggr.Criticalw("original signature does not match retry signature", "expectedSignature", sig, "receivedSignature", retrySig)
					}
				}()
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
			sigs := txm.txs.ListAll()

			// exit switch if not txs to confirm
			if len(sigs) == 0 {
				break
			}

			// get client
			client, err := txm.client.Get()
			if err != nil {
				txm.lggr.Errorw("failed to get client in soltxm.confirm", "error", err)
				break // exit switch
			}

			// batch sigs no more than MaxSigsToConfirm each
			sigsBatch, err := utils.BatchSplit(sigs, MaxSigsToConfirm)
			if err != nil { // this should never happen
				txm.lggr.Criticalw("failed to batch signatures", "error", err)
				break // exit switch
			}

			// process signatures
			processSigs := func(s []solanaGo.Signature, res []*rpc.SignatureStatusesResult) {
				for i := 0; i < len(res); i++ {
					// if status is nil (sig not found), continue polling
					// sig not found could mean invalid tx or not picked up yet
					if res[i] == nil {
						txm.lggr.Debugw("tx state: not found",
							"signature", s[i],
						)

						// check confirm timeout exceeded
						if txm.txs.Expired(s[i], txm.cfg.TxConfirmTimeout()) {
							txm.txs.OnError(s[i], TxFailDrop)
							txm.lggr.Warnw("failed to find transaction within confirm timeout", "signature", s[i], "timeoutSeconds", txm.cfg.TxConfirmTimeout())
						}
						continue
					}

					// if signature has an error, end polling
					if res[i].Err != nil {
						txm.lggr.Errorw("tx state: failed",
							"signature", s[i],
							"error", res[i].Err,
							"status", res[i].ConfirmationStatus,
						)
						txm.txs.OnError(s[i], TxFailRevert)
						continue
					}

					// if signature is processed, keep polling
					if res[i].ConfirmationStatus == rpc.ConfirmationStatusProcessed {
						txm.lggr.Debugw("tx state: processed",
							"signature", s[i],
						)

						// check confirm timeout exceeded
						if txm.txs.Expired(s[i], txm.cfg.TxConfirmTimeout()) {
							txm.txs.OnError(s[i], TxFailDrop)
							txm.lggr.Warnw("tx failed to move beyond 'processed' within confirm timeout", "signature", s[i], "timeoutSeconds", txm.cfg.TxConfirmTimeout())
						}
						continue
					}

					// if signature is confirmed/finalized, end polling
					if res[i].ConfirmationStatus == rpc.ConfirmationStatusConfirmed || res[i].ConfirmationStatus == rpc.ConfirmationStatusFinalized {
						txm.lggr.Debugw(fmt.Sprintf("tx state: %s", res[i].ConfirmationStatus),
							"signature", s[i],
						)
						txm.txs.OnSuccess(s[i])
						continue
					}
				}
			}

			// waitgroup for processing
			var wg sync.WaitGroup
			wg.Add(len(sigsBatch))

			// loop through batch
			for i := 0; i < len(sigsBatch); i++ {
				// fetch signature statuses
				statuses, err := client.SignatureStatuses(ctx, sigsBatch[i])
				if err != nil {
					txm.lggr.Errorw("failed to get signature statuses in soltxm.confirm", "error", err)
					wg.Done() // don't block if exit early
					break     // exit for loop
				}

				// nonblocking: process batches as soon as they come in
				go func(index int) {
					defer wg.Done()
					processSigs(sigsBatch[index], statuses)
				}(i)
			}
			wg.Wait() // wait for processing to finish
		}
		tick = time.After(utils.WithJitter(txm.cfg.ConfirmPollPeriod()))
	}
}

// goroutine that simulates tx (use a bounded number of goroutines to pick from queue?)
// simulate can cancel the send retry function early in the tx management process
// additionally, it can provide reasons for why a tx failed in the logs
func (txm *Txm) simulate(ctx context.Context) {
	defer txm.done.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-txm.chSim:
			// get client
			client, err := txm.client.Get()
			if err != nil {
				txm.lggr.Errorw("failed to get client in soltxm.simulate", "error", err)
				continue
			}

			res, err := client.SimulateTx(ctx, msg.tx, nil) // use default options
			if err != nil {
				// this error can occur if endpoint goes down or if invalid signature (invalid signature should occur further upstream in sendWithRetry)
				// allow retry to continue in case temporary endpoint failure (if still invalid, confirm or timeout will cleanup)
				txm.lggr.Errorw("failed to simulate tx", "signature", msg.signature, "error", err)
				continue
			}

			// continue if simulation does not return error continue
			if res.Err == nil {
				continue
			}

			// handle various errors
			// https://github.com/solana-labs/solana/blob/master/sdk/src/transaction/error.rs
			// ---
			errStr := fmt.Sprintf("%v", res.Err) // convert to string to handle various interfaces
			switch {
			// blockhash not found when simulating, occurs when network bank has not seen the given blockhash or tx is too old
			// let simulation process/clean up
			case strings.Contains(errStr, "BlockhashNotFound"):
				txm.lggr.Warnw("simulate: BlockhashNotFound", "signature", msg.signature, "result", res)
				continue
			// transaction will encounter execution error/revert, mark as reverted to remove from confirmation + retry
			case strings.Contains(errStr, "InstructionError"):
				txm.txs.OnError(msg.signature, TxFailSimRevert) // cancel retry
				txm.lggr.Warnw("simulate: InstructionError", "signature", msg.signature, "result", res)
				continue
			// transaction is already processed in the chain, letting txm confirmation handle
			case strings.Contains(errStr, "AlreadyProcessed"):
				txm.lggr.Debugw("simulate: AlreadyProcessed", "signature", msg.signature, "result", res)
				continue
			// unrecognized errors (indicates more concerning failures)
			default:
				txm.txs.OnError(msg.signature, TxFailSimOther) // cancel retry
				txm.lggr.Errorw("simulate: unrecognized error", "signature", msg.signature, "result", res)
				continue
			}
		}
	}
}

// Enqueue enqueue a msg destined for the solana chain.
func (txm *Txm) Enqueue(accountID string, tx *solanaGo.Transaction) error {
	// validate nil pointer
	if tx == nil {
		return errors.New("error in soltxm.Enqueue: tx is nil pointer")
	}
	// validate account keys slice
	if len(tx.Message.AccountKeys) == 0 {
		return errors.New("error in soltxm.Enqueue: not enough account keys in tx")
	}

	// get key
	// fee payer account is index 0 account
	// https://github.com/gagliardetto/solana-go/blob/main/transaction.go#L252
	key, err := txm.ks.Get(tx.Message.AccountKeys[0].String())
	if err != nil {
		return errors.Wrap(err, "error in soltxm.Enqueue.GetKey")
	}
	txMsg, err := tx.Message.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "error in soltxm.Enqueue.MarshalBinary")
	}
	// sign tx
	sigBytes, err := key.Sign(txMsg)
	if err != nil {
		return errors.Wrap(err, "error in soltxm.Enqueue.Sign")
	}
	var finalSig [64]byte
	copy(finalSig[:], sigBytes)
	tx.Signatures = append(tx.Signatures, finalSig)

	msg := pendingTx{
		tx:      tx,
		timeout: txm.cfg.TxRetryTimeout(),
	}

	select {
	case txm.chSend <- msg:
	default:
		txm.lggr.Errorw("failed to enqeue tx", "queueFull", len(txm.chSend) == MaxQueueLen, "tx", msg)
		return errors.Errorf("failed to enqueue transaction for %s", accountID)
	}
	return nil
}

func (txm *Txm) InflightTxs() int {
	return len(txm.txs.ListAll())
}

// Close close service
func (txm *Txm) Close() error {
	return txm.starter.StopOnce("solanatxm", func() error {
		close(txm.chStop)
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
