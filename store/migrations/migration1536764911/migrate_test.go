package migration1536764911_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536764911"
	"github.com/stretchr/testify/require"
)

func TestMigration_ConvertTaskSpec(t *testing.T) {
	input := cltest.LoadJSON("../../../internal/fixtures/bolt/old_job_without_task_params.json")
	var js migration0.JobSpec
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
