package cron_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/cron"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	pipelinemocks "github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
)

func TestCronV2Pipeline(t *testing.T) {
	runner := pipelinemocks.NewRunner(t)
	cfg := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	lggr := logger.TestLogger(t)
	orm := pipeline.NewORM(db, lggr, cfg.JobPipeline().MaxSuccessfulRuns())
	btORM := bridges.NewORM(db)
	jobORM := job.NewORM(db, orm, btORM, keyStore, lggr)

	jb := &job.Job{
		Type:          job.Cron,
		SchemaVersion: 1,
		CronSpec:      &job.CronSpec{CronSchedule: "@every 1s"},
		PipelineSpec:  &pipeline.Spec{},
		ExternalJobID: uuid.New(),
	}
	delegate := cron.NewDelegate(runner, lggr)

	require.NoError(t, jobORM.CreateJob(testutils.Context(t), jb))
	serviceArray, err := delegate.ServicesForSpec(testutils.Context(t), *jb)
	require.NoError(t, err)
	assert.Len(t, serviceArray, 1)
	service := serviceArray[0]

	err = service.Start(testutils.Context(t))
	require.NoError(t, err)
	defer func() { assert.NoError(t, service.Close()) }()
}

func TestCronV2Schedule(t *testing.T) {
	t.Parallel()

	spec := job.Job{
		Type:          job.Cron,
		SchemaVersion: 1,
		CronSpec:      &job.CronSpec{CronSchedule: "@every 1s"},
		PipelineSpec:  &pipeline.Spec{},
	}
	runner := pipelinemocks.NewRunner(t)
	awaiter := cltest.NewAwaiter()
	runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { awaiter.ItHappened() }).
		Return(false, nil).
		Once()

	service, err := cron.NewCronFromJobSpec(spec, runner, logger.TestLogger(t))
	require.NoError(t, err)
	err = service.Start(testutils.Context(t))
	require.NoError(t, err)
	defer func() { assert.NoError(t, service.Close()) }()

	awaiter.AwaitOrFail(t)
}
