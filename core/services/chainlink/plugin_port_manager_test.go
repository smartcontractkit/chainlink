package chainlink_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func TestPluginPortManager(t *testing.T) {
	m := chainlink.NewPluginPortManager()
	pFoo := m.Register("foo")
	require.Equal(t, pFoo, chainlink.PluginDefaultPort)
	pSame := m.Register("foo")
	require.Equal(t, pFoo, pSame)
	pBar := m.Register("bar")
	require.Greater(t, pBar, pFoo)
}
