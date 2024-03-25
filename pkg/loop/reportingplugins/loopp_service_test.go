package reportingplugins_test

import (
	"os/exec"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	testcore "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/core"
	testapi "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/core/api"
	testreportingplugin "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/ocr2/reporting_plugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type HelperProcessCommand test.HelperProcessCommand

func (h *HelperProcessCommand) New() *exec.Cmd {
	h.CommandLocation = "../internal/test/cmd/main.go"
	return (test.HelperProcessCommand)(*h).New()
}

func NewHelperProcessCommand(command string) *exec.Cmd {
	h := HelperProcessCommand{
		Command: command,
	}
	return h.New()
}

func TestLOOPPService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Plugin string
	}{
		// A generic plugin with a median provider
		{Plugin: testapi.MedianID},
		// A generic plugin with a plugin provider
		{Plugin: reportingplugins.PluginServiceName},
	}
	for _, ts := range tests {
		looppSvc := reportingplugins.NewLOOPPService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
			return NewHelperProcessCommand(ts.Plugin)
		},
			types.ReportingPluginServiceConfig{},
			testcore.MockConn{},
			testcore.PipelineRunner,
			testcore.Telemetry,
			&testcore.ErrorLog)
		hook := looppSvc.XXXTestHook()
		servicetest.Run(t, looppSvc)

		t.Run("control", func(t *testing.T) {
			testreportingplugin.RunFactory(t, looppSvc)
		})

		t.Run("Kill", func(t *testing.T) {
			hook.Kill()

			// wait for relaunch
			time.Sleep(2 * goplugin.KeepAliveTickDuration)

			testreportingplugin.RunFactory(t, looppSvc)
		})

		t.Run("Reset", func(t *testing.T) {
			hook.Reset()

			// wait for relaunch
			time.Sleep(2 * goplugin.KeepAliveTickDuration)

			testreportingplugin.RunFactory(t, looppSvc)
		})
	}
}

func TestLOOPPService_recovery(t *testing.T) {
	t.Parallel()
	var limit atomic.Int32
	looppSvc := reportingplugins.NewLOOPPService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		h := HelperProcessCommand{
			Command: testapi.MedianID,
			Limit:   int(limit.Add(1)),
		}
		return h.New()
	},
		types.ReportingPluginServiceConfig{},
		testcore.MockConn{},
		testcore.PipelineRunner,
		testcore.Telemetry,
		&testcore.ErrorLog)
	servicetest.Run(t, looppSvc)

	testreportingplugin.RunFactory(t, looppSvc)
}
