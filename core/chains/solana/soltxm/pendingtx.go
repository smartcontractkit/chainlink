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
	id         uuid.UUID
	key        solkey.Key
	baseTx     *solana.Transaction // original transaction (should not contain fee information)
	timestamp  time.Time           // when the current tx is broadcast
	signatures []solana.Signature  // broadcasted tx signatures
	currentFee uint64              // current fee for inflight tx
	broadcast  bool                // check to indicate if already broadcast before
}

// SetComputeUnitPrice sets the compute unit price in micro-lamports, returns new tx
// add fee as the last instruction
// add fee program as last account key
// recreates some of the logic from: https://github.com/gagliardetto/solana-go/blob/main/transaction.go#L313
func (tx *PendingTx) SetComputeUnitPrice(base, min, max uint64) (*solana.Transaction, uint64, error) {
	// input validation
	if base < min || base > max || min > max {
		return nil, 0, fmt.Errorf("invalid inputs: %d <= %d <= %d (not true)", min, base, max)
	}

	txWithFee := *tx.baseTx // make copy

	// find ComputeBudget program to accounts if it exists
	// reimplements HasAccount to retrieve index: https://github.com/gagliardetto/solana-go/blob/main/message.go#L228
	var exists bool
	var programIdx uint16
	price := ComputeUnitPrice(base)
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

	// double fee if already successfully broadcast and this is a retry
	if tx.broadcast {
		price = ComputeUnitPrice(2 * tx.currentFee)

		// handle 0 case
		if tx.currentFee == 0 {
			price = 1
		}

		// handle 1 case
		if tx.currentFee == 1 {
			price = 2
		}
	}

	// handle bounds
	if uint64(price) < min {
		price = ComputeUnitPrice(min)
	}
	if uint64(price) > max {
		price = ComputeUnitPrice(max)
	}

	// get instruction data
	data, err := price.Data()
	if err != nil {
		return nil, 0, err
	}

	// build tx
	txWithFee.Message.Instructions = append([]solana.CompiledInstruction{{
		ProgramIDIndex: programIdx,
		Data:           data,
	}}, txWithFee.Message.Instructions...)

	// track current fee by passing it out ad using it in PendingTxs.Add
	return &txWithFee, uint64(price), nil
}

type PendingTxs interface {
	New(tx PendingTx) uuid.UUID                                   // save pendingTx
	Add(id uuid.UUID, sig solana.Signature, txprice uint64) error // save signature after broadcasting
	Remove(id uuid.UUID) error
	ListSignatures() []solana.Signature // get all signatures for pending txs
	ListIDs() []uuid.UUID
	GetBySignature(sig solana.Signature) (PendingTx, bool) // get tx from signature
	GetByID(id uuid.UUID) (PendingTx, bool)
	// state change hooks
	OnSuccess(sig solana.Signature) (PendingTx, error)
	OnError(sig solana.Signature, errType int) (PendingTx, error) // match err type using enum
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
	tx.id = id
	txs.idMap[id] = tx
	return id
}

func (txs *pendingTxMemory) Add(id uuid.UUID, sig solana.Signature, price uint64) error {
	if id == uuid.Nil {
		return fmt.Errorf("uuid is nil")
	}

	if sig.IsZero() {
		return fmt.Errorf("signature is zero")
	}

	checkExists := func() error {
		if _, exists := txs.idMap[id]; !exists {
			return fmt.Errorf("ID does not exist: %s", id)
		}
		if _, exists := txs.sigMap[sig]; exists {
			return fmt.Errorf("signature already exists: %s", sig)
		}
		return nil
	}

	// check exists
	txs.lock.RLock()
	if err := checkExists(); err != nil {
		txs.lock.RUnlock()
		return err
	}
	txs.lock.RUnlock()

	// upgrade to write lock
	txs.lock.Lock()
	defer txs.lock.Unlock()

	// redo check within write lock
	if err := checkExists(); err != nil {
		return err
	}

	// save signatures (within write lock to ensure no other signatures were added)
	tx := txs.idMap[id]
	tx.signatures = append(tx.signatures, sig)
	tx.broadcast = true
	tx.timestamp = time.Now()
	tx.currentFee = price
	txs.idMap[id] = tx
	txs.sigMap[sig] = id
	return nil
}

func (txs *pendingTxMemory) Remove(id uuid.UUID) error {
	// check if already removed
	txs.lock.RLock()
	if _, exists := txs.idMap[id]; !exists {
		txs.lock.RUnlock()
		return fmt.Errorf("tx ID does not exist: %s", id)
	}
	txs.lock.RUnlock()

	// upgrade to write lock if ID exists
	txs.lock.Lock()
	defer txs.lock.Unlock()
	// redo check within write lock
	if _, exists := txs.idMap[id]; !exists {
		return fmt.Errorf("tx ID does not exist: %s", id)
	}
	for _, s := range txs.idMap[id].signatures {
		delete(txs.sigMap, s)
	}
	delete(txs.idMap, id)
	return nil
}

