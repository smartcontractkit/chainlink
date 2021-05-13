package pipeline_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres/mocks"
	"github.com/stretchr/testify/require"
)

func Test_PipelineORM_FindRun(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB

	eventBroadcaster := new(mocks.EventBroadcaster)
	orm := pipeline.NewORM(db, store.Config, eventBroadcaster)

	require.NoError(t, db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`).Error)
	expected := cltest.MustInsertPipelineRun(t, db)

	run, err := orm.FindRun(expected.ID)
	require.NoError(t, err)

	require.Equal(t, expected.ID, run.ID)
}
