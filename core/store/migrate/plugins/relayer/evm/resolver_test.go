package evm

import (
	"bytes"
	"embed"
	"path/filepath"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_resolve(t *testing.T) {
	type args struct {
		in  string
		val Cfg
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "evm template",
			args: args{
				val: Cfg{
					Schema:  "evm",
					ChainID: 3266,
				},
				in: "schema={{.Schema}}, chainID={{.ChainID}}",
			},
			wantOut: "schema=evm, chainID=3266",
		},
		{
			name: "unknown template",
			args: args{
				val: Cfg{
					Schema:  "evm",
					ChainID: 3266,
				},
				in: "schema={{.WrongField}}, chainID={{.ChainID}}",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := resolve(out, tt.args.in, tt.args.val)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantOut, out.String())
			}
		})
	}
}

func Test_generateMigrations(t *testing.T) {
	tDir := t.TempDir()
	type args struct {
		fsys    embed.FS
		rootDir string
		tmpDir  string
		val     Cfg
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "evm template",
			args: args{
				fsys:    embeddedTmplFS,
				rootDir: MigrationRootDir,
				tmpDir:  filepath.Join(tDir, "evm_42"),
				val: Cfg{
					Schema:  "evm_42",
					ChainID: 42,
				},
			},
			want: []string{
				filepath.Join(tDir, "evm_42/0001_create_schema.sql"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateMigrations(tt.args.fsys, tt.args.rootDir, tt.args.tmpDir, tt.args.val)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
