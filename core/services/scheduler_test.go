package services_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestScheduler_Start_LoadingRecurringJobs(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	jobWCron := cltest.NewJobWithSchedule("* * * * * *")
	require.NoError(t, store.CreateJob(&jobWCron))
	jobWoCron := cltest.NewJob()
	require.NoError(t, store.CreateJob(&jobWoCron))

	executeJobChannel := make(chan struct{})
	runManager := new(mocks.RunManager)
	runManager.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).
		Twice().
		Run(func(mock.Arguments) {
			executeJobChannel <- struct{}{}
		})

	sched := services.NewScheduler(store, runManager)
	require.NoError(t, sched.Start())

	cltest.CallbackOrTimeout(t, "Create", func() {
		<-executeJobChannel
		<-executeJobChannel
	}, 3*time.Second)

	sched.Stop()

	runManager.AssertExpectations(t)
}

func TestRecurring_AddJob(t *testing.T) {
	executeJobChannel := make(chan struct{}, 1)
	runManager := new(mocks.RunManager)
	runManager.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).
		Run(func(mock.Arguments) {
			executeJobChannel <- struct{}{}
		}).
		Twice()

	r := services.NewRecurring(runManager)
	cron := cltest.NewMockCron()
	r.Cron = cron

	job := cltest.NewJobWithSchedule("* * * * *")
	r.AddJob(job)

	cron.RunEntries()

	cltest.CallbackOrTimeout(t, "Create", func() {
		<-executeJobChannel
	}, 3*time.Second)

	cron.RunEntries()

	cltest.CallbackOrTimeout(t, "Create", func() {
		<-executeJobChannel
	}, 3*time.Second)

	r.Stop()

	runManager.AssertExpectations(t)
}

func TestRecurring_AddJob_PastEnd(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)

	r := services.NewRecurring(runManager)
	cron := cltest.NewMockCron()
	r.Cron = cron

	j := cltest.NewJobWithSchedule("* * * * *")
	j.EndAt = null.TimeFrom(time.Now().Add(-1 * time.Second))
	require.Nil(t, store.CreateJob(&j))

	r.AddJob(j)
	cron.RunEntries()

	// Sleep for some time to make sure no calls are made
	time.Sleep(1 * time.Second)

	r.Stop()

	runManager.AssertExpectations(t)
}

func TestRecurring_AddJob_BeforeStart(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)

	r := services.NewRecurring(runManager)
	cron := cltest.NewMockCron()
	r.Cron = cron

	j := cltest.NewJobWithSchedule("* * * * *")
	j.StartAt = null.TimeFrom(time.Now().Add(1 * time.Hour))
	require.Nil(t, store.CreateJob(&j))

	r.AddJob(j)
	cron.RunEntries()

	// Sleep for some time to make sure no calls are made
	time.Sleep(1 * time.Second)

	r.Stop()

	runManager.AssertExpectations(t)
}

func TestOneTime_AddJob(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	executeJobChannel := make(chan struct{})
	runManager := new(mocks.RunManager)
	runManager.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).
		Once().
		Run(func(mock.Arguments) {
			executeJobChannel <- struct{}{}
		})

	clock := cltest.NewTriggerClock(t)

	ot := services.OneTime{
		Clock:      clock,
		Store:      store,
		RunManager: runManager,
	}
	require.NoError(t, ot.Start())

	j := cltest.NewJobWithRunAtInitiator(time.Now())
	require.Nil(t, store.CreateJob(&j))

	ot.AddJob(j)

	clock.Trigger()

	cltest.CallbackOrTimeout(t, "ws client restarts", func() {
		<-executeJobChannel
	}, 3*time.Second)

	// This should block because if OneTime works it won't listen on the channel again
	go clock.TriggerWithoutTimeout()

	// Sleep for some time to make sure another call isn't made
	time.Sleep(1 * time.Second)

	ot.Stop()

	runManager.AssertExpectations(t)
}

func TestOneTime_AddJob_PastEnd(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)

	clock := cltest.NewTriggerClock(t)

	ot := services.OneTime{
		Clock:      clock,
		Store:      store,
		RunManager: runManager,
	}
	require.NoError(t, ot.Start())

	j := cltest.NewJobWithRunAtInitiator(time.Now())
	j.EndAt = null.TimeFrom(clock.Now().Add(-1 * time.Second))
	require.Nil(t, store.CreateJob(&j))

	ot.AddJob(j)

	clock.Trigger()

	// Sleep for some time to make sure no calls are made
	time.Sleep(1 * time.Second)

	ot.Stop()

	runManager.AssertExpectations(t)
}

func TestOneTime_AddJob_BeforeStart(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)

	clock := cltest.NewTriggerClock(t)

	ot := services.OneTime{
		Clock:      clock,
		Store:      store,
		RunManager: runManager,
	}
	require.NoError(t, ot.Start())

	j := cltest.NewJobWithRunAtInitiator(time.Now())
	j.StartAt = null.TimeFrom(clock.Now().Add(1 * time.Hour))
	require.Nil(t, store.CreateJob(&j))

	ot.AddJob(j)

	clock.Trigger()

	// Sleep for some time to make sure no calls are made
	time.Sleep(1 * time.Second)

	ot.Stop()

	runManager.AssertExpectations(t)
}

func TestExpectedRecurringScheduleJobError(t *testing.T) {
	t.Parallel()

	assert.True(t, services.ExpectedRecurringScheduleJobError(services.RecurringScheduleJobError{}))
	assert.False(t, services.ExpectedRecurringScheduleJobError(errors.New("recurring scheduler job error, but wrong type")))
}
