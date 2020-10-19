package pipeline_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestRunner(t *testing.T) {
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "pipeline_runner", true, true)
	defer cleanupDB()
	db := oldORM.DB

	pipelineORM := pipeline.NewORM(db, config)
	runner := pipeline.NewRunner(pipelineORM, config)
	jobORM := job.NewORM(db, config, pipelineORM)
	defer jobORM.Close()

	runner.Start()
	defer runner.Stop()

	t.Run("gets the election result winner", func(t *testing.T) {
		var httpURL string
		{
			mockElectionWinner, cleanupElectionWinner := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `Hal Finney`)
			defer cleanupElectionWinner()
			mockVoterTurnout, cleanupVoterTurnout := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"data": {"result": 62.57}}`)
			defer cleanupVoterTurnout()
			mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"turnout": 61.942}`)
			defer cleanupHTTP()

			_, bridgeER := cltest.NewBridgeType(t, "election_winner", mockElectionWinner.URL)
			err := db.Create(bridgeER).Error
			require.NoError(t, err)

			_, bridgeVT := cltest.NewBridgeType(t, "voter_turnout", mockVoterTurnout.URL)
			err = db.Create(bridgeVT).Error
			require.NoError(t, err)

			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		ocrSpec, dbSpec := makeVoterTurnoutOCRJobSpecWithHTTPURL(t, db, httpURL)
		err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
		require.NoError(t, err)

		runID, err := runner.CreateRun(context.Background(), dbSpec.ID, nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = runner.AwaitRun(ctx, runID)
		require.NoError(t, err)

		// Verify the final pipeline results
		results, err := runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)

		assert.Len(t, results, 2)
		assert.NoError(t, results[0].Error)
		assert.NoError(t, results[1].Error)
		assert.Equal(t, "6225.6", results[0].Value)
		assert.Equal(t, "Hal Finney", results[1].Value)

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Preload("PipelineTaskSpec").
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		assert.Len(t, runs, 9)

		for _, run := range runs {
			if run.DotID() == "answer2" {
				assert.Equal(t, "Hal Finney", run.Output.Val)
			} else if run.DotID() == "ds2" {
				assert.Equal(t, `{"turnout": 61.942}`, run.Output.Val)
			} else if run.DotID() == "ds2_parse" {
				assert.Equal(t, float64(61.942), run.Output.Val)
			} else if run.DotID() == "ds2_multiply" {
				assert.Equal(t, "6194.2", run.Output.Val)
			} else if run.DotID() == "ds1" {
				assert.Equal(t, `{"data": {"result": 62.57}}`, run.Output.Val)
			} else if run.DotID() == "ds1_parse" {
				assert.Equal(t, float64(62.57), run.Output.Val)
			} else if run.DotID() == "ds1_multiply" {
				assert.Equal(t, "6257", run.Output.Val)
			} else if run.DotID() == "answer1" {
				assert.Equal(t, "6225.6", run.Output.Val)
			} else if run.DotID() == "__result__" {
				assert.Equal(t, []interface{}{"6225.6", "Hal Finney"}, run.Output.Val)
			} else {
				t.Fatalf("unknown task '%v'", run.DotID())
			}
		}
	})

	t.Run("handles the case where the parsed value is literally null", func(t *testing.T) {
		var httpURL string
		resp := `{"USD": null}`
		{
			mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			defer cleanupHTTP()
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		ocrSpec, dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, httpURL)
		err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
		require.NoError(t, err)

		runID, err := runner.CreateRun(context.Background(), dbSpec.ID, nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = runner.AwaitRun(ctx, runID)
		require.NoError(t, err)

		// Verify the final pipeline results
		results, err := runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		assert.EqualError(t, results[0].Error, "type <nil> cannot be converted to decimal.Decimal")
		assert.Nil(t, results[0].Value)

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Preload("PipelineTaskSpec").
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		require.Len(t, runs, 4)

		for _, run := range runs {
			if run.DotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				assert.Equal(t, resp, run.Output.Val)
			} else if run.DotID() == "ds1_parse" {
				assert.True(t, run.Error.IsZero())
				// FIXME: Shouldn't it be the Val that is null?
				assert.Nil(t, run.Output)
			} else if run.DotID() == "ds1_multiply" {
				assert.Equal(t, "type <nil> cannot be converted to decimal.Decimal", run.Error.ValueOrZero())
				assert.Nil(t, run.Output)
			} else if run.DotID() == "__result__" {
				assert.Equal(t, []interface{}{nil}, run.Output.Val)
				assert.Equal(t, "[\"type \\u003cnil\\u003e cannot be converted to decimal.Decimal\"]", run.Error.ValueOrZero())
			} else {
				t.Fatalf("unknown task '%v'", run.DotID())
			}
		}
	})

	t.Run("handles the case where the jsonparse lookup path is missing from the http response", func(t *testing.T) {
		var httpURL string
		resp := "{\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}"
		{
			mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			defer cleanupHTTP()
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		ocrSpec, dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, httpURL)
		err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
		require.NoError(t, err)

		runID, err := runner.CreateRun(context.Background(), dbSpec.ID, nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = runner.AwaitRun(ctx, runID)
		require.NoError(t, err)

		// Verify the final pipeline results
		results, err := runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		assert.EqualError(t, results[0].Error, "type <nil> cannot be converted to decimal.Decimal")
		assert.Nil(t, results[0].Value)

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Preload("PipelineTaskSpec").
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		require.Len(t, runs, 4)

		for _, run := range runs {
			if run.DotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				assert.Equal(t, resp, run.Output.Val)
			} else if run.DotID() == "ds1_parse" {
				assert.True(t, run.Error.IsZero())
				// FIXME: Shouldn't it be the Val that is null?
				assert.Nil(t, run.Output)
			} else if run.DotID() == "ds1_multiply" {
				assert.Equal(t, "type <nil> cannot be converted to decimal.Decimal", run.Error.ValueOrZero())
				assert.Nil(t, run.Output)
			} else if run.DotID() == "__result__" {
				assert.Equal(t, []interface{}{nil}, run.Output.Val)
				assert.Equal(t, "[\"type \\u003cnil\\u003e cannot be converted to decimal.Decimal\"]", run.Error.ValueOrZero())
			} else {
				t.Fatalf("unknown task '%v'", run.DotID())
			}
		}
	})
}
