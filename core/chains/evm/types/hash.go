package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

type hash struct {
	commontypes.Hashable
	nativeHash common.Hash
}

var _ commontypes.Hashable = &hash{}

func (a *hash) MarshalText() (text []byte, err error) {
	return a.nativeHash.MarshalText()
}

func (a *hash) UnmarshalText(text []byte) error {
	return a.nativeHash.UnmarshalText(text)
}

func (a *hash) String() string {
	return a.nativeHash.String()
}

func (a *hash) Equals(h commontypes.Hashable) bool {
	return bytes.Equal(a.nativeHash.Bytes(), h.(*Address).nativeAddress.Bytes())
}

func (a *hash) IsEmpty() bool {
	return a.nativeHash == common.Hash{}
}

func (a *hash) NativeHash() *common.Hash {
	return &a.nativeHash
}

type BlockHash = hash
type TxHash = hash

func NewBlockHash(h common.Hash) *BlockHash {
	return &hash{
		nativeHash: h,
	}
}
func NewTxHash(h common.Hash) *TxHash {
	return &hash{
		nativeHash: h,
	}
}
