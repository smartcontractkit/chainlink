package migration1537223654_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1537223654"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1537223654/old"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate1537223654_updatesJobSpecsBucket(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	input := cltest.LoadJSON("../../../internal/fixtures/migrations/1537223654_job_without_initiator_params.json")
	var js1 old.JobSpec
	require.NoError(t, json.Unmarshal(input, &js1))
	require.NoError(t, store.ORM.DB.Save(&js1))

	migration := migration1537223654.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var js2 migration1537223654.JobSpec
	require.NoError(t, store.One("ID", js1.ID, &js2))

	assert.Equal(t, js1.Initiators[0].Schedule, js2.Initiators[0].Schedule)
}

func TestMigrate1537223654_updatesInitiatorsBucket(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	input := cltest.LoadJSON("../../../internal/fixtures/migrations/1537223654_initiator_without_params.json")
	var i1 old.Initiator
	require.NoError(t, json.Unmarshal(input, &i1))
	require.NoError(t, store.ORM.DB.Save(&i1))

	migration := migration1537223654.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var i2 migration1537223654.Initiator
	require.NoError(t, store.One("ID", i1.ID, &i2))

	assert.Equal(t, i1.ID, i2.ID)
	assert.Equal(t, i1.Schedule, i2.Schedule)
}

func TestMigrate1537223654_updatesJobRunsBucket(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	input := cltest.LoadJSON("../../../internal/fixtures/migrations/1537223654_jobrun_without_initiator_params.json")
	var jr1 old.JobRun
	require.NoError(t, json.Unmarshal(input, &jr1))
	require.NoError(t, store.ORM.DB.Save(&jr1))

	migration := migration1537223654.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var jr2 migration1537223654.JobRun
	require.NoError(t, store.One("ID", jr1.ID, &jr2))

	assert.Equal(t, jr1.ID, jr2.ID)
	assert.Equal(t, jr1.Initiator.ID, jr2.Initiator.ID)
	assert.Equal(t, jr1.Initiator.JobID, jr2.Initiator.JobID)
	assert.Equal(t, jr1.Initiator.Schedule, jr2.Initiator.Schedule)
}
