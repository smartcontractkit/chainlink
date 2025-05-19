package config

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var v1FeedId = [32]uint8{00, 01, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
var v2FeedId = [32]uint8{00, 02, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}

func Test_PluginConfig(t *testing.T) {
	t.Run("Mercury v1", func(t *testing.T) {
		t.Run("with valid values", func(t *testing.T) {
			rawToml := `
				ServerURL = "example.com:80"
				ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
				InitialBlockNumber = 1234
			`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			assert.Equal(t, "example.com:80", mc.RawServerURL)
			assert.Equal(t, "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.ServerPubKey.String())
			assert.Equal(t, int64(1234), mc.InitialBlockNumber.Int64)

			err = ValidatePluginConfig(mc, v1FeedId)
			require.NoError(t, err)
		})
		t.Run("with multiple server URLs", func(t *testing.T) {
			t.Run("if no ServerURL/ServerPubKey is specified", func(t *testing.T) {
				rawToml := `
					Servers = { "example.com:80" = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", "example2.invalid:1234" = "524ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93" }
			`

				var mc PluginConfig
				err := toml.Unmarshal([]byte(rawToml), &mc)
				require.NoError(t, err)

				assert.Len(t, mc.Servers, 2)
				assert.Equal(t, "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.Servers["example.com:80"].String())
				assert.Equal(t, "524ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", mc.Servers["example2.invalid:1234"].String())

				err = ValidatePluginConfig(mc, v1FeedId)
				require.NoError(t, err)
			})
			t.Run("if ServerURL or ServerPubKey is specified", func(t *testing.T) {
				rawToml := `
					Servers = { "example.com:80" = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", "example2.invalid:1234" = "524ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93" }
					ServerURL = "example.com:80"
			`
				var mc PluginConfig
				err := toml.Unmarshal([]byte(rawToml), &mc)
				require.NoError(t, err)

				err = ValidatePluginConfig(mc, v1FeedId)
				require.EqualError(t, err, "Mercury: Servers and RawServerURL/ServerPubKey may not be specified together")

				rawToml = `
					Servers = { "example.com:80" = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", "example2.invalid:1234" = "524ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93" }
					ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
			`
				err = toml.Unmarshal([]byte(rawToml), &mc)
				require.NoError(t, err)

				err = ValidatePluginConfig(mc, v1FeedId)
				require.EqualError(t, err, "Mercury: Servers and RawServerURL/ServerPubKey may not be specified together")
			})
		})

		t.Run("with invalid values", func(t *testing.T) {
			rawToml := `
				InitialBlockNumber = "invalid"
			`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.Error(t, err)
			assert.EqualError(t, err, `toml: strconv.ParseInt: parsing "invalid": invalid syntax`)

			rawToml = `
				ServerURL = "http://example.com"
				ServerPubKey = "4242"
			`

			err = toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			err = ValidatePluginConfig(mc, v1FeedId)
			require.Error(t, err)
			assert.Contains(t, err.Error(), `Mercury: invalid scheme specified for MercuryServer, got: "http://example.com" (scheme: "http") but expected a websocket url e.g. "192.0.2.2:4242" or "wss://192.0.2.2:4242"`)
			assert.Contains(t, err.Error(), `If RawServerURL is specified, ServerPubKey is also required and must be a 32-byte hex string`)
		})

		t.Run("with unnecessary values", func(t *testing.T) {
			rawToml := `
				ServerURL = "example.com:80"
				ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
				LinkFeedID = "0x00026b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472"
			`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			err = ValidatePluginConfig(mc, v1FeedId)
			assert.Contains(t, err.Error(), `linkFeedID may not be specified for v1 jobs`)
		})
	})

	t.Run("Mercury v2/v3", func(t *testing.T) {
		t.Run("with valid values", func(t *testing.T) {
			rawToml := `
				ServerURL = "example.com:80"
				ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
				LinkFeedID = "0x00026b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472"
				NativeFeedID = "0x00036b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472"
			`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			err = ValidatePluginConfig(mc, v2FeedId)
			require.NoError(t, err)

			require.NotNil(t, mc.LinkFeedID)
			require.NotNil(t, mc.NativeFeedID)
			assert.Equal(t, "0x00026b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472", (*mc.LinkFeedID).String())
			assert.Equal(t, "0x00036b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472", (*mc.NativeFeedID).String())
		})

		t.Run("with invalid values", func(t *testing.T) {
			var mc PluginConfig

			rawToml := `LinkFeedID = "test"`
			err := toml.Unmarshal([]byte(rawToml), &mc)
			assert.Contains(t, err.Error(), "toml: hash: expected a hex string starting with '0x'")

			rawToml = `LinkFeedID = "0xtest"`
			err = toml.Unmarshal([]byte(rawToml), &mc)
			assert.Contains(t, err.Error(), `toml: hash: UnmarshalText failed: encoding/hex: invalid byte: U+0074 't'`)

			rawToml = `
				ServerURL = "example.com:80"
				ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
				LinkFeedID = "0x00026b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472"
			`
			err = toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			err = ValidatePluginConfig(mc, v2FeedId)
			assert.Contains(t, err.Error(), "nativeFeedID must be specified for v2 jobs")
		})

		t.Run("with unnecessary values", func(t *testing.T) {
			rawToml := `
				ServerURL = "example.com:80"
				ServerPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
				InitialBlockNumber = 1234
			`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			err = ValidatePluginConfig(mc, v2FeedId)
			assert.Contains(t, err.Error(), `initialBlockNumber may not be specified for v2 jobs`)
		})
	})
}

func Test_PluginConfig_GetServers(t *testing.T) {
	t.Run("with single server", func(t *testing.T) {
		pubKey := utils.PlainHexBytes([]byte{1, 2, 3})
		pc := PluginConfig{RawServerURL: "example.com", ServerPubKey: pubKey}
		require.Len(t, pc.GetServers(), 1)
		assert.Equal(t, "example.com", pc.GetServers()[0].URL)
		assert.Equal(t, pubKey, pc.GetServers()[0].PubKey)

		pc = PluginConfig{RawServerURL: "wss://example.com", ServerPubKey: pubKey}
		require.Len(t, pc.GetServers(), 1)
		assert.Equal(t, "example.com", pc.GetServers()[0].URL)
		assert.Equal(t, pubKey, pc.GetServers()[0].PubKey)

		pc = PluginConfig{RawServerURL: "example.com:1234/foo", ServerPubKey: pubKey}
		require.Len(t, pc.GetServers(), 1)
		assert.Equal(t, "example.com:1234/foo", pc.GetServers()[0].URL)
		assert.Equal(t, pubKey, pc.GetServers()[0].PubKey)

		pc = PluginConfig{RawServerURL: "wss://example.com:1234/foo", ServerPubKey: pubKey}
		require.Len(t, pc.GetServers(), 1)
		assert.Equal(t, "example.com:1234/foo", pc.GetServers()[0].URL)
		assert.Equal(t, pubKey, pc.GetServers()[0].PubKey)
	})

	t.Run("with multiple servers", func(t *testing.T) {
		servers := map[string]utils.PlainHexBytes{
			"example.com:80":                 utils.PlainHexBytes([]byte{1, 2, 3}),
			"mercuryserver.invalid:1234/foo": utils.PlainHexBytes([]byte{4, 5, 6}),
		}
		pc := PluginConfig{Servers: servers}

		require.Len(t, pc.GetServers(), 2)
		assert.Equal(t, "example.com:80", pc.GetServers()[0].URL)
		assert.Equal(t, utils.PlainHexBytes{1, 2, 3}, pc.GetServers()[0].PubKey)
		assert.Equal(t, "mercuryserver.invalid:1234/foo", pc.GetServers()[1].URL)
		assert.Equal(t, utils.PlainHexBytes{4, 5, 6}, pc.GetServers()[1].PubKey)
	})
}
