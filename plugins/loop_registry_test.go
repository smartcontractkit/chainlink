package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v2 "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestPluginPortManager(t *testing.T) {
	m := plugins.NewLoopRegistry()
	pFoo := m.Register("foo", v2.NewTestGeneralConfig(t))
	require.Equal(t, pFoo, plugins.PluginDefaultPort)
	pSame := m.Register("foo", v2.NewTestGeneralConfig(t))
	require.Equal(t, pFoo, pSame)
	pBar := m.Register("bar", v2.NewTestGeneralConfig(t))
	require.Greater(t, pBar, pFoo)
}
