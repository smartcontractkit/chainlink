package loop_test

import (
	"os/exec"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils/tests"
)

func TestMedianService(t *testing.T) {
	t.Parallel()
	median := loop.NewMedianService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		return helperProcess(loop.PluginMedianName)
	}, test.StaticMedianProvider{}, test.StaticDataSource(), test.StaticJuelsPerFeeCoinDataSource(), &test.StaticErrorLog{})
	hook := median.TestHook()
	require.NoError(t, median.Start(tests.Context(t)))
	t.Cleanup(func() { assert.NoError(t, median.Close()) })

	t.Run("control", func(t *testing.T) {
		test.TestReportingPluginFactory(t, median)
	})

	t.Run("Kill", func(t *testing.T) {
		hook.Kill()

		// wait for relaunch
		time.Sleep(2 * loop.KeepAliveTickDuration)

		test.TestReportingPluginFactory(t, median)
	})

	t.Run("Reset", func(t *testing.T) {
		hook.Reset()

		// wait for relaunch
		time.Sleep(2 * loop.KeepAliveTickDuration)

		test.TestReportingPluginFactory(t, median)
	})
}

func TestMedianService_recovery(t *testing.T) {
	t.Parallel()
	var limit atomic.Int32
	median := loop.NewMedianService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		return helperProcess(loop.PluginMedianName, strconv.Itoa(int(limit.Add(1))))
	}, test.StaticMedianProvider{}, test.StaticDataSource(), test.StaticJuelsPerFeeCoinDataSource(), &test.StaticErrorLog{})
	require.NoError(t, median.Start(tests.Context(t)))
	t.Cleanup(func() { assert.NoError(t, median.Close()) })

	test.TestReportingPluginFactory(t, median)
}
