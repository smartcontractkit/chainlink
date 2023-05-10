package utils

import (
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// EarliestInChain recurses through parents until it finds the earliest one
func EarliestInChain[H types.Head[HASH], HASH types.Hashable](h H) H {
	for h.GetParent() != nil {
		h = h.GetParent().(H)
	}
	return h
}

// IsInChain returns true if the given hash matches the hash of a head in the chain
func IsInChain[H types.Head[HASH], HASH types.Hashable](h H, blockHash HASH) bool {
	for {
		if h.BlockHash() == blockHash {
			return true
		}
		if h.GetParent() != nil {
			h = h.GetParent().(H)
		} else {
			break
		}
	}
	return false
}

// HashAtHeight returns the hash of the block at the given height, if it is in the chain.
// If not in chain, returns the zero hash
func HashAtHeight[H types.Head[HASH], HASH types.Hashable](h H, blockNum int64) interface{} {
	for {
		if h.BlockNumber() == blockNum {
			return h.BlockHash()
		}
		if h.GetParent() != nil {
			h = h.GetParent().(H)
		} else {
			break
		}
	}
	return nil
}

// ChainLength returns the length of the chain followed by recursively looking up parents
func ChainLength[H types.Head[HASH], HASH types.Hashable](h H) uint32 {
	if h.Equals(nil) {
		return 0
	}
	l := uint32(1)

	for {
		parent := h.GetParent()
		if parent != nil {
			l++
			if h.Equals(parent) {
				panic("circular reference detected")
			}
			h = parent.(H)
		} else {
			break
		}
	}
	return l
}

// ChainHashes returns an array of block hashes by recursively looking up parents
func ChainHashes[H types.Head[HASH], HASH types.Hashable](h H) []HASH {
	var hashes []HASH

	for {
		hashes = append(hashes, h.BlockHash())
		if h.GetParent() != nil {
			if h.Equals(h.GetParent()) {
				panic("circular reference detected")
			}
			h = h.GetParent().(H)
		} else {
			break
		}
	}
	return hashes
}
