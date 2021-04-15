package cron

import (
	"testing"
)

func TestCronJobV2Pipeline(t *testing.T) {
	/* Causes Import Cycle

	TODO: Resolve test case

	runner := new(pipeline_mocks.Runner)
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "delegate_services_listener_handlelog", true, true)
	db := oldORM.DB
	orm, eventBroadcaster, cleanupPipeline := cltest.NewPipelineORM(t, config, db)
	jobORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})

	cleanup := func() {
		cleanupDB()
		cleanupPipeline()
	}
	defer cleanup()

	spec := &job.Job{
		Type:          job.CronJob,
		SchemaVersion: 1,
		CronSpec:      &job.CronSpec{},
		Pipeline:      *pipeline.NewTaskDAG(),
		PipelineSpec:  &pipeline.Spec{},
	}

	delegate := NewDelegate(runner)

	err := jobORM.CreateJob(context.Background(), spec, spec.Pipeline)
	require.NoError(t, err)
	serviceArray, err := delegate.ServicesForSpec(*spec)
	require.NoError(t, err)
	assert.Len(t, serviceArray, 1)
	service := serviceArray[0]

	err = service.Start()
	require.NoError(t, err)
	*/

}
