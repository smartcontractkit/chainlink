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

func TestStaticCommitStore(t *testing.T) {
	t.Parallel()

	// static test implementation is self consistent
	ctx := context.Background()
	assert.NoError(t, CommitStoreReader.Evaluate(ctx, CommitStoreReader))

	// error when the test implementation is evaluates something that differs from the static implementation
	botched := CommitStoreReader
	botched.changeConfigResponse = "not the right conifg"
	err := CommitStoreReader.Evaluate(ctx, botched)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not the right conifg")
}

func TestCommitStoreGRPC(t *testing.T) {
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
	offRamp, err := ccip.NewCommitStoreReaderGRPCServer(CommitStoreReader, brokerExt)
	require.NoError(t, err)
	offRamp = offRamp.WithCloser(closer)

	ccippb.RegisterCommitStoreReaderServer(testServer, offRamp)
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
	client := ccip.NewCommitStoreReaderGRPCClient(conn, brokerExt)

	// test the client
	roundTripCommitStoreTests(ctx, t, client)
	// closing the client executes the shutdown callback
	// which stops the server.  the wg.Wait() below ensures
	// that the server has stopped, which is what we care about.
	cerr := client.Close()
	require.NoError(t, cerr, "failed to close client %T, %v", cerr, status.Code(cerr))
	wg.Wait()
}

// roundTripCommitStoreTests tests the round trip of the client<->server.
// it should exercise all the methods of the client.
// do not add client.Close to this test, test that from the driver test
func roundTripCommitStoreTests(ctx context.Context, t *testing.T, client *ccip.CommitStoreGRPCClient) {
	t.Run("ChangeConfig", func(t *testing.T) {
		gotAddr, err := client.ChangeConfig(ctx, CommitStoreReader.changeConfigRequest.onchainConfig, CommitStoreReader.changeConfigRequest.offchainConfig)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.changeConfigResponse, gotAddr)
	})

	t.Run("DecodeCommitReport", func(t *testing.T) {
		report, err := client.DecodeCommitReport(ctx, CommitStoreReader.decodeCommitReportRequest)
		require.NoError(t, err)
		if !reflect.DeepEqual(CommitStoreReader.decodeCommitReportResponse, report) {
			t.Errorf("expected %v, got %v", CommitStoreReader.decodeCommitReportResponse, report)
		}
	})

	// reuse the test data for the encode method
	t.Run("EncodeCommtReport", func(t *testing.T) {
		report, err := client.EncodeCommitReport(ctx, CommitStoreReader.decodeCommitReportResponse)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.decodeCommitReportRequest, report)
	})

	// exercise all the gas price estimator methods
	t.Run("GasPriceEstimator", func(t *testing.T) {
		estimator, err := client.GasPriceEstimator(ctx)
		require.NoError(t, err)

		t.Run("GetGasPrice", func(t *testing.T) {
			price, err := estimator.GetGasPrice(ctx)
			require.NoError(t, err)
			assert.Equal(t, GasPriceEstimatorCommit.getGasPriceResponse, price)
		})

		t.Run("DenoteInUSD", func(t *testing.T) {
			usd, err := estimator.DenoteInUSD(
				GasPriceEstimatorCommit.denoteInUSDRequest.p,
				GasPriceEstimatorCommit.denoteInUSDRequest.wrappedNativePrice,
			)
			require.NoError(t, err)
			assert.Equal(t, GasPriceEstimatorCommit.denoteInUSDResponse.result, usd)
		})

		t.Run("Deviates", func(t *testing.T) {
			deviates, err := estimator.Deviates(
				GasPriceEstimatorCommit.deviatesRequest.p1,
				GasPriceEstimatorCommit.deviatesRequest.p2,
			)
			require.NoError(t, err)
			assert.Equal(t, GasPriceEstimatorCommit.deviatesResponse, deviates)
		})

		t.Run("Median", func(t *testing.T) {
			median, err := estimator.Median(GasPriceEstimatorCommit.medianRequest.gasPrices)
			require.NoError(t, err)
			assert.Equal(t, GasPriceEstimatorCommit.medianResponse, median)
		})
	})

	t.Run("GetAcceptedCommitReportGteTimestamp", func(t *testing.T) {
		report, err := client.GetAcceptedCommitReportsGteTimestamp(ctx,
			CommitStoreReader.getAcceptedCommitReportsGteTimestampRequest.timestamp,
			CommitStoreReader.getAcceptedCommitReportsGteTimestampRequest.confirmations)
		require.NoError(t, err)
		if !reflect.DeepEqual(CommitStoreReader.getAcceptedCommitReportsGteTimestampResponse, report) {
			t.Errorf("expected %v, got %v", CommitStoreReader.getAcceptedCommitReportsGteTimestampResponse, report)
		}
	})

	t.Run("GetCommitReportMatchingSeqNum", func(t *testing.T) {
		report, err := client.GetCommitReportMatchingSeqNum(ctx,
			CommitStoreReader.getCommitReportMatchingSeqNumRequest.seqNum,
			CommitStoreReader.getCommitReportMatchingSeqNumRequest.confirmations)
		require.NoError(t, err)
		// use the same response as the reportsGteTimestamp for simplicity
		if !reflect.DeepEqual(CommitStoreReader.getAcceptedCommitReportsGteTimestampResponse, report) {
			t.Errorf("expected %v, got %v", CommitStoreReader.getAcceptedCommitReportsGteTimestampRequest, report)
		}
	})

	t.Run("GetCommitStoreStaticConfig", func(t *testing.T) {
		config, err := client.GetCommitStoreStaticConfig(ctx)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.getCommitStoreStaticConfigResponse, config)
	})

	t.Run("GetExpectedNextSequenceNumber", func(t *testing.T) {
		seq, err := client.GetExpectedNextSequenceNumber(ctx)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.getExpectedNextSequenceNumberResponse, seq)
	})

	t.Run("GetLatestPriceEpochAndRound", func(t *testing.T) {
		got, err := client.GetLatestPriceEpochAndRound(ctx)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.getLatestPriceEpochAndRoundResponse, got)
	})

	t.Run("IsBlessed", func(t *testing.T) {
		got, err := client.IsBlessed(ctx, CommitStoreReader.isBlessedRequest)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.isBlessedResponse, got)
	})

	t.Run("IsDestChainHealthy", func(t *testing.T) {
		got, err := client.IsDestChainHealthy(ctx)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.isDestChainHealthyResponse, got)
	})

	t.Run("IsDown", func(t *testing.T) {
		got, err := client.IsDown(ctx)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.isDownResponse, got)
	})

	t.Run("OffchainConfig", func(t *testing.T) {
		config, err := client.OffchainConfig(ctx)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.offchainConfigResponse, config)
	})

	t.Run("VerifyExecutionReport", func(t *testing.T) {
		got, err := client.VerifyExecutionReport(ctx, CommitStoreReader.verifyExecutionReportRequest)
		require.NoError(t, err)
		assert.Equal(t, CommitStoreReader.verifyExecutionReportResponse, got)
	})
}
