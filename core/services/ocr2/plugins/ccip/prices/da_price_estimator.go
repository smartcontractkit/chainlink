package prices

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-integrations/evm/gas"
	"github.com/smartcontractkit/chainlink-integrations/evm/gas/rollups"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

const (
	// DANoDeviationThresholdUSD is the lower bound no deviation threshold for DA gas. If the DA gas price is less
	// than this value, we should never trigger a deviation. This is set to 20 gwei in USD terms.
	DANoDeviationThresholdUSD = 20e9
)

type DAGasPriceEstimator struct {
	execEstimator       GasPriceEstimator
	l1Oracle            rollups.L1Oracle
	feeEstimatorConfig  ccipdata.FeeEstimatorConfigReader
	priceEncodingLength uint
	daDeviationPPB      int64
}

func NewDAGasPriceEstimator(
	estimator gas.EvmFeeEstimator,
	maxGasPrice *big.Int,
	deviationPPB int64,
	daDeviationPPB int64,
	feeEstimatorConfig ccipdata.FeeEstimatorConfigReader, // DA Config Cache updates in the onRamp reader and shares the state
) *DAGasPriceEstimator {
	return &DAGasPriceEstimator{
		execEstimator:       NewExecGasPriceEstimator(estimator, maxGasPrice, deviationPPB),
		l1Oracle:            estimator.L1Oracle(),
		priceEncodingLength: daGasPriceEncodingLength,
		daDeviationPPB:      daDeviationPPB,
		feeEstimatorConfig:  feeEstimatorConfig,
	}
}

func (g DAGasPriceEstimator) GetGasPrice(ctx context.Context) (*big.Int, error) {
	execGasPrice, err := g.execEstimator.GetGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	var gasPrice *big.Int = execGasPrice
	if gasPrice.BitLen() > int(g.priceEncodingLength) {
		return nil, fmt.Errorf("native gas price exceeded max range %+v", gasPrice)
	}

	if g.l1Oracle == nil {
		return gasPrice, nil
	}

	daGasPriceWei, err := g.l1Oracle.GasPrice(ctx)
	if err != nil {
		return nil, err
	}

	daGasPrice := daGasPriceWei.ToInt()

	gasPrice, daGasPrice, err = g.feeEstimatorConfig.ModifyGasPriceComponents(ctx, gasPrice, daGasPrice)
	if err != nil {
		return nil, fmt.Errorf("gasPrice modification failed: %w", err)
	}

	if daGasPrice.Cmp(big.NewInt(0)) > 0 {
		if daGasPrice.BitLen() > int(g.priceEncodingLength) {
			return nil, fmt.Errorf("data availability gas price exceeded max range %+v", daGasPrice)
		}

		daGasPrice = new(big.Int).Lsh(daGasPrice, g.priceEncodingLength)
		gasPrice = new(big.Int).Add(gasPrice, daGasPrice)
	}

	return gasPrice, nil
}

func (g DAGasPriceEstimator) DenoteInUSD(ctx context.Context, p *big.Int, wrappedNativePrice *big.Int) (*big.Int, error) {
	daGasPrice, execGasPrice, err := g.parseEncodedGasPrice(p)
	if err != nil {
		return nil, err
	}

	// This assumes l1GasPrice is priced using the same native token as l2 native
	daUSD := ccipcalc.CalculateUsdPerUnitGas(daGasPrice, wrappedNativePrice)
	if daUSD.BitLen() > int(g.priceEncodingLength) {
		return nil, fmt.Errorf("data availability gas price USD exceeded max range %+v", daUSD)
	}
	execUSD := ccipcalc.CalculateUsdPerUnitGas(execGasPrice, wrappedNativePrice)
	if execUSD.BitLen() > int(g.priceEncodingLength) {
		return nil, fmt.Errorf("exec gas price USD exceeded max range %+v", execUSD)
	}

	daUSD = new(big.Int).Lsh(daUSD, g.priceEncodingLength)
	return new(big.Int).Add(daUSD, execUSD), nil
}

