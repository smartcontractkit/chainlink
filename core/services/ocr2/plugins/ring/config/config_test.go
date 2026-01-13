package config

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginConfig_Unmarshal(t *testing.T) {
	t.Run("unmarshals valid JSON", func(t *testing.T) {
		raw := `{"shardConfigAddr": "0x1234567890abcdef"}`
		var cfg PluginConfig
		err := cfg.Unmarshal([]byte(raw))
		require.NoError(t, err)
		assert.Equal(t, "0x1234567890abcdef", cfg.ShardConfigAddr)
	})

	t.Run("unmarshals empty config", func(t *testing.T) {
		raw := `{}`
		var cfg PluginConfig
		err := cfg.Unmarshal([]byte(raw))
		require.NoError(t, err)
		assert.Empty(t, cfg.ShardConfigAddr)
	})

	t.Run("fails on invalid JSON", func(t *testing.T) {
		raw := `{invalid}`
		var cfg PluginConfig
		err := cfg.Unmarshal([]byte(raw))
		require.Error(t, err)
	})
}

func TestPluginConfig_UnmarshalTOML(t *testing.T) {
	t.Run("unmarshals from TOML", func(t *testing.T) {
		rawToml := `shardConfigAddr = "0xdeadbeef"`
		var cfg PluginConfig
		err := toml.Unmarshal([]byte(rawToml), &cfg)
		require.NoError(t, err)
		assert.Equal(t, "0xdeadbeef", cfg.ShardConfigAddr)
	})
}

func TestPluginConfig_Validate(t *testing.T) {
	t.Run("valid address passes", func(t *testing.T) {
		cfg := PluginConfig{
			ShardConfigAddr: "0x1234567890123456789012345678901234567890",
		}
		err := cfg.Validate()
		require.NoError(t, err)
	})

	t.Run("empty address fails", func(t *testing.T) {
		cfg := PluginConfig{
			ShardConfigAddr: "",
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "shardConfigAddr is required")
	})

	t.Run("invalid address fails", func(t *testing.T) {
		cfg := PluginConfig{
			ShardConfigAddr: "not-an-address",
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a valid hex address")
	})

	t.Run("short address fails", func(t *testing.T) {
		cfg := PluginConfig{
			ShardConfigAddr: "0x1234",
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a valid hex address")
	})
}
