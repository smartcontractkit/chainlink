package features

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	de "github.com/smartcontractkit/chainlink/devenv"
)

func TestReorgHeadTrackerFinalityViolation(t *testing.T) {
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	l := framework.L

	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	rpcClient := rpc.New(in.Blockchains[0].Out.Nodes[0].ExternalHTTPUrl, nil)
	clNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	// wait until we've reached depth
	waitForBlocks := 50
	// see ../products/ocr2/basic.toml, default finality depth is 5 for local env
	// twice max reorg depth
	reorgForBlocks := 10

	require.Eventually(t, func() bool {
		bn, err := rpcClient.BlockNumber()
		require.NoError(t, err)
		l.Info().Int64("blockNumber", bn).Int("targetBlockNumber", waitForBlocks).Msg("Waiting for chain to progress above target block number")
		return bn >= int64(waitForBlocks)
	}, 5*time.Minute, 3*time.Second, "timeout exceeded: chain did not progress above the target block number")

	// reorg
	err = rpcClient.GethSetHead(reorgForBlocks)
	require.NoError(t, err)

	// verify all the nodes are reporting finality violation correctly
	require.Eventually(t, func() bool {
		violated := 0
		for _, node := range clNodes {
			resp, _, err := node.Health()
			require.NoError(t, err)
			for _, d := range resp.Data {
				if d.Attributes.Name == "EVM.1337.HeadTracker" &&
					strings.Contains(d.Attributes.Output, "finality violated") &&
					d.Attributes.Status == "failing" {
					violated++
				}
			}
			l.Debug().Msgf("Resp: %v", resp)
		}

		l.Info().Int("violatedResponses", violated).Int("nodes", len(clNodes)).Msg("Checking if all nodes reported finality violation")
		return violated == len(clNodes)
	}, 3*time.Minute, 5*time.Second, "not all the nodes report finality violation")
	l.Info().Msg("All nodes reported finality violation")
}
