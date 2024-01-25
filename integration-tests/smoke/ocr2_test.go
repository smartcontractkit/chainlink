package smoke

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
)

// Tests a basic OCRv2 median feed
func TestOCRv2Basic(t *testing.T) {
	t.Parallel()

	noMedianPlugin := map[string]string{string(env.MedianPlugin.Cmd): ""}
	medianPlugin := map[string]string{string(env.MedianPlugin.Cmd): "chainlink-feeds"}
	for _, test := range []struct {
		name                string
		env                 map[string]string
		chainReaderAndCodec bool
	}{
		{"legacy", noMedianPlugin, false},
		{"legacy-chain-reader", noMedianPlugin, true},
		{"plugins", medianPlugin, false},
		{"plugins-chain-reader", medianPlugin, true},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			l := logging.GetTestLogger(t)

			config, err := tc.GetConfig("Smoke", tc.OCR2)
			if err != nil {
				t.Fatal(err)
			}

			network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
			require.NoError(t, err, "Error building ethereum network config")

			env, err := test_env.NewCLTestEnvBuilder().
				WithTestInstance(t).
				WithTestConfig(&config).
				WithPrivateEthereumNetwork(network).
				WithMockAdapter().
				WithCLNodeConfig(node.NewConfig(node.NewBaseConfig(),
					node.WithOCR2(),
					node.WithP2Pv2(),
					node.WithTracing(),
				)).
				WithCLNodeOptions(test_env.WithNodeEnvVars(test.env)).
				WithCLNodes(6).
				WithFunding(big.NewFloat(.1)).
				WithStandardCleanup().
				Build()
			require.NoError(t, err)

			env.ParallelTransactions(true)

			nodeClients := env.ClCluster.NodeAPIs()
			bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

			linkToken, err := env.ContractDeployer.DeployLinkTokenContract()
			require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

			err = actions.FundChainlinkNodesLocal(workerNodes, env.EVMClient, big.NewFloat(.05))
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
			aggregatorContracts, err := actions.DeployOCRv2Contracts(1, linkToken, env.ContractDeployer, transmitters, env.EVMClient, ocrOffchainOptions)
			require.NoError(t, err, "Error deploying OCRv2 aggregator contracts")

			err = actions.CreateOCRv2JobsLocal(aggregatorContracts, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 5, env.EVMClient.GetChainID().Uint64(), false, test.chainReaderAndCodec)
			require.NoError(t, err, "Error creating OCRv2 jobs")

			ocrv2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
			require.NoError(t, err, "Error building OCRv2 config")

			err = actions.ConfigureOCRv2AggregatorContracts(env.EVMClient, ocrv2Config, aggregatorContracts)
			require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

			err = actions.WatchNewOCR2Round(1, aggregatorContracts, env.EVMClient, time.Minute*5, l)
			require.NoError(t, err, "Error watching for new OCR2 round")
			roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(1))
			require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
			require.Equal(t, int64(5), roundData.Answer.Int64(),
				"Expected latest answer from OCR contract to be 5 but got %d",
				roundData.Answer.Int64(),
			)

			err = env.MockAdapter.SetAdapterBasedIntValuePath("ocr2", []string{http.MethodGet, http.MethodPost}, 10)
			require.NoError(t, err)
			err = actions.WatchNewOCR2Round(2, aggregatorContracts, env.EVMClient, time.Minute*5, l)
			require.NoError(t, err)

			roundData, err = aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(2))
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

	config, err := tc.GetConfig("Smoke", tc.ForwarderOcr)
	if err != nil {
		t.Fatal(err)
	}

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(network).
		WithMockAdapter().
		WithCLNodeConfig(node.NewConfig(node.NewBaseConfig(),
			node.WithOCR2(),
			node.WithP2Pv2(),
			node.WithTracing(),
		)).
		WithCLNodes(6).
		WithFunding(big.NewFloat(.1)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	env.ParallelTransactions(true)

	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	linkToken, err := env.ContractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	err = actions.FundChainlinkNodesLocal(workerNodes, env.EVMClient, big.NewFloat(.05))
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
	aggregatorContracts, err := actions.DeployOCRv2Contracts(1, linkToken, env.ContractDeployer, transmitters, env.EVMClient, ocrOffchainOptions)
	require.NoError(t, err, "Error deploying OCRv2 aggregator contracts")

	err = actions.CreateOCRv2JobsLocal(aggregatorContracts, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 5, env.EVMClient.GetChainID().Uint64(), false, false)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	ocrv2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
	require.NoError(t, err, "Error building OCRv2 config")

	err = actions.ConfigureOCRv2AggregatorContracts(env.EVMClient, ocrv2Config, aggregatorContracts)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

	err = actions.WatchNewOCR2Round(1, aggregatorContracts, env.EVMClient, time.Minute*5, l)
	require.NoError(t, err, "Error watching for new OCR2 round")
	roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(1))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 5 but got %d",
		roundData.Answer.Int64(),
	)

	// Keep the mockserver value the same and continually request new rounds
	for round := 2; round <= 4; round++ {
		err = actions.StartNewOCR2Round(int64(round), aggregatorContracts, env.EVMClient, time.Minute*5, l)
		require.NoError(t, err, "Error starting new OCR2 round")
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

	config, err := tc.GetConfig("Smoke", tc.OCR2)
	if err != nil {
		t.Fatal(err)
	}

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithGeth().
		WithMockAdapter().
		WithCLNodeConfig(node.NewConfig(node.NewBaseConfig(),
			node.WithOCR2(),
			node.WithP2Pv2(),
			node.WithTracing(),
		)).
		WithCLNodes(6).
		WithFunding(big.NewFloat(.1)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	env.ParallelTransactions(true)

	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	linkToken, err := env.ContractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	err = actions.FundChainlinkNodesLocal(workerNodes, env.EVMClient, big.NewFloat(.05))
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
	aggregatorContracts, err := actions.DeployOCRv2Contracts(1, linkToken, env.ContractDeployer, transmitters, env.EVMClient, ocrOffchainOptions)
	require.NoError(t, err, "Error deploying OCRv2 aggregator contracts")

	err = actions.CreateOCRv2JobsLocal(aggregatorContracts, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 5, env.EVMClient.GetChainID().Uint64(), false, false)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	ocrv2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
	require.NoError(t, err, "Error building OCRv2 config")

	err = actions.ConfigureOCRv2AggregatorContracts(env.EVMClient, ocrv2Config, aggregatorContracts)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

	err = actions.WatchNewOCR2Round(1, aggregatorContracts, env.EVMClient, time.Minute*5, l)
	require.NoError(t, err, "Error watching for new OCR2 round")
	roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(1))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 5 but got %d",
		roundData.Answer.Int64(),
	)

	err = env.MockAdapter.SetAdapterBasedIntValuePath("ocr2", []string{http.MethodGet, http.MethodPost}, 10)
	require.NoError(t, err)
	err = actions.WatchNewOCR2Round(2, aggregatorContracts, env.EVMClient, time.Minute*5, l)
	require.NoError(t, err)

	roundData, err = aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(2))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 10 but got %d",
		roundData.Answer.Int64(),
	)

	err = actions.DeleteJobs(nodeClients)
	require.NoError(t, err)

	err = actions.DeleteBridges(nodeClients)
	require.NoError(t, err)

	err = actions.CreateOCRv2JobsLocal(aggregatorContracts, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 15, env.EVMClient.GetChainID().Uint64(), false, false)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	err = actions.WatchNewOCR2Round(3, aggregatorContracts, env.EVMClient, time.Minute*3, l)
	require.NoError(t, err, "Error watching for new OCR2 round")
	roundData, err = aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(3))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(15), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 15 but got %d",
		roundData.Answer.Int64(),
	)
}
