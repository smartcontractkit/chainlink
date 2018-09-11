package migrations_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models/migrations"
	"github.com/smartcontractkit/chainlink/store/models/migrations/migration1536521223"
	"github.com/smartcontractkit/chainlink/store/models/migrations/migration1536696950"
	"github.com/smartcontractkit/chainlink/store/models/orm"
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

func TestMigrate_ConvertRunResultAmount(t *testing.T) {
	input := cltest.LoadJSON("../../../internal/fixtures/bolt/jobrun_bigint_amount.json")
	var jr migration1536521223.JobRun
	require.NoError(t, json.Unmarshal(input, &jr))

	migration := migration1536696950.Migration{}
	jr2 := migration.Convert(jr)

	require.Equal(t, jr2.Result.Amount, assets.NewLink(1000000000000000000))
	require.Equal(t, jr2.TaskRuns[0].Result.Amount, assets.NewLink(1000000000000000000))
}

func TestMigrate_MigrateRunResultAmount1536521223(t *testing.T) {
	input := cltest.LoadJSON("../../../internal/fixtures/bolt/jobrun_bigint_amount.json")
	var jr migration1536521223.JobRun
	require.NoError(t, json.Unmarshal(input, &jr))

	store, cleanup := cltest.NewStore()
	defer cleanup()

	require.NoError(t, store.Save(&jr))

	migration := migration1536696950.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var jr2 migration1536696950.JobRun
	require.NoError(t, store.One("ID", jr.ID, &jr2))
	require.Equal(t, jr2.Result.Amount, assets.NewLink(1000000000000000000))
	require.Equal(t, jr2.TaskRuns[0].Result.Amount, assets.NewLink(1000000000000000000))
}
