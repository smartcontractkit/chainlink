package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gotoml "github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var (
	setInFile = "set in config file"
	setInEnv  = "set in env"

	testEnvContents = fmt.Sprintf("P2P.V2.AnnounceAddresses = ['%s']", setInEnv)

	testConfigFileContents = chainlink.Config{
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

	testSecretsFileContents = chainlink.Secrets{
		Secrets: toml.Secrets{
			Prometheus: toml.PrometheusSecrets{
				AuthToken: models.NewSecret("PROM_TOKEN"),
			},
		},
	}

	testSecretsRedactedContents = chainlink.Secrets{
		Secrets: toml.Secrets{
			Prometheus: toml.PrometheusSecrets{
				AuthToken: models.NewSecret("xxxxx"),
			},
		},
	}

	allow                           = false
	dbURL                           = "postgres://chainlink:mysecretpassword@172.17.0.1:5432/primary"
	backupDbURL                     = "postgres://chainlink:mysecretpassword@172.17.0.1:5433/replica"
	testSecretsFileContentsComplete = chainlink.Secrets{
		Secrets: v2.Secrets{
			Database: v2.DatabaseSecrets{
				URL:                  models.NewSecretURL(models.MustParseURL(dbURL)),
				BackupURL:            models.NewSecretURL(models.MustParseURL(backupDbURL)),
				AllowSimplePasswords: &allow,
			},
			Explorer: v2.ExplorerSecrets{
				AccessKey: models.NewSecret("EXPLORER_ACCESS_KEY"),
				Secret:    models.NewSecret("EXPLORER_TOKEN"),
			},
			Password: v2.Passwords{
				Keystore: models.NewSecret("mysecretpassword"),
				VRF:      models.NewSecret("mysecretvrfpassword"),
			},
			Pyroscope: v2.PyroscopeSecrets{
				AuthToken: models.NewSecret("PYROSCOPE_TOKEN"),
			},
			Prometheus: v2.PrometheusSecrets{
				AuthToken: models.NewSecret("PROM_TOKEN"),
			},
			Mercury: v2.MercurySecrets{
				Credentials: map[string]v2.MercuryCredentials{
					"key1": {
						URL:      models.NewSecretURL(models.MustParseURL("https://mercury.stage.link")),
						Username: models.NewSecret("user"),
						Password: models.NewSecret("user_pass"),
					},
					"key2": {
						URL:      models.NewSecretURL(models.MustParseURL("https://mercury.stage.link")),
						Username: models.NewSecret("user"),
						Password: models.NewSecret("user_pass"),
					},
				},
			},
			Threshold: v2.ThresholdKeyShareSecrets{
				ThresholdKeyShare: models.NewSecret("THRESHOLD_SECRET"),
			},
		},
	}

	testSecretsRedactedContentsComplete = chainlink.Secrets{
		Secrets: v2.Secrets{
			Database: v2.DatabaseSecrets{
				URL:                  models.NewSecretURL(models.MustParseURL("xxxxx")),
				BackupURL:            models.NewSecretURL(models.MustParseURL("xxxxx")),
				AllowSimplePasswords: &allow,
			},
			Explorer: v2.ExplorerSecrets{
				AccessKey: models.NewSecret("xxxxx"),
				Secret:    models.NewSecret("xxxxx"),
			},
			Password: v2.Passwords{
				Keystore: models.NewSecret("xxxxx"),
				VRF:      models.NewSecret("xxxxx"),
			},
			Pyroscope: v2.PyroscopeSecrets{
				AuthToken: models.NewSecret("xxxxx"),
			},
			Prometheus: v2.PrometheusSecrets{
				AuthToken: models.NewSecret("xxxxx"),
			},
			Mercury: v2.MercurySecrets{
				Credentials: map[string]v2.MercuryCredentials{
					"key1": {
						URL:      models.NewSecretURL(models.MustParseURL("xxxxx")),
						Username: models.NewSecret("xxxxx"),
						Password: models.NewSecret("xxxxx"),
					},
					"key2": {
						URL:      models.NewSecretURL(models.MustParseURL("xxxxx")),
						Username: models.NewSecret("xxxxx"),
						Password: models.NewSecret("xxxxx"),
					},
					"key3": {
						URL:      models.NewSecretURL(models.MustParseURL("xxxxx")),
						Username: models.NewSecret("xxxxx"),
						Password: models.NewSecret("xxxxx"),
					},
					"key4": {
						URL:      models.NewSecretURL(models.MustParseURL("xxxxx")),
						Username: models.NewSecret("xxxxx"),
						Password: models.NewSecret("xxxxx"),
					},
				},
			},
			Threshold: v2.ThresholdKeyShareSecrets{
				ThresholdKeyShare: models.NewSecret("xxxxx"),
			},
		},
	}

	additionalMercurySecrets = v2.MercurySecrets{
		Credentials: map[string]v2.MercuryCredentials{
			"key3": {
				URL:      models.NewSecretURL(models.MustParseURL("https://mercury.stage.link")),
				Username: models.NewSecret("user"),
				Password: models.NewSecret("user_pass"),
			},
			"key4": {
				URL:      models.NewSecretURL(models.MustParseURL("https://mercury.stage.link")),
				Username: models.NewSecret("user"),
				Password: models.NewSecret("user_pass"),
			},
		},
	}
)

