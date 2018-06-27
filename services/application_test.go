// +build !windows

package services_test

import (
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/tevino/abool"
)

func TestServices_ApplicationSignalShutdown(t *testing.T) {
	config, cleanup := cltest.NewConfig()
	defer cleanup()
	app, _ := cltest.NewApplicationWithConfig(config)

	completed := abool.New()
	app.Exiter = func(code int) {
		completed.Set()
	}

	app.Start()
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return completed.IsSet()
	}).Should(gomega.BeTrue())
}

func TestRunManager_Start_ResumeSleepingRuns(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()
	rm := services.NewRunManager(store)

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	json := fmt.Sprintf(`{"until":"%v"}`, utils.ISO8601UTC(time.Now().Add(time.Second)))
	j.Tasks = []models.TaskSpec{cltest.NewTask("sleep", json)}
	assert.NoError(t, store.Save(&j))

	jr := j.NewRun(i)
	jr.Status = models.RunStatusPendingSleep
	assert.NoError(t, store.Save(&jr))

	assert.NoError(t, rm.ResumeSleepingRuns())
	rr, open := <-store.RunQueue
	assert.Equal(t, jr.ID, rr.Input.JobRunID)
	assert.True(t, open)
}

func TestRunManager_WorkerChannelFor_equalityBetweenRuns(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	rm := services.NewRunManager(store)

	job, initr := cltest.NewJobWithWebInitiator()
	run1 := job.NewRun(initr)
	run2 := job.NewRun(initr)

	chan1a := rm.WorkerChannelFor(run1.ID)
	chan2 := rm.WorkerChannelFor(run2.ID)
	chan1b := rm.WorkerChannelFor(run1.ID)

	assert.NotEqual(t, chan1a, chan2)
	assert.Equal(t, chan1a, chan1a)
	assert.NotEqual(t, chan2, chan1b)
}

func TestRunManager_WorkerChannelFor_equalityAfterClosing(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	rm := services.NewRunManager(s)
	assert.NoError(t, rm.Start())

	j, initr := cltest.NewJobWithWebInitiator()
	assert.NoError(t, s.SaveJob(&j))
	jr := j.NewRun(initr)
	assert.NoError(t, s.Save(&jr))

	chan1 := rm.WorkerChannelFor(jr.ID)
	chan2 := rm.WorkerChannelFor(jr.ID)
	assert.Equal(t, chan1, chan2)

	chan1 <- store.RunRequest{}
	cltest.WaitForJobRunToComplete(t, s, jr)

	chan2 = rm.WorkerChannelFor(jr.ID)
	assert.NotEqual(t, chan1, chan2)
}

func TestRunManager_WorkerChannelFor_equalityWithoutClosing(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	rm := services.NewRunManager(s)
	assert.NoError(t, rm.Start())

	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask("nooppend")}
	assert.NoError(t, s.SaveJob(&j))
	jr := j.NewRun(initr)
	assert.NoError(t, s.Save(&jr))

	chan1 := rm.WorkerChannelFor(jr.ID)

	chan1 <- store.RunRequest{}
	cltest.WaitForJobRunToPendConfirmations(t, s, jr)

	chan2 := rm.WorkerChannelFor(jr.ID)
	assert.Equal(t, chan1, chan2)
}

func TestRunManager_Stop(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	rm := services.NewRunManager(s)

	j, initr := cltest.NewJobWithWebInitiator()
	jr := j.NewRun(initr)
	rc := rm.WorkerChannelFor(jr.ID)

	rm.Stop()

	_, open := <-rc
	assert.False(t, open)
}
