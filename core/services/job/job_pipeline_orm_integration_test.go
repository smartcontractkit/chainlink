package job_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func clearJobsDb(t *testing.T, db *gorm.DB) {
	err := db.Exec(`TRUNCATE jobs, pipeline_runs, pipeline_specs, pipeline_task_runs, pipeline_task_specs CASCADE`).Error
	require.NoError(t, err)
}

func TestPipelineORM_Integration(t *testing.T) {
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "pipeline_orm", true, true)
	defer cleanupDB()
	db := oldORM.DB

	key := cltest.MustInsertRandomKey(t, db)
	transmitterAddress := key.Address.Address()

	var specID int32

	u, err := url.Parse("https://chain.link/voter_turnout/USA-2020")
	require.NoError(t, err)

	result := &pipeline.ResultTask{
		BaseTask: pipeline.NewBaseTask("__result__", nil, 0),
	}
	answer1 := &pipeline.MedianTask{
		BaseTask: pipeline.NewBaseTask("answer1", nil, 0),
	}
	answer2 := &pipeline.BridgeTask{
		Name:     "election_winner",
		BaseTask: pipeline.NewBaseTask("answer2", nil, 1),
	}
	ds1_multiply := &pipeline.MultiplyTask{
		Times:    decimal.NewFromFloat(1.23),
		BaseTask: pipeline.NewBaseTask("ds1_multiply", answer1, 0),
	}
	ds1_parse := &pipeline.JSONParseTask{
		Path:     []string{"one", "two"},
		BaseTask: pipeline.NewBaseTask("ds1_parse", ds1_multiply, 0),
	}
	ds1 := &pipeline.BridgeTask{
		Name:     "voter_turnout",
		BaseTask: pipeline.NewBaseTask("ds1", ds1_parse, 0),
	}
	ds2_multiply := &pipeline.MultiplyTask{
		Times:    decimal.NewFromFloat(4.56),
		BaseTask: pipeline.NewBaseTask("ds2_multiply", answer1, 0),
	}
	ds2_parse := &pipeline.JSONParseTask{
		Path:     []string{"three", "four"},
		BaseTask: pipeline.NewBaseTask("ds2_parse", ds2_multiply, 0),
	}
	ds2 := &pipeline.HTTPTask{
		URL:         models.WebURL(*u),
		Method:      "GET",
		RequestData: pipeline.HttpRequestData{"hi": "hello"},
		BaseTask:    pipeline.NewBaseTask("ds2", ds2_parse, 0),
	}
	expectedTasks := []pipeline.Task{result, answer1, answer2, ds1_multiply, ds1_parse, ds1, ds2_multiply, ds2_parse, ds2}
	var expectedTaskSpecs []pipeline.TaskSpec
	for _, task := range expectedTasks {
		expectedTaskSpecs = append(expectedTaskSpecs, pipeline.TaskSpec{
			DotID:          task.DotID(),
			PipelineSpecID: specID,
			Type:           task.Type(),
			JSON:           pipeline.JSONSerializable{Val: task},
			Index:          task.OutputIndex(),
		})
	}

	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "blah")
	require.NoError(t, db.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "blah")
	require.NoError(t, db.Create(bridge2).Error)

	t.Run("creates task DAGs", func(t *testing.T) {
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

		var taskSpecs []pipeline.TaskSpec
		err = db.Find(&taskSpecs).Error
		require.NoError(t, err)
		require.Len(t, taskSpecs, len(expectedTaskSpecs))

		type equalser interface {
			ExportedEquals(otherTask pipeline.Task) bool
		}

		for _, taskSpec := range taskSpecs {
			taskSpec.JSON.Val.(map[string]interface{})["index"] = taskSpec.Index
			taskSpec.JSON.Val, err = pipeline.UnmarshalTaskFromMap(taskSpec.Type, taskSpec.JSON.Val, taskSpec.DotID, nil, nil, nil)
			require.NoError(t, err)

			var found bool
			for _, expected := range expectedTaskSpecs {
				if taskSpec.PipelineSpecID == specID &&
					taskSpec.Type == expected.Type &&
					taskSpec.Index == expected.Index &&
					taskSpec.JSON.Val.(equalser).ExportedEquals(expected.JSON.Val.(pipeline.Task)) {
					found = true
					break
				}
			}
			require.True(t, found)
		}

		require.NoError(t, db.Exec(`DELETE FROM pipeline_specs`).Error)
	})

	var runID int64
	t.Run("creates runs", func(t *testing.T) {
		orm, eventBroadcaster, cleanup := cltest.NewPipelineORM(t, config, db)
		defer cleanup()
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

		var taskSpecs []pipeline.TaskSpec
		err = db.Find(&taskSpecs).Error
		require.NoError(t, err)

		var taskSpecIDs []int32
		for _, taskSpec := range taskSpecs {
			taskSpecIDs = append(taskSpecIDs, taskSpec.ID)
		}

		// Create the run
		runID, err = orm.CreateRun(context.Background(), dbSpec.ID, nil)
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
		require.Len(t, taskRuns, len(taskSpecIDs))

		for _, taskRun := range taskRuns {
			require.Equal(t, runID, taskRun.PipelineRunID)
			require.Contains(t, taskSpecIDs, taskRun.PipelineTaskSpecID)
			require.Nil(t, taskRun.Output)
			require.True(t, taskRun.Error.IsZero())
		}
	})

	t.Run("processes runs and awaits their completion", func(t *testing.T) {
		tests := []struct {
			name       string
			answers    map[string]pipeline.Result
			runOutputs interface{}
			runErrors  interface{}
		}{
			{
				"all succeeded",
				map[string]pipeline.Result{
					"ds1":          {Value: float64(1)},
					"ds1_parse":    {Value: float64(2)},
					"ds1_multiply": {Value: float64(3)},
					"ds2":          {Value: float64(4)},
					"ds2_parse":    {Value: float64(5)},
					"ds2_multiply": {Value: float64(6)},
					"answer1":      {Value: float64(7)},
					"answer2":      {Value: float64(8)},
					"__result__":   {Value: []interface{}{float64(7), float64(8)}, Error: pipeline.FinalErrors{{}, {}}},
				},
				[]interface{}{float64(7), float64(8)},
				[]interface{}{nil, nil},
			},
			{
				"all failed",
				map[string]pipeline.Result{
					"ds1":          {Error: errors.New("fail 1")},
					"ds1_parse":    {Error: errors.New("fail 2")},
					"ds1_multiply": {Error: errors.New("fail 3")},
					"ds2":          {Error: errors.New("fail 4")},
					"ds2_parse":    {Error: errors.New("fail 5")},
					"ds2_multiply": {Error: errors.New("fail 6")},
					"answer1":      {Error: errors.New("fail 7")},
					"answer2":      {Error: errors.New("fail 8")},
					"__result__":   {Value: []interface{}{nil, nil}, Error: pipeline.FinalErrors{null.StringFrom("fail 7"), null.StringFrom("fail 8")}},
				},
				[]interface{}{nil, nil},
				[]interface{}{"fail 7", "fail 8"},
			},
			{
				"some succeeded, some failed",
				map[string]pipeline.Result{
					"ds1":          {Value: float64(1)},
					"ds1_parse":    {Error: errors.New("fail 1")},
					"ds1_multiply": {Error: errors.New("fail 2")},
					"ds2":          {Value: float64(2)},
					"ds2_parse":    {Value: float64(3)},
					"ds2_multiply": {Value: float64(4)},
					"answer1":      {Error: errors.New("fail 3")},
					"answer2":      {Value: float64(5)},
					"__result__":   {Value: []interface{}{nil, float64(5)}, Error: pipeline.FinalErrors{null.StringFrom("fail 3"), {}}},
				},
				[]interface{}{nil, float64(5)},
				[]interface{}{"fail 3", nil},
			},
		}

		for _, test := range tests {
			clearJobsDb(t, db)

			test := test
			t.Run(test.name, func(t *testing.T) {
				orm, eventBroadcaster, cleanup := cltest.NewPipelineORM(t, config, db)
				defer cleanup()
				ORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})
				defer ORM.Close()

				dbSpec := makeVoterTurnoutOCRJobSpec(t, db, transmitterAddress)

				// Need a job in order to create a run
				err := ORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
				require.NoError(t, err)

				// Create two runs
				// One will be processed, the other will be "locked" by another thread
				runID, err = orm.CreateRun(context.Background(), dbSpec.ID, nil)
				require.NoError(t, err)
				runID2, err := orm.CreateRun(context.Background(), dbSpec.ID, nil)
				require.NoError(t, err)

				// Set up a goroutine to await the run's completion
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				chRunComplete := make(chan struct{})
				go func() {
					err2 := orm.AwaitRun(ctx, runID)
					require.NoError(t, err2)
					close(chRunComplete)
				}()

				// First, delete one of the runs to implicitly lock and ensure that `ProcessNextUnfinishedRun` doesn't return it
				var (
					chClaimed = make(chan struct{})
					chBlock   = make(chan struct{})
					chDeleted = make(chan struct{})
				)
				go func() {
					err2 := postgres.GormTransaction(context.Background(), db, func(tx *gorm.DB) error {
						err2 := tx.Exec(`DELETE FROM pipeline_runs WHERE id = ?`, runID2).Error
						assert.NoError(t, err2)

						close(chClaimed)
						select {
						case <-chBlock:
						case <-time.After(30 * time.Second):
							t.Fatal("timed out unblocking")
						}
						return nil
					})
					close(chDeleted)
					require.NoError(t, err2)
				}()
				<-chClaimed

				// Process the run
				{
					var anyRemaining bool
					anyRemaining, err = orm.ProcessNextUnfinishedRun(context.Background(), func(_ context.Context, db *gorm.DB, run pipeline.Run, l logger.Logger) (trrs pipeline.TaskRunResults, err error) {
						for dotID, result := range test.answers {
							var tr pipeline.TaskRun
							require.NoError(t, db.
								Joins("INNER JOIN pipeline_task_specs ON pipeline_task_specs.id = pipeline_task_runs.pipeline_task_spec_id AND dot_id = ?", dotID).
								Where("pipeline_run_id = ? ", runID).
								First(&tr).Error)
							trr := pipeline.TaskRunResult{
								ID:         tr.ID,
								Result:     result,
								FinishedAt: time.Now(),
								IsTerminal: dotID == "__result__",
							}
							trrs = append(trrs, trr)
						}
						return trrs, nil
					})
					require.NoError(t, err)
					require.True(t, anyRemaining)
				}

				// Ensure that the ORM doesn't think there are more runs
				{
					anyRemaining, err2 := orm.ProcessNextUnfinishedRun(context.Background(), func(_ context.Context, db *gorm.DB, run pipeline.Run, l logger.Logger) (pipeline.TaskRunResults, error) {
						t.Fatal("this callback should never be reached")
						return nil, nil
					})
					require.NoError(t, err2)
					require.False(t, anyRemaining)
				}

				// Allow the extra run to be deleted
				close(chBlock)
				select {
				case <-chDeleted:
				case <-time.After(30 * time.Second):
					t.Fatal("timed out waiting for delete")
				}

				// Ensure that the run is now considered complete
				{
					select {
					case <-time.After(5 * time.Second):
						t.Fatal("run did not complete as expected")
					case <-chRunComplete:
					}

					var finishedTaskRuns []pipeline.TaskRun
					err = db.Preload("PipelineTaskSpec").Find(&finishedTaskRuns, "pipeline_run_id = ?", runID).Error
					require.NoError(t, err)
					require.Len(t, finishedTaskRuns, len(expectedTasks))

					for _, run := range finishedTaskRuns {
						require.True(t, run.Output != nil || !run.Error.IsZero())
						if run.Output != nil {
							require.Equal(t, test.answers[run.DotID()].Value, run.Output.Val)
						} else if !run.Error.IsZero() {
							require.Equal(t, test.answers[run.DotID()].Error.Error(), run.Error.ValueOrZero())
						}
					}

					var pipelineRun pipeline.Run
					err = db.First(&pipelineRun).Error
					require.NoError(t, err)

					require.NotNil(t, pipelineRun.Errors.Val)
					require.Equal(t, test.runErrors, pipelineRun.Errors.Val)
					require.NotNil(t, pipelineRun.Outputs.Val)
					require.Equal(t, test.runOutputs, pipelineRun.Outputs.Val)
				}

				// Ensure that we can retrieve the correct results by calling .ResultsForRun
				results, err := orm.ResultsForRun(context.Background(), runID)
				require.NoError(t, err)
				require.Len(t, results, 2)

				if test.answers["answer1"].Value != nil {
					require.Equal(t, test.answers["answer1"], results[0])
				} else {
					require.Equal(t, test.answers["answer1"].Error.Error(), results[0].Error.Error())
				}

				if test.answers["answer2"].Value != nil {
					require.Equal(t, test.answers["answer2"], results[1])
				} else {
					require.Equal(t, test.answers["answer2"].Error.Error(), results[1].Error.Error())
				}
			})
		}
	})

}

func TestORM_CreateRunWhenJobDeleted(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	db := store.DB

	orm, _, cleanup := cltest.NewPipelineORM(t, config, db)
	defer cleanup()

	// Use non-existent job ID to simulate situation if a job is deleted between runs
	_, err := orm.CreateRun(context.Background(), -1, nil)
	require.EqualError(t, err, "no job found with id -1 (most likely it was deleted)")
}
