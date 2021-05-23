package job_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
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
	t.Skip()
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "pipeline_orm", true, true)
	config.Set("DEFAULT_HTTP_TIMEOUT", "30ms")
	config.Set("MAX_HTTP_ATTEMPTS", "1")
	defer cleanupDB()
	db := oldORM.DB

	key := cltest.MustInsertRandomKey(t, db)
	transmitterAddress := key.Address.Address()

	var specID int32

	answer1 := &pipeline.MedianTask{
		BaseTask: pipeline.NewBaseTask("answer1", nil, 0, 0),
	}
	answer2 := &pipeline.BridgeTask{
		Name:     "election_winner",
		BaseTask: pipeline.NewBaseTask("answer2", nil, 1, 0),
	}
	ds1_multiply := &pipeline.MultiplyTask{
		Times:    "1.23",
		BaseTask: pipeline.NewBaseTask("ds1_multiply", answer1, 0, 0),
	}
	ds1_parse := &pipeline.JSONParseTask{
		Path:     "one,two",
		BaseTask: pipeline.NewBaseTask("ds1_parse", ds1_multiply, 0, 0),
	}
	ds1 := &pipeline.BridgeTask{
		Name:     "voter_turnout",
		BaseTask: pipeline.NewBaseTask("ds1", ds1_parse, 0, 0),
	}
	ds2_multiply := &pipeline.MultiplyTask{
		Times:    "4.56",
		BaseTask: pipeline.NewBaseTask("ds2_multiply", answer1, 0, 0),
	}
	ds2_parse := &pipeline.JSONParseTask{
		Path:     "three,four",
		BaseTask: pipeline.NewBaseTask("ds2_parse", ds2_multiply, 0, 0),
	}
	ds2 := &pipeline.HTTPTask{
		URL:         "https://chain.link/voter_turnout/USA-2020",
		Method:      "GET",
		RequestData: `{"hi": "hello"}`,
		BaseTask:    pipeline.NewBaseTask("ds2", ds2_parse, 0, 0),
	}
	expectedTasks := []pipeline.Task{answer1, answer2, ds1_multiply, ds1_parse, ds1, ds2_multiply, ds2_parse, ds2}
	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, db.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, db.Create(bridge2).Error)

	t.Run("creates task DAGs", func(t *testing.T) {
		clearJobsDb(t, db)
		orm, _, cleanup := cltest.NewPipelineORM(t, config, db)
		defer cleanup()

		g := pipeline.NewTaskDAG()
		err := g.UnmarshalText([]byte(pipeline.DotStr))
		require.NoError(t, err)

		specID, err = orm.CreateSpec(context.Background(), db, *g, models.Interval(0))
		require.NoError(t, err)

		var specs []pipeline.Spec
		err = db.Find(&specs).Error
		require.NoError(t, err)
		require.Len(t, specs, 1)
		require.Equal(t, specID, specs[0].ID)
		require.Equal(t, pipeline.DotStr, specs[0].DotDagSource)

		require.NoError(t, db.Exec(`DELETE FROM pipeline_specs`).Error)
	})

	t.Run("creates runs", func(t *testing.T) {
		clearJobsDb(t, db)
		orm, eventBroadcaster, cleanup := cltest.NewPipelineORM(t, config, db)
		defer cleanup()
		runner := pipeline.NewRunner(orm, config)
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
		runID, _, err := runner.ExecuteAndInsertFinishedRun(context.Background(), pipelineSpecs[0], nil, pipeline.JSONSerializable{}, *logger.Default, true)
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
