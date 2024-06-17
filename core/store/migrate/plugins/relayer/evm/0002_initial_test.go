package evm

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	type args struct {
		val Cfg
	}
	tests := []struct {
		name             string
		args             args
		wantErr          bool
		wantGoMigrations goose.Migrations
	}{
		{
			name: "evm template",
			args: args{
				val: Cfg{
					Schema:  "evm",
					ChainID: 3266,
				},
			},
			wantGoMigrations: goose.Migrations{
				&goose.Migration{
					Type:    "go",
					Version: 1,
					Source:  "ignore_this_prefix/0001_initial.go",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Register(tt.args.val)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tDir := t.TempDir()
				m, gerr := goose.CollectMigrations(tDir, 0, 1000)
				require.NoError(t, gerr)
				assert.Len(t, m, len(tt.wantGoMigrations))
				for i, m := range m {
					assert.Equal(t, tt.wantGoMigrations[i].Type, m.Type)
					assert.Equal(t, tt.wantGoMigrations[i].Version, m.Version)
					assert.Equal(t, filepath.Base(tt.wantGoMigrations[i].Source), filepath.Base(m.Source))
				}
			}
		})
	}
}

func Test_resolveup(t *testing.T) {
	type args struct {
		val Cfg
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "evm template",
			args: args{
				val: Cfg{
					Schema:  "evm",
					ChainID: 3266,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := resolveUp(out, tt.args.val)
			require.NoError(t, err)
			assert.NotEmpty(t, out.String())
		})
	}
}

func Test_init_functional(t *testing.T) {
	_, db := heavyweight.FullTestDBEmptyV2(t, nil)
	defer db.Close()

	b, err := os.ReadFile("./testdata/evm_initial_state.sql")
	require.NoError(t, err)

	_, err = db.DB.Exec(string(b))
	require.NoError(t, err, "failed to load initial state")

	type args struct {
		cfg Cfg
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "evm template",
			args: args{
				cfg: Cfg{
					Schema:  "evm_3266",
					ChainID: 3266,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = Register(tt.args.cfg)
			require.NoError(t, err)

			// we need a table to store the goose version for this cfg
			goose.SetTableName(fmt.Sprintf("goose_version_%s", tt.args.cfg.Schema))
			// run the migrations from the embedded templates
			goose.SetBaseFS(embeddedTmplFS)
			err = goose.Up(db.DB, MigrationRootDir)
			require.NoError(t, err, "failed to run migrations")

			// test that the migrations were applied
			goose.Status(db.DB, MigrationRootDir)

			rows, err := db.DB.Query("SELECT schemaname, tablename FROM pg_catalog.pg_tables where schemaname = $1", tt.args.cfg.Schema)
			require.NoError(t, err)
			defer rows.Close()
			for rows.Next() {
				var schema, table string
				err = rows.Scan(&schema, &table)
				t.Logf("schema: %s, table: %s", schema, table)
			}
			// check the error from rows
			err = rows.Err()
			require.NoError(t, err)
		})
	}
	t.FailNow()
}
