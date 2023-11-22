package loop_test

import (
	"os/exec"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestMedianService(t *testing.T) {
	t.Parallel()
	median := loop.NewMedianService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		return NewHelperProcessCommand(loop.PluginMedianName)
	}, test.StaticMedianProvider{}, test.StaticDataSource(), test.StaticJuelsPerFeeCoinDataSource(), &test.StaticErrorLog{})
	hook := median.PluginService.XXXTestHook()
	require.NoError(t, median.Start(tests.Context(t)))
	t.Cleanup(func() { assert.NoError(t, median.Close()) })

	t.Run("control", func(t *testing.T) {
		test.ReportingPluginFactory(t, median)
	})

	t.Run("Kill", func(t *testing.T) {
		hook.Kill()

		// wait for relaunch
		time.Sleep(2 * internal.KeepAliveTickDuration)

		test.ReportingPluginFactory(t, median)
	})

	t.Run("Reset", func(t *testing.T) {
		hook.Reset()

		// wait for relaunch
		time.Sleep(2 * internal.KeepAliveTickDuration)

		test.ReportingPluginFactory(t, median)
	})
}

func TestMedianService_recovery(t *testing.T) {
	t.Parallel()
	var limit atomic.Int32
	median := loop.NewMedianService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		h := HelperProcessCommand{
			Command: loop.PluginMedianName,
			Limit:   int(limit.Add(1)),
		}
		return h.New()
	}, test.StaticMedianProvider{}, test.StaticDataSource(), test.StaticJuelsPerFeeCoinDataSource(), &test.StaticErrorLog{})
	require.NoError(t, median.Start(tests.Context(t)))
	t.Cleanup(func() { assert.NoError(t, median.Close()) })

	test.ReportingPluginFactory(t, median)
}
