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
	return a.nativeHash.Bytes(), nil
}

func (a *TxHash) UnmarshalText(text []byte) error {
	a.nativeHash = common.BytesToHash(text)
	return nil
}

func (a *TxHash) String() string {
	return a.nativeHash.String()
}

func (a *TxHash) Equals(h commontypes.Hashable) bool {
	return bytes.Equal(a.nativeHash.Bytes(), h.(*Address).nativeAddress.Bytes())
}

func (a *TxHash) Empty() bool {
	return a == nil || a.nativeHash == common.Hash{}
}

func (a *TxHash) NativeHash() *common.Hash {
	return &a.nativeHash
}

func NewTxHash(h common.Hash) *TxHash {
	return &TxHash{
		nativeHash: h,
	}
}
