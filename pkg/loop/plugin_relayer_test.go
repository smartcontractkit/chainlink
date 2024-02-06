package loop_test

import (
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestPluginRelayer(t *testing.T) {
	t.Parallel()

	stopCh := newStopCh(t)
	test.PluginTest(t, loop.PluginRelayerName, &loop.GRPCPluginRelayer{PluginServer: test.StaticPluginRelayer{}, BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}, test.RunPluginRelayer)
}

func TestPluginRelayerExec(t *testing.T) {
	t.Parallel()
	stopCh := newStopCh(t)

	pr := newPluginRelayerExec(t, true, stopCh)

	test.RunPluginRelayer(t, pr)
}

func FuzzPluginRelayer(f *testing.F) {
	testFunc := func(t *testing.T) loop.PluginRelayer {
		t.Helper()

		stopCh := newStopCh(t)
		relayer := newPluginRelayerExec(t, true, stopCh)

		return relayer
	}

	test.RunFuzzPluginRelayer(f, testFunc)
}

func FuzzRelayer(f *testing.F) {
	testFunc := func(t *testing.T) loop.Relayer {
		t.Helper()

		stopCh := newStopCh(t)
		p := newPluginRelayerExec(t, false, stopCh)
		ctx := tests.Context(t)
		relayer, err := p.NewRelayer(ctx, test.ConfigTOML, test.StaticKeystore{})

		require.NoError(t, err)

		return relayer
	}

	test.RunFuzzRelayer(f, testFunc)
}

func newPluginRelayerExec(t *testing.T, staticChecks bool, stopCh <-chan struct{}) loop.PluginRelayer {
	relayer := loop.GRPCPluginRelayer{BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}
	cc := relayer.ClientConfig()
	cc.Cmd = NewHelperProcessCommand(loop.PluginRelayerName, staticChecks)
	c := plugin.NewClient(cc)
	t.Cleanup(c.Kill)
	client, err := c.Client()
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	require.NoError(t, client.Ping())
	i, err := client.Dispense(loop.PluginRelayerName)
	require.NoError(t, err)
	return i.(loop.PluginRelayer)
}
