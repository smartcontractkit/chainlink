package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/stretchr/testify/require"
)

func TestMercuryHTTPLoad(t *testing.T) {
	const mercuryFeedId = "ETH-USD-Optimism-Goerli-1"

	_, isExistingTestEnv, testNetwork, chainlinkNodes,
		mercuryServerRemoteUrl,
		evmClient, mockServerClient, mercuryServerClient, msRpcPubKey := testsetups.SetupMercuryEnv(t)
	_ = isExistingTestEnv

	nodesWithoutBootstrap := chainlinkNodes[1:]
	ocrConfig := testsetups.BuildMercuryOCR2Config(t, nodesWithoutBootstrap)
	verifier, _, _, _ := testsetups.SetupMercuryContracts(t, evmClient,
		mercuryServerRemoteUrl, mercuryFeedId, ocrConfig)

	testsetups.SetupMercuryNodeJobs(t, chainlinkNodes, mockServerClient, verifier.Address(),
		mercuryFeedId, msRpcPubKey, testNetwork.ChainID, 0)

	err := verifier.SetConfig(ocrConfig)
	require.NoError(t, err)

	// Wait for the DON to start generating reports
	d := 160 * time.Second
	log.Info().Msgf("Sleeping for %s to wait for Mercury env to be ready..", d)
	time.Sleep(d)

	latestBlockNum, err := evmClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Err getting latest block number")
	report, _, err := mercuryServerClient.GetReports(mercuryFeedId, latestBlockNum-5)
	require.NoError(t, err, "Error getting report from Mercury Server")
	require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")
}
