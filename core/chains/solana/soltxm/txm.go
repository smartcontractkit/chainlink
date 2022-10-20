package soltxm

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	solanaGo "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink/core/chains/solana/fees"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	MaxQueueLen      = 1000
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
	chSend  chan PendingTx
	chStop  chan struct{}
	done    sync.WaitGroup
	cfg     config.Config
	txs     PendingTxs // interface so DB can be plugged in
	ks      keystore.Solana
	client  *utils.LazyLoad[solanaClient.ReaderWriter]

	// compute budget unit price parameters
	fee fees.Estimator
}

// NewTxm creates a txm. Uses simulation so should only be used to send txes to trusted contracts i.e. OCR.
func NewTxm(chainID string, tc func() (solanaClient.ReaderWriter, error), cfg config.Config, ks keystore.Solana, lggr logger.Logger) *Txm {
	lggr = lggr.Named("Txm")
	return &Txm{
		starter: utils.StartStopOnce{},
		lggr:    lggr,
		chSend:  make(chan PendingTx, MaxQueueLen), // queue can support 1000 pending txs
		chStop:  make(chan struct{}),
		cfg:     cfg,
		txs:     newPendingTxMemoryWithProm(chainID),
		ks:      ks,
		client:  utils.NewLazyLoad(tc),
	}
}

// Start subscribes to queuing channel and processes them.
func (txm *Txm) Start(ctx context.Context) error {
	return txm.starter.StartOnce("solana_txm", func() error {

		// determine estimator type
		var estimator fees.Estimator
		var err error
		switch strings.ToLower(txm.cfg.FeeEstimatorMode()) {
		case "fixed":
			estimator, err = fees.NewFixedPriceEstimator(txm.cfg)
		case "recentfees":
			estimator, err = fees.NewRecentFeeEstimator(txm.cfg)
		default:
			err = fmt.Errorf("unknown solana fee estimator type: %s", txm.cfg.FeeEstimatorMode())
		}
		if err != nil {
			return err
		}
		txm.fee = estimator
		if err := txm.fee.Start(ctx); err != nil {
			return err
		}

		txm.done.Add(2) // waitgroup: broadcaster, simulator
		go txm.run()
		return nil
	})
}

func (txm *Txm) run() {
	defer txm.done.Done()
	ctx, cancel := utils.ContextFromChan(txm.chStop)
	defer cancel()

	// start confirmer
	go txm.confirm(ctx)

	for {
		select {
		case msg := <-txm.chSend:
			var id uuid.UUID

			// how to determine ID
			if msg.broadcast { // if msg has been successfully broadcast, retrieve from tx
				// this should never happen
				if len(msg.signatures) == 0 {
					txm.lggr.Errorw("tx was already broadcast but has no signatures - dropping from queue", "tx", msg)
					continue
				}
				// this can happen if multiple txs are broadcast for 1 base tx, and a signature is confirmed but others are not
				tx, exists := txm.txs.GetBySignature(msg.signatures[0])
				if !exists {
					txm.lggr.Warnw("signature does not match a tx ID - it likely has been confirmed", "id", msg.id, "signatures", msg.signatures)
					continue
				}
				id = tx.id
				txm.lggr.Debugw("rebroadcasting tx (unconfirmed)", "id", msg.id, "previous_signatures", msg.signatures, "count", len(msg.signatures))
			} else if !msg.broadcast && msg.id != uuid.Nil { // if msg has not been successfully broadcast, but has been saved (indicates RPC failure)
				id = msg.id
				txm.lggr.Debugw("rebroadcasting tx (rpc rejection)", "id", msg.id)
			} else { // new transaction
				id = txm.txs.New(msg)
			}

			// Set compute unit price for transaction, returns a copy of the base tx
			tx, price, err := msg.SetComputeUnitPrice(
				txm.fee.BaseComputeUnitPrice(),
				txm.cfg.MinComputeUnitPrice(),
				txm.cfg.MaxComputeUnitPrice(),
			)
			if err != nil { // should never happen, skip tx if this occurs
				txm.lggr.Errorw("failed to set compute unit price in tx", "error", err, "tx", msg)
				continue
			}

			// marshal + sign transaction
			txMsg, err := tx.Message.MarshalBinary()
			if err != nil {
				txm.lggr.Errorw("failed to marshal transaction message", "error", err)
				continue
			}
			sigBytes, err := msg.key.Sign(txMsg) // sign with stored key
			if err != nil {
				txm.lggr.Errorw("failed to sign transaction", "error", err)
				continue
			}
			var finalSig [64]byte
			copy(finalSig[:], sigBytes)
			tx.Signatures = append(tx.Signatures, finalSig)

			// process tx
			sig, validBlockhash, err := txm.send(ctx, tx)
			if err != nil {
				txm.lggr.Errorw("failed to send transaction", "error", err, "id", id, "tx", tx)
				txm.client.Reset() // clear client if tx fails immediately (potentially bad RPC)

				// incrementing metric, signature will be 0... (won't match a tx)
				if _, storeErr := txm.txs.OnError(sig, TxRPCReject); storeErr != nil {
					txm.lggr.Errorw("failed to mark tx as errored", "id", id, "signature", sig, "error", storeErr)
				}
				msg.id = id // set ID on msg to indicate already tracked

				// if rpc rejected: retry tx and don't bump fee
				// (fee is not bumped if txs.Add is not called => `broadcast = false`)
				txm.chSend <- msg
				continue // skip remainining
			}

			// if invalid blockhash, remove and skip tx
			if !validBlockhash {
				txm.lggr.Warnw("removing tx - invalid blockhash", "id", id, "blockhash", tx.Message.RecentBlockhash)
				if err := txm.txs.Remove(id); err != nil {
					txm.lggr.Debugw("failed to remove tx ID - it has likely already been removed", "error", err)
					continue // skip incrementing error if tx already removed
				}
				txm.txs.OnError(sig, TxInvalidBlockhash)
				continue
			}

			txm.lggr.Debugw("transaction sent", "signature", sig.String())

			// store tx signature
			if err := txm.txs.Add(id, sig, price); err != nil {
				// this can occur if a duplicate transaction is broadcast
				txm.lggr.Errorw("failed to save tx signature", "signature", sig, "error", err)

				// handle duplicate transcations
				// check if TX has any associated signatures, remove if not (indicates a duplicate transaction)
				if tx, exists := txm.txs.GetByID(id); exists && len(tx.signatures) == 0 {
					txm.lggr.Debugw("removing tx - duplicate signature", "id", id, "signature", sig)
					if err := txm.txs.Remove(id); err != nil {
						txm.lggr.Debugw("failed to remove tx ID - it has likely already been removed", "error", err)
					}
				}
			}
		case <-txm.chStop:
			return
		}
	}
}

