package estimatorconfig_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	mocks2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/estimatorconfig/mocks"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/estimatorconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
)

func TestFeeEstimatorConfigService(t *testing.T) {
	svc := estimatorconfig.NewFeeEstimatorConfigService()
	ctx := context.Background()

	var expectedDestDataAvailabilityOverheadGas int64 = 1
	var expectedDestGasPerDataAvailabilityByte int64 = 2
	var expectedDestDataAvailabilityMultiplierBps int64 = 3

	onRampReader := mocks.NewOnRampReader(t)
	destDataAvailabilityOverheadGas, destGasPerDataAvailabilityByte, destDataAvailabilityMultiplierBps, err := svc.GetDataAvailabilityConfig(ctx)
	require.NoError(t, err) // if onRampReader not set, return nil error and 0 values
	require.EqualValues(t, 0, destDataAvailabilityOverheadGas)
	require.EqualValues(t, 0, destGasPerDataAvailabilityByte)
	require.EqualValues(t, 0, destDataAvailabilityMultiplierBps)
	svc.SetOnRampReader(onRampReader)

	onRampReader.On("GetDynamicConfig", ctx).
		Return(ccip.OnRampDynamicConfig{
			DestDataAvailabilityOverheadGas:   uint32(expectedDestDataAvailabilityOverheadGas),
			DestGasPerDataAvailabilityByte:    uint16(expectedDestGasPerDataAvailabilityByte),
			DestDataAvailabilityMultiplierBps: uint16(expectedDestDataAvailabilityMultiplierBps),
		}, nil).Once()

	destDataAvailabilityOverheadGas, destGasPerDataAvailabilityByte, destDataAvailabilityMultiplierBps, err = svc.GetDataAvailabilityConfig(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedDestDataAvailabilityOverheadGas, destDataAvailabilityOverheadGas)
	require.Equal(t, expectedDestGasPerDataAvailabilityByte, destGasPerDataAvailabilityByte)
	require.Equal(t, expectedDestDataAvailabilityMultiplierBps, destDataAvailabilityMultiplierBps)

	onRampReader.On("GetDynamicConfig", ctx).
		Return(ccip.OnRampDynamicConfig{}, errors.New("test")).Once()
	_, _, _, err = svc.GetDataAvailabilityConfig(ctx)
	require.Error(t, err)
}

func TestModifyGasPriceComponents(t *testing.T) {
	t.Run("success modification", func(t *testing.T) {
		svc := estimatorconfig.NewFeeEstimatorConfigService()
		ctx := context.Background()

		initialExecGasPrice, initialDaGasPrice := big.NewInt(10), big.NewInt(1)

		gpi1 := mocks2.NewGasPriceInterceptor(t)
		svc.AddGasPriceInterceptor(gpi1)

		// change in first interceptor
		firstModExecGasPrice, firstModDaGasPrice := big.NewInt(5), big.NewInt(2)
		gpi1.On("ModifyGasPriceComponents", ctx, initialExecGasPrice, initialDaGasPrice).
			Return(firstModExecGasPrice, firstModDaGasPrice, nil)

		gpi2 := mocks2.NewGasPriceInterceptor(t)
		svc.AddGasPriceInterceptor(gpi2)

		// change in second iterceptor
		secondModExecGasPrice, secondModDaGasPrice := big.NewInt(50), big.NewInt(20)
		gpi2.On("ModifyGasPriceComponents", ctx, firstModExecGasPrice, firstModDaGasPrice).
			Return(secondModExecGasPrice, secondModDaGasPrice, nil)

		// has to return second interceptor values
		resGasPrice, resDAGasPrice, err := svc.ModifyGasPriceComponents(ctx, initialExecGasPrice, initialDaGasPrice)
		require.NoError(t, err)
		require.Equal(t, secondModExecGasPrice.Int64(), resGasPrice.Int64())
		require.Equal(t, secondModDaGasPrice.Int64(), resDAGasPrice.Int64())
	})

	t.Run("error modification", func(t *testing.T) {
		svc := estimatorconfig.NewFeeEstimatorConfigService()
		ctx := context.Background()

		initialExecGasPrice, initialDaGasPrice := big.NewInt(10), big.NewInt(1)
		gpi1 := mocks2.NewGasPriceInterceptor(t)
		svc.AddGasPriceInterceptor(gpi1)
		gpi1.On("ModifyGasPriceComponents", ctx, initialExecGasPrice, initialDaGasPrice).
			Return(nil, nil, errors.New("test"))

		// has to return second interceptor values
		_, _, err := svc.ModifyGasPriceComponents(ctx, initialExecGasPrice, initialDaGasPrice)
		require.Error(t, err)
	})

	t.Run("without interceptors", func(t *testing.T) {
		svc := estimatorconfig.NewFeeEstimatorConfigService()
		ctx := context.Background()

		initialExecGasPrice, initialDaGasPrice := big.NewInt(10), big.NewInt(1)

		// has to return second interceptor values
		resGasPrice, resDAGasPrice, err := svc.ModifyGasPriceComponents(ctx, initialExecGasPrice, initialDaGasPrice)
		require.NoError(t, err)

		// values should not be modified
		require.Equal(t, initialExecGasPrice, resGasPrice)
		require.Equal(t, initialDaGasPrice, resDAGasPrice)
	})
}
