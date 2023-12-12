package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type Bindings map[string]methodBindings

func (b Bindings) addEvent(contractName, typeName string, evt common.Hash) error {
	ae, err := b.getBinding(contractName, typeName)
	if err != nil {
		return err
	}

	ae.evt = &evt
	return nil
}

func (b Bindings) getBinding(contractName, typeName string) (*addrEvtBinding, error) {
	typeNames, ok := b[contractName]
	if !ok {
		return nil, fmt.Errorf("%w: contract %s not found", types.ErrInvalidConfig, contractName)
	}

	ae, ok := typeNames[typeName]
	if !ok {
		return nil, fmt.Errorf("%w: method %s not found in contract %s", types.ErrInvalidConfig, typeName, contractName)
	}

	return ae, nil
}

type methodBindings map[string]*addrEvtBinding

func NewAddrEvtFromAddress(address common.Address) *addrEvtBinding {
	return &addrEvtBinding{addr: address}
}

type addrEvtBinding struct {
	addr common.Address
	evt  *common.Hash
}
