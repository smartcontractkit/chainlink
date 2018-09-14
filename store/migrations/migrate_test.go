package migrations_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/migrations"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testGarbageModel struct {
	Garbage string `json:"garbage" storm:"id"`
}

type testMigration0000000001 struct {
	run bool
}

func (m *testMigration0000000001) Migrate(orm *orm.ORM) error {
	m.run = true
	return orm.InitializeModel(testGarbageModel{})
}

func (m *testMigration0000000001) Timestamp() string {
	return "0000000001"
}

func TestMigrate_RunNewMigrations(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	tm := &testMigration0000000001{}
	migrations.ExportedRegisterMigration(tm)

	timestamps := migrations.ExportedAvailableMigrationTimestamps()
	assert.Equal(t, "0", timestamps[0], "Should have initial migration available")
	assert.Equal(t, tm.Timestamp(), timestamps[1], "New test migration should have been registered")

	var migrationTimestamps []migrations.MigrationTimestamp
	assert.NoError(t, store.AllByIndex("Timestamp", &migrationTimestamps))
	assert.Equal(t, migrationTimestamps[0].Timestamp, "0", "Initial migration should have run in NewStore")
	assert.NotContains(t, migrationTimestamps, migrations.MigrationTimestamp{tm.Timestamp()}, "Migration should have not yet run")

	err := migrations.Migrate(store.ORM)
	require.NoError(t, err)

	assert.True(t, tm.run, "Migration should have run")

	err = store.AllByIndex("Timestamp", &migrationTimestamps)
	assert.NoError(t, err)
	assert.Equal(t, tm.Timestamp(), migrationTimestamps[1].Timestamp, "Migration should have been registered as run")
}
