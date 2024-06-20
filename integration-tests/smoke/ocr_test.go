package smoke

import (
	"math/big"
	"testing"
	"time"

	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/seth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

const (
	ErrWatchingNewOCRRound = "Error watching for new OCR round"
)

func TestOCRBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, ocrInstances, sethClient := prepareORCv1SmokeTestEnv(t, l, 5)
	nodeClients := env.ClCluster.NodeAPIs()
	workerNodes := nodeClients[1:]

	err := actions.SetAllAdapterResponsesToTheSameValueLocal(10, ocrInstances, workerNodes, env.MockAdapter)
	require.NoError(t, err, "Error setting all adapter responses to the same value")

	err = actions.WatchNewOCRRound(l, sethClient, 2, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(3*time.Minute))
	require.NoError(t, err, ErrWatchingNewOCRRound)

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}

func TestOCRJobReplacement(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, ocrInstances, sethClient := prepareORCv1SmokeTestEnv(t, l, 5)
	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	err := actions.SetAllAdapterResponsesToTheSameValueLocal(10, ocrInstances, workerNodes, env.MockAdapter)
	require.NoError(t, err, "Error setting all adapter responses to the same value")
	err = actions.WatchNewOCRRound(l, sethClient, 2, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(3*time.Minute))
	require.NoError(t, err, ErrWatchingNewOCRRound)

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())

	err = actions.DeleteJobs(nodeClients)
	require.NoError(t, err, "Error deleting OCR jobs")

	err = actions.DeleteBridges(nodeClients)
	require.NoError(t, err, "Error deleting OCR bridges")

	//Recreate job
	err = actions.CreateOCRJobsLocal(ocrInstances, bootstrapNode, workerNodes, 5, env.MockAdapter, big.NewInt(sethClient.ChainID))
	require.NoError(t, err, "Error creating OCR jobs")

	err = actions.WatchNewOCRRound(l, sethClient, 1, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(3*time.Minute))
	require.NoError(t, err, ErrWatchingNewOCRRound)

	answer, err = ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}

func prepareORCv1SmokeTestEnv(t *testing.T, l zerolog.Logger, firstRoundResult int64) (*test_env.CLClusterTestEnv, []contracts.OffchainAggregator, *seth.Client) {
	config, err := tc.GetConfig("Smoke", tc.OCR)
	if err != nil {
		t.Fatal(err)
	}

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(network.EthereumNetworkConfig).
		WithMockAdapter().
		WithCLNodes(6).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	evmNetwork, err := env.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	sethClient, err := seth_utils.GetChainClient(config, *evmNetwork)
	require.NoError(t, err, "Error getting seth client")

	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	err = actions.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes), big.NewFloat(*config.Common.ChainlinkNodeFunding))
	require.NoError(t, err, "Error funding Chainlink nodes")

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		_ = actions.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()))
	})

	linkContract, err := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, err, "Error deploying link token contract")

	ocrInstances, err := actions.DeployOCRv1Contracts(l, sethClient, 1, common.HexToAddress(linkContract.Address()), contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes))
	require.NoError(t, err, "Error deploying OCR contracts")

	err = actions.CreateOCRJobsLocal(ocrInstances, bootstrapNode, workerNodes, 5, env.MockAdapter, big.NewInt(sethClient.ChainID))
	require.NoError(t, err, "Error creating OCR jobs")

	err = actions.WatchNewOCRRound(l, sethClient, 1, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(3*time.Minute))
	require.NoError(t, err, "Error watching for new OCR round")

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, firstRoundResult, answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

	return env, ocrInstances, sethClient
}
