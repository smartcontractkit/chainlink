package smoke

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

type ocr2test struct {
	name                string
	env                 map[string]string
	chainReaderAndCodec bool
}

func defaultTestData() ocr2test {
	return ocr2test{
		name:                "n/a",
		env:                 make(map[string]string),
		chainReaderAndCodec: false,
	}
}

// Tests a basic OCRv2 median feed
func TestOCRv2Basic(t *testing.T) {
	t.Parallel()

	noMedianPlugin := map[string]string{string(env.MedianPlugin.Cmd): ""}
	medianPlugin := map[string]string{string(env.MedianPlugin.Cmd): "chainlink-feeds"}
	for _, test := range []ocr2test{
		{"legacy", noMedianPlugin, false},
		{"legacy-chain-reader", noMedianPlugin, true},
		{"plugins", medianPlugin, false},
		{"plugins-chain-reader", medianPlugin, true},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)

			testEnv, aggregatorContracts, sethClient := prepareORCv2SmokeTestEnv(t, test, l, 5)

			err := testEnv.MockAdapter.SetAdapterBasedIntValuePath("ocr2", []string{http.MethodGet, http.MethodPost}, 10)
			require.NoError(t, err)
			err = actions_seth.WatchNewOCRRound(l, sethClient, 2, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*5)
			require.NoError(t, err)

			roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(2))
			require.NoError(t, err, "Error getting latest OCR answer")
			require.Equal(t, int64(10), roundData.Answer.Int64(),
				"Expected latest answer from OCR contract to be 10 but got %d",
				roundData.Answer.Int64(),
			)
		})
	}
}

// Tests that just calling requestNewRound() will properly induce more rounds
func TestOCRv2Request(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	_, aggregatorContracts, sethClient := prepareORCv2SmokeTestEnv(t, defaultTestData(), l, 5)

	// Keep the mockserver value the same and continually request new rounds
	for round := 2; round <= 4; round++ {
		err := actions_seth.StartNewRound(contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts))
		require.NoError(t, err, "Error starting new OCR2 round")
		err = actions_seth.WatchNewOCRRound(l, sethClient, int64(round), contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*5)
		require.NoError(t, err, "Error watching for new OCR2 round")
		roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(int64(round)))
		require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
		require.Equal(t, int64(5), roundData.Answer.Int64(),
			"Expected round %d answer from OCR contract to be 5 but got %d",
			round,
			roundData.Answer.Int64(),
		)
	}
}

func TestOCRv2JobReplacement(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, aggregatorContracts, sethClient := prepareORCv2SmokeTestEnv(t, defaultTestData(), l, 5)
	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	err := env.MockAdapter.SetAdapterBasedIntValuePath("ocr2", []string{http.MethodGet, http.MethodPost}, 10)
	require.NoError(t, err)
	err = actions_seth.WatchNewOCRRound(l, sethClient, 2, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*5)
	require.NoError(t, err, "Error watching for new OCR2 round")

	roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(2))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 10 but got %d",
		roundData.Answer.Int64(),
	)

	err = actions.DeleteJobs(nodeClients)
	require.NoError(t, err)

	err = actions.DeleteBridges(nodeClients)
	require.NoError(t, err)

	err = actions.CreateOCRv2JobsLocal(aggregatorContracts, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 15, uint64(sethClient.ChainID), false, false)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	err = actions_seth.WatchNewOCRRound(l, sethClient, 3, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*3)
	require.NoError(t, err, "Error watching for new OCR2 round")

	roundData, err = aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(3))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(15), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 15 but got %d",
		roundData.Answer.Int64(),
	)
}

