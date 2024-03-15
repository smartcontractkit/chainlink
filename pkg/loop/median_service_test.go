package loop_test

import (
	"os/exec"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal"
	median_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/median/test"
	testcore "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/core"
	testreportingplugin "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/ocr2/reporting_plugin"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
)

func TestMedianService(t *testing.T) {
	t.Parallel()

	median := loop.NewMedianService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		return NewHelperProcessCommand(loop.PluginMedianName, false, 0)
	}, median_test.MedianProvider, median_test.DataSource, median_test.JuelsPerFeeCoinDataSource, testcore.ErrorLog)
	hook := median.PluginService.XXXTestHook()
	servicetest.Run(t, median)

	t.Run("control", func(t *testing.T) {
		testreportingplugin.RunFactory(t, median)
	})

	t.Run("Kill", func(t *testing.T) {
		hook.Kill()

		// wait for relaunch
		time.Sleep(2 * internal.KeepAliveTickDuration)

		testreportingplugin.RunFactory(t, median)
	})

	t.Run("Reset", func(t *testing.T) {
		hook.Reset()

		// wait for relaunch
		time.Sleep(2 * internal.KeepAliveTickDuration)

		testreportingplugin.RunFactory(t, median)
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
	}, median_test.MedianProvider, median_test.DataSource, median_test.JuelsPerFeeCoinDataSource, &testcore.ErrorLog)
	servicetest.Run(t, median)

	testreportingplugin.RunFactory(t, median)
}
