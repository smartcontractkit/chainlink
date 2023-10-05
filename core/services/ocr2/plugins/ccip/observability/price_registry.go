package observability

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
)

type ObservedPriceRegistryV1_0_0 struct {
	*price_registry.PriceRegistry
	metric metricDetails
}

func NewObservedPriceRegistryV1_0_0(address common.Address, pluginName string, client client.Client) (*ObservedPriceRegistryV1_0_0, error) {
	priceRegistry, err := price_registry.NewPriceRegistry(address, client)
	if err != nil {
		return nil, err
	}

	return &ObservedPriceRegistryV1_0_0{
		PriceRegistry: priceRegistry,
		metric: metricDetails{
			histogram:  priceRegistryHistogram,
			pluginName: pluginName,
			chainId:    client.ConfiguredChainID(),
		},
	}, nil
}

func (o *ObservedPriceRegistryV1_0_0) GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error) {
	return withObservedContract(o.metric, "GetFeeTokens", func() ([]common.Address, error) {
		return o.PriceRegistry.GetFeeTokens(opts)
	})
}

func (o *ObservedPriceRegistryV1_0_0) GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]price_registry.InternalTimestampedPackedUint224, error) {
	return withObservedContract(o.metric, "GetTokenPrices", func() ([]price_registry.InternalTimestampedPackedUint224, error) {
		return o.PriceRegistry.GetTokenPrices(opts, tokens)
	})
}

func (o *ObservedPriceRegistryV1_0_0) ParseUsdPerUnitGasUpdated(log types.Log) (*price_registry.PriceRegistryUsdPerUnitGasUpdated, error) {
	return withObservedContract(o.metric, "ParseUsdPerUnitGasUpdated", func() (*price_registry.PriceRegistryUsdPerUnitGasUpdated, error) {
		return o.PriceRegistry.ParseUsdPerUnitGasUpdated(log)
	})
}

func (o *ObservedPriceRegistryV1_0_0) ParseUsdPerTokenUpdated(log types.Log) (*price_registry.PriceRegistryUsdPerTokenUpdated, error) {
	return withObservedContract(o.metric, "ParseUsdPerTokenUpdated", func() (*price_registry.PriceRegistryUsdPerTokenUpdated, error) {
		return o.PriceRegistry.ParseUsdPerTokenUpdated(log)
	})
}
