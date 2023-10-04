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

func TestRelayerService(t *testing.T) {
	t.Parallel()
	relayer := loop.NewRelayerService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		return helperProcess(loop.PluginRelayerName)
	}, test.ConfigTOML, test.StaticKeystore{})
	hook := relayer.TestHook()
	require.NoError(t, relayer.Start(tests.Context(t)))
	t.Cleanup(func() { assert.NoError(t, relayer.Close()) })

	t.Run("control", func(t *testing.T) {
		test.TestRelayer(t, relayer)
	})

	t.Run("Kill", func(t *testing.T) {
		hook.Kill()

		// wait for relaunch
		time.Sleep(2 * loop.KeepAliveTickDuration)

		test.TestRelayer(t, relayer)
	})

	t.Run("Reset", func(t *testing.T) {
		hook.Reset()

		// wait for relaunch
		time.Sleep(2 * loop.KeepAliveTickDuration)

		test.TestRelayer(t, relayer)
	})
}

func TestRelayerService_recovery(t *testing.T) {
	t.Parallel()
	var limit atomic.Int32
	relayer := loop.NewRelayerService(logger.Test(t), loop.GRPCOpts{}, func() *exec.Cmd {
		return helperProcess(loop.PluginRelayerName, strconv.Itoa(int(limit.Add(1))))
	}, test.ConfigTOML, test.StaticKeystore{})
	require.NoError(t, relayer.Start(tests.Context(t)))
	t.Cleanup(func() { assert.NoError(t, relayer.Close()) })

	test.TestRelayer(t, relayer)
}
