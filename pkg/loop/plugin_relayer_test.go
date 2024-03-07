package loop_test

import (
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	testcore "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/core"
	relayer_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/relayer"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestPluginRelayer(t *testing.T) {
	t.Parallel()

	stopCh := newStopCh(t)
	test.PluginTest(t, loop.PluginRelayerName,
		&loop.GRPCPluginRelayer{
			PluginServer: relayer_test.NewRelayerTester(false),
			BrokerConfig: loop.BrokerConfig{
				Logger: logger.Test(t),
				StopCh: stopCh}},
		relayer_test.RunPlugin)
}

func TestPluginRelayerExec(t *testing.T) {
	t.Parallel()
	stopCh := newStopCh(t)

	pr := newPluginRelayerExec(t, true, stopCh)

	relayer_test.RunPlugin(t, pr)
}

func FuzzPluginRelayer(f *testing.F) {
	testFunc := func(t *testing.T) loop.PluginRelayer {
		t.Helper()

		stopCh := newStopCh(t)
		relayer := newPluginRelayerExec(t, true, stopCh)

		return relayer
	}

	relayer_test.RunFuzzPluginRelayer(f, testFunc)
}

func FuzzRelayer(f *testing.F) {
	testFunc := func(t *testing.T) loop.Relayer {
		t.Helper()

		stopCh := newStopCh(t)
		p := newPluginRelayerExec(t, false, stopCh)
		ctx := tests.Context(t)
		relayer, err := p.NewRelayer(ctx, test.ConfigTOML, testcore.Keystore)

		require.NoError(t, err)

		return relayer
	}

	relayer_test.RunFuzzRelayer(f, testFunc)
}

func newPluginRelayerExec(t *testing.T, staticChecks bool, stopCh <-chan struct{}) loop.PluginRelayer {
	relayer := loop.GRPCPluginRelayer{BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}
	cc := relayer.ClientConfig()
	cc.Cmd = NewHelperProcessCommand(loop.PluginRelayerName, staticChecks, 0)
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
