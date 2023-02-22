package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/stretchr/testify/require"
)

func TestMercurySmoke(t *testing.T) {
	const mercuryFeedId = "ETH-USD-Optimism-Goerli-1"
	_, isExistingTestEnv, testNetwork, chainlinkNodes, mercuryServerInternalUrl,
		evmClient, mockServerClient, mercuryServerClient := testsetups.SetupMercuryEnv(t)

	if isExistingTestEnv {
		log.Info().Msg("Use existing Mercury test env")
	} else {
		log.Info().Msg("Creating new Mercury test env..")
		nodesWithoutBootstrap := chainlinkNodes[1:]
		ocrConfig := testsetups.BuildMercuryOCR2Config(t, nodesWithoutBootstrap)
		verifier, _ := testsetups.SetupMercuryContracts(t, evmClient, mercuryFeedId, ocrConfig)
		testsetups.SetupMercuryNodeJobs(t, chainlinkNodes, mockServerClient, verifier.Address(),
			mercuryFeedId, mercuryServerInternalUrl, testNetwork.ChainID, 0)
		// Set OCR2 config in the contract
		verifier.SetConfig(ocrConfig)
		// Wait for the DON to start generating reports
		d := 160 * time.Second
		log.Info().Msgf("Sleeping for %s to wait for Mercury env to be ready..", d)
		time.Sleep(d)
	}

	t.Run("test mercury server has report for the latest block number", func(t *testing.T) {
		latestBlockNum, err := evmClient.LatestBlockNumber(context.Background())
		require.NoError(t, err, "Err getting latest block number")
		report, _, err := mercuryServerClient.GetReports(mercuryFeedId, latestBlockNum)
		require.NoError(t, err, "Error getting report from Mercury Server")
		require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")
	})
}
