package testhelpers

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FakePriceRegistry struct {
	*mock_contracts.PriceRegistryInterface

	tokenPrices []price_registry.InternalTimestampedPackedUint224
	feeTokens   []common.Address

	mu sync.RWMutex
}

func NewFakePriceRegistry(t *testing.T) (*FakePriceRegistry, common.Address) {
	addr := utils.RandomAddress()
	mockPriceRegistry := mock_contracts.NewPriceRegistryInterface(t)
	mockPriceRegistry.On("Address").Return(addr).Maybe()

	priceRegistry := &FakePriceRegistry{PriceRegistryInterface: mockPriceRegistry}
	return priceRegistry, addr
}

func (p *FakePriceRegistry) SetTokenPrices(prices []price_registry.InternalTimestampedPackedUint224) {
	setPriceRegistryVal(p, func(p *FakePriceRegistry) { p.tokenPrices = prices })
}

func (p *FakePriceRegistry) GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]price_registry.InternalTimestampedPackedUint224, error) {
	return getPriceRegistryVal(p, func(p *FakePriceRegistry) ([]price_registry.InternalTimestampedPackedUint224, error) {
		return p.tokenPrices, nil
	})
}

func (p *FakePriceRegistry) SetFeeTokens(tokens []common.Address) {
	setPriceRegistryVal(p, func(p *FakePriceRegistry) { p.feeTokens = tokens })
}

func (p *FakePriceRegistry) GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error) {
	return getPriceRegistryVal(p, func(p *FakePriceRegistry) ([]common.Address, error) { return p.feeTokens, nil })
}

func getPriceRegistryVal[T any](p *FakePriceRegistry, getter func(p *FakePriceRegistry) (T, error)) (T, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return getter(p)
}

func setPriceRegistryVal(p *FakePriceRegistry, setter func(p *FakePriceRegistry)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	setter(p)
}
