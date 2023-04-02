package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

type BlockHash struct {
	commontypes.Hashable
	nativeHash common.Hash
}

var _ commontypes.Hashable = &BlockHash{}

func (a *BlockHash) MarshalText() (text []byte, err error) {
	return a.nativeHash.MarshalText()
}

func (a *BlockHash) UnmarshalText(text []byte) error {
	return a.nativeHash.UnmarshalText(text)
}

func (a *BlockHash) String() string {
	return a.nativeHash.String()
}

func (a *BlockHash) Equals(h commontypes.Hashable) bool {
	return bytes.Equal(a.nativeHash.Bytes(), h.(*Address).nativeAddress.Bytes())
}

func (a *BlockHash) IsEmpty() bool {
	return a.nativeHash == common.Hash{}
}

func (a *BlockHash) NativeHash() *common.Hash {
	return &a.nativeHash
}

func NewBlockHash(h common.Hash) *BlockHash {
	return &BlockHash{
		nativeHash: h,
	}
}
