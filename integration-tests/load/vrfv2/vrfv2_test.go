package loadvrfv2

import (
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func TestVRFV2Load(t *testing.T) {
	env, vrfv2Contracts, key, err := vrfv2_actions.SetupLocalLoadTestEnv(big.NewFloat(10), big.NewInt(100))
	require.NoError(t, err)

	labels := map[string]string{
		"branch": "vrfv2_healthcheck",
		"commit": "vrfv2_healthcheck",
	}

	singleFeedConfig := &wasp.Config{
		T:          t,
		LoadType:   wasp.RPS,
		GenName:    "single_feed_gun",
		Gun:        SingleFeedGun(vrfv2Contracts, key),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	}

	multiFeedConfig := &wasp.Config{
		T:          t,
		LoadType:   wasp.VU,
		GenName:    "job_volume_gun",
		VU:         NewJobVolumeVU(1*time.Second, 1, env.GetAPIs(), env.Geth.EthClient, vrfv2Contracts),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	}

	t.Run("vrfv2 soak test", func(t *testing.T) {
		singleFeedConfig.Schedule = wasp.Plain(1, 10*time.Minute)

		MonitorLoadStats(t, vrfv2Contracts, labels)
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("vrfv2 constant volume soak test", func(t *testing.T) {
		multiFeedConfig.Schedule = wasp.Plain(5, 30*time.Minute)

		MonitorLoadStats(t, vrfv2Contracts, labels)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(multiFeedConfig)).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("vrfv2 one feed load test", func(t *testing.T) {
		singleFeedConfig.Schedule = wasp.Line(10, 50, 5*time.Minute)

		MonitorLoadStats(t, vrfv2Contracts, labels)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("vrfv2 volume load test", func(t *testing.T) {
		multiFeedConfig.Schedule = wasp.Line(5, 55, 3*time.Minute)

		MonitorLoadStats(t, vrfv2Contracts, labels)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(multiFeedConfig)).
			Run(true)
		require.NoError(t, err)
	})
}
