package cron_test

import (
	"context"
	"testing"
	"time"

	pipelinemocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"

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
	runner := new(pipelinemocks.Runner)
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "service_cron_orm", true, true)
	db := oldORM.DB
	orm, eventBroadcaster, cleanupPipeline := cltest.NewPipelineORM(t, config, db)
	jobORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})

	cleanup := func() {
		cleanupDB()
		cleanupPipeline()
	}
	defer cleanup()

	spec := &job.Job{
		Type:          job.Cron,
		SchemaVersion: 1,
		CronSpec:      &job.CronSpec{CronSchedule: "@every 1s"},
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
	spec := job.Job{
		Type:          job.Cron,
		SchemaVersion: 1,
		CronSpec:      &job.CronSpec{CronSchedule: "@every 1s"},
		Pipeline:      *pipeline.NewTaskDAG(),
		PipelineSpec:  &pipeline.Spec{},
	}
	runner := new(pipelinemocks.Runner)

	runner.On("ExecuteAndInsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, false).
		Return(int64(0), pipeline.FinalResult{}, nil)

	service, err := cron.NewCronFromJobSpec(spec, runner)
	require.NoError(t, err)
	err = service.Start()
	require.NoError(t, err)
	defer service.Close()

	cltest.EventuallyExpectationsMet(t, runner, 10*time.Second, 1*time.Second)
}
