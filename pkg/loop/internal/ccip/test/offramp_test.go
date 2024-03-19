package test

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/ccip"
	loopnet "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	loopnettest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net/test"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
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
	ctx := tests.Context(t)
	// create a price registry server
	port := freeport.GetOne(t)
	addr := fmt.Sprintf("localhost:%d", port)
	lis, err := net.Listen("tcp", addr)
	require.NoError(t, err, "failed to listen on port %d", port)
	t.Cleanup(func() { lis.Close() })
	// we explicitly stop the server later, do not add a cleanup function here
	testServer := grpc.NewServer()
	// handle client close and server stop
	shutdown := make(chan struct{})
	closer := &serviceCloser{closeFn: func() error { close(shutdown); return nil }}
	lggr := logger.Test(t)
	broker := &loopnettest.Broker{T: t}
	brokerExt := &loopnet.BrokerExt{
		Broker:       broker,
		BrokerConfig: loopnet.BrokerConfig{Logger: lggr, StopCh: make(chan struct{})},
	}
	offRamp, err := ccip.NewOffRampReaderGRPCServer(OffRampReader, brokerExt)
	require.NoError(t, err)
	offRamp = offRamp.WithCloser(closer)

	ccippb.RegisterOffRampReaderServer(testServer, offRamp)
	// start the server and shutdown handler
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, testServer.Serve(lis))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-shutdown
		t.Log("shutting down server")
		testServer.Stop()
	}()
	// create a price registry client
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "failed to dial %s", addr)
	t.Cleanup(func() { conn.Close() })
	client := ccip.NewOffRampReaderGRPCClient(conn, brokerExt)

	// test the client
	roundTripOffRampTests(ctx, t, client)
	// closing the client executes the shutdown callback
	// which stops the server.  the wg.Wait() below ensures
	// that the server has stopped, which is what we care about.
	cerr := client.Close()
	require.NoError(t, cerr, "failed to close client %T, %v", cerr, status.Code(cerr))
	wg.Wait()
}

type serviceCloser struct {
	closeFn func() error
}

func (s *serviceCloser) Close() error { return s.closeFn() }

// roundTripOffRampTests tests the round trip of the client<->server.
// it should exercise all the methods of the client.
// do not add client.Close to this test, test that from the driver test
func roundTripOffRampTests(ctx context.Context, t *testing.T, client *ccip.OffRampReaderGRPCClient) {
	t.Run("Address", func(t *testing.T) {
		address, err := client.Address(ctx)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.addressResponse, address)
	})

	t.Run("ChangeConfig", func(t *testing.T) {
		gotAddr1, gotAddr2, err := client.ChangeConfig(ctx, OffRampReader.changeConfigRequest.onchainConfig, OffRampReader.changeConfigRequest.offchainConfig)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.changeConfigResponse.onchainConfigDigest, gotAddr1)
		assert.Equal(t, OffRampReader.changeConfigResponse.offchainConfigDigest, gotAddr2)
	})

	t.Run("CurrentRateLimiterState", func(t *testing.T) {
		state, err := client.CurrentRateLimiterState(ctx)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.currentRateLimiterStateResponse, state)
	})

	t.Run("DecodeExecutionReport", func(t *testing.T) {
		report, err := client.DecodeExecutionReport(ctx, OffRampReader.decodeExecutionReportRequest)
		require.NoError(t, err)
		if !reflect.DeepEqual(OffRampReader.decodeExecutionReportResponse.Messages, report.Messages) {
			t.Errorf("expected messages %v, got %v", OffRampReader.decodeExecutionReportResponse.Messages, report.Messages)
		}
	})

	t.Run("EncodeExecutionReport", func(t *testing.T) {
		report, err := client.EncodeExecutionReport(ctx, OffRampReader.encodeExecutionReportRequest)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.encodeExecutionReportResponse, report)
	})

	// exercise all the gas price estimator methods
	t.Run("GasPriceEstimator", func(t *testing.T) {
		estimator, err := client.GasPriceEstimator(ctx)
		require.NoError(t, err)

		t.Run("GetGasPrice", func(t *testing.T) {
			price, err := estimator.GetGasPrice(ctx)
			require.NoError(t, err)
			assert.Equal(t, GasPriceEstimatorExec.getGasPriceResponse, price)
		})

		t.Run("DenoteInUSD", func(t *testing.T) {
			usd, err := estimator.DenoteInUSD(
				GasPriceEstimatorExec.denoteInUSDRequest.p,
				GasPriceEstimatorExec.denoteInUSDRequest.wrappedNativePrice,
			)
			require.NoError(t, err)
			assert.Equal(t, GasPriceEstimatorExec.denoteInUSDResponse.result, usd)
		})

		t.Run("EstimateMsgCostUSD", func(t *testing.T) {
			cost, err := estimator.EstimateMsgCostUSD(
				GasPriceEstimatorExec.estimateMsgCostUSDRequest.p,
				GasPriceEstimatorExec.estimateMsgCostUSDRequest.wrappedNativePrice,
				GasPriceEstimatorExec.estimateMsgCostUSDRequest.msg,
			)
			require.NoError(t, err)
			assert.Equal(t, GasPriceEstimatorExec.estimateMsgCostUSDResponse, cost)
		})

		t.Run("Median", func(t *testing.T) {
			median, err := estimator.Median(GasPriceEstimatorExec.medianRequest.gasPrices)
			require.NoError(t, err)
			assert.Equal(t, GasPriceEstimatorExec.medianResponse, median)
		})
	})

	t.Run("GetExecutionState", func(t *testing.T) {
		state, err := client.GetExecutionState(ctx, OffRampReader.getExecutionStateRequest)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getExecutionStateResponse, state)
	})

	t.Run("GetExecutionStateChangesBetweenSeqNums", func(t *testing.T) {
		state, err := client.GetExecutionStateChangesBetweenSeqNums(ctx, OffRampReader.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMin, OffRampReader.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMax, OffRampReader.getExecutionStateChangesBetweenSeqNumsRequest.confirmations)
		require.NoError(t, err)
		if !reflect.DeepEqual(OffRampReader.getExecutionStateChangesBetweenSeqNumsResponse.executionStateChangedWithTxMeta, state) {
			t.Errorf("expected %v, got %v", OffRampReader.getExecutionStateChangesBetweenSeqNumsResponse, state)
		}
	})

	t.Run("GetSenderNonce", func(t *testing.T) {
		nonce, err := client.GetSenderNonce(ctx, OffRampReader.getSenderNonceRequest)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getSenderNonceResponse, nonce)
	})

	t.Run("GetSourceToDestTokensMapping", func(t *testing.T) {
		mapping, err := client.GetSourceToDestTokensMapping(ctx)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getSourceToDestTokensMappingResponse, mapping)
	})

	t.Run("GetStaticConfig", func(t *testing.T) {
		config, err := client.GetStaticConfig(ctx)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getStaticConfigResponse, config)
	})

	t.Run("GetTokens", func(t *testing.T) {
		tokens, err := client.GetTokens(ctx)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.getTokensResponse, tokens)
	})

	t.Run("OffchainConfig", func(t *testing.T) {
		config, err := client.OffchainConfig(ctx)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.offchainConfigResponse, config)
	})

	t.Run("OnchainConfig", func(t *testing.T) {
		config, err := client.OnchainConfig(ctx)
		require.NoError(t, err)
		assert.Equal(t, OffRampReader.onchainConfigResponse, config)
	})
}
