package migrations_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1559081901"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1559767166"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560433987"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560791143"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560881846"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560886530"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormigrate "gopkg.in/gormigrate.v1"
)

func bootstrapORM(t *testing.T) (*gorm.DB, func()) {
	tc, cleanup := cltest.NewConfig(t)
	cfg := tc.Depot

	require.NoError(t, os.MkdirAll(cfg.RootDir(), 0700))
	cltest.WipePostgresDatabase(t, cfg)

	url := orm.NormalizedDatabaseURL(cfg)
	dialect, err := orm.DeduceDialect(url)
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(string(dialect), url)
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		assert.NoError(t, db.Close())
		cleanup()
		os.RemoveAll(cfg.RootDir())
	}
}

func TestMigrate_Migrations(t *testing.T) {
	db, cleanup := bootstrapORM(t)
	defer cleanup()

	require.NoError(t, migration0.Migrate(db))
	require.NoError(t, migration1559081901.Migrate(db))

	assert.True(t, db.HasTable("bridge_types"))
	assert.True(t, db.HasTable("encumbrances"))
	assert.True(t, db.HasTable("external_initiators"))
	assert.True(t, db.HasTable("heads"))
	assert.True(t, db.HasTable("job_specs"))
	assert.True(t, db.HasTable("initiators"))
	assert.True(t, db.HasTable("job_runs"))
	assert.True(t, db.HasTable("keys"))
	assert.True(t, db.HasTable("run_requests"))
	assert.True(t, db.HasTable("run_results"))
	assert.True(t, db.HasTable("service_agreements"))
	assert.True(t, db.HasTable("sessions"))
	assert.True(t, db.HasTable("task_runs"))
	assert.True(t, db.HasTable("task_specs"))
	assert.True(t, db.HasTable("tx_attempts"))
	assert.True(t, db.HasTable("txes"))
	assert.True(t, db.HasTable("users"))
}

func TestMigrate_Migration1560791143(t *testing.T) {
	db, cleanup := bootstrapORM(t)
	defer cleanup()

	require.NoError(t, migration0.Migrate(db))

	tx := migration0.Tx{
		ID:       1337,
		Data:     make([]byte, 10),
		Value:    models.NewBig(big.NewInt(1)),
		GasPrice: models.NewBig(big.NewInt(127)),
	}
	require.NoError(t, db.Create(&tx).Error)

	require.NoError(t, migration1559081901.Migrate(db))

	txFound := models.Tx{}
	require.NoError(t, db.Where("id = ?", tx.ID).Find(&txFound).Error)

	require.NoError(t, migration1560791143.Migrate(db))

	txNoID := models.Tx{
		Data:     make([]byte, 10),
		Value:    models.NewBig(big.NewInt(2)),
		GasPrice: models.NewBig(big.NewInt(119)),
	}
	require.NoError(t, db.Create(&txNoID).Error)
	assert.Equal(t, uint64(1338), txNoID.ID)

	noIDTxFound := models.Tx{}
	require.NoError(t, db.Where("id = ?", tx.ID).Find(&noIDTxFound).Error)
}

func TestMigrate_Migration1560881846(t *testing.T) {
	db, cleanup := bootstrapORM(t)
	defer cleanup()

	require.NoError(t, migration0.Migrate(db))
	require.NoError(t, migration1559081901.Migrate(db))
	require.NoError(t, migration1559767166.Migrate(db))
	require.NoError(t, migration1560433987.Migrate(db))
	require.NoError(t, migration1560791143.Migrate(db))
	require.NoError(t, migration1560881846.Migrate(db))

	head := migration0.Head{
		HashRaw: "dad0000000000000000000000000000000000000000000000000000000000b0d",
		Number:  8616460799,
	}
	require.NoError(t, db.Create(&head).Error)

	require.NoError(t, migration1560886530.Migrate(db))

	headFound := models.Head{}
	require.NoError(t, db.Where("id = (SELECT MAX(id) FROM heads)").Find(&headFound).Error)
	assert.Equal(t, "0xdad0000000000000000000000000000000000000000000000000000000000b0d", headFound.Hash.Hex())
	assert.Equal(t, int64(8616460799), headFound.Number)
}

func TestMigrate_NewerVersionGuard(t *testing.T) {
	db, cleanup := bootstrapORM(t)
	defer cleanup()

	// Do full migrations
	require.NoError(t, migrations.Migrate(db))

	// Add a fictional future migration
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID:      "9223372036854775807",
			Migrate: migration0.Migrate,
		},
	})
	require.NoError(t, m.Migrate())

	// Run migrations again, should error
	require.Error(t, migrations.Migrate(db))
}
