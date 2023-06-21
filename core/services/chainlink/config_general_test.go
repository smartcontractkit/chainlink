//go:build !dev

package chainlink

import (
	_ "embed"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type TestSecrets struct {
	TestConfigSecrets
}

type TestConfigSecrets struct {
	Database   TestDatabaseSecrets          `toml:",omitempty"`
	Explorer   TestExplorerSecrets          `toml:",omitempty"`
	Password   TestPasswordsSecrets         `toml:",omitempty"`
	Pyroscope  TestPyroscopeSecrets         `toml:",omitempty"`
	Prometheus TestPrometheusSecrets        `toml:",omitempty"`
	Mercury    TestMercurySecrets           `toml:",omitempty"`
	Threshold  TestThresholdKeyShareSecrets `toml:",omitempty"`
}

type TestDatabaseSecrets struct {
	URL                  models.URL
	BackupURL            models.URL
	AllowSimplePasswords *bool
}

type TestExplorerSecrets struct {
	AccessKey string
	Secret    string
}

type TestPasswordsSecrets struct {
	Keystore string
	VRF      string
}

type TestPyroscopeSecrets struct {
	AuthToken string
}

type TestPrometheusSecrets struct {
	AuthToken string
}

type TestMercurySecrets struct {
	Credentials map[string]TestMercuryCredentials
}

type TestMercuryCredentials struct {
	URL      models.URL
	Username string
	Password string
}

type TestThresholdKeyShareSecrets struct {
	ThresholdKeyShare string
}

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

func TestConfig_SecretsMerging(t *testing.T) {
	setInFile := "set in config file"
	allow := false
	dbURL := "postgres://172.17.0.1:5432/primary"
	backupDbURL := "postgres://172.17.0.1:5433/replica"
	testConfigFileContents := Config{
		Core: toml.Core{
			RootDir: &setInFile,
			P2P: toml.P2P{
				V2: toml.P2PV2{
					AnnounceAddresses: &[]string{setInFile},
					ListenAddresses:   &[]string{setInFile},
				},
			},
		},
	}
	testSecretsFileContentsComplete := TestSecrets{
		TestConfigSecrets: TestConfigSecrets{
			Database: TestDatabaseSecrets{
				URL:                  *models.MustParseURL(dbURL),
				BackupURL:            *models.MustParseURL(backupDbURL),
				AllowSimplePasswords: &allow,
			},
			Explorer: TestExplorerSecrets{
				AccessKey: "EXPLORER_ACCESS_KEY",
				Secret:    "EXPLORER_TOKEN",
			},
			Password: TestPasswordsSecrets{
				Keystore: "mysecretpassword",
				VRF:      "mysecretvrfpassword",
			},
			Pyroscope: TestPyroscopeSecrets{
				AuthToken: "PYROSCOPE_TOKEN",
			},
			Prometheus: TestPrometheusSecrets{
				AuthToken: "PROM_TOKEN",
			},
			Mercury: TestMercurySecrets{
				Credentials: map[string]TestMercuryCredentials{
					"key1": {
						URL:      *models.MustParseURL("https://mercury.stage.link"),
						Username: "user",
						Password: "user_pass",
					},
					"key2": {
						URL:      *models.MustParseURL("https://mercury.stage.link"),
						Username: "user",
						Password: "user_pass",
					},
				},
			},
			Threshold: TestThresholdKeyShareSecrets{
				ThresholdKeyShare: "THRESHOLD_SECRET",
			},
		},
	}

	additionalMercurySecrets := TestMercurySecrets{
		Credentials: map[string]TestMercuryCredentials{
			"key3": {
				URL:      *models.MustParseURL("https://mercury.stage.link"),
				Username: "user",
				Password: "user_pass",
			},
			"key4": {
				URL:      *models.MustParseURL("https://mercury.stage.link"),
				Username: "user",
				Password: "user_pass",
			},
		},
	}
	t.Run("verify secrets merging", func(t *testing.T) {
		opts := new(GeneralConfigOpts)
		configFiles := []string{utils.MakeTestFile(t, testConfigFileContents, "test.toml")}
		secretsFiles := []string{
			utils.MakeTestFile(t, TestSecrets{TestConfigSecrets: TestConfigSecrets{Database: testSecretsFileContentsComplete.Database}}, "test_secrets1.toml"),
			utils.MakeTestFile(t, TestSecrets{TestConfigSecrets: TestConfigSecrets{Explorer: testSecretsFileContentsComplete.Explorer}}, "test_secrets2.toml"),
			utils.MakeTestFile(t, TestSecrets{TestConfigSecrets: TestConfigSecrets{Password: testSecretsFileContentsComplete.Password}}, "test_secrets3.toml"),
			utils.MakeTestFile(t, TestSecrets{TestConfigSecrets: TestConfigSecrets{Pyroscope: testSecretsFileContentsComplete.Pyroscope}}, "test_secrets4.toml"),
			utils.MakeTestFile(t, TestSecrets{TestConfigSecrets: TestConfigSecrets{Prometheus: testSecretsFileContentsComplete.Prometheus}}, "test_secrets5.toml"),
			utils.MakeTestFile(t, TestSecrets{TestConfigSecrets: TestConfigSecrets{Mercury: testSecretsFileContentsComplete.Mercury}}, "test_secrets6.toml"),
			utils.MakeTestFile(t, TestSecrets{TestConfigSecrets: TestConfigSecrets{Mercury: additionalMercurySecrets}}, "test_secrets6a.toml"),
			utils.MakeTestFile(t, TestSecrets{TestConfigSecrets: TestConfigSecrets{Threshold: testSecretsFileContentsComplete.Threshold}}, "test_secrets7.toml"),
		}
		err := opts.Setup(configFiles, secretsFiles)
		require.NoErrorf(t, err, "error: %s", err)

		err = opts.parse()
		require.NoErrorf(t, err, "error: %s", err)

		require.NoErrorf(t, err, "error testing: %s, %s", configFiles, secretsFiles)

		assert.Equal(t, dbURL, opts.Secrets.Database.URL.URL().String())
		assert.Equal(t, backupDbURL, opts.Secrets.Database.BackupURL.URL().String())
		assert.Equal(t, testSecretsFileContentsComplete.Explorer.AccessKey, opts.Secrets.Explorer.AccessKey.XXXTestingOnlyString())
		assert.Equal(t, testSecretsFileContentsComplete.Explorer.Secret, opts.Secrets.Explorer.Secret.XXXTestingOnlyString())
		assert.Equal(t, testSecretsFileContentsComplete.Password.Keystore, opts.Secrets.Password.Keystore.XXXTestingOnlyString())
		assert.Equal(t, testSecretsFileContentsComplete.Password.VRF, opts.Secrets.Password.VRF.XXXTestingOnlyString())
		assert.Equal(t, testSecretsFileContentsComplete.Pyroscope.AuthToken, opts.Secrets.Pyroscope.AuthToken.XXXTestingOnlyString())
		assert.Equal(t, testSecretsFileContentsComplete.Prometheus.AuthToken, opts.Secrets.Prometheus.AuthToken.XXXTestingOnlyString())
		assert.Equal(t, *merge(testSecretsFileContentsComplete.Mercury, additionalMercurySecrets), convertMercurySecrets(opts.Secrets.Mercury))
		assert.Equal(t, testSecretsFileContentsComplete.Threshold.ThresholdKeyShare, opts.Secrets.Threshold.ThresholdKeyShare.XXXTestingOnlyString())
	})
}

func convertMercurySecrets(mercurySecrets toml.MercurySecrets) TestMercurySecrets {
	testSecrets := TestMercurySecrets{
		Credentials: make(map[string]TestMercuryCredentials),
	}

	for key, credentials := range mercurySecrets.Credentials {
		testCredentials := TestMercuryCredentials{
			URL:      (models.URL)(*credentials.URL),
			Username: (string)(*credentials.Username),
			Password: (string)(*credentials.Password),
		}

		testSecrets.Credentials[key] = testCredentials
	}

	return testSecrets
}

func merge(map1 TestMercurySecrets, map2 TestMercurySecrets) *TestMercurySecrets {
	combinedMap := make(map[string]TestMercuryCredentials)

	for key, value := range map1.Credentials {
		combinedMap[key] = value
	}

	for key, value := range map2.Credentials {
		combinedMap[key] = value
	}

	return &TestMercurySecrets{Credentials: combinedMap}
}
