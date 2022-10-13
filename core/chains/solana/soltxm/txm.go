package soltxm

import (
	"context"
	"fmt"
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
	chSend  chan PendingTx
	chStop  chan struct{}
	done    sync.WaitGroup
	cfg     config.Config
	txs     PendingTxs // interface so DB can be plugged in
	ks      keystore.Solana
	client  *utils.LazyLoad[solanaClient.ReaderWriter]

	// compute budget unit price parameters
	fee     uint64
	feeLock sync.RWMutex
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
func (txm *Txm) Start(context.Context) error {
	return txm.starter.StartOnce("solana_txm", func() error {
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
			// check if tx has been broadcast before and fetch or create UUID
			if msg.broadcast {
				// this should never happen
				if len(msg.signatures) == 0 {
					txm.lggr.Criticalw("tx was already broadcast but has no signatures - dropping from queue", "tx", msg)
					continue
				}
				// this can happen if multiple txs are broadcast for 1 base tx, and a signature is confirmed but others are not
				tx, exists := txm.txs.Get(msg.signatures[0])
				if !exists {
					txm.lggr.Warnw("signature does not match a tx ID - it may have already been confirmed", "signatures", msg.signatures)
					continue
				}
				id = tx.id
			} else {
				id = txm.txs.New(msg)
			}

			// Set compute unit price for transaction, returns a copy of the base tx
			tx, price, err := msg.SetComputeUnitPrice(ComputeUnitPrice(txm.GetFee()))
			if err != nil { // should never happen
				txm.lggr.Criticalw("failed to set compute unit price in tx", "error", err)
				// TODO: requeue tx?
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
			sig, err := txm.send(ctx, id, tx)
			if err != nil {
				txm.lggr.Errorw("failed to send transaction", "error", err)
				txm.client.Reset() // clear client if tx fails immediately (potentially bad RPC)

				// TODO: incrementing metric will fail because of 0... signature
				if _, storeErr := txm.txs.OnError(sig, TxFailReject); storeErr != nil {
					txm.lggr.Errorw("failed to mark tx as errored", "id", id, "signature", sig, "error", storeErr)
				}
				// TODO: don't give up on tx
				// TODO: don't bump fee if RPC failed
				// TODO: handle failed b/c blockhash expired or invalid
				continue // skip remainining
			}

			txm.lggr.Debugw("transaction sent", "signature", sig.String())

			// store tx signature
			if err := txm.txs.Add(id, sig, price); err != nil {
				txm.lggr.Errorw("failed to save tx signature to inflight txs", "signature", sig, "error", err)
			}
		case <-txm.chStop:
			return
		}
	}
}

// sendWithExpBackoff broadcasts a transaction at an exponential backoff rate to increase chances of inclusion by the next validator by rebroadcasting more tx packets
func (txm *Txm) send(chanCtx context.Context, id uuid.UUID, tx *solanaGo.Transaction) (solanaGo.Signature, error) {
	// fetch client
	client, err := txm.client.Get()
	if err != nil {
		return solanaGo.Signature{}, errors.Wrap(err, "failed to get client in soltxm.sendWithExpBackoff")
	}

	// create timeout context
	ctx, cancel := context.WithTimeout(chanCtx, txm.cfg.TxTimeout())
	defer cancel()

	// send tx
	sig, err := client.SendTx(ctx, tx) // returns 000.. signature if errors
	if err != nil {
		return solanaGo.Signature{}, errors.Wrap(err, "tx failed transmit")
	}

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
				txm.lggr.Criticalw("failed to batch signatures", "error", err)
				break // exit switch
			}

			// track which signatures need to be rebroadcast
			// signature used to look up transaction, UUID used to coalesce many sigs => 1 base tx
			needsRebroadcast := map[uuid.UUID]solanaGo.Signature{}

			// process signatures
			processSigs := func(s []solanaGo.Signature, res []*rpc.SignatureStatusesResult) {
				for i := 0; i < len(s); i++ {
					// defensive: if len(res) < len(s), exit processing to prevent panic
					// this should never happen
					if i > len(res)-1 {
						txm.lggr.Criticalw("mismatch requested signatures and responses length: %d > %d", len(s), len(res))
						return
					}

					// if status is nil (sig not found), continue polling
					// sig not found could mean invalid tx or not picked up yet
					if res[i] == nil {
						txm.lggr.Debugw("tx state: not found",
							"signature", s[i],
						)

						// check confirm timeout exceeded, store for queuing again
						if tx, exists := txm.txs.Get(s[i]); exists && time.Since(tx.timestamp) > txm.cfg.TxConfirmTimeout() {
							// only set if it hasn't been set yet, coalesce signatures => tx ID (deduplication)
							if _, exists := needsRebroadcast[tx.id]; !exists {
								needsRebroadcast[tx.id] = s[i]
							}
						}
						continue
					}

					// if signature has an error, end polling
					if res[i].Err != nil {
						txm.lggr.Errorw("tx state: errored",
							"signature", s[i],
							"error", res[i].Err,
							"status", res[i].ConfirmationStatus,
						)
						tx, err := txm.txs.OnError(s[i], TxFailRevert)
						if err != nil {
							txm.lggr.Warnw("failed to mark tx as errored - tx likely completed by another signature", "error", err)
							continue
						}
						txm.lggr.Debugw("tx marked as errored", "id", tx.id, "signatures", tx.signatures)
						continue
					}

					// if signature is processed, keep polling, don't retry yet (either will become confirmed or dropped)
					// if signature is confirmed/finalized, end polling
					if res[i].ConfirmationStatus == rpc.ConfirmationStatusConfirmed || res[i].ConfirmationStatus == rpc.ConfirmationStatusFinalized {
						txm.lggr.Debugw(fmt.Sprintf("tx state: %s", res[i].ConfirmationStatus),
							"signature", s[i],
						)
						tx, err := txm.txs.OnSuccess(s[i])
						if err != nil {
							txm.lggr.Warnw("failed to mark tx as successful - tx likely completed by another signature", "error", err)
							continue
						}
						txm.lggr.Debugw("tx marked as success", "id", tx.id, "signatures", tx.signatures)
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

			// check to make sure tx still exists after all signatures are processed, then rebroadcast
			for _, sig := range maps.Values(needsRebroadcast) {
				if tx, exists := txm.txs.Get(sig); exists {
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

func (txm *Txm) InflightTxs() int {
	return len(txm.txs.ListSignatures())
}

func (txm *Txm) SetFee(v uint64) {
	txm.feeLock.Lock()
	defer txm.feeLock.Unlock()
	txm.fee = v
}

func (txm *Txm) GetFee() uint64 {
	txm.feeLock.RLock()
	defer txm.feeLock.RUnlock()
	return txm.fee
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
