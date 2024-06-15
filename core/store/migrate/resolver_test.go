package migrate

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_resolve(t *testing.T) {
	fixture := "./template/relayers/evm/0002_b.tmpl.sql"
	type args struct {
		val SQLConfig
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
				val: SQLConfig{
					Schema: "evm",
				},
			},
			wantOut: `-- +goose Up
CREATE TABLE evm.bcf_3266_01 (
    "id" TEXT PRIMARY KEY
);
-- +goose Down
DROP TABLE evm.bcf_3266_01;`,
		},

		{
			name: "optimism template",
			args: args{
				val: SQLConfig{
					Schema: "optimism",
				},
			},
			wantOut: `-- +goose Up
CREATE TABLE optimism.bcf_3266_01 (
    "id" TEXT PRIMARY KEY
);
-- +goose Down
DROP TABLE optimism.bcf_3266_01;`,
		},
	}
	for _, tt := range tests {
		testInput, err := os.ReadFile(fixture)
		require.NoError(t, err)
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := resolve(out, bytes.NewBuffer(testInput), tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantOut, out.String(), cmp.Diff(tt.wantOut, out.String()))
		})
	}
}

func Test_generateMigrations(t *testing.T) {
	tDir := t.TempDir()
	type args struct {
		rootDir string
		tmpDir  string
		val     SQLConfig
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
				rootDir: "./",
				tmpDir:  filepath.Join(tDir, "evm"),
				val: SQLConfig{
					Schema: "evm",
				},
			},
			want: []string{
				filepath.Join(tDir, "evm/0001_a.sql"),
				filepath.Join(tDir, "evm/0002_b.sql"),
				filepath.Join(tDir, "evm/0003_c.sql")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateMigrations(tt.args.rootDir, tt.args.tmpDir, tt.args.val)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
