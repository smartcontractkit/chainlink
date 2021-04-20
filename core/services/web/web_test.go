package web_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/web"

	pipeline_mocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebJobV2Pipeline(t *testing.T) {
	runner := new(pipeline_mocks.Runner)
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "services_web_orm", true, true)
	db := oldORM.DB
	orm, eventBroadcaster, cleanupPipeline := cltest.NewPipelineORM(t, config, db)
	jobORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})

	cleanup := func() {
		cleanupDB()
		cleanupPipeline()
	}
	defer cleanup()

	spec := &job.Job{
		Type:          job.Web,
		SchemaVersion: 1,
		WebSpec:       &job.WebSpec{},
		Pipeline:      *pipeline.NewTaskDAG(),
		PipelineSpec:  &pipeline.Spec{},
	}
	delegate := web.NewDelegate(runner)

	err := jobORM.CreateJob(context.Background(), spec, spec.Pipeline)
	require.NoError(t, err)
	serviceArray, err := delegate.ServicesForSpec(*spec)
	require.NoError(t, err)
	assert.Len(t, serviceArray, 1)
	service := serviceArray[0]

	err = service.Start()
	defer service.Close()

	require.NoError(t, err)
}
