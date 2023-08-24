package loadfunctions

import (
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFunctionsLoad(t *testing.T) {
	cfg, err := ReadConfig()
	require.NoError(t, err)
	env, functionContracts, err := SetupLocalLoadTestEnv(cfg)
	require.NoError(t, err)
	env.ParallelTransactions(true)

	labels := map[string]string{
		"branch": "functions_healthcheck",
		"commit": "functions_healthcheck",
	}

	singleFeedConfig := &wasp.Config{
		T:           t,
		LoadType:    wasp.RPS,
		GenName:     "gun",
		CallTimeout: 2 * time.Minute,
		Gun: NewSingleFunctionCallGun(
			functionContracts,
			"const response = await Functions.makeHttpRequest({ url: 'http://dummyjson.com/products/1' }); return Functions.encodeUint256(response.data.id)",
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
