package types

// Head provides access to a chain's head, as needed by the TxManager.
// This is a generic interface which ALL chains will implement.
//
//go:generate mockery --quiet --name Head --output ./mocks/ --case=underscore
type Head[BLOCK_HASH Hashable] interface {
	// BlockNumber is the head's block number
	BlockNumber() int64

	// ChainLength returns the length of the chain followed by recursively looking up parents
	ChainLength() uint32

	// EarliestInChain traverses through parents until it finds the earliest one
	EarliestHeadInChain() Head[BLOCK_HASH]

	// Hash is the head's block hash
	BlockHash() BLOCK_HASH

	// ParentHash is the head's parent block hash
	GetParentHash() BLOCK_HASH

	// HashAtHeight returns the hash of the block at the given height, if it is in the chain.
	// If not in chain, returns the zero hash
	HashAtHeight(blockNum int64) BLOCK_HASH

	// GetParent returns the parent head
	GetParent() Head[BLOCK_HASH]
}

type HeadTrackerHead[H Head[BLOCK_HASH], BLOCK_HASH Hashable] interface {
	Head[BLOCK_HASH]

	// SetParent sets the parent head
	SetParent(H)
}
