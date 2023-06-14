//go:build !dev

package chainlink

import (
	_ "embed"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
)

func TestTOMLGeneralConfig_Defaults(t *testing.T) {
	config, err := GeneralConfigOpts{}.New()
	require.NoError(t, err)
	assert.Equal(t, (*url.URL)(nil), config.WebServer().BridgeResponseURL())
	assert.Nil(t, config.DefaultChainID())
	assert.False(t, config.EVMRPCEnabled())
	assert.False(t, config.EVMEnabled())
	assert.False(t, config.CosmosEnabled())
	assert.False(t, config.SolanaEnabled())
	assert.False(t, config.StarkNetEnabled())
	assert.Equal(t, false, config.JobPipeline().ExternalInitiatorsEnabled())
	assert.Equal(t, 15*time.Minute, config.WebServer().SessionTimeout().Duration())
}

func TestTOMLGeneralConfig_InsecureConfig(t *testing.T) {
	t.Parallel()

	t.Run("all insecure configs are false by default", func(t *testing.T) {
		config, err := GeneralConfigOpts{}.New()
		require.NoError(t, err)

		assert.False(t, config.Insecure().DevWebServer())
		assert.False(t, config.Insecure().DisableRateLimiting())
		assert.False(t, config.Insecure().InfiniteDepthQueries())
		assert.False(t, config.Insecure().OCRDevelopmentMode())
	})

	t.Run("insecure config ignore override on non-dev builds", func(t *testing.T) {
		config, err := GeneralConfigOpts{
			OverrideFn: func(c *Config, s *Secrets) {
				*c.Insecure.DevWebServer = true
				*c.Insecure.DisableRateLimiting = true
				*c.Insecure.InfiniteDepthQueries = true
				*c.AuditLogger.Enabled = true
			}}.New()
		require.NoError(t, err)

		// Just asserting that override logic work on a safe config
		assert.True(t, config.AuditLogger().Enabled())

		assert.False(t, config.Insecure().DevWebServer())
		assert.False(t, config.Insecure().DisableRateLimiting())
		assert.False(t, config.Insecure().InfiniteDepthQueries())
	})

	t.Run("ValidateConfig fails if insecure config is set on non-dev builds", func(t *testing.T) {
		config := `
		  [insecure]
		  DevWebServer = true
		  DisableRateLimiting = false
		  InfiniteDepthQueries = false
		  OCRDevelopmentMode = false
		`
		opts := GeneralConfigOpts{
			ConfigStrings: []string{config},
		}
		cfg, err := opts.New()
		require.NoError(t, err)
		err = cfg.Validate()
		require.Contains(t, err.Error(), "invalid configuration: Insecure.DevWebServer: invalid value (true): insecure configs are not allowed on secure builds")
	})
}

func TestValidateDB(t *testing.T) {
	t.Setenv(string(v2.EnvConfig), "")

	t.Run("unset db url", func(t *testing.T) {
		t.Setenv(string(v2.EnvDatabaseURL), "")

		config, err := GeneralConfigOpts{}.New()
		require.NoError(t, err)
		err = config.ValidateDB()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidSecrets)
	})

	t.Run("garbage db url", func(t *testing.T) {
		t.Setenv(string(v2.EnvDatabaseURL), "garbage")

		config, err := GeneralConfigOpts{}.New()
		require.NoError(t, err)
		err = config.ValidateDB()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidSecrets)
	})

	t.Run("dev url", func(t *testing.T) {
		t.Setenv(string(v2.EnvDatabaseURL), "postgres://postgres:admin@localhost:5432/chainlink_dev_test?sslmode=disable")

		config, err := GeneralConfigOpts{}.New()
		require.NoError(t, err)
		err = config.ValidateDB()
		require.NoError(t, err)
	})

	t.Run("bad password url", func(t *testing.T) {
		t.Setenv(string(v2.EnvDatabaseURL), "postgres://postgres:pwdToShort@localhost:5432/chainlink_dev_prod?sslmode=disable")

		config, err := GeneralConfigOpts{}.New()
		require.NoError(t, err)
		err = config.ValidateDB()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidSecrets)
	})

}

func TestConfig_LogSQL(t *testing.T) {
	config, err := GeneralConfigOpts{}.New()
	require.NoError(t, err)

	config.SetLogSQL(true)
	assert.Equal(t, config.Database().LogSQL(), true)

	config.SetLogSQL(false)
	assert.Equal(t, config.Database().LogSQL(), false)
}
