package loop_test

import (
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/test"
)

func TestPluginRelayer(t *testing.T) {
	t.Parallel()

	stopCh := make(chan struct{})
	if d, ok := t.Deadline(); ok {
		time.AfterFunc(time.Until(d), func() { close(stopCh) })
	}
	testPlugin(t, loop.PluginRelayerName, &loop.GRPCPluginRelayer{Logger: logger.Test(t), PluginServer: test.StaticPluginRelayer{}, StopCh: stopCh}, test.TestPluginRelayer)
}

func TestPluginRelayerExec(t *testing.T) {
	t.Parallel()
	relayer := loop.GRPCPluginRelayer{Logger: logger.Test(t)}
	cc := relayer.ClientConfig()
	cc.Cmd = helperProcess(loop.PluginRelayerName)
	c := plugin.NewClient(cc)
	client, err := c.Client()
	require.NoError(t, err)
	defer client.Close()
	require.NoError(t, client.Ping())
	i, err := client.Dispense(loop.PluginRelayerName)
	require.NoError(t, err)

	test.TestPluginRelayer(t, i.(loop.PluginRelayer))
}
