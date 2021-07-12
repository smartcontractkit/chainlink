package job_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func clearJobsDb(t *testing.T, db *gorm.DB) {
	t.Helper()
	err := db.Exec(`TRUNCATE jobs, pipeline_runs, pipeline_specs, pipeline_task_runs CASCADE`).Error
	require.NoError(t, err)
}

func TestPipelineORM_Integration(t *testing.T) {
	const DotStr = `
        // data source 1
        ds1          [type=bridge name=voter_turnout];
        ds1_parse    [type=jsonparse path="one,two"];
        ds1_multiply [type=multiply times=1.23];

        // data source 2
        ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}>];
        ds2_parse    [type=jsonparse path="three,four"];
        ds2_multiply [type=multiply times=4.56];

        ds1 -> ds1_parse -> ds1_multiply -> answer1;
        ds2 -> ds2_parse -> ds2_multiply -> answer1;

        answer1 [type=median                      index=0];
        answer2 [type=bridge name=election_winner index=1];
    `

	config, oldORM, cleanupDB := heavyweight.FullTestORM(t, "pipeline_orm", true, true)
	config.Set("DEFAULT_HTTP_TIMEOUT", "30ms")
	config.Set("MAX_HTTP_ATTEMPTS", "1")
	defer cleanupDB()
	db := oldORM.DB

	key := cltest.MustInsertRandomKey(t, db)
	transmitterAddress := key.Address.Address()

	var specID int32

	answer1 := &pipeline.MedianTask{
		AllowedFaults: "",
	}
	answer2 := &pipeline.BridgeTask{
		Name: "election_winner",
	}
	ds1_multiply := &pipeline.MultiplyTask{
		Times: "1.23",
	}
	ds1_parse := &pipeline.JSONParseTask{
		Path: "one,two",
	}
	ds1 := &pipeline.BridgeTask{
		Name: "voter_turnout",
	}
	ds2_multiply := &pipeline.MultiplyTask{
		Times: "4.56",
	}
	ds2_parse := &pipeline.JSONParseTask{
		Path: "three,four",
	}
	ds2 := &pipeline.HTTPTask{
		URL:         "https://chain.link/voter_turnout/USA-2020",
		Method:      "GET",
		RequestData: `{"hi": "hello"}`,
	}

	answer1.BaseTask = pipeline.NewBaseTask(6, "answer1", []pipeline.Task{ds1_multiply, ds2_multiply}, nil, 0)
	answer2.BaseTask = pipeline.NewBaseTask(7, "answer2", nil, nil, 1)
	ds1_multiply.BaseTask = pipeline.NewBaseTask(2, "ds1_multiply", []pipeline.Task{ds1_parse}, []pipeline.Task{answer1}, 0)
	ds2_multiply.BaseTask = pipeline.NewBaseTask(5, "ds2_multiply", []pipeline.Task{ds2_parse}, []pipeline.Task{answer1}, 0)
	ds1_parse.BaseTask = pipeline.NewBaseTask(1, "ds1_parse", []pipeline.Task{ds1}, []pipeline.Task{ds1_multiply}, 0)
	ds2_parse.BaseTask = pipeline.NewBaseTask(4, "ds2_parse", []pipeline.Task{ds2}, []pipeline.Task{ds2_multiply}, 0)
	ds1.BaseTask = pipeline.NewBaseTask(0, "ds1", nil, []pipeline.Task{ds1_parse}, 0)
	ds2.BaseTask = pipeline.NewBaseTask(3, "ds2", nil, []pipeline.Task{ds2_parse}, 0)
	expectedTasks := []pipeline.Task{ds1, ds1_parse, ds1_multiply, ds2, ds2_parse, ds2_multiply, answer1, answer2}
	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, db.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, db.Create(bridge2).Error)

	t.Run("creates task DAGs", func(t *testing.T) {
		clearJobsDb(t, db)
		orm, _, cleanup := cltest.NewPipelineORM(t, config, db)
		defer cleanup()

		p, err := pipeline.Parse(DotStr)
		require.NoError(t, err)

		specID, err = orm.CreateSpec(context.Background(), db, *p, models.Interval(0))
		require.NoError(t, err)

		var specs []pipeline.Spec
		err = db.Find(&specs).Error
		require.NoError(t, err)
		require.Len(t, specs, 1)
		require.Equal(t, specID, specs[0].ID)
		require.Equal(t, DotStr, specs[0].DotDagSource)

		require.NoError(t, db.Exec(`DELETE FROM pipeline_specs`).Error)
	})

	t.Run("creates runs", func(t *testing.T) {
		clearJobsDb(t, db)
		orm, eventBroadcaster, cleanup := cltest.NewPipelineORM(t, config, db)
		defer cleanup()
		runner := pipeline.NewRunner(orm, config, nil, nil)
		defer runner.Close()
		jobORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})
		defer jobORM.Close()

		dbSpec := makeVoterTurnoutOCRJobSpec(t, db, transmitterAddress)

		// Need a job in order to create a run
		err := jobORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)

		var pipelineSpecs []pipeline.Spec
		err = db.Find(&pipelineSpecs).Error
		require.NoError(t, err)
		require.Len(t, pipelineSpecs, 1)
		require.Equal(t, dbSpec.PipelineSpecID, pipelineSpecs[0].ID)
		pipelineSpecID := pipelineSpecs[0].ID

		// Create the run
		runID, _, err := runner.ExecuteAndInsertFinishedRun(context.Background(), pipelineSpecs[0], pipeline.NewVarsFrom(nil), *logger.Default, true)
		require.NoError(t, err)

		// Check the DB for the pipeline.Run
		var pipelineRuns []pipeline.Run
		err = db.Where("id = ?", runID).Find(&pipelineRuns).Error
		require.NoError(t, err)
		require.Len(t, pipelineRuns, 1)
		require.Equal(t, pipelineSpecID, pipelineRuns[0].PipelineSpecID)
		require.Equal(t, runID, pipelineRuns[0].ID)

		// Check the DB for the pipeline.TaskRuns
		var taskRuns []pipeline.TaskRun
		err = db.Where("pipeline_run_id = ?", runID).Find(&taskRuns).Error
		require.NoError(t, err)
		require.Len(t, taskRuns, len(expectedTasks))

		for _, taskRun := range taskRuns {
			require.Equal(t, runID, taskRun.PipelineRunID)
			require.Nil(t, taskRun.Output)
			require.False(t, taskRun.Error.IsZero())
		}
	})
}
