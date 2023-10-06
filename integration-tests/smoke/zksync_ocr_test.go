package smoke

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/gauntlet"
	"github.com/smartcontractkit/chainlink/integration-tests/l2/zksync"
)

// WIP
func TestOCRZKSync(t *testing.T) {
	l := logging.GetTestLogger(t)

	testEnvironment, testNetwork := setupOCRTest(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	// Adding L1 URL to HTTPURLs
	testNetwork.HTTPURLs = append(testNetwork.HTTPURLs, os.Getenv("L1_RPC_URL"))

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment, l)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")

	// Gauntlet Setup
	zkClient, err := zksync.Setup(os.Getenv("ZK_SYNC_GOERLI_HTTP_URLS"), chainClient.GetDefaultWallet().PrivateKey(), chainClient)
	require.NoError(t, err, "Creating ZKSync client should not fail")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")

	var chainlinkClients []*client.ChainlinkClient
	for _, k8sClient := range chainlinkNodes {
		chainlinkClients = append(chainlinkClients, k8sClient.ChainlinkClient)
	}

	err = zkClient.CreateKeys(chainlinkClients)
	require.NoError(t, err, "Creating keys should not fail")

	err = zkClient.DeployContracts(chainClient, gauntlet.DefaultOcrContract(), gauntlet.DefaultOcrConfig(), l)
	require.NoError(t, err, "Deploying Contracts should not fail")

	mockServer, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Creating mockserver clients shouldn't fail")

	t.Cleanup(func() {
		err = actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.DebugLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})
	chainClient.ParallelTransactions(true)

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")
	ocrInstance := []contracts.OffchainAggregator{
		zkClient.OCRContract,
	}

	// Set Config
	transmitterAddresses, err := actions.ChainlinkNodeAddresses(chainlinkNodes[1:])
	if err != nil {
		require.NoError(t, err, "Error getting transmitters")
	}

	// Exclude the first node, which will be used as a bootstrapper
	err = ocrInstance[0].SetConfig(
		chainlinkNodes[1:],
		contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes[1:])),
		transmitterAddresses,
	)
	if err != nil {
		require.NoError(t, err, "Error setting config")
	}

	bootstrapNode, workerNodes := chainlinkNodes[0], chainlinkNodes[1:]

	err = actions.CreateOCRJobs(ocrInstance, bootstrapNode, workerNodes, 5, mockServer, "280")
	require.NoError(t, err)

	err = actions.StartNewRound(1, ocrInstance, chainClient, l)
	require.NoError(t, err)

	answer, err := ocrInstance[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

	err = actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstance, chainlinkNodes, mockServer)
	require.NoError(t, err)
	err = actions.StartNewRound(2, ocrInstance, chainClient, l)
	require.NoError(t, err)

	answer, err = ocrInstance[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}

var ocrEnvVars = map[string]any{}

func setupOCRTest(t *testing.T) (
	testEnvironment *environment.Environment,
	testNetwork blockchain.EVMNetwork,
) {
	testNetwork = networks.SelectedNetwork
	evmConfig := ethereum.New(nil)
	if !testNetwork.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
		// For if we end up using env vars
		ocrEnvVars["ETH_URL"] = testNetwork.URLs[0]
		ocrEnvVars["ETH_HTTP_URL"] = testNetwork.HTTPURLs[0]
		ocrEnvVars["ETH_CHAIN_ID"] = fmt.Sprint(testNetwork.ChainID)
	}
	chainlinkChart := chainlink.New(0, map[string]interface{}{
		"toml":     client.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, config.DefaultOCRNetworkDetailTomlConfig, testNetwork),
		"replicas": 6,
	})

	useEnvVars := strings.ToLower(os.Getenv("TEST_USE_ENV_VAR_CONFIG"))
	if useEnvVars == "true" {
		chainlinkChart = chainlink.NewVersioned(0, "0.0.11", map[string]any{
			"replicas": 6,
			"env":      ocrEnvVars,
		})
	}

	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-ocr-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlinkChart)
	err := testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")
	return testEnvironment, testNetwork
}
