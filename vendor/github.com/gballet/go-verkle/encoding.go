// This is free and unencumbered software released into the public domain.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// For more information, please refer to <https://unlicense.org>

package verkle

import (
	"errors"
	"fmt"

	"github.com/crate-crypto/go-ipa/banderwagon"
)

var (
	ErrInvalidNodeEncoding = errors.New("invalid node encoding")

	mask = [8]byte{0x80, 0x40, 0x20, 0x10, 0x8, 0x4, 0x2, 0x1}
)

const (
	nodeTypeSize = 1
	bitlistSize  = NodeWidth / 8

	// Shared between internal and leaf nodes.
	nodeTypeOffset = 0

	// Internal nodes offsets.
	internalBitlistOffset    = nodeTypeOffset + nodeTypeSize
	internalCommitmentOffset = internalBitlistOffset + bitlistSize

	// Leaf node offsets.
	leafSteamOffset        = nodeTypeOffset + nodeTypeSize
	leafBitlistOffset      = leafSteamOffset + StemSize
	leafCommitmentOffset   = leafBitlistOffset + bitlistSize
	leafC1CommitmentOffset = leafCommitmentOffset + banderwagon.UncompressedSize
	leafC2CommitmentOffset = leafC1CommitmentOffset + banderwagon.UncompressedSize
	leafChildrenOffset     = leafC2CommitmentOffset + banderwagon.UncompressedSize
)

func bit(bitlist []byte, nr int) bool {
	if len(bitlist)*8 <= nr {
		return false
	}
	return bitlist[nr/8]&mask[nr%8] != 0
}

var errSerializedPayloadTooShort = errors.New("verkle payload is too short")

// ParseNode deserializes a node into its proper VerkleNode instance.
// The serialized bytes have the format:
// - Internal nodes: <nodeType><bitlist><commitment>
// - Leaf nodes:     <nodeType><stem><bitlist><comm><c1comm><c2comm><children...>
func ParseNode(serializedNode []byte, depth byte) (VerkleNode, error) {
	// Check that the length of the serialized node is at least the smallest possible serialized node.
	if len(serializedNode) < nodeTypeSize+banderwagon.UncompressedSize {
		return nil, errSerializedPayloadTooShort
	}

	switch serializedNode[0] {
	case leafRLPType:
		return parseLeafNode(serializedNode, depth)
	case internalRLPType:
		return CreateInternalNode(serializedNode[internalBitlistOffset:internalCommitmentOffset], serializedNode[internalCommitmentOffset:], depth)
	default:
		return nil, ErrInvalidNodeEncoding
	}
}

func parseLeafNode(serialized []byte, depth byte) (VerkleNode, error) {
	bitlist := serialized[leafBitlistOffset : leafBitlistOffset+bitlistSize]
	var values [NodeWidth][]byte
	offset := leafChildrenOffset
	for i := 0; i < NodeWidth; i++ {
		if bit(bitlist, i) {
			if offset+LeafValueSize > len(serialized) {
				return nil, fmt.Errorf("verkle payload is too short, need at least %d and only have %d, payload = %x (%w)", offset+32, len(serialized), serialized, errSerializedPayloadTooShort)
			}
			values[i] = serialized[offset : offset+LeafValueSize]
			offset += LeafValueSize
		}
	}
	ln := NewLeafNodeWithNoComms(serialized[leafSteamOffset:leafSteamOffset+StemSize], values[:])
	ln.setDepth(depth)
	ln.c1 = new(Point)

	// Sanity check that we have at least 3*banderwagon.UncompressedSize bytes left in the serialized payload.
	if len(serialized[leafCommitmentOffset:]) < 3*banderwagon.UncompressedSize {
		return nil, fmt.Errorf("leaf node commitments are not the correct size, expected at least %d, got %d", 3*banderwagon.UncompressedSize, len(serialized[leafC1CommitmentOffset:]))
	}

	if err := ln.c1.SetBytesUncompressed(serialized[leafC1CommitmentOffset:leafC1CommitmentOffset+banderwagon.UncompressedSize], true); err != nil {
		return nil, fmt.Errorf("setting c1 commitment: %w", err)
	}
	ln.c2 = new(Point)
	if err := ln.c2.SetBytesUncompressed(serialized[leafC2CommitmentOffset:leafC2CommitmentOffset+banderwagon.UncompressedSize], true); err != nil {
		return nil, fmt.Errorf("setting c2 commitment: %w", err)
	}
	ln.commitment = new(Point)
	if err := ln.commitment.SetBytesUncompressed(serialized[leafCommitmentOffset:leafC1CommitmentOffset], true); err != nil {
		return nil, fmt.Errorf("setting commitment: %w", err)
	}
	return ln, nil
}

func CreateInternalNode(bitlist []byte, raw []byte, depth byte) (*InternalNode, error) {
	// GetTreeConfig caches computation result, hence
	// this op has low overhead
	node := new(InternalNode)

	if len(bitlist) != bitlistSize {
		return nil, ErrInvalidNodeEncoding
	}

	// Create a HashNode placeholder for all values
	// corresponding to a set bit.
	node.children = make([]VerkleNode, NodeWidth)
	for i, b := range bitlist {
		for j := 0; j < 8; j++ {
			if b&mask[j] != 0 {
				node.children[8*i+j] = HashedNode{}
			} else {

				node.children[8*i+j] = Empty(struct{}{})
			}
		}
	}
	node.depth = depth
	if len(raw) != banderwagon.UncompressedSize {
		return nil, ErrInvalidNodeEncoding
	}

	node.commitment = new(Point)
	if err := node.commitment.SetBytesUncompressed(raw, true); err != nil {
		return nil, fmt.Errorf("setting commitment: %w", err)
	}
	return node, nil
}
