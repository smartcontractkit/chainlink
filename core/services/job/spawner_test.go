package job_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/job/mocks"
)

func mockORM(jobsInDB ...job.Spec) *mocks.ORM {
	orm := new(mocks.ORM)

	var jobsInDBMu sync.Mutex

	fnCall := orm.On("UnclaimedJobs", mock.Anything)
	fnCall.RunFn = func(args mock.Arguments) {
		jobsInDBMu.Lock()
		defer jobsInDBMu.Unlock()
		claimedJobs := make([]job.Spec, len(jobsInDB))
		copy(claimedJobs, jobsInDB)
		jobsInDB = nil
		fnCall.ReturnArguments = mock.Arguments{claimedJobs, nil}
	}
	fnCall.Maybe()

	orm.On("CreateJob", mock.Anything).
		Run(func(args mock.Arguments) {
			jobsInDBMu.Lock()
			defer jobsInDBMu.Unlock()
			jobsInDB = append(jobsInDB, args.Get(0).(job.Spec))
		}).
		Return(nil).
		Maybe()

	orm.On("DeleteJob", mock.Anything).
		Run(func(args mock.Arguments) {
			jobsInDBMu.Lock()
			defer jobsInDBMu.Unlock()
			jobID := args.Get(0).(int32)
			for i, job := range jobsInDB {
				if *job.JobID() == jobID {
					jobsInDB = append(jobsInDB[:i], jobsInDB[i+1:]...)
					break
				}
			}
		}).
		Return(nil).
		Maybe()

	return orm
}

func TestSpawner_CreateJobDeleteJob(t *testing.T) {
	t.Run("starts and stops job services when jobs are added and removed", func(t *testing.T) {
		orm := mockORM()

		spawner := job.NewSpawner(orm)
		spawner.Start()

		jobIDA := int32(1)
		jobSpecA := new(mocks.Spec)
		jobSpecA.On("JobType").Return(job.Type("AAA"))
		jobSpecA.On("JobID").Return(jobIDA)

		eventuallyA := cltest.NewAwaiter()
		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventuallyA.ItHappened() })
		spawner.RegisterJobType(job.Registration{
			JobType: "AAA",
			Spec:    nil,
			ServicesFactory: func(jobSpec job.Spec) ([]job.Service, error) {
				require.Equal(t, jobIDA, jobSpec.JobID())
				require.Equal(t, job.Type("AAA"), jobSpec.JobType())
				return []job.Service{serviceA1, serviceA2}, nil
			},
		})

		err := spawner.CreateJob(jobSpecA)
		require.NoError(t, err)

		eventuallyA.AwaitOrFail(t, 10*time.Second)
		mock.AssertExpectationsForObjects(t, orm, serviceA1, serviceA2)

		jobIDB := int32(2)
		jobSpecB := new(mocks.Spec)
		jobSpecB.On("JobType").Return(job.Type("BBB"))
		jobSpecB.On("JobID").Return(jobIDB)

		eventuallyB := cltest.NewAwaiter()
		serviceB1 := new(mocks.Service)
		serviceB2 := new(mocks.Service)
		serviceB1.On("Start").Return(nil).Once()
		serviceB2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventuallyB.ItHappened() })
		spawner.RegisterJobType(job.Registration{
			JobType: "BBB",
			Spec:    nil,
			ServicesFactory: func(jobSpec job.Spec) ([]job.Service, error) {
				require.Equal(t, jobIDB, jobSpec.JobID())
				require.Equal(t, job.Type("BBB"), jobSpec.JobType())
				return []job.Service{serviceB1, serviceB2}, nil
			},
		})

		err = spawner.CreateJob(jobSpecB)
		require.NoError(t, err)

		eventuallyB.AwaitOrFail(t, 10*time.Second)
		mock.AssertExpectationsForObjects(t, orm, serviceB1, serviceB2)

		serviceA1.On("Stop").Return(nil).Once()
		serviceA2.On("Stop").Return(nil).Once()
		spawner.DeleteJob(jobSpecA)

		serviceB1.On("Stop").Return(nil).Once()
		serviceB2.On("Stop").Return(nil).Once()
		spawner.DeleteJob(jobSpecB)

		spawner.Stop()
		orm.AssertExpectations(t)
		serviceA1.AssertExpectations(t)
		serviceA2.AssertExpectations(t)
		serviceB1.AssertExpectations(t)
		serviceB2.AssertExpectations(t)
	})

	t.Run("starts job services from the DB when .Start() is called", func(t *testing.T) {
		jobIDA := int32(1)
		jobSpecA := new(mocks.Spec)
		jobSpecA.On("JobType").Return(job.Type("AAA"))
		jobSpecA.On("JobID").Return(jobIDA)

		eventually := cltest.NewAwaiter()
		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventually.ItHappened() })

		orm := mockORM(jobSpecA)
		spawner := job.NewSpawner(orm)

		spawner.RegisterJobType(job.Registration{
			JobType: "AAA",
			Spec:    nil,
			ServicesFactory: func(jobSpec job.Spec) ([]job.Service, error) {
				require.Equal(t, jobIDA, jobSpec.JobID())
				require.Equal(t, job.Type("AAA"), jobSpec.JobType())
				return []job.Service{serviceA1, serviceA2}, nil
			},
		})

		spawner.Start()
		defer spawner.Stop()

		eventually.AwaitOrFail(t, 10*time.Second)
		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)

		serviceA1.On("Stop").Return(nil).Once()
		serviceA2.On("Stop").Return(nil).Once()
	})

	t.Run("stops job services when .Stop() is called", func(t *testing.T) {
		jobIDA := int32(1)
		jobSpecA := new(mocks.Spec)
		jobSpecA.On("JobType").Return(job.Type("AAA"))
		jobSpecA.On("JobID").Return(jobIDA)

		eventually := cltest.NewAwaiter()
		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventually.ItHappened() })

		orm := mockORM(jobSpecA)
		spawner := job.NewSpawner(orm)

		spawner.RegisterJobType(job.Registration{
			JobType: "AAA",
			Spec:    nil,
			ServicesFactory: func(jobSpec job.Spec) ([]job.Service, error) {
				require.Equal(t, jobIDA, jobSpec.JobID())
				require.Equal(t, job.Type("AAA"), jobSpec.JobType())
				return []job.Service{serviceA1, serviceA2}, nil
			},
		})

		spawner.Start()

		eventually.AwaitOrFail(t, 10*time.Second)
		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)

		serviceA1.On("Stop").Return(nil).Once()
		serviceA2.On("Stop").Return(nil).Once()

		spawner.Stop()

		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)
	})
}
