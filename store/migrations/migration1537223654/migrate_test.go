package migration1537223654_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536696950"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536764911"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1537223654"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate_updatesJobSpecsBucket(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	input := cltest.LoadJSON("../../../internal/fixtures/bolt/old_job_without_initiator_params.json")
	var js1 migration1536764911.JobSpec
	require.NoError(t, json.Unmarshal(input, &js1))
	require.NoError(t, store.Save(&js1))

	migration := migration1537223654.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var js2 migration1537223654.JobSpec
	require.NoError(t, store.One("ID", js1.ID, &js2))

	originalInitiators := []migration0.Initiator{}
	for _, uc := range js1.Initiators.([]interface{}) {
		ti, err := migration1537223654.UnchangedToInitiator(uc.(migration0.Unchanged))
		assert.NoError(t, err)
		originalInitiators = append(originalInitiators, ti)
	}
	assert.Equal(t, originalInitiators[0].Schedule, js2.Initiators[0].Schedule)
}

func TestMigrate_updatesInitiatorsBucket(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	input := cltest.LoadJSON("../../../internal/fixtures/bolt/old_initiator_without_params.json")
	var i1 migration0.Initiator
	require.NoError(t, json.Unmarshal(input, &i1))
	require.NoError(t, store.Save(&i1))

	migration := migration1537223654.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var i2 migration1537223654.Initiator
	require.NoError(t, store.One("ID", i1.ID, &i2))

	assert.Equal(t, i1.ID, i2.ID)
	assert.Equal(t, i1.Schedule, i2.Schedule)
}

func TestMigrate_updatesJobRunsBucket(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	input := cltest.LoadJSON("../../../internal/fixtures/bolt/old_jobrun_without_initiator_params.json")
	var jr1 migration1536696950.JobRun
	require.NoError(t, json.Unmarshal(input, &jr1))
	require.NoError(t, store.Save(&jr1))

	migration := migration1537223654.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var jr2 migration1537223654.JobRun
	require.NoError(t, store.One("ID", jr1.ID, &jr2))

	assert.Equal(t, jr1.ID, jr2.ID)
	oi, err := migration1537223654.UnchangedToInitiator(jr1.Initiator.(migration0.Unchanged))
	assert.NoError(t, err)
	assert.Equal(t, oi.ID, jr2.Initiator.ID)
	assert.Equal(t, oi.JobID, jr2.Initiator.JobID)
	assert.Equal(t, oi.Schedule, jr2.Initiator.Schedule)
}
