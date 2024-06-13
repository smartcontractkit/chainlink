package smoke

import (
	"math/big"
	"testing"
	"time"

	ctf_client "github.com/smartcontractkit/chainlink-testing-framework/client"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestReorgAboveFinality(t *testing.T) {
	t.Parallel()

	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Smoke", tc.OCR2)
	require.NoError(t, err, "Error getting config")

	privateNetworkConf, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err)

	nodeFinalityDepthInt := int64(10)

	testEnv, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetworkConf.EthereumNetworkConfig).
		WithMockAdapter().
		WithCLNodes(6).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithoutCleanup().
		WithSeth().
		Build()
	require.NoError(t, err)

	network := testEnv.EVMNetworks[0]
	client := ctf_client.NewRPCClient(network.HTTPURLs[0])

	// Wait for chain to progress
	targetBlockNumber := nodeFinalityDepthInt * 3
	require.Eventually(t, func() bool {
		bn, err := client.BlockNumber()
		require.NoError(t, err)
		l.Info().Int64("blockNumber", bn).Int64("targetBlockNumber", targetBlockNumber).Msg("Waiting for chain to progress above target block number")
		return bn > nodeFinalityDepthInt*3
	}, 3*time.Minute, 3*time.Second, "chain did not progress above the target block number")

	// Run reorg above finality depth
	reorgDepth := int(nodeFinalityDepthInt) + 20
	l.Info().
		Str("URL", client.URL).
		Int64("nodeFinalityDepth", nodeFinalityDepthInt).
		Int("reorgDepth", reorgDepth).
		Msg("Starting blockchain reorg on Simulated Geth chain")
	err = client.GethSetHead(reorgDepth)
	require.NoError(t, err, "Error starting blockchain reorg on Simulated Geth chain")

	l.Info().Msg("Waiting for all nodes to report finality violation")
	nodes := testEnv.ClCluster.NodeAPIs()
	require.Eventually(t, func() bool {
		violatedResponses := 0
		for _, node := range nodes {
			resp, _, err := node.Health()
			require.NoError(t, err)
			for _, d := range resp.Data {
				if d.Attributes.Name == "EVM.1337.LogPoller" && d.Attributes.Output == "finality violated" && d.Attributes.Status == "failing" {
					violatedResponses++
				}
			}
			l.Info().Msgf("Resp: %v", resp)
		}

		l.Info().Int("violatedResponses", violatedResponses).Int("nodes", len(nodes)).Msg("Checking if all nodes reported finality violation")
		return violatedResponses == len(nodes)
	}, 3*time.Minute, 5*time.Second, "not all the nodes report finality violation")
	l.Info().Msg("All nodes reported finality violation")
}
