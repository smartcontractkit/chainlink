package loadfunctions

import (
	"testing"

	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
)

func TestGatewayLoad(t *testing.T) {
	cfg, err := ReadConfig()
	require.NoError(t, err)
	ft, err := SetupLocalLoadTestEnv(cfg)
	require.NoError(t, err)
	ft.EVMClient.ParallelTransactions(false)

	labels := map[string]string{
		"branch": "gateway_healthcheck",
		"commit": "gateway_healthcheck",
	}

	secretsListCfg := &wasp.Config{
		LoadType: wasp.RPS,
		GenName:  functions.MethodSecretsList,
		Schedule: wasp.Plain(
			cfg.GatewayListSoak.RPS,
			cfg.GatewayListSoak.Duration.Duration(),
		),
		Gun: NewGatewaySecretsSetGun(
			cfg,
			functions.MethodSecretsList,
			ft.EthereumPrivateKey,
			ft.ThresholdPublicKey,
			ft.DONPublicKey,
		),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	}

	secretsSetCfg := &wasp.Config{
		LoadType: wasp.RPS,
		GenName:  functions.MethodSecretsSet,
		Schedule: wasp.Plain(
			cfg.GatewaySetSoak.RPS,
			cfg.GatewaySetSoak.Duration.Duration(),
		),
		Gun: NewGatewaySecretsSetGun(
			cfg,
			functions.MethodSecretsSet,
			ft.EthereumPrivateKey,
			ft.ThresholdPublicKey,
			ft.DONPublicKey,
		),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	}

	t.Run("gateway secrets list soak test", func(t *testing.T) {
		secretsListCfg.T = t
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(secretsListCfg)).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("gateway secrets set soak test", func(t *testing.T) {
		secretsListCfg.T = t
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(secretsSetCfg)).
			Run(true)
		require.NoError(t, err)
	})
}
