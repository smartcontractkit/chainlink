package loadfunctions

import (
	"testing"

	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
)

func TestGatewayLoad(t *testing.T) {
	listConfig, err := tc.GetConfig("GatewayList", tc.Functions)
	require.NoError(t, err)
	cfgl := listConfig.Logging.Loki

	require.NoError(t, err)
	ft, err := SetupLocalLoadTestEnv(&listConfig, &listConfig)
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
			*listConfig.Functions.Performance.RPS,
			listConfig.Functions.Performance.Duration.Duration,
		),
		Gun: NewGatewaySecretsSetGun(
			&listConfig,
			functions.MethodSecretsList,
			ft.EthereumPrivateKey,
			ft.ThresholdPublicKey,
			ft.DONPublicKey,
		),
		Labels:     labels,
		LokiConfig: wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
	}

	setConfig, err := tc.GetConfig("GatewaySet", tc.Functions)
	require.NoError(t, err)

	secretsSetCfg := &wasp.Config{
		LoadType: wasp.RPS,
		GenName:  functions.MethodSecretsSet,
		Schedule: wasp.Plain(
			*setConfig.Functions.Performance.RPS,
			setConfig.Functions.Performance.Duration.Duration,
		),
		Gun: NewGatewaySecretsSetGun(
			&setConfig,
			functions.MethodSecretsSet,
			ft.EthereumPrivateKey,
			ft.ThresholdPublicKey,
			ft.DONPublicKey,
		),
		Labels:     labels,
		LokiConfig: wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
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
