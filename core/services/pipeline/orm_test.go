package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

type testOnlyORM interface {
	pipeline.ORM
	AddJobPipelineSpecWithoutConstraints(ctx context.Context, jobID, pipelineSpecID int32) error
}

type testORM struct {
	pipeline.ORM
	ds sqlutil.DataSource
}

func (torm *testORM) AddJobPipelineSpecWithoutConstraints(ctx context.Context, jobID, pipelineSpecID int32) error {
	_, err := torm.ds.ExecContext(ctx, `SET CONSTRAINTS fk_job_pipeline_spec_job DEFERRED`)
	if err != nil {
		return err
	}
	_, err = torm.ds.ExecContext(ctx, `INSERT INTO job_pipeline_specs (job_id,pipeline_spec_id, is_primary) VALUES ($1, $2, false)`, jobID, pipelineSpecID)
	if err != nil {
		return err
	}
	return nil
}

func newTestORM(orm pipeline.ORM, ds sqlutil.DataSource) testOnlyORM {
	return &testORM{ORM: orm, ds: ds}
}

func setupORM(t *testing.T, heavy bool) (db *sqlx.DB, orm pipeline.ORM, jorm job.ORM) {
	t.Helper()

	if heavy {
		_, db = heavyweight.FullTestDBV2(t, nil)
	} else {
		db = pgtest.NewSqlxDB(t)
	}
	orm = pipeline.NewORM(db, logger.TestLogger(t), 123456)
	lggr := logger.TestLogger(t)
	keyStore := cltest.NewKeyStore(t, db)
	bridgeORM := bridges.NewORM(db)

	jorm = job.NewORM(db, orm, bridgeORM, keyStore, lggr)

	return
}

func setupHeavyORM(t *testing.T) (db *sqlx.DB, orm pipeline.ORM, jorm job.ORM) {
	return setupORM(t, true)
}

func setupLiteORM(t *testing.T) (db *sqlx.DB, orm pipeline.ORM, jorm job.ORM) {
	return setupORM(t, false)
}

func Test_PipelineORM_CreateSpec(t *testing.T) {
	ctx := testutils.Context(t)
	db, orm, _ := setupLiteORM(t)

	var (
		source          = ""
		maxTaskDuration = models.Interval(1 * time.Minute)
	)

	p := pipeline.Pipeline{
		Source: source,
	}

	id, err := orm.CreateSpec(ctx, p, maxTaskDuration)
	require.NoError(t, err)

	actual := pipeline.Spec{}
	err = db.Get(&actual, "SELECT * FROM pipeline_specs WHERE pipeline_specs.id = $1", id)
	require.NoError(t, err)
	assert.Equal(t, source, actual.DotDagSource)
	assert.Equal(t, maxTaskDuration, actual.MaxTaskDuration)
}

func Test_PipelineORM_FindRun(t *testing.T) {
	db, orm, _ := setupLiteORM(t)

	_, err := db.Exec(`SET CONSTRAINTS fk_pipeline_runs_pruning_key DEFERRED`)
	require.NoError(t, err)
	_, err = db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)
	require.NoError(t, err)
	expected := mustInsertPipelineRun(t, orm)

	ctx := testutils.Context(t)
	run, err := orm.FindRun(ctx, expected.ID)
	require.NoError(t, err)

	require.Equal(t, expected.ID, run.ID)
}

func mustInsertPipelineRun(t *testing.T, orm pipeline.ORM) pipeline.Run {
	t.Helper()

	run := pipeline.Run{
		State:       pipeline.RunStatusRunning,
		Outputs:     jsonserializable.JSONSerializable{},
		AllErrors:   pipeline.RunErrors{},
		FatalErrors: pipeline.RunErrors{},
		FinishedAt:  null.Time{},
	}

	ctx := testutils.Context(t)
	require.NoError(t, orm.InsertRun(ctx, &run))
	return run
}

