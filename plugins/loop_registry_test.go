package plugins_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestPluginPortManager(t *testing.T) {
	c := plugins.NewAppConfig(uuid.New(), zapcore.DebugLevel, false, false)
	// register one
	m := plugins.NewLoopRegistry()
	pFoo, err := m.Register("foo", c)
	require.NoError(t, err)
	require.Equal(t, "foo", pFoo.Name)
	require.Greater(t, pFoo.EnvCfg.PrometheusPort(), 0)
	require.Equal(t, c.JSONConsole(), pFoo.EnvCfg.JSONConsole())
	require.Equal(t, c.LogLevel(), pFoo.EnvCfg.LogLevel())
	require.Equal(t, c.LogUnixTimestamps(), pFoo.EnvCfg.LogUnixTimestamps())
	// test idempotent
	pSame, err := m.Register("foo", c)
	require.NoError(t, err)
	require.Equal(t, pFoo, pSame)
	// ensure increasing port assignment
	pBar, err := m.Register("bar", c)
	require.NoError(t, err)
	require.Equal(t, "bar", pBar.Name)
	require.Equal(t, pFoo.EnvCfg.PrometheusPort()+1, pBar.EnvCfg.PrometheusPort())
}
