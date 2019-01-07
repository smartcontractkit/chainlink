package migration1536696950_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536696950"
	"github.com/stretchr/testify/require"
)

func TestMigrate1536696950_ConvertRunResultAmount(t *testing.T) {
	input := cltest.LoadJSON("../../../internal/fixtures/migrations/1536696950_jobrun_bigint_amount.json")
	var jr migration0.JobRun
	require.NoError(t, json.Unmarshal(input, &jr))

	migration := migration1536696950.Migration{}
	jr2 := migration.Convert(jr)

	require.Equal(t, toBig(jr2.Result.Amount), big.NewInt(1000000000000000000))
	require.Equal(t, toBig(jr2.TaskRuns[0].Result.Amount), big.NewInt(1000000000000000000))
}

func TestMigrate1536696950_MigrateRunResultAmount1536521223(t *testing.T) {
	input := cltest.LoadJSON("../../../internal/fixtures/migrations/1536696950_jobrun_bigint_amount.json")
	var jr migration0.JobRun
	require.NoError(t, json.Unmarshal(input, &jr))

	store, cleanup := cltest.NewStore()
	defer cleanup()

	require.NoError(t, store.ORM.DB.Save(&jr))

	migration := migration1536696950.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var jr2 migration1536696950.JobRun
	require.NoError(t, store.One("ID", jr.ID, &jr2))
	require.Equal(t, toBig(jr2.Result.Amount), big.NewInt(1000000000000000000))
	require.Equal(t, toBig(jr2.TaskRuns[0].Result.Amount), big.NewInt(1000000000000000000))
}

func TestMigrate1536696950_MigrateRunResultAmount1536521223_asString(t *testing.T) {
	input := cltest.LoadJSON("../../../internal/fixtures/migrations/1536696950_jobrun_string_amount.json")
	var jr migration1536696950.JobRun
	require.NoError(t, json.Unmarshal(input, &jr))

	store, cleanup := cltest.NewStore()
	defer cleanup()

	require.NoError(t, store.ORM.DB.Save(&jr))

	migration := migration1536696950.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var jr2 migration1536696950.JobRun
	require.NoError(t, store.One("ID", jr.ID, &jr2))
	require.Equal(t, toBig(jr2.Result.Amount), big.NewInt(1000000000000000000))
	require.Equal(t, toBig(jr2.TaskRuns[0].Result.Amount), big.NewInt(1000000000000000000))
}

func toBig(l *migration0.Link) *big.Int {
	copy := *l
	return (*big.Int)(&copy)
}
