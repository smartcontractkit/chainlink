package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/stretchr/testify/require"
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
)

func makeTestConfigFile(t *testing.T) string {
	d := t.TempDir()
	p := filepath.Join(d, "test.toml")

	b, err := toml.Marshal(testConfigFileContents)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(p, b, 0777))
	return p
}

func Test_loadOpts(t *testing.T) {
	type args struct {
		opts      *chainlink.GeneralConfigOpts
		fileNames []string
		envVar    string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantOpts *chainlink.GeneralConfigOpts
	}{
		{
			name: "env only",
			args: args{
				opts:   new(chainlink.GeneralConfigOpts),
				envVar: testEnvContents,
			},
			wantOpts: &chainlink.GeneralConfigOpts{
				Config: chainlink.Config{
					Core: v2.Core{
						P2P: v2.P2P{
							V2: v2.P2PV2{
								AnnounceAddresses: &[]string{setInEnv},
							},
						},
					},
				},
			},
		},

		{
			name: "files only",
			args: args{
				opts:      new(chainlink.GeneralConfigOpts),
				fileNames: []string{makeTestConfigFile(t)},
			},
			wantOpts: &chainlink.GeneralConfigOpts{
				Config: testConfigFileContents,
			},
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
				fileNames: []string{makeTestConfigFile(t)},
				envVar:    testEnvContents,
			},
			wantOpts: &chainlink.GeneralConfigOpts{
				Config: chainlink.Config{
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
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.envVar != "" {
				t.Setenv(string(v2.EnvConfig), tt.args.envVar)
			}
			if err := loadOpts(tt.args.opts, tt.args.fileNames...); (err != nil) != tt.wantErr {
				t.Errorf("loadOpts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
