package features

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products"
)

func TestReorgHeadTrackerFinalityViolation(t *testing.T) {
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	l := framework.L

	t.Cleanup(func() {
		reorgMessage := products.NewAllowedLogMessage(
			"Got very old block. Either a very deep re-org occurred, one of the RPC nodes has gotten far out of sync, or the chain went backwards in block numbers.",
			"this test causes reorg so this message is expected",
			zapcore.DPanicLevel,
			products.WarnAboutAllowedMsgs_No,
		)
		cleanupErr := products.CleanupContainerLogs(products.DefaultSettings(reorgMessage))
		require.NoError(t, cleanupErr, "failed to process cleanup container logs")
	})

	rpcClient := rpc.New(in.Blockchains[0].Out.Nodes[0].ExternalHTTPUrl, nil)
	clNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	// wait until we've reached depth
	waitForBlocks := 60
	// see ../products/ocr2/basic.toml, default finality depth is 5 for local env
	reorgForBlocks := 50
	timeout := 3 * time.Minute
	poll := 3 * time.Second

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		bn, err := rpcClient.BlockNumber()
		assert.NoError(c, err)
		l.Info().Int64("blockNumber", bn).Int("targetBlockNumber", waitForBlocks).Msg("Waiting for chain to progress above target block number")
		assert.GreaterOrEqual(c, bn, int64(waitForBlocks))
	}, timeout, poll, "timeout exceeded: target block was not reached")

	// reorg
	err = rpcClient.GethSetHead(reorgForBlocks)
	require.NoError(t, err)

	// verify all the nodes are reporting finality violation correctly
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		violated := 0
		for _, node := range clNodes {
			resp, _, err := node.Health()
			assert.NoError(c, err)
			for _, d := range resp.Data {
				if d.Attributes.Name == "EVM.1337.HeadTracker" &&
					strings.Contains(d.Attributes.Output, "finality violated") &&
					d.Attributes.Status == "failing" {
					violated++
				}
			}
			l.Debug().Msgf("Resp: %v", resp)
		}
		l.Info().Int("Violated", violated).Int("Nodes", len(clNodes)).Msg("Checking if all nodes reported finality violation")
		assert.Len(c, clNodes, violated)
	}, timeout, poll, "not all the nodes report finality violation")

	l.Info().Msg("All nodes reported finality violation")
}
