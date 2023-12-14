package gas_test

import (
	"math/big"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	rollupMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestWrappedEvmEstimator(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	// fee values
	gasLimit := uint32(10)
	legacyFee := assets.NewWeiI(10)
	dynamicFee := gas.DynamicFee{
		FeeCap: assets.NewWeiI(20),
		TipCap: assets.NewWeiI(1),
	}
	est := mocks.NewEvmEstimator(t)
	est.On("GetDynamicFee", mock.Anything, mock.Anything, mock.Anything).
		Return(dynamicFee, gasLimit, nil).Twice()
	est.On("GetLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(legacyFee, gasLimit, nil).Twice()
	est.On("BumpDynamicFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(dynamicFee, gasLimit, nil).Once()
	est.On("BumpLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(legacyFee, gasLimit, nil).Once()
	getRootEst := func(logger.Logger) gas.EvmEstimator { return est }

	mockEstimatorName := "WrappedEvmEstimator"
	mockEvmEstimatorName := "WrappedEvmEstimator.MockEstimator"

	// L1Oracle returns the correct L1Oracle interface
	t.Run("L1Oracle", func(t *testing.T) {
		lggr := logger.Test(t)
		// expect nil
		estimator := gas.NewWrappedEvmEstimator(lggr, getRootEst, false, nil)
		l1Oracle := estimator.L1Oracle()
		assert.Nil(t, l1Oracle)

		// expect l1Oracle
		oracle := rollupMocks.NewL1Oracle(t)
		estimator = gas.NewWrappedEvmEstimator(lggr, getRootEst, false, oracle)
		l1Oracle = estimator.L1Oracle()
		assert.Equal(t, oracle, l1Oracle)
	})

	// GetFee returns gas estimation based on configuration value
	t.Run("GetFee", func(t *testing.T) {
		lggr := logger.Test(t)
		// expect legacy fee data
		dynamicFees := false
		estimator := gas.NewWrappedEvmEstimator(lggr, getRootEst, dynamicFees, nil)
		fee, max, err := estimator.GetFee(ctx, nil, 0, nil)
		require.NoError(t, err)
		assert.Equal(t, gasLimit, max)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.DynamicTipCap)
		assert.Nil(t, fee.DynamicFeeCap)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewWrappedEvmEstimator(lggr, getRootEst, dynamicFees, nil)
		fee, max, err = estimator.GetFee(ctx, nil, 0, nil)
		require.NoError(t, err)
		assert.Equal(t, gasLimit, max)
		assert.True(t, dynamicFee.FeeCap.Equal(fee.DynamicFeeCap))
		assert.True(t, dynamicFee.TipCap.Equal(fee.DynamicTipCap))
		assert.Nil(t, fee.Legacy)
	})

	// BumpFee returns bumped fee type based on original fee calculation
	t.Run("BumpFee", func(t *testing.T) {
		lggr := logger.Test(t)
		dynamicFees := false
		estimator := gas.NewWrappedEvmEstimator(lggr, getRootEst, dynamicFees, nil)

		// expect legacy fee data
		fee, max, err := estimator.BumpFee(ctx, gas.EvmFee{Legacy: assets.NewWeiI(0)}, 0, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, gasLimit, max)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.DynamicTipCap)
		assert.Nil(t, fee.DynamicFeeCap)

		// expect dynamic fee data
		fee, max, err = estimator.BumpFee(ctx, gas.EvmFee{
			DynamicFeeCap: assets.NewWeiI(0),
			DynamicTipCap: assets.NewWeiI(0),
		}, 0, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, gasLimit, max)
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
		estimator := gas.NewWrappedEvmEstimator(lggr, getRootEst, dynamicFees, nil)
		total, err := estimator.GetMaxCost(ctx, val, nil, gasLimit, nil)
		require.NoError(t, err)
		fee := new(big.Int).Mul(legacyFee.ToInt(), big.NewInt(int64(gasLimit)))
		assert.Equal(t, new(big.Int).Add(val.ToInt(), fee), total)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewWrappedEvmEstimator(lggr, getRootEst, dynamicFees, nil)
		total, err = estimator.GetMaxCost(ctx, val, nil, gasLimit, nil)
		require.NoError(t, err)
		fee = new(big.Int).Mul(dynamicFee.FeeCap.ToInt(), big.NewInt(int64(gasLimit)))
		assert.Equal(t, new(big.Int).Add(val.ToInt(), fee), total)
	})

	t.Run("Name", func(t *testing.T) {
		lggr := logger.Test(t)

		oracle := rollupMocks.NewL1Oracle(t)
		evmEstimator := mocks.NewEvmEstimator(t)
		evmEstimator.On("Name").Return(mockEvmEstimatorName, nil).Once()

		estimator := gas.NewWrappedEvmEstimator(lggr, func(logger.Logger) gas.EvmEstimator {
			return evmEstimator
		}, false, oracle)

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

		estimator := gas.NewWrappedEvmEstimator(lggr, getEst, false, nil)
		err := estimator.Start(ctx)
		require.NoError(t, err)
		err = estimator.Close()
		require.NoError(t, err)

		estimator = gas.NewWrappedEvmEstimator(lggr, getEst, false, oracle)
		err = estimator.Start(ctx)
		require.NoError(t, err)
		err = estimator.Close()
		require.NoError(t, err)
	})

	t.Run("Read calls both EVM estimator and L1Oracle", func(t *testing.T) {
		lggr := logger.Test(t)
		evmEstimator := mocks.NewEvmEstimator(t)
		oracle := rollupMocks.NewL1Oracle(t)

		evmEstimator.On("Ready").Return(nil).Twice()
		oracle.On("Ready").Return(nil).Once()
		getEst := func(logger.Logger) gas.EvmEstimator { return evmEstimator }

		estimator := gas.NewWrappedEvmEstimator(lggr, getEst, false, nil)
		err := estimator.Ready()
		require.NoError(t, err)

		estimator = gas.NewWrappedEvmEstimator(lggr, getEst, false, oracle)
		err = estimator.Ready()
		require.NoError(t, err)
	})

	t.Run("HealthReport merges report from EVM estimator and L1Oracle", func(t *testing.T) {
		lggr := logger.Test(t)
		evmEstimator := mocks.NewEvmEstimator(t)
		oracle := rollupMocks.NewL1Oracle(t)

		evmEstimatorKey := "evm"
		evmEstimatorError := errors.New("evm error")
		oracleKey := "oracle"
		oracleError := errors.New("oracle error")

		evmEstimator.On("HealthReport").Return(map[string]error{evmEstimatorKey: evmEstimatorError}).Twice()
		oracle.On("HealthReport").Return(map[string]error{oracleKey: oracleError}).Once()
		getEst := func(logger.Logger) gas.EvmEstimator { return evmEstimator }

		estimator := gas.NewWrappedEvmEstimator(lggr, getEst, false, nil)
		report := estimator.HealthReport()
		require.True(t, errors.Is(report[evmEstimatorKey], evmEstimatorError))
		require.Nil(t, report[oracleKey])
		require.NotNil(t, report[mockEstimatorName])

		estimator = gas.NewWrappedEvmEstimator(lggr, getEst, false, oracle)
		report = estimator.HealthReport()
		require.True(t, errors.Is(report[evmEstimatorKey], evmEstimatorError))
		require.True(t, errors.Is(report[oracleKey], oracleError))
		require.NotNil(t, report[mockEstimatorName])
	})
}
