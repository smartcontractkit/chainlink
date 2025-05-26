package config

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
				Servers = { "example.com:80" = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", "example2.invalid:1234" = "524ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93" }
				BenchmarkMode = true
				ChannelDefinitionsContractAddress = "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
				ChannelDefinitionsContractFromBlock = 1234
				ChannelDefinitions = """
%s
"""`, cdjson)

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			assert.Len(t, mc.Servers, 2)
			assert.Equal(t, map[string]utils.PlainHexBytes{"example.com:80": utils.PlainHexBytes{0x72, 0x4f, 0xf6, 0xea, 0xe9, 0xe9, 0x0, 0x27, 0xe, 0xdf, 0xff, 0x23, 0x3e, 0x16, 0x32, 0x2a, 0x70, 0xec, 0x6, 0xe1, 0xa6, 0xe6, 0x2a, 0x81, 0xef, 0x13, 0x92, 0x1f, 0x39, 0x8f, 0x6c, 0x93}, "example2.invalid:1234": utils.PlainHexBytes{0x52, 0x4f, 0xf6, 0xea, 0xe9, 0xe9, 0x0, 0x27, 0xe, 0xdf, 0xff, 0x23, 0x3e, 0x16, 0x32, 0x2a, 0x70, 0xec, 0x6, 0xe1, 0xa6, 0xe6, 0x2a, 0x81, 0xef, 0x13, 0x92, 0x1f, 0x39, 0x8f, 0x6c, 0x93}}, mc.Servers)
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
				Servers = { "example.com:80" = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93" }
				DonID = 12345
				ChannelDefinitions = """
%s
"""`, cdjson)

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			assert.Len(t, mc.Servers, 1)
			assert.JSONEq(t, cdjson, mc.ChannelDefinitions)
			assert.Equal(t, uint32(12345), mc.DonID)
			assert.False(t, mc.BenchmarkMode)

			err = mc.Validate()
			require.NoError(t, err)
		})
		t.Run("with only channelDefinitions contract details", func(t *testing.T) {
			rawToml := `
			Servers = { "example.com:80" = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93" }
			DonID = 12345
			ChannelDefinitionsContractAddress = "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			assert.Len(t, mc.Servers, 1)
			assert.Equal(t, "0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF", mc.ChannelDefinitionsContractAddress.Hex())
			assert.Equal(t, uint32(12345), mc.DonID)
			assert.False(t, mc.BenchmarkMode)

			err = mc.Validate()
			require.NoError(t, err)
		})
		t.Run("with missing ChannelDefinitionsContractAddress", func(t *testing.T) {
			rawToml := `
			DonID = 12345
			Servers = { "example.com:80" = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93" }
			`

			var mc PluginConfig
			err := toml.Unmarshal([]byte(rawToml), &mc)
			require.NoError(t, err)

			assert.Len(t, mc.Servers, 1)
			assert.Equal(t, uint32(12345), mc.DonID)
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
			assert.Contains(t, err.Error(), `DonID must be specified and not zero`)
			assert.Contains(t, err.Error(), `At least one Mercury server or Transmitter must be specified`)
			assert.Contains(t, err.Error(), `ChannelDefinitionsContractAddress is required if ChannelDefinitions is not specified`)
		})
	})
}

func Test_PluginConfig_Validate(t *testing.T) {
	t.Run("with invalid URLs or keys", func(t *testing.T) {
		servers := map[string]utils.PlainHexBytes{
			"not a valid url":                utils.PlainHexBytes([]byte{1, 2, 3}),
			"mercuryserver.invalid:1234/foo": nil,
		}
		pc := PluginConfig{Servers: servers}

		err := pc.Validate()
		assert.Contains(t, err.Error(), "ServerPubKey must be a 32-byte hex string")
		assert.Contains(t, err.Error(), "invalid value for ServerURL: llo: invalid value for ServerURL, got: \"not a valid url\"")
	})
}

func Test_PluginConfig_GetServers(t *testing.T) {
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
