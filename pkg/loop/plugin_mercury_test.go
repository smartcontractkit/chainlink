package loop_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	mercury_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	testcore "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/core"
	relayer_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/relayer"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestPluginMercury(t *testing.T) {
	t.Parallel()

	stopCh := newStopCh(t)
	test.PluginTest(t, loop.PluginMercuryName, &loop.GRPCPluginMercury{PluginServer: mercury_test.FactoryServer, BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}, mercury_test.PluginMercury)

	t.Run("proxy", func(t *testing.T) {
		test.PluginTest(t, loop.PluginRelayerName,
			&loop.GRPCPluginRelayer{
				PluginServer: relayer_test.NewRelayerTester(false),
				BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}},
			func(t *testing.T, pr loop.PluginRelayer) {
				p := newMercuryProvider(t, pr)
				pm := mercury_test.PluginMercuryTest{MercuryProvider: p}
				test.PluginTest(t, loop.PluginMercuryName,
					&loop.GRPCPluginMercury{PluginServer: mercury_test.FactoryServer,
						BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}},
					pm.TestPluginMercury)
			})
	})
}

func TestPluginMercuryExec(t *testing.T) {
	t.Parallel()
	stopCh := newStopCh(t)
	mercury := loop.GRPCPluginMercury{BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}
	cc := mercury.ClientConfig()
	cc.Cmd = NewHelperProcessCommand(loop.PluginMercuryName, true, 0)
	c := plugin.NewClient(cc)
	t.Cleanup(c.Kill)
	client, err := c.Client()
	require.NoError(t, err)
	defer client.Close()
	require.NoError(t, client.Ping())

	i, err := client.Dispense(loop.PluginMercuryName)
	require.NoError(t, err)
	require.NotNil(t, i)
	mercury_test.PluginMercury(t, i.(types.PluginMercury))

	t.Run("proxy", func(t *testing.T) {
		pr := newPluginRelayerExec(t, true, stopCh)
		p := newMercuryProvider(t, pr)
		pm := mercury_test.PluginMercuryTest{MercuryProvider: p}
		pm.TestPluginMercury(t, i.(types.PluginMercury))
	})
}

func newMercuryProvider(t *testing.T, pr loop.PluginRelayer) types.MercuryProvider {
	ctx := context.Background()
	r, err := pr.NewRelayer(ctx, test.ConfigTOML, testcore.Keystore)
	require.NoError(t, err)
	servicetest.Run(t, r)
	p, err := r.NewPluginProvider(ctx, mercury_test.RelayArgs, mercury_test.PluginArgs)
	mp, ok := p.(types.MercuryProvider)
	require.True(t, ok)
	require.NoError(t, err)
	servicetest.Run(t, mp)
	return mp
}
