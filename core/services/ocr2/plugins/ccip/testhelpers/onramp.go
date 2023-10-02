package testhelpers

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	FakeOnRampVersion = "1.2.0"
)

type FakeOnRamp struct {
	*mock_contracts.EVM2EVMOnRampInterface

	dynamicConfig evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig

	mu sync.RWMutex
}

func NewFakeOnRamp(t *testing.T) (*FakeOnRamp, common.Address) {
	addr := utils.RandomAddress()
	mockOnRamp := mock_contracts.NewEVM2EVMOnRampInterface(t)
	mockOnRamp.On("Address").Return(addr).Maybe()

	onRamp := &FakeOnRamp{EVM2EVMOnRampInterface: mockOnRamp}
	return onRamp, addr
}

func (o *FakeOnRamp) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	return fmt.Sprintf("%s %s", ccipconfig.EVM2EVMOnRamp, FakeOnRampVersion), nil
}

func (o *FakeOnRamp) SetDynamicCfg(cfg evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig) {
	setOnRampVal(o, func(o *FakeOnRamp) { o.dynamicConfig = cfg })
}

func (o *FakeOnRamp) GetDynamicConfig(opts *bind.CallOpts) (evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig, error) {
	return getOnRampVal(o, func(o *FakeOnRamp) (evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig, error) { return o.dynamicConfig, nil })
}

func getOnRampVal[T any](o *FakeOnRamp, getter func(o *FakeOnRamp) (T, error)) (T, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return getter(o)
}

func setOnRampVal(o *FakeOnRamp, setter func(o *FakeOnRamp)) {
	o.mu.Lock()
	defer o.mu.Unlock()
	setter(o)
}
