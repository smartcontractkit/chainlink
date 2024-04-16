package test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	loopnet "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ccip"
	looptest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestStaticOffRamp(t *testing.T) {
	t.Parallel()

	// static test implementation is self consistent
	ctx := context.Background()
	assert.NoError(t, OffRampReader.Evaluate(ctx, OffRampReader))

	// error when the test implementation is evaluates something that differs from the static implementation
	botched := OffRampReader
	botched.addressResponse = "oops"
	err := OffRampReader.Evaluate(ctx, botched)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "oops")
}

func TestOffRampGRPC(t *testing.T) {
	t.Parallel()

	scaffold := looptest.NewGRPCScaffold(t, setupOffRampServer, ccip.NewOffRampReaderGRPCClient)
	roundTripOffRampTests(t, scaffold.Client())

	// offramp implements dependency management, test that it closes properly
	t.Run("Dependency management", func(t *testing.T) {
		d := &looptest.MockDep{}
		scaffold.Server().AddDep(d)
		assert.False(t, d.IsClosed())
		scaffold.Client().Close()
		assert.True(t, d.IsClosed())
	})
}

// roundTripOffRampTests tests the round trip of the client<->server.
// it should exercise all the methods of the client.
// do not add client.Close to this test, test that from the driver test
// func roundTripOffRampTests(ctx context.Context, t *testing.T, client *ccip.OffRampReaderGRPCClient) {
func roundTripOffRampTests(t *testing.T, client cciptypes.OffRampReader) {
	t.Run("Address", func(t *testing.T) {
		address, err := client.Address(tests.Context(t))
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.addressResponse, address)
	})

	t.Run("ChangeConfig", func(t *testing.T) {
		gotAddr1, gotAddr2, err := client.ChangeConfig(tests.Context(t), OffRampReader.changeConfigRequest.onchainConfig, OffRampReader.changeConfigRequest.offchainConfig)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.changeConfigResponse.onchainConfigDigest, gotAddr1)
		assert.Equal(t, OffRampReader.changeConfigResponse.offchainConfigDigest, gotAddr2)
	})

	t.Run("CurrentRateLimiterState", func(t *testing.T) {
		state, err := client.CurrentRateLimiterState(tests.Context(t))
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.currentRateLimiterStateResponse, state)
	})

	t.Run("DecodeExecutionReport", func(t *testing.T) {
		report, err := client.DecodeExecutionReport(tests.Context(t), OffRampReader.decodeExecutionReportRequest)
		require.NoError(t, err)
		if !reflect.DeepEqual(OffRampReader.decodeExecutionReportResponse, report) {
			t.Errorf("expected messages %v, got %v", OffRampReader.decodeExecutionReportResponse, report)
		}
	})

	t.Run("EncodeExecutionReport", func(t *testing.T) {
		report, err := client.EncodeExecutionReport(tests.Context(t), OffRampReader.encodeExecutionReportRequest)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.encodeExecutionReportResponse, report)
	})

	// exercise all the gas price estimator methods
	t.Run("GasPriceEstimator", func(t *testing.T) {
		estimator, err := client.GasPriceEstimator(tests.Context(t))
		require.NoError(t, err)
		gasClient, ok := estimator.(*ccip.ExecGasEstimatorGRPCClient)
		require.True(t, ok, "expected GasPriceEstimatorGRPCClient")
		roundTripGasPriceEstimatorExecTests(t, gasClient)
	})

	t.Run("GetExecutionState", func(t *testing.T) {
		state, err := client.GetExecutionState(tests.Context(t), OffRampReader.getExecutionStateRequest)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getExecutionStateResponse, state)
	})

	t.Run("GetExecutionStateChangesBetweenSeqNums", func(t *testing.T) {
		state, err := client.GetExecutionStateChangesBetweenSeqNums(tests.Context(t), OffRampReader.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMin, OffRampReader.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMax, OffRampReader.getExecutionStateChangesBetweenSeqNumsRequest.confirmations)
		require.NoError(t, err)
		if !reflect.DeepEqual(OffRampReader.getExecutionStateChangesBetweenSeqNumsResponse.executionStateChangedWithTxMeta, state) {
			t.Errorf("expected %v, got %v", OffRampReader.getExecutionStateChangesBetweenSeqNumsResponse, state)
		}
	})

	t.Run("ListSenderNonces", func(t *testing.T) {
		nonce, err := client.ListSenderNonces(tests.Context(t), OffRampReader.listSenderNoncesRequest)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.listSenderNoncesResponse, nonce)
	})

	t.Run("GetSourceToDestTokensMapping", func(t *testing.T) {
		mapping, err := client.GetSourceToDestTokensMapping(tests.Context(t))
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getSourceToDestTokensMappingResponse, mapping)
	})

	t.Run("GetStaticConfig", func(t *testing.T) {
		config, err := client.GetStaticConfig(tests.Context(t))
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getStaticConfigResponse, config)
	})

	t.Run("GetTokens", func(t *testing.T) {
		tokens, err := client.GetTokens(tests.Context(t))
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getTokensResponse, tokens)
	})

	t.Run("GetRouter", func(t *testing.T) {
		router, err := client.GetRouter(tests.Context(t))
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getRouterResponse, router)
	})

	t.Run("OffchainConfig", func(t *testing.T) {
		config, err := client.OffchainConfig(tests.Context(t))
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.offchainConfigResponse, config)
	})

	t.Run("OnchainConfig", func(t *testing.T) {
		config, err := client.OnchainConfig(tests.Context(t))
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.onchainConfigResponse, config)
	})
}

func setupOffRampServer(t *testing.T, s *grpc.Server, b *loopnet.BrokerExt) *ccip.OffRampReaderGRPCServer {
	offRampProvider, err := ccip.NewOffRampReaderGRPCServer(OffRampReader, b)
	require.NoError(t, err)
	ccippb.RegisterOffRampReaderServer(s, offRampProvider)
	return offRampProvider
}
