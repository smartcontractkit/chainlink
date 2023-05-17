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

func TestPluginMedian(t *testing.T) {
	t.Parallel()

	stopCh := make(chan struct{})
	if d, ok := t.Deadline(); ok {
		time.AfterFunc(time.Until(d), func() { close(stopCh) })
	}
	testPlugin(t, loop.PluginMedianName, &loop.GRPCPluginMedian{Logger: logger.Test(t), PluginServer: test.StaticPluginMedian{}, StopCh: stopCh}, test.TestPluginMedian)
}

func TestPluginMedianExec(t *testing.T) {
	t.Parallel()
	median := loop.GRPCPluginMedian{Logger: logger.Test(t)}
	cc := median.ClientConfig()
	cc.Cmd = helperProcess(loop.PluginMedianName)
	c := plugin.NewClient(cc)
	client, err := c.Client()
	require.NoError(t, err)
	defer client.Close()
	require.NoError(t, client.Ping())
	i, err := client.Dispense(loop.PluginMedianName)
	require.NoError(t, err)

	test.TestPluginMedian(t, i.(loop.PluginMedian))
}
