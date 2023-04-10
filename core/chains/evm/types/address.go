package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

var _ commontypes.Hashable[*Address] = (*Address)(nil)

type Address struct{ common.Address }

func (a *Address) Equals(h *Address) bool {
	return bytes.Equal(a.Bytes(), h.Bytes())
}

func (a *Address) Empty() bool {
	return a == nil || bytes.Equal(a.Bytes(), common.Address{}.Bytes())
}

func NewAddress(h common.Address) *Address {
	return &Address{h}
}

func HexToAddress(hexStr string) *Address {
	return &Address{common.HexToAddress(hexStr)}
}
