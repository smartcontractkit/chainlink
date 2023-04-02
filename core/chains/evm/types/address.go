package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

var _ commontypes.Hashable = &Address{}

type Address struct {
	commontypes.Hashable
	nativeAddress common.Address
}

func (a *Address) MarshalText() (text []byte, err error) {
	return a.nativeAddress.Bytes(), nil
}

func (a *Address) UnmarshalText(text []byte) error {
	a.nativeAddress = common.BytesToAddress(text)
	return nil
}

func (a *Address) String() string {
	return a.nativeAddress.String()
}

func (a *Address) Equals(h commontypes.Hashable) bool {
	return bytes.Equal(a.nativeAddress.Bytes(), h.(*Address).nativeAddress.Bytes())
}

func (a *Address) IsEmpty() bool {
	return a == nil || a.nativeAddress == common.Address{}
}

func (a *Address) NativeAddress() *common.Address {
	return &a.nativeAddress
}

func NewAddress(h common.Address) *Address {
	return &Address{
		nativeAddress: h,
	}
}
