package ccipdata

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	_ PriceRegistryReader = &PriceRegistryV1_2_0{}
)

type PriceRegistryV1_2_0 struct {
	*PriceRegistryV1_0_0
	pr *price_registry.PriceRegistry
}

func NewPriceRegistryV1_2_0(lggr logger.Logger, priceRegistryAddr common.Address, lp logpoller.LogPoller, ec client.Client) (*PriceRegistryV1_2_0, error) {
	v100, err := NewPriceRegistryV1_0_0(lggr, priceRegistryAddr, lp, ec)
	if err != nil {
		return nil, err
	}
	priceRegistry, err := price_registry.NewPriceRegistry(priceRegistryAddr, ec)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryV1_2_0{
		PriceRegistryV1_0_0: v100,
		pr:                  priceRegistry,
	}, nil
}

// GetTokenPrices must be overridden to use the 1.2 ABI (return parameter changed from uint192 to uint224)
// See https://github.com/smartcontractkit/ccip/blob/ccip-develop/contracts/src/v0.8/ccip/PriceRegistry.sol#L141
func (p *PriceRegistryV1_2_0) GetTokenPrices(ctx context.Context, wantedTokens []common.Address) ([]TokenPriceUpdate, error) {
	// Make call using 224 ABI.
	tps, err := p.pr.GetTokenPrices(&bind.CallOpts{Context: ctx}, wantedTokens)
	if err != nil {
		return nil, err
	}
	var tpu []TokenPriceUpdate
	for i, tp := range tps {
		tpu = append(tpu, TokenPriceUpdate{
			TokenPrice: TokenPrice{
				Token: wantedTokens[i],
				Value: tp.Value,
			},
			TimestampUnixSec: big.NewInt(int64(tp.Timestamp)),
		})
	}
	return tpu, nil
}

func ApplyPriceRegistryUpdateV1_2_0(t *testing.T, user *bind.TransactOpts, addr common.Address, ec client.Client, gasPrices []GasPrice, tokenPrices []TokenPrice) common.Hash {
	require.True(t, len(gasPrices) <= 1)
	pr, err := price_registry.NewPriceRegistry(addr, ec)
	require.NoError(t, err)
	o, err := pr.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, user.From, o)
	var tps []price_registry.InternalTokenPriceUpdate
	for _, tp := range tokenPrices {
		tps = append(tps, price_registry.InternalTokenPriceUpdate{
			SourceToken: tp.Token,
			UsdPerToken: tp.Value,
		})
	}
	var gps []price_registry.InternalGasPriceUpdate
	for _, gp := range gasPrices {
		gps = append(gps, price_registry.InternalGasPriceUpdate{
			DestChainSelector: gp.DestChainSelector,
			UsdPerUnitGas:     gp.Value,
		})
	}
	tx, err := pr.UpdatePrices(user, price_registry.InternalPriceUpdates{
		TokenPriceUpdates: tps,
		GasPriceUpdates:   gps,
	})
	require.NoError(t, err)
	return tx.Hash()
}
