package loadfunctions

import (
	"testing"
	"time"

	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestFunctionsLoad(t *testing.T) {
	generalConfig, err := tc.GetConfig(tc.NoKey, tc.Functions)
	require.NoError(t, err, "failed to get config")

	ft, err := SetupLocalLoadTestEnv(&generalConfig, &generalConfig)
	require.NoError(t, err)

	labels := map[string]string{
		"branch": "functions_healthcheck",
		"commit": "functions_healthcheck",
	}

	MonitorLoadStats(t, ft, labels, &generalConfig)

	t.Run("mumbai functions soak test http", func(t *testing.T) {
		config, err := tc.GetConfig("Soak", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		cfgl := config.Logging.Loki
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_soak_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration,
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeHTTPPayload,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadHTTP,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions stress test http", func(t *testing.T) {
		config, err := tc.GetConfig("Stress", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		cfgl := config.Logging.Loki
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_stress_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration,
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeHTTPPayload,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadHTTP,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions soak test only secrets", func(t *testing.T) {
		config, err := tc.GetConfig("SecretsSoak", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		cfgl := config.Logging.Loki
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_soak_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration,
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeSecretsOnlyPayload,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadWithSecrets,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions stress test only secrets", func(t *testing.T) {
		config, err := tc.GetConfig("SecretsStress", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		cfgl := config.Logging.Loki
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_stress_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration,
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeSecretsOnlyPayload,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadWithSecrets,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions soak test real", func(t *testing.T) {
		config, err := tc.GetConfig("RealSoak", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		cfgl := config.Logging.Loki
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_soak_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration,
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeReal,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadReal,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{"1", "2", "3", "4"},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions stress test real", func(t *testing.T) {
		config, err := tc.GetConfig("RealStress", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		cfgl := config.Logging.Loki
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_stress_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration,
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeReal,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadReal,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{"1", "2", "3", "4"},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
			})).
			Run(true)
		require.NoError(t, err)
	})
}
