package migrations_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/migrations"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1551816486"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1551895034"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1551895034/old"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate_RunNewMigrations(t *testing.T) {
	migrations.ExportedClearRegisteredMigrations()

	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB
	tm := &testMigration0000000001{}
	migrations.ExportedRegisterMigration(tm)

	timestamps := migrations.ExportedAvailableMigrationTimestamps()
	assert.Len(t, timestamps, 1)
	assert.Equal(t, tm.Timestamp(), timestamps[0], "New test migration should have been registered")

	err := migrations.Migrate(orm)
	require.NoError(t, err)

	assert.True(t, tm.run, "Migration should have run")

	var migrationTimestamps []migrations.MigrationTimestamp
	err = db.Order("timestamp asc").Find(&migrationTimestamps).Error
	assert.NoError(t, err)
	assert.Equal(t, tm.Timestamp(), migrationTimestamps[0].Timestamp, "Migration should have been registered as run")
}

func TestMigrate_ErrorOnFailedMigration(t *testing.T) {
	migrations.ExportedClearRegisteredMigrations()

	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB
	tm := &testMigration0000000002{}
	migrations.ExportedRegisterMigration(tm)

	err := migrations.Migrate(orm)
	require.Error(t, err)

	assert.True(t, tm.run, "Migration should not have run")

	var migrationTimestamps []migrations.MigrationTimestamp
	err = db.Order("timestamp asc").Find(&migrationTimestamps).Error
	assert.NoError(t, err)

	assert.Len(t, migrationTimestamps, 0, "Migration should have been registered as run")
}

func TestMigrate_Migration0(t *testing.T) {
	migrations.ExportedClearRegisteredMigrations()

	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB
	tm := &migration0.Migration{}
	migrations.ExportedRegisterMigration(tm)

	timestamps := migrations.ExportedAvailableMigrationTimestamps()
	assert.Equal(t, "0", timestamps[0], "Should have migration 0 available")

	err := migrations.Migrate(orm)
	require.NoError(t, err)

	var migrationTimestamps []migrations.MigrationTimestamp
	err = db.Order("timestamp asc").Find(&migrationTimestamps).Error
	assert.NoError(t, err)
	assert.Equal(t, tm.Timestamp(), migrationTimestamps[0].Timestamp, "Migration should have been registered as run")
}

func bootstrapORM(t *testing.T) (*orm.ORM, func()) {
	tc, cleanup := cltest.NewConfig()
	config := tc.Config

	require.NoError(t, os.MkdirAll(config.RootDir(), 0700))
	cltest.WipePostgresDatabase(tc.Config)

	orm, err := orm.NewORM(config.NormalizedDatabaseURL(), config.DatabaseTimeout())
	require.NoError(t, err)

	return orm, func() {
		assert.NoError(t, orm.Close())
		cleanup()
		os.RemoveAll(config.RootDir())
	}
}

type testGarbageModel struct {
	Garbage int `json:"garbage" gorm:"primary_key"`
}

type testMigration0000000001 struct {
	run bool
}

func (m *testMigration0000000001) Migrate(orm *orm.ORM) error {
	m.run = true
	return orm.DB.AutoMigrate(&testGarbageModel{}).Error
}

func (m *testMigration0000000001) Timestamp() string {
	return "0000000001"
}

type testFailingModel struct{}

type testMigration0000000002 struct {
	run bool
}

func (m *testMigration0000000002) Migrate(orm *orm.ORM) error {
	m.run = true
	var result string
	err := orm.DB.Raw("SELECT * FROM non_existent_table;").Scan(&result).Error
	return err
}

func (m *testMigration0000000002) Timestamp() string {
	return "0000000002"
}

func TestMigrate1551816486(t *testing.T) {
	migrations.ExportedClearRegisteredMigrations()

	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	tm := &migration1551816486.Migration{}
	migrations.ExportedRegisterMigration(tm)

	// seed db w old table
	err := orm.DB.Exec(`
		CREATE TABLE "bridge_types" (
			"name" varchar(255),
			"url" varchar(255),
			"confirmations" bigint,
			"incoming_token" varchar(255),
			"outgoing_token" varchar(255),
			"minimum_contract_payment" varchar(255),
			UNIQUE (name));
	`).Error

	require.NoError(t, err)

	initial := models.BridgeType{
		Name: "someUniqueName",
		URL:  cltest.WebURL("http://someurl.com"),
	}

	require.NoError(t, orm.DB.Save(&initial).Error)
	require.NoError(t, migrations.Migrate(orm))

	migratedbt, err := orm.FindBridge(initial.Name.String())
	require.NoError(t, err)
	require.Equal(t, initial, migratedbt)
}

func TestMigrate1551895034(t *testing.T) {
	migrations.ExportedClearRegisteredMigrations()

	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	tm := &migration1551895034.Migration{}
	migrations.ExportedRegisterMigration(tm)

	height := models.NewBig(big.NewInt(1337))
	hash := common.HexToHash("0xde3fb1df888c6c7f77f3a8e9c2582f87e7ad5277d98bd06cfd17cd2d7ea49f42")

	previous := old.IndexableBlockNumber{
		Number: *height,
		Digits: 4,
		Hash:   hash,
	}
	// seed w old schema and data
	err := orm.DB.AutoMigrate(old.IndexableBlockNumber{}).Error
	require.NoError(t, err)
	err = orm.DB.Save(&previous).Error
	require.NoError(t, err)

	// migrate
	require.NoError(t, migrations.Migrate(orm))

	retrieved := models.Head{}
	err = orm.DB.First(&retrieved).Error
	require.NoError(t, err)

	require.Equal(t, height.ToInt(), retrieved.ToInt())
	require.Equal(
		t,
		hash.String(),
		retrieved.Hash().Hex())
}
