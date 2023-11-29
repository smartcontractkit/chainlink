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
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
)

// Tests a basic OCRv2 median feed
func TestOCRv2Basic(t *testing.T) {
	t.Parallel()
	chainConfig := ctf_test_env.EthereumChainConfig{
		SecondsPerSlot: 8,
		SlotsPerEpoch:  4,
	}

	networks := map[string]ctf_test_env.EthereumNetwork{
		"geth": func() ctf_test_env.EthereumNetwork {
			ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
			cfg, err := ethBuilder.
				WithConsensusType(ctf_test_env.ConsensusType_PoS).
				WithConsensusLayer(ctf_test_env.ConsensusLayer_Prysm).
				WithExecutionLayer(ctf_test_env.ExecutionLayer_Geth).
				WithEthereumChainConfig(chainConfig).
				Build()
			require.NoError(t, err, "Error building ethereum network config")
			return cfg
		}(),
		"besu": func() ctf_test_env.EthereumNetwork {
			ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
			cfg, err := ethBuilder.
				WithConsensusType(ctf_test_env.ConsensusType_PoS).
				WithConsensusLayer(ctf_test_env.ConsensusLayer_Prysm).
				WithExecutionLayer(ctf_test_env.ExecutionLayer_Besu).
				WithEthereumChainConfig(chainConfig).
				Build()
			require.NoError(t, err, "Error building ethereum network config")
			return cfg
		}(),
		"erigon": func() ctf_test_env.EthereumNetwork {
			ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
			cfg, err := ethBuilder.
				WithConsensusType(ctf_test_env.ConsensusType_PoS).
				WithConsensusLayer(ctf_test_env.ConsensusLayer_Prysm).
				WithExecutionLayer(ctf_test_env.ExecutionLayer_Erigon).
				WithEthereumChainConfig(chainConfig).
				Build()
			require.NoError(t, err, "Error building ethereum network config")
			return cfg
		}(),
		"nethermind": func() ctf_test_env.EthereumNetwork {
			ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
			cfg, err := ethBuilder.
				WithConsensusType(ctf_test_env.ConsensusType_PoS).
				WithConsensusLayer(ctf_test_env.ConsensusLayer_Prysm).
				WithExecutionLayer(ctf_test_env.ExecutionLayer_Nethermind).
				WithEthereumChainConfig(chainConfig).
				Build()
			require.NoError(t, err, "Error building ethereum network config")
			return cfg
		}(),
	}

	for k, v := range networks {
		name := k
		network := v
		t.Run(name, func(t *testing.T) {
			l := logging.GetTestLogger(t)

			env, err := test_env.NewCLTestEnvBuilder().
				WithTestLogger(t).
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

			err = actions.CreateOCRv2JobsLocal(aggregatorContracts, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 5, env.EVMClient.GetChainID().Uint64(), false)
			require.NoError(t, err, "Error creating OCRv2 jobs")

			ocrv2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
			require.NoError(t, err, "Error building OCRv2 config")

			err = actions.ConfigureOCRv2AggregatorContracts(env.EVMClient, ocrv2Config, aggregatorContracts)
			require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

			err = actions.StartNewOCR2Round(1, aggregatorContracts, env.EVMClient, time.Minute*5, l)
			require.NoError(t, err, "Error starting new OCR2 round")
			roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(1))
			require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
			require.Equal(t, int64(5), roundData.Answer.Int64(),
				"Expected latest answer from OCR contract to be 5 but got %d",
				roundData.Answer.Int64(),
			)

			err = env.MockAdapter.SetAdapterBasedIntValuePath("ocr2", []string{http.MethodGet, http.MethodPost}, 10)
			require.NoError(t, err)
			err = actions.StartNewOCR2Round(2, aggregatorContracts, env.EVMClient, time.Minute*5, l)
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

func TestOCRv2JobReplacement(t *testing.T) {
	l := logging.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
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

	err = actions.CreateOCRv2JobsLocal(aggregatorContracts, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 5, env.EVMClient.GetChainID().Uint64(), false)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	ocrv2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
	require.NoError(t, err, "Error building OCRv2 config")

	err = actions.ConfigureOCRv2AggregatorContracts(env.EVMClient, ocrv2Config, aggregatorContracts)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

	err = actions.StartNewOCR2Round(1, aggregatorContracts, env.EVMClient, time.Minute*5, l)
	require.NoError(t, err, "Error starting new OCR2 round")
	roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(1))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 5 but got %d",
		roundData.Answer.Int64(),
	)

	err = env.MockAdapter.SetAdapterBasedIntValuePath("ocr2", []string{http.MethodGet, http.MethodPost}, 10)
	require.NoError(t, err)
	err = actions.StartNewOCR2Round(2, aggregatorContracts, env.EVMClient, time.Minute*5, l)
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

	err = actions.CreateOCRv2JobsLocal(aggregatorContracts, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 15, env.EVMClient.GetChainID().Uint64(), false)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	err = actions.StartNewOCR2Round(3, aggregatorContracts, env.EVMClient, time.Minute*3, l)
	require.NoError(t, err, "Error starting new OCR2 round")
	roundData, err = aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(3))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(15), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 15 but got %d",
		roundData.Answer.Int64(),
	)
}
