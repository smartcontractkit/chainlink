package config

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var feedId = [32]uint8{00, 01, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}

func Test_PluginConfig(t *testing.T) {
	t.Run("with valid values", func(t *testing.T) {
		rawToml := `
ServerURL = "example.com:80"
ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
`

		var mc PluginConfig
		err := toml.Unmarshal([]byte(rawToml), &mc)
		require.NoError(t, err)

		assert.Equal(t, "example.com:80", mc.RawServerURL)
		assert.Equal(t, "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.ServerPubKey.String())

		err = ValidatePluginConfig(mc, feedId)
		require.NoError(t, err)
	})

	t.Run("invalid values", func(t *testing.T) {
		rawToml := `
ServerURL = "http://example.com"
ServerPubKey = "4242"
`

		var mc PluginConfig
		err := toml.Unmarshal([]byte(rawToml), &mc)
		require.NoError(t, err)

		err = ValidatePluginConfig(mc, feedId)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `Mercury: invalid scheme specified for MercuryServer, got: "http://example.com" (scheme: "http") but expected a websocket url e.g. "192.0.2.2:4242" or "wss://192.0.2.2:4242"`)
		assert.Contains(t, err.Error(), `mercury: ServerPubKey is required and must be a 32-byte hex string`)
	})

	t.Skip("TODO: linkFeedID, nativeFeedID, reportSchemaVersion, initialBlockNumber")
}

func Test_PluginConfig_ServerURL(t *testing.T) {
	pc := PluginConfig{RawServerURL: "example.com"}
	assert.Equal(t, "example.com", pc.ServerURL())
	pc = PluginConfig{RawServerURL: "wss://example.com"}
	assert.Equal(t, "example.com", pc.ServerURL())
	pc = PluginConfig{RawServerURL: "example.com:1234/foo"}
	assert.Equal(t, "example.com:1234/foo", pc.ServerURL())
	pc = PluginConfig{RawServerURL: "wss://example.com:1234/foo"}
	assert.Equal(t, "example.com:1234/foo", pc.ServerURL())
}
