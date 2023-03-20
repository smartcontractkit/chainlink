package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/common/types"
)

type Address struct {
	commontypes.Hashable
	nativeAddress common.Address
}

func (a *Address) ToBytes() []byte {
	return a.nativeAddress.Bytes()
}

func (a *Address) ToString() string {
	return a.nativeAddress.String()
}

func (a *Address) FromString(str string) {
	a.nativeAddress = common.HexToAddress(str)
}

func (a *Address) Equals(h commontypes.Hashable) bool {
	return bytes.Equal(a.nativeAddress.Bytes(), h.ToBytes())
}

func (a *Address) NativeAddress() *common.Address {
	return &a.nativeAddress
}

func NewAddress(h common.Address) *Address {
	return &Address{
		nativeAddress: h,
	}
}
