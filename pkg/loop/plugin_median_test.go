package loop_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	keystoretest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/keystore/test"
	mediantest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/median/test"
	relayertest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestPluginMedian(t *testing.T) {
	t.Parallel()

	stopCh := newStopCh(t)
	t.Run("no proxy", func(t *testing.T) {
		test.PluginTest(t, loop.PluginMedianName,
			&loop.GRPCPluginMedian{
				PluginServer: mediantest.MedianFactoryServer,
				BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh},
			},
			mediantest.PluginMedian)
	})

	t.Run("proxy", func(t *testing.T) {
		test.PluginTest(t, loop.PluginRelayerName,
			&loop.GRPCPluginRelayer{
				PluginServer: relayertest.NewRelayerTester(false),
				BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}},
			func(t *testing.T, pr loop.PluginRelayer) {
				p := newMedianProvider(t, pr)
				pm := mediantest.PluginMedianTest{MedianProvider: p}
				test.PluginTest(t, loop.PluginMedianName,
					&loop.GRPCPluginMedian{
						PluginServer: mediantest.MedianFactoryServer,
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

	mediantest.PluginMedian(t, i.(types.PluginMedian))

	t.Run("proxy", func(t *testing.T) {
		pr := newPluginRelayerExec(t, false, stopCh)
		p := newMedianProvider(t, pr)
		pm := mediantest.PluginMedianTest{MedianProvider: p}
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
	r, err := pr.NewRelayer(ctx, test.ConfigTOML, keystoretest.Keystore)
	require.NoError(t, err)
	servicetest.Run(t, r)
	p, err := r.NewPluginProvider(ctx, relayertest.RelayArgs, relayertest.PluginArgs)
	mp, ok := p.(types.MedianProvider)
	require.True(t, ok)
	require.NoError(t, err)
	servicetest.Run(t, mp)
	return mp
}

func newGenericPluginProvider(t *testing.T, pr loop.PluginRelayer) types.PluginProvider {
	ctx := context.Background()
	r, err := pr.NewRelayer(ctx, test.ConfigTOML, keystoretest.Keystore)
	require.NoError(t, err)
	servicetest.Run(t, r)
	ra := relayertest.RelayArgs
	ra.ProviderType = string(types.GenericPlugin)
	p, err := r.NewPluginProvider(ctx, ra, relayertest.PluginArgs)
	require.NoError(t, err)
	servicetest.Run(t, p)
	return p
}
