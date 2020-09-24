package pipeline_test

import (
	"context"
	"math/rand"
	"net/url"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	ormpkg "github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestORM(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	db, err := gorm.Open(string(ormpkg.DialectPostgres), config.DatabaseURL())
	require.NoError(t, err)
	defer db.Close()

	orm := pipeline.NewORM(db)

	var specID int32

	u, err := url.Parse("https://chain.link/voter_turnout/USA-2020")
	require.NoError(t, err)

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
	expectedTasks := []pipeline.Task{answer1, answer2, ds1_multiply, ds1_parse, ds1, ds2_multiply, ds2_parse, ds2}
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
		g := pipeline.NewTaskDAG()
		err := g.UnmarshalText([]byte(dotStr))
		require.NoError(t, err)

		specID, err = orm.CreateSpec(*g)
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
		jobORM := job.NewORM(db, config.DatabaseURL(), orm)
		defer jobORM.Close()

		ocrSpec, dbSpec := makeOCRJobSpec(t)

		// Need a job in order to create a run
		err := jobORM.CreateJob(dbSpec, ocrSpec.TaskDAG())
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
		runID, err = orm.CreateRun(dbSpec.ID)
		require.NoError(t, err)

		// Check the DB for the pipeline.Run
		var pipelineRuns []pipeline.Run
		err = db.Find(&pipelineRuns).Error
		require.NoError(t, err)
		require.Len(t, pipelineRuns, 1)
		require.Equal(t, pipelineSpecID, pipelineRuns[0].PipelineSpecID)
		require.Equal(t, runID, pipelineRuns[0].ID)

		// Check the DB for the pipeline.TaskRuns
		var taskRuns []pipeline.TaskRun
		err = db.Find(&taskRuns).Error
		require.NoError(t, err)
		require.Len(t, taskRuns, len(taskSpecIDs))

		for _, taskRun := range taskRuns {
			require.Equal(t, runID, taskRun.PipelineRunID)
			require.Contains(t, taskSpecIDs, taskRun.PipelineTaskSpecID)
			require.Nil(t, taskRun.Output)
			require.True(t, taskRun.Error.IsZero())
		}
	})

	var (
		answers      = make(map[string]float64)
		taskRuns     = make(map[string]pipeline.TaskRun)
		predecessors = make(map[string][]pipeline.TaskRun)
	)

	t.Run("processes runs and awaits their completion", func(t *testing.T) {
		// Set up a goroutine to await the run's completion
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		chRunComplete := make(chan struct{})
		go func() {
			err := orm.AwaitRun(ctx, runID)
			require.NoError(t, err)
			close(chRunComplete)
		}()

		// First, "claim" one of the output task runs to ensure that `WithNextUnclaimedTaskRun` doesn't return it
		chClaimed := make(chan struct{})
		chBlock := make(chan struct{})
		chUnlocked := make(chan struct{})
		var locked pipeline.TaskRun
		go func() {
			err := utils.GormTransaction(db, func(tx *gorm.DB) error {
				err := tx.Raw(`
                    SELECT * FROM pipeline_task_runs
                    INNER JOIN pipeline_task_specs on pipeline_task_runs.pipeline_task_spec_id = pipeline_task_specs.id
                    WHERE pipeline_task_specs.type = 'median'
                    FOR UPDATE OF pipeline_task_runs
                `).Scan(&locked).Error
				require.NoError(t, err)

				close(chClaimed)
				<-chBlock
				return nil
			})
			require.NoError(t, err)
			close(chUnlocked)
		}()
		<-chClaimed

		// Process all of the unclaimed task runs
		{
			var done bool
			for !done {
				done, err = orm.WithNextUnclaimedTaskRun(func(taskRun pipeline.TaskRun, predecessorRuns []pipeline.TaskRun) pipeline.Result {
					// Ensure we don't fetch the locked task run
					require.NotEqual(t, locked.ID, taskRun.ID)

					// Ensure the predecessors' answers match what we expect
					for _, p := range predecessorRuns {
						require.Equal(t, answers[p.DotID], p.Output.Val)
						require.True(t, p.Error.IsZero())
					}

					taskRuns[taskRun.DotID] = taskRun
					answers[taskRun.DotID] = rand.Float64()
					predecessors[taskRun.DotID] = predecessorRuns
					return pipeline.Result{Value: answers[taskRun.DotID]}
				})
				require.NoError(t, err)
			}
		}

		// Ensure the run isn't considered complete yet
		{
			time.Sleep(3 * time.Second)
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

			done, err := orm.WithNextUnclaimedTaskRun(func(taskRun pipeline.TaskRun, predecessorRuns []pipeline.TaskRun) pipeline.Result {
				// Ensure the predecessors' answers match what we expect
				for _, p := range predecessorRuns {
					require.Equal(t, answers[p.DotID], p.Output.Val)
					require.True(t, p.Error.IsZero())
				}

				val := rand.Float64()
				taskRuns[taskRun.DotID] = taskRun
				answers[taskRun.DotID] = val
				predecessors[taskRun.DotID] = predecessorRuns
				return pipeline.Result{Value: val}
			})
			require.NoError(t, err)
			require.False(t, done)
		}

		// Ensure that the ORM doesn't think there are more runs
		{
			done, err := orm.WithNextUnclaimedTaskRun(func(taskRun pipeline.TaskRun, predecessorRuns []pipeline.TaskRun) pipeline.Result {
				val := rand.Float64()
				taskRuns[taskRun.DotID] = taskRun
				answers[taskRun.DotID] = val
				predecessors[taskRun.DotID] = predecessorRuns
				return pipeline.Result{Value: val}
			})
			require.NoError(t, err)
			require.True(t, done)
		}

		// Ensure the run is now considered complete
		{
			select {
			case <-time.After(5 * time.Second):
				t.Fatal("run did not complete as expected")
			case <-chRunComplete:
			}

			var finishedRuns []pipeline.TaskRun
			err = db.Find(&finishedRuns).Error
			require.NoError(t, err)
			require.Len(t, finishedRuns, len(expectedTasks))

			for _, run := range finishedRuns {
				require.Equal(t, answers[run.DotID], run.Output.Val.(float64))
			}
		}
	})

	t.Run("it fetches run results", func(t *testing.T) {
		results, err := orm.ResultsForRun(runID)
		require.NoError(t, err)
		require.Len(t, results, 2)

		require.Equal(t, answers["answer1"], results[0].Value)
		require.Equal(t, answers["answer2"], results[1].Value)
		require.NoError(t, results[0].Error)
		require.NoError(t, results[1].Error)
	})
}
