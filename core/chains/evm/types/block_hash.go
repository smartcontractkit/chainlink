package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

var _ commontypes.Hashable[*BlockHash] = (*BlockHash)(nil)

type BlockHash struct{ common.Hash }

func (a *BlockHash) Equals(h *BlockHash) bool {
	return bytes.Equal(a.Hash.Bytes(), h.Hash.Bytes())
}

func (a *BlockHash) Empty() bool {
	return a == nil || bytes.Equal(a.Hash.Bytes(), common.Hash{}.Bytes())
}

func NewBlockHash(h common.Hash) *BlockHash {
	return &BlockHash{h}
}
