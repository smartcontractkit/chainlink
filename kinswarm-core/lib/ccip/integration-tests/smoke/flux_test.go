package smoke

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestFluxBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig([]string{"Smoke"}, tc.Flux)
	require.NoError(t, err, "Error getting config")

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithMockAdapter().
		WithCLNodes(3).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	nodeAddresses, err := env.ClCluster.NodeAddresses()
	require.NoError(t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")

	evmNetwork, err := env.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	sethClient, err := utils.TestAwareSethClient(t, config, evmNetwork)
	require.NoError(t, err, "Error getting seth client")

	adapterUUID := uuid.NewString()
	adapterPath := fmt.Sprintf("/variable-%s", adapterUUID)
	err = env.MockAdapter.SetAdapterBasedIntValuePath(adapterPath, []string{http.MethodPost}, 1e5)
	require.NoError(t, err, "Setting mock adapter value path shouldn't fail")

	lt, err := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	fluxInstance, err := contracts.DeployFluxAggregatorContract(sethClient, lt.Address(), contracts.DefaultFluxAggregatorOptions())
	require.NoError(t, err, "Deploying Flux Aggregator Contract shouldn't fail")

	err = lt.Transfer(fluxInstance.Address(), big.NewInt(1e18))
	require.NoError(t, err, "Funding Flux Aggregator Contract shouldn't fail")

	err = fluxInstance.UpdateAvailableFunds()
	require.NoError(t, err, "Updating the available funds on the Flux Aggregator Contract shouldn't fail")

	err = actions.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()), big.NewFloat(*config.Common.ChainlinkNodeFunding))
	require.NoError(t, err, "Failed to fund the nodes")

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		_ = actions.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()))
	})

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

	oracles, err := fluxInstance.GetOracles(testcontext.Get(t))
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
			EVMChainID:        fmt.Sprint(sethClient.ChainID),
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
	err = actions.WatchNewFluxRound(l, sethClient, 1, fluxInstance, fluxRoundTimeout)
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")
	data, err := fluxInstance.GetContractData(testcontext.Get(t))
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

	err = env.MockAdapter.SetAdapterBasedIntValuePath(adapterPath, []string{http.MethodPost}, 1e10)
	require.NoError(t, err, "Setting value path in mock server shouldn't fail")
	err = actions.WatchNewFluxRound(l, sethClient, 2, fluxInstance, fluxRoundTimeout)
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")
	data, err = fluxInstance.GetContractData(testcontext.Get(t))
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
		payment, _ := fluxInstance.WithdrawablePayment(testcontext.Get(t), oracleAddr)
		require.Equal(t, int64(2), payment.Int64(),
			"Expected flux aggregator contract's withdrawable payment to be %d, but found %d", int64(2), payment.Int64())
	}
}
