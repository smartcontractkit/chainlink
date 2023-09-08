package prices

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
)

type noOpGetter struct{}

func NewNoOpGetter() gas.PriceComponentGetter {
	return &noOpGetter{}
}

func (e *noOpGetter) RefreshComponents(_ context.Context) error {
	return nil
}

func (e *noOpGetter) GetPriceComponents(_ context.Context, gasPrice *assets.Wei) (prices []gas.PriceComponent, err error) {
	return []gas.PriceComponent{
		{
			Price:     gasPrice,
			PriceType: gas.GAS_PRICE,
		},
	}, nil
}
