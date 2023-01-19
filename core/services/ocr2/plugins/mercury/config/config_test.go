package config

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PluginConfig(t *testing.T) {
	t.Run("with valid values", func(t *testing.T) {
		rawToml := `
ServerHost = "example.com:80"
ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
`

		var mc PluginConfig
		err := toml.Unmarshal([]byte(rawToml), &mc)
		require.NoError(t, err)

		assert.Equal(t, "example.com:80", mc.ServerHost)
		assert.Equal(t, "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.ServerPubKey.String())

		err = ValidatePluginConfig(mc)
		require.NoError(t, err)
	})

	t.Run("invalid values", func(t *testing.T) {
		rawToml := `
ServerHost = "http://example.com"
ServerPubKey = "4242"
`

		var mc PluginConfig
		err := toml.Unmarshal([]byte(rawToml), &mc)
		require.NoError(t, err)

		err = ValidatePluginConfig(mc)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `Mercury: invalid value specified for MercuryServer, got :http://example.com but expected value in the form of "address:port" e.g. "192.0.2.2:4242"`)
		assert.Contains(t, err.Error(), `Mercury: ServerPubKey is required and must be a 32-byte hex string`)
	})
}
