package prices

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
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

func (g ExecGasPriceEstimator) GetGasPrice(ctx context.Context) (GasPrice, error) {
	gasPriceWei, _, err := g.estimator.GetFee(ctx, nil, 0, assets.NewWei(g.maxGasPrice))
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

func (g ExecGasPriceEstimator) DenoteInUSD(p GasPrice, wrappedNativePrice *big.Int) (GasPrice, error) {
	return ccipcalc.CalculateUsdPerUnitGas(p, wrappedNativePrice), nil
}

func (g ExecGasPriceEstimator) Median(gasPrices []GasPrice) (GasPrice, error) {
	var prices []*big.Int
	for _, p := range gasPrices {
		prices = append(prices, p)
	}

	return ccipcalc.BigIntSortedMiddle(prices), nil
}

func (g ExecGasPriceEstimator) Deviates(p1 GasPrice, p2 GasPrice) (bool, error) {
	return ccipcalc.Deviates(p1, p2, g.deviationPPB), nil
}

func (g ExecGasPriceEstimator) EstimateMsgCostUSD(p GasPrice, wrappedNativePrice *big.Int, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) (*big.Int, error) {
	execGasAmount := new(big.Int).Add(big.NewInt(feeBoostingOverheadGas), msg.GasLimit)
	execGasAmount = new(big.Int).Add(execGasAmount, new(big.Int).Mul(big.NewInt(int64(len(msg.Data))), big.NewInt(execGasPerPayloadByte)))
	execGasAmount = new(big.Int).Add(execGasAmount, new(big.Int).Mul(big.NewInt(int64(len(msg.TokenAmounts))), big.NewInt(execGasPerToken)))

	execGasCost := new(big.Int).Mul(execGasAmount, p)

	return ccipcalc.CalculateUsdPerUnitGas(execGasCost, wrappedNativePrice), nil
}

func (g ExecGasPriceEstimator) String(p GasPrice) string {
	var pi *big.Int = p
	return pi.String()
}
