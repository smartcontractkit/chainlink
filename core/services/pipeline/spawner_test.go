package pipeline_test

// import (
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"

// 	"github.com/smartcontractkit/chainlink/core/internal/mocks"
// 	"github.com/smartcontractkit/chainlink/core/services/pipeline"
// 	"github.com/smartcontractkit/chainlink/core/store/models"
// )

// func TestSpawner_AddJobRemoveJob(t *testing.T) {
// 	t.Run("starts and stops job services when jobs are added and removed", func(t *testing.T) {
// 		orm := new(mocks.JobSpawnerORM)
// 		orm.On("JobsAsInterfaces", mock.Anything).Return(nil).Once()

// 		spawner := job.NewSpawner(orm)
// 		err := spawner.Start()
// 		require.NoError(t, err)

// 		jobIDA := models.NewID()
// 		jobSpecA := new(mocks.JobSpec)
// 		jobSpecA.On("JobType").Return("AAA")
// 		jobSpecA.On("JobID").Return(jobIDA)

// 		serviceA1 := new(mocks.JobService)
// 		serviceA2 := new(mocks.JobService)
// 		serviceA1.On("Start").Return(nil).Once()
// 		serviceA2.On("Start").Return(nil).Once()
// 		spawner.RegisterJobType("AAA", func(jobSpec job.JobSpec) ([]job.JobService, error) {
// 			require.Equal(t, jobIDA, jobSpec.JobID())
// 			require.Equal(t, "AAA", jobSpec.JobType())
// 			return []job.JobService{serviceA1, serviceA2}, nil
// 		})

// 		err = spawner.AddJob(jobSpecA)
// 		require.NoError(t, err)
// 		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2) }, 5*time.Second, 100*time.Millisecond)

// 		jobIDB := models.NewID()
// 		jobSpecB := new(mocks.JobSpec)
// 		jobSpecB.On("JobType").Return("BBB")
// 		jobSpecB.On("JobID").Return(jobIDB)

// 		serviceB1 := new(mocks.JobService)
// 		serviceB2 := new(mocks.JobService)
// 		serviceB1.On("Start").Return(nil).Once()
// 		serviceB2.On("Start").Return(nil).Once()
// 		spawner.RegisterJobType("BBB", func(jobSpec job.JobSpec) ([]job.JobService, error) {
// 			require.Equal(t, jobIDB, jobSpec.JobID())
// 			require.Equal(t, "BBB", jobSpec.JobType())
// 			return []job.JobService{serviceB1, serviceB2}, nil
// 		})

// 		err = spawner.AddJob(jobSpecB)
// 		require.NoError(t, err)
// 		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceB1, serviceB2) }, 5*time.Second, 100*time.Millisecond)

// 		serviceA1.On("Stop").Return(nil).Once()
// 		serviceA2.On("Stop").Return(nil).Once()
// 		spawner.RemoveJob(jobSpecA.JobID())

// 		serviceB1.On("Stop").Return(nil).Once()
// 		serviceB2.On("Stop").Return(nil).Once()
// 		spawner.RemoveJob(jobSpecB.JobID())

// 		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2, serviceB1, serviceB2) }, 5*time.Second, 100*time.Millisecond)

// 		spawner.Stop()
// 	})

// 	t.Run("starts job services from the DB when .Start() is called", func(t *testing.T) {
// 		jobIDA := models.NewID()
// 		jobSpecA := new(mocks.JobSpec)
// 		jobSpecA.On("JobType").Return("AAA")
// 		jobSpecA.On("JobID").Return(jobIDA)

// 		serviceA1 := new(mocks.JobService)
// 		serviceA2 := new(mocks.JobService)
// 		serviceA1.On("Start").Return(nil).Once()
// 		serviceA2.On("Start").Return(nil).Once()

// 		orm := new(mocks.JobSpawnerORM)
// 		orm.On("JobsAsInterfaces", mock.Anything).
// 			Run(func(args mock.Arguments) {
// 				fn := args.Get(0).(func(job.JobSpec) bool)
// 				fn(jobSpecA)
// 			}).
// 			Return(nil).
// 			Once()

// 		spawner := job.NewSpawner(orm)

// 		spawner.RegisterJobType("AAA", func(jobSpec job.JobSpec) ([]job.JobService, error) {
// 			require.Equal(t, jobIDA, jobSpec.JobID())
// 			require.Equal(t, "AAA", jobSpec.JobType())
// 			return []job.JobService{serviceA1, serviceA2}, nil
// 		})

// 		err := spawner.Start()
// 		require.NoError(t, err)
// 		defer spawner.Stop()

// 		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2) }, 5*time.Second, 100*time.Millisecond)

// 		serviceA1.On("Stop").Return(nil).Once()
// 		serviceA2.On("Stop").Return(nil).Once()
// 	})

// 	t.Run("stops job services .Stop() is called", func(t *testing.T) {
// 		jobIDA := models.NewID()
// 		jobSpecA := new(mocks.JobSpec)
// 		jobSpecA.On("JobType").Return("AAA")
// 		jobSpecA.On("JobID").Return(jobIDA)

// 		serviceA1 := new(mocks.JobService)
// 		serviceA2 := new(mocks.JobService)
// 		serviceA1.On("Start").Return(nil).Once()
// 		serviceA2.On("Start").Return(nil).Once()

// 		orm := new(mocks.JobSpawnerORM)
// 		orm.On("JobsAsInterfaces", mock.Anything).
// 			Run(func(args mock.Arguments) {
// 				fn := args.Get(0).(func(job.JobSpec) bool)
// 				fn(jobSpecA)
// 			}).
// 			Return(nil).
// 			Once()

// 		spawner := job.NewSpawner(orm)

// 		spawner.RegisterJobType("AAA", func(jobSpec job.JobSpec) ([]job.JobService, error) {
// 			require.Equal(t, jobIDA, jobSpec.JobID())
// 			require.Equal(t, "AAA", jobSpec.JobType())
// 			return []job.JobService{serviceA1, serviceA2}, nil
// 		})

// 		err := spawner.Start()
// 		require.NoError(t, err)

// 		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2) }, 5*time.Second, 100*time.Millisecond)

// 		serviceA1.On("Stop").Return(nil).Once()
// 		serviceA2.On("Stop").Return(nil).Once()

// 		spawner.Stop()

// 		require.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, serviceA1, serviceA2) }, 5*time.Second, 100*time.Millisecond)
// 	})
// }
