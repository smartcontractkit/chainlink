package config

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			assert.Contains(t, err.Error(), `mercury: ServerPubKey is required and must be a 32-byte hex string`)
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
