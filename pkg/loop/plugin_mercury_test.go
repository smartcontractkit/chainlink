package loop_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	mercury_common_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/mercury/common/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestPluginMercury(t *testing.T) {
	t.Parallel()

	stopCh := newStopCh(t)
	test.PluginTest(t, loop.PluginMercuryName, &loop.GRPCPluginMercury{PluginServer: test.StaticPluginMercury{}, BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}, test.PluginMercury)

	t.Run("proxy", func(t *testing.T) {
		test.PluginTest(t, loop.PluginRelayerName, &loop.GRPCPluginRelayer{PluginServer: test.StaticPluginRelayer{}, BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}, func(t *testing.T, pr loop.PluginRelayer) {
			p := newMercuryProvider(t, pr)
			pm := test.PluginMercuryTest{MercuryProvider: p}
			test.PluginTest(t, loop.PluginMercuryName, &loop.GRPCPluginMercury{PluginServer: test.StaticPluginMercury{}, BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}, pm.TestPluginMercury)
		})
	})
}

func TestPluginMercuryExec(t *testing.T) {
	t.Parallel()
	stopCh := newStopCh(t)
	mercury := loop.GRPCPluginMercury{BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}
	cc := mercury.ClientConfig()
	cc.Cmd = NewHelperProcessCommand(loop.PluginMercuryName, true)
	c := plugin.NewClient(cc)
	t.Cleanup(c.Kill)
	client, err := c.Client()
	require.NoError(t, err)
	defer client.Close()
	require.NoError(t, client.Ping())

	i, err := client.Dispense(loop.PluginMercuryName)
	require.NoError(t, err)
	require.NotNil(t, i)
	test.PluginMercury(t, i.(types.PluginMercury))

	t.Run("proxy", func(t *testing.T) {
		pr := newPluginRelayerExec(t, true, stopCh)
		p := newMercuryProvider(t, pr)
		pm := test.PluginMercuryTest{MercuryProvider: p}
		pm.TestPluginMercury(t, i.(types.PluginMercury))
	})
}

func newMercuryProvider(t *testing.T, pr loop.PluginRelayer) types.MercuryProvider {
	ctx := context.Background()
	r, err := pr.NewRelayer(ctx, test.ConfigTOML, test.StaticKeystore{})
	require.NoError(t, err)
	servicetest.Run(t, r)
	p, err := r.NewPluginProvider(ctx, mercury_common_test.RelayArgs, mercury_common_test.PluginArgs)
	mp, ok := p.(types.MercuryProvider)
	require.True(t, ok)
	require.NoError(t, err)
	servicetest.Run(t, mp)
	return mp
}
