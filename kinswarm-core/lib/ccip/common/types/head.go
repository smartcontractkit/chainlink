package types

import (
	"math/big"
	"time"
)

// Head provides access to a chain's head, as needed by the TxManager.
// This is a generic interface which ALL chains will implement.
type Head[BLOCK_HASH Hashable] interface {
	// BlockNumber is the head's block number
	BlockNumber() int64

	// Timestamp the time of mining of the block
	GetTimestamp() time.Time

	// ChainLength returns the length of the chain followed by recursively looking up parents
	ChainLength() uint32

	// EarliestHeadInChain traverses through parents until it finds the earliest one
	EarliestHeadInChain() Head[BLOCK_HASH]

	// Parent is the head's parent block
	GetParent() Head[BLOCK_HASH]

	// Hash is the head's block hash
	BlockHash() BLOCK_HASH
	GetParentHash() BLOCK_HASH

	// HashAtHeight returns the hash of the block at the given height, if it is in the chain.
	// If not in chain, returns the zero hash
	HashAtHeight(blockNum int64) BLOCK_HASH

	// HeadAtHeight returns head at specified height or an error, if one does not exist in provided chain.
	HeadAtHeight(blockNum int64) (Head[BLOCK_HASH], error)

	// Returns the total difficulty of the block. For chains who do not have a concept of block
	// difficulty, return 0.
	BlockDifficulty() *big.Int
	// IsValid returns true if the head is valid.
	IsValid() bool

	// Returns the latest finalized based on finality tag or depth
	LatestFinalizedHead() Head[BLOCK_HASH]
}
