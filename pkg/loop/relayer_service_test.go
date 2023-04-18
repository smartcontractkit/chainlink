package loop_test

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
)

func TestPluginService(t *testing.T) {
	t.Parallel()
	ps := loop.NewRelayerService(logger.Test(t), func() *exec.Cmd {
		return helperProcess(loop.PluginRelayerName)
	}, configTOML, staticKeystore{})
	require.NoError(t, ps.Start(context.Background()))
	t.Cleanup(func() {
		assert.NoError(t, ps.Close())
	})

	t.Run("Start", func(t *testing.T) {
		relayer, err := ps.Relayer()
		require.NoError(t, err)
		testRelayer(t, relayer)
	})

	t.Run("Kill", func(t *testing.T) {
		ps.Kill()

		// wait for relaunch
		time.Sleep(2 * loop.KeepAliveTickDuration)

		relayer, err := ps.Relayer()
		require.NoError(t, err)
		testRelayer(t, relayer)
	})

	t.Run("Reset", func(t *testing.T) {
		ps.Reset()

		// wait for relaunch
		time.Sleep(2 * loop.KeepAliveTickDuration)

		relayer, err := ps.Relayer()
		require.NoError(t, err)
		testRelayer(t, relayer)
	})
}
