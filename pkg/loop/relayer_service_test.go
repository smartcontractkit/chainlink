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
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
)

func TestRelayerService(t *testing.T) {
	t.Parallel()
	relayer := loop.NewRelayerService(logger.Test(t), func() *exec.Cmd {
		return helperProcess(loop.PluginRelayerName)
	}, configTOML, staticKeystore{})
	hook := relayer.TestHook()
	require.NoError(t, relayer.Start(utils.Context(t)))
	t.Cleanup(func() { assert.NoError(t, relayer.Close()) })

	t.Run("control", func(t *testing.T) {
		testRelayer(t, relayer)
	})

	t.Run("Kill", func(t *testing.T) {
		hook.Kill()

		// wait for relaunch
		time.Sleep(2 * loop.KeepAliveTickDuration)

		testRelayer(t, relayer)
	})

	t.Run("Reset", func(t *testing.T) {
		hook.Reset()

		// wait for relaunch
		time.Sleep(2 * loop.KeepAliveTickDuration)

		testRelayer(t, relayer)
	})
}

func TestRelayerService_recovery(t *testing.T) {
	t.Parallel()
	var limit atomic.Int32
	relayer := loop.NewRelayerService(logger.Test(t), func() *exec.Cmd {
		return helperProcess(loop.PluginRelayerName, strconv.Itoa(int(limit.Add(1))))
	}, configTOML, staticKeystore{})
	require.NoError(t, relayer.Start(utils.Context(t)))
	t.Cleanup(func() { assert.NoError(t, relayer.Close()) })

	testRelayer(t, relayer)
}
