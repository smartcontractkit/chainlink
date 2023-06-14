package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestPluginPortManager(t *testing.T) {
	// register one
	m := plugins.NewLoopRegistry()
	pFoo, err := m.Register("foo")
	require.NoError(t, err)
	require.Equal(t, "foo", pFoo.Name)
	require.Greater(t, pFoo.EnvCfg.PrometheusPort(), 0)
	// test idempotent
	pSame, err := m.Register("foo")
	require.NoError(t, err)
	require.Equal(t, pFoo, pSame)
	// ensure increasing port assignment
	pBar, err := m.Register("bar")
	require.NoError(t, err)
	require.Equal(t, "bar", pBar.Name)
	require.Equal(t, pFoo.EnvCfg.PrometheusPort()+1, pBar.EnvCfg.PrometheusPort())
}
