package soltxm

import (
	"bytes"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

// https://github.com/solana-labs/solana/blob/master/sdk/src/compute_budget.rs
const (
	// deprecated: will not support for building instruction
	Instruction_RequestUnitsDeprecated uint8 = iota

	// Request a specific transaction-wide program heap region size in bytes.
	// The value requested must be a multiple of 1024. This new heap region
	// size applies to each program executed in the transaction, including all
	// calls to CPIs.
	// note: uses ag_binary.Varuint32
	Instruction_RequestHeapFrame

	// Set a specific compute unit limit that the transaction is allowed to consume.
	// note: uses ag_binary.Varuint32
	Instruction_SetComputeUnitLimit

	// Set a compute unit price in "micro-lamports" to pay a higher transaction
	// fee for higher transaction prioritization.
	// note: uses ag_binary.Uint64
	Instruction_SetComputeUnitPrice
)

// https://docs.solana.com/developing/programming-model/runtime
type SetComputeUnitPrice uint64

// returns the compute budget program
func (val SetComputeUnitPrice) ProgramID() solana.PublicKey {
	return solana.MustPublicKeyFromBase58("ComputeBudget111111111111111111111111111111") //TODO: temporary to avoid upgrading solana-go
}

// No accoutns needed
func (val SetComputeUnitPrice) Accounts() (accounts []*solana.AccountMeta) {
	return accounts
}

// simple encoding into program expected format
func (val SetComputeUnitPrice) Data() ([]byte, error) {
	buf := new(bytes.Buffer)

	// encode method identifier
	if err := buf.WriteByte(Instruction_SetComputeUnitPrice); err != nil {
		return []byte{}, err
	}

	// encode value
	if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
