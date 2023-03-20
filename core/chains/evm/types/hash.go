package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/common/types"
)

type hash struct {
	commontypes.Hashable
	nativeHash common.Hash
}

func (a *hash) ToBytes() []byte {
	return a.nativeHash.Bytes()
}

func (a *hash) ToString() string {
	return a.nativeHash.String()
}

func (a *hash) FromString(str string) {
	a.nativeHash = common.HexToHash(str)
}

func (a *hash) Equals(h commontypes.Hashable) bool {
	return bytes.Equal(a.nativeHash.Bytes(), h.ToBytes())
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
