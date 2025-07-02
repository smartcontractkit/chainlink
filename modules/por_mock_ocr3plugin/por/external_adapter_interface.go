package por

import (
	"context"
)

// A unique identifier for a chain, e.g., Chain.Selector from the chain-selectors package
// See https://github.com/smartcontractkit/chain-selectors
type ChainSelector uint64

// The External Adapter (EA) is the point of contact with PoR related events.
// An EA is expected to correspond to a single token, and track its corresponding chains.
// Its main purpose is to calculate the minting allowance per chain based on
// (1) on-chain (pre-)minting requests,
// (2) global supply across all chains, and
// (3) off-chain reserve information
//
// All methods are expected to:
// - be deterministic (given a *fixed* state of the EA) & thread-safe
// - finish within a reasonable time frame (max 500 ms)
// - return information pertaining to what is (expected to be) the *finalized* state of the chain(s).
// Determining what is considered to be the finalized state is up to the adapter implementation
type ExternalAdapter interface {
	// GetPayload returns the payload (ExternalAdapterPayload) for the queried blocks.
	//
	// The payload contains:
	// (1) the per-chain mintable information (BlockMintablePair) computed given the query (`blocks`) and the latest reserve information,
	// (2) the latest reserve information (ReserveInfo),
	// (3) the latest blocks for each chain (LatestBlocks).
	//
	// Mintables is either:
	// - a map with the mintable amount for each (chain, block) in `blocks`, or
	// - `nil` if the information for a chain or its block is not available. (If the information for a block is not available,
	// the correct mintable amount cannot be determined.)
	//
	// ReserveInfo is the same used to calculate the mintable amounts.
	//
	// LatestBlocks are the latest blocks for each chain. Specifically:
	// Specifically, given the on-chain events to PoR (namely, premints, mints, and burns):
	// - LatestBlocks includes, per-chain, the latest blockNumber known by the EA for that chain.
	// - LatestBlocks is monotonically non-decreasing in every chain across repeated calls to GetPayload.
	// - If the `blocks` argument contains a chain which the EA does not track (yet), LatestBlocks should
	// include that chain with block number of 0.
	//
	// Notes on usage:
	// - The plugin will call GetPayload periodically and compare the new generated mintables with the previous ones to determine
	//   if a new report should be generated. By avoiding changing the block numbers when there are no new PoR-related events,
	//   the plugin can avoid generating unnecessary reports.
	GetPayload(ctx context.Context, blocks Blocks) (ExternalAdapterPayload, error)
}
