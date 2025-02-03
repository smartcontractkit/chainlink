//go:build !dev

package chainlink

import (
	_ "embed"
	"fmt"
	"maps"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

func TestTOMLGeneralConfig_Defaults(t *testing.T) {
	config, err := GeneralConfigOpts{}.New()
	require.NoError(t, err)
	assert.Equal(t, (*url.URL)(nil), config.WebServer().BridgeResponseURL())
	assert.False(t, config.EVMConfigs().RPCEnabled())
	assert.False(t, config.EVMEnabled())
	assert.False(t, config.CosmosEnabled())
	assert.False(t, config.SolanaEnabled())
	assert.False(t, config.StarkNetEnabled())
	assert.False(t, config.TronEnabled())
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
	t.Setenv(string(env.Config), "")

	t.Run("unset db url", func(t *testing.T) {
		t.Setenv(string(env.DatabaseURL), "")

		config, err := GeneralConfigOpts{}.New()
		require.NoError(t, err)

		err = config.ValidateDB()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidSecrets)
	})

	t.Run("dev url", func(t *testing.T) {
		t.Setenv(string(env.DatabaseURL), "postgres://postgres:admin@localhost:5432/chainlink_dev_test?sslmode=disable")

		config, err := GeneralConfigOpts{}.New()
		require.NoError(t, err)
		err = config.ValidateDB()
		require.NoError(t, err)
	})

	t.Run("bad password url", func(t *testing.T) {
		t.Setenv(string(env.DatabaseURL), "postgres://postgres:pwdTooShort@localhost:5432/chainlink_dev_prod?sslmode=disable")
		t.Setenv(string(env.DatabaseAllowSimplePasswords), "false")

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

//go:embed testdata/mergingsecretsdata/secrets-database.toml
var databaseSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-password.toml
var passwordSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-pyroscope.toml
var pyroscopeSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-prometheus.toml
var prometheusSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-mercury-split-one.toml
var mercurySecretsTOMLSplitOne string

//go:embed testdata/mergingsecretsdata/secrets-mercury-split-two.toml
var mercurySecretsTOMLSplitTwo string

//go:embed testdata/mergingsecretsdata/secrets-threshold.toml
var thresholdSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-webserver-ldap.toml
var WebServerLDAPSecretsTOML string

func TestConfig_SecretsMerging(t *testing.T) {
	t.Run("verify secrets merging in GeneralConfigOpts.New()", func(t *testing.T) {
		databaseSecrets, err := parseSecrets(databaseSecretsTOML)
		require.NoErrorf(t, err, "error: %s", err)
		passwordSecrets, err2 := parseSecrets(passwordSecretsTOML)
		require.NoErrorf(t, err2, "error: %s", err2)
		pyroscopeSecrets, err3 := parseSecrets(pyroscopeSecretsTOML)
		require.NoErrorf(t, err3, "error: %s", err3)
		prometheusSecrets, err4 := parseSecrets(prometheusSecretsTOML)
		require.NoErrorf(t, err4, "error: %s", err4)
		mercurySecrets_a, err5 := parseSecrets(mercurySecretsTOMLSplitOne)
		require.NoErrorf(t, err5, "error: %s", err5)
		mercurySecrets_b, err6 := parseSecrets(mercurySecretsTOMLSplitTwo)
		require.NoErrorf(t, err6, "error: %s", err6)
		thresholdSecrets, err7 := parseSecrets(thresholdSecretsTOML)
		require.NoErrorf(t, err7, "error: %s", err7)
		webserverLDAPSecrets, err8 := parseSecrets(WebServerLDAPSecretsTOML)
		require.NoErrorf(t, err8, "error: %s", err8)

		opts := new(GeneralConfigOpts)
		configFiles := []string{
			"testdata/mergingsecretsdata/config.toml",
		}
		secretsFiles := []string{
			"testdata/mergingsecretsdata/secrets-database.toml",
			"testdata/mergingsecretsdata/secrets-password.toml",
			"testdata/mergingsecretsdata/secrets-pyroscope.toml",
			"testdata/mergingsecretsdata/secrets-prometheus.toml",
			"testdata/mergingsecretsdata/secrets-mercury-split-one.toml",
			"testdata/mergingsecretsdata/secrets-mercury-split-two.toml",
			"testdata/mergingsecretsdata/secrets-threshold.toml",
			"testdata/mergingsecretsdata/secrets-webserver-ldap.toml",
		}
		err = opts.Setup(configFiles, secretsFiles)
		require.NoErrorf(t, err, "error: %s", err)

		err = opts.parse()
		require.NoErrorf(t, err, "error testing: %s, %s", configFiles, secretsFiles)

		assert.Equal(t, databaseSecrets.Database.URL.URL().String(), opts.Secrets.Database.URL.URL().String())
		assert.Equal(t, databaseSecrets.Database.BackupURL.URL().String(), opts.Secrets.Database.BackupURL.URL().String())

		assert.Equal(t, (string)(*passwordSecrets.Password.Keystore), (string)(*opts.Secrets.Password.Keystore))
		assert.Equal(t, (string)(*passwordSecrets.Password.VRF), (string)(*opts.Secrets.Password.VRF))
		assert.Equal(t, (string)(*pyroscopeSecrets.Pyroscope.AuthToken), (string)(*opts.Secrets.Pyroscope.AuthToken))
		assert.Equal(t, (string)(*prometheusSecrets.Prometheus.AuthToken), (string)(*opts.Secrets.Prometheus.AuthToken))
		assert.Equal(t, (string)(*thresholdSecrets.Threshold.ThresholdKeyShare), (string)(*opts.Secrets.Threshold.ThresholdKeyShare))

		assert.Equal(t, webserverLDAPSecrets.WebServer.LDAP.ServerAddress.URL().String(), opts.Secrets.WebServer.LDAP.ServerAddress.URL().String())
		assert.Equal(t, webserverLDAPSecrets.WebServer.LDAP.ReadOnlyUserLogin, opts.Secrets.WebServer.LDAP.ReadOnlyUserLogin)
		assert.Equal(t, webserverLDAPSecrets.WebServer.LDAP.ReadOnlyUserPass, opts.Secrets.WebServer.LDAP.ReadOnlyUserPass)

		err = assertDeepEqualityMercurySecrets(*merge(mercurySecrets_a.Mercury, mercurySecrets_b.Mercury), opts.Secrets.Mercury)
		require.NoErrorf(t, err, "merged mercury secrets unequal")
	})
}

func parseSecrets(secrets string) (*Secrets, error) {
	var s Secrets
	if err := config.DecodeTOML(strings.NewReader(secrets), &s); err != nil {
		return nil, fmt.Errorf("failed to decode secrets TOML: %w", err)
	}

	return &s, nil
}

func assertDeepEqualityMercurySecrets(expected toml.MercurySecrets, actual toml.MercurySecrets) error {
	if len(expected.Credentials) != len(actual.Credentials) {
		return fmt.Errorf("maps are not equal in length: len(expected): %d, len(actual): %d", len(expected.Credentials), len(actual.Credentials))
	}

	for key, value := range expected.Credentials {
		equal := true
		actualValue := actual.Credentials[key]
		if (string)(*value.Username) != (string)(*actualValue.Username) {
			equal = false
		}
		if (string)(*value.Password) != (string)(*actualValue.Password) {
			equal = false
		}
		if value.URL.URL().String() != actualValue.URL.URL().String() {
			equal = false
		}
		if !equal {
			return fmt.Errorf("maps are not equal: expected[%s] = {%s, %s, %s}, actual[%s] = {%s, %s, %s}",
				key, (string)(*value.Username), (string)(*value.Password), value.URL.URL().String(),
				key, (string)(*actualValue.Username), (string)(*actualValue.Password), actualValue.URL.URL().String())
		}
	}
	return nil
}

func merge(map1 toml.MercurySecrets, map2 toml.MercurySecrets) *toml.MercurySecrets {
	combinedMap := make(map[string]toml.MercuryCredentials)
	maps.Copy(combinedMap, map1.Credentials)
	maps.Copy(combinedMap, map2.Credentials)
	return &toml.MercurySecrets{Credentials: combinedMap}
}
