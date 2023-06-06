package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var (
	setInFile = "set in config file"
	setInEnv  = "set in env"

	testEnvContents = fmt.Sprintf("P2P.V2.AnnounceAddresses = ['%s']", setInEnv)

	testConfigFileContents = chainlink.Config{
		Core: v2.Core{
			RootDir: &setInFile,
			P2P: v2.P2P{
				V2: v2.P2PV2{
					AnnounceAddresses: &[]string{setInFile},
					ListenAddresses:   &[]string{setInFile},
				},
			},
		},
	}

	testSecretsFileContents = chainlink.Secrets{
		Secrets: v2.Secrets{
			Prometheus: v2.PrometheusSecrets{
				AuthToken: models.NewSecret("PROM_TOKEN"),
			},
		},
	}

	testSecretsRedactedContents = chainlink.Secrets{
		Secrets: v2.Secrets{
			Prometheus: v2.PrometheusSecrets{
				AuthToken: models.NewSecret("xxxxx"),
			},
		},
	}
)

func makeTestFile(t *testing.T, contents any, fileName string) string {
	d := t.TempDir()
	p := filepath.Join(d, fileName)

	b, err := toml.Marshal(contents)
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
		opts        *chainlink.GeneralConfigOpts
		fileNames   []string
		secretsFile string
		envVar      string
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
				Core: v2.Core{
					P2P: v2.P2P{
						V2: v2.P2PV2{
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
				Core: v2.Core{
					RootDir: &setInFile,
					P2P: v2.P2P{
						V2: v2.P2PV2{
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
				opts:        new(chainlink.GeneralConfigOpts),
				fileNames:   []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFile: "/doesnt-exist",
			},
			wantErr: true,
		},
		{
			name: "reading secrets",
			args: args{
				opts:        new(chainlink.GeneralConfigOpts),
				fileNames:   []string{makeTestFile(t, testConfigFileContents, "test.toml")},
				secretsFile: makeTestFile(t, testSecretsFileContents, "test_secrets.toml"),
			},
			wantCfg: withDefaults(t, testConfigFileContents, testSecretsRedactedContents),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.envVar != "" {
				t.Setenv(string(v2.EnvConfig), tt.args.envVar)
			}
			cfg, err := initServerConfig(tt.args.opts, tt.args.fileNames, tt.args.secretsFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadOpts() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, cfg, tt.wantCfg)
		})
	}
}
