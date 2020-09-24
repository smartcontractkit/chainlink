package job_test

import (
	"context"
	// "sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type delegate struct {
	jobType                    job.Type
	services                   []job.Service
	jobID                      int32
	chContinueCreatingServices chan struct{}
	job.Delegate
}

func (d delegate) JobType() job.Type {
	return d.jobType
}

func (d delegate) ServicesForSpec(js job.Spec) ([]job.Service, error) {
	if js.JobType() != d.jobType {
		return nil, nil
	}
	return d.services, nil
}

func (d delegate) FromDBRow(dbRow models.JobSpecV2) job.Spec {
	if d.chContinueCreatingServices != nil {
		<-d.chContinueCreatingServices
	}
	if dbRow.ID != d.jobID {
		return nil
	}

	// Wrap
	inner := d.Delegate.FromDBRow(dbRow)
	return &spec{inner, d.jobType}
}

func (d delegate) ToDBRow(js job.Spec) models.JobSpecV2 {
	// Unwrap
	inner := js.(*spec).Spec.(*offchainreporting.OracleSpec)
	return d.Delegate.ToDBRow(inner)
}

type spec struct {
	job.Spec
	jobType job.Type
}

func (s spec) JobType() job.Type {
	return s.jobType
}

func TestSpawner_CreateJobDeleteJob(t *testing.T) {
	jobTypeA := job.Type("AAA")
	jobTypeB := job.Type("BBB")

	innerJobSpecA, _ := makeOCRJobSpec(t)
	innerJobSpecB, _ := makeOCRJobSpec(t)
	jobSpecA := &spec{innerJobSpecA, jobTypeA}
	jobSpecB := &spec{innerJobSpecB, jobTypeB}

	t.Run("starts and stops job services when jobs are added and removed", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		orm := job.NewORM(store.ORM.DB, store.Config.DatabaseURL(), pipeline.NewORM(store.ORM.DB))
		defer orm.Close()
		spawner := job.NewSpawner(orm)
		spawner.Start()

		eventuallyA := cltest.NewAwaiter()
		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventuallyA.ItHappened() })

		delegateA := &delegate{jobTypeA, []job.Service{serviceA1, serviceA2}, 0, make(chan struct{}), offchainreporting.NewJobSpawnerDelegate(nil, nil, nil, nil)}
		spawner.RegisterDelegate(delegateA)

		jobSpecIDA, err := spawner.CreateJob(jobSpecA)
		require.NoError(t, err)
		delegateA.jobID = jobSpecIDA
		close(delegateA.chContinueCreatingServices)

		eventuallyA.AwaitOrFail(t, 20*time.Second)
		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)

		eventuallyB := cltest.NewAwaiter()
		serviceB1 := new(mocks.Service)
		serviceB2 := new(mocks.Service)
		serviceB1.On("Start").Return(nil).Once()
		serviceB2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventuallyB.ItHappened() })

		delegateB := &delegate{jobTypeB, []job.Service{serviceB1, serviceB2}, 0, make(chan struct{}), offchainreporting.NewJobSpawnerDelegate(nil, nil, nil, nil)}
		spawner.RegisterDelegate(delegateB)

		jobSpecIDB, err := spawner.CreateJob(jobSpecB)
		require.NoError(t, err)
		delegateB.jobID = jobSpecIDB
		close(delegateB.chContinueCreatingServices)

		eventuallyB.AwaitOrFail(t, 20*time.Second)
		mock.AssertExpectationsForObjects(t, serviceB1, serviceB2)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		serviceA1.On("Stop").Return(nil).Once()
		serviceA2.On("Stop").Return(nil).Once()
		spawner.DeleteJob(ctx, jobSpecA)

		serviceB1.On("Stop").Return(nil).Once()
		serviceB2.On("Stop").Return(nil).Once()
		spawner.DeleteJob(ctx, jobSpecB)

		spawner.Stop()
		serviceA1.AssertExpectations(t)
		serviceA2.AssertExpectations(t)
		serviceB1.AssertExpectations(t)
		serviceB2.AssertExpectations(t)
	})

	t.Run("starts job services from the DB when .Start() is called", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		eventually := cltest.NewAwaiter()
		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventually.ItHappened() })

		orm := job.NewORM(store.ORM.DB, store.Config.DatabaseURL(), pipeline.NewORM(store.ORM.DB))
		defer orm.Close()
		spawner := job.NewSpawner(orm)

		delegateA := &delegate{jobTypeA, []job.Service{serviceA1, serviceA2}, 0, nil, offchainreporting.NewJobSpawnerDelegate(nil, nil, nil, nil)}
		spawner.RegisterDelegate(delegateA)

		jobSpecIDA, err := spawner.CreateJob(jobSpecA)
		require.NoError(t, err)
		delegateA.jobID = jobSpecIDA

		spawner.Start()
		defer spawner.Stop()

		eventually.AwaitOrFail(t, 10*time.Second)
		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)

		serviceA1.On("Stop").Return(nil).Once()
		serviceA2.On("Stop").Return(nil).Once()
	})

	t.Run("stops job services when .Stop() is called", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		eventually := cltest.NewAwaiter()
		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventually.ItHappened() })

		orm := job.NewORM(store.ORM.DB, store.Config.DatabaseURL(), pipeline.NewORM(store.ORM.DB))
		defer orm.Close()
		spawner := job.NewSpawner(orm)

		delegateA := &delegate{jobTypeA, []job.Service{serviceA1, serviceA2}, 0, nil, offchainreporting.NewJobSpawnerDelegate(nil, nil, nil, nil)}
		spawner.RegisterDelegate(delegateA)

		jobSpecIDA, err := spawner.CreateJob(jobSpecA)
		require.NoError(t, err)
		delegateA.jobID = jobSpecIDA

		spawner.Start()

		eventually.AwaitOrFail(t, 10*time.Second)
		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)

		serviceA1.On("Stop").Return(nil).Once()
		serviceA2.On("Stop").Return(nil).Once()

		spawner.Stop()

		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)
	})
}