func (g DAGasPriceEstimator) Median(ctx context.Context, gasPrices []*big.Int) (*big.Int, error) {
	daPrices := make([]*big.Int, len(gasPrices))
	execPrices := make([]*big.Int, len(gasPrices))

	for i := range gasPrices {
		daGasPrice, execGasPrice, err := g.parseEncodedGasPrice(gasPrices[i])
		if err != nil {
			return nil, err
		}

		daPrices[i] = daGasPrice
		execPrices[i] = execGasPrice
	}

	daMedian := ccipcalc.BigIntSortedMiddle(daPrices)
	execMedian := ccipcalc.BigIntSortedMiddle(execPrices)

	daMedian = new(big.Int).Lsh(daMedian, g.priceEncodingLength)
	return new(big.Int).Add(daMedian, execMedian), nil
}

func (g DAGasPriceEstimator) Deviates(ctx context.Context, p1, p2 *big.Int) (bool, error) {
	p1DAGasPrice, p1ExecGasPrice, err := g.parseEncodedGasPrice(p1)
	if err != nil {
		return false, err
	}
	p2DAGasPrice, p2ExecGasPrice, err := g.parseEncodedGasPrice(p2)
	if err != nil {
		return false, err
	}

	execDeviates, err := g.execEstimator.Deviates(ctx, p1ExecGasPrice, p2ExecGasPrice)
	if err != nil {
		return false, err
	}
	if execDeviates {
		return execDeviates, nil
	}

	return ccipcalc.DeviatesOnCurve(p1DAGasPrice, p2DAGasPrice, big.NewInt(DANoDeviationThresholdUSD), g.daDeviationPPB), nil
}

func (g DAGasPriceEstimator) EstimateMsgCostUSD(ctx context.Context, p *big.Int, wrappedNativePrice *big.Int, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) (*big.Int, error) {
	daGasPrice, execGasPrice, err := g.parseEncodedGasPrice(p)
	if err != nil {
		return nil, err
	}

	execCostUSD, err := g.execEstimator.EstimateMsgCostUSD(ctx, execGasPrice, wrappedNativePrice, msg)
	if err != nil {
		return nil, err
	}

	// If there is data availability price component, then include data availability cost in fee estimation
	if daGasPrice.Cmp(big.NewInt(0)) > 0 {
		daGasCostUSD, err := g.estimateDACostUSD(daGasPrice, wrappedNativePrice, msg)
		if err != nil {
			return nil, err
		}
		execCostUSD = new(big.Int).Add(daGasCostUSD, execCostUSD)
	}
	return execCostUSD, nil
}

func (g DAGasPriceEstimator) parseEncodedGasPrice(p *big.Int) (*big.Int, *big.Int, error) {
	if p.BitLen() > int(g.priceEncodingLength*2) {
		return nil, nil, fmt.Errorf("encoded gas price exceeded max range %+v", p)
	}

	daGasPrice := new(big.Int).Rsh(p, g.priceEncodingLength)

	daStart := new(big.Int).Lsh(big.NewInt(1), g.priceEncodingLength)
	execGasPrice := new(big.Int).Mod(p, daStart)

	return daGasPrice, execGasPrice, nil
}

func (g DAGasPriceEstimator) estimateDACostUSD(daGasPrice *big.Int, wrappedNativePrice *big.Int, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) (*big.Int, error) {
	var sourceTokenDataLen int
	for _, tokenData := range msg.SourceTokenData {
		sourceTokenDataLen += len(tokenData)
	}

	daOverheadGas, gasPerDAByte, daMultiplier, err := g.feeEstimatorConfig.GetDataAvailabilityConfig(context.Background())
	if err != nil {
		return nil, err
	}

	dataLen := evmMessageFixedBytes + len(msg.Data) + len(msg.TokenAmounts)*evmMessageBytesPerToken + sourceTokenDataLen
	dataGas := big.NewInt(int64(dataLen)*gasPerDAByte + daOverheadGas)

	dataGasEstimate := new(big.Int).Mul(dataGas, daGasPrice)
	dataGasEstimate = new(big.Int).Div(new(big.Int).Mul(dataGasEstimate, big.NewInt(daMultiplier)), big.NewInt(daMultiplierBase))

	return ccipcalc.CalculateUsdPerUnitGas(dataGasEstimate, wrappedNativePrice), nil
}
