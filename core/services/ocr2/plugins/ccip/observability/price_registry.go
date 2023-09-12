package observability

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
)

type ObservedPriceRegistry struct {
	price_registry.PriceRegistryInterface
	metric metricDetails
}

func NewObservedPriceRegistry(address common.Address, pluginName string, client client.Client) (price_registry.PriceRegistryInterface, error) {
	priceRegistry, err := price_registry.NewPriceRegistry(address, client)
	if err != nil {
		return nil, err
	}

	return &ObservedPriceRegistry{
		PriceRegistryInterface: priceRegistry,
		metric: metricDetails{
			histogram:  priceRegistryHistogram,
			pluginName: pluginName,
			chainId:    client.ConfiguredChainID(),
		},
	}, nil
}

func (o *ObservedPriceRegistry) GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error) {
	return withObservedContract(o.metric, "GetFeeTokens", func() ([]common.Address, error) {
		return o.PriceRegistryInterface.GetFeeTokens(opts)
	})
}

func (o *ObservedPriceRegistry) GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]price_registry.InternalTimestampedPackedUint224, error) {
	return withObservedContract(o.metric, "GetTokenPrices", func() ([]price_registry.InternalTimestampedPackedUint224, error) {
		return o.PriceRegistryInterface.GetTokenPrices(opts, tokens)
	})
}

func (o *ObservedPriceRegistry) ParseUsdPerUnitGasUpdated(log types.Log) (*price_registry.PriceRegistryUsdPerUnitGasUpdated, error) {
	return withObservedContract(o.metric, "ParseUsdPerUnitGasUpdated", func() (*price_registry.PriceRegistryUsdPerUnitGasUpdated, error) {
		return o.PriceRegistryInterface.ParseUsdPerUnitGasUpdated(log)
	})
}

func (o *ObservedPriceRegistry) ParseUsdPerTokenUpdated(log types.Log) (*price_registry.PriceRegistryUsdPerTokenUpdated, error) {
	return withObservedContract(o.metric, "ParseUsdPerTokenUpdated", func() (*price_registry.PriceRegistryUsdPerTokenUpdated, error) {
		return o.PriceRegistryInterface.ParseUsdPerTokenUpdated(log)
	})
}
