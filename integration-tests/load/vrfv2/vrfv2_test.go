package loadvrfv2

import (
	"testing"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
)

func TestVRFV2Load(t *testing.T) {
	cfg, err := ReadConfig()
	require.NoError(t, err)
	env, vrfv2Contracts, key, err := vrfv2_actions.SetupLocalLoadTestEnv(cfg.Common.NodeFunds, cfg.Common.SubFunds)
	require.NoError(t, err)

	labels := map[string]string{
		"branch": "vrfv2_healthcheck",
		"commit": "vrfv2_healthcheck",
	}

	singleFeedConfig := &wasp.Config{
		T:          t,
		LoadType:   wasp.RPS,
		GenName:    "gun",
		Gun:        SingleFeedGun(vrfv2Contracts, key),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	}

	multiFeedConfig := &wasp.Config{
		T:          t,
		LoadType:   wasp.VU,
		GenName:    "vu",
		VU:         NewJobVolumeVU(cfg.SoakVolume.Pace.Duration(), 1, env.ClCluster.NodeAPIs(), env.EVMClient, vrfv2Contracts),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	}

	MonitorLoadStats(t, vrfv2Contracts, labels)

	// is our "job" stable at all, no memory leaks, no flaking performance under some RPS?
	t.Run("vrfv2 soak test", func(t *testing.T) {
		singleFeedConfig.Schedule = wasp.Plain(
			cfg.Soak.RPS,
			cfg.Soak.Duration.Duration(),
		)
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)
	})

	// what are the limits for one "job", figuring out the max/optimal performance params by increasing RPS and varying configuration
	t.Run("vrfv2 load test", func(t *testing.T) {
		singleFeedConfig.Schedule = wasp.Steps(
			cfg.Load.RPSFrom,
			cfg.Load.RPSIncrease,
			cfg.Load.RPSSteps,
			cfg.Load.Duration.Duration(),
		)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)
	})

	// how many "jobs" of the same type we can run at once at a stable load with optimal configuration?
	t.Run("vrfv2 volume soak test", func(t *testing.T) {
		multiFeedConfig.Schedule = wasp.Plain(
			cfg.SoakVolume.Products,
			cfg.SoakVolume.Duration.Duration(),
		)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(multiFeedConfig)).
			Run(true)
		require.NoError(t, err)
	})

	// what are the limits if we add more and more "jobs/products" of the same type, each "job" have a stable RPS we vary only amount of jobs
	t.Run("vrfv2 volume load test", func(t *testing.T) {
		multiFeedConfig.Schedule = wasp.Steps(
			cfg.LoadVolume.ProductsFrom,
			cfg.LoadVolume.ProductsIncrease,
			cfg.LoadVolume.ProductsSteps,
			cfg.LoadVolume.Duration.Duration(),
		)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(multiFeedConfig)).
			Run(true)
		require.NoError(t, err)
	})
}
