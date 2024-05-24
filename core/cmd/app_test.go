package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
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
)

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
				fileNames: []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
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
				fileNames: []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
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
				fileNames:    []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{"/doesnt-exist"},
			},
			wantErr: true,
		},
		{
			name: "reading secrets",
			args: args{
				opts:         new(chainlink.GeneralConfigOpts),
				fileNames:    []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{configtest.WriteTOMLFile(t, testSecretsFileContents, "test_secrets.toml")},
			},
			wantCfg: withDefaults(t, testConfigFileContents, testSecretsRedactedContents),
		},
		{
			name: "reading multiple secrets",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					"../services/chainlink/testdata/mergingsecretsdata/secrets-database.toml",
					"../services/chainlink/testdata/mergingsecretsdata/secrets-password.toml",
					"../services/chainlink/testdata/mergingsecretsdata/secrets-pyroscope.toml",
					"../services/chainlink/testdata/mergingsecretsdata/secrets-prometheus.toml",
					"../services/chainlink/testdata/mergingsecretsdata/secrets-mercury-split-one.toml",
					"../services/chainlink/testdata/mergingsecretsdata/secrets-mercury-split-two.toml",
					"../services/chainlink/testdata/mergingsecretsdata/secrets-threshold.toml",
					"../services/chainlink/testdata/mergingsecretsdata/secrets-webserver-ldap.toml",
				},
			},
			wantErr: false,
		},
		{
			name: "reading multiple secrets with overrides: Database",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					"../testdata/mergingsecretsdata/secrets-database.toml",
					"../testdata/mergingsecretsdata/secrets-database.toml",
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Password",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					"../testdata/mergingsecretsdata/secrets-password.toml",
					"../testdata/mergingsecretsdata/secrets-password.toml",
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Pyroscope",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					"../testdata/mergingsecretsdata/secrets-pyroscope.toml",
					"../testdata/mergingsecretsdata/secrets-pyroscope.toml",
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Prometheus",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					"../testdata/mergingsecretsdata/secrets-prometheus.toml",
					"../testdata/mergingsecretsdata/secrets-prometheus.toml",
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Mercury",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					"../testdata/mergingsecretsdata/secrets-mercury-split-one.toml",
					"../testdata/mergingsecretsdata/secrets-mercury-split-one.toml",
				},
			},
			wantErr: true,
		},
		{
			name: "reading multiple secrets with overrides: Threshold",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{configtest.WriteTOMLFile(t, testConfigFileContents, "test.toml")},
				secretsFiles: []string{
					"../testdata/mergingsecretsdata/secrets-threshold.toml",
					"../testdata/mergingsecretsdata/secrets-threshold.toml",
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
			if tt.wantCfg != nil {
				assert.Equal(t, tt.wantCfg, cfg)
			}
		})
	}
}
