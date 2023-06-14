package plugins_test

import (
	"testing"

	"github.com/test-go/testify/require"

	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestHostname(t *testing.T) {
	t.Run("override host", func(t *testing.T) {
		t.Setenv(string(v2.EnvPluginPromTarget), "fakeHost")
		c := plugins.NewEnvConfig(22)
		require.Equal(t, 22, c.PrometheusPort())
		require.Equal(t, "fakeHost", c.Hostname())
	})

	t.Run("localhost", func(t *testing.T) {
		t.Setenv("HOSTNAME", "localhost")
		c := plugins.NewEnvConfig(22)
		require.Equal(t, 22, c.PrometheusPort())
		require.Equal(t, "localhost", c.Hostname())
	})

}
