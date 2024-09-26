package gas_test

import (
	"errors"
	"math/big"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	rollupMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
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
		Return(dynamicFee, nil).Times(6)
	est.On("GetLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(legacyFee, gasLimit, nil).Times(6)
	est.On("BumpDynamicFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(dynamicFee, nil).Once()
	est.On("BumpLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(legacyFee, gasLimit, nil).Once()
	getRootEst := func(logger.Logger) gas.EvmEstimator { return est }
	geCfg := gas.NewMockGasConfig()
	geCfg.LimitMultiplierF = limitMultiplier

	mockEstimatorName := "WrappedEvmEstimator"
	mockEvmEstimatorName := "WrappedEvmEstimator.MockEstimator"

	fromAddress := testutils.NewAddress()
	toAddress := testutils.NewAddress()

	// L1Oracle returns the correct L1Oracle interface
	t.Run("L1Oracle", func(t *testing.T) {
		lggr := logger.Test(t)

		evmEstimator := mocks.NewEvmEstimator(t)
		evmEstimator.On("L1Oracle").Return(nil).Once()

		getEst := func(logger.Logger) gas.EvmEstimator { return evmEstimator }

		// expect nil
		estimator := gas.NewEvmFeeEstimator(lggr, getEst, false, nil, nil)
		l1Oracle := estimator.L1Oracle()

		assert.Nil(t, l1Oracle)

		// expect l1Oracle
		oracle, err := rollups.NewL1GasOracle(lggr, nil, chaintype.ChainOptimismBedrock)
		require.NoError(t, err)
		// cast oracle to L1Oracle interface
		estimator = gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg, nil)

		evmEstimator.On("L1Oracle").Return(oracle).Once()
		l1Oracle = estimator.L1Oracle()
		assert.Equal(t, oracle, l1Oracle)
	})

	// GetFee returns gas estimation based on configuration value
	t.Run("GetFee", func(t *testing.T) {
		lggr := logger.Test(t)
		// expect legacy fee data
		dynamicFees := false
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, nil)
		fee, max, err := estimator.GetFee(ctx, nil, 0, nil, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(gasLimit)*limitMultiplier), max)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.DynamicTipCap)
		assert.Nil(t, fee.DynamicFeeCap)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, nil)
		fee, max, err = estimator.GetFee(ctx, nil, gasLimit, nil, nil, nil)
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
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, nil)

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
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, nil)
		total, err := estimator.GetMaxCost(ctx, val, nil, gasLimit, nil, nil, nil)
		require.NoError(t, err)
		fee := new(big.Int).Mul(legacyFee.ToInt(), big.NewInt(int64(gasLimit)))
		fee, _ = new(big.Float).Mul(new(big.Float).SetInt(fee), big.NewFloat(float64(limitMultiplier))).Int(nil)
		assert.Equal(t, new(big.Int).Add(val.ToInt(), fee), total)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, nil)
		total, err = estimator.GetMaxCost(ctx, val, nil, gasLimit, nil, nil, nil)
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
		}, false, geCfg, nil)

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

		estimator := gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg, nil)
		err := estimator.Start(ctx)
		require.NoError(t, err)
		err = estimator.Close()
		require.NoError(t, err)

		evmEstimator.On("L1Oracle", mock.Anything).Return(oracle).Twice()

		estimator = gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg, nil)
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

		estimator := gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg, nil)
		err := estimator.Ready()
		require.NoError(t, err)

		estimator = gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg, nil)
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

		estimator := gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg, nil)
		report := estimator.HealthReport()
		require.True(t, pkgerrors.Is(report[evmEstimatorKey], evmEstimatorError))
		require.Nil(t, report[oracleKey])
		require.NotNil(t, report[mockEstimatorName])

		evmEstimator.On("L1Oracle").Return(oracle).Once()

		estimator = gas.NewEvmFeeEstimator(lggr, getEst, false, geCfg, nil)
		report = estimator.HealthReport()
		require.True(t, pkgerrors.Is(report[evmEstimatorKey], evmEstimatorError))
		require.True(t, pkgerrors.Is(report[oracleKey], oracleError))
		require.NotNil(t, report[mockEstimatorName])
	})

	t.Run("GetFee, estimate gas limit enabled, succeeds", func(t *testing.T) {
		estimatedGasLimit := uint64(5)
		lggr := logger.Test(t)
		// expect legacy fee data
		dynamicFees := false
		geCfg.EstimateLimitF = true
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ethClient.On("EstimateGas", mock.Anything, mock.Anything).Return(estimatedGasLimit, nil).Twice()
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		fee, limit, err := estimator.GetFee(ctx, []byte{}, gasLimit, nil, &fromAddress, &toAddress)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(estimatedGasLimit)*gas.EstimateGasBuffer), limit)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.DynamicTipCap)
		assert.Nil(t, fee.DynamicFeeCap)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		fee, limit, err = estimator.GetFee(ctx, []byte{}, gasLimit, nil, &fromAddress, &toAddress)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(estimatedGasLimit)*gas.EstimateGasBuffer), limit)
		assert.True(t, dynamicFee.FeeCap.Equal(fee.DynamicFeeCap))
		assert.True(t, dynamicFee.TipCap.Equal(fee.DynamicTipCap))
		assert.Nil(t, fee.Legacy)
	})

	t.Run("GetFee, estimate gas limit enabled, estimate exceeds provided limit, returns error", func(t *testing.T) {
		estimatedGasLimit := uint64(100)
		lggr := logger.Test(t)
		// expect legacy fee data
		dynamicFees := false
		geCfg.EstimateLimitF = true
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ethClient.On("EstimateGas", mock.Anything, mock.Anything).Return(estimatedGasLimit, nil).Twice()
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		_, _, err := estimator.GetFee(ctx, []byte{}, gasLimit, nil, &fromAddress, &toAddress)
		require.ErrorIs(t, err, commonfee.ErrFeeLimitTooLow)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		_, _, err = estimator.GetFee(ctx, []byte{}, gasLimit, nil, &fromAddress, &toAddress)
		require.ErrorIs(t, err, commonfee.ErrFeeLimitTooLow)
	})

	t.Run("GetFee, estimate gas limit enabled, buffer exceeds provided limit, fallsback to provided limit", func(t *testing.T) {
		estimatedGasLimit := uint64(15) // same as provided limit
		lggr := logger.Test(t)
		dynamicFees := false // expect legacy fee data
		geCfg.EstimateLimitF = true
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ethClient.On("EstimateGas", mock.Anything, mock.Anything).Return(estimatedGasLimit, nil).Twice()
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		fee, limit, err := estimator.GetFee(ctx, []byte{}, gasLimit, nil, &fromAddress, &toAddress)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(gasLimit)*limitMultiplier), limit)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.DynamicTipCap)
		assert.Nil(t, fee.DynamicFeeCap)

		dynamicFees = true // expect dynamic fee data
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		fee, limit, err = estimator.GetFee(ctx, []byte{}, gasLimit, nil, &fromAddress, &toAddress)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(gasLimit)*limitMultiplier), limit)
		assert.True(t, dynamicFee.FeeCap.Equal(fee.DynamicFeeCap))
		assert.True(t, dynamicFee.TipCap.Equal(fee.DynamicTipCap))
		assert.Nil(t, fee.Legacy)
	})

	t.Run("GetFee, estimate gas limit enabled, RPC fails and fallsback to provided gas limit", func(t *testing.T) {
		lggr := logger.Test(t)
		// expect legacy fee data
		dynamicFees := false
		geCfg.EstimateLimitF = true
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ethClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(0), errors.New("something broke")).Twice()
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		fee, limit, err := estimator.GetFee(ctx, []byte{}, gasLimit, nil, &fromAddress, &toAddress)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(gasLimit)*limitMultiplier), limit)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.DynamicTipCap)
		assert.Nil(t, fee.DynamicFeeCap)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		fee, limit, err = estimator.GetFee(ctx, []byte{}, gasLimit, nil, &fromAddress, &toAddress)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(gasLimit)*limitMultiplier), limit)
		assert.True(t, dynamicFee.FeeCap.Equal(fee.DynamicFeeCap))
		assert.True(t, dynamicFee.TipCap.Equal(fee.DynamicTipCap))
		assert.Nil(t, fee.Legacy)
	})

	t.Run("GetFee, estimate gas limit enabled, provided fee limit 0, returns uncapped estimation", func(t *testing.T) {
		est.On("GetDynamicFee", mock.Anything, mock.Anything).
			Return(dynamicFee, nil).Once()
		est.On("GetLegacyGas", mock.Anything, mock.Anything, uint64(0), mock.Anything).
			Return(legacyFee, uint64(0), nil).Once()
		estimatedGasLimit := uint64(100) // same as provided limit
		lggr := logger.Test(t)
		// expect legacy fee data
		dynamicFees := false
		geCfg.EstimateLimitF = true
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ethClient.On("EstimateGas", mock.Anything, mock.Anything).Return(estimatedGasLimit, nil).Twice()
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		fee, limit, err := estimator.GetFee(ctx, []byte{}, uint64(0), nil, &fromAddress, &toAddress)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(estimatedGasLimit)*gas.EstimateGasBuffer), limit)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.DynamicTipCap)
		assert.Nil(t, fee.DynamicFeeCap)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		fee, limit, err = estimator.GetFee(ctx, []byte{}, 0, nil, &fromAddress, &toAddress)
		require.NoError(t, err)
		assert.Equal(t, uint64(float32(estimatedGasLimit)*gas.EstimateGasBuffer), limit)
		assert.True(t, dynamicFee.FeeCap.Equal(fee.DynamicFeeCap))
		assert.True(t, dynamicFee.TipCap.Equal(fee.DynamicTipCap))
		assert.Nil(t, fee.Legacy)
	})

	t.Run("GetFee, estimate gas limit enabled, provided fee limit 0, returns error on failure", func(t *testing.T) {
		est.On("GetDynamicFee", mock.Anything, mock.Anything).
			Return(dynamicFee, nil).Once()
		est.On("GetLegacyGas", mock.Anything, mock.Anything, uint64(0), mock.Anything).
			Return(legacyFee, uint64(0), nil).Once()
		lggr := logger.Test(t)
		// expect legacy fee data
		dynamicFees := false
		geCfg.EstimateLimitF = true
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ethClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(0), errors.New("something broke")).Twice()
		estimator := gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		_, _, err := estimator.GetFee(ctx, []byte{}, 0, nil, &fromAddress, &toAddress)
		require.Error(t, err)

		// expect dynamic fee data
		dynamicFees = true
		estimator = gas.NewEvmFeeEstimator(lggr, getRootEst, dynamicFees, geCfg, ethClient)
		_, _, err = estimator.GetFee(ctx, []byte{}, 0, nil, &fromAddress, &toAddress)
		require.Error(t, err)
	})
}
