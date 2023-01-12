package reasonablegasprice

import (
	"context"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/types"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
)

// reasonableGasPriceProvider provides an estimate for the average gas price
type reasonableGasPriceProvider struct {
	estimator          gas.Estimator
	timeout            time.Duration
	maxGasPrice        *assets.Wei
	supportsDynamicFee bool
}

var _ types.ReasonableGasPrice = (*reasonableGasPriceProvider)(nil)

func NewReasonableGasPriceProvider(
	estimator gas.Estimator,
	timeout time.Duration,
	maxGasPrice *assets.Wei,
	supportsDynamicFee bool,
) types.ReasonableGasPrice {
	return &reasonableGasPriceProvider{
		estimator:          estimator,
		timeout:            timeout,
		maxGasPrice:        maxGasPrice,
		supportsDynamicFee: supportsDynamicFee,
	}
}

func (r *reasonableGasPriceProvider) ReasonableGasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	if r.supportsDynamicFee {
		fee, _, err := r.estimator.GetDynamicFee(ctx, 0, r.maxGasPrice)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch reasonable gas price")
		}
		return fee.FeeCap.ToInt(), nil
	}
	fee, _, err := r.estimator.GetLegacyGas(ctx, []byte{}, 0, r.maxGasPrice)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch reasonable gas price")
	}
	return fee.ToInt(), nil
}
