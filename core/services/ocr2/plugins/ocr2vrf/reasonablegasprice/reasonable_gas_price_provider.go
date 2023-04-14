package reasonablegasprice

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/ocr2vrf/types"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// reasonableGasPriceProvider provides an estimate for the average gas price
type reasonableGasPriceProvider struct {
	estimator          txmgrtypes.FeeEstimator[*evmtypes.Head, gas.EvmFee, *assets.Wei, evmtypes.TxHash]
	timeout            time.Duration
	maxGasPrice        *assets.Wei
	supportsDynamicFee bool
}

var _ types.ReasonableGasPrice = (*reasonableGasPriceProvider)(nil)

func NewReasonableGasPriceProvider(
	estimator txmgrtypes.FeeEstimator[*evmtypes.Head, gas.EvmFee, *assets.Wei, evmtypes.TxHash],
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
