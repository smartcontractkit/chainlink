package pipeline_test

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest2 "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type ormconfig struct {
	pg.QConfig
}

func (ormconfig) JobPipelineMaxSuccessfulRuns() uint64 { return 123456 }

func setupORM(t *testing.T, name string) (db *sqlx.DB, orm pipeline.ORM) {
	t.Helper()

	if name != "" {
		_, db = heavyweight.FullTestDBV2(t, name, nil)
	} else {
		db = pgtest.NewSqlxDB(t)
	}
	cfg := ormconfig{pgtest.NewQConfig(true)}
	orm = pipeline.NewORM(db, logger.TestLogger(t), cfg)

	return
}

func setupHeavyORM(t *testing.T, name string) (db *sqlx.DB, orm pipeline.ORM) {
	return setupORM(t, name)
}

func setupLiteORM(t *testing.T) (db *sqlx.DB, orm pipeline.ORM) {
	return setupORM(t, "")
}

func Test_PipelineORM_CreateSpec(t *testing.T) {
	db, orm := setupLiteORM(t)

	var (
		source          = ""
		maxTaskDuration = models.Interval(1 * time.Minute)
	)

	p := pipeline.Pipeline{
		Source: source,
	}

	id, err := orm.CreateSpec(p, maxTaskDuration)
	require.NoError(t, err)

	actual := pipeline.Spec{}
	err = db.Get(&actual, "SELECT * FROM pipeline_specs WHERE pipeline_specs.id = $1", id)
	require.NoError(t, err)
	assert.Equal(t, source, actual.DotDagSource)
	assert.Equal(t, maxTaskDuration, actual.MaxTaskDuration)
}

func Test_PipelineORM_FindRun(t *testing.T) {
	db, orm := setupLiteORM(t)

	_, err := db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)
	require.NoError(t, err)
	expected := mustInsertPipelineRun(t, orm)

	run, err := orm.FindRun(expected.ID)
	require.NoError(t, err)

	require.Equal(t, expected.ID, run.ID)
}

func mustInsertPipelineRun(t *testing.T, orm pipeline.ORM) pipeline.Run {
	t.Helper()

	run := pipeline.Run{
		State:       pipeline.RunStatusRunning,
		Outputs:     pipeline.JSONSerializable{},
		AllErrors:   pipeline.RunErrors{},
		FatalErrors: pipeline.RunErrors{},
		FinishedAt:  null.Time{},
	}

	require.NoError(t, orm.InsertRun(&run))
	return run
}

func mustInsertAsyncRun(t *testing.T, orm pipeline.ORM) *pipeline.Run {
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
	specID, err := orm.CreateSpec(*p, maxTaskDuration)
	require.NoError(t, err)

	run := &pipeline.Run{
		PipelineSpecID: specID,
		State:          pipeline.RunStatusRunning,
		Outputs:        pipeline.JSONSerializable{},
		CreatedAt:      time.Now(),
	}

	err = orm.CreateRun(run)
	require.NoError(t, err)
	return run
}

func TestInsertFinishedRuns(t *testing.T) {
	db, orm := setupLiteORM(t)

	_, err := db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)
	require.NoError(t, err)

	ps := cltest.MustInsertPipelineSpec(t, db)

	var runs []*pipeline.Run
	for i := 0; i < 3; i++ {
		now := time.Now()
		r := pipeline.Run{
			PipelineSpecID: ps.ID,
			State:          pipeline.RunStatusRunning,
			AllErrors:      pipeline.RunErrors{},
			FatalErrors:    pipeline.RunErrors{},
			CreatedAt:      now,
			FinishedAt:     null.Time{},
			Outputs:        pipeline.JSONSerializable{},
		}

		require.NoError(t, orm.InsertRun(&r))

		r.PipelineTaskRuns = []pipeline.TaskRun{
			{
				ID:            uuid.NewV4(),
				PipelineRunID: r.ID,
				Type:          "bridge",
				DotID:         "ds1",
				CreatedAt:     now,
				FinishedAt:    null.TimeFrom(now.Add(100 * time.Millisecond)),
			},
			{
				ID:            uuid.NewV4(),
				PipelineRunID: r.ID,
				Type:          "median",
				DotID:         "answer2",
				Output:        pipeline.JSONSerializable{Val: 1, Valid: true},
				CreatedAt:     now,
				FinishedAt:    null.TimeFrom(now.Add(200 * time.Millisecond)),
			},
		}
		r.FinishedAt = null.TimeFrom(now.Add(300 * time.Millisecond))
		r.Outputs = pipeline.JSONSerializable{
			Val:   "stuff",
			Valid: true,
		}
		r.AllErrors = append(r.AllErrors, null.NewString("", false))
		r.State = pipeline.RunStatusCompleted
		runs = append(runs, &r)
	}

	err = orm.InsertFinishedRuns(runs, true)
	require.NoError(t, err)

}