func makeTestFile(t *testing.T, contents any, fileName string) string {
	d := t.TempDir()
	p := filepath.Join(d, fileName)

	b, err := gotoml.Marshal(contents)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(p, b, 0777))
	return p
}

func withDefaults(t *testing.T, c chainlink.Config, s chainlink.Secrets) chainlink.GeneralConfig {
	cfg, err := chainlink.GeneralConfigOpts{Config: c, Secrets: s}.New()
	require.NoError(t, err)
	return cfg
}

func Test_initServerConfig(t *testing.T) {
	type args struct {
		opts         *chainlink.GeneralConfigOpts
		fileNames    []string
		secretsFiles []string
		envVar       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantCfg chainlink.GeneralConfig
	}{
		{
			name: "env only",
			args: args{
				opts:   new(chainlink.GeneralConfigOpts),
				envVar: testEnvContents,
			},
			wantCfg: withDefaults(t, chainlink.Config{
				Core: toml.Core{
					P2P: toml.P2P{
						V2: toml.P2PV2{
							AnnounceAddresses: &[]string{setInEnv},
						},
					},
				},
			}, chainlink.Secrets{}),
		},
		{
			name: "files only",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
			},
			wantCfg: withDefaults(t, testConfigFileContents, chainlink.Secrets{}),
		},
		{
			name: "file error",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{"notexist"},
			},
			wantErr: true,
		},
		{
			name: "env overlay of file",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				envVar:    testEnvContents,
			},
			wantCfg: withDefaults(t, chainlink.Config{
				Core: toml.Core{
					RootDir: &setInFile,
					P2P: toml.P2P{
						V2: toml.P2PV2{
							// env should override this specific field
							AnnounceAddresses: &[]string{setInEnv},
							ListenAddresses:   &[]string{setInFile},
						},
					},
				},
			}, chainlink.Secrets{}),
		},
		{
			name: "failed to read secrets",
			args: args{
				opts:         new(chainlink.GeneralConfigOpts),
				fileNames:    []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{"/doesnt-exist"},
			},
			wantErr: true,
		},
		{
			name: "reading secrets",
			args: args{
				opts:         new(chainlink.GeneralConfigOpts),
				fileNames:    []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{makeTestFile(t, testSecretsFileContents, "test_secrets.toml")},
			},
			wantCfg: withDefaults(t, testConfigFileContents, testSecretsRedactedContents),
		},
		{
			name: "reading multiple secrets",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Database: testSecretsFileContentsComplete.Database}}, "test_secrets1.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Explorer: testSecretsFileContentsComplete.Explorer}}, "test_secrets2.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Password: testSecretsFileContentsComplete.Password}}, "test_secrets3.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Pyroscope: testSecretsFileContentsComplete.Pyroscope}}, "test_secrets4.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Prometheus: testSecretsFileContentsComplete.Prometheus}}, "test_secrets5.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Mercury: testSecretsFileContentsComplete.Mercury}}, "test_secrets6.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Mercury: additionalMercurySecrets}}, "test_secrets6a.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Threshold: testSecretsFileContentsComplete.Threshold}}, "test_secrets7.toml"),
				},
			},
			wantCfg: withDefaults(t, testConfigFileContents, testSecretsRedactedContentsComplete),
		},
		{
			name: "reading multiple secrets with overrides: Database",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Database: testSecretsFileContentsComplete.Database}}, "test_secrets1.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Database: testSecretsFileContentsComplete.Database}}, "test_secrets1a.toml"),
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Explorer",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Explorer: testSecretsFileContentsComplete.Explorer}}, "test_secrets1.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Explorer: testSecretsFileContentsComplete.Explorer}}, "test_secrets1a.toml"),
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Password",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Password: testSecretsFileContentsComplete.Password}}, "test_secrets1.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Password: testSecretsFileContentsComplete.Password}}, "test_secrets1a.toml"),
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Pyroscope",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Pyroscope: testSecretsFileContentsComplete.Pyroscope}}, "test_secrets1.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Pyroscope: testSecretsFileContentsComplete.Pyroscope}}, "test_secrets1a.toml"),
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Prometheus",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Prometheus: testSecretsFileContentsComplete.Prometheus}}, "test_secrets1.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Prometheus: testSecretsFileContentsComplete.Prometheus}}, "test_secrets1a.toml"),
				},
			},
			wantErr: true,
		},
				{
			name: "reading multiple secrets with overrides: Mercury",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Mercury: testSecretsFileContentsComplete.Mercury}}, "test_secrets1.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Mercury: testSecretsFileContentsComplete.Mercury}}, "test_secrets1a.toml"),
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Threshold",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Threshold: testSecretsFileContentsComplete.Threshold}}, "test_secrets1.toml"),
					makeTestFile(t, chainlink.Secrets{Secrets: v2.Secrets{Threshold: testSecretsFileContentsComplete.Threshold}}, "test_secrets1a.toml"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.envVar != "" {
				t.Setenv(string(env.Config), tt.args.envVar)
			}
			cfg, err := initServerConfig(tt.args.opts, tt.args.fileNames, tt.args.secretsFiles)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadOpts() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantCfg, cfg)
		})
	}
}
