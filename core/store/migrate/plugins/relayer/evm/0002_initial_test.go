package evm

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
					ChainID: big.NewI(int64(3266)),
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

	// must load the initial state, derivied from the core migrations at v244
	// because the evm migrations try to move try from the core schema to the evm schema
	b, err := os.ReadFile("./testutils/evm_initial_state.sql")
	require.NoError(t, err)

	_, err = db.DB.Exec(string(b))
	require.NoError(t, err, "failed to load initial state")

	type args struct {
		cfg Cfg
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantTables []string
	}{
		{
			name: "evm template",
			args: args{
				cfg: Cfg{
					Schema:  "evm_3266",
					ChainID: big.NewI(int64(3266)),
				},
			},
			wantTables: []string{
				"forwarders",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = Register0002(tt.args.cfg)
			require.NoError(t, err)

			// we need a table to store the goose version for this cfg
			goose.SetTableName(fmt.Sprintf("goose_version_%s", tt.args.cfg.Schema))
			// run the migrations from the embedded templates

			tDir := t.TempDir()
			_, err := generateMigrations(embeddedTmplFS, MigrationRootDir, tDir, tt.args.cfg)
			require.NoError(t, err, "failed to generate migrations")
			//goose.SetBaseFS(tDir)
			err = goose.UpTo(db.DB, tDir, 2)
			require.NoError(t, err, "failed to run migrations")

			// test that the migrations were applied
			goose.Status(db.DB, MigrationRootDir)

			rows, err := db.DB.Query("SELECT schemaname, tablename FROM pg_catalog.pg_tables where schemaname = $1", tt.args.cfg.Schema)
			require.NoError(t, err)
			defer rows.Close()
			var gotTables []string
			for rows.Next() {
				var schema, table string
				err = rows.Scan(&schema, &table)
				t.Logf("schema: %s, table: %s", schema, table)
				require.NoError(t, err)
				gotTables = append(gotTables, table)
			}
			// check the error from rows
			err = rows.Err()
			require.NoError(t, err)
			assert.Equal(t, tt.wantTables, gotTables)

			// check the goose version
			v, err := goose.EnsureDBVersion(db.DB)
			require.NoError(t, err)
			assert.Equal(t, int64(2), v)

			// run the down migrations
			err = goose.DownTo(db.DB, tDir, 0)
			require.NoError(t, err, "failed to run down migrations")
			v, err = goose.EnsureDBVersion(db.DB)
			require.NoError(t, err)
			rows, err = db.DB.Query("SELECT schemaname, tablename FROM pg_catalog.pg_tables where schemaname = $1", tt.args.cfg.Schema)
			require.NoError(t, err)
			defer rows.Close()
			var gotDownTables []string
			for rows.Next() {
				var schema, table string
				err = rows.Scan(&schema, &table)
				t.Logf("schema: %s, table: %s", schema, table)
				require.NoError(t, err)
				gotTables = append(gotTables, table)
			}
			assert.Len(t, gotDownTables, 0)
		})
	}
	t.FailNow()
}
