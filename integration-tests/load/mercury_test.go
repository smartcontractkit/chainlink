package load

import (
	"context"
	"os"
	"testing"
	"time"

	mercuryserver "github.com/smartcontractkit/chainlink-env/pkg/helm/mercury-server"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/load/tools"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

var (
	resources = &mercury.ResourcesConfig{
		DONResources: map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2048Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2048Mi",
			},
		},
		DONDBResources: map[string]interface{}{
			"stateful": "true",
			"capacity": "10Gi",
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "2000m",
					"memory": "2048Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "2000m",
					"memory": "2048Mi",
				},
			},
		},
		MercuryResources: map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2048Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2048Mi",
			},
		},
		MercuryDBResources: map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2048Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2048Mi",
			},
		},
	}
)

func TestMercuryHTTPLoad(t *testing.T) {
	feeds := [][32]byte{mercury.StringToByte32("feed-1")}
	testEnv, err := mercury.SetupMercuryMultiFeedEnv(t.Name(), "load", feeds, resources)
	require.NoError(t, err)
	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	bn, _ := testEnv.EvmClient.LatestBlockNumber(context.Background())

	gun := tools.NewHTTPGun(testEnv.Env.URLs[mercuryserver.URLsKey][1], testEnv.MSClient, feeds, bn)
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
		Schedule:    loadgen.Line(10, 300, 100*time.Second),
		Gun:         gun,
	})
	require.NoError(t, err)
	gen.Run()
	_, _ = gen.Wait()
}

func TestMercuryWSLoad(t *testing.T) {
	feeds := [][32]byte{mercury.StringToByte32("feed-1")}
	testEnv, err := mercury.SetupMercuryMultiFeedEnv(t.Name(), "load", feeds, resources)
	require.NoError(t, err)
	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})

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
