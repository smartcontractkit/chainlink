package prices

import (
	"context"
	"fmt"
	"math/big"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	"github.com/smartcontractkit/chainlink-integrations/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

const (
	// ExecNoDeviationThresholdUSD is the lower bound no deviation threshold for exec gas. If the exec gas price is
	// less than this value, we should never trigger a deviation. This is set to 10 gwei in USD terms.
	ExecNoDeviationThresholdUSD = 10e9
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
	gasPrice := gasPriceWei.GasPrice.ToInt()
	if gasPriceWei.GasFeeCap != nil {
		gasPrice = gasPriceWei.GasFeeCap.ToInt()
	}
	if gasPrice == nil {
		return nil, fmt.Errorf("missing gas price %+v", gasPriceWei)
	}

	return gasPrice, nil
}

func (g ExecGasPriceEstimator) DenoteInUSD(ctx context.Context, p *big.Int, wrappedNativePrice *big.Int) (*big.Int, error) {
	return ccipcalc.CalculateUsdPerUnitGas(p, wrappedNativePrice), nil
}

func (g ExecGasPriceEstimator) Median(ctx context.Context, gasPrices []*big.Int) (*big.Int, error) {
	return ccipcalc.BigIntSortedMiddle(gasPrices), nil
}

func (g ExecGasPriceEstimator) Deviates(ctx context.Context, p1 *big.Int, p2 *big.Int) (bool, error) {
	return ccipcalc.DeviatesOnCurve(p1, p2, big.NewInt(ExecNoDeviationThresholdUSD), g.deviationPPB), nil
}

func (g ExecGasPriceEstimator) EstimateMsgCostUSD(ctx context.Context, p *big.Int, wrappedNativePrice *big.Int, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) (*big.Int, error) {
	execGasAmount := new(big.Int).Add(big.NewInt(feeBoostingOverheadGas), msg.GasLimit)
	execGasAmount = new(big.Int).Add(execGasAmount, new(big.Int).Mul(big.NewInt(int64(len(msg.Data))), big.NewInt(execGasPerPayloadByte)))
	execGasAmount = new(big.Int).Add(execGasAmount, new(big.Int).Mul(big.NewInt(int64(len(msg.TokenAmounts))), big.NewInt(execGasPerToken)))

	execGasCost := new(big.Int).Mul(execGasAmount, p)

	return ccipcalc.CalculateUsdPerUnitGas(execGasCost, wrappedNativePrice), nil
}
