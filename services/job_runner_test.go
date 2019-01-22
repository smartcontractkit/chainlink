package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobRunner_resumeRunsSinceLastShutdown(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(store)
	defer cleanup()

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	json := fmt.Sprintf(`{"until":"%v"}`, utils.ISO8601UTC(time.Now().Add(time.Second)))
	j.Tasks = []models.TaskSpec{cltest.NewTask("sleep", json)}
	assert.NoError(t, store.SaveJob(&j))

	sleepingRun := j.NewRun(i)
	sleepingRun.Status = models.RunStatusPendingSleep
	sleepingRun.TaskRuns[0].Status = models.RunStatusPendingSleep
	assert.NoError(t, store.SaveJobRun(&sleepingRun))

	inProgressRun := j.NewRun(i)
	inProgressRun.Status = models.RunStatusInProgress
	assert.NoError(t, store.SaveJobRun(&inProgressRun))

	assert.NoError(t, services.ExportedResumeRunsSinceLastShutdown(rm))
	messages := []string{}

	rr, open := <-store.RunChannel.Receive()
	assert.True(t, open)
	messages = append(messages, rr.ID)

	rr, open = <-store.RunChannel.Receive()
	assert.True(t, open)
	messages = append(messages, rr.ID)

	expectedMessages := []string{sleepingRun.ID, inProgressRun.ID}
	assert.ElementsMatch(t, expectedMessages, messages)
}

func TestJobRunner_ChannelForRun_equalityBetweenRuns(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(store)
	defer cleanup()

	job, initr := cltest.NewJobWithWebInitiator()
	run1 := job.NewRun(initr)
	run2 := job.NewRun(initr)

	chan1a := services.ExportedChannelForRun(rm, run1.ID)
	chan2 := services.ExportedChannelForRun(rm, run2.ID)
	chan1b := services.ExportedChannelForRun(rm, run1.ID)

	assert.NotEqual(t, chan1a, chan2)
	assert.Equal(t, chan1a, chan1b)
	assert.NotEqual(t, chan2, chan1b)
}

func TestJobRunner_ChannelForRun_sendAfterClosing(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(s)
	defer cleanup()
	assert.NoError(t, rm.Start())

	j, initr := cltest.NewJobWithWebInitiator()
	assert.NoError(t, s.SaveJob(&j))
	jr := j.NewRun(initr)
	assert.NoError(t, s.SaveJobRun(&jr))

	chan1 := services.ExportedChannelForRun(rm, jr.ID)
	chan1 <- struct{}{}
	cltest.WaitForJobRunToComplete(t, s, jr)

	gomega.NewGomegaWithT(t).Eventually(func() chan<- struct{} {
		return services.ExportedChannelForRun(rm, jr.ID)
	}).Should(gomega.Not(gomega.Equal(chan1))) // eventually deletes the channel

	chan2 := services.ExportedChannelForRun(rm, jr.ID)
	chan2 <- struct{}{} // does not panic
}

func TestJobRunner_ChannelForRun_equalityWithoutClosing(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(s)
	defer cleanup()
	assert.NoError(t, rm.Start())

	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask("nooppend")}
	assert.NoError(t, s.SaveJob(&j))
	jr := j.NewRun(initr)
	assert.NoError(t, s.SaveJobRun(&jr))

	chan1 := services.ExportedChannelForRun(rm, jr.ID)

	chan1 <- struct{}{}
	cltest.WaitForJobRunToPendConfirmations(t, s, jr)

	chan2 := services.ExportedChannelForRun(rm, jr.ID)
	assert.Equal(t, chan1, chan2)
}

func TestJobRunner_Stop(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(s)
	defer cleanup()
	j, initr := cltest.NewJobWithWebInitiator()
	jr := j.NewRun(initr)

	require.NoError(t, rm.Start())

	services.ExportedChannelForRun(rm, jr.ID)
	assert.Equal(t, 1, services.ExportedWorkerCount(rm))

	rm.Stop()

	gomega.NewGomegaWithT(t).Eventually(func() int {
		return services.ExportedWorkerCount(rm)
	}).Should(gomega.Equal(0))
}
