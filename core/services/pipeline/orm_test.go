package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func Test_PipelineORM_CreateSpec(t *testing.T) {
	db, orm := setupORM(t)

	var (
		source          = ""
		maxTaskDuration = models.Interval(1 * time.Minute)
	)

	p := pipeline.Pipeline{
		Source: source,
	}

	id, err := orm.CreateSpec(context.Background(), db, p, maxTaskDuration)
	require.NoError(t, err)

	actual := pipeline.Spec{}
	err = db.Find(&actual, id).Error
	require.NoError(t, err)
	assert.Equal(t, source, actual.DotDagSource)
	assert.Equal(t, maxTaskDuration, actual.MaxTaskDuration)
}

func Test_PipelineORM_FindRun(t *testing.T) {
	db, orm := setupORM(t)

	require.NoError(t, db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`).Error)
	expected := mustInsertPipelineRun(t, db)

	run, err := orm.FindRun(expected.ID)
	require.NoError(t, err)

	require.Equal(t, expected.ID, run.ID)
}

func mustInsertPipelineRun(t *testing.T, db *gorm.DB) pipeline.Run {
	t.Helper()

	run := pipeline.Run{
		State:      pipeline.RunStatusRunning,
		Outputs:    pipeline.JSONSerializable{Null: true},
		Errors:     pipeline.RunErrors{},
		FinishedAt: null.Time{},
	}
	require.NoError(t, db.Create(&run).Error)
	return run
}

func setupORM(t *testing.T) (*gorm.DB, pipeline.ORM) {
	t.Helper()

	db := pgtest.NewGormDB(t)
	orm := pipeline.NewORM(db)

	return db, orm
}

func mustInsertAsyncRun(t *testing.T, orm pipeline.ORM, db *gorm.DB) *pipeline.Run {
	t.Helper()

	s := `
ds1 [type=bridge async=true name="example-bridge" timeout=0 requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
ds1_parse [type=jsonparse lax=false  path="data,result"]
ds1_multiply [type=multiply times=1000000000000000000]

ds1->ds1_parse->ds1_multiply->answer1;

answer1 [type=median index=0];
answer2 [type=bridge name=election_winner index=1];
`

	p, err := pipeline.Parse(s)
	require.NoError(t, err)
	require.NotNil(t, p)

	maxTaskDuration := models.Interval(1 * time.Minute)
	specID, err := orm.CreateSpec(context.Background(), db, *p, maxTaskDuration)
	require.NoError(t, err)

	run := &pipeline.Run{
		PipelineSpecID: specID,
		State:          pipeline.RunStatusRunning,
		Errors:         nil,
		Outputs:        pipeline.JSONSerializable{Null: true},
		CreatedAt:      time.Now(),
	}

	err = orm.CreateRun(db, run)
	require.NoError(t, err)
	return run
}

// Tests that inserting run results, then later updating the run results via upsert will work correctly.
func Test_PipelineORM_StoreRun_ShouldUpsert(t *testing.T) {
	db, orm := setupORM(t)

	run := mustInsertAsyncRun(t, orm, db)

	sdb, err := orm.DB().DB()
	require.NoError(t, err)

	now := time.Now()

	run.PipelineTaskRuns = []pipeline.TaskRun{
		// pending task
		{
			ID:            uuid.NewV4(),
			PipelineRunID: run.ID,
			Type:          "bridge",
			DotID:         "ds1",
			CreatedAt:     now,
			FinishedAt:    null.Time{},
		},
		// finished task
		{
			ID:            uuid.NewV4(),
			PipelineRunID: run.ID,
			Type:          "median",
			DotID:         "answer2",
			Output:        &pipeline.JSONSerializable{Val: 1},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	restart, err := orm.StoreRun(sdb, run)
	require.NoError(t, err)
	// no new data, so we don't need a restart
	require.Equal(t, false, restart)
	// the run is paused
	require.Equal(t, pipeline.RunStatusSuspended, run.State)

	r, err := orm.FindRun(run.ID)
	require.NoError(t, err)
	run = &r
	// this is an incomplete run, so partial results should be present (regardless of saveSuccessfulTaskRuns)
	require.Equal(t, 2, len(run.PipelineTaskRuns))
	// and ds1 is not finished
	task := run.ByDotID("ds1")
	require.NotNil(t, task)
	require.False(t, task.FinishedAt.Valid)

	// now try setting the ds1 result: call store run again

	run.PipelineTaskRuns = []pipeline.TaskRun{
		// pending task
		{
			ID:            uuid.NewV4(),
			PipelineRunID: run.ID,
			Type:          "bridge",
			DotID:         "ds1",
			Output:        &pipeline.JSONSerializable{Val: 2},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	restart, err = orm.StoreRun(sdb, run)
	require.NoError(t, err)
	// no new data, so we don't need a restart
	require.Equal(t, false, restart)
	// the run is paused
	require.Equal(t, pipeline.RunStatusSuspended, run.State)

	r, err = orm.FindRun(run.ID)
	require.NoError(t, err)
	run = &r
	// this is an incomplete run, so partial results should be present (regardless of saveSuccessfulTaskRuns)
	require.Equal(t, 2, len(run.PipelineTaskRuns))
	// and ds1 is finished
	task = run.ByDotID("ds1")
	require.NotNil(t, task)
	require.NotNil(t, task.FinishedAt)
}

// Tests that trying to persist a partial run while new data became available (i.e. via /v2/restart)
// will detect a restart and update the result data on the Run.
func Test_PipelineORM_StoreRun_DetectsRestarts(t *testing.T) {
	db, orm := setupORM(t)

	run := mustInsertAsyncRun(t, orm, db)

	r, err := orm.FindRun(run.ID)
	require.NoError(t, err)
	require.Equal(t, run.Inputs, r.Inputs)

	now := time.Now()

	sdb, err := orm.DB().DB()
	require.NoError(t, err)

	ds1_id := uuid.NewV4()

	sqlxDb := postgres.WrapDbWithSqlx(sdb)

	// insert something for this pipeline_run to trigger an early resume while the pipeline is running
	_, err = sqlxDb.NamedQuery(`
	INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
	VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at, :finished_at)
	`, pipeline.TaskRun{
		ID:            ds1_id,
		PipelineRunID: run.ID,
		Type:          "bridge",
		DotID:         "ds1",
		Output:        &pipeline.JSONSerializable{Val: 2},
		CreatedAt:     now,
		FinishedAt:    null.TimeFrom(now),
	})
	require.NoError(t, err)

	run.PipelineTaskRuns = []pipeline.TaskRun{
		// pending task
		{
			ID:            ds1_id,
			PipelineRunID: run.ID,
			Type:          "bridge",
			DotID:         "ds1",
			CreatedAt:     now,
			FinishedAt:    null.Time{},
		},
		// finished task
		{
			ID:            uuid.NewV4(),
			PipelineRunID: run.ID,
			Type:          "median",
			DotID:         "answer2",
			Output:        &pipeline.JSONSerializable{Val: 1},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}

	restart, err := orm.StoreRun(sdb, run)
	require.NoError(t, err)
	// new data available! immediately restart the run
	require.Equal(t, true, restart)
	// the run is still in progress
	require.Equal(t, pipeline.RunStatusRunning, run.State)

	// confirm we now contain the latest restart data merged with local task data
	ds1 := run.ByDotID("ds1")
	require.Equal(t, ds1.Output.Val, float64(2))
	require.True(t, ds1.FinishedAt.Valid)

}

func Test_PipelineORM_StoreRun_UpdateTaskRunResult(t *testing.T) {
	db, orm := setupORM(t)

	run := mustInsertAsyncRun(t, orm, db)

	sdb, err := orm.DB().DB()
	require.NoError(t, err)

	now := time.Now()

	ds1_id := uuid.NewV4()
	run.PipelineTaskRuns = []pipeline.TaskRun{
		// pending task
		{
			ID:            ds1_id,
			PipelineRunID: run.ID,
			Type:          "bridge",
			DotID:         "ds1",
			CreatedAt:     now,
			FinishedAt:    null.Time{},
		},
		// finished task
		{
			ID:            uuid.NewV4(),
			PipelineRunID: run.ID,
			Type:          "median",
			DotID:         "answer2",
			Output:        &pipeline.JSONSerializable{Val: 1},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	// assert that run should be in "running" state
	require.Equal(t, pipeline.RunStatusRunning, run.State)

	// Now store a partial run
	restart, err := orm.StoreRun(sdb, run)
	require.NoError(t, err)
	require.False(t, restart)
	// assert that run should be in "paused" state
	require.Equal(t, pipeline.RunStatusSuspended, run.State)

	r, start, err := orm.UpdateTaskRunResult(sdb, ds1_id, "foo")
	run = &r
	require.NoError(t, err)
	require.Len(t, run.PipelineTaskRuns, 2)
	// assert that run should be in "running" state
	require.Equal(t, pipeline.RunStatusRunning, run.State)
	// assert that we get the start signal
	require.True(t, start)

	// assert that the task is now updated
	task := run.ByDotID("ds1")
	require.True(t, task.FinishedAt.Valid)
	require.Equal(t, &pipeline.JSONSerializable{Val: "foo"}, task.Output)
}
