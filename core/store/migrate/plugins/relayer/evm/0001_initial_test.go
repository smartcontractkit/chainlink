package evm

import (
	_ "embed"
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

func Test_init_functional(t *testing.T) {
	tDir := t.TempDir()
	c, db := heavyweight.FullTestDBEmptyV2(t, nil)
	defer db.Close()

	t.Logf("url %s", c.Database().URL().Path)
	cfg := Cfg{
		Schema:  "evm",
		ChainID: 3266,
	}
	t.Logf("%v", os.Environ())
	err := Register(cfg)
	require.NoError(t, err)

	goose.SetTableName("goose_evm_3266_version")
	md := filepath.Join(tDir, "migrations")
	require.NoError(t, os.Mkdir(md, os.ModePerm))
	err = goose.Up(db.DB, md)
	require.NoError(t, err, "failed to run migrations")

	// test that the migrations were applied
	goose.Status(db.DB, md)

	rows, err := db.DB.Query("SELECT schemaname, tablename FROM pg_catalog.pg_tables")
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

	t.FailNow()
}
