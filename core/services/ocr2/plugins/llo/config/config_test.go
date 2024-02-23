package config

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Config(t *testing.T) {
	t.Run("unmarshals from toml", func(t *testing.T) {
		cdjson := `{
	"42": {
		"reportFormat": 42,
		"chainSelector": 142,
		"streamIds": [1, 2]
	},
	"43": {
		"reportFormat": 42,
		"chainSelector": 142,
		"streamIds": [1, 3]
	},
	"44": {
		"reportFormat": 42,
		"chainSelector": 143,
		"streamIds": [1, 4]
	}
}`

		t.Run("with all possible values set", func(t *testing.T) {
			rawToml := fmt.Sprintf(`
				ServerURL = "example.com:80"
				ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
				BenchmarkMode = true
				ChannelDefinitionsContractAddress = "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
				ChannelDefinitionsContractFromBlock = 1234
				ChannelDefinitions = """
%s
"""`, cdjson)

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			assert.Equal(t, "example.com:80", mc.RawServerURL)
			assert.Equal(t, "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.ServerPubKey.String())
			assert.Equal(t, "0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF", mc.ChannelDefinitionsContractAddress.Hex())
			assert.Equal(t, int64(1234), mc.ChannelDefinitionsContractFromBlock)
			assert.JSONEq(t, cdjson, mc.ChannelDefinitions)
			assert.True(t, mc.BenchmarkMode)

			err = mc.Validate()
			require.Error(t, err)

			assert.Contains(t, err.Error(), "llo: ChannelDefinitionsContractAddress is not allowed if ChannelDefinitions is specified")
			assert.Contains(t, err.Error(), "llo: ChannelDefinitionsContractFromBlock is not allowed if ChannelDefinitions is specified")
		})

		t.Run("with only channelDefinitions", func(t *testing.T) {
			rawToml := fmt.Sprintf(`
				ServerURL = "example.com:80"
				ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
				ChannelDefinitions = """
%s
"""`, cdjson)

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			assert.Equal(t, "example.com:80", mc.RawServerURL)
			assert.Equal(t, "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.ServerPubKey.String())
			assert.JSONEq(t, cdjson, mc.ChannelDefinitions)
			assert.False(t, mc.BenchmarkMode)

			err = mc.Validate()
			require.NoError(t, err)
		})
		t.Run("with only channelDefinitions contract details", func(t *testing.T) {
			rawToml := `
			ServerURL = "example.com:80"
			ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
			ChannelDefinitionsContractAddress = "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			assert.Equal(t, "example.com:80", mc.RawServerURL)
			assert.Equal(t, "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.ServerPubKey.String())
			assert.Equal(t, "0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF", mc.ChannelDefinitionsContractAddress.Hex())
			assert.False(t, mc.BenchmarkMode)

			err = mc.Validate()
			require.NoError(t, err)
		})
		t.Run("with missing ChannelDefinitionsContractAddress", func(t *testing.T) {
			rawToml := `
			ServerURL = "example.com:80"
			ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			assert.Equal(t, "example.com:80", mc.RawServerURL)
			assert.Equal(t, "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.ServerPubKey.String())
			assert.False(t, mc.BenchmarkMode)

			err = mc.Validate()
			require.Error(t, err)
			assert.EqualError(t, err, "llo: ChannelDefinitionsContractAddress is required if ChannelDefinitions is not specified")
		})

		t.Run("with invalid values", func(t *testing.T) {
			rawToml := `
				ChannelDefinitionsContractFromBlock = "invalid"
			`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.Error(t, err)
			assert.EqualError(t, err, `toml: cannot decode TOML string into struct field config.PluginConfig.ChannelDefinitionsContractFromBlock of type int64`)
			assert.False(t, mc.BenchmarkMode)

			rawToml = `
				ServerURL = "http://example.com"
				ServerPubKey = "4242"
			`

			err = toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			err = mc.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), `invalid scheme specified for MercuryServer, got: "http://example.com" (scheme: "http") but expected a websocket url e.g. "192.0.2.2:4242" or "wss://192.0.2.2:4242"`)
			assert.Contains(t, err.Error(), `ServerPubKey is required and must be a 32-byte hex string`)
		})
	})
}
