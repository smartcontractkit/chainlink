package load

import (
	"context"
	"os"
	"testing"
	"time"

	mercuryserver "github.com/smartcontractkit/chainlink-env/pkg/helm/mercury-server"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/smartcontractkit/chainlink/integration-tests/load/tools"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
	"github.com/stretchr/testify/require"
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
	dbSettings map[string]interface{},
	serverResources map[string]interface{},
) (*mercury.TestEnv, uint64) {
	testEnv, err := mercury.SetupMercuryTestEnv("load", dbSettings, serverResources)
	require.NoError(t, err)

	t.Cleanup(func() {
		//nolint
		testEnv.Cleanup(t)
	})

	latestBlockNum, err := testEnv.EvmClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Err getting latest block number")
	report, _, err := testEnv.MSClient.GetReports(testEnv.FeedId, latestBlockNum-5)
	require.NoError(t, err, "Error getting report from Mercury Server")
	require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")

	return testEnv, latestBlockNum
}

func TestMercuryHTTPLoad(t *testing.T) {
	testEnv, latestBlockNumber := setupMercuryLoadEnv(t, dbSettings, serverResources)

	gun := tools.NewHTTPGun(testEnv.Env.URLs[mercuryserver.URLsKey][1], testEnv.MSClient, testEnv.FeedId, latestBlockNumber)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			bn, _ := testEnv.EvmClient.LatestBlockNumber(context.Background())
			gun.Bn.Store(bn - 5)
		}
	}()
	gen, err := loadgen.NewLoadGenerator(&loadgen.Config{
		T: t,
		LokiConfig: ctfClient.NewDefaultLokiConfig(
			os.Getenv("LOKI_URL"),
			os.Getenv("LOKI_TOKEN")),
		Labels: map[string]string{
			"test_group": "stress",
			"cluster":    "sdlc",
			"app":        "mercury-server",
			"namespace":  testEnv.Env.Cfg.Namespace,
			"test_id":    "http",
		},
		LoadType:    loadgen.RPSScheduleType,
		CallTimeout: 5 * time.Second,
		Schedule:    loadgen.Line(10, 800, 500*time.Second),
		Gun:         gun,
	})
	require.NoError(t, err)
	gen.Run()
	_, _ = gen.Wait()
}

func TestMercuryWSLoad(t *testing.T) {
	testEnv, _ := setupMercuryLoadEnv(t, dbSettings, serverResources)

	gen, err := loadgen.NewLoadGenerator(&loadgen.Config{
		T: t,
		LokiConfig: ctfClient.NewDefaultLokiConfig(
			os.Getenv("LOKI_URL"),
			os.Getenv("LOKI_TOKEN")),
		Labels: map[string]string{
			"test_group": "stress",
			"cluster":    "sdlc",
			"app":        "mercury-server",
			"namespace":  testEnv.Env.Cfg.Namespace,
			"test_id":    "ws",
		},
		LoadType: loadgen.InstancesScheduleType,
		Schedule: loadgen.Line(10, 300, 30*time.Second),
		Instance: tools.NewWSInstance(testEnv.MSClient),
	})
	require.NoError(t, err)
	gen.Run()
	_, _ = gen.Wait()
}
