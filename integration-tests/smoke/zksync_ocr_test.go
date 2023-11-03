package smoke

import (
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/L2/zksync"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

// WIP
func TestOCRZKSync(t *testing.T) {
	l := logging.GetTestLogger(t)

	testEnvironment, testNetwork, err := zksync.SetupOCRTest(t)
	require.NoError(t, err, "Deploying env should not fail")
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	l1RpcUrl, isSet := os.LookupEnv("L1_RPC_URL")
	require.Equal(t, isSet, true, "L1_RPC_URL should be defined")

	// Adding L1 URL to HTTPURLs
	testNetwork.HTTPURLs = append(testNetwork.HTTPURLs, l1RpcUrl)

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment, l)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")

	// Gauntlet Setup
	zkClient, err := zksync.Setup(os.Getenv("ZK_SYNC_GOERLI_HTTP_URLS"), chainClient.GetDefaultWallet().PrivateKey(), chainClient)
	require.NoError(t, err, "Creating ZKSync client should not fail")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")

	err = zkClient.DeployOCRFeed(testEnvironment, chainClient, chainlinkNodes, testNetwork, l)
	require.NoError(t, err, "Error deploying OCR FEED")

	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.DebugLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	answer, err := zkClient.RequestOCRRound(1, 10, l)
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())

	answer, err = zkClient.RequestOCRRound(2, 10, l)
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}
