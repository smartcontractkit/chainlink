// Code generated — DO NOT EDIT.

//go:build !wasip1

package capabilities_registry

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmmock "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/mock"
)

var (
	_ = errors.New
	_ = fmt.Errorf
	_ = big.NewInt
	_ = common.Big1
)

// CapabilitiesRegistryMock is a mock implementation of CapabilitiesRegistry for testing.
type CapabilitiesRegistryMock struct {
	TypeAndVersion func() (string, error)
}

// NewCapabilitiesRegistryMock creates a new CapabilitiesRegistryMock for testing.
func NewCapabilitiesRegistryMock(address common.Address, clientMock *evmmock.ClientCapability) *CapabilitiesRegistryMock {
	mock := &CapabilitiesRegistryMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["typeAndVersion"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TypeAndVersion == nil {
				return nil, errors.New("typeAndVersion method not mocked")
			}
			result, err := mock.TypeAndVersion()
			if err != nil {
				return nil, err
			}
			return abi.Methods["typeAndVersion"].Outputs.Pack(result)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
