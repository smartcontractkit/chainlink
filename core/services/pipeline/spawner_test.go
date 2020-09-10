package pipeline

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func TestSpawner_AddJobRemoveJob(t *testing.T) {
	t.Run("starts and stops job services when jobs are added and removed", func(t *testing.T) {
		orm := new(mocks.JobSpawnerORM)

		var jobs []JobSpec
		fnCall := orm.On("UnclaimedJobs")
		fn.RunFn = func(args mock.Arguments) {
			fnCall.ReturnArguments = mock.Arguments{jobs, nil}
		}

		spawner := NewSpawner(orm)
		err := spawner.Start()
		require.NoError(t, err)

		jobIDA := models.NewID()
		jobSpecA := new(mocks.JobSpec)
		jobSpecA.On("JobType").Return("AAA")
		jobSpecA.On("JobID").Return(jobIDA)

		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Run(func(mock.Arguments) { jobs = jobs[1:] }).Once()
		serviceA2.On("Start").Return(nil).Once()
		spawner.RegisterJobType("AAA", func(jobSpec JobSpec) ([]Service, error) {
			jobs = append(jobs, jobSpec)
			require.Equal(t, jobIDA, jobSpec.JobID())
			require.Equal(t, "AAA", jobSpec.JobType())
			return []Service{serviceA1, serviceA2}, nil
		})

		err = spawner.AddJob(jobSpecA)
		require.NoError(t, err)
		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2) }, 5*time.Second, 100*time.Millisecond)

		jobIDB := models.NewID()
		jobSpecB := new(mocks.JobSpec)
		jobSpecB.On("JobType").Return("BBB")
		jobSpecB.On("JobID").Return(jobIDB)

		serviceB1 := new(mocks.Service)
		serviceB2 := new(mocks.Service)
		serviceB1.On("Start").Return(nil).Run(func(mock.Arguments) { jobs = jobs[1:] }).Once()
		serviceB2.On("Start").Return(nil).Once()
		spawner.RegisterJobType("BBB", func(jobSpec JobSpec) ([]Service, error) {
			jobs = append(jobs, jobSpec)
			require.Equal(t, jobIDB, jobSpec.JobID())
			require.Equal(t, "BBB", jobSpec.JobType())
			return []Service{serviceB1, serviceB2}, nil
		})

		err = spawner.AddJob(jobSpecB)
		require.NoError(t, err)
		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceB1, serviceB2) }, 5*time.Second, 100*time.Millisecond)

		serviceA1.On("Stop").Return(nil).Once()
		serviceA2.On("Stop").Return(nil).Once()
		spawner.RemoveJob(jobSpecA.JobID())

		serviceB1.On("Stop").Return(nil).Once()
		serviceB2.On("Stop").Return(nil).Once()
		spawner.RemoveJob(jobSpecB.JobID())

		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2, serviceB1, serviceB2) }, 5*time.Second, 100*time.Millisecond)

		spawner.Stop()
	})

	t.Run("starts job services from the DB when .Start() is called", func(t *testing.T) {
		jobIDA := models.NewID()
		jobSpecA := new(mocks.JobSpec)
		jobSpecA.On("JobType").Return("AAA")
		jobSpecA.On("JobID").Return(jobIDA)

		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once()

		orm := new(mocks.JobSpawnerORM)
		orm.On("UnclaimedJobs").Return([]Spec{jobSpecA}, nil).Once()
		orm.On("UnclaimedJobs").Return(nil, nil)

		spawner := NewSpawner(orm)

		spawner.RegisterJobType("AAA", func(jobSpec JobSpec) ([]Service, error) {
			require.Equal(t, jobIDA, jobSpec.JobID())
			require.Equal(t, "AAA", jobSpec.JobType())
			return []Service{serviceA1, serviceA2}, nil
		})

		err := spawner.Start()
		require.NoError(t, err)
		defer spawner.Stop()

		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2) }, 5*time.Second, 100*time.Millisecond)

		serviceA1.On("Stop").Return(nil).Once()
		serviceA2.On("Stop").Return(nil).Once()
	})

	t.Run("stops job services when .Stop() is called", func(t *testing.T) {
		jobIDA := models.NewID()
		jobSpecA := new(mocks.JobSpec)
		jobSpecA.On("JobType").Return("AAA")
		jobSpecA.On("JobID").Return(jobIDA)

		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once()

		orm := new(mocks.JobSpawnerORM)
		orm.On("UnclaimedJobs").Return([]Spec{jobSpecA}, nil).Once()
		orm.On("UnclaimedJobs").Return(nil, nil)

		spawner := NewSpawner(orm)

		spawner.RegisterJobType("AAA", func(jobSpec JobSpec) ([]Service, error) {
			require.Equal(t, jobIDA, jobSpec.JobID())
			require.Equal(t, "AAA", jobSpec.JobType())
			return []Service{serviceA1, serviceA2}, nil
		})

		err := spawner.Start()
		require.NoError(t, err)

		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2) }, 5*time.Second, 100*time.Millisecond)

		serviceA1.On("Stop").Return(nil).Once()
		serviceA2.On("Stop").Return(nil).Once()

		spawner.Stop()

		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2) }, 5*time.Second, 100*time.Millisecond)
	})
}
