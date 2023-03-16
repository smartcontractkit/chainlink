package smoke

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func TestFluxBasic(t *testing.T) {
	t.Parallel()
	testEnvironment, testNetwork := setupFluxTest(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
	mockServer, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Creating mock server client shouldn't fail")
	// Register cleanup
	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	chainClient.ParallelTransactions(true)

	adapterUUID := uuid.NewV4().String()
	adapterPath := fmt.Sprintf("/variable-%s", adapterUUID)
	err = mockServer.SetValuePath(adapterPath, 1e5)
	require.NoError(t, err, "Setting mockserver value path shouldn't fail")

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
	fluxInstance, err := contractDeployer.DeployFluxAggregatorContract(linkToken.Address(), contracts.DefaultFluxAggregatorOptions())
	require.NoError(t, err, "Deploying Flux Aggregator Contract shouldn't fail")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for deployment of flux aggregator contract")

	err = linkToken.Transfer(fluxInstance.Address(), big.NewInt(1e18))
	require.NoError(t, err, "Funding Flux Aggregator Contract shouldn't fail")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for funding of flux aggregator contract")

	err = fluxInstance.UpdateAvailableFunds()
	require.NoError(t, err, "Updating the available funds on the Flux Aggregator Contract shouldn't fail")

	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.02))
	require.NoError(t, err, "Funding chainlink nodes with ETH shouldn't fail")

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
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")
	oracles, err := fluxInstance.GetOracles(context.Background())
	require.NoError(t, err, "Getting oracle details from the Flux aggregator contract shouldn't fail")
	log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("Oracles set")

	adapterFullURL := fmt.Sprintf("%s%s", mockServer.Config.ClusterURL, adapterPath)
	bta := client.BridgeTypeAttributes{
		Name: fmt.Sprintf("variable-%s", adapterUUID),
		URL:  adapterFullURL,
	}
	for i, n := range chainlinkNodes {
		err = n.MustCreateBridge(&bta)
		require.NoError(t, err, "Creating bridge shouldn't fail for node %d", i+1)

		fluxSpec := &client.FluxMonitorJobSpec{
			Name:              fmt.Sprintf("flux-monitor-%s", adapterUUID),
			ContractAddress:   fluxInstance.Address(),
			Threshold:         0,
			AbsoluteThreshold: 0,
			PollTimerPeriod:   15 * time.Second, // min 15s
			IdleTimerDisabled: true,
			ObservationSource: client.ObservationSourceSpecBridge(bta),
		}
		_, err = n.MustCreateJob(fluxSpec)
		require.NoError(t, err, "Creating flux job shouldn't fail for node %d", i+1)
	}

	// initial value set is performed before jobs creation
	fluxRoundTimeout := 2 * time.Minute
	fluxRound := contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(1), fluxRoundTimeout)
	chainClient.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")
	data, err := fluxInstance.GetContractData(context.Background())
	require.NoError(t, err, "Getting contract data from flux aggregator contract shouldn't fail")
	log.Info().Interface("Data", data).Msg("Round data")
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

	fluxRound = contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(2), fluxRoundTimeout)
	chainClient.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
	err = mockServer.SetValuePath(adapterPath, 1e10)
	require.NoError(t, err, "Setting value path in mock server shouldn't fail")
	err = chainClient.WaitForEvents()
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
	log.Info().Interface("data", data).Msg("Round data")

	for _, oracleAddr := range nodeAddresses {
		payment, _ := fluxInstance.WithdrawablePayment(context.Background(), oracleAddr)
		require.Equal(t, int64(2), payment.Int64(),
			"Expected flux aggregator contract's withdrawable payment to be %d, but found %d", int64(2), payment.Int64())
	}
}

func setupFluxTest(t *testing.T) (testEnvironment *environment.Environment, testNetwork blockchain.EVMNetwork) {
	testNetwork = networks.SelectedNetwork
	evmConf := ethereum.New(nil)
	if !testNetwork.Simulated {
		evmConf = ethereum.New(&ethereum.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}
	baseTOML := `[OCR]
Enabled = true`
	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-flux-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConf).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"toml":     client.AddNetworksConfig(baseTOML, testNetwork),
			"replicas": 3,
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")
	return testEnvironment, testNetwork
}
