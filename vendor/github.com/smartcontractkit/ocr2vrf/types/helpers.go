package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.dedis.ch/kyber/v3"
)

func (b Block) VRFHash(separator common.Hash, pKey kyber.Point) common.Hash {
	var heightBytes [8]byte
	binary.BigEndian.PutUint64(heightBytes[:], b.Height)
	var delayBytes [4]byte
	binary.BigEndian.PutUint32(delayBytes[:], b.ConfirmationDelay)
	key, err := pKey.MarshalBinary()
	if err != nil {
		panic("could not serialize key as domain separator")
	}
	hashMsg := bytes.Join(
		[][]byte{separator[:], key, delayBytes[:], heightBytes[:], b.Hash.Bytes()},
		nil,
	)
	return common.BytesToHash(crypto.Keccak256(hashMsg))
}

func (b Block) String() string {
	return fmt.Sprintf(
		"Block{Height: %d, ConfirmationDelay: %d, Hash: 0x%x}",
		b.Height, b.ConfirmationDelay, b.Hash,
	)
}

type Blocks []Block

var _ sort.Interface = (*Blocks)(nil)

func (b Blocks) Len() int { return len(b) }
func (b Blocks) Less(i, j int) bool {
	if b[i].Height < b[j].Height {
		return true
	}
	if b[i].Height > b[j].Height {
		return false
	}
	if b[i].ConfirmationDelay < b[j].ConfirmationDelay {
		return true
	}
	if b[i].ConfirmationDelay > b[j].ConfirmationDelay {
		return false
	}
	return b[i].Hash.Hex() < b[j].Hash.Hex()
}

func (b Blocks) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
