package v1_2_0

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
)

var (
	_ ccipdata.PriceRegistryReader = &PriceRegistry{}
)

type PriceRegistry struct {
	*v1_0_0.PriceRegistry
	pr *price_registry_1_2_0.PriceRegistry
}

func NewPriceRegistry(lggr logger.Logger, priceRegistryAddr common.Address, lp logpoller.LogPoller, ec client.Client, registerFilters bool) (*PriceRegistry, error) {
	v100, err := v1_0_0.NewPriceRegistry(lggr, priceRegistryAddr, lp, ec, registerFilters)
	if err != nil {
		return nil, err
	}
	priceRegistry, err := price_registry_1_2_0.NewPriceRegistry(priceRegistryAddr, ec)
	if err != nil {
		return nil, err
	}
	return &PriceRegistry{
		PriceRegistry: v100,
		pr:            priceRegistry,
	}, nil
}

// GetTokenPrices must be overridden to use the 1.2 ABI (return parameter changed from uint192 to uint224)
// See https://github.com/smartcontractkit/ccip/blob/ccip-develop/contracts/src/v0.8/ccip/PriceRegistry.sol#L141
func (p *PriceRegistry) GetTokenPrices(ctx context.Context, wantedTokens []cciptypes.Address) ([]cciptypes.TokenPriceUpdate, error) {
	evmAddrs, err := ccipcalc.GenericAddrsToEvm(wantedTokens...)
	if err != nil {
		return nil, err
	}

	// Make call using 224 ABI.
	tps, err := p.pr.GetTokenPrices(&bind.CallOpts{Context: ctx}, evmAddrs)
	if err != nil {
		return nil, err
	}
	var tpu []cciptypes.TokenPriceUpdate
	for i, tp := range tps {
		tpu = append(tpu, cciptypes.TokenPriceUpdate{
			TokenPrice: cciptypes.TokenPrice{
				Token: wantedTokens[i],
				Value: tp.Value,
			},
			TimestampUnixSec: big.NewInt(int64(tp.Timestamp)),
		})
	}
	return tpu, nil
}
