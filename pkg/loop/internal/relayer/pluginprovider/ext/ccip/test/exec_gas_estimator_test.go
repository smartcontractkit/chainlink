package test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	loopnet "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ccip"
	looptest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
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

	scaffold := looptest.NewGRPCScaffold(t, setupExecGasEstimatorServer, setupExecGasEstimatorClient)
	t.Cleanup(scaffold.Close)
	roundTripGasPriceEstimatorExecTests(t, scaffold.Client())
}

// roundTripGasPriceEstimatorExecTests tests the round trip of the client<->server.
// it should exercise all the methods of the client.
// do not add client.Close to this test, test that from the driver test
func roundTripGasPriceEstimatorExecTests(t *testing.T, client *ccip.ExecGasEstimatorGRPCClient) {
	t.Run("GetGasPrice", func(t *testing.T) {
		price, err := client.GetGasPrice(tests.Context(t))
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

func setupExecGasEstimatorServer(t *testing.T, s *grpc.Server, b *loopnet.BrokerExt) *ccip.ExecGasEstimatorGRPCServer {
	gasProvider := ccip.NewExecGasEstimatorGRPCServer(GasPriceEstimatorExec)
	ccippb.RegisterGasPriceEstimatorExecServer(s, gasProvider)
	return gasProvider
}

// adapt the client constructor so we can use it with the grpc scaffold
func setupExecGasEstimatorClient(b *loopnet.BrokerExt, conn grpc.ClientConnInterface) *ccip.ExecGasEstimatorGRPCClient {
	return ccip.NewExecGasEstimatorGRPCClient(conn)
}

var _ looptest.SetupGRPCServer[*ccip.ExecGasEstimatorGRPCServer] = setupExecGasEstimatorServer
var _ looptest.SetupGRPCClient[*ccip.ExecGasEstimatorGRPCClient] = setupExecGasEstimatorClient
