package test

import (
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

	scaffold := looptest.NewGRPCScaffold(t, setupCommitGasEstimatorServer, setupCommitGasEstimatorClient)
	t.Cleanup(scaffold.Close)
	// test the client
	roundTripGasPriceEstimatorCommitTests(t, scaffold.Client())
}

// roundTripGasPriceEstimatorCommitTests tests the round trip of the client<->server.
// it should exercise all the methods of the client.
// do not add client.Close to this test, test that from the driver test
func roundTripGasPriceEstimatorCommitTests(t *testing.T, client *ccip.CommitGasEstimatorGRPCClient) {
	t.Run("GetGasPrice", func(t *testing.T) {
		price, err := client.GetGasPrice(tests.Context(t))
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

func setupCommitGasEstimatorServer(t *testing.T, s *grpc.Server, b *loopnet.BrokerExt) *ccip.CommitGasEstimatorGRPCServer {
	gasProvider := ccip.NewCommitGasEstimatorGRPCServer(GasPriceEstimatorCommit)
	ccippb.RegisterGasPriceEstimatorCommitServer(s, gasProvider)
	return gasProvider
}

// adapt the client constructor so we can use it with the grpc scaffold
func setupCommitGasEstimatorClient(b *loopnet.BrokerExt, conn grpc.ClientConnInterface) *ccip.CommitGasEstimatorGRPCClient {
	return ccip.NewCommitGasEstimatorGRPCClient(conn)
}
