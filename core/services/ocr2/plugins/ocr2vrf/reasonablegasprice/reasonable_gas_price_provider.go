package reasonablegasprice

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-vrf/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
)

// reasonableGasPriceProvider provides an estimate for the average gas price
type reasonableGasPriceProvider struct {
	estimator          gas.EvmFeeEstimator
	timeout            time.Duration
	maxGasPrice        *assets.Wei
	supportsDynamicFee bool
}

var _ types.ReasonableGasPrice = (*reasonableGasPriceProvider)(nil)

func NewReasonableGasPriceProvider(
	estimator gas.EvmFeeEstimator,
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

// TODO: implement this function to use a gas estimator. This change can be rolled out
// to all nodes while the on-chain `useReasonableGasPrice` flag is disabled. Once reasonable
// gas prices reported by nodes become 'reasonable' the flag can be enabled.
func (r *reasonableGasPriceProvider) ReasonableGasPrice() (*big.Int, error) {
	return big.NewInt(0), nil
}
