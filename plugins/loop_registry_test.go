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
	pFoo, err := m.Register("foo", lc)
	require.NoError(t, err)
	require.Equal(t, "foo", pFoo.Name)
	require.Greater(t, pFoo.EnvCfg.PrometheusPort(), 0)
	require.Equal(t, lc.JSONConsole(), pFoo.EnvCfg.JSONConsole())
	require.Equal(t, lc.Level(), pFoo.EnvCfg.Level())
	require.Equal(t, lc.UnixTimestamps(), pFoo.EnvCfg.UnixTimestamps())
	// test idempotent
	pSame, err := m.Register("foo", lc)
	require.NoError(t, err)
	require.Equal(t, pFoo, pSame)
	// ensure increasing port assignment
	pBar, err := m.Register("bar", lc)
	require.NoError(t, err)
	require.Equal(t, "bar", pBar.Name)
	require.Equal(t, pFoo.EnvCfg.PrometheusPort()+1, pBar.EnvCfg.PrometheusPort())
}
