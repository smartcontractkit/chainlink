package migration1536764911_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536764911"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536764911/old"
	"github.com/stretchr/testify/require"
)

func TestMigrate1536764911_ConvertTaskSpec(t *testing.T) {
	input := cltest.LoadJSON("../../../internal/fixtures/migrations/1536764911_job_without_task_params.json")
	var js old.JobSpec
	require.NoError(t, json.Unmarshal(input, &js))

	migration := migration1536764911.Migration{}
	js2 := migration.Convert(js)

	require.Equal(t, "0x356a04bce728ba4c62a30294a55e6a8600a320b3", js2.Tasks[3].Params.Get("address").String())
	results := js2.Tasks[1].Params.Get("path").Array()
	var path []string
	for _, r := range results {
		path = append(path, r.String())
	}
	require.Equal(t, []string{"last"}, path)
}

func TestMigrate1536764911_JobRun(t *testing.T) {
	input := cltest.LoadJSON("../../../internal/fixtures/migrations/1536764911_jobrun_without_task_params.json")
	var jr old.JobRun
	require.NoError(t, json.Unmarshal(input, &jr))

	store, cleanup := cltest.NewStore()
	defer cleanup()

	require.NoError(t, store.ORM.DB.Save(&jr))

	migration := migration1536764911.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var jr2 migration1536764911.JobRun
	require.NoError(t, store.One("ID", jr.ID, &jr2))
	require.Equal(t, jr.TaskRuns[0].Task.Params.Get("address").String(), jr2.TaskRuns[0].Task.Params.Get("address").String())
	require.Equal(t, jr.TaskRuns[0].Task.Params.Get("times").Int(), jr2.TaskRuns[0].Task.Params.Get("times").Int())
}
