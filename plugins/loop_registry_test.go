package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestPluginPortManager(t *testing.T) {
	lc := plugins.NewLoggingConfig(zapcore.DebugLevel, false, false)
	// register one
	m := plugins.NewLoopRegistry()
	pFoo := m.Register("foo", lc)
	require.Equal(t, "foo", pFoo.Name)
	require.Equal(t, pFoo.EnvCfg.PrometheusPort(), plugins.PluginDefaultPort)
	require.Equal(t, lc.JSONConsole(), pFoo.EnvCfg.JSONConsole())
	require.Equal(t, lc.LogLevel(), pFoo.EnvCfg.LogLevel())
	require.Equal(t, lc.LogUnixTimestamps(), pFoo.EnvCfg.LogUnixTimestamps())
	// test idempotent
	pSame := m.Register("foo", lc)
	require.Equal(t, pFoo, pSame)
	// ensure increasing port assignment
	pBar := m.Register("bar", lc)
	require.Equal(t, "bar", pBar.Name)
	require.Greater(t, pBar.EnvCfg.PrometheusPort(), pFoo.EnvCfg.PrometheusPort())
}
