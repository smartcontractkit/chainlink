package rlphelpers

import (
	"bytes"

	"github.com/ethereum/go-ethereum/rlp"
)

// RawRLPOutput is a struct that implements "raw" RLP decoding.
// In particular, it is used to decode RLP-encoded data into a tree of byte slices.
// This is supposed to emulate the rlp.decode() function in the rlp npm package.
// See https://www.npmjs.com/package/@ethereumjs/rlp for details.
// See https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp/ for more information on RLP serialization.
type RawRLPOutput struct {
	// Data is set for a terminal node, i.e a node with no children.
	Data []byte
	// Children is set for a non-terminal node, i.e a node with children.
	Children []*RawRLPOutput
}

// NewRLPBuffers creates a new RLPBuffers struct that is ready to use.
func NewRLPBuffers() *RawRLPOutput {
	return &RawRLPOutput{
		Children: make([]*RawRLPOutput, 0),
	}
}

// DecodeRLP implements go-ethereum's rlp.Decoder interface.
func (r *RawRLPOutput) DecodeRLP(s *rlp.Stream) error {
	kind, _, err := s.Kind()
	if err != nil {
		return err
	}
	switch kind {
	case rlp.List:
		// recursively traverse the RLP list.
		_, err2 := s.List()
		if err2 != nil {
			return err2
		}

		for s.MoreDataInList() {
			newBuf := NewRLPBuffers()
			err2 = newBuf.DecodeRLP(s)
			if err2 != nil {
				return err2
			}
			r.Children = append(r.Children, newBuf)
		}

		if err2 = s.ListEnd(); err2 != nil {
			return err2
		}
	case rlp.Byte:
		b, err2 := s.Raw()
		if err2 != nil {
			return err2
		}
		// Don't trim here since it's a single byte.
		r.Data = b
	case rlp.String:
		b, err2 := s.Raw()
		if err2 != nil {
			return err2
		}
		// trim the first byte, which is the "type" byte.
		// this is what the rlp npm package does.
		r.Data = b[1:]
	}
	return nil
}

func (r *RawRLPOutput) Equal(o *RawRLPOutput) bool {
	if !bytes.Equal(r.Data, o.Data) || len(r.Children) != len(o.Children) {
		return false
	}
	for i := range r.Children {
		if !r.Children[i].Equal(o.Children[i]) {
			return false
		}
	}
	return true
}
