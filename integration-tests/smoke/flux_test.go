package smoke

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

func TestFluxBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithGeth().
		WithMockAdapter().
		WithCLNodes(3).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	nodeAddresses, err := env.ClCluster.NodeAddresses()
	require.NoError(t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	env.EVMClient.ParallelTransactions(true)

	adapterUUID := uuid.NewString()
	adapterPath := fmt.Sprintf("/variable-%s", adapterUUID)
	err = env.MockAdapter.SetAdapterBasedIntValuePath(adapterPath, []string{http.MethodPost}, 1e5)
	require.NoError(t, err, "Setting mock adapter value path shouldn't fail")

	lt, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
	fluxInstance, err := env.ContractDeployer.DeployFluxAggregatorContract(lt.Address(), contracts.DefaultFluxAggregatorOptions())
	require.NoError(t, err, "Deploying Flux Aggregator Contract shouldn't fail")
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for deployment of flux aggregator contract")

	err = lt.Transfer(fluxInstance.Address(), big.NewInt(1e18))
	require.NoError(t, err, "Funding Flux Aggregator Contract shouldn't fail")
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for funding of flux aggregator contract")

	err = fluxInstance.UpdateAvailableFunds()
	require.NoError(t, err, "Updating the available funds on the Flux Aggregator Contract shouldn't fail")

	err = env.FundChainlinkNodes(big.NewFloat(1))
	require.NoError(t, err, "Failed to fund the nodes")

	err = fluxInstance.SetOracles(
		contracts.FluxAggregatorSetOraclesOptions{
			AddList:            nodeAddresses,
			RemoveList:         []common.Address{},
			AdminList:          nodeAddresses,
			MinSubmissions:     3,
			MaxSubmissions:     3,
			RestartDelayRounds: 0,
		})
	require.NoError(t, err, "Setting oracle options in the Flux Aggregator contract shouldn't fail")

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")
	oracles, err := fluxInstance.GetOracles(context.Background())
	require.NoError(t, err, "Getting oracle details from the Flux aggregator contract shouldn't fail")
	l.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("Oracles set")

	adapterFullURL := fmt.Sprintf("%s%s", env.MockAdapter.InternalEndpoint, adapterPath)
	l.Info().Str("AdapterFullURL", adapterFullURL).Send()
	bta := &client.BridgeTypeAttributes{
		Name: fmt.Sprintf("variable-%s", adapterUUID),
		URL:  adapterFullURL,
	}
	for i, n := range env.ClCluster.Nodes {
		err = n.API.MustCreateBridge(bta)
		require.NoError(t, err, "Creating bridge shouldn't fail for node %d", i+1)

		fluxSpec := &client.FluxMonitorJobSpec{
			Name:              fmt.Sprintf("flux-monitor-%s", adapterUUID),
			ContractAddress:   fluxInstance.Address(),
			EVMChainID:        env.EVMClient.GetChainID().String(),
			Threshold:         0,
			AbsoluteThreshold: 0,
			PollTimerPeriod:   15 * time.Second, // min 15s
			IdleTimerDisabled: true,
			ObservationSource: client.ObservationSourceSpecBridge(bta),
		}
		_, err = n.API.MustCreateJob(fluxSpec)
		require.NoError(t, err, "Creating flux job shouldn't fail for node %d", i+1)
	}

	// initial value set is performed before jobs creation
	fluxRoundTimeout := 1 * time.Minute
	fluxRound := contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(1), fluxRoundTimeout, l)
	env.EVMClient.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")
	data, err := fluxInstance.GetContractData(context.Background())
	require.NoError(t, err, "Getting contract data from flux aggregator contract shouldn't fail")
	require.Equal(t, int64(1e5), data.LatestRoundData.Answer.Int64(),
		"Expected latest round answer to be %d, but found %d", int64(1e5), data.LatestRoundData.Answer.Int64())
	require.Equal(t, int64(1), data.LatestRoundData.RoundId.Int64(),
		"Expected latest round id to be %d, but found %d", int64(1), data.LatestRoundData.RoundId.Int64())
	require.Equal(t, int64(1), data.LatestRoundData.AnsweredInRound.Int64(),
		"Expected latest round's answered in round to be %d, but found %d", int64(1), data.LatestRoundData.AnsweredInRound.Int64())
	require.Equal(t, int64(999999999999999997), data.AvailableFunds.Int64(),
		"Expected available funds to be %d, but found %d", int64(999999999999999997), data.AvailableFunds.Int64())
	require.Equal(t, int64(3), data.AllocatedFunds.Int64(),
		"Expected allocated funds to be %d, but found %d", int64(3), data.AllocatedFunds.Int64())

	fluxRound = contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(2), fluxRoundTimeout, l)
	env.EVMClient.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
	err = env.MockAdapter.SetAdapterBasedIntValuePath(adapterPath, []string{http.MethodPost}, 1e10)
	require.NoError(t, err, "Setting value path in mock server shouldn't fail")
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")
	data, err = fluxInstance.GetContractData(context.Background())
	require.NoError(t, err, "Getting contract data from flux aggregator contract shouldn't fail")
	require.Equal(t, int64(1e10), data.LatestRoundData.Answer.Int64(),
		"Expected latest round answer to be %d, but found %d", int64(1e10), data.LatestRoundData.Answer.Int64())
	require.Equal(t, int64(2), data.LatestRoundData.RoundId.Int64(),
		"Expected latest round id to be %d, but found %d", int64(2), data.LatestRoundData.RoundId.Int64())
	require.Equal(t, int64(999999999999999994), data.AvailableFunds.Int64(),
		"Expected available funds to be %d, but found %d", int64(999999999999999994), data.AvailableFunds.Int64())
	require.Equal(t, int64(6), data.AllocatedFunds.Int64(),
		"Expected allocated funds to be %d, but found %d", int64(6), data.AllocatedFunds.Int64())
	l.Info().Interface("data", data).Msg("Round data")

	for _, oracleAddr := range nodeAddresses {
		payment, _ := fluxInstance.WithdrawablePayment(context.Background(), oracleAddr)
		require.Equal(t, int64(2), payment.Int64(),
			"Expected flux aggregator contract's withdrawable payment to be %d, but found %d", int64(2), payment.Int64())
	}
}
