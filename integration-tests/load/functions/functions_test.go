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
	ft, err := SetupLocalLoadTestEnv(cfg)
	require.NoError(t, err)
	ft.EVMClient.ParallelTransactions(false)

	labels := map[string]string{
		"branch": "functions_healthcheck",
		"commit": "functions_healthcheck",
	}

	MonitorLoadStats(t, ft, labels)

	t.Run("functions soak test", func(t *testing.T) {
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_soak_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					cfg.Soak.RPS,
					cfg.Soak.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					cfg.Soak.RequestsPerCall,
					cfg.Common.FunctionsCallPayload,
					cfg.Common.SecretsSlotID,
					cfg.Common.SecretsVersionID,
					[]string{},
					cfg.Common.SubscriptionID,
					StringToByte32(cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewEnvLokiConfig(),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("functions stress test", func(t *testing.T) {
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_stress_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					cfg.Stress.RPS,
					cfg.Stress.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					cfg.Soak.RequestsPerCall,
					cfg.Common.FunctionsCallPayload,
					cfg.Common.SecretsSlotID,
					cfg.Common.SecretsVersionID,
					[]string{},
					cfg.Common.SubscriptionID,
					StringToByte32(cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewEnvLokiConfig(),
			})).
			Run(true)
		require.NoError(t, err)
	})
}
