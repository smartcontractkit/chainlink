package test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/ccip"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestStaticCommitGasEstimator(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	// ensure GasPriceEstimatorCommit fixture is self consistent
	assert.NoError(t, GasPriceEstimatorCommit.Evaluate(ctx, GasPriceEstimatorCommit))

	// ensure
	biffed := GasPriceEstimatorCommit
	biffed.deviatesResponse = !GasPriceEstimatorCommit.deviatesResponse
	assert.NotEqual(t, biffed.denoteInUSDRequest, GasPriceEstimatorCommit.deviatesResponse)
	assert.Error(t, GasPriceEstimatorCommit.Evaluate(ctx, biffed))
}

func TestGasPriceEstimatorCommitGRPC(t *testing.T) {
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
	defer testServer.Stop()
	// handle client close and server stop

	gasPriceEstimatorCommit := ccip.NewCommitGasEstimatorGRPCServer(GasPriceEstimatorCommit)

	ccippb.RegisterGasPriceEstimatorCommitServer(testServer, gasPriceEstimatorCommit)
	// start the server and shutdown handler
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, testServer.Serve(lis))
	}()

	// create a price registry client
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "failed to dial %s", addr)
	t.Cleanup(func() { conn.Close() })
	client := ccip.NewCommitGasEstimatorGRPCClient(conn)

	// test the client
	roundTripGasPriceEstimatorCommitTests(ctx, t, client)
}

// roundTripGasPriceEstimatorCommitTests tests the round trip of the client<->server.
// it should exercise all the methods of the client.
// do not add client.Close to this test, test that from the driver test
func roundTripGasPriceEstimatorCommitTests(ctx context.Context, t *testing.T, client *ccip.CommitGasEstimatorGRPCClient) {
	t.Run("GetGasPrice", func(t *testing.T) {
		price, err := client.GetGasPrice(ctx)
		require.NoError(t, err)
		assert.Equal(t, GasPriceEstimatorCommit.getGasPriceResponse, price)
	})

	t.Run("DenoteInUSD", func(t *testing.T) {
		usd, err := client.DenoteInUSD(
			GasPriceEstimatorCommit.denoteInUSDRequest.p,
			GasPriceEstimatorCommit.denoteInUSDRequest.wrappedNativePrice,
		)
		require.NoError(t, err)
		assert.Equal(t, GasPriceEstimatorCommit.denoteInUSDResponse.result, usd)
	})

	t.Run("Deviates", func(t *testing.T) {
		isDeviant, err := client.Deviates(
			GasPriceEstimatorCommit.deviatesRequest.p1,
			GasPriceEstimatorCommit.deviatesRequest.p2,
		)
		require.NoError(t, err)
		assert.Equal(t, GasPriceEstimatorCommit.deviatesResponse, isDeviant)
	})

	t.Run("Median", func(t *testing.T) {
		median, err := client.Median(GasPriceEstimatorCommit.medianRequest.gasPrices)
		require.NoError(t, err)
		assert.Equal(t, GasPriceEstimatorCommit.medianResponse, median)
	})
}
