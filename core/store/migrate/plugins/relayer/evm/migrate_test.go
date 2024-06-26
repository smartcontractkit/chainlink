package evm_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate/plugins/relayer/evm"
)

func TestMigrate(t *testing.T) {
	t.Run("core migration with optional relayer migration", func(t *testing.T) {
		_, db := heavyweight.FullTestDBEmptyV2(t, nil)

		ctx := testutils.Context(t)
		cfg := evm.Cfg{
			Schema:  "evm_42",
			ChainID: big.NewI(int64(42)),
		}
		// the evm migrations only work if the core migrations have been run
		// because we are moving existing tables
		err := evm.Migrate(ctx, db, cfg)
		require.Error(t, err)
		err = migrate.Migrate(ctx, db.DB)
		require.NoError(t, err)

		err = evm.Migrate(ctx, db, cfg)
		require.NoError(t, err)

		v2, err := evm.Current(ctx, db, cfg)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, v2, int64(2))

		err = evm.Rollback(ctx, db, null.IntFrom(0), cfg)
		require.NoError(t, err)

		v2, err = evm.Current(ctx, db, cfg)
		require.NoError(t, err)

		assert.Equal(t, int64(0), v2)
	})
}

func TestGoDataMigration(t *testing.T) {
	var totalRecords = 3 // by convention, each legacy table is loaded with 3 records; 2 for chain 0 and 1 for chain 1
	type test struct {
		table string
	}
	var cases = []test{
		{
			table: "forwarders",
		},
		{
			table: "heads",
		},
		{
			table: "key_states",
		},
	}
	for _, tt := range cases {
		t.Run(tt.table+" data migration", func(t *testing.T) {
			ctx := testutils.Context(t)
			db := loadLegacyDatabase(t, ctx)

			type dataTest struct {
				name                    string
				cfg                     evm.Cfg
				wantMigratedRecordCount int
			}
			var dataCases = []dataTest{
				{
					name: "chain 731",
					cfg: evm.Cfg{
						Schema:  "evm_731",
						ChainID: big.NewI(int64(731)),
					},
					wantMigratedRecordCount: 0,
				},
				{
					name: "chain 0",
					cfg: evm.Cfg{
						Schema:  "evm_0",
						ChainID: big.NewI(int64(0)),
					},
					wantMigratedRecordCount: 2,
				},
				{
					name: "chain 1",
					cfg: evm.Cfg{
						Schema:  "evm_1",
						ChainID: big.NewI(int64(1)),
					},
					wantMigratedRecordCount: 1,
				},
			}
			for _, dataCase := range dataCases {
				t.Run(dataCase.name, func(t *testing.T) {
					err := evm.Migrate(ctx, db, dataCase.cfg)
					require.NoError(t, err)
					var moved int
					err = db.Get(&moved, fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", dataCase.cfg.Schema, tt.table))
					require.NoError(t, err)
					require.Equal(t, dataCase.wantMigratedRecordCount, moved)

					var remaining int
					err = db.Get(&remaining, fmt.Sprintf("SELECT COUNT(*) FROM evm.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, totalRecords-dataCase.wantMigratedRecordCount, remaining)

					err = evm.Rollback(ctx, db, null.IntFrom(0), dataCase.cfg)
					require.NoError(t, err)
					var rollbackTotal int
					err = db.Get(&rollbackTotal, fmt.Sprintf("SELECT COUNT(*) FROM evm.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, totalRecords, rollbackTotal)

					// cfg schema should be gone
					var schemaCount int
					err = db.Get(&schemaCount, fmt.Sprintf("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = '%s'", dataCase.cfg.Schema))
					require.NoError(t, err)
					require.Equal(t, 0, schemaCount)

				})
			}
			/*
							cfg := evm.Cfg{
					Schema:  "evm_731",
					ChainID: big.NewI(int64(731)),
				}

					err := evm.Migrate(ctx, db, cfg)
					require.NoError(t, err)
					// no data for chain 731 in the fixtures by convention
					var cnt int
					moved = -1
					err = db.Get(&moved, fmt.Sprintf("SELECT COUNT(*) FROM evm_731.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, 0, moved)
					// all 3 records should still be in the legacy table
					moved = -1
					err = db.Get(&moved, fmt.Sprintf("SELECT COUNT(*) FROM evm.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, 3, moved)

					// run the migration for chain 0 which has 2 records
					cfg = evm.Cfg{
						Schema:  "evm_0",
						ChainID: big.NewI(int64(0)),
					}
					err = evm.Migrate(ctx, db, cfg)
					require.NoError(t, err)
					err = db.Get(&moved, fmt.Sprintf("SELECT COUNT(*) FROM evm_0.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, 2, moved)
					// the 2 records should have been moved from the legacy table. leaving only 1 record
					err = db.Get(&moved, fmt.Sprintf("SELECT COUNT(*) FROM evm.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, 1, moved)

					// rollback of the migration for chain 0 should move the 2 records back
					err = evm.Rollback(ctx, db, null.IntFrom(0), cfg)
					require.NoError(t, err)
					moved = -1
					err = db.Get(&moved, fmt.Sprintf("SELECT COUNT(*) FROM evm.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, 3, moved)
					// evm_0 schema should be gone
					moved = -1
					err = db.Get(&moved, "SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = 'evm_0'")
					require.NoError(t, err)
					require.Equal(t, 0, moved)

					// run the migration for chain 1 which has 1 record
					cfg = evm.Cfg{
						Schema:  "evm_1",
						ChainID: big.NewI(int64(1)),
					}
					err = evm.Migrate(ctx, db, cfg)
					require.NoError(t, err)
					moved = -1
					err = db.Get(&moved, fmt.Sprintf("SELECT COUNT(*) FROM evm_1.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, 1, moved)
					// the 1 record should have been moved from the legacy table. leaving 2 records
					moved = -1
					err = db.Get(&moved, fmt.Sprintf("SELECT COUNT(*) FROM evm.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, 2, moved)

					// rollback of the migration for chain 1 should move the 1 record back
					err = evm.Rollback(ctx, db, null.IntFrom(0), cfg)
					require.NoError(t, err)
					// the 1 record should have been moved back to the legacy table resulting in 3 records
					moved = -1
					err = db.Get(&moved, fmt.Sprintf("SELECT COUNT(*) FROM evm.%s", tt.table))
					require.NoError(t, err)
					require.Equal(t, 3, moved)
					// evm_1 schema should be gone
					moved = -1
					err = db.Get(&moved, "SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = 'evm_1'")
					require.NoError(t, err)
					require.Equal(t, 0, moved)
			*/
		})
	}
}

// loadLegacyDatabase loads the legacy forwarder and heads data into the database
// as a matter of convenience and convention, each legacy table is loaded with 3 records: 2 for chain 0 and 1 for chain 1
func loadLegacyDatabase(t *testing.T, ctx context.Context) *sqlx.DB {
	t.Helper()

	_, db := heavyweight.FullTestDBEmptyV2(t, nil)
	err := migrate.Migrate(ctx, db.DB)
	require.NoError(t, err)
	// load the legacy forwarder data
	forwarderMigration, err := os.ReadFile("testdata/forwarders/initial.sql")
	require.NoError(t, err)
	_, err = db.Exec(string(forwarderMigration))
	require.NoError(t, err)
	var cnt int
	err = db.Get(&cnt, "SELECT COUNT(*) FROM evm.forwarders")
	require.NoError(t, err)
	require.Equal(t, 3, cnt)

	// load the legacy heads data
	headsMigration, err := os.ReadFile("testdata/heads/initial.sql")
	require.NoError(t, err)
	_, err = db.Exec(string(headsMigration))
	require.NoError(t, err)
	err = db.Get(&cnt, "SELECT COUNT(*) FROM evm.heads")
	require.NoError(t, err)
	require.Equal(t, 3, cnt)

	// load the legacy key_states data
	keyStatesMigration, err := os.ReadFile("testdata/key_states/initial.sql")
	require.NoError(t, err)
	_, err = db.Exec(string(keyStatesMigration))
	require.NoError(t, err)
	err = db.Get(&cnt, "SELECT COUNT(*) FROM evm.key_states")
	require.NoError(t, err)
	require.Equal(t, 3, cnt)
	return db
}
