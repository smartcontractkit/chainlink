package smoke

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"math/big"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/ocr2vrf_actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/ocr2vrf_actions/ocr2vrf_constants"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func TestOCR2VRFRedeemModel(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	testEnvironment, testNetwork := setupOCR2VRFEnvironment(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err, "Error connecting to blockchain")
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Error building contract deployer")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	chainClient.ParallelTransactions(true)

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Error deploying LINK token")

	mockETHLinkFeed, err := contractDeployer.DeployMockETHLINKFeed(ocr2vrf_constants.LinkEthFeedResponse)
	require.NoError(t, err, "Error deploying Mock ETH/LINK Feed")

	_, _, vrfBeaconContract, consumerContract, subID := ocr2vrf_actions.SetupOCR2VRFUniverse(
		t,
		linkToken,
		mockETHLinkFeed,
		contractDeployer,
		chainClient,
		nodeAddresses,
		chainlinkNodes,
		testNetwork,
	)

	//Request and Redeem Randomness
	requestID := ocr2vrf_actions.RequestAndRedeemRandomness(
		t,
		consumerContract,
		chainClient,
		vrfBeaconContract,
		ocr2vrf_constants.NumberOfRandomWordsToRequest,
		subID,
		ocr2vrf_constants.ConfirmationDelay,
		ocr2vrf_constants.RandomnessRedeemTransmissionEventTimeout,
	)

	for i := uint16(0); i < ocr2vrf_constants.NumberOfRandomWordsToRequest; i++ {
		randomness, err := consumerContract.GetRandomnessByRequestId(nil, requestID, big.NewInt(int64(i)))
		require.NoError(t, err)
		l.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness retrieved from Consumer contract")
		require.NotEqual(t, 0, randomness.Uint64(), "Randomness retrieved from Consumer contract give an answer other than 0")
	}
}

func TestOCR2VRFFulfillmentModel(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	testEnvironment, testNetwork := setupOCR2VRFEnvironment(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err, "Error connecting to blockchain")
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Error building contract deployer")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	chainClient.ParallelTransactions(true)

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Error deploying LINK token")

	mockETHLinkFeed, err := contractDeployer.DeployMockETHLINKFeed(ocr2vrf_constants.LinkEthFeedResponse)
	require.NoError(t, err, "Error deploying Mock ETH/LINK Feed")

	_, _, vrfBeaconContract, consumerContract, subID := ocr2vrf_actions.SetupOCR2VRFUniverse(
		t,
		linkToken,
		mockETHLinkFeed,
		contractDeployer,
		chainClient,
		nodeAddresses,
		chainlinkNodes,
		testNetwork,
	)

	requestID := ocr2vrf_actions.RequestRandomnessFulfillmentAndWaitForFulfilment(
		t,
		consumerContract,
		chainClient,
		vrfBeaconContract,
		ocr2vrf_constants.NumberOfRandomWordsToRequest,
		subID,
		ocr2vrf_constants.ConfirmationDelay,
		ocr2vrf_constants.RandomnessFulfilmentTransmissionEventTimeout,
	)

	for i := uint16(0); i < ocr2vrf_constants.NumberOfRandomWordsToRequest; i++ {
		randomness, err := consumerContract.GetRandomnessByRequestId(nil, requestID, big.NewInt(int64(i)))
		require.NoError(t, err, "Error getting Randomness result from Consumer Contract")
		l.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness Fulfillment retrieved from Consumer contract")
		require.NotEqual(t, 0, randomness.Uint64(), "Randomness Fulfillment retrieved from Consumer contract give an answer other than 0")
	}
}

func setupOCR2VRFEnvironment(t *testing.T) (testEnvironment *environment.Environment, testNetwork blockchain.EVMNetwork) {
	testNetwork = networks.SelectedNetwork
	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}

	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-ocr2vrf-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "6",
			"toml": client.AddNetworkDetailedConfig(
				config.BaseOCR2VRFTomlConfig,
				config.DefaultOCR2VRFNetworkDetailTomlConfig,
				testNetwork,
			),
		}))
	err := testEnvironment.Run()

	require.NoError(t, err, "Error running test environment")

	return testEnvironment, testNetwork
}
