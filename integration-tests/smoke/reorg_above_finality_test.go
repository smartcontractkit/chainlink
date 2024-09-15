package smoke

import (
	"testing"
	"time"

	ctf_client "github.com/smartcontractkit/chainlink-testing-framework/lib/client"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestReorgAboveFinality_FinalityTagDisabled(t *testing.T) {
	t.Parallel()

	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig([]string{t.Name()}, tc.LogPoller)
	require.NoError(t, err, "Error getting config")

	privateNetworkConf, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err)

	// Get values from the node config
	configMap := make(map[string]interface{})
	err = toml.Unmarshal([]byte(config.NodeConfig.CommonChainConfigTOML), &configMap)
	require.NoError(t, err, "Error unmarshaling TOML")
	nodeFinalityDepthInt, isFinalityDepthSet := configMap["FinalityDepth"].(int64)
	nodeFinalityTagEnabled := configMap["FinalityTagEnabled"].(bool)
	l.Info().Int64("nodeFinalityDepth", nodeFinalityDepthInt).Bool("nodeFinalityTagEnabled", nodeFinalityTagEnabled).Msg("Node reorg config")

	var reorgDepth int
	if isFinalityDepthSet {
		reorgDepth = int(nodeFinalityDepthInt) + 5
	} else {
		reorgDepth = 15
	}
	minChainBlockNumberBeforeReorg := reorgDepth + 10

	testEnv, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetworkConf.EthereumNetworkConfig).
		WithCLNodes(6).
		WithoutCleanup().
		Build()
	require.NoError(t, err)

	evmNetwork, err := testEnv.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	client := ctf_client.NewRPCClient(evmNetwork.HTTPURLs[0], nil)

	// Wait for chain to progress
	require.Eventually(t, func() bool {
		bn, err := client.BlockNumber()
		require.NoError(t, err)
		l.Info().Int64("blockNumber", bn).Int("targetBlockNumber", minChainBlockNumberBeforeReorg).Msg("Waiting for chain to progress above target block number")
		return bn >= int64(minChainBlockNumberBeforeReorg)
	}, 8*time.Minute, 3*time.Second, "timeout exceeded: chain did not progress above the target block number")

	// Run reorg above finality depth
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
