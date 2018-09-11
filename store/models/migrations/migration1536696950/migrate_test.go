package migration1536696950_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models/migrations/migration1536521223"
	"github.com/smartcontractkit/chainlink/store/models/migrations/migration1536696950"
	"github.com/stretchr/testify/require"
)

func TestMigrate_ConvertRunResultAmount(t *testing.T) {
	input := cltest.LoadJSON("../../../../internal/fixtures/bolt/jobrun_bigint_amount.json")
	var jr migration1536521223.JobRun
	require.NoError(t, json.Unmarshal(input, &jr))

	migration := migration1536696950.Migration{}
	jr2 := migration.Convert(jr)

	require.Equal(t, jr2.Result.Amount, assets.NewLink(1000000000000000000))
	require.Equal(t, jr2.TaskRuns[0].Result.Amount, assets.NewLink(1000000000000000000))
}

func TestMigrate_MigrateRunResultAmount1536521223(t *testing.T) {
	input := cltest.LoadJSON("../../../../internal/fixtures/bolt/jobrun_bigint_amount.json")
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

func TestMigrate_MigrateRunResultAmount1536521223_asString(t *testing.T) {
	input := cltest.LoadJSON("../../../../internal/fixtures/bolt/jobrun_string_amount.json")
	var jr migration1536696950.JobRun
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