func (txs *pendingTxMemory) ListSignatures() []solana.Signature {
	txs.lock.RLock()
	defer txs.lock.RUnlock()
	return maps.Keys(txs.sigMap)
}

func (txs *pendingTxMemory) ListIDs() []uuid.UUID {
	txs.lock.RLock()
	defer txs.lock.RUnlock()
	return maps.Keys(txs.idMap)
}

func (txs *pendingTxMemory) GetBySignature(sig solana.Signature) (PendingTx, bool) {
	txs.lock.RLock()
	id, exists := txs.sigMap[sig]
	txs.lock.RUnlock()

	if !exists {
		return PendingTx{}, false
	}

	return txs.GetByID(id)
}

func (txs *pendingTxMemory) GetByID(id uuid.UUID) (PendingTx, bool) {
	txs.lock.RLock()
	defer txs.lock.RUnlock()

	tx, exists := txs.idMap[id]
	return tx, exists
}

func (txs *pendingTxMemory) OnSuccess(sig solana.Signature) (PendingTx, error) {
	if tx, exists := txs.GetBySignature(sig); exists {
		return tx, txs.Remove(tx.id)
	}
	return PendingTx{}, fmt.Errorf("tx signature does not exist: %s", sig)
}

func (txs *pendingTxMemory) OnError(sig solana.Signature, _ int) (PendingTx, error) {
	if tx, exists := txs.GetBySignature(sig); exists {
		return tx, txs.Remove(tx.id)
	}
	return PendingTx{}, fmt.Errorf("tx signature does not exist: %s", sig)
}

var _ PendingTxs = &pendingTxMemoryWithProm{}

type pendingTxMemoryWithProm struct {
	pendingTx *pendingTxMemory
	chainID   string
}

const (
	TxFailRevert       = iota // execution revert
	TxRPCReject               // rpc rejected transaction
	TxInvalidBlockhash        // tx used an invalid blockhash
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

func (txs *pendingTxMemoryWithProm) Add(id uuid.UUID, sig solana.Signature, price uint64) error {
	return txs.pendingTx.Add(id, sig, price)
}

func (txs *pendingTxMemoryWithProm) Remove(id uuid.UUID) error {
	return txs.pendingTx.Remove(id)
}

func (txs *pendingTxMemoryWithProm) ListSignatures() []solana.Signature {
	sigs := txs.pendingTx.ListSignatures()
	promSolTxmPendingTxs.WithLabelValues(txs.chainID).Set(float64(len(sigs)))
	return sigs
}

func (txs *pendingTxMemoryWithProm) ListIDs() []uuid.UUID {
	return txs.pendingTx.ListIDs()
}

func (txs *pendingTxMemoryWithProm) GetBySignature(sig solana.Signature) (PendingTx, bool) {
	return txs.pendingTx.GetBySignature(sig)
}

func (txs *pendingTxMemoryWithProm) GetByID(id uuid.UUID) (PendingTx, bool) {
	return txs.pendingTx.GetByID(id)
}

// Success - tx included in block and confirmed
func (txs *pendingTxMemoryWithProm) OnSuccess(sig solana.Signature) (PendingTx, error) {
	// don't increment prom metrics if tx no longer exists
	tx, err := txs.pendingTx.OnSuccess(sig)
	if err != nil {
		return tx, err
	}
	promSolTxmSuccessTxs.WithLabelValues(txs.chainID).Add(1)
	return tx, nil
}

func (txs *pendingTxMemoryWithProm) OnError(sig solana.Signature, errType int) (PendingTx, error) {
	var tx PendingTx
	var err error
	switch errType {
	case TxFailRevert:
		// don't increment prom metrics if tx no longer exists
		tx, err = txs.pendingTx.OnError(sig, errType)
		if err != nil {
			return tx, err
		}
		promSolTxmRevertTxs.WithLabelValues(txs.chainID).Add(1)
	case TxRPCReject:
		// reject called when RPC rejects a tx, no valid signature to check
		promSolTxmRejectTxs.WithLabelValues(txs.chainID).Add(1)
	case TxInvalidBlockhash:
		// invalid blockhash called when tx has invalid blockhash, no valid signature to check
		promSolTxmInvalidBlockhash.WithLabelValues(txs.chainID).Add(1)
	}
	// increment total errors
	promSolTxmErrorTxs.WithLabelValues(txs.chainID).Add(1)
	return tx, nil
}
