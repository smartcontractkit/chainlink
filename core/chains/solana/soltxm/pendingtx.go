package soltxm

import (
	"fmt"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	"golang.org/x/exp/maps"
)

type PendingTx struct {
	key        solkey.Key
	baseTx     *solana.Transaction // original transaction (should not contain fee information)
	timestamp  time.Time           // when the current tx is broadcast
	signatures []solana.Signature
	currentFee uint64 // current fee for inflight tx
	broadcast  bool   // check to indicate if already broadcast
}

// SetComputeUnitPrice sets the compute unit price in micro-lamports, returns new tx
// add fee as the last instruction
// add fee program as last account key
// recreates some of the logic from: https://github.com/gagliardetto/solana-go/blob/main/transaction.go#L313
func (tx *PendingTx) SetComputeUnitPrice(price ComputeUnitPrice) (*solana.Transaction, error) {
	txWithFee := *tx.baseTx // make copy

	// find ComputeBudget program to accounts if it exists
	// reimplements HasAccount to retrieve index: https://github.com/gagliardetto/solana-go/blob/main/message.go#L228
	var exists bool
	var programIdx uint16
	for i, a := range txWithFee.Message.AccountKeys {
		if a.Equals(price.ProgramID()) {
			exists = true
			programIdx = uint16(i)
		}
	}
	// if it doesn't exist, add to account keys
	if !exists {
		txWithFee.Message.AccountKeys = append(txWithFee.Message.AccountKeys, price.ProgramID())
		programIdx = uint16(len(txWithFee.Message.AccountKeys) - 1) // last index of account keys

		// https://github.com/gagliardetto/solana-go/blob/main/transaction.go#L291
		txWithFee.Message.Header.NumReadonlyUnsignedAccounts++
	}

	// double fee if already broadcast and this is a retry
	if tx.broadcast {
		price = ComputeUnitPrice(tx.currentFee * tx.currentFee)

		// handle 0 case
		if tx.currentFee == 0 {
			price = 1
		}
	}

	// get instruction data
	data, err := price.Data()
	if err != nil {
		return nil, err
	}

	// build tx
	txWithFee.Message.Instructions = append([]solana.CompiledInstruction{{
		ProgramIDIndex: programIdx,
		Data:           data,
	}}, txWithFee.Message.Instructions...)

	// track current fee
	tx.currentFee = uint64(price)
	return &txWithFee, nil
}

type PendingTxs interface {
	New(tx PendingTx) uuid.UUID                   // save pendingTx
	Add(id uuid.UUID, sig solana.Signature) error // save signature after broadcasting
	Remove(id uuid.UUID)
	ListSignatures() []solana.Signature                    // get all signatures for pending txs
	Get(sig solana.Signature) (uuid.UUID, PendingTx, bool) // get tx from signature
	// state change hooks
	OnSuccess(sig solana.Signature)
	OnError(sig solana.Signature, errType int) // match err type using enum
}

var _ PendingTxs = &pendingTxMemory{}

// in memory version of PendingTxs
type pendingTxMemory struct {
	idMap  map[uuid.UUID]PendingTx        // map id to transaction data
	sigMap map[solana.Signature]uuid.UUID // map tx signature to id
	lock   sync.RWMutex
}

func newPendingTxMemory() *pendingTxMemory {
	return &pendingTxMemory{
		idMap:  map[uuid.UUID]PendingTx{},
		sigMap: map[solana.Signature]uuid.UUID{},
	}
}

func (txs *pendingTxMemory) New(tx PendingTx) uuid.UUID {
	id := uuid.New()

	txs.lock.Lock()
	defer txs.lock.Unlock()

	txs.idMap[id] = tx
	return id
}

