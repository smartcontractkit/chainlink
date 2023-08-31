package loadfunctions

import (
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGatewayLoad(t *testing.T) {
	cfg, err := ReadConfig()
	require.NoError(t, err)
	env, ft, err := SetupLocalLoadTestEnv(cfg)
	require.NoError(t, err)
	env.ParallelTransactions(false)

	labels := map[string]string{
		"branch": "gateway_healthcheck",
		"commit": "gateway_healthcheck",
	}

	gatewayGunConfig := &wasp.Config{
		T:           t,
		LoadType:    wasp.RPS,
		GenName:     "gun",
		CallTimeout: 2 * time.Minute,
		Gun: NewGatewaySecretsSetGun(
			cfg,
			ft.EthereumPrivateKey,
			ft.ThresholdPublicKey,
			ft.DONPublicKey,
		),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	}

	t.Run("gateway secrets set soak test", func(t *testing.T) {
		gatewayGunConfig.Schedule = wasp.Plain(
			cfg.Soak.RPS,
			cfg.Soak.Duration.Duration(),
		)
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(gatewayGunConfig)).
			Run(true)
		require.NoError(t, err)
	})

	//t.Run("functions load test", func(t *testing.T) {
	//	singleFeedConfig.Schedule = wasp.Steps(
	//		cfg.Load.RPSFrom,
	//		cfg.Load.RPSIncrease,
	//		cfg.Load.RPSSteps,
	//		cfg.Load.Duration.Duration(),
	//	)
	//	_, err = wasp.NewProfile().
	//		Add(wasp.NewGenerator(singleFeedConfig)).
	//		Run(true)
	//	require.NoError(t, err)
	//})
}
