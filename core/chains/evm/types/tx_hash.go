package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

var _ commontypes.Hashable[TxHash] = (*TxHash)(nil)

type TxHash struct{ common.Hash }

func (a TxHash) Equals(h TxHash) bool {
	return bytes.Equal(a.Hash.Bytes(), h.Hash.Bytes())
}

func (a TxHash) Empty() bool {
	return bytes.Equal(a.Hash.Bytes(), common.Hash{}.Bytes())
}

func NewTxHash(h common.Hash) TxHash {
	return TxHash{h}
}
