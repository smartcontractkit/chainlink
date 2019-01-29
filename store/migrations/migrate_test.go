package migrations_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/migrations"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestMigrate_RunNewMigrations(t *testing.T) {
	migrations.ExportedClearRegisteredMigrations()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	db := store.ORM.DB
	tm := &testMigration0000000001{}
	migrations.ExportedRegisterMigration(tm)

	timestamps := migrations.ExportedAvailableMigrationTimestamps()
	assert.Len(t, timestamps, 1)
	assert.Equal(t, tm.Timestamp(), timestamps[0], "New test migration should have been registered")

	var migrationTimestamps []migrations.MigrationTimestamp
	assert.NoError(t, db.Order("timestamp asc").Find(&migrationTimestamps).Error)
	assert.NotContains(t, migrationTimestamps, migrations.MigrationTimestamp{Timestamp: tm.Timestamp()}, "Migration should have not yet run")

	err := migrations.Migrate(store.ORM)
	require.NoError(t, err)

	assert.True(t, tm.run, "Migration should have run")

	err = db.Order("timestamp asc").Find(&migrationTimestamps).Error
	assert.NoError(t, err)
	assert.Equal(t, tm.Timestamp(), migrationTimestamps[0].Timestamp, "Migration should have been registered as run")
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

func TestMigrate_ErrorOnFailedMigration(t *testing.T) {
	migrations.ExportedClearRegisteredMigrations()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	db := store.ORM.DB
	tm := &testMigration0000000002{}
	migrations.ExportedRegisterMigration(tm)

	err := migrations.Migrate(store.ORM)
	require.Error(t, err)

	assert.True(t, tm.run, "Migration should not have run")

	var migrationTimestamps []migrations.MigrationTimestamp
	err = db.Order("timestamp asc").Find(&migrationTimestamps).Error
	assert.NoError(t, err)

	assert.Len(t, migrationTimestamps, 0, "Migration should have been registered as run")
}

func TestMigrate_Migration0(t *testing.T) {
	migrations.ExportedClearRegisteredMigrations()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	db := store.ORM.DB
	tm := &migration0.Migration{}
	migrations.ExportedRegisterMigration(tm)

	timestamps := migrations.ExportedAvailableMigrationTimestamps()
	assert.Equal(t, "0", timestamps[0], "Should have migration 0 available")

	err := migrations.Migrate(store.ORM)
	require.NoError(t, err)

	var migrationTimestamps []migrations.MigrationTimestamp
	err = db.Order("timestamp asc").Find(&migrationTimestamps).Error
	assert.NoError(t, err)
	assert.Equal(t, tm.Timestamp(), migrationTimestamps[0].Timestamp, "Migration should have been registered as run")
}
