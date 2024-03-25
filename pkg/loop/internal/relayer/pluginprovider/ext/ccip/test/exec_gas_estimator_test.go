package test

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"sync"
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestStaticExecGasEstimator(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	// ensure GasPriceEstimatorExec fixture is self consistent
	assert.NoError(t, GasPriceEstimatorExec.Evaluate(ctx, GasPriceEstimatorExec))

	// ensure
	biffed := GasPriceEstimatorExec
	biffed.estimateMsgCostUSDResponse = big.NewInt(131)
	assert.NotEqual(t, biffed.estimateMsgCostUSDResponse, GasPriceEstimatorExec.estimateMsgCostUSDResponse)
	assert.Error(t, GasPriceEstimatorExec.Evaluate(ctx, biffed))
}

func TestGasPriceEstimatorExecGRPC(t *testing.T) {
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

	gasPriceEstimatorExec := ccip.NewExecGasEstimatorGRPCServer(GasPriceEstimatorExec)

	ccippb.RegisterGasPriceEstimatorExecServer(testServer, gasPriceEstimatorExec)
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
	client := ccip.NewExecGasEstimatorGRPCClient(conn)

	// test the client
	roundTripGasPriceEstimatorExecTests(ctx, t, client)
}

// roundTripGasPriceEstimatorExecTests tests the round trip of the client<->server.
// it should exercise all the methods of the client.
// do not add client.Close to this test, test that from the driver test
func roundTripGasPriceEstimatorExecTests(ctx context.Context, t *testing.T, client *ccip.ExecGasEstimatorGRPCClient) {
	t.Run("GetGasPrice", func(t *testing.T) {
		price, err := client.GetGasPrice(ctx)
		require.NoError(t, err)
		assert.Equal(t, GasPriceEstimatorExec.getGasPriceResponse, price)
	})

	t.Run("DenoteInUSD", func(t *testing.T) {
		usd, err := client.DenoteInUSD(
			GasPriceEstimatorExec.denoteInUSDRequest.p,
			GasPriceEstimatorExec.denoteInUSDRequest.wrappedNativePrice,
		)
		require.NoError(t, err)
		assert.Equal(t, GasPriceEstimatorExec.denoteInUSDResponse.result, usd)
	})

	t.Run("EstimateMsgCostUSD", func(t *testing.T) {
		cost, err := client.EstimateMsgCostUSD(
			GasPriceEstimatorExec.estimateMsgCostUSDRequest.p,
			GasPriceEstimatorExec.estimateMsgCostUSDRequest.wrappedNativePrice,
			GasPriceEstimatorExec.estimateMsgCostUSDRequest.msg,
		)
		require.NoError(t, err)
		assert.Equal(t, GasPriceEstimatorExec.estimateMsgCostUSDResponse, cost)
	})

	t.Run("Median", func(t *testing.T) {
		median, err := client.Median(GasPriceEstimatorExec.medianRequest.gasPrices)
		require.NoError(t, err)
		assert.Equal(t, GasPriceEstimatorExec.medianResponse, median)
	})
}
