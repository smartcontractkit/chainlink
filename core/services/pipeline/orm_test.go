package pipeline_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func clearDB(t *testing.T, db *gorm.DB) {
	err := db.Exec(`TRUNCATE jobs, pipeline_runs, pipeline_specs, pipeline_task_runs, pipeline_task_specs CASCADE`).Error
	require.NoError(t, err)
}

func TestORM(t *testing.T) {
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "pipeline_orm", true, true)
	defer cleanupDB()
	db := oldORM.DB

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
			JSON:           pipeline.JSONSerializable{task},
			Index:          task.OutputIndex(),
		})
	}

	t.Run("creates task DAGs", func(t *testing.T) {
		eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
		defer eventBroadcaster.Stop()
		orm := pipeline.NewORM(db, config, eventBroadcaster)

		g := pipeline.NewTaskDAG()
		err := g.UnmarshalText([]byte(dotStr))
		require.NoError(t, err)

		specID, err = orm.CreateSpec(context.Background(), *g)
		require.NoError(t, err)

		var specs []pipeline.Spec
		err = db.Find(&specs).Error
		require.NoError(t, err)
		require.Len(t, specs, 1)
		require.Equal(t, specID, specs[0].ID)
		require.Equal(t, dotStr, specs[0].DotDagSource)

		var taskSpecs []pipeline.TaskSpec
		err = db.Find(&taskSpecs).Error
		require.NoError(t, err)
		require.Len(t, taskSpecs, len(expectedTaskSpecs))

		type equalser interface {
			ExportedEquals(otherTask pipeline.Task) bool
		}

		for _, taskSpec := range taskSpecs {
			taskSpec.JSON.Val.(map[string]interface{})["index"] = taskSpec.Index
			taskSpec.JSON.Val, err = pipeline.UnmarshalTaskFromMap(taskSpec.Type, taskSpec.JSON.Val, taskSpec.DotID, nil, nil)
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
		eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
		defer eventBroadcaster.Stop()
		orm := pipeline.NewORM(db, config, eventBroadcaster)
		jobORM := job.NewORM(db, config, orm, eventBroadcaster)
		defer jobORM.Close()

		ocrSpec, dbSpec := makeVoterTurnoutOCRJobSpec(t, db)

		// Need a job in order to create a run
		err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
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
			name    string
			answers map[string]pipeline.Result
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
			},
		}

		for _, test := range tests {
			clearDB(t, db)

			test := test
			t.Run(test.name, func(t *testing.T) {
				eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
				defer eventBroadcaster.Stop()
				orm := pipeline.NewORM(db, config, eventBroadcaster)
				jobORM := job.NewORM(db, config, orm, eventBroadcaster)
				defer jobORM.Close()

				var (
					taskRuns     = make(map[string]pipeline.TaskRun)
					predecessors = make(map[string][]pipeline.TaskRun)
				)

				ocrSpec, dbSpec := makeVoterTurnoutOCRJobSpec(t, db)

				// Need a job in order to create a run
				err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
				require.NoError(t, err)

				// Create the run
				runID, err = orm.CreateRun(context.Background(), dbSpec.ID, nil)
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

				// First, "claim" one of the output task runs to ensure that `ProcessNextUnclaimedTaskRun` doesn't return it
				var (
					chClaimed  = make(chan struct{})
					chBlock    = make(chan struct{})
					chUnlocked = make(chan struct{})
					locked     pipeline.TaskRun
				)
				go func() {
					err2 := postgres.GormTransaction(context.Background(), db, func(tx *gorm.DB) error {
						err2 := tx.Raw(`
                            SELECT * FROM pipeline_task_runs
                            INNER JOIN pipeline_task_specs on pipeline_task_runs.pipeline_task_spec_id = pipeline_task_specs.id
                            WHERE pipeline_task_specs.type = 'result'
                            FOR UPDATE OF pipeline_task_runs
                        `).Scan(&locked).Error
						require.NoError(t, err2)

						close(chClaimed)
						<-chBlock
						return nil
					})
					require.NoError(t, err2)
					close(chUnlocked)
				}()
				<-chClaimed

				// Process all of the unclaimed task runs
				{
					anyRemaining := true
					for anyRemaining {
						anyRemaining, err = orm.ProcessNextUnclaimedTaskRun(context.Background(), func(jobID int32, taskRun pipeline.TaskRun, predecessorRuns []pipeline.TaskRun) pipeline.Result {
							// Ensure we don't fetch the locked task run
							require.NotEqual(t, locked.ID, taskRun.ID)

							// Ensure the predecessors' answers match what we expect
							for _, p := range predecessorRuns {
								_, exists := test.answers[p.DotID()]
								require.True(t, exists)
								require.True(t, p.Output != nil || !p.Error.IsZero())
								if p.Output != nil {
									require.Equal(t, test.answers[p.DotID()].Value, p.Output.Val)
								} else if !p.Error.IsZero() {
									require.Equal(t, test.answers[p.DotID()].Error.Error(), p.Error.ValueOrZero())
								}
							}

							taskRuns[taskRun.DotID()] = taskRun
							predecessors[taskRun.DotID()] = predecessorRuns
							return test.answers[taskRun.DotID()]
						})
						require.NoError(t, err)
					}
				}

				// Ensure the run isn't considered complete yet
				{
					time.Sleep(5 * time.Second)
					require.Len(t, taskRuns, len(expectedTasks)-1)
					select {
					case <-chRunComplete:
						t.Fatal("run completed too early")
					default:
					}
				}

				// Now, release the claim and make sure we can process the final task run
				{
					close(chBlock)
					<-chUnlocked
					time.Sleep(3 * time.Second)

					anyRemaining, err2 := orm.ProcessNextUnclaimedTaskRun(context.Background(), func(jobID int32, taskRun pipeline.TaskRun, predecessorRuns []pipeline.TaskRun) pipeline.Result {
						fmt.Println(taskRun.DotID())
						// Ensure the predecessors' answers match what we expect
						for _, p := range predecessorRuns {
							_, exists := test.answers[p.DotID()]
							require.True(t, exists)
							require.True(t, p.Output != nil || !p.Error.IsZero())
							if p.Output != nil {
								require.Equal(t, test.answers[p.DotID()].Value, p.Output.Val)
							} else if !p.Error.IsZero() {
								require.Equal(t, test.answers[p.DotID()].Error.Error(), p.Error.ValueOrZero())
							}
						}

						taskRuns[taskRun.DotID()] = taskRun
						predecessors[taskRun.DotID()] = predecessorRuns
						return test.answers[taskRun.DotID()]
					})
					require.NoError(t, err2)
					require.True(t, anyRemaining)
				}

				// Ensure that the ORM doesn't think there are more runs
				{
					anyRemaining, err2 := orm.ProcessNextUnclaimedTaskRun(context.Background(), func(jobID int32, taskRun pipeline.TaskRun, predecessorRuns []pipeline.TaskRun) pipeline.Result {
						t.Fatal("this callback should never be reached")
						return pipeline.Result{}
					})
					require.NoError(t, err2)
					require.False(t, anyRemaining)
				}

				// Ensure that the run is now considered complete
				{
					select {
					case <-time.After(5 * time.Second):
						t.Fatal("run did not complete as expected")
					case <-chRunComplete:
					}

					var finishedRuns []pipeline.TaskRun
					err = db.Preload("PipelineTaskSpec").Find(&finishedRuns).Error
					require.NoError(t, err)
					require.Len(t, finishedRuns, len(expectedTasks))

					for _, run := range finishedRuns {
						require.True(t, run.Output != nil || !run.Error.IsZero())
						if run.Output != nil {
							require.Equal(t, test.answers[run.DotID()].Value, run.Output.Val)
						} else if !run.Error.IsZero() {
							require.Equal(t, test.answers[run.DotID()].Error.Error(), run.Error.ValueOrZero())
						}
					}
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