// Tests that inserting run results, then later updating the run results via upsert will work correctly.
func Test_PipelineORM_StoreRun_ShouldUpsert(t *testing.T) {
	_, orm := setupLiteORM(t)

	run := mustInsertAsyncRun(t, orm)

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
			Output:        pipeline.JSONSerializable{Val: 1, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	restart, err := orm.StoreRun(run)
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
			Output:        pipeline.JSONSerializable{Val: 2, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	restart, err = orm.StoreRun(run)
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
	db, orm := setupLiteORM(t)

	run := mustInsertAsyncRun(t, orm)

	r, err := orm.FindRun(run.ID)
	require.NoError(t, err)
	require.Equal(t, run.Inputs, r.Inputs)

	now := time.Now()

	ds1_id := uuid.NewV4()

	// insert something for this pipeline_run to trigger an early resume while the pipeline is running
	_, err = db.NamedQuery(`
	INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
	VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at, :finished_at)
	`, pipeline.TaskRun{
		ID:            ds1_id,
		PipelineRunID: run.ID,
		Type:          "bridge",
		DotID:         "ds1",
		Output:        pipeline.JSONSerializable{Val: 2, Valid: true},
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
			Output:        pipeline.JSONSerializable{Val: 1, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}

	restart, err := orm.StoreRun(run)
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
	_, orm := setupLiteORM(t)

	run := mustInsertAsyncRun(t, orm)

	ds1_id := uuid.NewV4()
	now := time.Now()
	address, err := utils.TryParseHex("0x8bd112d3f8f92e41c861939545ad387307af9703")
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
			ID:            uuid.NewV4(),
			PipelineRunID: run.ID,
			Type:          "cbor_parse",
			DotID:         "ds2",
			Output:        pipeline.JSONSerializable{Val: cborOutput, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
		// finished task
		{
			ID:            uuid.NewV4(),
			PipelineRunID: run.ID,
			Type:          "median",
			DotID:         "answer2",
			Output:        pipeline.JSONSerializable{Val: 1, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	// assert that run should be in "running" state
	require.Equal(t, pipeline.RunStatusRunning, run.State)

	// Now store a partial run
	restart, err := orm.StoreRun(run)
	require.NoError(t, err)
	require.False(t, restart)
	// assert that run should be in "paused" state
	require.Equal(t, pipeline.RunStatusSuspended, run.State)

	r, start, err := orm.UpdateTaskRunResult(ds1_id, pipeline.Result{Value: "foo"})
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
	require.Equal(t, pipeline.JSONSerializable{Val: "foo", Valid: true}, task.Output)

	// assert correct task run serialization
	task2 := run.ByDotID("ds2")
	cborOutput["contractAddress"] = "0x8bd112d3f8f92e41c861939545ad387307af9703"
	require.Equal(t, pipeline.JSONSerializable{Val: cborOutput, Valid: true}, task2.Output)
}

func Test_PipelineORM_DeleteRun(t *testing.T) {
	_, orm := setupLiteORM(t)

	run := mustInsertAsyncRun(t, orm)

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
			Output:        pipeline.JSONSerializable{Val: 1, Valid: true},
			CreatedAt:     now,
			FinishedAt:    null.TimeFrom(now),
		},
	}
	restart, err := orm.StoreRun(run)
	require.NoError(t, err)
	// no new data, so we don't need a restart
	require.Equal(t, false, restart)
	// the run is paused
	require.Equal(t, pipeline.RunStatusSuspended, run.State)

	err = orm.DeleteRun(run.ID)
	require.NoError(t, err)

	_, err = orm.FindRun(run.ID)
	require.Error(t, err, "not found")
}

func Test_PipelineORM_DeleteRunsOlderThan(t *testing.T) {
	_, orm := setupHeavyORM(t, "pipeline_runs_reaper")

	var runsIds []int64

	for i := 1; i <= 2000; i++ {
		run := mustInsertAsyncRun(t, orm)

		now := time.Now()

		run.PipelineTaskRuns = []pipeline.TaskRun{
			// finished task
			{
				ID:            uuid.NewV4(),
				PipelineRunID: run.ID,
				Type:          "median",
				DotID:         "answer2",
				Output:        pipeline.JSONSerializable{Val: 1, Valid: true},
				CreatedAt:     now,
				FinishedAt:    null.TimeFrom(now.Add(-1 * time.Second)),
			},
		}
		run.State = pipeline.RunStatusCompleted
		run.FinishedAt = null.TimeFrom(now.Add(-1 * time.Second))
		run.Outputs = pipeline.JSONSerializable{Val: 1, Valid: true}
		run.AllErrors = pipeline.RunErrors{null.StringFrom("SOMETHING")}

		restart, err := orm.StoreRun(run)
		assert.NoError(t, err)
		// no new data, so we don't need a restart
		assert.Equal(t, false, restart)

		runsIds = append(runsIds, run.ID)
	}

	err := orm.DeleteRunsOlderThan(testutils.Context(t), 1*time.Second)
	assert.NoError(t, err)

	for _, runId := range runsIds {
		_, err := orm.FindRun(runId)
		require.Error(t, err, "not found")
	}
}

func Test_GetUnfinishedRuns_Keepers(t *testing.T) {
	t.Parallel()

	// The test configures single Keeper job with two running tasks.
	// GetUnfinishedRuns() expects to catch both running tasks.

	config := configtest2.NewTestGeneralConfig(t)
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	porm := pipeline.NewORM(db, lggr, config)
	bridgeORM := bridges.NewORM(db, lggr, config)

	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	jorm := job.NewORM(db, cc, porm, bridgeORM, keyStore, lggr, config)
	defer func() { assert.NoError(t, jorm.Close()) }()

	timestamp := time.Now()
	var keeperJob = job.Job{
		ID: 1,
		KeeperSpec: &job.KeeperSpec{
			ContractAddress: cltest.NewEIP55Address(),
			FromAddress:     cltest.NewEIP55Address(),
			CreatedAt:       timestamp,
			UpdatedAt:       timestamp,
			EVMChainID:      (*utils.Big)(&cltest.FixtureChainID),
		},
		ExternalJobID: uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
		PipelineSpec: &pipeline.Spec{
			ID:           1,
			DotDagSource: "",
		},
		Type:            job.Keeper,
		SchemaVersion:   1,
		Name:            null.StringFrom("test"),
		MaxTaskDuration: models.Interval(1 * time.Minute),
	}

	err := jorm.CreateJob(&keeperJob)
	require.NoError(t, err)
	require.Equal(t, job.Keeper, keeperJob.Type)

	runID1 := uuid.NewV4()
	runID2 := uuid.NewV4()

	err = porm.CreateRun(&pipeline.Run{
		PipelineSpecID: keeperJob.PipelineSpecID,
		State:          pipeline.RunStatusRunning,
		Outputs:        pipeline.JSONSerializable{},
		CreatedAt:      time.Now(),
		PipelineTaskRuns: []pipeline.TaskRun{{
			ID:        runID1,
			Type:      pipeline.TaskTypeETHTx,
			Index:     0,
			Output:    pipeline.JSONSerializable{},
			CreatedAt: time.Now(),
			DotID:     "perform_upkeep_tx",
		}},
	})
	require.NoError(t, err)

	err = porm.CreateRun(&pipeline.Run{
		PipelineSpecID: keeperJob.PipelineSpecID,
		State:          pipeline.RunStatusRunning,
		Outputs:        pipeline.JSONSerializable{},
		CreatedAt:      time.Now(),
		PipelineTaskRuns: []pipeline.TaskRun{{
			ID:        runID2,
			Type:      pipeline.TaskTypeETHCall,
			Index:     1,
			Output:    pipeline.JSONSerializable{},
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

	// The test configures single DR job with two task runs: one is running and one is suspended.
	// GetUnfinishedRuns() expects to catch the one that is running.

	config := configtest2.NewTestGeneralConfig(t)
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	porm := pipeline.NewORM(db, lggr, config)
	bridgeORM := bridges.NewORM(db, lggr, config)

	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	jorm := job.NewORM(db, cc, porm, bridgeORM, keyStore, lggr, config)
	defer func() { assert.NoError(t, jorm.Close()) }()

	timestamp := time.Now()
	var drJob = job.Job{
		ID: 1,
		DirectRequestSpec: &job.DirectRequestSpec{
			ContractAddress: cltest.NewEIP55Address(),
			CreatedAt:       timestamp,
			UpdatedAt:       timestamp,
			EVMChainID:      (*utils.Big)(&cltest.FixtureChainID),
		},
		ExternalJobID: uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
		PipelineSpec: &pipeline.Spec{
			ID:           1,
			DotDagSource: `ds1 [type=http method=GET url="https://pricesource1.com"`,
		},
		Type:            job.DirectRequest,
		SchemaVersion:   1,
		Name:            null.StringFrom("test"),
		MaxTaskDuration: models.Interval(1 * time.Minute),
	}

	err := jorm.CreateJob(&drJob)
	require.NoError(t, err)
	require.Equal(t, job.DirectRequest, drJob.Type)

	runningID := uuid.NewV4()

	err = porm.CreateRun(&pipeline.Run{
		PipelineSpecID: drJob.PipelineSpecID,
		State:          pipeline.RunStatusRunning,
		Outputs:        pipeline.JSONSerializable{},
		CreatedAt:      time.Now(),
		PipelineTaskRuns: []pipeline.TaskRun{{
			ID:        runningID,
			Type:      pipeline.TaskTypeHTTP,
			Index:     0,
			Output:    pipeline.JSONSerializable{},
			CreatedAt: time.Now(),
			DotID:     "ds1",
		}},
	})
	require.NoError(t, err)

	err = porm.CreateRun(&pipeline.Run{
		PipelineSpecID: drJob.PipelineSpecID,
		State:          pipeline.RunStatusSuspended,
		Outputs:        pipeline.JSONSerializable{},
		CreatedAt:      time.Now(),
		PipelineTaskRuns: []pipeline.TaskRun{{
			ID:        uuid.NewV4(),
			Type:      pipeline.TaskTypeHTTP,
			Index:     1,
			Output:    pipeline.JSONSerializable{},
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

	cfg := configtest2.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.MaxSuccessfulRuns = &n
	})
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	db := pgtest.NewSqlxDB(t)
	porm := pipeline.NewORM(db, lggr, cfg)

	ps1 := cltest.MustInsertPipelineSpec(t, db)

	t.Run("when there are no runs to prune, does nothing", func(t *testing.T) {
		porm.Prune(db, ps1.ID)

		// no error logs; it did nothing
		assert.Empty(t, observed.All())
	})

	// ps1 has:
	// - 20 completed runs
	for i := 0; i < 20; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps1.ID, pipeline.RunStatusCompleted)
	}

	ps2 := cltest.MustInsertPipelineSpec(t, db)

	// ps2 has:
	// - 12 completed runs
	// - 3 errored runs
	// - 3 running run
	// - 3 suspended run
	for i := 0; i < 12; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps2.ID, pipeline.RunStatusCompleted)
	}
	for i := 0; i < 3; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps2.ID, pipeline.RunStatusErrored)
	}
	for i := 0; i < 3; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps2.ID, pipeline.RunStatusRunning)
	}
	for i := 0; i < 3; i++ {
		cltest.MustInsertPipelineRunWithStatus(t, db, ps2.ID, pipeline.RunStatusSuspended)
	}

	porm.Prune(db, ps2.ID)

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
