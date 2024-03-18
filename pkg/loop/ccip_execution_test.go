package loop_test

import (
	"context"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal"
	ccip_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/ccip/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	testcore "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/core"
	testreportingplugin "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/ocr2/reporting_plugin"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestExecService(t *testing.T) {
	t.Parallel()

	exec := loop.NewExecutionService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		return NewHelperProcessCommand(loop.CCIPExecutionLOOPName, false, 0)
	}, ccip_test.ExecutionProvider, ccip_test.ExecutionConfig)
	hook := exec.PluginService.XXXTestHook()
	servicetest.Run(t, exec)

	t.Run("control", func(t *testing.T) {
		testreportingplugin.RunFactory(t, exec)
	})

	t.Run("Kill", func(t *testing.T) {
		hook.Kill()

		// wait for relaunch
		time.Sleep(2 * internal.KeepAliveTickDuration)

		testreportingplugin.RunFactory(t, exec)
	})

	t.Run("Reset", func(t *testing.T) {
		hook.Reset()

		// wait for relaunch
		time.Sleep(2 * internal.KeepAliveTickDuration)

		testreportingplugin.RunFactory(t, exec)
	})
}

func TestExecService_recovery(t *testing.T) {
	t.Parallel()
	var limit atomic.Int32
	exec := loop.NewExecutionService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		h := HelperProcessCommand{
			Command: loop.CCIPExecutionLOOPName,
			Limit:   int(limit.Add(1)),
		}
		return h.New()
	}, ccip_test.ExecutionProvider, ccip_test.ExecutionConfig)
	servicetest.Run(t, exec)

	testreportingplugin.RunFactory(t, exec)
}

func TestExecLOOP(t *testing.T) {
	// launch the exec loop via the main program
	t.Parallel()
	stopCh := newStopCh(t)
	exec := loop.ExecutionLoop{BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}
	cc := exec.ClientConfig()
	cc.Cmd = NewHelperProcessCommand(loop.CCIPExecutionLOOPName, false, 0)
	c := plugin.NewClient(cc)
	// make sure to kill the exec loop
	t.Cleanup(c.Kill)
	client, err := c.Client()
	require.NoError(t, err)
	defer client.Close()
	require.NoError(t, client.Ping())
	// get a concrete instance of the exec loop
	instance, err := client.Dispense(loop.CCIPExecutionLOOPName)
	remoteExecFactory := instance.(types.CCIPExecutionFactoryGenerator)
	require.NoError(t, err)

	ccip_test.RunExecutionLOOP(t, remoteExecFactory)

	t.Run("proxy: exec loop <--> relayer loop", func(t *testing.T) {
		// launch the relayer as external process via the main program
		pr := newPluginRelayerExec(t, false, stopCh)
		remoteProvider, err := newExecutionProvider(t, pr)
		require.Error(t, err, "expected error")
		assert.Contains(t, err.Error(), "BCF-3061")
		if err == nil {
			// test to run when BCF-3061 is fixed
			ccip_test.ExecutionLOOPTester{CCIPExecProvider: remoteProvider}.Run(t, remoteExecFactory)
		}
	})
}

func newExecutionProvider(t *testing.T, pr loop.PluginRelayer) (types.CCIPExecProvider, error) {
	ctx := context.Background()
	r, err := pr.NewRelayer(ctx, test.ConfigTOML, testcore.Keystore)
	require.NoError(t, err)
	servicetest.Run(t, r)

	// TODO: fix BCF-3061. we expect an error here until then.
	p, err := r.NewPluginProvider(ctx, ccip_test.ExecutionRelayArgs, ccip_test.ExecutionPluginArgs)
	if err != nil {
		return nil, err
	}
	// TODO: this shouldn't run until BCF-3061 is fixed
	require.NoError(t, err)
	execProvider, ok := p.(types.CCIPExecProvider)
	require.True(t, ok, "got %T", p)
	servicetest.Run(t, execProvider)
	return execProvider, nil
}