func mustInsertAsyncRun(t *testing.T, orm pipeline.ORM, jobORM job.ORM) *pipeline.Run {
	t.Helper()
	ctx := testutils.Context(t)

	s := `
ds1 [type=bridge async=true name="example-bridge" timeout=0 requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
ds1_parse [type=jsonparse lax=false  path="data,result"]
ds1_multiply [type=multiply times=1000000000000000000]

ds1->ds1_parse->ds1_multiply->answer1;

answer1 [type=median index=0];
answer2 [type=bridge name=election_winner index=1];
`
	jb := job.Job{
		Type:            job.DirectRequest,
		SchemaVersion:   1,
		MaxTaskDuration: models.Interval(1 * time.Minute),
		DirectRequestSpec: &job.DirectRequestSpec{
			ContractAddress: cltest.NewEIP55Address(),
			EVMChainID:      (*big.Big)(&cltest.FixtureChainID),
		},
		PipelineSpec: &pipeline.Spec{
			DotDagSource: s,
		},
	}
	err := jobORM.CreateJob(ctx, &jb)
	require.NoError(t, err)

	run := &pipeline.Run{
		PipelineSpecID: jb.PipelineSpecID,
		PruningKey:     jb.ID,
		State:          pipeline.RunStatusRunning,
		Outputs:        jsonserializable.JSONSerializable{},
		CreatedAt:      time.Now(),
	}

	err = orm.CreateRun(ctx, run)
	require.NoError(t, err)
	return run
}

func TestInsertFinishedRuns(t *testing.T) {
	ctx := testutils.Context(t)
	db, orm, _ := setupLiteORM(t)

	_, err := db.Exec(`SET CONSTRAINTS fk_pipeline_runs_pruning_key DEFERRED`)
	require.NoError(t, err)
	_, err = db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)
	require.NoError(t, err)

	ps := cltest.MustInsertPipelineSpec(t, db)

	var runs []*pipeline.Run
	for i := 0; i < 3; i++ {
		now := time.Now()
		r := pipeline.Run{
			PipelineSpecID: ps.ID,
			PruningKey:     ps.ID, // using the spec ID as the pruning key for test purposes, this is supposed to be the job ID
			State:          pipeline.RunStatusRunning,
			AllErrors:      pipeline.RunErrors{},
			FatalErrors:    pipeline.RunErrors{},
			CreatedAt:      now,
			FinishedAt:     null.Time{},
			Outputs:        jsonserializable.JSONSerializable{},
		}

		require.NoError(t, orm.InsertRun(ctx, &r))

		r.PipelineTaskRuns = []pipeline.TaskRun{
			{
				ID:            uuid.New(),
				PipelineRunID: r.ID,
				Type:          "bridge",
				DotID:         "ds1",
				CreatedAt:     now,
				FinishedAt:    null.TimeFrom(now.Add(100 * time.Millisecond)),
			},
			{
				ID:            uuid.New(),
				PipelineRunID: r.ID,
				Type:          "median",
				DotID:         "answer2",
				Output:        jsonserializable.JSONSerializable{Val: 1, Valid: true},
				CreatedAt:     now,
				FinishedAt:    null.TimeFrom(now.Add(200 * time.Millisecond)),
			},
		}
		r.FinishedAt = null.TimeFrom(now.Add(300 * time.Millisecond))
		r.Outputs = jsonserializable.JSONSerializable{
			Val:   "stuff",
			Valid: true,
		}
		r.AllErrors = append(r.AllErrors, null.NewString("", false))
		r.State = pipeline.RunStatusCompleted
		runs = append(runs, &r)
	}

	err = orm.InsertFinishedRuns(ctx, runs, true)
	require.NoError(t, err)
}

