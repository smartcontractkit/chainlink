package fees

import (
	"bytes"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

// https://github.com/solana-labs/solana/blob/60858d043ca612334de300805d93ea3014e8ab37/sdk/src/compute_budget.rs#L25
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

const (
	COMPUTE_BUDGET_PROGRAM = "ComputeBudget111111111111111111111111111111"
)

// https://docs.solana.com/developing/programming-model/runtime
type ComputeUnitPrice uint64

// returns the compute budget program
func (val ComputeUnitPrice) ProgramID() solana.PublicKey {
	return solana.MustPublicKeyFromBase58(COMPUTE_BUDGET_PROGRAM)
}

// No accounts needed
func (val ComputeUnitPrice) Accounts() (accounts []*solana.AccountMeta) {
	return accounts
}

// simple encoding into program expected format
func (val ComputeUnitPrice) Data() ([]byte, error) {
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

// modifies passed in tx to set compute unit price
func SetComputeUnitPrice(tx *solana.Transaction, price ComputeUnitPrice) error {
	// find ComputeBudget program to accounts if it exists
	// reimplements HasAccount to retrieve index: https://github.com/gagliardetto/solana-go/blob/618f56666078f8131a384ab27afd918d248c08b7/message.go#L233
	var exists bool
	var programIdx uint16
	for i, a := range tx.Message.AccountKeys {
		if a.Equals(price.ProgramID()) {
			exists = true
			programIdx = uint16(i)
			break
		}
	}
	// if it doesn't exist, add to account keys
	if !exists {
		tx.Message.AccountKeys = append(tx.Message.AccountKeys, price.ProgramID())
		programIdx = uint16(len(tx.Message.AccountKeys) - 1) // last index of account keys

		// https://github.com/gagliardetto/solana-go/blob/618f56666078f8131a384ab27afd918d248c08b7/transaction.go#L293
		tx.Message.Header.NumReadonlyUnsignedAccounts++
	}

	// get instruction data
	data, err := price.Data()
	if err != nil {
		return err
	}

	// compiled instruction
	instruction := solana.CompiledInstruction{
		ProgramIDIndex: programIdx,
		Data:           data,
	}

	// check if there is an instruction for setcomputeunitprice
	var found bool
	var instructionIdx int
	for i := range tx.Message.Instructions {
		if tx.Message.Instructions[i].ProgramIDIndex == programIdx &&
			len(tx.Message.Instructions[i].Data) > 0 &&
			tx.Message.Instructions[i].Data[0] == Instruction_SetComputeUnitPrice {
			found = true
			instructionIdx = i
			break
		}
	}

	if found {
		tx.Message.Instructions[instructionIdx] = instruction
	} else {
		// build with first instruction as set compute unit price
		tx.Message.Instructions = append([]solana.CompiledInstruction{instruction}, tx.Message.Instructions...)
	}

	return nil
}
