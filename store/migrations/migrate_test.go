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

type testMigration1536545228 struct{}

func (m testMigration1536545228) Migrate(orm *orm.ORM) error {
	return orm.InitializeModel(testGarbageModel{})
}

func (m testMigration1536545228) Timestamp() string {
	return "0000000001"
}

func TestMigrate_RunNewMigrations(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	migrations.ExportedRegisterMigration(testMigration1536545228{})

	timestamps := migrations.ExportedAvailableMigrationTimestamps()
	assert.Equal(t, testMigration1536545228{}.Timestamp(), timestamps[0])

	err := migrations.Migrate(store.ORM)
	require.NoError(t, err)

	var migrationTimestamps []migrations.MigrationTimestamp
	err = store.AllByIndex("Timestamp", &migrationTimestamps)
	assert.NoError(t, err)
	assert.Equal(t, testMigration1536545228{}.Timestamp(), migrationTimestamps[0].Timestamp)
}
