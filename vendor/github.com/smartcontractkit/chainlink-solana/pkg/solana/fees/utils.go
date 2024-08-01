package fees

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// returns new fee based on number of times bumped
func CalculateFee(base, max, min uint64, count uint) uint64 {
	amount := base

	for i := uint(0); i < count; i++ {
		if base == 0 && i == 0 {
			amount = 1
		} else {
			next := amount + amount
			if next <= amount {
				// overflowed
				amount = max
				break
			}
			amount = next
		}
	}

	// respect bounds
	if amount < min {
		return min
	}
	if amount > max {
		return max
	}
	return amount
}

type BlockData struct {
	Fees   []uint64           // total fee
	Prices []ComputeUnitPrice // price per unit
}

// ParseBlock parses the fee calculations from all the transactions within a block
func ParseBlock(res *rpc.GetBlockResult) (out BlockData, err error) {
	if res == nil {
		return out, fmt.Errorf("GetBlockResult was nil")
	}

	for _, tx := range res.Transactions {
		if tx.Meta == nil {
			continue
		}
		baseTx, getTxErr := tx.GetTransaction()
		if getTxErr != nil {
			// exit on GetTransaction error
			// if this occurs, solana-go was unable to parse a transaction
			// further investigation is required to determine if there is incompatibility
			return out, fmt.Errorf("failed to GetTransaction (blockhash: %s): %w", res.Blockhash, err)
		}
		if baseTx == nil {
			continue
		}

		// filter out consensus vote transactions
		// consensus messages are included as txs within blocks
		if len(baseTx.Message.Instructions) == 1 &&
			baseTx.Message.AccountKeys[baseTx.Message.Instructions[0].ProgramIDIndex] == solana.VoteProgramID {
			continue
		}

		var price ComputeUnitPrice // default 0
		for _, instruction := range baseTx.Message.Instructions {
			// find instructions for compute budget program
			if baseTx.Message.AccountKeys[instruction.ProgramIDIndex] == solana.MustPublicKeyFromBase58(ComputeBudgetProgram) {
				parsed, parseErr := ParseComputeUnitPrice(instruction.Data)
				// if compute unit price found, break instruction loop
				// only one compute unit price tx is allowed
				// err returned if not SetComputeUnitPrice instruction
				if parseErr == nil {
					price = parsed
					break
				}
			}
		}
		out.Prices = append(out.Prices, price)
		out.Fees = append(out.Fees, tx.Meta.Fee)
	}
	return out, nil
}
