package gas_test

import (
	"math/big"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	rollupMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
)

func TestWrappedEvmEstimator(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)

	// fee values
	gasLimit := uint64(10)
	legacyFee := assets.NewWeiI(10)
	dynamicFee := gas.DynamicFee{
		FeeCap: assets.NewWeiI(20),
		TipCap: assets.NewWeiI(1),
	}
	limitMultiplier := float32(1.5)

	est := mocks.NewEvmEstimator(t)
	est.On("GetDynamicFee", mock.Anything, mock.Anything).
		Return(dynamicFee, nil).Twice()
	est.On("GetLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(legacyFee, gasLimit, nil).Twice()
	est.On("BumpDynamicFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(dynamicFee, nil).Once()
	est.On("BumpLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(legacyFee, gasLimit, nil).Once()
	getRootEst := func(logger.Logger) gas.EvmEstimator { return est }
	geCfg := gas.NewMockGasConfig()
	geCfg.LimitMultiplierF = limitMultiplier

	mockEstimatorName := "WrappedEvmEstimator"
	mockEvmEstimatorName := "WrappedEvmEstimator.MockEstimator"

	// L1Oracle returns the correct L1Oracle interface
	t.Run("L1Oracle", func(t *testing.T) {
		lggr := logger.Test(t)

		evmEstimator := mocks.NewEvmEstimator(t)
		evmEstimator.On("L1Oracle").Return(nil).Once()

		getEst := func(logger.Logger) gas.EvmEstimator { return evmEstimator }

		// expect nil
		estimator := gas.NewEvmFeeEstimator(lggr, getEst, false, nil)
		l1Oracle := estimator.L1Oracle()

		assert.Nil(t, l1Oracle)

		// expect l1Oracle
		oracle := rollups.NewL1GasOracle(lggr, nil, chaintype.ChainOptimismBedrock)
		// cast oracle to L1Oracle interface
		estimator = gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg)

		evmEstimator.On("L1Oracle").Return(oracle).Once()
		l1Oracle = estimator.L1Oracle()
		assert.Equal(t, oracle, l1Oracle)
	})

	// GetFee returns gas estimation based on configuration value
	t.Run("GetFee", func(t *testing.T) {
		lggr := logger.Test(t)
		// expect legacy fee data
		dynamicFees := false
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg)
		fee, max, err := estimator.GetFee(ctx, nil, 0, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(gasLimit)*limitMultiplier), max)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.DynamicTipCap)
		assert.Nil(t, fee.DynamicFeeCap)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg)
		fee, max, err = estimator.GetFee(ctx, nil, gasLimit, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(gasLimit)*limitMultiplier), max)
		assert.True(t, dynamicFee.FeeCap.Equal(fee.DynamicFeeCap))
		assert.True(t, dynamicFee.TipCap.Equal(fee.DynamicTipCap))
		assert.Nil(t, fee.Legacy)
	})

	// BumpFee returns bumped fee type based on original fee calculation
	t.Run("BumpFee", func(t *testing.T) {
		lggr := logger.Test(t)
		dynamicFees := false
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg)

		// expect legacy fee data
		fee, max, err := estimator.BumpFee(ctx, gas.EvmFee{Legacy: assets.NewWeiI(0)}, 0, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(gasLimit)*limitMultiplier), max)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.DynamicTipCap)
		assert.Nil(t, fee.DynamicFeeCap)

		// expect dynamic fee data
		fee, max, err = estimator.BumpFee(ctx, gas.EvmFee{
			DynamicFeeCap: assets.NewWeiI(0),
			DynamicTipCap: assets.NewWeiI(0),
		}, gasLimit, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(gasLimit)*limitMultiplier), max)
		assert.True(t, dynamicFee.FeeCap.Equal(fee.DynamicFeeCap))
		assert.True(t, dynamicFee.TipCap.Equal(fee.DynamicTipCap))
		assert.Nil(t, fee.Legacy)

		// expect error
		_, _, err = estimator.BumpFee(ctx, gas.EvmFee{}, 0, nil, nil)
		assert.Error(t, err)
		_, _, err = estimator.BumpFee(ctx, gas.EvmFee{
			Legacy:        legacyFee,
			DynamicFeeCap: dynamicFee.FeeCap,
			DynamicTipCap: dynamicFee.TipCap,
		}, 0, nil, nil)
		assert.Error(t, err)
	})

	t.Run("GetMaxCost", func(t *testing.T) {
		lggr := logger.Test(t)
		val := assets.NewEthValue(1)

		// expect legacy fee data
		dynamicFees := false
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg)
		total, err := estimator.GetMaxCost(ctx, val, nil, gasLimit, nil)
		require.NoError(t, err)
		fee := new(big.Int).Mul(legacyFee.ToInt(), big.NewInt(int64(gasLimit)))
		fee, _ = new(big.Float).Mul(new(big.Float).SetInt(fee), big.NewFloat(float64(limitMultiplier))).Int(nil)
		assert.Equal(t, new(big.Int).Add(val.ToInt(), fee), total)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg)
		total, err = estimator.GetMaxCost(ctx, val, nil, gasLimit, nil)
		require.NoError(t, err)
		fee = new(big.Int).Mul(dynamicFee.FeeCap.ToInt(), big.NewInt(int64(gasLimit)))
		fee, _ = new(big.Float).Mul(new(big.Float).SetInt(fee), big.NewFloat(float64(limitMultiplier))).Int(nil)
		assert.Equal(t, new(big.Int).Add(val.ToInt(), fee), total)
	})

	t.Run("Name", func(t *testing.T) {
		lggr := logger.Test(t)

		evmEstimator := mocks.NewEvmEstimator(t)
		evmEstimator.On("Name").Return(mockEvmEstimatorName, nil).Once()

		estimator := gas.NewEvmFeeEstimator(lggr, func(logger.Logger) gas.EvmEstimator {
			return evmEstimator
		}, false, geCfg)

		require.Equal(t, mockEstimatorName, estimator.Name())
		require.Equal(t, mockEvmEstimatorName, evmEstimator.Name())
	})

	t.Run("Start and stop calls both EVM estimator and L1Oracle", func(t *testing.T) {
		lggr := logger.Test(t)
		oracle := rollupMocks.NewL1Oracle(t)
		evmEstimator := mocks.NewEvmEstimator(t)

		evmEstimator.On("Start", mock.Anything).Return(nil).Twice()
		evmEstimator.On("Close").Return(nil).Twice()
		oracle.On("Start", mock.Anything).Return(nil).Once()
		oracle.On("Close").Return(nil).Once()
		getEst := func(logger.Logger) gas.EvmEstimator { return evmEstimator }

		evmEstimator.On("L1Oracle", mock.Anything).Return(nil).Twice()

		estimator := gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg)
		err := estimator.Start(ctx)
		require.NoError(t, err)
		err = estimator.Close()
		require.NoError(t, err)

		evmEstimator.On("L1Oracle", mock.Anything).Return(oracle).Twice()

		estimator = gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg)
		err = estimator.Start(ctx)
		require.NoError(t, err)
		err = estimator.Close()
		require.NoError(t, err)
	})

	t.Run("Read calls both EVM estimator and L1Oracle", func(t *testing.T) {
		lggr := logger.Test(t)
		evmEstimator := mocks.NewEvmEstimator(t)
		oracle := rollupMocks.NewL1Oracle(t)

		evmEstimator.On("L1Oracle").Return(oracle).Twice()
		evmEstimator.On("Ready").Return(nil).Twice()
		oracle.On("Ready").Return(nil).Twice()
		getEst := func(logger.Logger) gas.EvmEstimator { return evmEstimator }

		estimator := gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg)
		err := estimator.Ready()
		require.NoError(t, err)

		estimator = gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg)
		err = estimator.Ready()
		require.NoError(t, err)
	})

	t.Run("HealthReport merges report from EVM estimator and L1Oracle", func(t *testing.T) {
		lggr := logger.Test(t)
		evmEstimator := mocks.NewEvmEstimator(t)
		oracle := rollupMocks.NewL1Oracle(t)

		evmEstimatorKey := "evm"
		evmEstimatorError := pkgerrors.New("evm error")
		oracleKey := "oracle"
		oracleError := pkgerrors.New("oracle error")

		evmEstimator.On("L1Oracle").Return(nil).Once()
		evmEstimator.On("HealthReport").Return(map[string]error{evmEstimatorKey: evmEstimatorError}).Twice()

		oracle.On("HealthReport").Return(map[string]error{oracleKey: oracleError}).Once()
		getEst := func(logger.Logger) gas.EvmEstimator { return evmEstimator }

		estimator := gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg)
		report := estimator.HealthReport()
		require.True(t, pkgerrors.Is(report[evmEstimatorKey], evmEstimatorError))
		require.Nil(t, report[oracleKey])
		require.NotNil(t, report[mockEstimatorName])

		evmEstimator.On("L1Oracle").Return(oracle).Once()

		estimator = gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg)
		report = estimator.HealthReport()
		require.True(t, pkgerrors.Is(report[evmEstimatorKey], evmEstimatorError))
		require.True(t, pkgerrors.Is(report[oracleKey], oracleError))
		require.NotNil(t, report[mockEstimatorName])
	})
}
