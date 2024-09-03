package prices

import (
	"context"
	"fmt"
	"math/big"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

type ExecGasPriceEstimator struct {
	estimator    gas.EvmFeeEstimator
	maxGasPrice  *big.Int
	deviationPPB int64
}

func NewExecGasPriceEstimator(estimator gas.EvmFeeEstimator, maxGasPrice *big.Int, deviationPPB int64) ExecGasPriceEstimator {
	return ExecGasPriceEstimator{
		estimator:    estimator,
		maxGasPrice:  maxGasPrice,
		deviationPPB: deviationPPB,
	}
}

func (g ExecGasPriceEstimator) GetGasPrice(ctx context.Context) (*big.Int, error) {
	gasPriceWei, _, err := g.estimator.GetFee(ctx, nil, 0, assets.NewWei(g.maxGasPrice), nil, nil)
	if err != nil {
		return nil, err
	}
	// Use legacy if no dynamic is available.
	gasPrice := gasPriceWei.Legacy.ToInt()
	if gasPriceWei.DynamicFeeCap != nil {
		gasPrice = gasPriceWei.DynamicFeeCap.ToInt()
	}
	if gasPrice == nil {
		return nil, fmt.Errorf("missing gas price %+v", gasPriceWei)
	}

	return gasPrice, nil
}

func (g ExecGasPriceEstimator) DenoteInUSD(p *big.Int, wrappedNativePrice *big.Int) (*big.Int, error) {
	return ccipcalc.CalculateUsdPerUnitGas(p, wrappedNativePrice), nil
}

func (g ExecGasPriceEstimator) Median(gasPrices []*big.Int) (*big.Int, error) {
	return ccipcalc.BigIntSortedMiddle(gasPrices), nil
}

func (g ExecGasPriceEstimator) Deviates(p1 *big.Int, p2 *big.Int) (bool, error) {
	return ccipcalc.Deviates(p1, p2, g.deviationPPB), nil
}

func (g ExecGasPriceEstimator) EstimateMsgCostUSD(p *big.Int, wrappedNativePrice *big.Int, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) (*big.Int, error) {
	execGasAmount := new(big.Int).Add(big.NewInt(feeBoostingOverheadGas), msg.GasLimit)
	execGasAmount = new(big.Int).Add(execGasAmount, new(big.Int).Mul(big.NewInt(int64(len(msg.Data))), big.NewInt(execGasPerPayloadByte)))
	execGasAmount = new(big.Int).Add(execGasAmount, new(big.Int).Mul(big.NewInt(int64(len(msg.TokenAmounts))), big.NewInt(execGasPerToken)))

	execGasCost := new(big.Int).Mul(execGasAmount, p)

	return ccipcalc.CalculateUsdPerUnitGas(execGasCost, wrappedNativePrice), nil
}
