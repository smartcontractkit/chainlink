package smoke

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/ocr2vrf_actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/ocr2vrf_actions/ocr2vrf_constants"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
)

var ocr2vrfSmokeConfig *testconfig.TestConfig

func TestOCR2VRFRedeemModel(t *testing.T) {
	t.Parallel()
	t.Skip("VRFv3 is on pause, skipping")
	l := logging.GetTestLogger(t)
	config, err := testconfig.GetConfig("Smoke", testconfig.OCR2)
	if err != nil {
		t.Fatal(err)
	}

	testEnvironment, testNetwork := setupOCR2VRFEnvironment(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	testNetwork = utils.MustReplaceSimulatedNetworkUrlWithK8(l, testNetwork, *testEnvironment)
	chainClient, err := actions_seth.GetChainClientWithConfigFunction(config, testNetwork, actions_seth.OneEphemeralKeysLiveTestnetCheckFn)
	require.NoError(t, err, "Error creating seth client")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	t.Cleanup(func() {
		err := actions_seth.TeardownSuite(t, chainClient, testEnvironment, chainlinkNodes, nil, zapcore.ErrorLevel, &config)
		require.NoError(t, err, "Error tearing down environment")
	})

	linkToken, err := contracts.DeployLinkTokenContract(l, chainClient)
	require.NoError(t, err, "Error deploying LINK token")

	mockETHLinkFeed, err := contracts.DeployMockETHLINKFeed(chainClient, ocr2vrf_constants.LinkEthFeedResponse)
	require.NoError(t, err, "Error deploying Mock ETH/LINK Feed")

	_, _, vrfBeaconContract, consumerContract, subID := ocr2vrf_actions.SetupOCR2VRFUniverse(
		t,
		linkToken,
		mockETHLinkFeed,
		chainClient,
		nodeAddresses,
		chainlinkNodes,
		testNetwork,
	)

	//Request and Redeem Randomness
	requestID := ocr2vrf_actions.RequestAndRedeemRandomness(
		t,
		consumerContract,
		vrfBeaconContract,
		ocr2vrf_constants.NumberOfRandomWordsToRequest,
		subID,
		ocr2vrf_constants.ConfirmationDelay,
		ocr2vrf_constants.RandomnessRedeemTransmissionEventTimeout,
	)

	for i := uint16(0); i < ocr2vrf_constants.NumberOfRandomWordsToRequest; i++ {
		randomness, err := consumerContract.GetRandomnessByRequestId(testcontext.Get(t), requestID, big.NewInt(int64(i)))
		require.NoError(t, err)
		l.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness retrieved from Consumer contract")
		require.NotEqual(t, 0, randomness.Uint64(), "Randomness retrieved from Consumer contract give an answer other than 0")
	}
}

func TestOCR2VRFFulfillmentModel(t *testing.T) {
	t.Parallel()
	t.Skip("VRFv3 is on pause, skipping")
	l := logging.GetTestLogger(t)
	config, err := testconfig.GetConfig("Smoke", testconfig.OCR2)
	if err != nil {
		t.Fatal(err)
	}

	testEnvironment, testNetwork := setupOCR2VRFEnvironment(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	testNetwork = utils.MustReplaceSimulatedNetworkUrlWithK8(l, testNetwork, *testEnvironment)
	chainClient, err := actions_seth.GetChainClientWithConfigFunction(config, testNetwork, actions_seth.OneEphemeralKeysLiveTestnetCheckFn)
	require.NoError(t, err, "Error creating seth client")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	t.Cleanup(func() {
		err := actions_seth.TeardownSuite(t, chainClient, testEnvironment, chainlinkNodes, nil, zapcore.ErrorLevel, &config)
		require.NoError(t, err, "Error tearing down environment")
	})

	linkToken, err := contracts.DeployLinkTokenContract(l, chainClient)
	require.NoError(t, err, "Error deploying LINK token")

	mockETHLinkFeed, err := contracts.DeployMockETHLINKFeed(chainClient, ocr2vrf_constants.LinkEthFeedResponse)
	require.NoError(t, err, "Error deploying Mock ETH/LINK Feed")

	_, _, vrfBeaconContract, consumerContract, subID := ocr2vrf_actions.SetupOCR2VRFUniverse(
		t,
		linkToken,
		mockETHLinkFeed,
		chainClient,
		nodeAddresses,
		chainlinkNodes,
		testNetwork,
	)

	requestID := ocr2vrf_actions.RequestRandomnessFulfillmentAndWaitForFulfilment(
		t,
		consumerContract,
		vrfBeaconContract,
		ocr2vrf_constants.NumberOfRandomWordsToRequest,
		subID,
		ocr2vrf_constants.ConfirmationDelay,
		ocr2vrf_constants.RandomnessFulfilmentTransmissionEventTimeout,
	)

	for i := uint16(0); i < ocr2vrf_constants.NumberOfRandomWordsToRequest; i++ {
		randomness, err := consumerContract.GetRandomnessByRequestId(testcontext.Get(t), requestID, big.NewInt(int64(i)))
		require.NoError(t, err, "Error getting Randomness result from Consumer Contract")
		l.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness Fulfillment retrieved from Consumer contract")
		require.NotEqual(t, 0, randomness.Uint64(), "Randomness Fulfillment retrieved from Consumer contract give an answer other than 0")
	}
}

func setupOCR2VRFEnvironment(t *testing.T) (testEnvironment *environment.Environment, testNetwork blockchain.EVMNetwork) {
	if ocr2vrfSmokeConfig == nil {
		c, err := testconfig.GetConfig("Smoke", testconfig.OCR2VRF)
		if err != nil {
			t.Fatal(err)
		}
		ocr2vrfSmokeConfig = &c
	}

	testNetwork = networks.MustGetSelectedNetworkConfig(ocr2vrfSmokeConfig.Network)[0]
	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}

	var overrideFn = func(_ interface{}, target interface{}) {
		ctf_config.MustConfigOverrideChainlinkVersion(ocr2vrfSmokeConfig.GetChainlinkImageConfig(), target)
		ctf_config.MightConfigOverridePyroscopeKey(ocr2vrfSmokeConfig.GetPyroscopeConfig(), target)
	}

	cd := chainlink.NewWithOverride(0, map[string]interface{}{
		"replicas": 6,
		"toml": networks.AddNetworkDetailedConfig(
			config.BaseOCR2Config,
			ocr2vrfSmokeConfig.Pyroscope,
			config.DefaultOCR2VRFNetworkDetailTomlConfig,
			testNetwork,
		),
	}, ocr2vrfSmokeConfig.ChainlinkImage, overrideFn)

	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-ocr2vrf-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(evmConfig).
		AddHelm(cd)
	err := testEnvironment.Run()

	require.NoError(t, err, "Error running test environment")

	return testEnvironment, testNetwork
}
