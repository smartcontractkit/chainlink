package loop_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	median_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/median/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	testcore "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/core"
	relayer_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/relayer"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestPluginMedian(t *testing.T) {
	t.Parallel()

	stopCh := newStopCh(t)
	test.PluginTest(t, loop.PluginMedianName,
		&loop.GRPCPluginMedian{
			PluginServer: median_test.MedianFactoryServer,
			BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh},
		},
		median_test.PluginMedian)

	t.Run("proxy", func(t *testing.T) {
		test.PluginTest(t, loop.PluginRelayerName,
			&loop.GRPCPluginRelayer{
				PluginServer: relayer_test.NewRelayerTester(false),
				BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}},
			func(t *testing.T, pr loop.PluginRelayer) {
				p := newMedianProvider(t, pr)
				pm := median_test.PluginMedianTest{MedianProvider: p}
				test.PluginTest(t, loop.PluginMedianName,
					&loop.GRPCPluginMedian{
						PluginServer: median_test.MedianFactoryServer,
						BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}},
					pm.TestPluginMedian)
			})
	})
}

func TestPluginMedianExec(t *testing.T) {
	t.Parallel()
	stopCh := newStopCh(t)
	median := loop.GRPCPluginMedian{BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}
	cc := median.ClientConfig()
	cc.Cmd = NewHelperProcessCommand(loop.PluginMedianName, false, 0)
	c := plugin.NewClient(cc)
	t.Cleanup(c.Kill)
	client, err := c.Client()
	require.NoError(t, err)
	defer client.Close()
	require.NoError(t, client.Ping())
	i, err := client.Dispense(loop.PluginMedianName)
	require.NoError(t, err)

	median_test.PluginMedian(t, i.(types.PluginMedian))

	t.Run("proxy", func(t *testing.T) {
		pr := newPluginRelayerExec(t, false, stopCh)
		p := newMedianProvider(t, pr)
		pm := median_test.PluginMedianTest{MedianProvider: p}
		pm.TestPluginMedian(t, i.(types.PluginMedian))
	})
}

func newStopCh(t *testing.T) <-chan struct{} {
	stopCh := make(chan struct{})
	if d, ok := t.Deadline(); ok {
		time.AfterFunc(time.Until(d), func() { close(stopCh) })
	}
	return stopCh
}

func newMedianProvider(t *testing.T, pr loop.PluginRelayer) types.MedianProvider {
	ctx := context.Background()
	r, err := pr.NewRelayer(ctx, test.ConfigTOML, testcore.Keystore)
	require.NoError(t, err)
	servicetest.Run(t, r)
	p, err := r.NewPluginProvider(ctx, relayer_test.RelayArgs, relayer_test.PluginArgs)
	mp, ok := p.(types.MedianProvider)
	require.True(t, ok)
	require.NoError(t, err)
	servicetest.Run(t, mp)
	return mp
}

func newGenericPluginProvider(t *testing.T, pr loop.PluginRelayer) types.PluginProvider {
	ctx := context.Background()
	r, err := pr.NewRelayer(ctx, test.ConfigTOML, testcore.Keystore)
	require.NoError(t, err)
	servicetest.Run(t, r)
	ra := relayer_test.RelayArgs
	ra.ProviderType = string(types.GenericPlugin)
	p, err := r.NewPluginProvider(ctx, ra, relayer_test.PluginArgs)
	require.NoError(t, err)
	servicetest.Run(t, p)
	return p
}
