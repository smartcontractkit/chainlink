package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

type TxHash struct {
	commontypes.Hashable
	nativeHash common.Hash
}

var _ commontypes.Hashable = &TxHash{}

func (a *TxHash) MarshalText() (text []byte, err error) {
	return a.nativeHash.MarshalText()
}

func (a *TxHash) UnmarshalText(text []byte) error {
	return a.nativeHash.UnmarshalText(text)
}

func (a *TxHash) String() string {
	return a.nativeHash.String()
}

func (a *TxHash) Equals(h commontypes.Hashable) bool {
	return bytes.Equal(a.nativeHash.Bytes(), h.(*Address).nativeAddress.Bytes())
}

func (a *TxHash) IsEmpty() bool {
	return a.nativeHash == common.Hash{}
}

func (a *TxHash) NativeHash() *common.Hash {
	return &a.nativeHash
}

func NewTxHash(h common.Hash) *TxHash {
	return &TxHash{
		nativeHash: h,
	}
}
