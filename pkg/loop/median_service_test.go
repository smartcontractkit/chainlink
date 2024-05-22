package loop_test

import (
	"os/exec"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	errorlogtest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/errorlog/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	mediantest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/median/test"
	reportingplugintest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/reportingplugin/test"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
)

func TestMedianService(t *testing.T) {
	t.Parallel()

	median := loop.NewMedianService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		return NewHelperProcessCommand(loop.PluginMedianName, false, 0)
	}, mediantest.MedianProvider, mediantest.DataSource, mediantest.JuelsPerFeeCoinDataSource, mediantest.GasPriceSubunitsDataSource, errorlogtest.ErrorLog)
	hook := median.PluginService.XXXTestHook()
	servicetest.Run(t, median)

	t.Run("control", func(t *testing.T) {
		reportingplugintest.RunFactory(t, median)
	})

	t.Run("Kill", func(t *testing.T) {
		hook.Kill()

		// wait for relaunch
		time.Sleep(2 * goplugin.KeepAliveTickDuration)

		reportingplugintest.RunFactory(t, median)
	})

	t.Run("Reset", func(t *testing.T) {
		hook.Reset()

		// wait for relaunch
		time.Sleep(2 * goplugin.KeepAliveTickDuration)

		reportingplugintest.RunFactory(t, median)
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
	}, mediantest.MedianProvider, mediantest.DataSource, mediantest.JuelsPerFeeCoinDataSource, mediantest.GasPriceSubunitsDataSource, errorlogtest.ErrorLog)
	servicetest.Run(t, median)

	reportingplugintest.RunFactory(t, median)
}
