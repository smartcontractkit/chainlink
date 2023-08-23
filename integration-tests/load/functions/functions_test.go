package loadfunctions

import (
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFunctionsLoad(t *testing.T) {
	cfg, err := ReadConfig()
	require.NoError(t, err)
	_, functionContracts, err := SetupLocalLoadTestEnv(cfg)
	require.NoError(t, err)

	labels := map[string]string{
		"branch": "functions_healthcheck",
		"commit": "functions_healthcheck",
	}

	singleFeedConfig := &wasp.Config{
		T:        t,
		LoadType: wasp.RPS,
		GenName:  "gun",
		Gun: NewSingleFunctionCallGun(
			functionContracts,
			"return Functions.encodeUint256(1)",
			[]byte{},
			[]string{},
			cfg.Common.SubscriptionID,
			StringToByte32(cfg.Common.DONID),
		),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	}

	t.Run("functions soak test", func(t *testing.T) {
		singleFeedConfig.Schedule = wasp.Plain(
			cfg.Soak.RPS,
			cfg.Soak.Duration.Duration(),
		)
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("functions load test", func(t *testing.T) {
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
}