func (txs *pendingTxMemory) Add(id uuid.UUID, sig solana.Signature) error {
	var tx PendingTx
	var exists bool

	// check exists
	txs.lock.RLock()
	if tx, exists = txs.idMap[id]; !exists {
		txs.lock.RUnlock()
		return fmt.Errorf("ID does not exist: %s", id)
	}
	txs.lock.RUnlock()

	// save signatures
	tx.signatures = append(tx.signatures, sig)

	// upgrade to write lock
	txs.lock.Lock()
	defer txs.lock.Unlock()
	txs.idMap[id] = tx
	txs.sigMap[sig] = id
	return nil
}

func (txs *pendingTxMemory) Remove(id uuid.UUID) {
	// check if already removed
	txs.lock.RLock()
	if _, exists := txs.idMap[id]; !exists {
		txs.lock.RUnlock()
		return
	}
	txs.lock.RUnlock()

	// upgrade to write lock if ID exists
	txs.lock.Lock()
	defer txs.lock.Unlock()
	for _, s := range txs.idMap[id].signatures {
		delete(txs.sigMap, s)
	}
	delete(txs.idMap, id)
}

func (txs *pendingTxMemory) ListSignatures() []solana.Signature {
	txs.lock.RLock()
	defer txs.lock.RUnlock()
	return maps.Keys(txs.sigMap)
}

func (txs *pendingTxMemory) Get(sig solana.Signature) (uuid.UUID, PendingTx, bool) {
	txs.lock.RLock()
	defer txs.lock.RUnlock()

	if id, idExists := txs.sigMap[sig]; idExists {
		if tx, txExists := txs.idMap[id]; txExists {
			return id, tx, true
		}
	}
	return uuid.UUID{}, PendingTx{}, false
}

func (txs *pendingTxMemory) OnSuccess(sig solana.Signature) {
	if id, _, exists := txs.Get(sig); exists {
		txs.Remove(id)
	}
}

func (txs *pendingTxMemory) OnError(sig solana.Signature, _ int) {
	if id, _, exists := txs.Get(sig); exists {
		txs.Remove(id)
	}
}

var _ PendingTxs = &pendingTxMemoryWithProm{}

type pendingTxMemoryWithProm struct {
	pendingTx *pendingTxMemory
	chainID   string
}

const (
	TxFailRevert = iota // execution revert
	TxFailReject        // rpc rejected transaction
)

func newPendingTxMemoryWithProm(id string) *pendingTxMemoryWithProm {
	return &pendingTxMemoryWithProm{
		chainID:   id,
		pendingTx: newPendingTxMemory(),
	}
}

func (txs *pendingTxMemoryWithProm) New(tx PendingTx) uuid.UUID {
	return txs.pendingTx.New(tx)
}

func (txs *pendingTxMemoryWithProm) Add(id uuid.UUID, sig solana.Signature) error {
	return txs.pendingTx.Add(id, sig)
}

func (txs *pendingTxMemoryWithProm) Remove(id uuid.UUID) {
	txs.pendingTx.Remove(id)
}

func (txs *pendingTxMemoryWithProm) ListSignatures() []solana.Signature {
	sigs := txs.pendingTx.ListSignatures()
	promSolTxmPendingTxs.WithLabelValues(txs.chainID).Set(float64(len(sigs)))
	return sigs
}

func (txs *pendingTxMemoryWithProm) Get(sig solana.Signature) (uuid.UUID, PendingTx, bool) {
	return txs.pendingTx.Get(sig)
}

// Success - tx included in block and confirmed
func (txs *pendingTxMemoryWithProm) OnSuccess(sig solana.Signature) {
	promSolTxmSuccessTxs.WithLabelValues(txs.chainID).Add(1)
	txs.pendingTx.OnSuccess(sig)
}

func (txs *pendingTxMemoryWithProm) OnError(sig solana.Signature, errType int) {
	switch errType {
	case TxFailRevert:
		promSolTxmRevertTxs.WithLabelValues(txs.chainID).Add(1)
	case TxFailReject:
		promSolTxmRejectTxs.WithLabelValues(txs.chainID).Add(1)
	}
	// increment total errors
	promSolTxmErrorTxs.WithLabelValues(txs.chainID).Add(1)
	txs.pendingTx.OnError(sig, errType)
}
