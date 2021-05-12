package cron_test

import (
	"context"
	"testing"
	"time"

	pipeline_mocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCronV2Pipeline(t *testing.T) {
	runner := new(pipeline_mocks.Runner)
	config, cleanup := cltest.NewConfig(t)
	t.Cleanup(cleanup)
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	t.Cleanup(cleanup)
	db := store.DB
	orm, eventBroadcaster, cleanupPipeline := cltest.NewPipelineORM(t, config, db)
	t.Cleanup(cleanupPipeline)
	jobORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})

	spec := &job.Job{
		Type:          job.Cron,
		SchemaVersion: 1,
		CronSpec:      &job.CronSpec{CronSchedule: "0 0 1 1 *"},
		Pipeline:      *pipeline.NewTaskDAG(),
		PipelineSpec:  &pipeline.Spec{},
	}
	delegate := cron.NewDelegate(runner)

	err := jobORM.CreateJob(context.Background(), spec, spec.Pipeline)
	require.NoError(t, err)
	serviceArray, err := delegate.ServicesForSpec(*spec)
	require.NoError(t, err)
	assert.Len(t, serviceArray, 1)
	service := serviceArray[0]

	err = service.Start()
	require.NoError(t, err)
	defer service.Close()
}

func TestCronV2Schedule(t *testing.T) {
	t.Parallel()

	t.Skip("FIXME: This test is flaky and needs rethinking")

	type tc struct {
		name             string
		schedule         string
		expectedNumCalls int
		waitMinutes      time.Duration
	}
	for _, tc := range []tc{
		{
			name:             "1_min_cron_no_execution",
			schedule:         "* * * * *",
			expectedNumCalls: 0,
			waitMinutes:      time.Second * 30,
		},
		{
			name:             "1_min_cron_one_execution",
			schedule:         "* * * * *",
			expectedNumCalls: 1,
			waitMinutes:      time.Minute,
		},
		{
			name:             "1_min_cron_two_executions",
			schedule:         "* * * * *",
			expectedNumCalls: 2,
			waitMinutes:      2 * time.Minute,
		},
		{
			name:             "2_min_cron_one_execution",
			schedule:         "*/2 * * * *",
			expectedNumCalls: 1,
			waitMinutes:      2 * time.Minute,
		},
	} {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			config, cleanup := cltest.NewConfig(t)
			t.Cleanup(cleanup)
			store, cleanup := cltest.NewStoreWithConfig(t, config)
			t.Cleanup(cleanup)
			db := store.DB
			orm, eventBroadcaster, cleanupPipeline := cltest.NewPipelineORM(t, config, db)
			t.Cleanup(cleanupPipeline)
			jobORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})

			spec := &job.Job{
				Type:          job.Cron,
				SchemaVersion: 1,
				CronSpec:      &job.CronSpec{CronSchedule: testCase.schedule},
				Pipeline:      *pipeline.NewTaskDAG(),
				PipelineSpec:  &pipeline.Spec{},
			}
			runner := new(pipeline_mocks.Runner)
			delegate := cron.NewDelegate(runner)
			err := jobORM.CreateJob(context.Background(), spec, spec.Pipeline)
			require.NoError(t, err)
			serviceArray, err := delegate.ServicesForSpec(*spec)
			require.NoError(t, err)
			assert.Len(t, serviceArray, 1)
			service := serviceArray[0]
			err = service.Start()
			require.NoError(t, err)
			defer service.Close()
			if testCase.expectedNumCalls > 0 {
				runner.On("ExecuteAndInsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Times(testCase.expectedNumCalls).
					Return(int64(0), pipeline.FinalResult{}, nil)
			}
			// Wait for cron schedules to execute given test case + buffer
			time.Sleep(testCase.waitMinutes + (10 * time.Second))
			runner.AssertExpectations(t)
		})
	}
}
