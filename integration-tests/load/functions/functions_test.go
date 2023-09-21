package loadfunctions

import (
	"github.com/rs/zerolog/log"
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

	t.Run("functions soak test http", func(t *testing.T) {
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
					ModeHTTPPayload,
					cfg.Soak.RequestsPerCall,
					cfg.Common.FunctionsCallPayloadHTTP,
					cfg.Common.SecretsSlotID,
					cfg.Common.SecretsVersionID,
					[]string{},
					cfg.SelectedNetwork.SubscriptionID,
					StringToByte32(cfg.SelectedNetwork.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewEnvLokiConfig(),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("functions stress test http", func(t *testing.T) {
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
					ModeHTTPPayload,
					cfg.Stress.RequestsPerCall,
					cfg.Common.FunctionsCallPayloadHTTP,
					cfg.Common.SecretsSlotID,
					cfg.Common.SecretsVersionID,
					[]string{},
					cfg.SelectedNetwork.SubscriptionID,
					StringToByte32(cfg.SelectedNetwork.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewEnvLokiConfig(),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("functions soak test only secrets", func(t *testing.T) {
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_soak_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					cfg.SecretsSoak.RPS,
					cfg.SecretsSoak.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeSecretsOnlyPayload,
					cfg.SecretsSoak.RequestsPerCall,
					cfg.Common.FunctionsCallPayloadWithSecrets,
					cfg.Common.SecretsSlotID,
					cfg.Common.SecretsVersionID,
					[]string{},
					cfg.SelectedNetwork.SubscriptionID,
					StringToByte32(cfg.SelectedNetwork.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewEnvLokiConfig(),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("functions stress test only secrets", func(t *testing.T) {
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_stress_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					cfg.SecretsStress.RPS,
					cfg.SecretsStress.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeSecretsOnlyPayload,
					cfg.SecretsStress.RequestsPerCall,
					cfg.Common.FunctionsCallPayloadWithSecrets,
					cfg.Common.SecretsSlotID,
					cfg.Common.SecretsVersionID,
					[]string{},
					cfg.SelectedNetwork.SubscriptionID,
					StringToByte32(cfg.SelectedNetwork.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewEnvLokiConfig(),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("functions soak test real", func(t *testing.T) {
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_soak_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					cfg.RealSoak.RPS,
					cfg.RealSoak.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeReal,
					cfg.RealSoak.RequestsPerCall,
					cfg.Common.FunctionsCallPayloadReal,
					cfg.Common.SecretsSlotID,
					cfg.Common.SecretsVersionID,
					[]string{"1", "2", "3", "4"},
					cfg.SelectedNetwork.SubscriptionID,
					StringToByte32(cfg.SelectedNetwork.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewEnvLokiConfig(),
			})).
			Run(true)
		log.Warn().Msg("Awaiting all requests to be processed")
		time.Sleep(10 * time.Minute)
		require.NoError(t, err)
	})

	t.Run("functions stress test real", func(t *testing.T) {
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_stress_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					cfg.RealStress.RPS,
					cfg.RealStress.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeReal,
					cfg.RealStress.RequestsPerCall,
					cfg.Common.FunctionsCallPayloadReal,
					cfg.Common.SecretsSlotID,
					cfg.Common.SecretsVersionID,
					[]string{"1", "2", "3", "4"},
					cfg.SelectedNetwork.SubscriptionID,
					StringToByte32(cfg.SelectedNetwork.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewEnvLokiConfig(),
			})).
			Run(true)
		require.NoError(t, err)
		log.Warn().Msg("Awaiting all requests to be processed")
		time.Sleep(5 * time.Minute)
	})
}
