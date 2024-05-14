package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type SequenceTracker[
	// Represents an account address, in native chain format.
	ADDR types.Hashable,
	// Represents the sequence type for a chain. For example, nonce for EVM.
	SEQ types.Sequence,
] interface {
	// Load the next sequence needed for transactions for all enabled addresses
	LoadNextSequences(context.Context, []ADDR)
	// Get the next sequence to assign to a transaction
	GetNextSequence(context.Context, ADDR) (SEQ, error)
	// Signals the existing sequence has been used so generates and stores the next sequence
	// Can be a no-op depending on the chain
	GenerateNextSequence(ADDR, SEQ)
	// Syncs the local sequence with the one on-chain in case the address as been used externally
	// Can be a no-op depending on the chain
	SyncSequence(context.Context, ADDR, services.StopChan)
}
