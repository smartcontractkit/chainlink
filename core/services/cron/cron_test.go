package cron_test

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

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
		CronSpec:      &job.CronSpec{CronSchedule: "@every 1s"},
		PipelineSpec:  &pipeline.Spec{},
		ExternalJobID: uuid.NewV4(),
	}
	delegate := cron.NewDelegate(runner)

	jb, err := jobORM.CreateJob(context.Background(), spec, spec.Pipeline)
	require.NoError(t, err)
	serviceArray, err := delegate.ServicesForSpec(jb)
	require.NoError(t, err)
	assert.Len(t, serviceArray, 1)
	service := serviceArray[0]

	err = service.Start()
	require.NoError(t, err)
	defer service.Close()
}

func TestCronV2Schedule(t *testing.T) {
	t.Parallel()

	spec := job.Job{
		Type:          job.Cron,
		SchemaVersion: 1,
		CronSpec:      &job.CronSpec{CronSchedule: "@every 1s"},
		PipelineSpec:  &pipeline.Spec{},
	}
	runner := new(pipelinemocks.Runner)

	runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything).
		Return(false, nil).Once()

	service, err := cron.NewCronFromJobSpec(spec, runner)
	require.NoError(t, err)
	err = service.Start()
	require.NoError(t, err)
	defer service.Close()

	cltest.EventuallyExpectationsMet(t, runner, 10*time.Second, 1*time.Second)
}
