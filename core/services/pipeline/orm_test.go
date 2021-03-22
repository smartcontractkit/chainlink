package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres/mocks"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func Test_PipelineORM_CreateRun(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB

	eventBroadcaster := new(mocks.EventBroadcaster)
	orm := pipeline.NewORM(db, store.Config, eventBroadcaster)

	job := cltest.MustInsertSampleDirectRequestJob(t, db)
	meta := make(map[string]interface{})

	runID, err := orm.CreateRun(context.Background(), job.ID, meta)
	require.NoError(t, err)

	// Check that JobRun, TaskRuns were created

	var prs []pipeline.Run
	var trs []pipeline.TaskRun

	require.NoError(t, db.Find(&prs).Error)
	require.NoError(t, db.Find(&trs).Error)

	require.Len(t, prs, 1)
	require.Equal(t, runID, prs[0].ID)
	require.Len(t, trs, 4)
}

func Test_PipelineORM_UpdatePipelineRun(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB

	require.NoError(t, db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`).Error)

	eventBroadcaster := new(mocks.EventBroadcaster)
	orm := pipeline.NewORM(db, store.Config, eventBroadcaster)

	t.Run("saves errored run with string error correctly", func(t *testing.T) {
		run := cltest.MustInsertPipelineRun(t, db)
		trrs := pipeline.TaskRunResults{
			pipeline.TaskRunResult{
				IsTerminal: true,
				Result: pipeline.Result{
					Value: []interface{}{nil},
					Error: pipeline.FinalErrors{null.StringFrom("Random: String, foo")},
				},
				FinishedAt: time.Now(),
			},
		}

		err := orm.UpdatePipelineRun(db, &run, trrs.FinalResult())
		require.NoError(t, err)

		require.Equal(t, []interface{}{nil}, run.Outputs.Val)
		require.Equal(t, []interface{}{"Random: String, foo"}, run.Errors.Val)
		require.NotNil(t, run.FinishedAt)
	})

	t.Run("saves errored run with final errors correctly", func(t *testing.T) {
		run := cltest.MustInsertPipelineRun(t, db)
		trrs := pipeline.TaskRunResults{
			pipeline.TaskRunResult{
				IsTerminal: true,
				Result: pipeline.Result{
					Value: []interface{}{1, nil},
					Error: pipeline.FinalErrors([]null.String{
						null.String{},
						null.StringFrom(`Random: String, foo`),
					}),
				},
				FinishedAt: time.Now(),
			},
		}

		err := orm.UpdatePipelineRun(db, &run, trrs.FinalResult())
		require.NoError(t, err)

		require.Equal(t, []interface{}{float64(1), nil}, run.Outputs.Val)
		require.Equal(t, []interface{}{nil, "Random: String, foo"}, run.Errors.Val)
		require.NotNil(t, run.FinishedAt)
	})
}

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
