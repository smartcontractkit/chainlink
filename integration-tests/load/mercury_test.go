package load

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	mercuryserver "github.com/smartcontractkit/chainlink-env/pkg/helm/mercury-server"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

const (
	mercuryFeedID = "mock-feed"
)

var (
	dbSettings = map[string]interface{}{
		"stateful": "true",
		"capacity": "10Gi",
		"resources": map[string]interface{}{
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2048Mi",
			},
		},
	}

	serverResources = map[string]interface{}{
		"limits": map[string]interface{}{
			"cpu":    "2000m",
			"memory": "2048Mi",
		},
	}
)

func setupMercuryLoadEnv(
	t *testing.T,
	feedID string,
	dbSettings map[string]interface{},
	serverResources map[string]interface{},
) (*environment.Environment, *client.MercuryServer, blockchain.EVMClient, uint64) {
	env, isExistingTestEnv, testNetwork, chainlinkNodes,
		mercuryServerRemoteUrl,
		evmClient, mockServerClient, mercuryServerClient, msRpcPubKey := testsetups.SetupMercuryEnv(t, dbSettings, serverResources)
	_ = isExistingTestEnv

	nodesWithoutBootstrap := chainlinkNodes[1:]
	ocrConfig := testsetups.BuildMercuryOCR2Config(t, nodesWithoutBootstrap)
	verifier, _, _, _ := testsetups.SetupMercuryContracts(t, evmClient,
		mercuryServerRemoteUrl, feedID, ocrConfig)

	testsetups.SetupMercuryNodeJobs(t, chainlinkNodes, mockServerClient, verifier.Address(),
		feedID, msRpcPubKey, testNetwork.ChainID, 0)

	err := verifier.SetConfig(ocrConfig)
	require.NoError(t, err)

	// Wait for the DON to start generating reports
	d := 160 * time.Second
	log.Info().Msgf("Sleeping for %s to wait for Mercury env to be ready..", d)
	time.Sleep(d)

	latestBlockNum, err := evmClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Err getting latest block number")
	report, _, err := mercuryServerClient.GetReports(feedID, latestBlockNum-5)
	require.NoError(t, err, "Error getting report from Mercury Server")
	require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")
	return env, mercuryServerClient, evmClient, latestBlockNum
}

func TestMercuryHTTPLoad(t *testing.T) {
	env, msClient, evmClient, latestBlockNumber := setupMercuryLoadEnv(t, mercuryFeedID, dbSettings, serverResources)

	gun := NewHTTPGun(env.URLs[mercuryserver.URLsKey][1], msClient, mercuryFeedID, latestBlockNumber)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			bn, _ := evmClient.LatestBlockNumber(context.Background())
			log.Warn().Uint64("Block number", bn).Send()
			gun.bn.Store(bn - 5)
		}
	}()
	gen, err := loadgen.NewLoadGenerator(&loadgen.LoadGeneratorConfig{
		T: t,
		LokiConfig: ctfClient.NewDefaultLokiConfig(
			os.Getenv("LOKI_URL"),
			os.Getenv("LOKI_TOKEN")),
		Labels: map[string]string{
			"test_group": "stress",
			"cluster":    "sdlc",
			"app":        "mercury-server",
			"namespace":  env.Cfg.Namespace,
			"test_id":    "http",
		},
		Duration: 1200 * time.Second,
		Schedule: &loadgen.LoadSchedule{
			Type:          loadgen.RPSScheduleType,
			StartFrom:     10,
			Increase:      5,
			StageInterval: 20 * time.Second,
			Limit:         1000,
		},
		Gun: gun,
	})
	require.NoError(t, err)
	gen.Run()
	_, _ = gen.Wait()
}

func TestMercuryWSLoad(t *testing.T) {
	env, msClient, _, _ := setupMercuryLoadEnv(t, mercuryFeedID, dbSettings, serverResources)

	gen, err := loadgen.NewLoadGenerator(&loadgen.LoadGeneratorConfig{
		T: t,
		LokiConfig: ctfClient.NewDefaultLokiConfig(
			os.Getenv("LOKI_URL"),
			os.Getenv("LOKI_TOKEN")),
		Labels: map[string]string{
			"test_group": "stress",
			"cluster":    "sdlc",
			"app":        "mercury-server",
			"namespace":  env.Cfg.Namespace,
			"test_id":    "ws",
		},
		Duration: 1200 * time.Second,
		Schedule: &loadgen.LoadSchedule{
			Type:          loadgen.InstancesScheduleType,
			StartFrom:     10,
			Increase:      20,
			StageInterval: 10 * time.Second,
			Limit:         500,
		},
		Instance: NewWSInstance(msClient),
	})
	require.NoError(t, err)
	gen.Run()
	_, _ = gen.Wait()
}
