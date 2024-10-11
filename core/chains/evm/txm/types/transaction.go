package types

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
)

type TxState string

const (
	TxUnstarted   = TxState("unstarted")
	TxUnconfirmed = TxState("unconfirmed")
	TxConfirmed   = TxState("confirmed")

	TxFatalError = TxState("fatal")
	TxFinalized  = TxState("finalized")
)

type Transaction struct {
	ID                uint64
	IdempotencyKey    *string
	ChainID           *big.Int
	Nonce             uint64
	FromAddress       common.Address
	ToAddress         common.Address
	Value             *big.Int
	Data              []byte
	SpecifiedGasLimit uint64

	CreatedAt       time.Time
	LastBroadcastAt time.Time

	State        TxState
	IsPurgeable  bool
	Attempts     []*Attempt
	AttemptCount uint16 // AttempCount is strictly kept inMemory and prevents indefinite retrying
	// Meta, ForwarderAddress, Strategy
}

func (t *Transaction) FindAttemptByHash(attemptHash common.Hash) (*Attempt, error) {
	for _, a := range t.Attempts {
		if a.Hash == attemptHash {
			return a, nil
		}
	}
	return nil, fmt.Errorf("attempt with hash: %v was not found", attemptHash)
}

func (t *Transaction) DeepCopy() *Transaction {
	copy := *t
	var attemptsCopy []*Attempt
	for _, attempt := range t.Attempts {
		attemptsCopy = append(attemptsCopy, attempt.DeepCopy())
	}
	copy.Attempts = attemptsCopy
	return &copy
}

type Attempt struct {
	ID                uint64
	TxID              uint64
	Hash              common.Hash
	Fee               gas.EvmFee
	GasLimit          uint64
	Type              byte
	SignedTransaction *types.Transaction

	CreatedAt   time.Time
	BroadcastAt time.Time
}

func (a *Attempt) DeepCopy() *Attempt {
	copy := *a
	if a.SignedTransaction != nil {
		signedTransactionCopy := *a.SignedTransaction
		copy.SignedTransaction = &signedTransactionCopy
	}
	return &copy
}