func Test_PipelineORM_InsertFinishedRunWithSpec(t *testing.T) {
	ctx := testutils.Context(t)
	db, orm, jorm := setupLiteORM(t)

	s := `
ds1 [type=bridge async=true name="example-bridge" timeout=0 requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
ds1_parse [type=jsonparse lax=false  path="data,result"]
ds1_multiply [type=multiply times=1000000000000000000]

ds1->ds1_parse->ds1_multiply->answer1;

answer1 [type=median index=0];
answer2 [type=bridge name=election_winner index=1];
`
	jb := job.Job{
		Type:            job.DirectRequest,
		SchemaVersion:   1,
		MaxTaskDuration: models.Interval(1 * time.Minute),
		DirectRequestSpec: &job.DirectRequestSpec{
			ContractAddress: cltest.NewEIP55Address(),
			EVMChainID:      (*big.Big)(&cltest.FixtureChainID),
		},
		PipelineSpec: &pipeline.Spec{
			DotDagSource: s,
		},
	}
	err := jorm.CreateJob(ctx, &jb)
	require.NoError(t, err)
	spec := pipeline.Spec{
		DotDagSource:    s,
		CreatedAt:       time.Now(),
		MaxTaskDuration: models.Interval(1 * time.Minute),
		JobID:           jb.ID,
		JobName:         jb.Name.ValueOrZero(),
		JobType:         string(jb.Type),
	}
	defaultVars := map[string]interface{}{
		"jb": map[string]interface{}{
			"databaseID":    jb.ID,
			"externalJobID": jb.ExternalJobID,
			"name":          jb.Name.ValueOrZero(),
		},
	}
	now := time.Now()
	run := pipeline.NewRun(spec, pipeline.NewVarsFrom(defaultVars))
	run.PipelineTaskRuns = []pipeline.TaskRun{
		{
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "bridge",
			DotID:         "ds1",
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now.Add(100 * time.Millisecond)),
		},
		{
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "median",
			DotID:         "answer2",
			Output:        jsonserializable.JSONSerializable{Val: 1, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now.Add(200 * time.Millisecond)),
		},
	}
	run.FinishedAt = null.TimeFrom(now.Add(300 * time.Millisecond))
	run.Outputs = jsonserializable.JSONSerializable{
		Val:   "stuff",
		Valid: true,
	}
	run.AllErrors = append(run.AllErrors, null.NewString("", false))
	run.State = pipeline.RunStatusCompleted

	err = orm.InsertFinishedRunWithSpec(ctx, run, true)
	require.NoError(t, err)

	var pipelineSpec pipeline.Spec
	err = db.Get(&pipelineSpec, "SELECT pipeline_specs.* FROM pipeline_specs JOIN job_pipeline_specs ON (pipeline_specs.id = job_pipeline_specs.pipeline_spec_id) WHERE job_pipeline_specs.job_id = $1 AND pipeline_specs.id = $2", jb.ID, run.PipelineSpecID)
	require.NoError(t, err)
	var jobPipelineSpec job.PipelineSpec
	err = db.Get(&jobPipelineSpec, "SELECT * FROM job_pipeline_specs WHERE job_id = $1 AND pipeline_spec_id = $2", jb.ID, pipelineSpec.ID)
	require.NoError(t, err)

	assert.Equal(t, run.PipelineSpecID, pipelineSpec.ID)
	assert.False(t, jobPipelineSpec.IsPrimary)
}