func prepareORCv2SmokeTestEnv(t *testing.T, testData ocr2test, l zerolog.Logger, firstRoundResult int) (*test_env.CLClusterTestEnv, []contracts.OffchainAggregatorV2, *seth.Client) {
	config, err := tc.GetConfig("Smoke", tc.OCR2)
	if err != nil {
		t.Fatal(err)
	}

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	clNodeCount := 6

	testEnv, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithMockAdapter().
		WithCLNodes(clNodeCount).
		WithCLNodeOptions(test_env.WithNodeEnvVars(testData.env)).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithStandardCleanup().
		WithSeth().
		Build()
	require.NoError(t, err)

	selectedNetwork := networks.MustGetSelectedNetworkConfig(config.Network)[0]
	sethClient, err := testEnv.GetSethClient(selectedNetwork.ChainID)
	require.NoError(t, err, "Error getting seth client")

	nodeClients := testEnv.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	linkContract, err := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, err, "Error deploying link token contract")

	err = actions_seth.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes), big.NewFloat(.05))
	require.NoError(t, err, "Error funding Chainlink nodes")

	// Gather transmitters
	var transmitters []string
	for _, node := range workerNodes {
		addr, err := node.PrimaryEthAddress()
		if err != nil {
			require.NoError(t, fmt.Errorf("error getting node's primary ETH address: %w", err))
		}
		transmitters = append(transmitters, addr)
	}

	ocrOffchainOptions := contracts.DefaultOffChainAggregatorOptions()
	aggregatorContracts, err := actions_seth.DeployOCRv2Contracts(l, sethClient, 1, common.HexToAddress(linkContract.Address()), transmitters, ocrOffchainOptions)
	require.NoError(t, err, "Error deploying OCRv2 aggregator contracts")

	err = actions.CreateOCRv2JobsLocal(aggregatorContracts, bootstrapNode, workerNodes, testEnv.MockAdapter, "ocr2", 5, uint64(sethClient.ChainID), false, testData.chainReaderAndCodec)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	ocrv2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
	require.NoError(t, err, "Error building OCRv2 config")

	err = actions_seth.ConfigureOCRv2AggregatorContracts(ocrv2Config, aggregatorContracts)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

	assertCorrectNodeConfiguration(t, l, clNodeCount, testData, testEnv)

	err = actions_seth.WatchNewOCRRound(l, sethClient, 1, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*5)
	require.NoError(t, err, "Error watching for new OCR2 round")
	roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(1))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(firstRoundResult), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 5 but got %d",
		roundData.Answer.Int64(),
	)

	return testEnv, aggregatorContracts, sethClient
}

func assertCorrectNodeConfiguration(t *testing.T, l zerolog.Logger, totalNodeCount int, testData ocr2test, testEnv *test_env.CLClusterTestEnv) {
	expectedNodesWithConfiguration := totalNodeCount - 1 // minus bootstrap node
	expectedPatterns := []string{}

	if testData.env[string(env.MedianPlugin.Cmd)] != "" {
		expectedPatterns = append(expectedPatterns, "Registered loopp.*OCR2.*Median.*")
	}

	if testData.chainReaderAndCodec {
		expectedPatterns = append(expectedPatterns, "relayConfig\\.chainReader")
	} else {
		expectedPatterns = append(expectedPatterns, "ChainReader missing from RelayConfig; falling back to internal MedianContract")
	}

	// make sure that nodes are correctly configured by scanning the logs
	for _, pattern := range expectedPatterns {
		l.Info().Msgf("Checking for pattern: '%s' in CL node logs", pattern)
		var correctlyConfiguredNodes []string
		for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
			logProcessor, processFn, err := logstream.GetRegexMatchingProcessor(testEnv.LogStream, pattern)
			require.NoError(t, err, "Error getting regex matching processor")

			count, err := logProcessor.ProcessContainerLogs(testEnv.ClCluster.Nodes[i].ContainerName, processFn)
			require.NoError(t, err, "Error processing container logs")
			if *count >= 1 {
				correctlyConfiguredNodes = append(correctlyConfiguredNodes, testEnv.ClCluster.Nodes[i].ContainerName)
			}
		}
		require.Equal(t, expectedNodesWithConfiguration, len(correctlyConfiguredNodes), "expected correct plugin config to be applied to %d cl-nodes, but only following ones had it: %s; regexp used: %s", expectedNodesWithConfiguration, strings.Join(correctlyConfiguredNodes, ", "), string(pattern))
	}
}