// sendWithExpBackoff broadcasts a transaction at an exponential backoff rate to increase chances of inclusion by the next validator by rebroadcasting more tx packets
func (txm *Txm) send(chanCtx context.Context, tx *solanaGo.Transaction) (sig solanaGo.Signature, validBlockhash bool, err error) {
	// fetch client
	client, err := txm.client.Get()
	if err != nil {
		return solanaGo.Signature{}, validBlockhash, errors.Wrap(err, "failed to get client in soltxm.sendWithExpBackoff")
	}

	// create timeout context
	ctx, cancel := context.WithTimeout(chanCtx, txm.cfg.TxTimeout())
	defer cancel()

	// validate block hash
	validBlockhash, err = client.IsBlockhashValid(ctx, tx.Message.RecentBlockhash)
	if err != nil {
		return sig, validBlockhash, errors.Wrap(err, "err in txm.send.IsBlockhashValid")
	}
	if !validBlockhash {
		return sig, validBlockhash, nil
	}

	// send tx
	sig, err = client.SendTx(ctx, tx) // returns 000.. signature if errors
	if err != nil {
		return sig, validBlockhash, errors.Wrap(err, "tx failed transmit")
	}

	// return signature for use in simulation
	return sig, validBlockhash, nil
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
			sigs := txm.txs.ListSignatures()

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
				txm.lggr.Errorw("failed to batch signatures", "error", err)
				break // exit switch
			}

			// track which signatures statuses
			rebroadcast := map[uuid.UUID]solanaGo.Signature{} // map used to coalesce many sigs => 1 base tx
			var rebroadcastLock sync.Mutex
			success := []solanaGo.Signature{}
			var successLock sync.Mutex
			reverted := []solanaGo.Signature{}
			var revertedLock sync.Mutex

			// process signatures
			processSigs := func(s []solanaGo.Signature, res []*rpc.SignatureStatusesResult) {
				for i := 0; i < len(s); i++ {
					// defensive: if len(res) < len(s), exit processing to prevent panic
					// this should never happen
					if i > len(res)-1 {
						txm.lggr.Error(fmt.Sprintf("mismatch requested signatures and responses length: %d > %d", len(s), len(res)))
						return
					}

					// if status is nil (sig not found), continue polling
					// sig not found could mean invalid tx or not picked up yet
					if res[i] == nil {
						txm.lggr.Debugw("tx state: not found",
							"signature", s[i],
						)

						// check confirm timeout exceeded, store for queuing again
						if tx, exists := txm.txs.GetBySignature(s[i]); exists && time.Since(tx.timestamp) > txm.cfg.TxConfirmTimeout() {
							// only set if it hasn't been set yet, coalesce signatures => tx ID (deduplication)
							rebroadcastLock.Lock()
							if _, exists := rebroadcast[tx.id]; !exists {
								rebroadcast[tx.id] = s[i]
							}
							rebroadcastLock.Unlock()
						}
						continue
					}

					// if signature has an error, end polling
					if res[i].Err != nil {
						txm.lggr.Debugw("tx state: errored",
							"signature", s[i],
							"error", res[i].Err,
							"status", res[i].ConfirmationStatus,
						)
						revertedLock.Lock()
						reverted = append(reverted, s[i])
						revertedLock.Unlock()
						continue
					}

					// if signature is processed, keep polling, don't retry yet (either will become confirmed or dropped)
					// if signature is confirmed/finalized, end polling
					if res[i].ConfirmationStatus == rpc.ConfirmationStatusConfirmed || res[i].ConfirmationStatus == rpc.ConfirmationStatusFinalized {
						txm.lggr.Debugw(fmt.Sprintf("tx state: %s", res[i].ConfirmationStatus),
							"signature", s[i],
						)
						successLock.Lock()
						success = append(success, s[i])
						successLock.Unlock()
						continue
					}
				}
			}

			// waitgroup for processing
			var wg sync.WaitGroup
			wg.Add(len(sigsBatch))

			// loop through batch
			for i := 0; i < len(sigsBatch); i++ {
				// nonblocking: process batches in parallel as soon as they come in
				go func(index int) {
					defer wg.Done()

					// fetch signature statuses
					statuses, err := client.SignatureStatuses(ctx, sigsBatch[index])
					if err != nil {
						txm.lggr.Errorw("failed to get signature statuses in soltxm.confirm", "error", err)
						return
					}
					processSigs(sigsBatch[index], statuses)
				}(i)
			}
			wg.Wait() // wait for processing to finish

			// process successful first then reverted TXs for proper metric incrementing
			for _, s := range success {
				tx, err := txm.txs.OnSuccess(s)
				if err != nil {
					txm.lggr.Warnw("failed to mark tx as successful - tx likely completed by another signature", "error", err)
					continue
				}
				txm.lggr.Debugw("tx marked as success", "id", tx.id, "signatures", tx.signatures)
			}
			for _, s := range reverted {
				tx, err := txm.txs.OnError(s, TxFailRevert)
				if err != nil {
					txm.lggr.Warnw("failed to mark tx as errored - tx likely completed by another signature", "error", err)
					continue
				}
				txm.lggr.Debugw("tx marked as errored", "id", tx.id, "signatures", tx.signatures)
			}

			// check to make sure tx still exists after all signatures are processed, then rebroadcast
			for _, sig := range maps.Values(rebroadcast) {
				if tx, exists := txm.txs.GetBySignature(sig); exists {
					select {
					case txm.chSend <- tx:
					default:
						txm.lggr.Errorw("failed to enqeue tx", "queueFull", len(txm.chSend) == MaxQueueLen, "tx", tx)
					}
				}
			}
		}
		tick = time.After(utils.WithJitter(txm.cfg.ConfirmPollPeriod()))
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

	// get signer key
	// fee payer account is index 0 account
	// https://github.com/gagliardetto/solana-go/blob/main/transaction.go#L252
	key, err := txm.ks.Get(tx.Message.AccountKeys[0].String())
	if err != nil {
		return errors.Wrap(err, "error in soltxm.Enqueue.GetKey")
	}

	msg := PendingTx{
		baseTx: tx,
		key:    key,
	}

	select {
	case txm.chSend <- msg:
	default:
		txm.lggr.Errorw("failed to enqeue tx", "queueFull", len(txm.chSend) == MaxQueueLen, "tx", msg)
		return errors.Errorf("failed to enqueue transaction for %s", accountID)
	}
	return nil
}

// returns number of unique TXs + broadcasted signatures associated to unique TXs
func (txm *Txm) InflightTxs() (int, int) {
	return len(txm.txs.ListIDs()), len(txm.txs.ListSignatures())
}

// Close close service
func (txm *Txm) Close() error {
	return txm.starter.StopOnce("solanatxm", func() error {
		close(txm.chStop)
		txm.done.Wait()
		return txm.fee.Close()
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
