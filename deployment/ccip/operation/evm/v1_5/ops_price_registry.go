package v1_5

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type PriceRegistryGetAllFeeTokensIn struct {
	Address       common.Address
	ChainSelector uint64
}

var (
	PriceRegistryGetAllFeeTokensOps = operations.NewOperation(
		"GetAllFeeTokensOps",
		semver.MustParse("1.0.0"),
		"Gets the FeeTokens for a price Registry",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input PriceRegistryGetAllFeeTokensIn) ([]common.Address, error) {
			priceRegistry, err := price_registry.NewPriceRegistry(input.Address, deps.Chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create PriceRegistry contract binding on source chain %d: %w", deps.Chain.Selector, err)
			}

			allFeeTokens, err2 := priceRegistry.GetFeeTokens(nil)
			if err2 != nil {
				return nil, fmt.Errorf("failed to all tokens on PriceRegistry %s for  source chain %d: %w", input.Address.Hex(), deps.Chain.Selector, err2)
			}

			return allFeeTokens, nil
		})
)