// Tests that inserting run results, then later updating the run results via upsert will work correctly.
func Test_PipelineORM_StoreRun_ShouldUpsert(t *testing.T) {
	ctx := testutils.Context(t)
	_, orm, jorm := setupLiteORM(t)

	run := mustInsertAsyncRun(t, orm, jorm)

	now := time.Now()

	run.PipelineTaskRuns = []pipeline.TaskRun{
		// pending task
		{
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "bridge",
			DotID:         "ds1",
			CreatedAt:     now,
			FinishedAt:    null.Time{},
		},
		// finished task
		{
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "median",
			DotID:         "answer2",
			Output:        jsonserializable.JSONSerializable{Val: 1, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	restart, err := orm.StoreRun(ctx, run)
	require.NoError(t, err)
	// no new data, so we don't need a restart
	require.Equal(t, false, restart)
	// the run is paused
	require.Equal(t, pipeline.RunStatusSuspended, run.State)

	r, err := orm.FindRun(ctx, run.ID)
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
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "bridge",
			DotID:         "ds1",
			Output:        jsonserializable.JSONSerializable{Val: 2, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	restart, err = orm.StoreRun(ctx, run)
	require.NoError(t, err)
	// no new data, so we don't need a restart
	require.Equal(t, false, restart)
	// the run is paused
	require.Equal(t, pipeline.RunStatusSuspended, run.State)

	r, err = orm.FindRun(ctx, run.ID)
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
	ctx := testutils.Context(t)
	db, orm, jorm := setupLiteORM(t)

	run := mustInsertAsyncRun(t, orm, jorm)

	r, err := orm.FindRun(ctx, run.ID)
	require.NoError(t, err)
	require.Equal(t, run.Inputs, r.Inputs)

	now := time.Now()

	ds1_id := uuid.New()

	// insert something for this pipeline_run to trigger an early resume while the pipeline is running
	rows, err := db.NamedQuery(`
	INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
	VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at, :finished_at)
	`, pipeline.TaskRun{
		ID:            ds1_id,
		PipelineRunID: run.ID,
		Type:          "bridge",
		DotID:         "ds1",
		Output:        jsonserializable.JSONSerializable{Val: 2, Valid: true},
		CreatedAt:     now,
		FinishedAt:    null.TimeFrom(now),
	})
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, rows.Close()) })

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
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "median",
			DotID:         "answer2",
			Output:        jsonserializable.JSONSerializable{Val: 1, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}

	restart, err := orm.StoreRun(ctx, run)
	require.NoError(t, err)
	// new data available! immediately restart the run
	require.Equal(t, true, restart)
	// the run is still in progress
	require.Equal(t, pipeline.RunStatusRunning, run.State)

	// confirm we now contain the latest restart data merged with local task data
	ds1 := run.ByDotID("ds1")
	require.Equal(t, ds1.Output.Val, int64(2))
	require.True(t, ds1.FinishedAt.Valid)
}

func Test_PipelineORM_StoreRun_UpdateTaskRunResult(t *testing.T) {
	ctx := testutils.Context(t)
	_, orm, jorm := setupLiteORM(t)

	run := mustInsertAsyncRun(t, orm, jorm)

	ds1_id := uuid.New()
	now := time.Now()
	address, err := hex.DecodeString("0x8bd112d3f8f92e41c861939545ad387307af9703")
	require.NoError(t, err)
	cborOutput := map[string]interface{}{
		"blockNum":        "0x13babbd",
		"confirmations":   int64(10),
		"contractAddress": address,
		"libraryVersion":  int64(1),
		"remoteChainId":   int64(106),
	}

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
		// finished task with json output
		{
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "cbor_parse",
			DotID:         "ds2",
			Output:        jsonserializable.JSONSerializable{Val: cborOutput, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
		// finished task
		{
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "median",
			DotID:         "answer2",
			Output:        jsonserializable.JSONSerializable{Val: 1, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	// assert that run should be in "running" state
	require.Equal(t, pipeline.RunStatusRunning, run.State)

	// Now store a partial run
	restart, err := orm.StoreRun(ctx, run)
	require.NoError(t, err)
	require.False(t, restart)
	// assert that run should be in "paused" state
	require.Equal(t, pipeline.RunStatusSuspended, run.State)

	r, start, err := orm.UpdateTaskRunResult(ctx, ds1_id, pipeline.Result{Value: "foo"})
	run = &r
	require.NoError(t, err)
	assert.Greater(t, run.ID, int64(0))
	assert.Greater(t, run.PipelineSpec.ID, int32(0)) // Make sure it actually loaded everything

	require.Len(t, run.PipelineTaskRuns, 3)
	// assert that run should be in "running" state
	require.Equal(t, pipeline.RunStatusRunning, run.State)
	// assert that we get the start signal
	require.True(t, start)

	// assert that the task is now updated
	task := run.ByDotID("ds1")
	require.True(t, task.FinishedAt.Valid)
	require.Equal(t, jsonserializable.JSONSerializable{Val: "foo", Valid: true}, task.Output)

	// assert correct task run serialization
	task2 := run.ByDotID("ds2")
	cborOutput["contractAddress"] = "0x8bd112d3f8f92e41c861939545ad387307af9703"
	require.Equal(t, jsonserializable.JSONSerializable{Val: cborOutput, Valid: true}, task2.Output)
}

func Test_PipelineORM_DeleteRun(t *testing.T) {
	ctx := testutils.Context(t)
	_, orm, jorm := setupLiteORM(t)

	run := mustInsertAsyncRun(t, orm, jorm)

	now := time.Now()

	run.PipelineTaskRuns = []pipeline.TaskRun{
		// pending task
		{
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "bridge",
			DotID:         "ds1",
			CreatedAt:     now,
			FinishedAt:    null.Time{},
		},
		// finished task
		{
			ID:            uuid.New(),
			PipelineRunID: run.ID,
			Type:          "median",
			DotID:         "answer2",
			Output:        jsonserializable.JSONSerializable{Val: 1, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	restart, err := orm.StoreRun(ctx, run)
	require.NoError(t, err)
	// no new data, so we don't need a restart
	require.Equal(t, false, restart)
	// the run is paused
	require.Equal(t, pipeline.RunStatusSuspended, run.State)

	err = orm.DeleteRun(ctx, run.ID)
	require.NoError(t, err)

	_, err = orm.FindRun(ctx, run.ID)
	require.Error(t, err, "not found")
}

func Test_PipelineORM_DeleteRunsOlderThan(t *testing.T) {
	ctx := testutils.Context(t)
	_, orm, jorm := setupHeavyORM(t)

	var runsIds []int64

	for i := 1; i <= 2000; i++ {
		run := mustInsertAsyncRun(t, orm, jorm)

		now := time.Now()

		run.PipelineTaskRuns = []pipeline.TaskRun{
			// finished task
			{
				ID:            uuid.New(),
				PipelineRunID: run.ID,
				Type:          "median",
				DotID:         "answer2",
				Output:        jsonserializable.JSONSerializable{Val: 1, Valid: true},
				CreatedAt:     now,
				FinishedAt:    null.TimeFrom(now.Add(-1 * time.Second)),
			},
		}
		run.State = pipeline.RunStatusCompleted
		run.FinishedAt = null.TimeFrom(now.Add(-1 * time.Second))
		run.Outputs = jsonserializable.JSONSerializable{Val: 1, Valid: true}
		run.AllErrors = pipeline.RunErrors{null.StringFrom("SOMETHING")}

		restart, err := orm.StoreRun(ctx, run)
		assert.NoError(t, err)
		// no new data, so we don't need a restart
		assert.Equal(t, false, restart)

		runsIds = append(runsIds, run.ID)
	}

	err := orm.DeleteRunsOlderThan(testutils.Context(t), 1*time.Second)
	assert.NoError(t, err)

	for _, runId := range runsIds {
		_, err := orm.FindRun(ctx, runId)
		require.Error(t, err, "not found")
	}
}

func Test_GetUnfinishedRuns_Keepers(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	// The test configures single Keeper job with two running tasks.
	// GetUnfinishedRuns() expects to catch both running tasks.

	config := configtest.NewTestGeneralConfig(t)
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	porm := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgeORM := bridges.NewORM(db)

	jorm := job.NewORM(db, porm, bridgeORM, keyStore, lggr)
	defer func() { assert.NoError(t, jorm.Close()) }()

	timestamp := time.Now()
	var keeperJob = job.Job{
		ID: 1,
		KeeperSpec: &job.KeeperSpec{
			ContractAddress: cltest.NewEIP55Address(),
			FromAddress:     cltest.NewEIP55Address(),
			CreatedAt:       timestamp,
			UpdatedAt:       timestamp,
			EVMChainID:      (*big.Big)(&cltest.FixtureChainID),
		},
		ExternalJobID: uuid.MustParse("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
		PipelineSpec: &pipeline.Spec{
			ID:           1,
			DotDagSource: "",
		},
		Type:            job.Keeper,
		SchemaVersion:   1,
		Name:            null.StringFrom("test"),
		MaxTaskDuration: models.Interval(1 * time.Minute),
	}

	err := jorm.CreateJob(ctx, &keeperJob)
	require.NoError(t, err)
	require.Equal(t, job.Keeper, keeperJob.Type)

	runID1 := uuid.New()
	runID2 := uuid.New()

	err = porm.CreateRun(ctx, &pipeline.Run{
		PipelineSpecID: keeperJob.PipelineSpecID,
		PruningKey:     keeperJob.ID,
		State:          pipeline.RunStatusRunning,
		Outputs:        jsonserializable.JSONSerializable{},
		CreatedAt:      time.Now(),
		PipelineTaskRuns: []pipeline.TaskRun{{
			ID:        runID1,
			Type:      pipeline.TaskTypeETHTx,
			Index:     0,
			Output:    jsonserializable.JSONSerializable{},
			CreatedAt: time.Now(),
			DotID:     "perform_upkeep_tx",
		}},
	})
	require.NoError(t, err)

	err = porm.CreateRun(ctx, &pipeline.Run{
		PipelineSpecID: keeperJob.PipelineSpecID,
		PruningKey:     keeperJob.ID,
		State:          pipeline.RunStatusRunning,
		Outputs:        jsonserializable.JSONSerializable{},
		CreatedAt:      time.Now(),
		PipelineTaskRuns: []pipeline.TaskRun{{
			ID:        runID2,
			Type:      pipeline.TaskTypeETHCall,
			Index:     1,
			Output:    jsonserializable.JSONSerializable{},
			CreatedAt: time.Now(),
			DotID:     "check_upkeep_tx",
		}},
	})
	require.NoError(t, err)

	var counter int

	err = porm.GetUnfinishedRuns(testutils.Context(t), time.Now(), func(run pipeline.Run) error {
		counter++

		require.Equal(t, job.Keeper.String(), run.PipelineSpec.JobType)
		require.Equal(t, pipeline.KeepersObservationSource, run.PipelineSpec.DotDagSource)
		require.NotEmpty(t, run.PipelineTaskRuns)

		switch run.PipelineTaskRuns[0].ID {
		case runID1:
			trun := run.ByDotID("perform_upkeep_tx")
			require.NotNil(t, trun)
		case runID2:
			trun := run.ByDotID("check_upkeep_tx")
			require.NotNil(t, trun)
		}

		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 2, counter)
}

func Test_GetUnfinishedRuns_DirectRequest(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	// The test configures single DR job with two task runs: one is running and one is suspended.
	// GetUnfinishedRuns() expects to catch the one that is running.

	config := configtest.NewTestGeneralConfig(t)
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	porm := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgeORM := bridges.NewORM(db)

	jorm := job.NewORM(db, porm, bridgeORM, keyStore, lggr)
	defer func() { assert.NoError(t, jorm.Close()) }()

	timestamp := time.Now()
	var drJob = job.Job{
		ID: 1,
		DirectRequestSpec: &job.DirectRequestSpec{
			ContractAddress: cltest.NewEIP55Address(),
			CreatedAt:       timestamp,
			UpdatedAt:       timestamp,
			EVMChainID:      (*big.Big)(&cltest.FixtureChainID),
		},
		ExternalJobID: uuid.MustParse("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
		PipelineSpec: &pipeline.Spec{
			ID:           1,
			DotDagSource: `ds1 [type=http method=GET url="https://pricesource1.com"`,
		},
		Type:            job.DirectRequest,
		SchemaVersion:   1,
		Name:            null.StringFrom("test"),
		MaxTaskDuration: models.Interval(1 * time.Minute),
	}

	err := jorm.CreateJob(ctx, &drJob)
	require.NoError(t, err)
	require.Equal(t, job.DirectRequest, drJob.Type)

	runningID := uuid.New()

	err = porm.CreateRun(ctx, &pipeline.Run{
		PipelineSpecID: drJob.PipelineSpecID,
		PruningKey:     drJob.ID,
		State:          pipeline.RunStatusRunning,
		Outputs:        jsonserializable.JSONSerializable{},
		CreatedAt:      time.Now(),
		PipelineTaskRuns: []pipeline.TaskRun{{
			ID:        runningID,
			Type:      pipeline.TaskTypeHTTP,
			Index:     0,
			Output:    jsonserializable.JSONSerializable{},
			CreatedAt: time.Now(),
			DotID:     "ds1",
		}},
	})
	require.NoError(t, err)

	err = porm.CreateRun(ctx, &pipeline.Run{
		PipelineSpecID: drJob.PipelineSpecID,
		PruningKey:     drJob.ID,
		State:          pipeline.RunStatusSuspended,
		Outputs:        jsonserializable.JSONSerializable{},
		CreatedAt:      time.Now(),
		PipelineTaskRuns: []pipeline.TaskRun{{
			ID:        uuid.New(),
			Type:      pipeline.TaskTypeHTTP,
			Index:     1,
			Output:    jsonserializable.JSONSerializable{},
			CreatedAt: time.Now(),
			DotID:     "ds1",
		}},
	})
	require.NoError(t, err)

	var counter int

	err = porm.GetUnfinishedRuns(testutils.Context(t), time.Now(), func(run pipeline.Run) error {
		counter++

		require.Equal(t, job.DirectRequest.String(), run.PipelineSpec.JobType)
		require.NotEmpty(t, run.PipelineTaskRuns)
		require.Equal(t, runningID, run.PipelineTaskRuns[0].ID)

		trun := run.ByDotID("ds1")
		require.NotNil(t, trun)

		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 1, counter)
}

func Test_Prune(t *testing.T) {
	t.Parallel()

	n := uint64(2)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.MaxSuccessfulRuns = &n
	})
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	db := pgtest.NewSqlxDB(t)
	porm := pipeline.NewORM(db, lggr, cfg.JobPipeline().MaxSuccessfulRuns())
	torm := newTestORM(porm, db)

	ps1 := cltest.MustInsertPipelineSpec(t, db)

	// We need a job_pipeline_specs entry to test the pruning mechanism
	err := torm.AddJobPipelineSpecWithoutConstraints(testutils.Context(t), ps1.ID, ps1.ID)
	require.NoError(t, err)

	jobID := ps1.ID

	t.Run("when there are no runs to prune, does nothing", func(t *testing.T) {
		ctx := tests.Context(t)
		porm.Prune(ctx, jobID)

		// no error logs; it did nothing
		assert.Empty(t, observed.All())
	})

	_, err = db.Exec(`SET CONSTRAINTS fk_pipeline_runs_pruning_key DEFERRED`)
	require.NoError(t, err)

	// ps1 has:
	// - 20 completed runs
	for i := 0; i < 20; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps1.ID, pipeline.RunStatusCompleted, jobID)
	}

	ps2 := cltest.MustInsertPipelineSpec(t, db)

	jobID2 := ps2.ID
	// ps2 has:
	// - 12 completed runs
	// - 3 errored runs
	// - 3 running runs
	// - 3 suspended run
	for i := 0; i < 12; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps2.ID, pipeline.RunStatusCompleted, jobID2)
	}
	for i := 0; i < 3; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps2.ID, pipeline.RunStatusErrored, jobID2)
	}
	for i := 0; i < 3; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps2.ID, pipeline.RunStatusRunning, jobID2)
	}
	for i := 0; i < 3; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps2.ID, pipeline.RunStatusSuspended, jobID2)
	}

	porm.Prune(tests.Context(t), jobID2)

	cnt := pgtest.MustCount(t, db, "SELECT count(*) FROM pipeline_runs WHERE pipeline_spec_id = $1 AND state = $2", ps1.ID, pipeline.RunStatusCompleted)
	assert.Equal(t, cnt, 20)

	cnt = pgtest.MustCount(t, db, "SELECT count(*) FROM pipeline_runs WHERE pipeline_spec_id = $1 AND state = $2", ps2.ID, pipeline.RunStatusCompleted)
	assert.Equal(t, 2, cnt)
	cnt = pgtest.MustCount(t, db, "SELECT count(*) FROM pipeline_runs WHERE pipeline_spec_id = $1 AND state = $2", ps2.ID, pipeline.RunStatusErrored)
	assert.Equal(t, 3, cnt)
	cnt = pgtest.MustCount(t, db, "SELECT count(*) FROM pipeline_runs WHERE pipeline_spec_id = $1 AND state = $2", ps2.ID, pipeline.RunStatusRunning)
	assert.Equal(t, 3, cnt)
	cnt = pgtest.MustCount(t, db, "SELECT count(*) FROM pipeline_runs WHERE pipeline_spec_id = $1 AND state = $2", ps2.ID, pipeline.RunStatusSuspended)
	assert.Equal(t, 3, cnt)
}
