package types

import "github.com/ethereum/go-ethereum/common"

// Provides a minimal access to a chain's head, as needed by the TxManager.
// This is a generic interface whcih ALL chains will implement.
type HeadView interface {
	// The head's block number
	BlockNumber() int64

	// ChainLength returns the length of the chain followed by recursively looking up parents
	ChainLength() uint32

	// EarliestInChain recurses through parents until it finds the earliest one
	EarliestInChain() *HeadView

	Hash() common.Hash

	Parent() *HeadView

	// HashAtHeight returns the hash of the block at the given height, if it is in the chain.
	// If not in chain, returns the zero hash
	HashAtHeight(blockNum int64) common.Hash
}
